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
#include <thread>
#include <chrono>

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

using grpc::ServerReaderWriter;
using gnmi::SubscribeResponse;
using gnmi::SubscribeRequest;
using gnmi::SubscriptionList;
using gnmi::Subscription;

using json = nlohmann::json;

static std::string cap_array[] =
    {"IETF NETMOD Working Group?module=ietf-aug-base-1&revision=2016-07-01",
     "IETF NETMOD Working Group?module=ietf-aug-base-2&revision=2016-07-01",
     "YDK?module=oc-pattern&revision=2015-11-17",
     "YDK?module=main&revision=2015-11-17",
     "YDK?module=main-aug1&revision=2015-11-17",
     "YDK?module=ydktest-aug-ietf-5&revision=2017-07-26",
     "YDK?module=ydktest-aug-ietf-4&revision=2016-06-27",
     "YDK?module=ydktest-aug-ietf-1&revision=2016-06-17",
     "YDK?module=ydktest-aug-ietf-2&revision=2016-06-22",
     "YDK?module=ydktest-filterread&revision=2015-11-17",
     "YDK?module=ydktest-sanity&revision=2015-11-17",		//&features=ipv6-privacy-autoconf,ipv4-non-contiguous-netmasks",
     "YDK?module=ydktest-sanity-augm&revision=2015-11-17",
     "YDK?module=ydktest-sanity-types&revision=2016-04-11",
     "YDK?module=ydktest-types&revision=2016-05-23",		//&features=crypto",
     "OpenConfig working group?module=openconfig-bgp&revision=2016-06-21",
     "OpenConfig working group?module=openconfig-bgp-policy&revision=2016-06-21",
     "OpenConfig working group?module=openconfig-bgp-types&revision=2016-06-21",
     "OpenConfig working group?module=openconfig-interfaces&revision=2016-05-26",
     "OpenConfig working group?module=openconfig-if-ethernet&revision=2016-05-26",
     "OpenConfig working group?module=openconfig-extensions&revision=2015-10-09",
     "OpenConfig working group?module=openconfig-types&revision=2016-05-31",
     "OpenConfig working group?module=openconfig-platform&revision=2016-06-06",		//&deviations=cisco-xr-openconfig-platform-deviations",
     "OpenConfig working group?module=openconfig-platform-types&revision=2016-06-06",
     "OpenConfig working group?module=openconfig-platform-transceiver&revision=2016-05-24",
     "OpenConfig working group?module=openconfig-policy-types&revision=2016-05-12",
     "OpenConfig working group?module=openconfig-routing-policy&revision=2016-05-12",
     "OpenConfig working group?module=openconfig-terminal-device&revision=2016-06-17",
     "OpenConfig working group?module=openconfig-transport-types&revision=2016-06-17"
    };

static std::string bgp_payload =
		"{\"global\":{\"config\":{\"as\":65172} },\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}}";
static std::string int_payload =
		"{\"interface\":[{\"name\":\"Loopback10\",\"config\":{\"name\":\"Loopback10\",\"description\":\"Test\"}}]}";
static std::string null_payload =
		"{\"value\":\"null\"}";

// Logic and data behind the server's behavior.
class gNMIImpl final : public gNMI::Service
{    
    bool int_set = false;
    bool bgp_set = false;

    Status Capabilities(ServerContext* context, const CapabilityRequest* request, CapabilityResponse* response) override 
    {
        std::string module_delim = "?module=";
        std::string rev_delim = "&revision=";

        for (auto cap : cap_array)
        {
            ::gnmi::ModelData* modeldata = response->add_supported_models();

            auto module_pos = cap.find(module_delim);
            auto rev_pos = cap.find(rev_delim);
            if (module_pos != std::string::npos && rev_pos != std::string::npos) {
                std::string org_name = cap.substr(0, module_pos);
                modeldata->set_organization(org_name);

                std::string module_name = cap.substr(module_pos+module_delim.length(), rev_pos-module_pos-module_delim.length());
                modeldata->set_name(module_name);

                std::string revision_number = cap.substr(rev_pos + rev_delim.length());
                modeldata->set_version(revision_number);
            }
        }

        response->add_supported_encodings(::gnmi::Encoding::JSON);
        response->add_supported_encodings(::gnmi::Encoding::JSON_IETF);

        response->set_gnmi_version("0.4.0");

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

    const std::string get_response_payload(const std::string origin, const std::string last_elem)
    {
        std:: string response_payload = null_payload;
        if (origin == "openconfig-bgp" && bgp_set) {
            if (last_elem == "bgp")
                response_payload = bgp_payload;
            else if (last_elem == "global")
                response_payload = "{\"config\":{\"as\":65172} }";
            else if (last_elem == "as")
                response_payload = "65172";
     	    else if (last_elem == "neighbors")
                response_payload = "{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}";
     	    else if (last_elem == "neighbor")
                response_payload = "[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]";
     	    else if (last_elem == "config")
                response_payload = "{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}";
     	    else if (last_elem == "neighbor-address")
                response_payload = "\"172.16.255.2\"";
     	    else if (last_elem == "peer-as")
                response_payload = "65172";
        }
        else if (origin == "openconfig-interfaces" && int_set) {
            if (last_elem == "interfaces")
            	response_payload = int_payload;
            else if (last_elem == "interface")
            	response_payload = "[{\"name\":\"Loopback10\",\"config\":{\"name\":\"Loopback10\",\"description\":\"Test\"}}]";
            else if (last_elem == "config")
                response_payload = "{\"name\":\"Loopback10\",\"description\":\"Test\"}";
            else if (last_elem == "description")
                response_payload = "\"Test\"";
        }
        return response_payload;
    }

    Status Get(ServerContext* context, const GetRequest* request, GetResponse* response) override 
    {
        //std::cout << "=========== Get Request Received ===========" << std::endl;
        //std::cout << request->DebugString() << std::endl;

        for(int j = 0; j < request->path_size(); ++j) 
        {
            ::gnmi::Path* response_path = new ::gnmi::Path;
            auto notification = response->add_notification();
            std::time_t timestamp_val = std::time(nullptr);
            notification->set_timestamp(static_cast< ::google::protobuf::int64>(timestamp_val));

            auto req_path = request->path(j);
            std::string last_elem{};
            std::string origin = req_path.origin();
            if (origin.length() > 0) {
                response_path->set_origin(origin);
            }
            for (int i = 0; i < req_path.elem_size(); ++i)
            {
                gnmi::PathElem* path_elem = response_path->add_elem();
                path_elem->CopyFrom(req_path.elem(i));
                last_elem = path_elem->name();
            }
            auto update = notification->add_update();
            update->set_allocated_path(response_path);

            ::gnmi::TypedValue* value = new ::gnmi::TypedValue;

            //std::vector<std::string> path_container;
            //json json_payload = json::parse(response_payload);
            //get_value(path_container, json_payload);

            auto response_payload = get_response_payload(origin, last_elem);
            value->set_json_ietf_val(response_payload);
            update->set_allocated_val(value);
        }

        //std::cout << "=========== Get Response Sent ===========" << std::endl;
        //std::cout << response->DebugString() << std::endl;
        return Status::OK;
    } 

    Status Set(ServerContext* context, const SetRequest* request, SetResponse* response) override 
    {
        //std::cout << "=========== Set Request Received ===========" << std::endl;
        //std::cout << request->DebugString() << std::endl;

        std::time_t timestamp_val = std::time(nullptr);
        response->set_timestamp(static_cast< ::google::protobuf::int64>(timestamp_val));

        for (int i = 0; i < request->delete__size(); ++i)
        {
            ::gnmi::UpdateResult* update_response = response->add_response();

            ::gnmi::Path* response_path = new ::gnmi::Path();
            ::gnmi::Path delete_path = request->delete_(i);
            std::string origin = delete_path.origin();
            if (origin.length() > 0) {
                response_path->set_origin(origin);
                if (origin == "openconfig-bgp" && bgp_set) {
                    bgp_set = false;
                }
                else if (origin == "openconfig-interfaces" && int_set) {
                    int_set = false;
                }
            }
            for(int j = 0; j < delete_path.elem_size(); ++j)
            {
                gnmi::PathElem* path_elem = response_path->add_elem();
                path_elem->CopyFrom(delete_path.elem(j));
            }
            update_response->set_allocated_path(response_path);
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_DELETE);
        } 

        for (int i = 0; i < request->replace_size(); ++i)
        {
            ::gnmi::UpdateResult* update_response = response->add_response();

            ::gnmi::Path* response_path = new ::gnmi::Path;
            ::gnmi::Update request_update= request->replace(i);
            auto replace_path = request_update.path();
            std::string origin = replace_path.origin();
            if (origin.length() > 0) {
                response_path->set_origin(origin);
                if (origin == "openconfig-bgp") {
                    bgp_set = true;
                }
                else if (origin == "openconfig-interfaces") {
                    int_set = true;
                }
            }
            for(int j = 0; j < replace_path.elem_size(); ++j)
            {
                gnmi::PathElem* path_elem = response_path->add_elem();
                path_elem->CopyFrom(replace_path.elem(j));
            }
            update_response->set_allocated_path(response_path);
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_REPLACE);
        }

        for (int i = 0; i < request->update_size(); ++i)
        {
            ::gnmi::UpdateResult* update_response = response->add_response();

            ::gnmi::Path* response_path = new ::gnmi::Path;
            ::gnmi::Update request_update= request->update(i);
            auto replace_path = request_update.path();
            std::string origin = replace_path.origin();
            if (origin.length() > 0) {
                response_path->set_origin(origin);
                if (origin == "openconfig-bgp") {
                    bgp_set = true;
                }
                else if (origin == "openconfig-interfaces") {
                    int_set = true;
                }
            }
            for(int j = 0; j < replace_path.elem_size(); ++j)
            {
                gnmi::PathElem* path_elem = response_path->add_elem();
                path_elem->CopyFrom(replace_path.elem(j));
            }
            update_response->set_allocated_path(response_path);
            update_response->set_op(::gnmi::UpdateResult_Operation::UpdateResult_Operation_UPDATE);
        }
        //std::cout << "=========== Set Response Sent ===========" << std::endl;
        //std::cout << response->DebugString() << std::endl;
        return Status::OK;
    }

    Status Subscribe(ServerContext* context, ServerReaderWriter<SubscribeResponse, SubscribeRequest>* stream) override
    {
        SubscribeRequest request{};
        SubscriptionList subscription_list{};

        while (stream->Read(&request))
        {
            //std::cout << "=========== Subscribe Request Received ===========" << std::endl;
            //std::cout << request.DebugString() << std::endl;

            if (request.has_subscribe())
                subscription_list = request.subscribe();

            if (subscription_list.subscription_size() > 0)
            {
                if (subscription_list.mode() == SubscriptionList::ONCE)
                {
                    SubscribeResponse response = get_subscribe_response( subscription_list);

                    // Write message to the wire
                    ::grpc::WriteOptions options{};
                    options.set_last_message();
                    options.set_write_through();
                    stream->Write( response, options);
                    break;
                }
                else
                if (subscription_list.mode() == SubscriptionList::STREAM)
                {
                    Subscription sub = subscription_list.subscription(0);
                    ::google::protobuf::uint64 timer = 0;
                    while (timer < sub.heartbeat_interval())
                    {
                        SubscribeResponse response = get_subscribe_response( subscription_list);
                        // Write message to the wire
                        ::grpc::WriteOptions options{};
                        options.set_write_through();
                        timer += sub.sample_interval();
                        if (timer >= sub.heartbeat_interval())
                            options.set_last_message();
                        stream->Write( response, options);
                        std::this_thread::sleep_for(std::chrono::milliseconds(sub.sample_interval()/1000000));
                    }
                    break;
                }

                if (subscription_list.mode() == SubscriptionList::POLL && request.has_poll())
                {
                    request.release_poll();
                    SubscribeResponse response = get_subscribe_response( subscription_list);

                    // Write message to the wire
                    ::grpc::WriteOptions options{};
                    options.set_write_through();
                    stream->Write( response, options);
                }
            }
        }
        return Status::OK;
    }

 private:
    SubscribeResponse get_subscribe_response(SubscriptionList& subscription_list)
    {
        SubscribeResponse response{};

        ::gnmi::Notification* notification = new ::gnmi::Notification;
        std::time_t timestamp_val = std::time(nullptr);
        notification->set_timestamp(static_cast< ::google::protobuf::int64>(timestamp_val));

        for (int s = 0; s < subscription_list.subscription_size(); ++s)
        {
            Subscription sub = subscription_list.subscription(s);

            // Build path
            auto req_path = sub.path();
            ::gnmi::Path* response_path = new ::gnmi::Path;
            std::string origin = req_path.origin();
            std::string last_elem{};
            if (origin.length() > 0) {
                response_path->set_origin(origin);
            }
            for (int j = 0; j < req_path.elem_size(); ++j)
            {
                gnmi::PathElem* path_elem = response_path->add_elem();
                path_elem->CopyFrom(req_path.elem(j));
                last_elem = path_elem->name();
            }

            // Build Update
            ::gnmi::Update* update_response = notification->add_update();

            update_response->set_allocated_path(response_path);

            ::gnmi::TypedValue* value = new ::gnmi::TypedValue;
            auto response_payload = get_response_payload(origin, last_elem);
            value->set_json_ietf_val(response_payload);
            update_response->set_allocated_val(value);
        }
        response.set_allocated_update(notification);
        //std::cout << "=========== Subscribe Response Sent ===========" << std::endl;
        //std::cout << response.DebugString() << std::endl;
        return response;
    }
};

void RunServer(bool is_secure)
{
    std::string server_address("127.0.0.1:50051");

    gNMIImpl service;
    ServerBuilder builder;
    if (is_secure)
    {
        std::cout << "Starting YDK gNMI Server in secure mode" << std::endl;

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

        //std::cout << "server key: " << server_key << std::endl;
        //std::cout << "server cert: " << server_cert << std::endl;

        builder.AddListeningPort(server_address, grpc::SslServerCredentials(ssl_opts));
    } else {
        std::cout << "Starting YDK gNMI Server in non-secure mode" << std::endl;
        builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());  
    }
    
    builder.RegisterService(&service);
    std::unique_ptr<Server> server(builder.BuildAndStart());
    std::cout << "Server listening on " << server_address << std::endl;
    server->Wait();
}

int main(int argc, char** argv) 
{
    bool enable_ssl = false;
    if (argc > 1) {
        std::string param = argv[1];
        if (param == "-s")
            enable_ssl = true;
    }
    RunServer(enable_ssl);
    return 0;
}
