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

func ExecuteRpc(provider types.ServiceProvider, entity types.Entity, Filter string, dataTag string, setConfigFlag bool) types.DataNode {
	state := provider.GetState()
	cstate := getCState(state)
	wrappedProvider := provider.GetPrivate().(types.CServiceProvider)
	realProvider := wrappedProvider.Private.(C.ServiceProvider)
	rootSchema := C.ServiceProviderGetRootSchema(*cstate, realProvider)
	panicOnCStateError(cstate)

	ydkRpc := C.RootSchemaNodeRpc(*cstate, rootSchema, C.CString(Filter))
	panicOnCStateError(cstate)

	if rootSchema == nil {
		panic(1)
	}

	data := getDataPayload(state, entity, rootSchema, provider)
	defer C.free(unsafe.Pointer(data))

	input := C.RpcInput(*cstate, ydkRpc)
	panicOnCStateError(cstate)

	if setConfigFlag {
		C.DataNodeCreate(*cstate, input, C.CString("only-config"), C.CString(""))
		panicOnCStateError(cstate)
	}

	C.DataNodeCreate(*cstate, input, C.CString(dataTag), data)
	panicOnCStateError(cstate)

	dataNode := types.DataNode{C.RpcExecute(*cstate, ydkRpc, realProvider)}
	panicOnCStateError(cstate)

	return dataNode
}

func getDataPayload(state *types.State, entity types.Entity, rootSchema C.RootSchemaNode, provider types.ServiceProvider) *C.char {
	datanode := getDataNodeFromEntity(state, entity, rootSchema)

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

func ReadDatanode(filter types.Entity, readDataNode types.DataNode) types.Entity {
	if readDataNode.Private == nil {
		return nil
	}

	topEntity := getTopEntityFromFilter(filter)
	fmt.Printf("Reading top entity: '%s'\n", topEntity.GetSegmentPath())

	cchildren := C.DataNodeGetChildren(readDataNode.Private.(C.DataNode))
	children := (*[1 << 30]C.DataNode)(unsafe.Pointer(cchildren.datanodes))[:cchildren.count:cchildren.count]
	getEntityFromDataNode(children[0], topEntity)
	return topEntity
}

// ConnectToNetconfProvider connects to NETCONF service provider and returns types.CServiceProvider
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

// DisconnectFromNetconfProvider disconnects from NETCONF device and frees types.CServiceProvider passed in
func DisconnectFromNetconfProvider(provider types.CServiceProvider) {
	realProvider := provider.Private.(C.ServiceProvider)
	C.NetconfServiceProviderFree(realProvider)
}

// CleanUpErrorState cleans up memory for CState
func CleanUpErrorState(state *types.State) {
	realState := getCState(state)
	C.YDKStateFree(*realState)
}

// ConnectToRestconfProvider connects to RESTCONF device and returns types.CServiceProvider
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

// DisconnectFromRestconfProvider disconnects from RESTCONF device and frees types.CServiceProvider passed in
func DisconnectFromRestconfProvider(provider types.CServiceProvider) {
	realProvider := provider.Private.(C.ServiceProvider)
	C.RestconfServiceProviderFree(realProvider)
}

// InitCodecServiceProvider initializes CodecServiceProvider and returns root schema node parsed from repository
func InitCodecServiceProvider(state *types.State, entity types.Entity, repo types.Repository) types.RootSchemaNode {
	caps := entity.GetAugmentCapabilitiesFunction()()

	var repoPath *C.char
	defer C.free(unsafe.Pointer(repoPath))

	if len(repo.Path) > 0 {
		fmt.Printf("CodecServiceProvider using YANG models in %v\n", repo.Path)
		repoPath = C.CString(repo.Path)
	} else {
		yangPath := entity.GetBundleYangModelsLocation()
		fmt.Printf("CodecServiceProvider using YANG models in %v\n", yangPath)
		repoPath = C.CString(yangPath)
	}

	realCaps := make([]C.Capability, 0)
	var realCap C.Capability
	for mod, rev := range caps {
		realCap = C.CapabilityCreate(*getCState(state), C.CString(mod), C.CString(rev))
		panicOnCStateError(getCState(state))
		defer C.CapabilityFree(realCap)
		realCaps = append(realCaps, realCap)
	}

	realRepo := C.RepositoryInitWithPath(*getCState(state), repoPath)
	panicOnCStateError(getCState(state))

	repo.Private = realRepo
	rootSchemaWrapper := C.RepositoryCreateRootSchemaWrapper(*getCState(state), realRepo, &realCaps[0], C.int(len(realCaps)))
	panicOnCStateError(getCState(state))

	rootSchemaNode := types.RootSchemaNode{Private: rootSchemaWrapper}
	return rootSchemaNode
}

// CodecServiceEncode encodes entity to XML/JSON payloads based on encoding format passed in
func CodecServiceEncode(state *types.State, entity types.Entity, rootSchema types.RootSchemaNode, encoding types.EncodingFormat) string {
	rootSchemaWrapper := rootSchema.Private.(C.RootSchemaWrapper)
	realRootSchema := C.RootSchemaWrapperUnwrap(rootSchemaWrapper)

	dataNode := getDataNodeFromEntity(state, entity, realRootSchema)

	if dataNode == nil {
		return ""
	}

	codec := C.CodecInit()
	defer C.CodecFree(codec)

	var payload *C.char
	defer C.free(unsafe.Pointer(payload))

	switch encoding {
	case types.XML:
		payload = C.CodecEncode(*getCState(state), codec, dataNode, C.XML, 1)
		panicOnCStateError(getCState(state))
	case types.JSON:
		payload = C.CodecEncode(*getCState(state), codec, dataNode, C.JSON, 1)
		panicOnCStateError(getCState(state))
	}

	return C.GoString(payload)
}

// CodecServiceDecode decodes XML/JSON payloads passed in to entity
func CodecServiceDecode(state *types.State, rootSchema types.RootSchemaNode, payload string, encoding types.EncodingFormat, topEntity types.Entity) types.Entity {
	rootSchemaWrapper := rootSchema.Private.(C.RootSchemaWrapper)
	realRootSchema := C.RootSchemaWrapperUnwrap(rootSchemaWrapper)

	codec := C.CodecInit()
	defer C.CodecFree(codec)

	var realPayload = C.CString(payload)
	defer C.free(unsafe.Pointer(realPayload))
	var realDataNode C.DataNode

	switch encoding {
	case types.XML:
		realDataNode = C.CodecDecode(*getCState(state), codec, realRootSchema, realPayload, C.XML)
		panicOnCStateError(getCState(state))
	case types.JSON:
		realDataNode = C.CodecDecode(*getCState(state), codec, realRootSchema, realPayload, C.JSON)
		panicOnCStateError(getCState(state))
	}

	var dataNode = types.DataNode{Private: realDataNode}

	return ReadDatanode(topEntity, dataNode)
}

// ConnectToOpenDaylightProvider connects to OpenDaylight device and returns types.COpenDaylightServiceProvier
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

// DisconnectFromOpenDaylightProvider disconnects from OpenDaylight device and frees memory allocated
func DisconnectFromOpenDaylightProvider(provider types.COpenDaylightServiceProvider) {
	realProvider := provider.Private.(C.OpenDaylightServiceProvider)
	C.OpenDaylightServiceProviderFree(realProvider)
}

// OpenDaylightServiceProviderGetNodeIDS returns node ids available
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

// OpenDaylightServiceProviderGetNodeProvider returns service provider based on node id passed in
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
func getDataNodeFromEntity(state *types.State, entity types.Entity, rootSchema C.RootSchemaNode) C.DataNode {
	if entity == nil {
		return nil
	}
	for parent := entity.GetParent(); parent != nil; parent = parent.GetParent() {
		entity = parent
	}

	rootPath := entity.GetEntityPath(nil)
	path := C.CString(rootPath.Path)
	defer C.free(unsafe.Pointer(path))

	rootDataNode := C.RootSchemaNodeCreate(*getCState(state), rootSchema, path)
	panicOnCStateError(getCState(state))

	if types.IsSet(entity.GetFilter()) {
		p1 := C.CString(string(entity.GetFilter()))
		defer C.free(unsafe.Pointer(p1))
		C.DataNodeAddAnnotation(rootDataNode, p1)
	}

	populateNameValues(state, rootDataNode, rootPath)
	walkChildren(state, entity, rootDataNode)
	return rootDataNode
}

func walkChildren(state *types.State, entity types.Entity, dataNode C.DataNode) {
	children := entity.GetChildren()

	fmt.Printf("Got %d entity children\n", len(children))

	for childName := range children {

		fmt.Printf("Lookin at entity child '%s'\n", children[childName].GetSegmentPath())

		if children[childName].HasDataOrFilter() {
			populateDataNode(state, children[childName], dataNode)
		}
	}
	fmt.Println()
}

func populateDataNode(state *types.State, entity types.Entity, parentDataNode C.DataNode) {
	path := entity.GetEntityPath(entity.GetParent())
	p := C.CString(path.Path)
	defer C.free(unsafe.Pointer(p))
	ep := C.CString("")
	defer C.free(unsafe.Pointer(ep))

	dataNode := C.DataNodeCreate(*getCState(state), parentDataNode, p, ep)
	panicOnCStateError(getCState(state))

	if dataNode == nil {
		panic("Datanode could not be created for: " + path.Path)
	}

	if types.IsSet(entity.GetFilter()) {
		p1 := C.CString(string(entity.GetFilter()))
		defer C.free(unsafe.Pointer(p1))
		C.DataNodeAddAnnotation(dataNode, p1)
	}

	populateNameValues(state, dataNode, path)
	walkChildren(state, entity, dataNode)
}

func populateNameValues(state *types.State, dataNode C.DataNode, path types.EntityPath) {
	for _, nameValue := range path.ValuePaths {
		var result C.DataNode
		leafData := nameValue.Data
		p := C.CString(nameValue.Name)
		fmt.Printf("got leaf {%s: %s}\n", nameValue.Name, nameValue.Data.Value)

		if leafData.IsSet {
			p1 := C.CString(leafData.Value)
			result = C.DataNodeCreate(*getCState(state), dataNode, p, p1)
			panicOnCStateError(getCState(state))
			C.DataNodeCreate(*getCState(state), dataNode, p, p1)
			panicOnCStateError(getCState(state))
			C.free(unsafe.Pointer(p1))
		}

		if types.IsSet(leafData.Filter) {
			p1 := C.CString(string(nameValue.Data.Filter))
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

	cchildren := C.DataNodeGetChildren(node)
	children := (*[1 << 30]C.DataNode)(unsafe.Pointer(cchildren.datanodes))[:cchildren.count:cchildren.count]
	fmt.Printf("Got %d datanode children\n", cchildren.count)

	for _, childDataNode := range children {
		childName := C.GoString(C.DataNodeGetArgument(childDataNode))
		fmt.Printf("Lookin at child datanode: '%s'\n", childName)

		if dataNodeIsLeaf(childDataNode) {

			value := C.GoString(C.DataNodeGetValue(childDataNode))
			fmt.Printf("Creating leaf '%s' with value '%s'\n", childName, value)
			entity.SetValue(childName, value)
		} else {

			var childEntity types.Entity
			if dataNodeIsList(childDataNode) {
				segmentPath := C.GoString(C.DataNodeGetSegmentPath(childDataNode))
				fmt.Printf("Creating child list instance '%s'\n", segmentPath)
				childEntity = entity.GetChildByName(childName, segmentPath)
			} else {
				fmt.Printf("Creating child node '%s'\n", childName)
				childEntity = entity.GetChildByName(childName, "")
			}
			if childEntity == nil {
				panic("Could not create child entity!")
			}
			childEntity.SetParent(entity)
			getEntityFromDataNode(childDataNode, childEntity)
		}
	}
}

func dataNodeIsLeaf(dataNode C.DataNode) bool {
	return C.GoString(C.DataNodeGetKeyword(dataNode)) == "leaf" ||
		C.GoString(C.DataNodeGetKeyword(dataNode)) == "leaf-list"
}

func dataNodeIsList(dataNode C.DataNode) bool {
	return C.GoString(C.DataNodeGetKeyword(dataNode)) == "list"
}

//////////////////////////////////////////////////////////////////////////
// Error Handling
//////////////////////////////////////////////////////////////////////////

// AddCState creates and add cstate to *types.State
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
