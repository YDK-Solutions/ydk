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
package path

// #cgo CXXFLAGS: -g -std=c++11
// #cgo LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #include <ydk/ydk.h>
// #include <stdlib.h>
import "C"

import (
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"unsafe"
)

func ExecuteRpc(provider types.ServiceProvider, entity types.Entity, operation string, data_tag string, set_config_flag bool) types.DataNode {
	wrapped_provider := provider.GetPrivate().(types.CServiceProvider)
	real_provider := wrapped_provider.Private.(C.ServiceProvider)
	root_schema := C.ServiceProviderGetRootSchema(real_provider)

    //fmt.Println("got root & prov", root_schema, real_provider, operation, data_tag, &entity)

	ydk_rpc := C.RootSchemaNodeRpc(root_schema, C.CString(operation))
    if root_schema == nil {
        panic(1)
    }
    var data *C.char
    if data_tag == "entity" {
        data = getDataPayload(entity, root_schema)
    } else {
        data = C.CString("<bgp xmlns=\"http://openconfig.net/yang/bgp\"/>")
    }

    //fmt.Println("deeee!", C.GoString(data))

	input := C.RpcInput(ydk_rpc)

	if set_config_flag {
		C.DataNodeCreate(input, C.CString("only-config"), C.CString(""))
	}

	C.DataNodeCreate(input, C.CString(data_tag), data)
	return types.DataNode{C.RpcExecute(ydk_rpc, real_provider)}
}

func getDataPayload(entity types.Entity, root_schema C.RootSchemaNode) *C.char {
	datanode := getDataNodeFromEntity(entity, root_schema)

    //fmt.Println("hoooo>> ", entity.GetSegmentPath(), datanode)

	if datanode == nil {
		return nil
	}

	//for datanode != nil && C.DataNodeGetParent(datanode) != nil {
	//	datanode = C.DataNodeGetParent(datanode)
	//}

	codec := C.CodecServiceInit()
	defer C.CodecServiceFree(codec)
	var data *C.char = C.CodecServiceEncode(codec, datanode, C.XML, 1)
	//defer C.free(unsafe.Pointer(data))

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
	//children := C.DataNodeGetChildren(read_data_node)

	//types.getEntityFromDataNode(children, top_entity)
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
    cprovider := types.CServiceProvider{Private:p}
	return cprovider
}

func DisconnectFromProvider(provider types.CServiceProvider) {
	real_provider := provider.Private.(C.ServiceProvider)
	C.NetconfServiceProviderFree(real_provider)
}

//////////////////////////////////////////////////////////////////////////
// DataNode from Entity
//////////////////////////////////////////////////////////////////////////
func getDataNodeFromEntity(entity types.Entity, root_schema C.RootSchemaNode) C.DataNode {
	if entity == nil {
		return nil
	}
	root_path := entity.GetEntityPath(nil)
	path := C.CString(root_path.Path)
	defer C.free(unsafe.Pointer(path))

	root_data_node := C.RootSchemaNodeCreate(root_schema, path)

    /*
    HACK!!!
     */
    res := C.DataNodeCreate(root_data_node, C.CString("global/config/as"), C.CString("65001"))
    res = C.DataNodeCreate(root_data_node, C.CString("global/config/router-id"), C.CString("1.2.3.4"))

    res = C.DataNodeCreate(root_data_node, C.CString("global/afi-safis/afi-safi[afi-safi-name=\"openconfig-bgp-types:IPV6_UNICAST\"]/config/afi-safi-name"), C.CString("openconfig-bgp-types:IPV6_UNICAST"))
    res = C.DataNodeCreate(root_data_node, C.CString("global/afi-safis/afi-safi[afi-safi-name=\"openconfig-bgp-types:IPV6_UNICAST\"]/config/enabled"), C.CString("true"))

    res = C.DataNodeCreate(root_data_node, C.CString("global/afi-safis/afi-safi[afi-safi-name=\"openconfig-bgp-types:IPV4_UNICAST\"]/config/afi-safi-name"), C.CString("openconfig-bgp-types:IPV4_UNICAST"))
    res = C.DataNodeCreate(root_data_node, C.CString("global/afi-safis/afi-safi[afi-safi-name=\"openconfig-bgp-types:IPV4_UNICAST\"]/config/enabled"), C.CString("true"))

    if(res == nil) {
        panic("Could not create data!")
    }

    /*
    HACK!!!
     */

    //fmt.Println("gooo ", root_path.Path, root_data_node)
	//if IsSet(entity.GetOperation()) {
	//    p1 := C.CString(entity.GetOperation())
	//    defer C.free(unsafe.Pointer(p1))
	//    //C.DataNodeAddAnnotation(p1, root_data_node)
	//}

	//populateNameValues(root_data_node, root_path) TODO
	//walkChildren(entity, root_data_node)
	return root_data_node
}

func walkChildren(entity types.Entity, data_node C.DataNode) {
	children := entity.GetChildren()

    //fmt.Println("yikes. getchildren!1", len(children))

	for child_name := range children {

		if children[child_name].HasOperation() || children[child_name].HasData() {
			populateDataNode(children[child_name], data_node)
		}
	}
}

func populateDataNode(entity types.Entity, parent_data_node C.DataNode) {
	path := entity.GetEntityPath(entity.GetParent())
	p := C.CString(path.Path)
	defer C.free(unsafe.Pointer(p))
	ep := C.CString("")
	defer C.free(unsafe.Pointer(ep))

	data_node := C.DataNodeCreate(parent_data_node, p, ep)

    //fmt.Println("yeee populate", path.Path)

	//if(IsSet(entity.GetOperation())) {
	//    p1 := C.CString(entity.GetOperation())
	//    defer C.free(unsafe.Pointer(p1))
	//    //C.DataNodeAddAnnotation(p1, data_node)
	//}

	populateNameValues(data_node, path)
	walkChildren(entity, data_node)
}

func populateNameValues(data_node C.DataNode, path types.EntityPath) {
	for _, name_value := range path.ValuePaths {
		//var result C.DataNode
		leaf_data := name_value.Data
		p := C.CString(name_value.Name)

		if leaf_data.IsSet {
			p1 := C.CString(leaf_data.Value)
			//result = C.DataNodeCreate(data_node, p, p1)
			C.DataNodeCreate(data_node, p, p1)
			C.free(unsafe.Pointer(p1))
		}

		//if(IsSet(leaf_data.Operation)) {
		//    p1 := C.CString(name_value.Data.Operation)
		//    defer C.free(unsafe.Pointer(p1))
		//    //C.DataNodeAddAnnotation(p1, result)
		//}
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

	//c_children := C.DataNodeGetChildren(node)
	//children := []C.DataNode{unsafe.Pointer(c_children.datanodes)}
	children := []C.DataNode{}
	//real_data_node := node.private.(C.DataNode)
	//c_children := C.DataNodeGetChildren(real_data_node)

	for _, child_data_node := range children {
		child_name := C.DataNodeGetArgument(child_data_node)

		if dataNodeIsLeaf(child_data_node) {
			//C.DataNodeGetValue(child_data_node), C.DataNodeGetPath(node)
			entity.SetValue(C.GoString(child_name), C.GoString(C.DataNodeGetValue(child_data_node)))

		} else {

			var child_entity types.Entity
			if dataNodeIsList(child_data_node) {
				child_entity = entity.GetChildByName(C.GoString(child_name),
					getSegmentPath(C.GoString(
						C.DataNodeGetPath(child_data_node))))
			} else {
				child_entity = entity.GetChildByName(C.GoString(child_name), "")
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

func getSegmentPath(path string) string {
	segments := segmentalize(path)
	return segments[len(segments)-1]
}

func segmentalize(path string) []string {
	return make([]string, 2)
}
