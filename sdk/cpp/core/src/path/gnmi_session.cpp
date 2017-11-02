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

static shared_ptr<path::Rpc> create_rpc_instance(path::RootSchemaNode & root_schema, string rpc_name);
static path::DataNode& create_rpc_input(path::Rpc & gnmi_rpc);

static bool is_candidate_supported(vector<string> capabilities);
static void create_input_target(path::DataNode & input, bool candidate_supported);
static void create_input_source(path::DataNode & input, bool config);
static void create_input_error_option(path::DataNode & input);
static string get_read_rpc_name(bool config);
static string get_commit_rpc_payload();
static shared_ptr<path::DataNode> handle_edit_reply(const string & reply, gNMIClient & client, bool candidate_supported, const string & operation);

static bool is_config(path::Rpc & rpc);
static string get_filter_payload(path::Rpc & ydk_rpc);
static string get_gnmi_payload(path::DataNode & input, string data_tag, string data_value);
static string get_config_payload(path::RootSchemaNode & root_schema, path::Rpc & rpc);
static gNMISession::SecureChannelArguments get_channel_credentials();

const char* TEMP_CANDIDATE = "urn:ietf:params:netconf:capability:candidate:1.0";

// Secure
// Create a default SSL ChannelCredentials object
gNMISession::SecureChannelArguments input_args = get_channel_credentials();
gNMISession::gNMISession(const std::string& address, bool is_secure)
    : client(make_unique<gNMIClient>(grpc::CreateCustomChannel(address, input_args.channel_creds, *(input_args.args))))
{
    path::Repository repo;
    initialize(repo, address, is_secure);
    YLOG_DEBUG("Connected to {} using Secure Channel", address);
}

gNMISession::gNMISession(path::Repository & repo, const std::string& address, bool is_secure)
    : client(make_unique<gNMIClient>(grpc::CreateCustomChannel(address, input_args.channel_creds, *(input_args.args))))
{
    initialize(repo, address, is_secure);
    YLOG_DEBUG("Connected to {} using Secure Channel", address);
}

// Unsecure
gNMISession::gNMISession(const std::string& address)
    : client(make_unique<gNMIClient>(grpc::CreateChannel(address, grpc::InsecureChannelCredentials())))
{
    path::Repository repo;
    bool is_secure = false;
    YLOG_DEBUG("Connected to {}", address);
    initialize(repo, address, is_secure);
}

gNMISession::gNMISession(path::Repository & repo, const std::string& address)
    : client(make_unique<gNMIClient>(grpc::CreateChannel(address, grpc::InsecureChannelCredentials())))
{
    bool is_secure = false;
    YLOG_DEBUG("Connected to {}", address);
    initialize(repo, address, is_secure);
}

gNMISession::~gNMISession() = default;

void gNMISession::initialize(path::Repository & repo, const std::string& address, bool is_secure)
{
    IetfCapabilitiesParser capabilities_parser{};
    client->connect(address,is_secure);
    server_capabilities = client->get_capabilities();

    root_schema = repo.create_root_schema(capabilities_parser.parse(server_capabilities));

    if(root_schema.get() == nullptr)
    {
        YLOG_ERROR("Root schema cannot be obtained");
        throw(YCPPIllegalStateError{"Root schema cannot be obtained"});
    }
}

gNMISession::SecureChannelArguments get_channel_credentials()
{
    // Authenticate Server at Client
    string server_cert;
    ifstream rf("ems.pem");

    server_cert.assign((istreambuf_iterator<char>(rf)),(istreambuf_iterator<char>()));

    YLOG_DEBUG("In gnmi_session server cert: {}", server_cert);

    grpc::SslCredentialsOptions ssl_opts;
    auto args = std::make_shared<grpc::ChannelArguments>();
    gNMISession::SecureChannelArguments input_args;
    ssl_opts.pem_root_certs = server_cert;
    args->SetSslTargetNameOverride("ems.cisco.com");

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
        YLOG_ERROR("RPC is not supported");
        throw(YCPPOperationNotSupportedError{"RPC is not supported!"});
    }
    return datanode;
}

shared_ptr<path::DataNode> gNMISession::handle_edit(path::Rpc& ydk_rpc, const std::string & operation) const
{
    bool candidate_supported = is_candidate_supported(server_capabilities);
    auto gnmi_rpc = create_rpc_instance(*root_schema, "ietf-netconf:edit-config");
    auto & input = create_rpc_input(*gnmi_rpc);
    create_input_target(input, candidate_supported);
    create_input_error_option(input);
    string config_payload = get_config_payload(*root_schema, ydk_rpc);

    ly_verb(LY_LLSILENT); // turn off libyang logging at the beginning
    string gnmi_payload = get_gnmi_payload(input, "config", config_payload);
    ly_verb(LY_LLVRB); // enable libyang logging after payload has been created

    return handle_edit_reply(execute_payload(gnmi_payload, operation), *client, candidate_supported, operation);
}

static shared_ptr<path::DataNode> handle_edit_reply(const string & reply, gNMIClient & client, bool candidate_supported, const std::string & operation)
{
    if(reply.find("Success") == string::npos)
    {
        YLOG_ERROR("No ok in reply received from device");
        throw(YCPPServiceProviderError{reply});
    }

    // No error no output for edit-config
    return nullptr;
}


shared_ptr<path::DataNode> gNMISession::handle_read(path::Rpc& ydk_rpc, const std::string & operation) const
{
    //for now we only support crud rpc's
    bool config = is_config(ydk_rpc);
    auto gnmi_rpc = create_rpc_instance(*root_schema, get_read_rpc_name(config));
    auto & input = create_rpc_input(*gnmi_rpc);
    create_input_source(input, config);
    string filter_value = get_filter_payload(ydk_rpc);
    string gnmi_payload = get_gnmi_payload(input, "filter", filter_value);
    shared_ptr<path::DataNode> datanode = nullptr;
    return handle_read_reply(execute_payload(gnmi_payload, operation), *root_schema);
}

string gNMISession::execute_payload(const string & payload, const string & operation) const
{
    string reply = client->execute_wrapper(payload, operation);
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
        throw(YCPPServiceProviderError{reply});
    }
    data_start+= sizeof("\"data\":") - 1;
    auto data_end = reply.find_last_of("}");
    if(data_end == string::npos)
    {
        YLOG_ERROR( "No end data tag found in reply sent by device {}", reply.c_str());
        throw(YCPPError{"No end data tag found"});
    }

    string data = reply.substr(data_start, data_end-data_start + 1);

    auto datanode = shared_ptr<path::DataNode>(codec_service.decode(root_schema, data, EncodingFormat::JSON));
    if(!datanode)
    {
        YLOG_ERROR( "Codec service failed to decode datanode");
        throw(YCPPError{"Problems deserializing output"});
    }
    return datanode;
}

static void create_input_target(path::DataNode & input, bool candidate_supported)
{
    if(candidate_supported){
        input.create_datanode("target/candidate", "");
    }
    else {
        input.create_datanode("target/running", "");
    }
}

static bool is_candidate_supported(vector<string> capabilities)
{
    if(find(capabilities.begin(), capabilities.end(), TEMP_CANDIDATE) != capabilities.end()){
        //candidate is supported
        return true;
    }
    return false;
}

static void create_input_error_option(path::DataNode & input)
{
    input.create_datanode("error-option", "rollback-on-error");
}

static string get_config_payload(path::RootSchemaNode & root_schema,
    path::Rpc & rpc)
{
    path::Codec codec_service{};
    auto entity = rpc.get_input_node().find("entity");
    if(entity.empty()){
        YLOG_ERROR("Failed to get entity node");
        throw(YCPPInvalidArgumentError{"Failed to get entity node"});
    }

    path::DataNode* entity_node = entity[0].get();
    string entity_value = entity_node->get_value();

    //deserialize the entity_value
    auto datanode = codec_service.decode(root_schema, entity_value, EncodingFormat::JSON);

    if(!datanode){
        YLOG_ERROR("Failed to decode entity node");
        throw(YCPPInvalidArgumentError{"Failed to decode entity node"});
    }

    string config_payload {};

    for(auto const & child : datanode->get_children())
    {
        config_payload += codec_service.encode(*child, EncodingFormat::JSON, false);
    }
    return config_payload;
}

static string get_commit_rpc_payload()
{
    return "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
           "<commit/>"
           "</rpc>";
}

static bool is_config(path::Rpc & rpc)
{
    if(!rpc.get_input_node().find("only-config").empty())
    {
        return true;
    }
    return false;
}

static shared_ptr<path::Rpc> create_rpc_instance(path::RootSchemaNode & root_schema, string rpc_name)
{
    auto rpc = shared_ptr<path::Rpc>(root_schema.create_rpc(rpc_name));
    if(rpc == nullptr)
    {
        YLOG_ERROR("Cannot create payload for RPC: {}", rpc_name);
        throw(YCPPIllegalStateError{"Cannot create payload for RPC: "+ rpc_name});
    }
    return rpc;
}

static string get_read_rpc_name(bool config)
{
    if(config)
    {
        return "ietf-netconf:get-config";
    }
    return "ietf-netconf:get";
}

static path::DataNode& create_rpc_input(path::Rpc & gnmi_rpc)
{
    return gnmi_rpc.get_input_node();
}

static void create_input_source(path::DataNode & input, bool config)
{
    if(config)
    {
        input.create_datanode("source/running");
    }
}

static string get_filter_payload(path::Rpc & ydk_rpc)
{
    auto entity = ydk_rpc.get_input_node().find("filter");
    if(entity.empty())
    {
        YLOG_ERROR("Failed to get entity node.");
        throw(YCPPInvalidArgumentError{"Failed to get entity node"});
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
        throw(YCPPIllegalStateError{"CRUD read rpc schema not found!"});
    }
    return c[0];
}
}
}
