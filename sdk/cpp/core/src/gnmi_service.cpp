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
#include "gnmi_util.hpp"
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
    auto yfilter = entity.yfilter;
    entity.yfilter = YFilter::not_set;
    path::DataNode& datanode = get_data_node_from_entity(entity, root_schema);
    entity.yfilter = yfilter;

    string payload = codec_service.encode(datanode, EncodingFormat::JSON, false);
    YLOG_DEBUG("===========Generating Target Payload============");
    YLOG_DEBUG("{}", payload);
    return payload;
}

gNMIService::gNMIService()
{
}

gNMIService::~gNMIService()
{
}

//get
shared_ptr<Entity>
gNMIService::get(gNMIServiceProvider& provider, Entity& filter, const string & operation) const
{
    YLOG_INFO("Executing get gRPC");

    gnmi::Path* path = new gnmi::Path;
    parse_entity_to_path(filter, path);

    gNMIRequest request{};
    request.alias = filter.get_segment_path();
    request.payload = get_data_payload(provider, filter);
    request.path = path;
    request.type = "get";
    request.operation = operation;

    vector<gNMIRequest> get_request_list{};
    get_request_list.push_back(request);

    auto & gnmi_session = dynamic_cast<const path::gNMISession&> (provider.get_session());
    auto & client = gnmi_session.get_client();
    vector<string> reply = client.execute_get_operation(get_request_list, operation);

    auto root_dn = gnmi_session.handle_get_reply(reply);
    if (root_dn) {
        return read_datanode(filter, root_dn->get_children()[0]);
    }
    return nullptr;
}

vector<shared_ptr<Entity>>
gNMIService::get(gNMIServiceProvider & provider, vector<Entity*> & filter_list, const string & operation)
{
    YLOG_INFO("Executing get gRPC for multiple entities");

    vector<gNMIRequest> get_request_list{};
    vector<shared_ptr<Entity>> response_list{};
    for (auto filter : filter_list) {
        gnmi::Path* path = new gnmi::Path;
        parse_entity_to_path(*filter, path);

        gNMIRequest request{};
        request.alias = filter->get_segment_path();
        request.payload = get_data_payload(provider, *filter);
        request.path = path;
        request.type = "get";
        request.operation = operation;
        get_request_list.push_back(request);
    }
    auto & gnmi_session = dynamic_cast<const path::gNMISession&> (provider.get_session());
    auto & client = gnmi_session.get_client();
    vector<string> reply = client.execute_get_operation(get_request_list, operation);

    auto root_dn = gnmi_session.handle_get_reply(reply);
    if (root_dn) {
        // Build map of data nodes in order to retain filter list order
        map<string,std::shared_ptr<path::DataNode>> path_to_dn{};
        for (auto dn : root_dn->get_children()) {
            string path = dn->get_path();
            if (path.find("/") == 0)
                path = path.substr(1);
            path_to_dn[path] = dn;
        }

        // Build output list
        for (auto filter : filter_list) {
            auto dn = path_to_dn[filter->get_segment_path()];
            if (dn) {
                auto entity = read_datanode(*filter, dn);
                response_list.push_back(entity);
            }
            else {
                response_list.push_back((std::shared_ptr<Entity>) filter);
            }
        }
    }
    return response_list;
}

//set
static gNMIRequest build_set_request(gNMIServiceProvider& provider, Entity& entity)
{
	string operation = to_string(entity.yfilter);
	if (operation != "replace" && operation != "update" && operation != "delete")
    {
        YLOG_ERROR("gNMIService::set: {} operation not supported", operation );
        throw(YServiceProviderError{operation + " operation not supported"});
    }
    YLOG_INFO("Executing set gRPC operation '{}' on entity '{}'", operation, entity.get_segment_path());

    gnmi::Path* path = new gnmi::Path;
    string payload{};

    if (operation == "delete") {
        parse_entity_to_path(entity, path);
    }
    else {
        parse_entity_prefix(entity, path);
        payload = get_data_payload(provider, entity);
        auto pos = payload.find("{", 4);
        if (pos != string::npos)
            payload = payload.substr(pos, payload.length()-pos-1);
    }

    gNMIRequest request{};
    request.alias = entity.get_segment_path();
    request.path = path;
    request.payload = payload;
    request.type = "set";
    request.operation = operation;
    return request;
}

//set
bool gNMIService::set(gNMIServiceProvider& provider, Entity& entity) const
{
    vector<gNMIRequest> set_request_list{};
    gNMIRequest request = build_set_request(provider, entity);
    set_request_list.push_back(request);

    auto & gnmi_session = dynamic_cast<const path::gNMISession&> (provider.get_session());
    auto & client = gnmi_session.get_client();

    return client.execute_set_operation(set_request_list);
}

bool gNMIService::set(gNMIServiceProvider& provider, vector<Entity*> & entity_list) const
{
    vector<gNMIRequest> set_request_list{};
    for (auto entity : entity_list) {
        gNMIRequest request = build_set_request(provider, *entity);
        set_request_list.push_back(request);
    }
    auto & gnmi_session = dynamic_cast<const path::gNMISession&> (provider.get_session());
    auto & client = gnmi_session.get_client();

    return client.execute_set_operation(set_request_list);
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
        YLOG_ERROR("gNMIService::subscribe: list mode '{}' is not supported", mode);
        throw(YServiceProviderError{list_mode + " list mode is not supported"});
    }

    if (mode != "ON_CHANGE" && mode != "SAMPLE" && mode!= "TARGET_DEFINED")
    {
        YLOG_ERROR("gNMIService::subscribe: mode '{}' is not supported", mode);
        throw(YServiceProviderError{mode + " mode is not supported"});
    }

    YLOG_INFO("gNMIService::subscribe: Executing subscribe RPC in '{}' list mode", list_mode);
    std::pair<std::string, std::string> prefix_pair;
    std::vector<PathElem> path_container;
    parse_entity_to_prefix_and_paths(filter, prefix_pair, path_container);

    auto & gnmi_session = dynamic_cast<const path::gNMISession&> (provider.get_session());
    auto & client = gnmi_session.get_client();
    client.execute_subscribe_operation(prefix_pair, path_container, list_mode, qos, sample_interval, mode, func);
}

std::string
gNMIService::capabilities(gNMIServiceProvider & provider)
{
    auto & gnmi_session = dynamic_cast<const path::gNMISession&> (provider.get_session());

    auto cap_rdn = gnmi_session.handle_get_capabilities();

    ydk::path::Codec codec{};
    auto json_caps = codec.encode(*cap_rdn, ydk::EncodingFormat::JSON, true);
    return json_caps;
}

}
