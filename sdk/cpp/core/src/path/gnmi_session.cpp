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
#include "../path_api.hpp"
#include "../ydk_yang.hpp"

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

//static shared_ptr<path::Rpc> create_rpc_instance(path::RootSchemaNode & root_schema, string rpc_name);
//static path::DataNode& create_rpc_input(path::Rpc & gnmi_rpc);
//static string get_read_rpc_name(bool config);

static shared_ptr<path::DataNode> handle_edit_reply(const string & reply);

static bool is_config(path::Rpc & rpc);
static string get_filter_payload(path::Rpc & ydk_rpc);
static string get_config_payload(path::RootSchemaNode & root_schema, path::Rpc & rpc);
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

shared_ptr<path::DataNode> gNMISession::invoke(path::Rpc& rpc) const
{
    path::SchemaNode* create_schema = get_schema_for_operation(*root_schema, "ydk:create");
    path::SchemaNode* read_schema = get_schema_for_operation(*root_schema, "ydk:read");
    path::SchemaNode* update_schema = get_schema_for_operation(*root_schema, "ydk:update");
    path::SchemaNode* delete_schema = get_schema_for_operation(*root_schema, "ydk:delete");

    path::SchemaNode* rpc_schema = &(rpc.get_schema_node());
    shared_ptr<path::DataNode> datanode = nullptr;

    if(rpc_schema == create_schema || rpc_schema == delete_schema || rpc_schema == update_schema)
    {
        if (rpc_schema == create_schema)
        {
            return handle_edit(rpc, "create");
        }
        else if (rpc_schema == delete_schema)
        {
            return handle_edit(rpc, "delete");
        }
        else
        {
            return handle_edit(rpc, "update");
        }
    }
    else if(rpc_schema == read_schema)
    {
        return handle_read(rpc, "read");
    }
    else
    {
        YLOG_ERROR("gNMISession::invoke: RPC is not supported");
        throw(YOperationNotSupportedError{"RPC is not supported!"});
    }
    return datanode;
}

shared_ptr<path::DataNode> gNMISession::invoke(path::DataNode& datanode) const
{
    throw(YOperationNotSupportedError{"gNMISession::invoke: action datanode is not supported!"});
    return nullptr;
}

shared_ptr<path::DataNode> gNMISession::handle_edit(path::Rpc& ydk_rpc, const std::string & operation) const
{
    string config_payload = get_config_payload(*root_schema, ydk_rpc);

    return handle_edit_reply(execute_payload(config_payload, operation, false));
}

static shared_ptr<path::DataNode> handle_edit_reply(const string & reply)
{
    if(reply.find("Success") == string::npos)
    {
        YLOG_ERROR("No OK in reply received from device");
        throw(YServiceProviderError{reply});
    }

    // No error no output for edit-config
    return nullptr;
}


shared_ptr<path::DataNode> gNMISession::handle_read(path::Rpc& ydk_rpc, const std::string & operation) const
{
    bool config = is_config(ydk_rpc);

    string filter_value = get_filter_payload(ydk_rpc);

    return handle_read_reply(execute_payload(filter_value, operation, config), *root_schema);
}

string gNMISession::execute_payload(const string & payload, const string & operation, bool is_config) const
{
    string reply = client->execute_wrapper(payload, operation, is_config);
    YLOG_DEBUG("=============Reply payload received from device=============");
    YLOG_DEBUG("{}", reply.c_str());
    YLOG_DEBUG("\n");
    return reply;
}

shared_ptr<path::DataNode> gNMISession::handle_read_reply(string reply, path::RootSchemaNode & root_schema) const
{
    path::Codec codec_service{};
    auto empty_data = reply.find("data");
    if(empty_data == string::npos)
    {
        YLOG_INFO("Found empty data tag");
        return nullptr;
    }

    auto data_start = reply.find("\"data\":");
    if(data_start == string::npos)
    {
        YLOG_ERROR( "Can't find data tag in reply sent by device {}", reply.c_str());
        throw(YServiceProviderError{reply});
    }
    data_start+= sizeof("\"data\":") - 1;
    auto data_end = reply.find_last_of("}");
    if(data_end == string::npos)
    {
        YLOG_ERROR( "No end data tag found in reply sent by device {}", reply.c_str());
        throw(YError{"No end data tag found"});
    }

    string data = reply.substr(data_start, data_end-data_start + 1);

    auto datanode = shared_ptr<path::DataNode>(codec_service.decode(root_schema, data, EncodingFormat::JSON));
    if(!datanode)
    {
        YLOG_ERROR( "Codec service failed to decode datanode");
        throw(YError{"Problems deserializing output"});
    }
    return datanode;
}

gNMIClient & gNMISession::get_client() const
{
    return *client;
}

static string get_config_payload(path::RootSchemaNode & root_schema,
    path::Rpc & rpc)
{
    path::Codec codec_service{};
    auto entity = rpc.get_input_node().find("entity");
    if(entity.empty()){
        YLOG_ERROR("Failed to get entity node");
        throw(YInvalidArgumentError{"Failed to get entity node"});
    }

    path::DataNode* entity_node = entity[0].get();
    string entity_value = entity_node->get_value();

    //deserialize the entity_value
    auto datanode = codec_service.decode(root_schema, entity_value, EncodingFormat::JSON);

    if(!datanode){
        YLOG_ERROR("Failed to decode entity node");
        throw(YInvalidArgumentError{"Failed to decode entity node"});
    }

    string config_payload {};

    for(auto const & child : datanode->get_children())
    {
        config_payload += codec_service.encode(*child, EncodingFormat::JSON, false);
    }
    return config_payload;
}

static bool is_config(path::Rpc & rpc)
{
    if(!rpc.get_input_node().find("only-config").empty())
    {
        return true;
    }
    return false;
}

static string get_filter_payload(path::Rpc & ydk_rpc)
{
    auto entity = ydk_rpc.get_input_node().find("filter");
    if(entity.empty())
    {
        YLOG_ERROR("Failed to get entity node.");
        throw(YInvalidArgumentError{"Failed to get entity node"});
    }

    auto datanode = entity[0];
    return datanode->get_value();
}

static string get_gnmi_payload(path::DataNode & input, string data_tag, string data_value)
{
    path::Codec codec_service{};
    input.create_datanode(data_tag, data_value);
    string payload{"\"rpc\":"};
    payload+=codec_service.encode(input, EncodingFormat::JSON, false);
    YLOG_DEBUG("===========Generating Target Payload============");
    YLOG_DEBUG("{}", payload.c_str());
    YLOG_DEBUG("\n");
    return payload;
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
}
}
