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

DataNodeWrapper* wrap(ydk::path::DataNode* datanode)
{
    return (new DataNodeWrapper(shared_ptr<ydk::path::DataNode>(datanode)));
}

DataNodeWrapper* wrap(shared_ptr<ydk::path::DataNode> datanode)
{
    return (new DataNodeWrapper(datanode));
}

ydk::path::DataNode* unwrap(DataNodeWrapper* datanode_wrapper)
{
    return datanode_wrapper->priv.get();
}

RpcWrapper* wrap(ydk::path::Rpc* rpc)
{
    return (new RpcWrapper(shared_ptr<ydk::path::Rpc>(rpc)));
}

RpcWrapper* wrap(shared_ptr<ydk::path::Rpc> rpc)
{
    return (new RpcWrapper(rpc));
}

ydk::path::Rpc* unwrap(RpcWrapper* rpc_wrapper)
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
    ydk::path::Repository* real_repo = (ydk::path::Repository*)repo;
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
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);

    std::string payload = real_codec->encode(*real_datanode, get_real_encoding(encoding), pretty);
    char * cstr = new char [payload.length()+1];
    std::strcpy (cstr, payload.c_str());
    return cstr;
}

DataNode CodecServiceDecode(CodecService codec, RootSchemaNode root_schema, const char* payload, EncodingFormat encoding)
{
    ydk::path::CodecService * real_codec = (ydk::path::CodecService*)codec;
    ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
    shared_ptr<ydk::path::DataNode> datanode = real_codec->decode(*real_root_schema, payload, get_real_encoding(encoding));

    return (void*)(wrap(datanode));
}

DataNode RootSchemaNodeCreate(RootSchemaNode root_schema, const char* path)
{
    ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
    ydk::path::DataNode * datanode = &real_root_schema->create(path);

    return (void*)(wrap(datanode));
}

Rpc RootSchemaNodeRpc(RootSchemaNode root_schema, const char* path)
{
    ydk::path::RootSchemaNode * real_root_schema = (ydk::path::RootSchemaNode*)root_schema;
    shared_ptr<ydk::path::Rpc> rpc = real_root_schema->rpc(path);

    return (void*)(wrap(rpc));
}

DataNode RpcInput(Rpc rpc)
{
    RpcWrapper* rpc_wrapper = (RpcWrapper*)rpc;
    ydk::path::Rpc* real_rpc = unwrap(rpc_wrapper);

    ydk::path::DataNode * input = &real_rpc->input();

    return (void*)(wrap(input));
}

DataNode RpcExecute(Rpc rpc, ServiceProvider provider)
{
    RpcWrapper* rpc_wrapper = (RpcWrapper*)rpc;
    ydk::path::Rpc* real_rpc = unwrap(rpc_wrapper);

    ydk::path::ServiceProvider * real_provider = (ydk::path::ServiceProvider *) provider;
    std::shared_ptr<ydk::path::DataNode> result = (*real_rpc)(*real_provider);

    return (void*)(wrap(result));
}

DataNode DataNodeCreate(DataNode datanode, const char* path, const char* value)
{
    DataNodeWrapper* datanode_wrapper = (DataNodeWrapper*)datanode;
    ydk::path::DataNode* real_datanode = unwrap(datanode_wrapper);
    ydk::path::DataNode * result = &real_datanode->create(path, value);

    return (void*)(wrap(result));
}

void EnableLogging(void)
{
    auto console = spdlog::stdout_color_mt("ydk");
    console->set_level(spdlog::level::info);
}
