//////////////////////////////////////////////////////////////////
// @file golang.cpp
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

#include <ydk/path_api.hpp>
#include <ydk/types.hpp>
#include <ydk/netconf_provider.hpp>
#include <iostream>
#include <spdlog/spdlog.h>

#include "golang.h"

using namespace std;

typedef struct RpcWrapper
{
    RpcWrapper(shared_ptr<ydk::path::Rpc> rpc)
        : private_rpc(rpc)
    {
    }

    shared_ptr<ydk::path::Rpc> private_rpc;
} RpcWrapper;

typedef struct DataNodeWrapper
{
    DataNodeWrapper(shared_ptr<ydk::path::DataNode> data)
        : private_datanode(data)
    {
    }

    shared_ptr<ydk::path::DataNode> private_datanode;
} DataNodeWrapper;

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

Repository RepositoryInitWithPath(const char* path)
{
    ydk::path::Repository* real_repo = new ydk::path::Repository(path);
    return (void*)real_repo;
}

Repository RepositoryInit()
{
    ydk::path::Repository* real_repo = new ydk::path::Repository();
    return (void*)real_repo;
}

void RepositoryFree(Repository repo)
{
    ydk::path::Repository* real_repo = new ydk::path::Repository();
    delete real_repo;
}

NetconfServiceProvider NetconfServiceProviderInit(Repository repo, const char * address, const char * username, const char * password, int port)
{
    ydk::path::Repository* real_repo = (ydk::path::Repository*)repo;
    ydk::NetconfServiceProvider * real_provider = new ydk::NetconfServiceProvider(*real_repo, address, username, password, port);
    return (void*)real_provider;
}

void NetconfServiceProviderFree(NetconfServiceProvider provider)
{
    ydk::NetconfServiceProvider * real_provider = (ydk::NetconfServiceProvider*)provider;
    delete real_provider;
}

RootSchemaNode NetconfServiceProviderGetRootSchema(NetconfServiceProvider provider)
{
    ydk::NetconfServiceProvider * real_provider = (ydk::NetconfServiceProvider*)provider;
    ydk::path::RootSchemaNode* root_schema = &real_provider->get_root_schema();
    return (void*)root_schema;
}


CodecService CodecServiceInit(void)
{
    ydk::path::CodecService * codec = new ydk::path::CodecService();
    return (void*)codec;
}

void CodecServiceFree(CodecService codec)
{
    ydk::path::CodecService * real_codec = (ydk::path::CodecService*)codec;
    delete real_codec;
}

const char* CodecServiceEncode(CodecService codec, DataNode datanode, EncodingFormat encoding, int pretty)
{
    ydk::path::CodecService * real_codec = (ydk::path::CodecService*)codec;
    ydk::path::DataNode * real_datanode = (ydk::path::DataNode *)datanode;
    std::string payload = real_codec->encode(*real_datanode, get_real_encoding(encoding), pretty);
    return payload.c_str();
}

DataNode CodecServiceDecode(CodecService codec, RootSchemaNode root_schema, const char* payload, EncodingFormat encoding)
{
    ydk::path::CodecService * real_codec = (ydk::path::CodecService*)codec;
    ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
    shared_ptr<ydk::path::DataNode> datanode = real_codec->decode(*real_root_schema, payload, get_real_encoding(encoding));
    DataNodeWrapper * datanode_wrapper = new DataNodeWrapper(datanode);
    return (void*)datanode_wrapper;
}

DataNode RootSchemaNodeCreate(RootSchemaNode root_schema, const char* path)
{
    ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
    ydk::path::DataNode * datanode = &real_root_schema->create(path);
    return (void*)datanode;
}

Rpc RootSchemaNodeRpc(RootSchemaNode root_schema, const char* path)
{
    ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
    shared_ptr<ydk::path::Rpc> rpc = real_root_schema->rpc(path);
    RpcWrapper * rpc_wrapper = new RpcWrapper(rpc);
    return (void*)rpc_wrapper;
}

DataNode RpcInput(Rpc rpc)
{
    RpcWrapper* real_rpc = (RpcWrapper*)rpc;
    ydk::path::DataNode * input = &real_rpc->private_rpc->input();
    return (void*)input;
}

DataNode RpcExecute(Rpc rpc, ServiceProvider provider)
{
    RpcWrapper* real_rpc = (RpcWrapper*)rpc;
    ydk::path::ServiceProvider * real_provider = (ydk::path::ServiceProvider *) provider;
    std::shared_ptr<ydk::path::DataNode> result = (*real_rpc->private_rpc)(*real_provider);
    DataNodeWrapper * datanode_wrapper = new DataNodeWrapper(result);
    return (void*)datanode_wrapper;
}

DataNode DataNodeCreate(DataNode datanode, const char* path, const char* value)
{
    ydk::path::DataNode * real_datanode = (ydk::path::DataNode *)datanode;
    ydk::path::DataNode * result = &real_datanode->create(path, value);
    return (void*)result;
}

void EnableLogging(void)
{
    auto console = spdlog::stdout_color_mt("ydk");
    console->set_level(spdlog::level::info);
}
