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

// Logic and data behind the server's behavior.
class gNMIImpl final : public gNMI::Service 
{    
    int set_counter = 0;
    int delete_counter = 0;
    bool is_secure = true;

    Status Capabilities(ServerContext* context, const CapabilityRequest* request, CapabilityResponse* response) override 
    {
        response->set_gnmi_version("0.2.2");

        // read cap file and populate model data
        std::ifstream capfile("/usr/local/share/ydk/0.0.0.0:50051/capabilities.txt");

        ::gnmi::ModelData* modeldata;
        ::gnmi::Encoding encoding;

        std::string cap;
        std::vector<std::string> capabilities;
        std::string org_name;
        std::string module_name;
        std::string revision_number;
        std::string name_delim = "?module=";
        std::string rev_delim = "&revision=";

        while (std::getline(capfile, cap)) 
        {
            capabilities.push_back(cap);
            if((cap.find(name_delim)==std::string::npos) && (cap.find(rev_delim)==std::string::npos)) 
            {
                org_name = cap.substr(0,cap.find(name_delim));
            } 
            else 
            {
                org_name = cap.substr(0,cap.find(name_delim));
                cap.erase(0,cap.find(name_delim) + name_delim.length());
                module_name = cap.substr(0, cap.find(rev_delim));
                revision_number = cap.substr(cap.find(rev_delim) + rev_delim.length());
            }

            modeldata = response->add_supported_models();
            modeldata->set_name(module_name);
            std::cout << "Organization: " << modeldata->organization() << std::endl;
            if(!(modeldata->organization()).empty())
            {
                std::cout << "Organization: " << modeldata->organization() << std::endl;
            }
            if(!module_name.empty()) 
            { 
                modeldata->set_organization(org_name);
                std::cout << "Module: " << modeldata->name() << std::endl;
            }
            if(!revision_number.empty()) 
            {
                modeldata->set_version(revision_number);
                std::cout << "Revision: " << modeldata->version() << std::endl;
                std::cout << "          -----------             \n" << std::endl;
            }
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
        ::gnmi::Notification* notification;
        std::string response_payload;
        ::gnmi::Update* update;
        ::gnmi::TypedValue* value = new ::gnmi::TypedValue;
        ::gnmi::Path* path = new ::gnmi::Path;
        ::gnmi::Path prefix;
        std::string prefix_element;
        std::string path_element;
        std::vector<std::string> path_container;
        int element_size;
        int is_bgp = 0;
        json json_payload;

        notification = response->add_notification();
        std::time_t timestamp_val = std::time(nullptr);
        notification->set_timestamp(static_cast< ::google::protobuf::int64>(timestamp_val));
        notification->mutable_prefix()->CopyFrom(request->prefix());

        update = notification->add_update();

        for(int i = 0; i < notification->prefix().element_size(); ++i) 
        {
          prefix_element.append("\"");
          prefix_element.append(notification->prefix().element(i));
          prefix_element.append("\"");
        }

        for(int j = 0; j < request->path_size(); ++j) 
        {
          for(int i = 0; i < request->path(j).element_size(); ++i) 
          {
            path_element.append(request->path(j).element(i));
            path->add_element(request->path(j).element(i));
            path_container.push_back(path_element);
            if(path_element == "openconfig-bgp:bgp") {
                is_bgp = 1;
            }
            path_element.clear();
          }
        }

        update->set_allocated_path(path);

        if (set_counter == 1 && delete_counter == 0 && is_bgp == 1){
            std::cout << "DEBUG: Update Value Set by Create Request\n";
            response_payload = "{\"global\": {\"config\": {\"as\":65172} }, \"neighbors\": {\"neighbor\": [{\"neighbor-address\":\"172.16.255.2\", \"config\": {\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}}";
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


        std::cout << "===========Get Request Received===========" << std::endl;
        std::cout << request->DebugString() << std::endl;
        std::cout << "===========Get Response Sent===========" << std::endl;
        std::cout << response->DebugString() << std::endl;
        return Status::OK;
    } 

    Status Set(ServerContext* context, const SetRequest* request, SetResponse* response) override 
    {
        ::grpc::Status status;
        ::gnmi::Update update;
        ::gnmi::Update replace;
        ::gnmi::UpdateResult* update_response;
        ::gnmi::Path* path = new ::gnmi::Path;
        ::gnmi::UpdateResult_Operation operation;
        
        update_response = response->add_response();

        std::time_t timestamp_val = std::time(nullptr);
        update_response->set_timestamp(static_cast< ::google::protobuf::int64>(timestamp_val));

        if (request->delete__size() >= 1)
        {
            delete_counter = 1;
            for(int i = 0; i < request->delete__size(); ++i)
            {
                ::gnmi::Path delete_path = request->delete_(i); 
                delete_path.add_element(delete_path.element(i));
            }
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_DELETE);
            std::cout << "===========Set Request Received===========" << std::endl;
            std::cout << request->DebugString() << std::endl;
            std::cout << "===========Set Response Sent===========" << std::endl;
            std::cout << response->DebugString() << std::endl;
            return Status::OK;
        } 
        else if (request->replace_size() >= 1) {
            for(int i = 0; i < request->replace_size(); ++i)
            {
                replace = request->replace(i); 
                for(int j = 0; j < replace.path().element_size(); ++j) 
                {
                    path->add_element(replace.path().element(j));
                }
            }
            update_response->set_allocated_path(path);
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_REPLACE);
            std::cout << "===========Set Request Received===========" << std::endl;
            std::cout << request->DebugString() << std::endl;
            std::cout << "===========Set Response Sent===========" << std::endl;
            std::cout << response->DebugString() << std::endl;
            return Status::OK;
        } 
        else if (request->update_size() >= 1) {
            set_counter = 1;
            delete_counter = 0;
            for(int i = 0; i < request->update_size(); ++i)
            {
                update = request->update(i); 
                for(int j = 0; j < update.path().element_size(); ++j) 
                {
                    path->add_element(update.path().element(j));
                }
            }
            update_response->set_allocated_path(path);
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_UPDATE);
            std::cout << "===========Set Request Received===========" << std::endl;
            std::cout << request->DebugString() << std::endl;
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
    bool is_secure = true;

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
