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
 *   http:www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *----------------------------------------------
 */
package providers

import (
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk/path"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
)

type NetconfServiceProvider struct {
	Repo     types.Repository
	Address  string
	Username string
	Password string
	Port     int

	Private types.CServiceProvider
}

type RestconfServiceProvider struct {
	Path     string
	Address  string
	Username string
	Password string
	Port     int

	Private types.CServiceProvider
}

func (provider *NetconfServiceProvider) GetPrivate() interface{} {
	return provider.Private
}

func (provider *NetconfServiceProvider) Connect() {
	provider.Private = path.ConnectToProvider(provider.Repo, provider.Address, provider.Username, provider.Password, provider.Port)
}

func (provider *NetconfServiceProvider) Disconnect() {
	if provider.Private.Private == nil {
		return
	}
	path.DisconnectFromProvider(provider.Private)
}

func (provider *RestconfServiceProvider) GetPrivate() interface{} {
	return provider.Private
}

func (provider *RestconfServiceProvider) Connect() {
	provider.Private = path.ConnectToRestconfProvider(provider.Path, provider.Address, provider.Username, provider.Password, provider.Port)
}

func (provider *RestconfServiceProvider) Disconnect() {
	if provider.Private.Private == nil {
		return
	}
	path.DisconnectFromRestconfProvider(provider.Private)
}

type CodecServiceProvider struct {
	Repo     types.Repository
	Encoding types.EncodingFormat

	RootSchemaTable map[string]types.RootSchemaNode
}

func (provider *CodecServiceProvider) Initialize(entity types.Entity) {
	bundle_name := entity.GetBundleName()
	if len(provider.RootSchemaTable) == 0 {
		provider.RootSchemaTable = make(map[string]types.RootSchemaNode)
	}
	_, ok := provider.RootSchemaTable[bundle_name]
	if !ok {
		fmt.Printf("CodecServiceProvider initialize with %v bundle\n", bundle_name)
		root_schema_node := path.InitCodecServiceProvider(entity, provider.Repo)
		provider.RootSchemaTable[bundle_name] = root_schema_node
	}
}

func (provider *CodecServiceProvider) GetEncoding() types.EncodingFormat {
	return provider.Encoding
}

func (provider *CodecServiceProvider) GetRootSchemaNode(entity types.Entity) types.RootSchemaNode {
	root_schema_node, ok := provider.RootSchemaTable[entity.GetBundleName()]
	if !ok {
		panic("Root schema node not found in provider!")
	}
	return root_schema_node
}
