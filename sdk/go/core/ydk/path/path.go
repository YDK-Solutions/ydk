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
	"strings"
	"unsafe"
)

func ExecuteRpc(provider types.ServiceProvider, entity types.Entity, Filter string, data_tag string, set_config_flag bool) types.DataNode {
	state := provider.GetState()
	cstate := getCState(state)
	wrapped_provider := provider.GetPrivate().(types.CServiceProvider)
	real_provider := wrapped_provider.Private.(C.ServiceProvider)
	root_schema := C.ServiceProviderGetRootSchema(*cstate, real_provider)
	panicOnCStateError(cstate)

	ydk_rpc := C.RootSchemaNodeRpc(*cstate, root_schema, C.CString(Filter))
	panicOnCStateError(cstate)

	if root_schema == nil {
		panic(1)
	}

	data := getDataPayload(state, entity, root_schema, provider)
	defer C.free(unsafe.Pointer(data))

	input := C.RpcInput(*cstate, ydk_rpc)
	panicOnCStateError(cstate)

	if set_config_flag {
		C.DataNodeCreate(*cstate, input, C.CString("only-config"), C.CString(""))
		panicOnCStateError(cstate)
	}

	C.DataNodeCreate(*cstate, input, C.CString(data_tag), data)
	panicOnCStateError(cstate)

	dataNode := types.DataNode{C.RpcExecute(*cstate, ydk_rpc, real_provider)}
	panicOnCStateError(cstate)

	return dataNode
}

func getDataPayload(state *types.State, entity types.Entity, root_schema C.RootSchemaNode, provider types.ServiceProvider) *C.char {
	datanode := getDataNodeFromEntity(state, entity, root_schema)

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
	var data *C.char = C.CodecEncode(*getCState(state), codec, datanode, cencoding, 1)
	panicOnCStateError(getCState(state))

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

func ConnectToNetconfProvider(state *types.State, repo types.Repository, Address, Username, Password string, port int, Protocol string) types.CServiceProvider {
	var address *C.char = C.CString(Address)
	defer C.free(unsafe.Pointer(address))
	var username *C.char = C.CString(Username)
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString(Password)
	defer C.free(unsafe.Pointer(password))
	var cport C.int = C.int(port)

	var cprotocol *C.char = C.CString(Protocol)
	defer C.free(unsafe.Pointer(cprotocol))

	AddCState(state)
	cstate := getCState(state)

	var p C.ServiceProvider

	if len(repo.Path) > 0 {
		var path *C.char = C.CString(repo.Path)
		repo := C.RepositoryInitWithPath(*cstate, path)
		panicOnCStateError(cstate)
		p = C.NetconfServiceProviderInitWithRepo(*cstate, repo, address, username, password, cport, cprotocol)
		panicOnCStateError(cstate)
	} else {
		p = C.NetconfServiceProviderInit(*cstate, address, username, password, cport, cprotocol)
		panicOnCStateError(cstate)
	}

	cprovider := types.CServiceProvider{Private: p}
	return cprovider
}

func DisconnectFromNetconfProvider(provider types.CServiceProvider) {
	real_provider := provider.Private.(C.ServiceProvider)
	C.NetconfServiceProviderFree(real_provider)
}

func CleanUpErrorState(state *types.State) {
	real_state := getCState(state)
	C.YDKStateFree(*real_state)
}

func ConnectToRestconfProvider(state *types.State, Path, Address, Username, Password string, port int, encoding types.EncodingFormat, stateURLRoot, configURLRoot string) types.CServiceProvider {
	var path *C.char = C.CString(Path)
	defer C.free(unsafe.Pointer(path))
	var address *C.char = C.CString(Address)
	defer C.free(unsafe.Pointer(address))
	var username *C.char = C.CString(Username)
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString(Password)
	defer C.free(unsafe.Pointer(password))
	var cport C.int = C.int(port)

	cencoding := getCEncoding(encoding)

	var cstateURLRoot *C.char = C.CString(stateURLRoot)
	defer C.free(unsafe.Pointer(cstateURLRoot))

	var cconfigURLRoot *C.char = C.CString(configURLRoot)
	defer C.free(unsafe.Pointer(cconfigURLRoot))

	AddCState(state)
	cstate := getCState(state)

	var p C.ServiceProvider

	crepo := C.RepositoryInitWithPath(*getCState(state), path)
	panicOnCStateError(cstate)
	p = C.RestconfServiceProviderInitWithRepo(*cstate, crepo, address, username, password, cport, cencoding, cstateURLRoot, cconfigURLRoot)
	panicOnCStateError(cstate)

	cprovider := types.CServiceProvider{Private: p}
	return cprovider
}

func DisconnectFromRestconfProvider(provider types.CServiceProvider) {
	real_provider := provider.Private.(C.ServiceProvider)
	C.RestconfServiceProviderFree(real_provider)
}

func InitCodecServiceProvider(state *types.State, entity types.Entity, repo types.Repository) types.RootSchemaNode {
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
		real_cap = C.CapabilityCreate(*getCState(state), C.CString(mod), C.CString(rev))
		panicOnCStateError(getCState(state))
		defer C.CapabilityFree(real_cap)
		real_caps = append(real_caps, real_cap)
	}

	real_repo := C.RepositoryInitWithPath(*getCState(state), repo_path)
	panicOnCStateError(getCState(state))

	repo.Private = real_repo
	root_schema_wrapper := C.RepositoryCreateRootSchemaWrapper(*getCState(state), real_repo, &real_caps[0], C.int(len(real_caps)))
	panicOnCStateError(getCState(state))

	root_schema_node := types.RootSchemaNode{Private: root_schema_wrapper}
	return root_schema_node
}

func CodecServiceEncode(state *types.State, entity types.Entity, root_schema types.RootSchemaNode, encoding types.EncodingFormat) string {
	root_schema_wrapper := root_schema.Private.(C.RootSchemaWrapper)
	real_root_schema := C.RootSchemaWrapperUnwrap(root_schema_wrapper)

	data_node := getDataNodeFromEntity(state, entity, real_root_schema)

	if data_node == nil {
		return ""
	}

	codec := C.CodecInit()
	defer C.CodecFree(codec)

	var payload *C.char
	defer C.free(unsafe.Pointer(payload))

	switch encoding {
	case types.XML:
		payload = C.CodecEncode(*getCState(state), codec, data_node, C.XML, 1)
		panicOnCStateError(getCState(state))
	case types.JSON:
		payload = C.CodecEncode(*getCState(state), codec, data_node, C.JSON, 1)
		panicOnCStateError(getCState(state))
	}

	return C.GoString(payload)
}

func CodecServiceDecode(state *types.State, root_schema types.RootSchemaNode, payload string, encoding types.EncodingFormat, top_entity types.Entity) types.Entity {
	root_schema_wrapper := root_schema.Private.(C.RootSchemaWrapper)
	real_root_schema := C.RootSchemaWrapperUnwrap(root_schema_wrapper)

	codec := C.CodecInit()
	defer C.CodecFree(codec)

	var real_payload = C.CString(payload)
	defer C.free(unsafe.Pointer(real_payload))
	var real_data_node C.DataNode

	switch encoding {
	case types.XML:
		real_data_node = C.CodecDecode(*getCState(state), codec, real_root_schema, real_payload, C.XML)
		panicOnCStateError(getCState(state))
	case types.JSON:
		real_data_node = C.CodecDecode(*getCState(state), codec, real_root_schema, real_payload, C.JSON)
		panicOnCStateError(getCState(state))
	}

	var data_node = types.DataNode{Private: real_data_node}

	return ReadDatanode(top_entity, data_node)
}

func ConnectToOpenDaylightProvider(state *types.State, Path, Address, Username, Password string, port int, encoding types.EncodingFormat, protocol types.Protocol) types.COpenDaylightServiceProvider {
	var path *C.char = C.CString(Path)
	defer C.free(unsafe.Pointer(path))
	var address *C.char = C.CString(Address)
	defer C.free(unsafe.Pointer(address))
	var username *C.char = C.CString(Username)
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString(Password)
	defer C.free(unsafe.Pointer(password))
	var cport C.int = C.int(port)

	AddCState(state)
	cstate := getCState(state)

	var p C.OpenDaylightServiceProvider
	crepo := C.RepositoryInitWithPath(*getCState(state), path)
	panicOnCStateError(getCState(state))

	cencoding := getCEncoding(encoding)

	cprotocol := getCProtocol(protocol)

	p = C.OpenDaylightServiceProviderInitWithRepo(*cstate, crepo, address, username, password, cport, cencoding, cprotocol)
	panicOnCStateError(cstate)

	cprovider := types.COpenDaylightServiceProvider{Private: p}
	return cprovider
}

func DisconnectFromOpenDaylightProvider(provider types.COpenDaylightServiceProvider) {
	real_provider := provider.Private.(C.OpenDaylightServiceProvider)
	C.OpenDaylightServiceProviderFree(real_provider)
}

func OpenDaylightServiceProviderGetNodeIDs(state *types.State, provider types.COpenDaylightServiceProvider) []string {
	cprovider := provider.Private.(C.OpenDaylightServiceProvider)
	var ids []string
	id := 0
	for {
		cid := C.int(id)
		nodeID := C.OpenDaylightServiceProviderGetNodeIDByIndex(*getCState(state), cprovider, cid)
		panicOnCStateError(getCState(state))
		defer C.free(unsafe.Pointer(nodeID))
		if nodeID != nil {
			ids = append(ids, C.GoString(nodeID))
			id++
		} else {
			break
		}
	}
	return ids
}

func OpenDaylightServiceProviderGetNodeProvider(state *types.State, provider types.COpenDaylightServiceProvider, nodeID string) types.CServiceProvider {
	realProvider := provider.Private.(C.OpenDaylightServiceProvider)
	cnodeID := C.CString(nodeID)
	defer C.free(unsafe.Pointer(cnodeID))
	var nodeProvider C.ServiceProvider
	nodeProvider = C.OpenDaylightServiceProviderGetNodeProvider(*getCState(state), realProvider, cnodeID)
	panicOnCStateError(getCState(state))

	cnodeProvider := types.CServiceProvider{Private: nodeProvider}
	return cnodeProvider
}

//////////////////////////////////////////////////////////////////////////
// DataNode from Entity
//////////////////////////////////////////////////////////////////////////
func getDataNodeFromEntity(state *types.State, entity types.Entity, root_schema C.RootSchemaNode) C.DataNode {
	if entity == nil {
		return nil
	}
	for parent := entity.GetParent(); parent != nil; parent = parent.GetParent() {
		entity = parent
	}

	root_path := entity.GetEntityPath(nil)
	path := C.CString(root_path.Path)
	defer C.free(unsafe.Pointer(path))

	root_data_node := C.RootSchemaNodeCreate(*getCState(state), root_schema, path)
	panicOnCStateError(getCState(state))

	if types.IsSet(entity.GetFilter()) {
		p1 := C.CString(string(entity.GetFilter()))
		defer C.free(unsafe.Pointer(p1))
		C.DataNodeAddAnnotation(root_data_node, p1)
	}

	populateNameValues(state, root_data_node, root_path)
	walkChildren(state, entity, root_data_node)
	return root_data_node
}

func walkChildren(state *types.State, entity types.Entity, data_node C.DataNode) {
	children := entity.GetChildren()

	fmt.Printf("Got %d entity children\n", len(children))

	for child_name := range children {

		fmt.Printf("Lookin at entity child '%s'\n", children[child_name].GetSegmentPath())

		if children[child_name].HasDataOrFilter() {
			populateDataNode(state, children[child_name], data_node)
		}
	}
	fmt.Println()
}

func populateDataNode(state *types.State, entity types.Entity, parent_data_node C.DataNode) {
	path := entity.GetEntityPath(entity.GetParent())
	p := C.CString(path.Path)
	defer C.free(unsafe.Pointer(p))
	ep := C.CString("")
	defer C.free(unsafe.Pointer(ep))

	data_node := C.DataNodeCreate(*getCState(state), parent_data_node, p, ep)
	panicOnCStateError(getCState(state))

	if data_node == nil {
		panic("Datanode could not be created for: " + path.Path)
	}

	if types.IsSet(entity.GetFilter()) {
		p1 := C.CString(string(entity.GetFilter()))
		defer C.free(unsafe.Pointer(p1))
		C.DataNodeAddAnnotation(data_node, p1)
	}

	populateNameValues(state, data_node, path)
	walkChildren(state, entity, data_node)
}

func populateNameValues(state *types.State, data_node C.DataNode, path types.EntityPath) {
	for _, name_value := range path.ValuePaths {
		var result C.DataNode
		leaf_data := name_value.Data
		p := C.CString(name_value.Name)
		fmt.Printf("got leaf {%s: %s}\n", name_value.Name, name_value.Data.Value)

		if leaf_data.IsSet {
			p1 := C.CString(leaf_data.Value)
			result = C.DataNodeCreate(*getCState(state), data_node, p, p1)
			panicOnCStateError(getCState(state))
			C.DataNodeCreate(*getCState(state), data_node, p, p1)
			panicOnCStateError(getCState(state))
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

//////////////////////////////////////////////////////////////////////////
// Error Handling
//////////////////////////////////////////////////////////////////////////
func AddCState(state *types.State) {
	cstate := C.YDKStateCreate()
	state.Private = types.CState{Private: cstate}
}

func checkState(cstate *C.YDKStatePtr) error {
	var cerrorOccurred C.boolean
	cerrorOccurred = C.YDKStateErrorOccurred(*cstate)
	if cerrorOccurred == 0 {
		return nil
	}

	rawErrMsg := C.GoString(C.YDKStateGetErrorMessage(*cstate))
	i := strings.Index(rawErrMsg, ":")
	errMsg := rawErrMsg[i+1:]

	var cerrorType C.YDKErrorType
	cerrorType = C.YDKStateGetErrorType(*cstate)

	switch cerrorType {
	case C.YDK_CLIENT_ERROR:
		return &types.YGOClientError{Msg: errMsg}
	case C.YDK_SERVICE_PROVIDER_ERROR:
		return &types.YGOServiceProviderError{Msg: errMsg}
	case C.YDK_SERVICE_ERROR:
		return &types.YGOServiceError{Msg: errMsg}
	case C.YDK_ILLEGAL_STATE_ERROR:
		return &types.YGOIllegalStateError{Msg: errMsg}
	case C.YDK_INVALID_ARGUMENT_ERROR:
		return &types.YGOInvalidArgumentError{Msg: errMsg}
	case C.YDK_OPERATION_NOTSUPPORTED_ERROR:
		return &types.YGOOperationNotSupportedError{Msg: errMsg}
	case C.YDK_MODEL_ERROR:
		return &types.YGOModelError{Msg: errMsg}
	case C.YDK_CORE_ERROR:
		return &types.YGOCoreError{Msg: errMsg}
	case C.YDK_CODEC_ERROR:
		return &types.YGOCodecError{Msg: errMsg}
	default:
		return &types.YGOError{Msg: errMsg}
	}
}

func getCProtocol(protocol types.Protocol) C.Protocol {
	if protocol == types.Netconf {
		return C.Netconf
	} else {
		return C.Restconf
	}
}

func getCState(state *types.State) *C.YDKStatePtr {
	statePtr := state.Private.(types.CState).Private.(C.YDKStatePtr)
	return &statePtr
}

func getCEncoding(encoding types.EncodingFormat) C.EncodingFormat {
	if encoding == types.XML {
		return C.XML
	} else {
		return C.JSON
	}
}

func panicOnCStateError(cstate *C.YDKStatePtr) {
	err := checkState(cstate)
	if err != nil {
		panic(err.Error())
	}
}
