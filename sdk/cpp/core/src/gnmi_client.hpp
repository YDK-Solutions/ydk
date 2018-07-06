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
using gnmi::SubscribeRequest;
using gnmi::SubscribeResponse;
using grpc::Status;

namespace ydk 
{

class Entity;

struct gNMIRequest {
	std::string alias;
	std::string payload;
	gnmi::Path* path;
	std::string type;
	std::string operation;
};

struct gNMIModelData {
	std::string name;
	std::string organization;
	std::string version;
};

struct gNMICapabilityResponse {
    std::string gnmi_version;
    std::vector<gNMIModelData>  supported_models;
    std::vector<std::string> supported_encodings;
};

class gNMIClient
{
public:
    typedef struct PathPrefixValueFlags
    {
        bool path_has_value;
        bool prefix_has_value;
    } PathPrefixValueFlags;

    gNMIClient(std::shared_ptr<Channel> channel, const std::string & username, const std::string & password);
    gNMIClient(std::shared_ptr<Channel> channel);
    ~gNMIClient();

    int connect();

    std::vector<std::string> execute_get_operation(const std::vector<gNMIRequest> get_request_list, const std::string& operation);

    bool execute_set_operation(const std::vector<gNMIRequest> get_request_list);

    void execute_subscribe_operation(std::pair<std::string, std::string> & prefix,
                                    std::vector<PathElem> & path_container,
                                    const std::string & list_mode,
                                    long long qos,
                                    int sample_interval,
                                    const std::string & mode,
                                    std::function<void(const std::string &)> func);

    std::vector<std::string> get_capabilities();
    gNMICapabilityResponse execute_get_capabilities();

private:

    void parse_capabilities_modeldata(::gnmi::CapabilityResponse* response);
    void parse_capabilities(::gnmi::CapabilityResponse* response);

    std::vector<std::string> execute_get_payload(const ::gnmi::GetRequest& request, ::gnmi::GetResponse* response);
    bool execute_set_payload(const ::gnmi::SetRequest& request, ::gnmi::SetResponse* response);

    std::vector<std::string> capabilities;
    std::unique_ptr<gNMI::Stub> stub_;
    std::string username;
    std::string password;
    bool is_secure;
};
}

#endif /* _YDK_GNMI_CLIENT_H_ */
