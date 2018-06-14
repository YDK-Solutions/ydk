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

#include "entity_util.hpp"
#include "errors.hpp"
#include "json.hpp"
#include "logger.hpp"
#include "types.hpp"

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

vector<string> capabilities;
gNMIClient::PathPrefixValueFlags flag;


static bool check_capabilities_status(Status status);
static json parse_gnmi_set_request_payload(const string & payload);
static string check_if_path_has_value(string element);
static string format_notification_response(string prefix_to_prepend, const std::string& path_to_prepend, const std::string& value);
static void allocate_set_request_path(const string & operation, gnmi::SetRequest & request, vector<string> root_path, json config_payload);

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

static bool check_capabilities_status(Status status)
{
    if (!(status.ok()))
    {
      YLOG_ERROR("Capabilities RPC Failed.");
      throw(YError{"Capabilities RPC Failed"});
      return false;
    }
    return true;
}

static void parse_data_payload_to_paths(const string & payload_filter, vector<string> & path_container)
{
    string path_elem;
    istringstream payload_filter_value(payload_filter);
    while(getline(payload_filter_value, path_elem, '"'))
    {
        if(path_elem.find("{")==string::npos && path_elem.find("}")==string::npos
             && path_elem.find(",")==string::npos  && path_elem != ":")
            path_container.push_back(path_elem);
    }
}

//static void parse_get_request_payload(const string & payload, vector<string> & path_container)
//{
//    auto payload_to_parse = json::parse("{" + payload + "}");
//    string payload_filter = payload_to_parse.value("/rpc/ietf-netconf:get-config/filter"_json_pointer, "null");
//    if (payload_filter == "null") {
//        payload_filter = payload_to_parse.value("/rpc/ietf-netconf:get/filter"_json_pointer, "null");
//        if(payload_filter == "null") {
//            string delim = ":{\"";
//            int filter_len = payload.find_last_of("\"") - (payload.find(delim) + delim.length());
//            auto filter = payload.substr(payload.find(delim) + delim.length(), filter_len);
//            payload_filter = filter;
//        }
//    }
//    parse_data_payload_to_paths(payload, path_container);
//}

static json parse_set_request_payload(const string & payload)
{
    auto payload_to_parse = json::parse("{" + payload + "}");
    json config_payload = json::parse(payload_to_parse.value("/rpc/ietf-netconf:edit-config/config"_json_pointer, "Empty Config"));
    YLOG_DEBUG("JSON Payload: {}", config_payload.dump());
    return config_payload;
}

static json parse_gnmi_set_request_payload(const string & payload)
{
    auto payload_to_parse = json::parse("{" + payload + "}");
    YLOG_DEBUG("JSON Payload: {}", payload_to_parse.dump());
//    json default_payload = json::parse("{\"value\":\"null\"}");
//    json config_payload = payload_to_parse.value("/filter"_json_pointer, default_payload);
    return payload_to_parse;
}

static void allocate_set_request_path(const string & operation, gnmi::SetRequest & request, vector<string> root_path, json config_payload)
{
    gnmi::Path* path = new gnmi::Path;
    gnmi::Update* update;
    gnmi::TypedValue* value = new gnmi::TypedValue;

    if((operation == "delete")||(operation == "gnmi_delete"))
    {
        gnmi::Path* delete_path = request.add_delete_();
        for(size_t i = 0; i < root_path.size(); ++i)
        {
            delete_path->add_element(root_path[i]);
        }
    } else if ((operation == "create")||(operation == "gnmi_create"))
    {
        update = request.add_update();
        for(size_t i = 0; i < root_path.size(); ++i)
        {
            path->add_element(root_path[i]);
        }
        update->set_allocated_path(path);
        for (auto& config_value : config_payload)
        {
            value->set_json_ietf_val(config_value.dump());
            update->set_allocated_val(value);
        }
    }
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

static gnmi::Path* get_prefix(pair<string, string> & prefix_pair)
{
    gnmi::Path* prefix = new gnmi::Path;
    prefix->set_origin(prefix_pair.first);
    auto prefix_elem = prefix->add_elem();
    prefix_elem->set_name(prefix_pair.second);
    return prefix;
}

static gnmi::Path* get_gnmi_path(std::vector<PathElem> & path_container)
{
    gnmi::Path* path = new gnmi::Path;
    for(size_t i = 0; i < path_container.size(); ++i)
    {
        auto path_elem = path->add_elem();
        path_elem->set_name(path_container[i].path);
        for(auto k : path_container[i].keys)
        {
            auto key = path_elem->mutable_key();
            (*key)[k.name] = k.value;
        }
        YLOG_DEBUG("Key size {}", path_elem->key_size());
    }
    return path;
}

static void populate_subscribe_request(gnmi::SubscribeRequest & request,
                                    std::pair<std::string, std::string> & prefix_pair,
                                    std::vector<PathElem> & path_container,
                                    const std::string & list_mode,
                                    long long qos,
                                    int sample_interval,
                                    const std::string & mode)
{
    gnmi::SubscriptionList* sl = new gnmi::SubscriptionList;

    sl->set_mode(get_sublist_mode(list_mode));
    sl->set_allocated_prefix(get_prefix(prefix_pair));
    sl->set_encoding(gnmi::Encoding::PROTO);

    gnmi::QOSMarking* qosmarking = new gnmi::QOSMarking;
    qosmarking->set_marking(qos);
    sl->set_allocated_qos(qosmarking);

    gnmi::Subscription* sub = sl->add_subscription();

    sub->set_allocated_path(get_gnmi_path(path_container));
    sub->set_sample_interval(sample_interval);

    if(mode == "SAMPLE")
    {
        sub->set_mode(gnmi::SubscriptionMode::SAMPLE);
    }
    else if(mode == "ON_CHANGE")
    {
        sub->set_mode(gnmi::SubscriptionMode::ON_CHANGE);
    }
    else
    {
        sub->set_mode(gnmi::SubscriptionMode::TARGET_DEFINED);
    }

    request.set_allocated_subscribe(sl);
}

static string format_notification_response(string prefix_to_prepend, const std::string& path_to_prepend, const std::string& value)
{
    string reply_to_parse;

    // TODO: Update again when payload format from the server(IOS XR) is made consistent with different request paths
    if (flag.path_has_value)
        reply_to_parse.append("\"data\":{" + prefix_to_prepend + path_to_prepend + "[" + value + "]" + "}}");
    else if (flag.prefix_has_value)
        reply_to_parse.append("\"data\":{" + prefix_to_prepend + ":[" + path_to_prepend + value + "]" + "}}");
    else if ((flag.prefix_has_value) && (flag.path_has_value))
        reply_to_parse.append("\"data\":{" + prefix_to_prepend + ":[" + path_to_prepend + "[" + value + "]]" + "}}");
    else
        reply_to_parse.append("\"data\":{" + prefix_to_prepend + path_to_prepend + value + "}");
    size_t pos = 0;
    while ((pos = path_to_prepend.find ("{", ++pos)) != string::npos) {
    	reply_to_parse.append("}");
    }
    return reply_to_parse;
}

static void populate_get_request(gnmi::GetRequest & request, const std::string& payload, bool is_config)
{
    vector<string> path_container;
    parse_data_payload_to_paths(payload, path_container);

    gnmi::Path* path = request.add_path();
    for(size_t i = 0; i < path_container.size(); ++i)
        path->add_element(path_container[i]);
    if(is_config)
        request.set_type(gnmi::GetRequest::CONFIG);
    else
        request.set_type(gnmi::GetRequest::STATE);
    request.set_encoding(gnmi::Encoding::JSON_IETF);
}

static void populate_get_request(gnmi::GetRequest & request, std::pair<std::string, std::string> & prefix,
                                    std::vector<PathElem> & path_container, bool only_config)
{
    request.set_allocated_prefix(get_prefix(prefix));
    auto path = request.add_path();
    for(size_t i = 0; i < path_container.size(); ++i)
    {
        auto path_elem = path->add_elem();
        path_elem->set_name(path_container[i].path);
        for(auto k : path_container[i].keys)
        {
            auto key = path_elem->mutable_key();
            (*key)[k.name] = k.value;
        }
        //YLOG_DEBUG("Key size {}", path_elem->key_size());
    }
    if(only_config)
        request.set_type(gnmi::GetRequest::CONFIG);
    else
        request.set_type(gnmi::GetRequest::STATE);
    request.set_encoding(gnmi::Encoding::JSON_IETF);
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

bool gNMIClient::has_gnmi_version(gnmi::CapabilityResponse* response) 
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

void gNMIClient::parse_capabilities_encodings(gnmi::CapabilityResponse* response)
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

    check_capabilities_status(status);

    if(has_gnmi_version(&response))
    {
        parse_capabilities(&response);
    }
    return capabilities;
}

void gNMIClient::populate_set_request(gnmi::SetRequest & request, const std::string& payload, const std::string& operation)
{
    YLOG_DEBUG("Payload: {}\n Operation: {}", payload, operation);
    json config_payload = parse_gnmi_set_request_payload(payload);
//    if ((operation == "gnmi_create")||(operation == "gnmi_delete")){
//        config_payload = parse_gnmi_set_request_payload(payload);
//    } else {
//        config_payload = parse_set_request_payload(payload);
//    }
    string path_elem;
    vector<string> root_path;
    istringstream payload_config_value(config_payload.dump());
    while(getline(payload_config_value, path_elem, '"'))
    {
        if(path_elem.find("{")==string::npos && path_elem.find("}")==string::npos)
        {
            root_path.push_back(path_elem);
            break;
        }
    }
    allocate_set_request_path(operation, request, root_path, payload);
}

string gNMIClient::get_prefix_from_notification(gnmi::Notification notification)
{
    gnmi::Path prefix = notification.prefix();
    string prefix_to_prepend;

    string origin = prefix.origin();
    if (origin.length() > 0) {
    	prefix_to_prepend = "\"";
    	prefix_to_prepend.append(origin);
    }

    if (prefix.element_size() > 0) {
        for (int i = 0; i < prefix.element_size(); i++)
        {
            if (prefix.element(i) != origin) {
            	if (origin.length() > 0)
            		prefix_to_prepend.append(":");
            	else
            		prefix_to_prepend = "\"";
                prefix_to_prepend.append(prefix.element(i));
                break;
            }
        }
        if (prefix_to_prepend.length() > 0)
            prefix_to_prepend.append("\":");
    }
    else if (prefix.elem_size() > 0) {
        for (int i = 0; i < prefix.elem_size(); i++)
        {
      	    auto path_elem = prefix.elem(i);
            if (path_elem.name() != origin) {
            	if (origin.length() > 0)
            		prefix_to_prepend.append(":");
            	else
            		prefix_to_prepend = "\"";
                prefix_to_prepend.append(path_elem.name());
                break;
            }
        }
        if (prefix_to_prepend.length() > 0)
            prefix_to_prepend.append("\":");
    }
    return prefix_to_prepend;
}

static string check_if_path_has_value(string element)
{
    if(element.find("[")==string::npos)
    {
        flag.path_has_value = false;
        return element;
    } else {
        flag.path_has_value = true;
        return element.substr(0,element.find("["));
    }
}

string gNMIClient::get_path_from_update(gnmi::Update update)
{
    string path_to_prepend;
    string path_element_to_add;

    if(update.has_path())
    {
        gnmi::Path path = update.path();
        int element_size = path.element_size();
        if (element_size > 0) {
            for(int l = 0; l < element_size - 1; ++l)
            {
                path_element_to_add = check_if_path_has_value(path.element(l));
                path_to_prepend.append("\"" + path_element_to_add + "\":{");
            }
            path_element_to_add = check_if_path_has_value(path.element(element_size - 1));
            path_to_prepend.append("\"" + path_element_to_add + "\":");
        }
    }
    else {
        path_to_prepend.clear();
    }
    return path_to_prepend;
}

string gNMIClient::get_value_from_update(gnmi::Update update)
{
    gnmi::TypedValue value;
    string value_for_payload;

    if(update.has_val())
        value = update.val();
        value_for_payload.append(value.json_ietf_val());
    return value_for_payload;
}

string gNMIClient::parse_get_response(gnmi::GetResponse* response)
{
    gnmi::Notification notification;
    gnmi::Update update;
    string value;
    string reply_to_parse;
    string prefix_to_prepend;
    string path_to_prepend;

    int notification_size = response->notification_size();

    for(int i = 0; i < notification_size; ++i)
    {
        notification = response->notification(i);

        if (notification.has_prefix()) {
            prefix_to_prepend = get_prefix_from_notification(notification);
        }
        else {
            prefix_to_prepend.clear();
        }

        if(notification.update_size() != 0) {
            for(int k = 0; k < notification.update_size(); ++k)
            {
                update = notification.update(k);
                path_to_prepend = get_path_from_update(update);
                value.append(get_value_from_update(update));
                reply_to_parse = format_notification_response(prefix_to_prepend, path_to_prepend, value);
            }
        }
    }
    return reply_to_parse;
}

string gNMIClient::parse_set_response(gnmi::SetResponse* response)
{
    gnmi::UpdateResult update_response;
    string reply_to_parse;
    for(int i = 0; i < response->response_size(); ++i)
    {
        update_response = response->response(i);
        reply_to_parse.append(UpdateResult_Operation_Name(update_response.op()) + " Success");
    }
    return reply_to_parse;
}

string gNMIClient::execute_wrapper(const string & payload, const std::string& operation, bool is_config)
{
    if (operation == "read")
    {
        gnmi::GetRequest request;
        gnmi::GetResponse response;

        populate_get_request(request, payload, is_config);
        YLOG_INFO("\n===============Get Request Sent================\n{}\n", request.DebugString().c_str());
        string reply = execute_get_payload(request, &response);
        YLOG_INFO("Get Operation {} Succeeded", operation);
        return reply;
    }
    else if ((operation == "create")||(operation == "update")||(operation == "delete")||(operation == "gnmi_create")||(operation == "gnmi_delete"))
    {   
        gnmi::SetRequest request;
        gnmi::SetResponse response;

        populate_set_request(request, payload, operation);
        YLOG_INFO("\n===============Set Request Sent================\n{}\n", request.DebugString().c_str());
        string reply = execute_set_payload(request, &response);
        YLOG_INFO("Set Operation {} Succeeded", operation);
        return reply;
    }
    return "";
}

string gNMIClient::execute_get_payload(const GetRequest& request, GetResponse* response)
{
    grpc::ClientContext context;
    grpc::Status status;
    context.AddMetadata("username", username);
    context.AddMetadata("password", password);
    status = stub_->Get(&context, request, response);
    YLOG_INFO("\n=============Get Response Received=============\n{}\n", response->DebugString().c_str());
    if (!(status.ok())) 
    {
        YLOG_ERROR("Get RPC Status Not OK:\n{}", status.error_message());
        throw(YError{status.error_message()});
    }
    else 
    {
        YLOG_DEBUG("Get RPC OK");
        string reply = parse_get_response(response);
        return reply;
    }
}

string gNMIClient::execute_set_payload(const SetRequest& request, SetResponse* response)
{
    grpc::ClientContext context;
    grpc::Status status;
    context.AddMetadata("username", username);
    context.AddMetadata("password", password);
    status = stub_->Set(&context, request, response);
    YLOG_INFO("\n=============Set Response Received=============\n{}\n", response->DebugString().c_str());
    if (!(status.ok())) 
    {
        YLOG_ERROR("Set RPC Status not OK");
        throw(YError{status.error_message()});
    } 

    YLOG_DEBUG("Set RPC OK");
    return parse_set_response(response);
}

string gNMIClient::execute_get_operation(std::pair<std::string, std::string> & prefix,
                                    std::vector<PathElem> & path_container, bool only_config)
{
    gnmi::GetRequest request;
    gnmi::GetResponse response;

    populate_get_request(request, prefix, path_container, only_config);
    YLOG_INFO("\n===============Get Request Sent================\n{}\n", request.DebugString().c_str());
    string reply = execute_get_payload(request, &response);
    YLOG_INFO("Get Operation Succeeded ");
    return reply;
}

void gNMIClient::execute_subscribe_operation(std::pair<std::string, std::string> & prefix,
                                    std::vector<PathElem> & path_container,
                                    const std::string & list_mode,
                                    long long qos,
                                    int sample_interval,
                                    const std::string & mode,
                                    std::function<void(const std::string &)> func)
{
    grpc::ClientContext context;
    gnmi::SubscribeResponse response;
    auto request = make_shared<gnmi::SubscribeRequest>();

    populate_subscribe_request(*request, prefix, path_container, list_mode, qos, sample_interval, mode);

    context.AddMetadata("username", username);
    context.AddMetadata("password", password);

    YLOG_INFO("\n===============Subscribing================");
    auto a = stub_->Subscribe(&context);
    std::shared_ptr<grpc::ClientReaderWriter<gnmi::SubscribeRequest, ::gnmi::SubscribeResponse>> client_reader_writer = move(a);

    YLOG_INFO("\n===============Sending SubscribeRequest================\n{}\n", request->DebugString().c_str());
    client_reader_writer->Write(*request);
//    client_reader_writer->WritesDone();
    YLOG_DEBUG("Done sending request");

    YLOG_INFO("Subscribe Operation Succeeded");
    YLOG_DEBUG("\n=============Receiving Subscribe Response=============");
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
                    YLOG_INFO("\n===============Sending SubscribeRequest================\n{}\n", req.DebugString().c_str());
                    client_reader_writer->Write(req);
                }
            }
        });
        writer.detach();
    }

    while (client_reader_writer->Read(&response))
    {
        if(func!=nullptr)
        {
            string s;
            google::protobuf::TextFormat::PrintToString(response, &s);
            func(s);
        }
    }
}

}
