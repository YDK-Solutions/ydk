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

#include <fstream>
#include <libyang/libyang.h>

#include "../gnmi_provider.hpp"
#include "../ietf_parser.hpp"
#include "../logger.hpp"
#include "../gnmi_path_api.hpp"
#include "../ydk_yang.hpp"
#include "path_private.hpp"
#include "gnmi_util.hpp"

using grpc::Channel;
using grpc::ChannelArguments;
using grpc::ChannelCredentials;
using grpc::SslCredentialsOptions;

using namespace std;
using namespace ydk;

namespace ydk
{
namespace path
{
static path::SchemaNode* get_schema_for_operation(path::RootSchemaNode& root_schema, string operation);

static gNMISession::SecureChannelArguments get_channel_credentials();

const char* TEMP_CANDIDATE = "urn:ietf:params:netconf:capability:candidate:1.0";

// Secure
// Create a default SSL ChannelCredentials object
gNMISession::gNMISession(const std::string& address,
                   const std::string& username,
                   const std::string& password,
                   int port)
{
    is_secure = true;
	path::Repository repo;
    initialize(repo, address, username, password, port);
    YLOG_DEBUG("gNMISession: Connected to {} using Secure Channel", address, username, password, port);
}

gNMISession::gNMISession(Repository & repo,
                   const std::string& address,
                   const std::string& username,
                   const std::string& password,
                   int port)
{
    is_secure = true;
    initialize(repo, address, username, password, port);
    YLOG_DEBUG("gNMISession: Connected to {} using Secure Channel", address, username, password, port);
}

gNMISession::gNMISession(Repository & repo,
                   const std::string& address,
                   int port)
{
    is_secure = false;
    initialize(repo, address, "", "", port);
    YLOG_DEBUG("gNMISession: Connected to {} using Insecure Channel", address);
}

gNMISession::~gNMISession() = default;

void gNMISession::initialize(path::Repository & repo, const std::string& address, const std::string& username, const std::string& password, int port)
{
    IetfCapabilitiesParser capabilities_parser{};

    std::ostringstream address_buffer{};
    address_buffer << address << ":" << port;

    if (is_secure) {
        gNMISession::SecureChannelArguments input_args = get_channel_credentials();
        client = make_unique<gNMIClient>(grpc::CreateCustomChannel(address_buffer.str(), input_args.channel_creds, *(input_args.args)), username, password);
        client->connect();
    }
    else {
    	client = make_unique<gNMIClient>(grpc::CreateChannel(address_buffer.str(), grpc::InsecureChannelCredentials()));
    }
    server_capabilities = client->get_capabilities();
    std::vector<std::string> empty_caps;

    root_schema = repo.create_root_schema(capabilities_parser.parse(empty_caps));

    if(root_schema.get() == nullptr)
    {
        YLOG_ERROR("gNMISession::initialize: Root schema cannot be obtained");
        throw(YIllegalStateError{"Root schema cannot be obtained"});
    }
}

gNMISession::SecureChannelArguments get_channel_credentials()
{
    // Authenticate Server at Client
    string server_cert;
    ifstream rf("ems.pem");

    server_cert.assign((istreambuf_iterator<char>(rf)),(istreambuf_iterator<char>()));

    YLOG_DEBUG("gNMISession::get_channel_credentials: Server certificate:\n{}", server_cert);

    grpc::SslCredentialsOptions ssl_opts;
    auto args = std::make_shared<grpc::ChannelArguments>();
    gNMISession::SecureChannelArguments input_args;
    ssl_opts.pem_root_certs = server_cert;

    args->SetSslTargetNameOverride("2001:420:1101:1::b");
    //args->SetSslTargetNameOverride("ems.cisco.com");

    // TODO: Authenticate client at server
    /* string client_key, client_cert;
    ifstream kf("client.key");
    ifstream cf("client.pem");
    client_key.assign((istreambuf_iterator<char>(kf)),(istreambuf_iterator<char>()));
    client_cert.assign((istreambuf_iterator<char>(cf)),(istreambuf_iterator<char>()));
    ssl_opts = {server_cert, client_key, client_cert};*/

    auto channel_creds = grpc::SslCredentials(grpc::SslCredentialsOptions(ssl_opts));
    input_args.channel_creds = channel_creds;
    input_args.args = args;
    return input_args;
}

EncodingFormat gNMISession::get_encoding() const
{
    return EncodingFormat::JSON;
}

std::vector<std::string> gNMISession::get_capabilities() const
{
    return server_capabilities;
}

path::RootSchemaNode& gNMISession::get_root_schema() const
{
    return *root_schema;
}

void gNMISession::invoke(path::Rpc& rpc, std::function<void(const vector<string> &)> func) const
{
    path::SchemaNode* gnmi_sub = get_schema_for_operation(*root_schema, "ydk:gnmi-subscribe");

    path::SchemaNode* rpc_schema = &(rpc.get_schema_node());
    auto rpc_name = ((SchemaNodeImpl*)rpc_schema)->m_node->name;
    if(rpc_schema == gnmi_sub) {
        // TODO
        YLOG_ERROR("gNMISession::invoke: RPC '{}' is not supported", rpc_name);
    }
    else {
        YLOG_ERROR("gNMISession::invoke: RPC '{}' is not supported", rpc_name);
        throw(YOperationNotSupportedError{"RPC is not supported!"});
    }
}

shared_ptr<path::DataNode> gNMISession::invoke(path::Rpc& rpc) const
{
    path::SchemaNode* gnmi_get = get_schema_for_operation(*root_schema, "ydk:gnmi-get");
    path::SchemaNode* gnmi_set = get_schema_for_operation(*root_schema, "ydk:gnmi-set");
    path::SchemaNode* gnmi_cap = get_schema_for_operation(*root_schema, "ydk:gnmi-caps");

    path::SchemaNode* rpc_schema = &(rpc.get_schema_node());
    if(rpc_schema == gnmi_set)
    {
        handle_set(rpc);
        return nullptr;
    }
    else if(rpc_schema == gnmi_get)
    {
        return handle_get(rpc);
    }
    else if(rpc_schema == gnmi_cap)
    {
        return handle_get_capabilities();
    }
    else {
        auto rpc_name = ((SchemaNodeImpl*)rpc_schema)->m_node->name;
        YLOG_ERROR("gNMISession::invoke: RPC '{}' is not supported", rpc_name);
        throw(YOperationNotSupportedError{"RPC is not supported!"});
    }
    return nullptr;
}

shared_ptr<path::DataNode> gNMISession::invoke(path::DataNode& datanode) const
{
    throw(YOperationNotSupportedError{"gNMISession::invoke: action datanode is not supported!"});
    return nullptr;
}

static void populate_path_from_payload(gnmi::Path* path, const string & payload, RootSchemaNode & root_schema)
{
    Codec s{};
    auto root_dn = s.decode(root_schema, payload, EncodingFormat::JSON);
    if(!root_dn || root_dn->get_children().empty()) {
        YLOG_ERROR( "Codec service failed to decode datanode from JSON payload");
        throw(YError{"Problems deserializing JSON payload"});
    }
    auto child = (root_dn->get_children())[0].get();
    parse_datanode_to_path(child, path);
}

static gNMIRequest build_set_request(path::RootSchemaNode & root_schema, DataNode* request, const string & operation)
{
    gNMIRequest one_request{};
    one_request.type = "set";
    one_request.operation = operation;
    auto entity = request->find("entity");
    if (entity.empty()) {
        YLOG_ERROR("Failed to get 'entity' node from set RPC");
        throw(YInvalidArgumentError{"Failed to get 'entity' node from set RPC"});
    }
    path::DataNode* entity_node = entity[0].get();
    one_request.payload = entity_node->get_value();

    one_request.path = new gnmi::Path();
    if (operation == "delete") {
        populate_path_from_payload(one_request.path, one_request.payload, root_schema);
    }
    else {
        auto pos = one_request.payload.find("{", 4);
        if (pos != string::npos) {
            string prefix = one_request.payload.substr(2, pos-4);
            parse_prefix_to_path(prefix, one_request.path);
            one_request.payload = one_request.payload.substr(pos, one_request.payload.length()-pos-1);
        }
    }
    return one_request;
}

bool gNMISession::handle_set(path::Rpc& ydk_rpc) const
{
	vector<gNMIRequest> setRequest{};

    //path::SchemaNode* rpc_schema = &(ydk_rpc.get_schema_node());
    //auto rpc_name = ((SchemaNodeImpl*)rpc_schema)->m_node->name;

    auto delete_list = ydk_rpc.get_input_node().find("delete");
    if (!delete_list.empty()) {
        for (auto request : delete_list) {
            gNMIRequest one_request = build_set_request(get_root_schema(), request.get(), "delete");
            setRequest.push_back(one_request);
        }
    }

    auto replace_list = ydk_rpc.get_input_node().find("replace");
    if (!replace_list.empty()) {
        for (auto request : replace_list) {
            gNMIRequest one_request = build_set_request(get_root_schema(), request.get(), "replace");
            setRequest.push_back(one_request);
        }
    }

    auto update_list = ydk_rpc.get_input_node().find("update");
    if (!update_list.empty()) {
        for (auto request : update_list) {
            gNMIRequest one_request = build_set_request(get_root_schema(), request.get(), "update");
            setRequest.push_back(one_request);
        }
    }

    return client->execute_set_operation(setRequest);
}

shared_ptr<path::DataNode>
gNMISession::handle_get(path::Rpc& rpc) const
{
	vector<gNMIRequest> getRequest{};

    path::SchemaNode* rpc_schema = &(rpc.get_schema_node());
    auto rpc_name = ((SchemaNodeImpl*)rpc_schema)->m_node->name;
    auto request_list = rpc.get_input_node().find("request");
    if (request_list.empty()) {
        YLOG_ERROR("Failed to get 'request' node from '{}' RPC", rpc_name);
        throw(YInvalidArgumentError{"Failed to get 'request' node from RPC"});
    }

    string operation = "CONFIG";
    auto type = rpc.get_input_node().find("type");
    if (!type.empty()) {
        path::DataNode* type_node = type[0].get();
        operation = type_node->get_value();
    }

    for (auto request : request_list) {
        gNMIRequest one_request{};
        one_request.type = "get";
        one_request.operation = operation;

        auto alias = request.get()->find("alias");
        if (!alias.empty()) {
            path::DataNode* alias_node = alias[0].get();
            one_request.alias = alias_node->get_value();
        }

        auto entity = request.get()->find("entity");
        if (entity.empty()) {
            YLOG_ERROR("Failed to get 'entity' node from '{}' RPC", rpc_name);
            throw(YInvalidArgumentError{"Failed to get 'entity' node from RPC"});
        }
        path::DataNode* entity_node = entity[0].get();
        one_request.payload = entity_node->get_value();

        one_request.path = new gnmi::Path();
        populate_path_from_payload(one_request.path, one_request.payload, get_root_schema());

        getRequest.push_back(one_request);
    }

    vector<string> reply = client->execute_get_operation(getRequest, operation);
    YLOG_DEBUG("============= Reply payload received from device =============");
    for (auto response : reply) {
        YLOG_DEBUG("\n{}", response);
    }

    return handle_get_reply(reply);
}

shared_ptr<path::DataNode>
gNMISession::handle_get_reply(vector<string> reply_val) const
{
    path::Codec codec{};

    shared_ptr<path::DataNode> root_dn = codec.decode_json_output(get_root_schema(), reply_val);
    if (!root_dn) {
        YLOG_ERROR( "Codec service failed to decode JSON values from GetResponse");
        throw(YError{"Problems deserializing JSON output"});
    }
    return root_dn;
}

shared_ptr<path::DataNode>
gNMISession::handle_get_capabilities() const
{
	gNMICapabilityResponse reply = client->execute_get_capabilities();

	RootSchemaNodeImpl & rs_impl = dynamic_cast<ydk::path::RootSchemaNodeImpl &> (*root_schema);

    ydk::path::RootDataImpl* rd = new ydk::path::RootDataImpl{rs_impl, rs_impl.m_ctx, "/"};

	auto & output_dn = rd->create_datanode("ydk:gnmi-capabilities", "");

	for (auto model : reply.supported_models) {
		ostringstream so, sv;
	    so << "supported-models[name='" << model.name << "']/organization";
	    output_dn.create_datanode(so.str(), model.organization);

	    sv << "supported-models[name='" << model.name << "']/version";
	    output_dn.create_datanode(sv.str(), model.version);
	}

	for (auto e : reply.supported_encodings) {
		ostringstream s;
		s << "supported-encodings[.='" << e << "']";
		output_dn.create_datanode(s.str());
	}

	output_dn.create_datanode("gnmi-version", reply.gnmi_version);

    return shared_ptr<path::DataNode> (rd);
}

gNMIClient & gNMISession::get_client() const
{
    return *client;
}

static path::SchemaNode* get_schema_for_operation(path::RootSchemaNode & root_schema, string operation)
{
    auto c = root_schema.find(operation);
    if(c.empty())
    {
        YLOG_ERROR("'{}' rpc schema not found!", operation);
        throw(YIllegalStateError{"CRUD read rpc schema not found!"});
    }
    return c[0];
}
}    // namespace path
}  // namespace ydk
