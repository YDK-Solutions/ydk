//////////////////////////////////////////////////////////////////
// @file golang.h
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

#ifndef _GOLANG_H_
#define _GOLANG_H_

#ifdef __cplusplus
extern "C" {
#endif

typedef void* DataNode;
typedef void* Rpc;
typedef void* SchemaNode;
typedef void* RootSchemaNode;
typedef void* CodecService;
typedef void* ServiceProvider;
typedef void* Capability;
typedef void* Repository;

typedef enum EncodingFormat
{
    XML   = 0,
    JSON
} EncodingFormat;

Repository RepositoryInitWithPath(const char*);
Repository RepositoryInit();
void RepositoryFree(Repository);

ServiceProvider NetconfServiceProviderInit(Repository repo, const char * address, const char * username, const char * password, int port);
void NetconfServiceProviderFree(ServiceProvider);
RootSchemaNode NetconfServiceProviderGetRootSchema(ServiceProvider);

CodecService CodecServiceInit(void);
void CodecServiceFree(CodecService);
const char* CodecServiceEncode(CodecService, DataNode, EncodingFormat, int);
DataNode CodecServiceDecode(CodecService, RootSchemaNode, const char*, EncodingFormat);

DataNode RootSchemaNodeCreate(RootSchemaNode, const char*);
Rpc RootSchemaNodeRpc(RootSchemaNode, const char*);

DataNode RpcInput(Rpc);
DataNode RpcExecute(Rpc, ServiceProvider);

DataNode DataNodeCreate(DataNode, const char*, const char*);

void EnableLogging(void);

#ifdef __cplusplus
}
#endif

#endif /* _GOLANG_H_ */
