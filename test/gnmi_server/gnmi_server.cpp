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

#include <ctime>
#include <fstream>
#include <grpc++/grpc++.h>
#include <iostream>
#include <memory>
#include <string>

#include "json.hpp"
#include "gnmi.grpc.pb.h"

using gnmi::CapabilityRequest;
using gnmi::CapabilityResponse;
using gnmi::Encoding;
using gnmi::GetRequest;
using gnmi::GetResponse;
using gnmi::SetRequest;
using gnmi::SetResponse;
using gnmi::gNMI;
using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::SslServerCredentialsOptions;
using grpc::SslServerCredentials;
using grpc::Status;

using json = nlohmann::json;

std::string cap_array[] =
    {"ietf-aug-base-1?module=ietf-aug-base-1&revision=2016-07-01",
     "ietf-aug-base-2?module=ietf-aug-base-2&revision=2016-07-01",
     "main?module=main&revision=2015-11-17",
     "main-aug1?module=main-aug1&revision=2015-11-17",
     "oc-pattern?module=oc-pattern&revision=2015-11-17",
     "yaug-five?module=ydktest-aug-ietf-5&revision=2017-07-26",
     "yaug-four?module=ydktest-aug-ietf-4&revision=2016-06-27",
     "yaug-one?module=ydktest-aug-ietf-1&revision=2016-06-17",
     "yaug-two?module=ydktest-aug-ietf-2&revision=2016-06-22",
     "ydk-filter?module=ydktest-filterread&revision=2015-11-17",
     "ydktest-sanity?module=ydktest-sanity&revision=2015-11-17&features=ipv6-privacy-autoconf,ipv4-non-contiguous-netmasks",
     "ydktest-sanity-augm?module=ydktest-sanity-augm&revision=2015-11-17",
     "ydktest-sanity-types?module=ydktest-sanity-types&revision=2016-04-11",
     "ydktest-types?module=ydktest-types&revision=2016-05-23&features=crypto",
     "bgp?module=openconfig-bgp&revision=2016-06-21",
     "bgp-policy?module=openconfig-bgp-policy&revision=2016-06-21",
     "bgp-types?module=openconfig-bgp-types&revision=2016-06-21",
     "interfaces?module=openconfig-interfaces&revision=2016-05-26",
     "interfaces/ethernet?module=openconfig-if-ethernet&revision=2016-05-26",
     "openconfig-ext?module=openconfig-extensions&revision=2015-10-09",
     "openconfig-types?module=openconfig-types&revision=2016-05-31",
     "platform?module=openconfig-platform&revision=2016-06-06&deviations=cisco-xr-openconfig-platform-deviations",
     "platform-types?module=openconfig-platform-types&revision=2016-06-06",
     "platform/transceiver?module=openconfig-platform-transceiver&revision=2016-05-24",
     "policy-types?module=openconfig-policy-types&revision=2016-05-12",
     "routing-policy?module=openconfig-routing-policy&revision=2016-05-12",
     "terminal-device?module=openconfig-terminal-device&revision=2016-06-17",
     "transport-types?module=openconfig-transport-types&revision=2016-06-17"
    };

// Logic and data behind the server's behavior.
class gNMIImpl final : public gNMI::Service 
{    
    int set_counter = 0;
    int delete_counter = 0;
    bool is_secure = true;

    Status Capabilities(ServerContext* context, const CapabilityRequest* request, CapabilityResponse* response) override 
    {
        response->set_gnmi_version("0.2.2");

        ::gnmi::ModelData* modeldata;
        ::gnmi::Encoding encoding;

        std::string cap;
        std::vector<std::string> capabilities;
        std::string org_name{};
        std::string module_name{};
        std::string revision_number{};
        std::string features{};
        std::string module_delim = "?module=";
        std::string rev_delim = "&revision=";

        for (auto cap : cap_array)
        {
            capabilities.push_back(cap);
            modeldata = response->add_supported_models();

            auto module_pos = cap.find(module_delim);
            if (module_pos != std::string::npos) {
                org_name = cap.substr(0, module_pos);
                cap.erase(0, module_pos + module_delim.length());
                modeldata->set_organization(org_name);
            }
            auto rev_pos = cap.find(rev_delim);
            if (rev_pos != std::string::npos) {
                module_name = cap.substr(0, rev_pos);
                modeldata->set_name(module_name);
                revision_number = cap.substr(rev_pos + rev_delim.length());
                modeldata->set_version(revision_number);
            }

            if(!(modeldata->organization()).empty())
            {
                std::cout << "Prefix: " << modeldata->organization() << std::endl;
            }
            if(!module_name.empty()) 
            { 
                std::cout << "Module: " << modeldata->name() << std::endl;
            }
            if(!revision_number.empty()) 
            {
                std::cout << "Revision: " << modeldata->version() << std::endl;
            }
            std::cout << std::endl;
        }
        response->add_supported_encodings(::gnmi::Encoding::JSON);
        return Status::OK;
    }

    void get_value(std::vector<std::string> path_container, json json_payload)
    {
        std::string path_to_string;
        for(int i = 0; i < path_container.size(); i++)
        {
            path_to_string.append("/");
            path_to_string.append(path_container.at(i));
        }
        json::json_pointer path_ptr(path_to_string);
    }

    Status Get(ServerContext* context, const GetRequest* request, GetResponse* response) override 
    {
        std::string response_payload;
        ::gnmi::TypedValue* value = new ::gnmi::TypedValue;
        std::string path_element;
        std::vector<std::string> path_container;
        int element_size;
        bool is_bgp = false;
        bool is_int = false;
        json json_payload;

        std::cout << "===========Get Request Received===========" << std::endl;
        std::cout << request->DebugString() << std::endl;

        auto notification = response->add_notification();
        std::time_t timestamp_val = std::time(nullptr);
        notification->set_timestamp(static_cast< ::google::protobuf::int64>(timestamp_val));
        std::string origin = request->prefix().origin();
        if (origin.length() > 0) {
            notification->mutable_prefix()->CopyFrom(request->prefix());
            if (origin == "openconfig-bgp") {
                is_bgp = true;
            }
            else if (origin == "openconfig-interfaces") {
                is_int = true;
            }
        }
        for(int j = 0; j < request->path_size(); ++j) 
        {
            ::gnmi::Path* response_path = new ::gnmi::Path;

            auto req_path = request->path(j);
            origin = req_path.origin();
            if (origin.length() > 0) {
                response_path->set_origin(origin);
                if (origin == "openconfig-bgp") {
                    is_bgp = true;
                }
                else if (origin == "openconfig-interfaces") {
                    is_int = true;
                }
            }
            for(int i = 0; i < req_path.elem_size(); ++i)
            {
                gnmi::PathElem* path_elem = response_path->add_elem();
                path_elem->CopyFrom(req_path.elem(i));
            }
            auto update = notification->add_update();
            update->set_allocated_path(response_path);

            if (set_counter == 1 && delete_counter == 0 && (is_bgp || is_int)) {
                std::cout << "DEBUG: Update Value Set by Create Request\n";
                if (is_bgp)
                    response_payload = "{\"global\": {\"config\": {\"as\":65172} }, \"neighbors\": {\"neighbor\": [{\"neighbor-address\":\"172.16.255.2\", \"config\": {\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}}";
                else if (is_int)
                    response_payload = "{\"interface\":[{\"name\":\"Loopback10\",\"config\":{\"name\":\"Loopback10\",\"description\":\"Test\"}}]}";
                json_payload = json::parse(response_payload);
                get_value(path_container, json_payload);
                value->set_json_ietf_val(response_payload);
                update->set_allocated_val(value);
            }
            else if (delete_counter == 1) {
                std::cout << "DEBUG: Update Value Deleted by Delete Request\n";
                response_payload = "{\"value\":\"null\"}";
            }
            else {
                response_payload = "{\"value\":\"null\"}";
            }
        }

        std::cout << "===========Get Response Sent===========" << std::endl;
        std::cout << response->DebugString() << std::endl;
        return Status::OK;
    } 

    Status Set(ServerContext* context, const SetRequest* request, SetResponse* response) override 
    {
        ::grpc::Status status;
        ::gnmi::Update update;
        ::gnmi::UpdateResult* update_response;
        ::gnmi::Path* path = new ::gnmi::Path;
        ::gnmi::UpdateResult_Operation operation;
        
        update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_DELETE);
        std::cout << "===========Set Request Received===========" << std::endl;
        std::cout << request->DebugString() << std::endl;

        update_response = response->add_response();
        std::time_t timestamp_val = std::time(nullptr);
        update_response->set_timestamp(static_cast< ::google::protobuf::int64>(timestamp_val));

        if (request->delete__size() >= 1)
        {
            delete_counter = 1;
            for(int i = 0; i < request->delete__size(); ++i)
            {
                ::gnmi::Path* response_path = new ::gnmi::Path;

                ::gnmi::Path delete_path = request->delete_(i);
                std::string origin = delete_path.origin();
                if (origin.length() > 0) {
                    response_path->set_origin(origin);
                }
                for(int j = 0; j < delete_path.elem_size(); ++j)
                {
                    gnmi::PathElem* path_elem = response_path->add_elem();
                    path_elem->CopyFrom(delete_path.elem(j));
                }
                update_response->set_allocated_path(response_path);
            }
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_DELETE);
            std::cout << "===========Set Response Sent===========" << std::endl;
            std::cout << response->DebugString() << std::endl;
            return Status::OK;
        } 
        else if (request->replace_size() > 0) {
            set_counter = 1;
            delete_counter = 0;
            for(int i = 0; i < request->replace_size(); ++i)
            {
                ::gnmi::Path* response_path = new ::gnmi::Path;

                auto replace_path = request->replace(i).path();
                std::string origin = replace_path.origin();
                if (origin.length() > 0) {
                    response_path->set_origin(origin);
                }
                for(int j = 0; j < replace_path.elem_size(); ++j)
                {
                    gnmi::PathElem* path_elem = response_path->add_elem();
                    path_elem->CopyFrom(replace_path.elem(j));
                }
                update_response->set_allocated_path(response_path);
            }
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_REPLACE);
            std::cout << "===========Set Response Sent===========" << std::endl;
            std::cout << response->DebugString() << std::endl;
            return Status::OK;
        } 
        else if (request->update_size() >= 1) {
            set_counter = 1;
            delete_counter = 0;
            for(int i = 0; i < request->update_size(); ++i)
            {
                ::gnmi::Path* response_path = new ::gnmi::Path;

                auto update_path = request->update(i).path();
                std::string origin = update_path.origin();
                if (origin.length() > 0) {
                    response_path->set_origin(origin);
                }
                for(int j = 0; j < update_path.elem_size(); ++j)
                {
                    gnmi::PathElem* path_elem = response_path->add_elem();
                    path_elem->CopyFrom(update_path.elem(j));
                }
                update_response->set_allocated_path(response_path);
            }
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_UPDATE);
            std::cout << "===========Set Response Sent===========" << std::endl;
            std::cout << response->DebugString() << std::endl;
            return Status::OK;
        } else {
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_INVALID); 
            std::cout << status.error_message() << std::endl;
            return Status::CANCELLED;
        }
    } 
};

void RunServer() 
{
    std::string server_address("0.0.0.0:50051");
    bool is_secure = false;

    gNMIImpl service;
    ServerBuilder builder; 
    if (is_secure)
    {
        // Secure Channel
        std::string server_key;
        std::string server_cert;
    
        std::ifstream kf("../keys/ems-key.pem");
        std::ifstream cf("../keys/ems.pem");
    
        server_key.assign((std::istreambuf_iterator<char>(kf)),(std::istreambuf_iterator<char>()));
        server_cert.assign((std::istreambuf_iterator<char>(cf)),(std::istreambuf_iterator<char>()));
    
        grpc::SslServerCredentialsOptions::PemKeyCertPair pkcp = {server_key, server_cert};
    
        grpc::SslServerCredentialsOptions ssl_opts;
    
        ssl_opts.pem_root_certs = "";
        ssl_opts.pem_key_cert_pairs.push_back(pkcp);
    
        std::cout << "server key: " << server_key << std::endl;
        std::cout << "server cert: " << server_cert << std::endl;
    
        builder.AddListeningPort(server_address, grpc::SslServerCredentials(ssl_opts));
    } else {
        builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());  
    }
    
    builder.RegisterService(&service);
    std::unique_ptr<Server> server(builder.BuildAndStart());
    std::cout << "Server listening on " << server_address << std::endl;
    server->Wait();
}

int main(int argc, char** argv) 
{
    RunServer();
    return 0;
}
