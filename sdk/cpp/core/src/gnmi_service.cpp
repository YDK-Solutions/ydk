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
#include "entity_util.hpp"
#include "errors.hpp"
#include "logger.hpp"
#include "gnmi_provider.hpp"
#include "gnmi_service.hpp"
#include "common_utilities.hpp"

using namespace std;

using grpc::Channel;
using grpc::ChannelArguments;
using grpc::ChannelCredentials;
using grpc::SslCredentialsOptions;

namespace ydk {

static string get_data_payload(gNMIServiceProvider& provider, Entity & entity)
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

gNMIService::gNMIService()
{
}

gNMIService::~gNMIService()
{
}

//get
shared_ptr<Entity> gNMIService::get(gNMIServiceProvider& provider, Entity& filter, bool only_config) const
{
    YLOG_INFO("Executing get RPC");
    std::pair<std::string, std::string> prefix_pair;
    std::vector<PathElem> path_container;
    parse_entity_to_prefix_and_paths(filter, prefix_pair, path_container);

    path::RootSchemaNode & root_schema = provider.get_session().get_root_schema();

    auto & client = (dynamic_cast<const path::gNMISession&>(provider.get_session())).get_client();
    string reply = client.execute_get_operation(prefix_pair, path_container, only_config);
    auto output = (dynamic_cast<const path::gNMISession&>(provider.get_session())).handle_read_reply(reply, root_schema);
    return read_datanode(filter, output);
}

//set
bool gNMIService::set(gNMIServiceProvider& provider, Entity& filter, const string & operation) const
{
    if(operation != "gnmi_create" && operation != "gnmi_delete")
    {
        YLOG_ERROR("{} operation not supported", operation );
        throw(YServiceProviderError{operation + " operation not supported"});
    }
    YLOG_INFO("Executing set RPC with {} operation", operation);
    string data_payload = get_data_payload(provider, filter);

    string reply = (dynamic_cast<const path::gNMISession&>(provider.get_session())).execute_payload(data_payload, operation, true);
    if(!reply.empty())
        return true;
    return false;
}

//subscribe
void gNMIService::subscribe(gNMIServiceProvider& provider,
                            Entity& filter,
                            const std::string & list_mode,
                            long long qos,
                            const std::string & mode,
                            int sample_interval,
                            std::function<void(const std::string &)> func) const
{
    if(list_mode != "ONCE" && list_mode != "STREAM" && list_mode != "POLL")
    {
        YLOG_ERROR("{} list mode not supported", mode);
        throw(YServiceProviderError{list_mode + " list mode not supported"});
    }

    if(mode != "ON_CHANGE" && mode != "SAMPLE" && mode!= "TARGET_DEFINED")
    {
        YLOG_ERROR("{} mode not supported", mode);
        throw(YServiceProviderError{mode + " mode not supported"});
    }

    YLOG_INFO("Executing subscribe RPC in {} list mode", list_mode);
    std::pair<std::string, std::string> prefix_pair;
    std::vector<PathElem> path_container;
    parse_entity_to_prefix_and_paths(filter, prefix_pair, path_container);
    auto & client = (dynamic_cast<const path::gNMISession&>(provider.get_session())).get_client();
    client.execute_subscribe_operation(prefix_pair, path_container, list_mode, qos, sample_interval, mode, func);
}

}
