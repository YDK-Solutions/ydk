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

// OpenDaylightServiceProvider A service provider to be used to communicate with an OpenDaylight instance: https://www.opendaylight.org
type OpenDaylightServiceProvider struct {
	Path           string
	Address        string
	Username       string
	Password       string
	Port           int
	EncodingFormat types.EncodingFormat
	Protocol       types.Protocol

	Private types.COpenDaylightServiceProvider
	// keep alive
	ProvidersHolder []types.ServiceProvider
}

// NetconfServiceProvider Implementation of ServiceProvider for the NETCONF protocol: https://tools.ietf.org/html/rfc6241
type NetconfServiceProvider struct {
	Repo     types.Repository
	Address  string
	Username string
	Password string
	Port     int

	Private types.CServiceProvider
}

// RestconfServiceProvider Implementation of ServiceProvider for the RESTCONF protocol: https://tools.ietf.org/html/draft-ietf-netconf-restconf-18
type RestconfServiceProvider struct {
	Path     string
	Address  string
	Username string
	Password string
	Port     int

	Private types.CServiceProvider
}

// GetPrivate returns private pointer for OpenDaylightServiceProvider
func (provider *OpenDaylightServiceProvider) GetPrivate() interface{} {
	return provider.Private
}

// Connect to OpenDaylightServiceProvider using Path/Address/Username/Password/Port
func (provider *OpenDaylightServiceProvider) Connect() {
	provider.Private = path.ConnectToOpenDaylightProvider(provider.Path, provider.Address, provider.Username, provider.Password, provider.Port, provider.EncodingFormat, provider.Protocol)
}

// GetNodeIDs returns OpenDaylightServiceProvider Node IDs
func (provider *OpenDaylightServiceProvider) GetNodeIDs() []string {
	return path.OpenDaylightServiceProviderGetNodeIDs(provider.Private)
}

// GetNodeProvider returns Node provider by ID
func (provider *OpenDaylightServiceProvider) GetNodeProvider(nodeID string) types.ServiceProvider {
	p := path.OpenDaylightServiceProviderGetNodeProvider(provider.Private, nodeID)
	if provider.Protocol == types.Restconf {
		nodeProvider := RestconfServiceProvider{Path: provider.Path, Address: provider.Address, Password: provider.Password, Username: provider.Username, Port: provider.Port}
		nodeProvider.Private = p
		provider.ProvidersHolder = append(provider.ProvidersHolder, &nodeProvider)
		return &nodeProvider
	}
	repo := types.Repository{}
	repo.Path = provider.Path
	nodeProvider := NetconfServiceProvider{Repo: repo, Address: provider.Address, Password: provider.Password, Username: provider.Username, Port: provider.Port}
	nodeProvider.Private = p
	provider.ProvidersHolder = append(provider.ProvidersHolder, &nodeProvider)
	return &nodeProvider
}

// Disconnect from OpenDaylightServiceProvider
func (provider *OpenDaylightServiceProvider) Disconnect() {
	if provider.Private.Private == nil {
		return
	}
	path.DisconnectFromOpenDaylightProvider(provider.Private)
}

// GetPrivate returns private pointer for NetconfServiceProvider
func (provider *NetconfServiceProvider) GetPrivate() interface{} {
	return provider.Private
}

// Connect to NetconfServiceProvider using Repo/Address/Username/Password/Port
func (provider *NetconfServiceProvider) Connect() {
	provider.Private = path.ConnectToProvider(provider.Repo, provider.Address, provider.Username, provider.Password, provider.Port)
}

// Disconnect from NetconfServiceProvider
func (provider *NetconfServiceProvider) Disconnect() {
	if provider.Private.Private == nil {
		return
	}
	path.DisconnectFromProvider(provider.Private)
}

// GetPrivate returns private pointer for RestconfServiceProvider
func (provider *RestconfServiceProvider) GetPrivate() interface{} {
	return provider.Private
}

// Connect to RestconfServiceProvider using Path/Address/Username/Password/Port
func (provider *RestconfServiceProvider) Connect() {
	provider.Private = path.ConnectToRestconfProvider(provider.Path, provider.Address, provider.Username, provider.Password, provider.Port)
}

// Disconnect from RestconfServiceProvider
func (provider *RestconfServiceProvider) Disconnect() {
	if provider.Private.Private == nil {
		return
	}
	path.DisconnectFromRestconfProvider(provider.Private)
}

// CodecServiceProvider Encode and decode to XML/JSON format
type CodecServiceProvider struct {
	Repo     types.Repository
	Encoding types.EncodingFormat

	RootSchemaTable map[string]types.RootSchemaNode
}

// Initialize CodecServiceProvider
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

// GetEncoding returns encoding format for CodecServiceProvider
func (provider *CodecServiceProvider) GetEncoding() types.EncodingFormat {
	return provider.Encoding
}

// GetRootSchemaNode returns root schema node for entity
func (provider *CodecServiceProvider) GetRootSchemaNode(entity types.Entity) types.RootSchemaNode {
	root_schema_node, ok := provider.RootSchemaTable[entity.GetBundleName()]
	if !ok {
		panic("Root schema node not found in provider!")
	}
	return root_schema_node
}
