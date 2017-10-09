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

#include "types.hpp"
#include "json.hpp"
#include "errors.hpp"
#include "logger.hpp"

#include "gnmi_client.hpp"

using namespace std;
using namespace ydk;

using json = nlohmann::json;
using grpc::SslCredentialsOptions;

vector<string> capabilities;
gNMIClient::PathPrefixValueFlags flag;

gNMIClient::gNMIClient(shared_ptr<Channel> channel) : stub_(gNMI::NewStub(channel))
{
}

static bool check_capabilities_status(Status status);
static string parse_get_request_payload(string payload);
static json parse_set_request_payload(string payload);
static json parse_gnmi_set_request_payload(string payload);
static vector<string> get_path_from_payload_filter(string payload_filter, vector<string> container);
static string check_if_path_has_value(string element);
static string format_notification_response(string prefix_to_prepend, string path_to_prepend, string value);
static ::gnmi::SetRequest allocate_set_request_path(string operation, ::gnmi::SetRequest request, vector<string> root_path, json config_payload);


int gNMIClient::connect(string address, bool is_secure) 
{
    if(is_secure) 
    {
        // Authenticate Server at Client
        string server_cert;
        ifstream rf("ems.pem");
    
        server_cert.assign((istreambuf_iterator<char>(rf)),(istreambuf_iterator<char>()));
    
        grpc::SslCredentialsOptions ssl_opts;
        grpc::ChannelArguments      args;
    
        ssl_opts.pem_root_certs = server_cert;
        args.SetSslTargetNameOverride("ems.cisco.com");

        YLOG_DEBUG("In gnmi_connect server cert: {}", server_cert);
    
        // TODO: Authenticate Client at Server
        /* string client_key, client_cert;
        ifstream kf("client.key");
        ifstream cf("client.pem");
        client_key.assign((istreambuf_iterator<char>(kf)),(istreambuf_iterator<char>()));
        client_cert.assign((istreambuf_iterator<char>(cf)),(istreambuf_iterator<char>()));
        ssl_opts = {server_cert, client_key, client_cert};*/
    
        auto channel_creds = grpc::SslCredentials(grpc::SslCredentialsOptions(ssl_opts));
        grpc::CreateCustomChannel(address, channel_creds, args);

    } else {
        grpc::CreateChannel(address, grpc::InsecureChannelCredentials());
    }
    get_capabilities();
    return EXIT_SUCCESS;
}

static bool check_capabilities_status(Status status) 
{
    if (!(status.ok())) 
    {
      YLOG_ERROR("Capabilities RPC Failed.");
      throw(YCPPError{"Capabilities RPC Failed"});
      return false;
    }
    return true;
}

bool gNMIClient::has_gnmi_version(::gnmi::CapabilityResponse* response) 
{
    if (!(response->gnmi_version().size() > 0)) 
    {
        YLOG_ERROR("Capabilities Received Without gNMI Version");
        throw(YCPPError{"Capabilities Received Without gNMI Version"});
        return false;
    }
    return true; 
}

void gNMIClient::parse_capabilities_modeldata(::gnmi::CapabilityResponse* response)
{
    ::gnmi::ModelData modeldata;
    string cap, name, organization, version; 

    YLOG_DEBUG("====Capabilities Response Received====");
    for (int i = 0, n = response->supported_models_size(); i < n; i++) 
    {
        cap.clear();
        modeldata = response->supported_models(i);
        name = modeldata.name();
        organization = modeldata.organization();
        version = modeldata.version();
        if(!(modeldata.name().empty()))
        {
            YLOG_DEBUG("Name: {}", name.c_str());
            cap.append("?module=" + name);
        }
        if(!(modeldata.organization().empty())) 
            YLOG_DEBUG("Organization: {}", organization.c_str());
        if(!(modeldata.version().empty()))
        {
            string revision = version.c_str();
            revision = revision.substr(revision.find(":") + 1,string::npos);
            YLOG_DEBUG("Revision: {}", revision);
            cap.append("&revision=" + revision);
            YLOG_DEBUG("              ------------       ");
            YLOG_DEBUG("");
        }
        capabilities.push_back(cap);
    }   
}

void gNMIClient::parse_capabilities_encodings(::gnmi::CapabilityResponse* response)
{
    ::gnmi::Encoding encoding;
    string encoding_value;

    for (int i = 0, n = response->supported_encodings_size(); i < n; i++) 
    {
        encoding = response->supported_encodings(i);
        encoding_value = ::gnmi::Encoding_Name(encoding);
        YLOG_DEBUG("Encoding {}", encoding_value.c_str());
    }
}

bool gNMIClient::parse_capabilities(::gnmi::CapabilityResponse* response)
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
    grpc::Status status = stub_->Capabilities(&context, request, &response);

    check_capabilities_status(status);
    
    if(has_gnmi_version(&response)) 
    {
        parse_capabilities(&response);
    }
    return capabilities;
}

static vector<string> parse_get_request_payload(string payload, vector<string> path_container) 
{
    auto payload_to_parse = json::parse("{" + payload + "}");
    string payload_filter = payload_to_parse.value("/rpc/ietf-netconf:get-config/filter"_json_pointer, "null");
    if (payload_filter == "null") {
        payload_filter = payload_to_parse.value("/rpc/ietf-netconf:get/filter"_json_pointer, "null");
        if(payload_filter == "null") {
            string delim = ":{\"";
            int filter_len = payload.find_last_of("\"") - (payload.find(delim) + delim.length());
            auto filter = payload.substr(payload.find(delim) + delim.length(), filter_len);
            payload_filter = filter;
        }
    } 

    string path_elem;
    istringstream payload_filter_value(payload_filter);
    while(getline(payload_filter_value, path_elem, '"')) 
    {
        if(path_elem.find("{")==string::npos && path_elem.find("}")==string::npos) 
            path_container.push_back(path_elem);
    }
    return path_container;
}

::gnmi::GetRequest gNMIClient::populate_get_request(::gnmi::GetRequest request, string payload) 
{   
    vector<string> container;   
    vector<string> path_container = parse_get_request_payload(payload, container);
    
    ::gnmi::Path* path = request.add_path();  
    for(int i = 0; i < path_container.size(); ++i) 
        path->add_element(path_container[i]);
    request.set_type(::gnmi::GetRequest::CONFIG);
    request.set_encoding(::gnmi::Encoding::JSON_IETF);
    return request;
}

json parse_set_request_payload(string payload) 
{
    auto payload_to_parse = json::parse("{" + payload + "}");
    json config_payload = json::parse(payload_to_parse.value("/rpc/ietf-netconf:edit-config/config"_json_pointer, "Empty Config"));
    YLOG_DEBUG("JSON Payload: {}", config_payload.dump());
    return config_payload;
}

json parse_gnmi_set_request_payload(string payload) 
{
    auto payload_to_parse = json::parse("{" + payload + "}");
    json default_payload = json::parse("{\"value\":\"null\"}");
    json config_payload = payload_to_parse.value("/filter"_json_pointer, default_payload);
    return config_payload;
}

::gnmi::SetRequest allocate_set_request_path(string operation, ::gnmi::SetRequest request, vector<string> root_path, json config_payload)
{
    ::gnmi::Path* path = new ::gnmi::Path;
    ::gnmi::Update* update;
    ::gnmi::TypedValue* value = new ::gnmi::TypedValue;

    if((operation == "delete")||(operation == "gnmi_delete"))
    {
        ::gnmi::Path* delete_path = request.add_delete_(); 
        for(int i = 0; i < root_path.size(); ++i)
        {
            delete_path->add_element(root_path[i]);  
        }   
    } else if ((operation == "create")||(operation == "gnmi_create"))
    {
        update = request.add_update();
        for(int i = 0; i < root_path.size(); ++i)
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
    return request;  
}

::gnmi::SetRequest gNMIClient::populate_set_request(::gnmi::SetRequest request, string payload, string operation) 
{   
    YLOG_DEBUG("Payload {}\n Operation {}", payload, operation);
    json config_payload;
    if ((operation == "gnmi_create")||(operation == "gnmi_delete")){
        config_payload = parse_gnmi_set_request_payload(payload);
    } else {
        config_payload = parse_set_request_payload(payload);
    }
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
    return allocate_set_request_path(operation, request, root_path, config_payload);
}

string gNMIClient::get_prefix_from_notification(::gnmi::Notification notification)
{
    ::gnmi::Path prefix = notification.prefix();
    string element;
    string prefix_to_prepend;

    int element_size = prefix.element_size();
    for(int j = 0; j < element_size; ++j)
    {
        element.append("\"" + prefix.element(j) + "\" ");
        prefix_to_prepend.append(prefix.element(j));
        prefix_to_prepend.append(":");
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

string gNMIClient::get_path_from_update(::gnmi::Update update)
{
    string path_to_prepend;
    string path_element_to_add;

    if(update.has_path()) 
    {
        ::gnmi::Path path = update.path();
        int element_size = path.element_size();

        for(int l = 0; l < element_size - 1; ++l) 
        {
            path_element_to_add = check_if_path_has_value(path.element(l));
            path_to_prepend.append("\"" + path_element_to_add + "\":{");
        }
        path_element_to_add = check_if_path_has_value(path.element(element_size - 1));
        path_to_prepend.append("\"" + path_element_to_add + "\":");
    } else 
        path_to_prepend.clear();
    return path_to_prepend;
}

string gNMIClient::get_value_from_update(::gnmi::Update update)
{
    ::gnmi::TypedValue value;
    string value_for_payload;

    if(update.has_val()) 
        value = update.val();
        value_for_payload.append(value.json_ietf_val());
    return value_for_payload;
}

static string format_notification_response(string prefix_to_prepend, string path_to_prepend, string value)
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
    return reply_to_parse;
}

string gNMIClient::parse_get_response(::gnmi::GetResponse* response) 
{
    ::gnmi::Notification notification;
    ::gnmi::Update update;
    string value;
    string reply_to_parse;
    string prefix_to_prepend;
    string path_to_prepend;

    int notification_size = response->notification_size();

    for(int i = 0; i < notification_size; ++i) 
    {
        notification = response->notification(i);

        if(notification.has_prefix()) prefix_to_prepend = get_prefix_from_notification(notification);
        else prefix_to_prepend.clear();
        
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

string gNMIClient::parse_set_response(::gnmi::SetResponse* response) 
{
    ::gnmi::UpdateResult update_response;
    string reply_to_parse;
    for(int i = 0; i < response->response_size(); ++i) 
    {
        update_response = response->response(i);
        reply_to_parse.append(UpdateResult_Operation_Name(update_response.op()) + " Success");
    }
    return reply_to_parse;
}

string gNMIClient::execute_wrapper(const string & payload, string operation)
{
    if (operation == "read")
    {
        ::gnmi::GetRequest request;
        ::gnmi::GetResponse response;

        request = populate_get_request(request, payload);
        YLOG_DEBUG("\n===============Get Request Sent================\n{}\n", request.DebugString().c_str());
        string reply = execute_get_payload(request, &response);
        YLOG_DEBUG("Get Operation {} Succeeded", operation);
        return reply;
    } else if ((operation == "create")||(operation == "update")||(operation == "delete")||(operation == "gnmi_create")||(operation == "gnmi_delete"))
    {   
        ::gnmi::SetRequest request;
        ::gnmi::SetResponse response;

        request = populate_set_request(request, payload, operation);
        YLOG_DEBUG("\n===============Set Request Sent================\n{}\n", request.DebugString().c_str());
        string reply = execute_set_payload(request, &response);
        YLOG_DEBUG("Set Operation {} Succeeded", operation);
        return reply;

    }
}

string gNMIClient::execute_get_payload(const GetRequest& request, GetResponse* response)
{
    grpc::ClientContext context;
    grpc::Status status;
    status = stub_->Get(&context, request, response);
    YLOG_DEBUG("\n=============Get Response Received=============\n{}\n", response->DebugString().c_str());
    if (!(status.ok())) 
    {
        YLOG_ERROR("Get RPC Status Not OK");
        throw(YCPPError{status.error_message()});
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
    status = stub_->Set(&context, request, response);
    YLOG_DEBUG("\n=============Set Response Received=============\n{}\n", response->DebugString().c_str());
    if (!(status.ok())) 
    {
        YLOG_ERROR("Set RPC Status not OK");
        throw(YCPPError{status.error_message()});
    } 
    else 
    {
        YLOG_DEBUG("Set RPC OK");
        string reply = parse_set_response(response); 
        return reply;
    }
}

gNMIClient::~gNMIClient() {}
