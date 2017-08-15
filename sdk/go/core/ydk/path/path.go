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
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *----------------------------------------------
 */
package path

// #cgo CXXFLAGS: -g -std=c++11
// #cgo darwin LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #cgo linux LDFLAGS:  -fprofile-arcs -ftest-coverage --coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lstdc++ -lpython2.7 -ldl
// #include <ydk/ydk.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"unsafe"
)

func ExecuteRpc(provider types.ServiceProvider, entity types.Entity, Filter string, data_tag string, set_config_flag bool) types.DataNode {
	wrapped_provider := provider.GetPrivate().(types.CServiceProvider)
	real_provider := wrapped_provider.Private.(C.ServiceProvider)
	root_schema := C.ServiceProviderGetRootSchema(real_provider)

	ydk_rpc := C.RootSchemaNodeRpc(root_schema, C.CString(Filter))
	if root_schema == nil {
		panic(1)
	}

	data := getDataPayload(entity, root_schema, provider)
	defer C.free(unsafe.Pointer(data))

	input := C.RpcInput(ydk_rpc)

	if set_config_flag {
		C.DataNodeCreate(input, C.CString("only-config"), C.CString(""))
	}

	C.DataNodeCreate(input, C.CString(data_tag), data)
	return types.DataNode{C.RpcExecute(ydk_rpc, real_provider)}
}

func getDataPayload(entity types.Entity, root_schema C.RootSchemaNode, provider types.ServiceProvider) *C.char {
	datanode := getDataNodeFromEntity(entity, root_schema)

	if datanode == nil {
		return nil
	}

	//for datanode != nil && C.DataNodeGetParent(datanode) != nil {
	//	datanode = C.DataNodeGetParent(datanode)
	//}
	cprovider := provider.GetPrivate().(types.CServiceProvider).Private.(C.ServiceProvider)
	cencoding := C.ServiceProviderGetEncoding(cprovider)

	codec := C.CodecInit()
	defer C.CodecFree(codec)
	var data *C.char = C.CodecEncode(codec, datanode, cencoding, 1)

	return (data)
}

func getTopEntityFromFilter(filter types.Entity) types.Entity {
	if filter.GetParent() == nil {
		return filter
	}

	return getTopEntityFromFilter(filter.GetParent())
}

func ReadDatanode(filter types.Entity, read_data_node types.DataNode) types.Entity {
	if read_data_node.Private == nil {
		return nil
	}

	top_entity := getTopEntityFromFilter(filter)
	fmt.Printf("Reading top entity: '%s'\n", top_entity.GetSegmentPath())

	c_children := C.DataNodeGetChildren(read_data_node.Private.(C.DataNode))
	children := (*[1 << 30]C.DataNode)(unsafe.Pointer(c_children.datanodes))[:c_children.count:c_children.count]
	getEntityFromDataNode(children[0], top_entity)
	return top_entity
}

func ConnectToProvider(repo types.Repository, Address, Username, Password string, port int) types.CServiceProvider {
	var address *C.char = C.CString(Address)
	defer C.free(unsafe.Pointer(address))
	var username *C.char = C.CString(Username)
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString(Password)
	defer C.free(unsafe.Pointer(password))
	var cport C.int = C.int(port)

	var p C.ServiceProvider

	if len(repo.Path) > 0 {
		var path *C.char = C.CString(repo.Path)
		repo := C.RepositoryInitWithPath(path)
		p = C.NetconfServiceProviderInitWithRepo(repo, address, username, password, cport)
	} else {
		p = C.NetconfServiceProviderInit(address, username, password, cport)
	}
	if p == nil {
		panic("Could not connect to " + Address)
	}
	cprovider := types.CServiceProvider{Private: p}
	return cprovider
}

func DisconnectFromProvider(provider types.CServiceProvider) {
	real_provider := provider.Private.(C.ServiceProvider)
	C.NetconfServiceProviderFree(real_provider)
}

func ConnectToRestconfProvider(Path, Address, Username, Password string, port int) types.CServiceProvider {
	var path *C.char = C.CString(Path)
	defer C.free(unsafe.Pointer(path))
	var address *C.char = C.CString(Address)
	defer C.free(unsafe.Pointer(address))
	var username *C.char = C.CString(Username)
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString(Password)
	defer C.free(unsafe.Pointer(password))
	var cport C.int = C.int(port)

	var p C.ServiceProvider

	crepo := C.RepositoryInitWithPath(path)
	p = C.RestconfServiceProviderInitWithRepo(crepo, address, username, password, cport)

	if p == nil {
		panic("Could not connect to " + Address)
	}

	cprovider := types.CServiceProvider{Private: p}
	return cprovider
}

func DisconnectFromRestconfProvider(provider types.CServiceProvider) {
	real_provider := provider.Private.(C.ServiceProvider)
	C.RestconfServiceProviderFree(real_provider)
}

func InitCodecServiceProvider(entity types.Entity, repo types.Repository) types.RootSchemaNode {
	caps := entity.GetAugmentCapabilitiesFunction()()

	var repo_path *C.char
	defer C.free(unsafe.Pointer(repo_path))

	if len(repo.Path) > 0 {
		fmt.Printf("CodecServiceProvider using YANG models in %v\n", repo.Path)
		repo_path = C.CString(repo.Path)
	} else {
		yang_path := entity.GetBundleYangModelsLocation()
		fmt.Printf("CodecServiceProvider using YANG models in %v\n", yang_path)
		repo_path = C.CString(yang_path)
	}

	real_caps := make([]C.Capability, 0)
	var real_cap C.Capability
	for mod, rev := range caps {
		real_cap = C.CapabilityCreate(C.CString(mod), C.CString(rev))
		defer C.CapabilityFree(real_cap)
		real_caps = append(real_caps, real_cap)
	}

	real_repo := C.RepositoryInitWithPath(repo_path)

	repo.Private = real_repo
	root_schema_wrapper := C.RepositoryCreateRootSchemaWrapper(real_repo, &real_caps[0], C.int(len(real_caps)))

	root_schema_node := types.RootSchemaNode{Private: root_schema_wrapper}
	return root_schema_node
}

func CodecServiceEncode(entity types.Entity, root_schema types.RootSchemaNode, encoding types.EncodingFormat) string {
	root_schema_wrapper := root_schema.Private.(C.RootSchemaWrapper)
	real_root_schema := C.RootSchemaWrapperUnwrap(root_schema_wrapper)

	data_node := getDataNodeFromEntity(entity, real_root_schema)

	if data_node == nil {
		return ""
	}

	codec := C.CodecInit()
	defer C.CodecFree(codec)

	var payload *C.char
	defer C.free(unsafe.Pointer(payload))

	switch encoding {
	case types.XML:
		payload = C.CodecEncode(codec, data_node, C.XML, 1)
	case types.JSON:
		payload = C.CodecEncode(codec, data_node, C.JSON, 1)
	}

	return C.GoString(payload)
}

func CodecServiceDecode(root_schema types.RootSchemaNode, payload string, encoding types.EncodingFormat, top_entity types.Entity) types.Entity {
	root_schema_wrapper := root_schema.Private.(C.RootSchemaWrapper)
	real_root_schema := C.RootSchemaWrapperUnwrap(root_schema_wrapper)

	codec := C.CodecInit()
	defer C.CodecFree(codec)

	var real_payload = C.CString(payload)
	defer C.free(unsafe.Pointer(real_payload))
	var real_data_node C.DataNode

	switch encoding {
	case types.XML:
		real_data_node = C.CodecDecode(codec, real_root_schema, real_payload, C.XML)
	case types.JSON:
		real_data_node = C.CodecDecode(codec, real_root_schema, real_payload, C.JSON)
	}

	var data_node = types.DataNode{Private: real_data_node}

	return ReadDatanode(top_entity, data_node)
}

func ConnectToOpenDaylightProvider(Path, Address, Username, Password string, port int, encoding types.EncodingFormat, protocol types.Protocol) types.COpenDaylightServiceProvider {
	var path *C.char = C.CString(Path)
	defer C.free(unsafe.Pointer(path))
	var address *C.char = C.CString(Address)
	defer C.free(unsafe.Pointer(address))
	var username *C.char = C.CString(Username)
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString(Password)
	defer C.free(unsafe.Pointer(password))
	var cport C.int = C.int(port)

	var p C.OpenDaylightServiceProvider
	crepo := C.RepositoryInitWithPath(path)

	var cencoding C.EncodingFormat
	if encoding == types.XML {
		cencoding = C.XML
	} else {
		cencoding = C.JSON
	}

	var cprotocol C.Protocol
	if protocol == types.Netconf {
		cprotocol = C.Netconf
	} else {
		cprotocol = C.Restconf
	}

	p = C.OpenDaylightServiceProviderInitWithRepo(crepo, address, username, password, cport, cencoding, cprotocol)

	if p == nil {
		panic("Could not connect to " + Address)
	}

	cprovider := types.COpenDaylightServiceProvider{Private: p}
	return cprovider
}

func DisconnectFromOpenDaylightProvider(provider types.COpenDaylightServiceProvider) {
	real_provider := provider.Private.(C.OpenDaylightServiceProvider)
	C.OpenDaylightServiceProviderFree(real_provider)
}

func OpenDaylightServiceProviderGetNodeIds(provider types.COpenDaylightServiceProvider) []string {
	cprovider := provider.Private.(C.OpenDaylightServiceProvider)
	var ids []string
	id := 0
	for {
		cid := C.int(id)
		nodeId := C.OpenDaylightServiceProviderGetNodeIdByIndex(cprovider, cid)
		defer C.free(unsafe.Pointer(nodeId))
		if nodeId != nil {
			ids = append(ids, C.GoString(nodeId))
			id ++
		} else {
			break
		}
	}
	return ids
}

func OpenDaylightServiceProviderGetNodeProvider(provider types.COpenDaylightServiceProvider, nodeId string) types.CServiceProvider {
	realProvider := provider.Private.(C.OpenDaylightServiceProvider)
	cnodeId := C.CString(nodeId)
	defer C.free(unsafe.Pointer(cnodeId))
	var nodeProvider C.ServiceProvider
	nodeProvider = C.OpenDaylightServiceProviderGetNodeProvider(realProvider, cnodeId)
	cnodeProvider := types.CServiceProvider{Private: nodeProvider}
	return cnodeProvider
}

//////////////////////////////////////////////////////////////////////////
// DataNode from Entity
//////////////////////////////////////////////////////////////////////////
func getDataNodeFromEntity(entity types.Entity, root_schema C.RootSchemaNode) C.DataNode {
	if entity == nil {
		return nil
	}
	for parent := entity.GetParent(); parent != nil; parent = parent.GetParent() {
		entity = parent
	}

	root_path := entity.GetEntityPath(nil)
	path := C.CString(root_path.Path)
	defer C.free(unsafe.Pointer(path))

	root_data_node := C.RootSchemaNodeCreate(root_schema, path)

	if types.IsSet(entity.GetFilter()) {
		p1 := C.CString(string(entity.GetFilter()))
		defer C.free(unsafe.Pointer(p1))
		C.DataNodeAddAnnotation(root_data_node, p1)
	}

	populateNameValues(root_data_node, root_path)
	walkChildren(entity, root_data_node)
	return root_data_node
}

func walkChildren(entity types.Entity, data_node C.DataNode) {
	children := entity.GetChildren()

	fmt.Printf("Got %d entity children\n", len(children))

	for child_name := range children {

		fmt.Printf("Lookin at entity child '%s'\n", children[child_name].GetSegmentPath())

		if children[child_name].HasDataOrFilter() {
			populateDataNode(children[child_name], data_node)
		}
	}
	fmt.Println()
}

func populateDataNode(entity types.Entity, parent_data_node C.DataNode) {
	path := entity.GetEntityPath(entity.GetParent())
	p := C.CString(path.Path)
	defer C.free(unsafe.Pointer(p))
	ep := C.CString("")
	defer C.free(unsafe.Pointer(ep))

	data_node := C.DataNodeCreate(parent_data_node, p, ep)
	if data_node == nil {
		panic("Datanode could not be created for: " + path.Path)
	}

	if types.IsSet(entity.GetFilter()) {
		p1 := C.CString(string(entity.GetFilter()))
		defer C.free(unsafe.Pointer(p1))
		C.DataNodeAddAnnotation(data_node, p1)
	}

	populateNameValues(data_node, path)
	walkChildren(entity, data_node)
}

func populateNameValues(data_node C.DataNode, path types.EntityPath) {
	for _, name_value := range path.ValuePaths {
		var result C.DataNode
		leaf_data := name_value.Data
		p := C.CString(name_value.Name)
		fmt.Printf("got leaf {%s: %s}\n", name_value.Name, name_value.Data.Value)

		if leaf_data.IsSet {
			p1 := C.CString(leaf_data.Value)
			result = C.DataNodeCreate(data_node, p, p1)
			C.DataNodeCreate(data_node, p, p1)
			C.free(unsafe.Pointer(p1))
		}

		if types.IsSet(leaf_data.Filter) {
			p1 := C.CString(string(name_value.Data.Filter))
			defer C.free(unsafe.Pointer(p1))
			C.DataNodeAddAnnotation(result, p1)
		}
		C.free(unsafe.Pointer(p))
	}
}

//////////////////////////////////////////////////////////////////////////
// Entity from DataNode
//////////////////////////////////////////////////////////////////////////
func getEntityFromDataNode(node C.DataNode, entity types.Entity) {
	if entity == nil || node == nil {
		return
	}

	c_children := C.DataNodeGetChildren(node)
	children := (*[1 << 30]C.DataNode)(unsafe.Pointer(c_children.datanodes))[:c_children.count:c_children.count]
	fmt.Printf("Got %d datanode children\n", c_children.count)

	for _, child_data_node := range children {
		child_name := C.GoString(C.DataNodeGetArgument(child_data_node))
		fmt.Printf("Lookin at child datanode: '%s'\n", child_name)

		if dataNodeIsLeaf(child_data_node) {

			value := C.GoString(C.DataNodeGetValue(child_data_node))
			fmt.Printf("Creating leaf '%s' with value '%s'\n", child_name, value)
			entity.SetValue(child_name, value)
		} else {

			var child_entity types.Entity
			if dataNodeIsList(child_data_node) {
				segment_path := C.GoString(C.DataNodeGetSegmentPath(child_data_node))
				fmt.Printf("Creating child list instance '%s'\n", segment_path)
				child_entity = entity.GetChildByName(child_name, segment_path)
			} else {
				fmt.Printf("Creating child node '%s'\n", child_name)
				child_entity = entity.GetChildByName(child_name, "")
			}
			if child_entity == nil {
				panic("Could not create child entity!")
			}
			child_entity.SetParent(entity)
			getEntityFromDataNode(child_data_node, child_entity)
		}
	}
}

func dataNodeIsLeaf(data_node C.DataNode) bool {
	return C.GoString(C.DataNodeGetKeyword(data_node)) == "leaf" ||
		C.GoString(C.DataNodeGetKeyword(data_node)) == "leaf-list"
}

func dataNodeIsList(data_node C.DataNode) bool {
	return C.GoString(C.DataNodeGetKeyword(data_node)) == "list"
}
