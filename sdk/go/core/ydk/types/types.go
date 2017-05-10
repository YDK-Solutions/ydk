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
 * "License") you may not use this file except in compliance
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
package types

// #cgo CXXFLAGS: -g -std=c++11
// #cgo LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #include <ydk/ydk.h>
// #include <stdlib.h>
import "C"

import (
    "fmt"
    "sort"
    "unsafe"
)

type EditOperation int

const (
	Merge EditOperation = iota
	Create
	Remove
	Delete_
	Replace
	NotSet
)

type Empty struct {
	set bool
}

type LeafData struct {
	Value     string
	Operation EditOperation
	IsSet     bool
}

type NameLeafData struct {
	Name string
	Data LeafData
}

type EntityPath struct {
	Path       string
	ValuePaths []NameLeafData
}

type AugmentCapabilitiesFunction func()

type Entity interface {
	GetEntityPath(Entity) EntityPath
	GetSegmentPath() string

	HasData() bool
	HasOperation() bool

	SetValue(string, string)
	GetChildByName(string, string) Entity

	GetChildren() map[string]Entity

	SetParent(Entity)
	GetParent() Entity

	GetAugmentCapabilitiesFunction() AugmentCapabilitiesFunction
	GetBundleYangModelsLocation() string
	GetBundleName() string

    GetYangName() string
    GetParentYangName() string

    GetOperation() EditOperation
}

type Bits map[string]bool

type Decimal64 struct {
	value string
}

type Identity interface {
	to_string() string
}

type EnumYLeaf struct {
	value int
	name  string
}

type Enum struct {
	EnumYLeaf
}

type YType int

const (
	Uint_8 YType = iota
	Uint_16
	Uint_32
	Uint_64
	Int_8
	Int_16
	Int_32
	Int_64
	Empty_
	Identityref
	Str
	Boolean
	Enumeration
	Bits_
	Decimal64_
)

type YLeaf struct {
	name  string
	value string

	is_set     bool
	operation  EditOperation
	leaf_type  YType
	bits_value Bits
}

func (y *YLeaf) Get() string {
	return y.value
}

func (y *YLeaf) GetNameLeafdata() NameLeafData {
	return NameLeafData{y.name, LeafData{y.value, y.operation, y.is_set}}
}

type YLeafList struct {
	name   string
	values []YLeaf

	operation EditOperation
	leaf_type YType
}

func (y *YLeafList) GetYLeafs() []YLeaf {
	return y.values
}

func (y *YLeafList) GetNameLeafdata() [](NameLeafData) {
	result := make([]NameLeafData, len(y.values))
	for i := 0; i < len(y.values); i++ {
		result = append(result, NameLeafData{y.values[i].name, LeafData{y.values[i].value, y.values[i].operation, y.values[i].is_set}})
	}
	return result
}

type EncodingFormat int

const (
	XML EncodingFormat = iota
	JSON
)

func (e EditOperation) String() string {
	return fmt.Sprintf("%v", e)
}

type Protocol int

const (
	Restconf Protocol = iota
	Netconf
)

//////////////////////////////////////////////////////////////////////////
// Exported utility functions
//////////////////////////////////////////////////////////////////////////

func segmentalize(path string) []string {
    return make([]string, 2)
}

type EntitySlice []Entity

func (s EntitySlice) Len() int {
    return len(s)
}

func (s EntitySlice) Less(i, j int) bool {
    return s[i].GetSegmentPath() < s[j].GetSegmentPath()
}

func (s EntitySlice) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func GetRelativeEntityPath(current_node Entity, ancestor Entity, path string) string {
    path_buffer := path

    if(ancestor == nil) {
        return ""
    }
    p := current_node.GetParent()
    parents := EntitySlice{}
    for p != nil && p != ancestor {
        //append(parents, p)
        p = p.GetParent()
    }

    if p == nil {
        return ""
    }

    sort.Reverse(parents)

    p = nil
    for _, p1 := range parents {
    if p!=nil {
        path_buffer += "/"
    } else {
        p = p1
    }
        path_buffer += p1.GetSegmentPath()
    }
    if(p != nil) {
        path_buffer += "/"
    }
    path_buffer += current_node.GetSegmentPath()
    return path_buffer

}

func IsSet(operation EditOperation) bool {
    return operation != NotSet
}

//////////////////////////////////////////////////////////////////////////
// DataNode from Entity
//////////////////////////////////////////////////////////////////////////
func GetDataNodeFromEntity(entity Entity, root_schema C.RootSchemaNode)  C.DataNode {
    root_path := entity.GetEntityPath(nil)
    path := C.CString(root_path.Path)
    defer C.free(unsafe.Pointer(path))

    root_data_node := C.RootSchemaNodeCreate(root_schema, path)
    //if IsSet(entity.GetOperation()) {
    //    p1 := C.CString(entity.GetOperation())
    //    defer C.free(unsafe.Pointer(p1))
    //    //C.DataNodeAddAnnotation(p1, root_data_node)
    //}

    populateNameValues(root_data_node, root_path)
    walkChildren(entity, root_data_node)
    return root_data_node
}

func walkChildren(entity Entity, data_node C.DataNode) {
    children := entity.GetChildren()

    for child_name := range children {

        if (children[child_name].HasOperation() || children[child_name].HasData()) {
            populateDataNode(children[child_name], data_node)
        }
    }
}

func populateDataNode(entity Entity, parent_data_node C.DataNode) {
    path := entity.GetEntityPath(entity.GetParent())
    p := C.CString(path.Path)
    defer C.free(unsafe.Pointer(p))
    ep := C.CString("")
    defer C.free(unsafe.Pointer(ep))

    data_node := C.DataNodeCreate(parent_data_node, p, ep)

    //if(IsSet(entity.GetOperation())) {
    //    p1 := C.CString(entity.GetOperation())
    //    defer C.free(unsafe.Pointer(p1))
    //    //C.DataNodeAddAnnotation(p1, data_node)
    //}

    populateNameValues(data_node, path)
    walkChildren(entity, data_node)
}

func populateNameValues(data_node C.DataNode, path EntityPath) {
    for _, name_value := range path.ValuePaths {
        //var result C.DataNode
        leaf_data := name_value.Data
        p := C.CString(name_value.Name)
        defer C.free(unsafe.Pointer(p))

        if(leaf_data.IsSet) {
            p1 := C.CString(leaf_data.Value)
            defer C.free(unsafe.Pointer(p1))
            //result = C.DataNodeCreate(data_node, p, p1)
            C.DataNodeCreate(data_node, p, p1)
        }

        //if(IsSet(leaf_data.Operation)) {
        //    p1 := C.CString(name_value.Data.Operation)
        //    defer C.free(unsafe.Pointer(p1))
        //    //C.DataNodeAddAnnotation(p1, result)
        //}
    }
}

//////////////////////////////////////////////////////////////////////////
// Entity from DataNode
//////////////////////////////////////////////////////////////////////////
func GetEntityFromDataNode(node C.DataNode, entity Entity) {

    if entity == nil || node == nil {
        return
    }

    //c_children := C.DataNodeGetChildren(node)
    //children := []C.DataNode{unsafe.Pointer(c_children.datanodes)}
    children := []C.DataNode{}

    for _, child_data_node := range children {
        child_name := C.DataNodeGetArgument(child_data_node)

        if(dataNodeIsLeaf(child_data_node)) {
            //C.DataNodeGetValue(child_data_node), C.DataNodeGetPath(node)
            entity.SetValue(C.GoString(child_name), C.GoString(C.DataNodeGetValue(child_data_node)))

        } else {

            var child_entity Entity
            if(dataNodeIsList(child_data_node)) {
                child_entity = entity.GetChildByName(C.GoString(child_name),
                                                        getSegmentPath(C.GoString(
                                                                        C.DataNodeGetPath(child_data_node))))
            } else {
                child_entity = entity.GetChildByName(C.GoString(child_name), "")
            }
            child_entity.SetParent(entity)
            GetEntityFromDataNode(child_data_node, child_entity)
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
