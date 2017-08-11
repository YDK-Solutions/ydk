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

#include <iostream>
#include <sstream>

#include "gnmi_service.hpp"
#include "entity_data_node_walker.hpp"
#include "errors.hpp"
#include "path_api.hpp"
#include "codec_service.hpp"
#include "logger.hpp"

using namespace std;

using grpc::Channel;
using grpc::ChannelArguments;
using grpc::ChannelCredentials;
using grpc::SslCredentialsOptions;

namespace ydk {

static string get_gnmi_payload(gNMIServiceProvider& provider, Entity & entity);
static string set_gnmi_payload(gNMIServiceProvider& provider, Entity & entity);

gNMIService::gNMIService(string address)
: client(make_unique<gNMIClient>(grpc::CreateChannel(address, grpc::InsecureChannelCredentials())))
{
}

gNMIService::~gNMIService()
{
}

//get
string gNMIService::get(gNMIServiceProvider& provider, Entity& filter)
{
    YLOG_INFO("Executing get RPC");
    string gnmi_payload = get_gnmi_payload(provider, filter);
    string reply = client->execute_wrapper(gnmi_payload, "read");
    YLOG_DEBUG("=============Reply payload received from device=============");
    YLOG_DEBUG("{}", reply);
    YLOG_DEBUG("\n");
    return reply;
}

//set
string gNMIService::set(gNMIServiceProvider& provider, Entity& filter, string operation)
{
    YLOG_INFO("Executing set RPC");
    string gnmi_payload = set_gnmi_payload(provider, filter);
    string reply = client->execute_wrapper(gnmi_payload, operation);
    YLOG_DEBUG("=============Reply payload received from device=============");
    YLOG_DEBUG("{}", reply);
    YLOG_DEBUG("\n");
    return reply;
}

static string get_gnmi_payload(gNMIServiceProvider& provider, Entity & entity)
{
    path::Codec codec_service{};

    path::RootSchemaNode & root_schema = provider.get_session().get_root_schema();
    YLOG_DEBUG("Created root_schema");
    path::DataNode& datanode = get_data_node_from_entity(entity, root_schema);
    YLOG_DEBUG("Created datanode");

    string payload{"\"filter\":"};
    YLOG_DEBUG("Payload: {}", payload);
    payload+=codec_service.encode(datanode, EncodingFormat::JSON, false);
    YLOG_DEBUG("===========Generating Target Payload============");
    YLOG_DEBUG("{}", payload.c_str());
    YLOG_DEBUG("\n");
    return payload;
}

static string set_gnmi_payload(gNMIServiceProvider& provider, Entity & entity)
{
    path::Codec codec_service{};

    path::RootSchemaNode & root_schema = provider.get_session().get_root_schema();
    YLOG_DEBUG("Created root_schema");
    path::DataNode& datanode = get_data_node_from_entity(entity, root_schema);
    YLOG_DEBUG("Created datanode");

    string payload{"\"filter\":"};
    YLOG_DEBUG("Payload: {}", payload);
    payload+=codec_service.encode(datanode, EncodingFormat::JSON, false);
    YLOG_DEBUG("===========Generating Target Payload============");
    YLOG_DEBUG("{}", payload.c_str());
    YLOG_DEBUG("\n");
    return payload;
}
}