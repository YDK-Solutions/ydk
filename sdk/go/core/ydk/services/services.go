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
package services

// #cgo CXXFLAGS: -g -std=c++11
// #cgo darwin LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #cgo linux LDFLAGS:  -fprofile-arcs -ftest-coverage --coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lstdc++ -lpython2.7 -ldl
// #include <ydk/ydk.h>
// #include <stdlib.h>
import "C"

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/path"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
)

// CrudService supports CRUD operations on entities.
type CrudService struct {
}

// Create the entity.
// Returns whether the operation was successful or not (bool)
func (c *CrudService) Create(provider types.ServiceProvider, entity types.Entity) bool {
	return operationSucceeded(path.ExecuteRPC(provider, entity, "ydk:create", "entity", false))
}

// Update the entity.
// Returns whether the operation was successful or not (bool)
func (c *CrudService) Update(provider types.ServiceProvider, entity types.Entity) bool {
	return operationSucceeded(path.ExecuteRPC(provider, entity, "ydk:update", "entity", false))
}

// Delete the entity.
// Returns whether the operation was successful or not (bool)
func (c *CrudService) Delete(provider types.ServiceProvider, entity types.Entity) bool {
	return operationSucceeded(path.ExecuteRPC(provider, entity, "ydk:delete", "entity", false))
}

// Read the entity.
// Returns the entity as identified by the given filter (types.Entity)
func (c *CrudService) Read(provider types.ServiceProvider, filter types.Entity) types.Entity {
	return path.ReadDatanode(filter, path.ExecuteRPC(provider, filter, "ydk:read", "filter", true))
}

// ReadConfig only reads config.
// Returns the entity as identified by the given filter (types.Entity)
func (c *CrudService) ReadConfig(provider types.ServiceProvider, filter types.Entity) types.Entity {
	return path.ReadDatanode(filter, path.ExecuteRPC(provider, filter, "ydk:read", "filter", false))
}

func operationSucceeded(node types.DataNode) bool {
	return node.Private != nil
}

// CodecService supports encoding and decoding Go model API objects of type
type CodecService struct {
}

// Encode converts entity object to XML/JSON payload
func (c *CodecService) Encode(provider types.CodecServiceProvider, entity types.Entity) string {
	// 1. initialize provider, set root schema node for entity
	provider.Initialize(entity)
	// 2. get data node from root schema
	rootSchemaNode := provider.GetRootSchemaNode(entity)
	// 3. encode and return payload
	return path.CodecServiceEncode(provider.GetState(), entity, rootSchemaNode, provider.GetEncoding())
}

// Decode converts XML/JSON object to entity object
func (c *CodecService) Decode(provider types.CodecServiceProvider, payload string) types.Entity {
	// 1. parse payload, get topEntity
	nmsp := getEntityLookupKey(provider, payload)
	topEntity := ydk.GetTopEntity(nmsp)
	// 2. initialize repository, fetch rootSchema
	provider.Initialize(topEntity)
	rootSchema := provider.GetRootSchemaNode(topEntity)
	// 3. populate and return entity
	return path.CodecServiceDecode(provider.GetState(), rootSchema, payload, provider.GetEncoding(), topEntity)
}

func getEntityLookupKey(provider types.CodecServiceProvider, payload string) string {
	var nmsp string
	encoding := provider.GetEncoding()

	switch encoding {

	case types.XML:
		fmt.Println("Using XML encoding...")

		type XMLObj struct {
			XMLName xml.Name
		}

		var xmlObj XMLObj
		err := xml.Unmarshal([]byte(payload), &xmlObj)
		if err != nil {
			fmt.Println("Error parsing XML: ", err)
			panic(err)
		}

		nmsp = fmt.Sprintf("%v", xmlObj.XMLName)

	case types.JSON:
		fmt.Println("Using JSON encoding...")

		var jsonObj interface{}
		err := json.Unmarshal([]byte(payload), &jsonObj)
		if err != nil {
			fmt.Println("Error parsing JSON: ", err)
			panic(err)
		}

		items := jsonObj.(map[string]interface{})
		if len(items) != 1 {
			fmt.Println("List of payload not supported")
			panic("List of payload not supported")
		}

		for k := range items {
			nmsp = k
		}

	default:
		panic("Encoding not supported!")
	}

	return nmsp
}
