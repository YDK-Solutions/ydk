//////////////////////////////////////////////////////////////////
// @file ydk.cpp
//
// YANG Development Kit
// Copyright 2017 Cisco Systems. All rights reserved
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

#include "crud_service.hpp"
#include "netconf_provider.hpp"
#include "opendaylight_provider.hpp"
#include "restconf_provider.hpp"
#include "path_api.hpp"
#include "path/path_private.hpp"
#include "types.hpp"

#include <iostream>
#include <spdlog/spdlog.h>

#include "ydk.h"

using namespace std;

//////////////////////////////////////////////////////////////////////////
// structs
//////////////////////////////////////////////////////////////////////////
typedef struct RpcWrapper
{
    RpcWrapper(shared_ptr<ydk::path::Rpc> rpc)
        : priv(rpc)
    {
    }

    shared_ptr<ydk::path::Rpc> priv;
} RpcWrapper;

typedef struct DataNodeWrapper
{
    DataNodeWrapper(shared_ptr<ydk::path::DataNode> data)
        : priv(data)
    {
    }

    shared_ptr<ydk::path::DataNode> priv;
} DataNodeWrapper;

typedef struct RootSchemaNodeWrapper
{
    RootSchemaNodeWrapper(shared_ptr<ydk::path::RootSchemaNode> node)
        :priv(node)
    {
    }
    shared_ptr<ydk::path::RootSchemaNode> priv;
} RootSchemaNodeWrapper;

//////////////////////////////////////////////////////////////////////////
// Utility functions
//////////////////////////////////////////////////////////////////////////
static DataNodeWrapper* wrap(ydk::path::DataNode* datanode)
{
    return (new DataNodeWrapper(shared_ptr<ydk::path::DataNode>(datanode)));
}

static DataNodeWrapper* wrap(shared_ptr<ydk::path::DataNode> datanode)
{
    return (new DataNodeWrapper(datanode));
}

static ydk::path::DataNode* unwrap(DataNodeWrapper* datanode_wrapper)
{
    return datanode_wrapper->priv.get();
}

static RootSchemaNodeWrapper* wrap(ydk::path::RootSchemaNode* node)
{
    return (new RootSchemaNodeWrapper(shared_ptr<ydk::path::RootSchemaNode>(node)));
}

static RootSchemaNodeWrapper* wrap(shared_ptr<ydk::path::RootSchemaNode> node)
{
    return (new RootSchemaNodeWrapper(node));
}

static ydk::path::RootSchemaNode* unwrap(RootSchemaNodeWrapper* node_wrapper)
{
    return node_wrapper->priv.get();
}

static RpcWrapper* wrap(shared_ptr<ydk::path::Rpc> rpc)
{
    return (new RpcWrapper(rpc));
}

static ydk::path::Rpc* unwrap(RpcWrapper* rpc_wrapper)
{
    return rpc_wrapper->priv.get();
}

static ydk::EncodingFormat get_real_encoding(EncodingFormat encoding)
{
    switch(encoding)
    {
        case XML:
            return ydk::EncodingFormat::XML;
        case JSON:
                return ydk::EncodingFormat::JSON;
    }
}

static const char* string_to_array(const string & str)
{
    char * cstr = new char [str.length()+1];
    std::strcpy (cstr, str.c_str());
    return cstr;
}

//////////////////////////////////////////////////////////////////////////
// Exported functions
//////////////////////////////////////////////////////////////////////////

Repository RepositoryInitWithPath(const char* path)
{
    ydk::path::Repository* real_repo = new ydk::path::Repository(path);
    return static_cast<void*>(real_repo);
}

Repository RepositoryInit()
{
    ydk::path::Repository* real_repo = new ydk::path::Repository();
    return static_cast<void*>(real_repo);
}

RootSchemaWrapper RepositoryCreateRootSchemaWrapper(Repository repo, const Capability caps[], int caps_size) {
    try
    {
        std::vector<ydk::path::Capability> real_caps;
        for(int i = 0; i < caps_size; i++) {
            real_caps.push_back(*static_cast<ydk::path::Capability*>(caps[i]));
        }
        ydk::path::Repository* real_repo = static_cast<ydk::path::Repository*>(repo);
        std::shared_ptr<ydk::path::RootSchemaNode> schema_node = real_repo->create_root_schema(real_caps);
        return static_cast<void*>(wrap(schema_node));
    }
    catch(...)
    {
        return NULL;
    }
}

void RepositoryFree(Repository repo)
{
    ydk::path::Repository* real_repo = static_cast<ydk::path::Repository*>(repo);
    if(real_repo != NULL)
    {
        delete real_repo;
    }
}

Capability CapabilityCreate(const char* mod, const char* rev) {
    try
    {
        ydk::path::Capability * real_cap = new ydk::path::Capability(mod, rev);
        return static_cast<void*>(real_cap);
    }
    catch(...)
    {
        return NULL;
    }
}

void CapabilityFree(Capability cap)
{
    ydk::path::Capability * real_cap = static_cast<ydk::path::Capability*>(cap);
    if (real_cap != NULL)
    {
        delete real_cap;
    }
}

ServiceProvider NetconfServiceProviderInit(const char * address, const char * username, const char * password, int port)
{
    try
    {
        ydk::NetconfServiceProvider * real_provider = new ydk::NetconfServiceProvider(address, username, password, port);
        return static_cast<void*>(real_provider);
    }
    catch(...)
    {
        return NULL;
    }
}

ServiceProvider NetconfServiceProviderInitWithRepo(Repository repo, const char * address, const char * username, const char * password, int port)
{
    try
    {
        ydk::path::Repository* real_repo = static_cast<ydk::path::Repository*>(repo);
        ydk::NetconfServiceProvider * real_provider = new ydk::NetconfServiceProvider(*real_repo, address, username, password, port);
        return static_cast<void*>(real_provider);
    }
    catch(...)
    {
        return NULL;
    }
}

void NetconfServiceProviderFree(ServiceProvider provider)
{
    ydk::NetconfServiceProvider * real_provider = static_cast<ydk::NetconfServiceProvider*>(provider);
    if(real_provider != NULL)
    {
        delete real_provider;
    }
}

ServiceProvider RestconfServiceProviderInitWithRepo(Repository repo, const char * address, const char * username, const char * password, int port)
{
    try
    {
        ydk::path::Repository* real_repo = static_cast<ydk::path::Repository*>(repo);
        ydk::RestconfServiceProvider * real_provider = new ydk::RestconfServiceProvider(*real_repo, address, username, password, port);
        return static_cast<void*>(real_provider);
    }
    catch(...)
    {
        return NULL;
    }
}

void RestconfServiceProviderFree(ServiceProvider provider)
{
    ydk::RestconfServiceProvider * real_provider = static_cast<ydk::RestconfServiceProvider*>(provider);
    if(real_provider != NULL)
    {
        delete real_provider;
    }
}

OpenDaylightServiceProvider OpenDaylightServiceProviderInitWithRepo(Repository repo, const char * address, const char * username, const char * password, int port, EncodingFormat encoding, Protocol protocol)
{
    try
    {
        ydk::path::Repository* real_repo = static_cast<ydk::path::Repository*>(repo);
        ydk::EncodingFormat real_encoding;
        if (encoding == XML) {
            real_encoding = ydk::EncodingFormat::XML;
        }
        else {
            real_encoding = ydk::EncodingFormat::JSON;
        }
        ydk::Protocol real_protocol;
        if (protocol == Netconf) {
            real_protocol = ydk::Protocol::netconf;
        }
        else {
            real_protocol = ydk::Protocol::restconf;
        }
        ydk::OpenDaylightServiceProvider * real_provider = new ydk::OpenDaylightServiceProvider(*real_repo, address, username, password, port, real_encoding, real_protocol);
        return static_cast<void*>(real_provider);
    }
    catch(...)
    {
        return NULL;
    }
}

void OpenDaylightServiceProviderFree(OpenDaylightServiceProvider provider)
{
    ydk::OpenDaylightServiceProvider * real_provider = static_cast<ydk::OpenDaylightServiceProvider*>(provider);
    if(real_provider != NULL)
    {
        delete real_provider;
    }
}

ServiceProvider OpenDaylightServiceProviderGetNodeProvider(OpenDaylightServiceProvider provider, const char * node_id)
{
    ydk::OpenDaylightServiceProvider * real_provider = static_cast<ydk::OpenDaylightServiceProvider*>(provider);
    ydk::path::ServiceProvider * node_provider = &real_provider->get_node_provider(string(node_id));
    return static_cast<void*>(node_provider);
}

const char* OpenDaylightServiceProviderGetNodeIDByIndex(OpenDaylightServiceProvider provider, int idx)
{
    ydk::OpenDaylightServiceProvider * real_provider = static_cast<ydk::OpenDaylightServiceProvider*>(provider);
    if ((size_t)idx < real_provider->get_node_ids().size())
    {
        return string_to_array(real_provider->get_node_ids()[idx]);
    }
    else
    {
        return NULL;
    }
}

RootSchemaNode ServiceProviderGetRootSchema(ServiceProvider provider)
{
    try
    {
    ydk::path::ServiceProvider * real_provider = static_cast<ydk::path::ServiceProvider*>(provider);
    ydk::path::RootSchemaNode* root_schema = &real_provider->get_root_schema();
    return static_cast<void*>(root_schema);
    }
    catch(...)
    {
        return NULL;
    }
}

EncodingFormat ServiceProviderGetEncoding(ServiceProvider provider)
{
    ydk::path::ServiceProvider * real_provider = static_cast<ydk::path::ServiceProvider*>(provider);
    auto encoding = real_provider->get_encoding();
    if (encoding == ydk::EncodingFormat::XML) {
        return XML;
    }
    else {
        return JSON;
    }
}

Codec CodecInit(void)
{
    ydk::path::Codec * codec = new ydk::path::Codec();
    return static_cast<void*>(codec);
}

void CodecFree(Codec codec)
{
    ydk::path::Codec * real_codec = (ydk::path::Codec*)codec;
    if(real_codec != NULL)
    {
        delete real_codec;
    }
}

const char* CodecEncode(Codec codec, DataNode datanode, EncodingFormat encoding, boolean pretty)
{
    try
    {
        ydk::path::Codec * real_codec = (ydk::path::Codec*)codec;
        DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
        ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);

        std::string payload = real_codec->encode(*real_datanode, get_real_encoding(encoding), pretty);
        return string_to_array(payload);
    }
    catch(...)
    {
        return NULL;
    }
}

DataNode CodecDecode(Codec codec, RootSchemaNode root_schema, const char* payload, EncodingFormat encoding)
{
    try
    {
        ydk::path::Codec * real_codec = (ydk::path::Codec*)codec;
        ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
        shared_ptr<ydk::path::DataNode> datanode = real_codec->decode(*real_root_schema, payload, get_real_encoding(encoding));

        return static_cast<void*>(wrap(datanode));
    }
    catch(...)
    {
        return NULL;
    }
}

DataNode RootSchemaNodeCreate(RootSchemaNode root_schema, const char* path)
{
    try
    {
        ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
        ydk::path::DataNode * datanode = &real_root_schema->create(path);

        return static_cast<void*>(wrap(datanode));
    }
    catch(...)
    {
        return NULL;
    }
}

RootSchemaNode RootSchemaWrapperUnwrap(RootSchemaWrapper wrapper)
{
    RootSchemaNodeWrapper* real_wrapper = (RootSchemaNodeWrapper*) wrapper;
    return static_cast<void*>(unwrap(real_wrapper));
}

Rpc RootSchemaNodeRpc(RootSchemaNode root_schema, const char* path)
{
    try
    {
        ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
        shared_ptr<ydk::path::Rpc> rpc = real_root_schema->rpc(path);

        return static_cast<void*>(wrap(rpc));
    }
    catch(...)
    {
        return NULL;
    }
}

DataNode RpcInput(Rpc rpc)
{
    try
    {
        RpcWrapper* rpc_wrapper = (RpcWrapper*)rpc;
        ydk::path::Rpc* real_rpc = unwrap(rpc_wrapper);

        ydk::path::DataNode * input = &real_rpc->input();

        return static_cast<void*>(wrap(input));
    }
    catch(...)
    {
        return NULL;
    }
}

DataNode RpcExecute(Rpc rpc, ServiceProvider provider)
{
    try
    {
        RpcWrapper* rpc_wrapper = (RpcWrapper*)rpc;
        ydk::path::Rpc* real_rpc = unwrap(rpc_wrapper);

        ydk::path::ServiceProvider * real_provider = (ydk::path::ServiceProvider *) provider;
        std::shared_ptr<ydk::path::DataNode> result = (*real_rpc)(*real_provider);

        return static_cast<void*>(wrap(result));
    }
    catch(...)
    {
        return NULL;
    }
}

DataNode DataNodeCreate(DataNode datanode, const char* path, const char* value)
{
    try
    {
        DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
        ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
        ydk::path::DataNode * result = &real_datanode->create(path, value);

        return static_cast<void*>(wrap(result));
    }
    catch(...)
    {
        return NULL;
    }
}

const char* DataNodeGetArgument(DataNode datanode)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    string s = real_datanode->schema().statement().arg;
    return string_to_array(s);
}

const char* DataNodeGetKeyword(DataNode datanode)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    string s = real_datanode->schema().statement().keyword;
    return string_to_array(s);
}

const char* DataNodeGetPath(DataNode datanode)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    string s = real_datanode->path();
    return string_to_array(s);
}

const char* DataNodeGetValue(DataNode datanode)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    string s = real_datanode->get();
    return string_to_array(s);
}

DataNodeChildren DataNodeGetChildren(DataNode datanode)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);

    std::vector<shared_ptr<ydk::path::DataNode>> children = real_datanode->children();
    DataNode* child_array = new DataNode[children.size()];
    for(size_t i=0; i < children.size(); i++)
    {
        child_array[i] = wrap(children[i]);
    }
    return {child_array, static_cast<int>(children.size())};
}

void DataNodeAddAnnotation(DataNode datanode, const char* operation)
{
    string oper{operation};
    if(oper.size() == 0)
    {
        return;
    }
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    ydk::path::Annotation annotation{"ietf-netconf", "operation", operation};
    real_datanode->add_annotation(annotation);
}

DataNode DataNodeGetParent(DataNode datanode)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    ydk::path::DataNode* parent = real_datanode->parent();
    return static_cast<void*>(wrap(parent));
}

const char* DataNodeGetSegmentPath(DataNode datanode)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    string path = real_datanode->path();
    std::vector<std::string> segments = ydk::path::segmentalize(path);
    return string_to_array(segments.back());
}

void EnableLogging(LogLevel level)
{
    auto console = spdlog::stdout_color_mt("ydk");
    switch(level)
    {
        case DEBUG:
            console->set_level(spdlog::level::debug);
            return;

        case INFO:
            console->set_level(spdlog::level::info);
            return;

        case WARNING:
            console->set_level(spdlog::level::warn);
            return;

        case ERROR:
            console->set_level(spdlog::level::err);
            return;

        case OFF:
            console->set_level(spdlog::level::off);
            return;
    }
}

LogLevel GetLoggingLevel(void)
{
    auto console = spdlog::get("ydk");
    if(console == nullptr)
    {
        return OFF;
    }

    auto level = console->level();
    switch(level)
    {
        case spdlog::level::off:
            return OFF;
        case spdlog::level::trace:
        case spdlog::level::debug:
            return DEBUG;

        case spdlog::level::info:
            return INFO;

        case spdlog::level::warn:
            return WARNING;

        case spdlog::level::critical:
        case spdlog::level::err:
            return ERROR;
    }
}
