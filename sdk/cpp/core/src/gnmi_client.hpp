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

#ifndef _YDK_GNMI_CLIENT_H_
#define _YDK_GNMI_CLIENT_H_
#include <fstream>
#include <iostream>
#include <map>
#include <sstream>
#include <string>
#include <vector>

#include <grpc++/client_context.h>
#include <grpc++/create_channel.h>
#include <grpc++/grpc++.h>
#include "gnmi.grpc.pb.h"
#include "gnmi.pb.h"

using namespace gnmi;

using gnmi::gNMI;
using gnmi::CapabilityRequest;
using gnmi::CapabilityResponse;
using grpc::Channel;
using grpc::ClientContext;
using gnmi::GetRequest;
using gnmi::GetResponse;
using gnmi::SetRequest;
using gnmi::SetResponse;
using grpc::Status;

namespace ydk 
{

    class gNMIClient 
    {
    public:
        typedef struct PathPrefixValueFlags
        {
            bool path_has_value;
            bool prefix_has_value;
        } PathPrefixValueFlags;

        gNMIClient(std::shared_ptr<Channel> channel);
        ~gNMIClient();

        int connect(std::string address, bool is_secure);
        std::string execute_wrapper(const std::string & payload, const std::string& operation);
        std::string execute_get_payload(const ::gnmi::GetRequest& request, ::gnmi::GetResponse* response);
        std::string execute_set_payload(const ::gnmi::SetRequest& request, ::gnmi::SetResponse* response);
        std::vector<std::string> get_capabilities();
    
    private:
        std::string get_path_from_update(::gnmi::Update update);
        std::string get_prefix_from_notification(::gnmi::Notification notification);
        std::string get_value_from_update(::gnmi::Update update);
        bool has_gnmi_version(::gnmi::CapabilityResponse* response);
        void parse_capabilities_modeldata(::gnmi::CapabilityResponse* response);
        void parse_capabilities_encodings(::gnmi::CapabilityResponse* response);
        bool parse_capabilities(::gnmi::CapabilityResponse* response);
        std::string parse_get_response(::gnmi::GetResponse* response); 
        std::string parse_set_response(::gnmi::SetResponse* response); 
        ::gnmi::GetRequest populate_get_request(::gnmi::GetRequest request, const std::string& payload);
        ::gnmi::SetRequest populate_set_request(::gnmi::SetRequest request, const std::string& payload, const std::string& operation);
        std::unique_ptr<gNMI::Stub> stub_;
    };
}

#endif /* _YDK_GNMI_CLIENT_H_ */
