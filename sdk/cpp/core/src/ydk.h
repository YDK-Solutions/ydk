//////////////////////////////////////////////////////////////////
// @file ydk.h
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

#ifndef _YDK_H_
#define _YDK_H_

#ifdef __cplusplus
extern "C" {
#endif

typedef void* DataNode;
typedef void* Rpc;
typedef void* SchemaNode;
typedef void* RootSchemaNode;
typedef void* RootSchemaWrapper;
typedef void* Codec;
typedef void* ServiceProvider;
typedef void* OpenDaylightServiceProvider;
typedef void* Capability;
typedef void* Repository;

typedef struct DataNodeChildren
{
    DataNode* datanodes;
    int count;
} DataNodeChildren;

typedef int boolean;

typedef enum EncodingFormat
{
    XML   = 0,
    JSON
} EncodingFormat;

typedef enum LogLevel
{
    OFF    = 0,
    DEBUG,
    INFO,
    WARNING,
    ERROR
} LogLevel;

typedef enum Protocol
{
    Netconf = 0,
    Restconf
} Protocol;

Capability CapabilityCreate(const char* mod, const char* rev);
void CapabilityFree(Capability);

Repository RepositoryInitWithPath(const char*);
Repository RepositoryInit(void);
RootSchemaWrapper RepositoryCreateRootSchemaWrapper(Repository, const Capability caps[], int caps_size);
void RepositoryFree(Repository);

ServiceProvider NetconfServiceProviderInitWithRepo(Repository repo, const char * address, const char * username, const char * password, int port);
ServiceProvider NetconfServiceProviderInit(const char * address, const char * username, const char * password, int port);
RootSchemaNode ServiceProviderGetRootSchema(ServiceProvider);
EncodingFormat ServiceProviderGetEncoding(ServiceProvider);
void NetconfServiceProviderFree(ServiceProvider);

ServiceProvider RestconfServiceProviderInitWithRepo(Repository repo, const char * address, const char * username, const char * password, int port);
void RestconfServiceProviderFree(ServiceProvider);

OpenDaylightServiceProvider OpenDaylightServiceProviderInitWithRepo(Repository repo, const char * address, const char * username, const char * password, int port, EncodingFormat encoding, Protocol protocol);
void OpenDaylightServiceProviderFree(OpenDaylightServiceProvider);
ServiceProvider OpenDaylightServiceProviderGetNodeProvider(OpenDaylightServiceProvider provider, const char * node_id);
const char* OpenDaylightServiceProviderGetNodeIdByIndex(OpenDaylightServiceProvider provider, int idx);

Codec CodecInit(void);
void CodecFree(Codec);
const char* CodecEncode(Codec, DataNode, EncodingFormat, boolean);
DataNode CodecDecode(Codec, RootSchemaNode, const char*, EncodingFormat);

DataNode RootSchemaNodeCreate(RootSchemaNode, const char*);
Rpc RootSchemaNodeRpc(RootSchemaNode, const char*);
RootSchemaNode RootSchemaWrapperUnwrap(RootSchemaWrapper);

DataNode RpcInput(Rpc);
DataNode RpcExecute(Rpc, ServiceProvider);

DataNode DataNodeCreate(DataNode, const char*, const char*);
const char* DataNodeGetArgument(DataNode);
const char* DataNodeGetKeyword(DataNode);
const char* DataNodeGetPath(DataNode);
const char* DataNodeGetValue(DataNode);
DataNode DataNodeGetParent(DataNode);
void DataNodeAddAnnotation(DataNode, const char*);
DataNodeChildren DataNodeGetChildren(DataNode);
const char* DataNodeGetSegmentPath(DataNode);

void EnableLogging(LogLevel);
LogLevel GetLoggingLevel(void);

#ifdef __cplusplus
}
#endif

#endif /* _YDK_H_ */
