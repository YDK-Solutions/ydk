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
#include "codec_service.hpp"
#include "entity_data_node_walker.hpp"
#include "errors.hpp"
#include "logger.hpp"
#include "gnmi_provider.hpp"
#include "gnmi_service.hpp"

using namespace std;

using grpc::Channel;
using grpc::ChannelArguments;
using grpc::ChannelCredentials;
using grpc::SslCredentialsOptions;

namespace ydk {

static shared_ptr<Entity> execute_get_gnmi(gNMIServiceProvider& provider, Entity & entity);
static string execute_set_gnmi(gNMIServiceProvider& provider, Entity & entity);
static shared_ptr<Entity> read_datanode(Entity & filter, shared_ptr<path::DataNode> read_data_node);
static shared_ptr<Entity> get_top_entity_from_filter(Entity & filter);

gNMIService::gNMIService()
{
}

gNMIService::~gNMIService()
{
}

//get
shared_ptr<Entity> gNMIService::get(gNMIServiceProvider& provider, Entity& filter)
{
    YLOG_INFO("Executing get RPC");
    return execute_get_gnmi(provider, filter);
}

//set
bool gNMIService::set(gNMIServiceProvider& provider, Entity& filter, const string & operation)
{
    if(operation != "gnmi_create" && operation != "gnmi_delete")
    {
        YLOG_ERROR("{} operation not supported", operation );
        throw(YCPPServiceProviderError{operation + " operation not supported"});
    }
    YLOG_INFO("Executing set RPC with {} operation", operation);
    string gnmi_payload = execute_set_gnmi(provider, filter);
    string reply = (dynamic_cast<const path::gNMISession&>(provider.get_session())).execute_payload(gnmi_payload, operation);
    YLOG_DEBUG("=============Reply payload received from device=============");
    YLOG_DEBUG("{}\n", reply);
    if(!reply.empty()) return true;
    else return false;
}

static shared_ptr<Entity> execute_get_gnmi(gNMIServiceProvider& provider, Entity & entity)
{
    path::Codec codec_service{};
    path::RootSchemaNode & root_schema = provider.get_session().get_root_schema();
    path::DataNode& datanode = get_data_node_from_entity(entity, root_schema);

    string payload{"\"filter\":"};
    payload+=codec_service.encode(datanode, EncodingFormat::JSON, false);
    YLOG_DEBUG("===========Generating Target Payload============");
    YLOG_DEBUG("{}\n", payload.c_str());
    string reply = (dynamic_cast<const path::gNMISession&>(provider.get_session())).execute_payload(payload, "read");
    auto output = (dynamic_cast<const path::gNMISession&>(provider.get_session())).handle_read_reply(reply, root_schema);
    return read_datanode(entity, output);
}

static string execute_set_gnmi(gNMIServiceProvider& provider, Entity & entity)
{
    path::Codec codec_service{};
    path::RootSchemaNode & root_schema = provider.get_session().get_root_schema();
    path::DataNode& datanode = get_data_node_from_entity(entity, root_schema);

    string payload{"\"filter\":"};
    payload+=codec_service.encode(datanode, EncodingFormat::JSON, false);
    YLOG_DEBUG("===========Generating Target Payload============");
    YLOG_DEBUG("{}", payload.c_str());
    return payload;
}

static shared_ptr<Entity> read_datanode(Entity & filter, shared_ptr<path::DataNode> read_data_node)
{
    if (read_data_node == nullptr)
        return {};
    shared_ptr<Entity> top_entity = get_top_entity_from_filter(filter);
    get_entity_from_data_node(read_data_node->get_children()[0].get(), top_entity);
    return top_entity;
}

static shared_ptr<Entity> get_top_entity_from_filter(Entity & filter)
{
    if(filter.parent == nullptr)
        return filter.clone_ptr();

    return get_top_entity_from_filter(*(filter.parent));
}
}
