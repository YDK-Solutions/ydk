/*
 * ------------------------------------------------------------------
 * YANG Development Kit
 * Copyright 2017 Cisco Systems. All rights reserved
 *
 *----------------------------------------------
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http:*www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *----------------------------------------------
 */
package services

// #cgo CXXFLAGS: -g -std=c++11
// #cgo LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #include <ydk/ydk.h>
import "C"

import (
    "github.com/CiscoDevNet/ydk-go/providers"
    "github.com/CiscoDevNet/ydk-go/types"
    "unsafe"
)

type CrudService struct {
}

func (c *CrudService) Create(provider providers.ServiceProvider, entity types.Entity) bool {

    return operationSucceeded( executeRpc(provider, entity, "ydk:create", "entity", false) )
}

func (c *CrudService) Update(provider providers.ServiceProvider, entity types.Entity) bool {
    return operationSucceeded( executeRpc(provider, entity, "ydk:update", "entity", false) )
}

func (c *CrudService) Delete(provider providers.ServiceProvider, entity types.Entity) bool {
    return operationSucceeded( executeRpc(provider, entity, "ydk:delete", "entity", false) )
}


func (c *CrudService) Read(provider providers.ServiceProvider, filter types.Entity) types.Entity {
    return readDatanode( filter, executeRpc(provider, filter, "ydk:read", "filter", true) )
}

func (c *CrudService) ReadConfig(provider providers.ServiceProvider, filter types.Entity) types.Entity {
    return readDatanode( filter, executeRpc(provider, filter, "ydk:read", "filter", false) )
}

func readDatanode(filter types.Entity, read_data_node C.DataNode) types.Entity {
    if (read_data_node == nil) {
        return nil
    }

    top_entity := getTopEntityFromFilter(filter);
    children := C.DataNodeGetChildren(read_data_node)

    types.GetEntityFromDataNode(children, top_entity)
    return top_entity
}

func operationSucceeded(node C.DataNode) bool {
    return node == nil
}

func getTopEntityFromFilter(filter types.Entity) types.Entity {
    if(filter.GetParent() == nil) {
        return filter
    }

    return getTopEntityFromFilter(filter.GetParent())
}

func executeRpc(provider C.ServiceProvider, entity types.Entity, operation string, data_tag string, set_config_flag bool) C.DataNode {
    root_schema := C.ServiceProviderGetRootSchema(provider)

    ydk_rpc := C.RootSchemaNodeRpc(root_schema, C.CString(operation))
    data := getDataPayload(entity, provider)
    input := C.RpcInput(ydk_rpc)

    if(set_config_flag) {
        C.DataNodeCreate(input, C.CString("only-config"), C.CString(""))
    }

    C.DataNodeCreate(input, C.CString(data_tag), C.CString(data))
    return C.RpcExecute(ydk_rpc, provider)
}

func getDataPayload(entity types.Entity, provider C.ServiceProvider) string {
    root_schema := C.ServiceProviderGetRootSchema(provider)
    datanode := types.GetDataNodeFromEntity(entity, root_schema)

    for datanode!= nil && datanode.GetParent() != nil {
        datanode = datanode.GetParent()
    }

    codec := C.CodecServiceInit()
    defer C.CodecServiceFree(codec)
    var data *C.char = C.CodecServiceEncode(codec, datanode, C.XML, 1)
    defer C.free(unsafe.Pointer(data))

    return data
}
