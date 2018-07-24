/// YANG Development Kit
// Copyright 2016 Cisco Systems. All rights reserved
//
////////////////////////////////////////////////////////////////
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
//////////////////////////////////////////////////////////////////

#include <thread>
#include <google/protobuf/text_format.h>

#include <ydk/errors.hpp>
#include <ydk/json.hpp>
#include <ydk/logger.hpp>

#include "gnmi_util.hpp"
#include "gnmi_client.hpp"

using namespace std;
using namespace grpc;

using json = nlohmann::json;
using grpc::SslCredentialsOptions;


namespace ydk
{
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Utility functions
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

gNMIClient::PathPrefixValueFlags flag;

static void check_status(grpc::Status status, string message)
{
    if (!status.ok()) {
        ostringstream s;
        s << message << ":\n" << status.error_message();
        string er_msg = s.str();
        YLOG_ERROR(er_msg.c_str());
        throw(YServiceProviderError{er_msg});
    }
}

static std::shared_ptr<Channel> connect_to_server()
{
    grpc::SslCredentialsOptions ssl_opts;
    grpc::ChannelArguments      args;

    string server_cert;
    string server_key;

    ifstream k("ems.key");
    ifstream p("ems.pem");

    server_cert.assign((istreambuf_iterator<char>(p)),(istreambuf_iterator<char>()));
    server_key.assign((istreambuf_iterator<char>(k)),(istreambuf_iterator<char>()));

    ssl_opts.pem_root_certs = server_cert;
    ssl_opts.pem_private_key = server_key;

    args.SetSslTargetNameOverride("2001:420:1101:1::b");
    //args.SetSslTargetNameOverride("ems.cisco.com");

    YLOG_DEBUG("In gnmi_connect server cert: {}", server_cert);

    auto channel_creds = grpc::SslCredentials(grpc::SslCredentialsOptions(ssl_opts));
    return grpc::CreateCustomChannel("", channel_creds, args);
}

static gnmi::SubscriptionList_Mode get_sublist_mode(const string & list_mode)
{
    if(list_mode == "STREAM")
    {
        return gnmi::SubscriptionList_Mode_STREAM;
    }
    else if(list_mode == "POLL")
    {
        return gnmi::SubscriptionList_Mode_POLL;
    }
    return gnmi::SubscriptionList_Mode_ONCE;
}

static void populate_subscribe_request(gnmi::SubscriptionList* sl, gNMISubscription subscription)
{
    gnmi::Subscription* sub = sl->add_subscription();
    if (subscription.path)
        sub->set_allocated_path(subscription.path);
    sub->set_sample_interval(subscription.sample_interval);
    sub->set_suppress_redundant(subscription.suppress_redundant);
    sub->set_heartbeat_interval(subscription.heartbeat_interval);

    if(subscription.subscription_mode == "SAMPLE")
    {
        sub->set_mode(gnmi::SubscriptionMode::SAMPLE);
    }
    else if(subscription.subscription_mode == "ON_CHANGE")
    {
        sub->set_mode(gnmi::SubscriptionMode::ON_CHANGE);
    }
    else {
        sub->set_mode(gnmi::SubscriptionMode::TARGET_DEFINED);
    }
}

static string format_notification_response(const std::string& path_to_prepend, const std::string& value)
{
    string reply_val;

    // TODO: Update again when payload format from the server(IOS XR) is made consistent with different request paths
    if (flag.path_has_value)
        reply_val = "{" + path_to_prepend + "{[" + value + "]}}";
//    else if (flag.prefix_has_value)
//        reply_to_parse.append("\"data\":{" + prefix_to_prepend + ":[" + path_to_prepend + value + "]" + "}}");
//    else if ((flag.prefix_has_value) && (flag.path_has_value))
//        reply_to_parse.append("\"data\":{" + prefix_to_prepend + ":[" + path_to_prepend + "[" + value + "]]" + "}}");
    else {
        string val = (value == "{\"value\":\"null\"}") ? "{}" : value;
        if (path_to_prepend.length()==0 && val.find("{")==0)
            reply_val = val;
        else
            reply_val = "{" + path_to_prepend + val + "}";
    }
    size_t pos = 0;
    while ((pos = path_to_prepend.find ("{", ++pos)) != string::npos) {
        reply_val.append("}");
    }
    return reply_val;
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// gNMIClient
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

gNMIClient::gNMIClient(shared_ptr<Channel> channel, const std::string & username, const  std::string & password)
 : stub_(gNMI::NewStub(channel)), username(username), password(password), is_secure(true)
{
}

gNMIClient::gNMIClient(shared_ptr<Channel> channel)
 : stub_(gNMI::NewStub(channel)), username(""), password(""), is_secure(false)
{
}

gNMIClient::~gNMIClient()
{
}

int gNMIClient::connect()
{
    // Authenticate Server at Client
	if (is_secure) {
        connect_to_server();
	}
    return EXIT_SUCCESS;
}

static bool has_gnmi_version(gnmi::CapabilityResponse* response)
{
    if (!(response->gnmi_version().size() > 0)) 
    {
        YLOG_ERROR("Capabilities Received Without gNMI Version");
        throw(YError{"Capabilities Received Without gNMI Version"});
        return false;
    }
    return true;
}

void gNMIClient::parse_capabilities_modeldata(gnmi::CapabilityResponse* response)
{
    gnmi::ModelData modeldata;
    string cap, name, organization, version; 

    YLOG_DEBUG("====Capabilities Response Received====");
    for (size_t i = 0, n = response->supported_models_size(); i < n; i++)
    {
        ostringstream os;
        cap.clear();
        modeldata = response->supported_models(i);
        name = modeldata.name();
        organization = modeldata.organization();
        version = modeldata.version();
        if(!(modeldata.name().empty()))
        {
            os << name << ", ";
            cap.append("?module=" + name);
        }
        if(!(modeldata.organization().empty()))
            os << organization << ", ";
        if(!(modeldata.version().empty()))
        {
            string revision = version.c_str();
            revision = revision.substr(revision.find(":") + 1,string::npos);
            os << revision;
        }
        YLOG_DEBUG("{}", os.str().c_str());
        capabilities.push_back(cap);
    }
}

static void parse_capabilities_encodings(gnmi::CapabilityResponse* response)
{
    gnmi::Encoding encoding;
    string encoding_value;

    for (int i = 0, n = response->supported_encodings_size(); i < n; i++)
    {
        encoding = response->supported_encodings(i);
        encoding_value = gnmi::Encoding_Name(encoding);
        YLOG_DEBUG("Encoding {}", encoding_value.c_str());
    }
}

void gNMIClient::parse_capabilities(gnmi::CapabilityResponse* response)
{
    YLOG_DEBUG("Capabilities Received:");
    YLOG_DEBUG("");
    YLOG_DEBUG("==============gNMI Version==============");
    YLOG_DEBUG("gNMI Version: {}", response->gnmi_version().c_str());

    YLOG_DEBUG("============Supported Models============");
    parse_capabilities_modeldata(response);

    YLOG_DEBUG("===========Supported Encodings===========");
    parse_capabilities_encodings(response);

    YLOG_DEBUG("=========================================");
    YLOG_DEBUG("");
}

vector<string> gNMIClient::get_capabilities()
{
    CapabilityRequest request;
    CapabilityResponse response;
    grpc::ClientContext context;
    context.AddMetadata("username", username);
    context.AddMetadata("password", password);
    grpc::Status status = stub_->Capabilities(&context, request, &response);
    check_status(status, "CapabilityRequest failed with error");

    if(has_gnmi_version(&response))
    {
        parse_capabilities(&response);
    }
    return capabilities;
}

gNMICapabilityResponse gNMIClient::execute_get_capabilities()
{
    gNMICapabilityResponse reply{};

    CapabilityRequest request;
    CapabilityResponse response;
    grpc::ClientContext context;
    context.AddMetadata("username", username);
    context.AddMetadata("password", password);
    grpc::Status status = stub_->Capabilities(&context, request, &response);
    check_status(status, "CapabilityRequest failed with error");

    reply.gnmi_version = response.gnmi_version();
    for (size_t i = 0, n = response.supported_models_size(); i < n; i++)
    {
        auto model = response.supported_models(i);

        gNMIModelData model_data{};
        model_data.name = model.name();
        model_data.organization = model.organization();
        model_data.version = model.version();

        reply.supported_models.push_back(model_data);
    }
    for (int i = 0, n = response.supported_encodings_size(); i < n; i++)
    {
        gnmi::Encoding encoding = response.supported_encodings(i);
        string encoding_value = gnmi::Encoding_Name(encoding);
        reply.supported_encodings.push_back(encoding_value);
    }
    return reply;
}

static string get_path_from_update(gnmi::Update update)
{
    string path_to_prepend;
    string path_element_to_add;

    if(update.has_path())
    {
        gnmi::Path path = update.path();
        int elem_size = path.elem_size();
        if (elem_size > 0) {
            string origin = path.origin();
            if (origin.length() > 0) {
                path_to_prepend.append("\"" + origin + ":");
            }
            int l;
            for (l = 0; l < elem_size - 1; l++)
            {
                path_element_to_add = path.elem(l).name();	//check_if_path_has_key_values(path.elem(l));
                if (l>0 || path_to_prepend.length()==0)
                    path_to_prepend.append("\"");
                path_to_prepend.append(path_element_to_add + "\":{");
            }
            path_element_to_add = path.elem(l).name();	//check_if_path_has_key_values(path.elem(elem_size-1));
            if (elem_size>1 || path_to_prepend.length()==0)
                path_to_prepend.append("\"");
            path_to_prepend.append(path_element_to_add + "\":");
        }
    }
    else {
        path_to_prepend.clear();
    }
    return path_to_prepend;
}

static string get_value_from_update(gnmi::Update update)
{
    gnmi::TypedValue value;
    string value_for_payload;

    if(update.has_val())
        value = update.val();
        value_for_payload.append(value.json_ietf_val());
    return value_for_payload;
}

static vector<string> parse_get_response(gnmi::GetResponse* response)
{
    vector<string> response_list{};

    for(int i = 0; i < response->notification_size(); ++i)
    {
        gnmi::Notification notification = response->notification(i);
        string path_to_prepend;
        string reply_to_parse;
        string value;
        if(notification.update_size() != 0) {
            for(int k = 0; k < notification.update_size(); ++k)
            {
                gnmi::Update update = notification.update(k);
                path_to_prepend = get_path_from_update(update);
                value.append(get_value_from_update(update));
                reply_to_parse = format_notification_response(path_to_prepend, value);
            }
        }
        response_list.push_back(reply_to_parse);
    }
    return response_list;
}

vector<string>
gNMIClient::execute_get_operation(const std::vector<gNMIRequest> get_request_list, const std::string& operation)
{
    gnmi::GetRequest gnmi_get_request;
    gnmi::GetResponse gnmi_get_response;

    for (auto ydk_request : get_request_list) {
        // Populate gnmi::GetRequest
        gnmi::Path* path = gnmi_get_request.add_path();
        if (ydk_request.path != nullptr)
            path->CopyFrom(*ydk_request.path);
    }

    if (operation == "ALL")
        gnmi_get_request.set_type(gnmi::GetRequest::ALL);
    else if (operation == "OPERATIONAL")
        gnmi_get_request.set_type(gnmi::GetRequest::OPERATIONAL);
    else if (operation == "STATE")
        gnmi_get_request.set_type(gnmi::GetRequest::STATE);
    else
        gnmi_get_request.set_type(gnmi::GetRequest::CONFIG);
    gnmi_get_request.set_encoding(gnmi::Encoding::JSON_IETF);

    YLOG_INFO("\n=============== Get Request Sent ================\n{}\n", gnmi_get_request.DebugString());
    auto reply = execute_get_payload(gnmi_get_request, &gnmi_get_response);
    YLOG_INFO("Get Operation Succeeded");
    return reply;
}

bool
gNMIClient::execute_set_operation(const std::vector<ydk::gNMIRequest> set_request_list)
{
    gnmi::SetRequest gnmi_set_request;
    gnmi::SetResponse gnmi_set_response;

    // Populate gnmi::SetRequest
    // delete
    for (auto request : set_request_list) {
        if (request.operation == "delete") {
            gnmi::Path* path = gnmi_set_request.add_delete_();
            if (request.path != nullptr)
                path->CopyFrom(*request.path);
        }
    }
    // replace
    for (auto request : set_request_list) {
        if (request.operation == "replace") {
            ::gnmi::Update* update = gnmi_set_request.add_replace();
            if (request.path != nullptr) {
                update->set_allocated_path(request.path);
                if (!request.payload.empty()) {
                    ::gnmi::TypedValue* value = new ::gnmi::TypedValue;
                    value->set_json_ietf_val(request.payload);
                    update->set_allocated_val(value);
                }
            }
        }
    }
    // update
    for (auto request : set_request_list) {
        if (request.operation == "update") {
            ::gnmi::Update* update = gnmi_set_request.add_update();
            if (request.path != nullptr) {
                update->set_allocated_path(request.path);
                if (!request.payload.empty()) {
                    ::gnmi::TypedValue* value = new ::gnmi::TypedValue;
                    value->set_json_ietf_val(request.payload);
                    update->set_allocated_val(value);
                }
            }
        }
    }

    YLOG_INFO("\n=============== Set Request Sent ================\n{}\n", gnmi_set_request.DebugString());
    auto reply = execute_set_payload(gnmi_set_request, &gnmi_set_response);
    YLOG_INFO("Set Operation Succeeded");
    return reply;
}

vector<string>
gNMIClient::execute_get_payload(const GetRequest& request, GetResponse* response)
{
    grpc::ClientContext context;
    context.AddMetadata("username", username);
    context.AddMetadata("password", password);
    grpc::Status status = stub_->Get(&context, request, response);
    check_status(status, "GetRequest failed with error");
    YLOG_INFO("\n============= Get Response Received =============\n{}\n", response->DebugString().c_str());
    YLOG_DEBUG("Get RPC succeeded");
    return parse_get_response(response);
}

bool
gNMIClient::execute_set_payload(const SetRequest& request, SetResponse* response)
{
    grpc::ClientContext context;
    context.AddMetadata("username", username);
    context.AddMetadata("password", password);
    grpc::Status status = stub_->Set(&context, request, response);
    check_status(status, "SetRequest failed with error");
    YLOG_INFO("\n============= Set Response Received =============\n{}\n", response->DebugString().c_str());
    return true;
}

void
gNMIClient::execute_subscribe_operation(std::vector<gNMISubscription> subscription_list,
		                                uint32 qos, const std::string & list_mode,
                                        std::function<void(const std::string &)> func)
{
    grpc::ClientContext context;
    gnmi::SubscribeResponse response;

    auto request = make_shared<gnmi::SubscribeRequest>();

    gnmi::SubscriptionList* sl = new gnmi::SubscriptionList;
    sl->set_mode(get_sublist_mode(list_mode));
    sl->set_encoding(gnmi::Encoding::JSON_IETF);

    gnmi::QOSMarking* qosmarking = new gnmi::QOSMarking;
    qosmarking->set_marking(qos);
    sl->set_allocated_qos(qosmarking);

    for (auto subscription : subscription_list) {
        populate_subscribe_request(sl, subscription);
    }
    request->set_allocated_subscribe(sl);

    YLOG_INFO("\n=============== Sending SubscribeRequest ================\n{}\n", request->DebugString().c_str());

    context.AddMetadata("username", username);
    context.AddMetadata("password", password);

    auto a = stub_->Subscribe(&context);
    std::shared_ptr<grpc::ClientReaderWriter<gnmi::SubscribeRequest, ::gnmi::SubscribeResponse>> client_reader_writer = move(a);

    client_reader_writer->Write(*request);

    YLOG_INFO("Subscribe Operation Succeeded");
    YLOG_INFO("Invoking callback function to receive the subscription data");

    if(list_mode == "POLL")
    {
        std::thread writer([client_reader_writer]() {
            string i;
            cout<<"Start POLL request monitor....."<<endl;
            cout<<"Enter 'poll' to poll for update: ";
            while(cin >> i)
            {
                if(i == "poll")
                {
                    gnmi::SubscribeRequest req;
                    gnmi::Poll* p = new gnmi::Poll;
                    req.set_allocated_poll(p);
                    YLOG_INFO("\n=============== Sending SubscribeRequest ================\n{}\n", req.DebugString().c_str());
                    client_reader_writer->Write(req);
                }
            }
        });
        writer.detach();
    }

    while (client_reader_writer->Read(&response))
    {
        YLOG_INFO("\n=============== Received SubscribeResponse ================\n{}\n", response.DebugString());

        if (func != nullptr) {
            string s;
            google::protobuf::TextFormat::PrintToString(response, &s);
            func(s);
        }
        if (list_mode == "ONCE") {
            break;
        }
    }
    client_reader_writer->Finish();
}

}
