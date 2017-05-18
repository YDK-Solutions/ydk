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

import (
    "github.com/CiscoDevNet/ydk-go/ydk/path"
    "github.com/CiscoDevNet/ydk-go/ydk/types"
)

type CrudService struct {
}

func (c *CrudService) Create(provider types.ServiceProvider, entity types.Entity) bool {
    return operationSucceeded( path.ExecuteRpc(provider, entity, "ydk:create", "entity", false) )
}

func (c *CrudService) Update(provider types.ServiceProvider, entity types.Entity) bool {
    return operationSucceeded( path.ExecuteRpc(provider, entity, "ydk:update", "entity", false) )
}

func (c *CrudService) Delete(provider types.ServiceProvider, entity types.Entity) bool {
    return operationSucceeded( path.ExecuteRpc(provider, entity, "ydk:delete", "entity", false) )
}

func (c *CrudService) Read(provider types.ServiceProvider, filter types.Entity) types.Entity {
    return path.ReadDatanode( filter, path.ExecuteRpc(provider, filter, "ydk:read", "filter", true) )
}

func (c *CrudService) ReadConfig(provider types.ServiceProvider, filter types.Entity) types.Entity {
    return path.ReadDatanode( filter, path.ExecuteRpc(provider, filter, "ydk:read", "filter", false) )
}

func operationSucceeded(node types.DataNode) bool {
    return node.Private != nil
}

type CodecService struct {
}

func (c *CodecService) Encode(provider types.CodecServiceProvider, entity types.Entity) string {
    return ""
}

func (c *CodecService) Decode(provider types.CodecServiceProvider, payload string) types.Entity {
    return nil
}
