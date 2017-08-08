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
package types

import (
	"reflect"
	"sort"
)

type YFilter int

const (
	NotSet YFilter = iota
	Read
	Merge
	Create
	Remove
	Delete
	Replace
)

type Empty struct {
	set bool
}

type LeafData struct {
	Value  string
	Filter YFilter
	IsSet  bool
}

type NameLeafData struct {
	Name string
	Data LeafData
}

type nameLeafDataList []NameLeafData

func (p nameLeafDataList) Len() int {
	return len(p)
}

func (p nameLeafDataList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p nameLeafDataList) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

type EntityPath struct {
	Path       string
	ValuePaths []NameLeafData
}

type AugmentCapabilitiesFunction func() map[string]string

type Entity interface {
	GetEntityPath(Entity) EntityPath
	GetSegmentPath() string

	HasDataOrFilter() bool

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

	GetFilter() YFilter
}

type Bits map[string]bool

type Decimal64 struct {
	value string
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
	name string

	leaf_type  YType
	bits_value Bits

	Value  string
	IsSet  bool
	Filter YFilter
}

func (y *YLeaf) GetNameLeafdata() NameLeafData {
	return NameLeafData{y.name, LeafData{y.Value, y.Filter, y.IsSet}}
}

type YLeafList struct {
	name   string
	values []YLeaf

	Filter    YFilter
	leaf_type YType
}

func (y *YLeafList) GetYLeafs() []YLeaf {
	return y.values
}

func (y *YLeafList) GetNameLeafdata() [](NameLeafData) {
	result := make([]NameLeafData, len(y.values))
	for i := 0; i < len(y.values); i++ {
		result = append(result, NameLeafData{y.values[i].name, LeafData{y.values[i].Value, y.values[i].Filter, y.values[i].IsSet}})
	}
	return result
}

type EncodingFormat int

const (
	XML EncodingFormat = iota
	JSON
)

func (e YFilter) String() string {
	switch e {
	case Read:
		return "read"
	case Replace:
		return "replace"
	case Delete:
		return "delete"
	case Merge:
		return "merge"
	case Create:
		return "create"
    case Remove:
        return "remove"
	case NotSet:
		return ""
	}
	return ""
}

type ServiceProvider interface {
	GetPrivate() interface{}
	Connect()
	Disconnect()
}

type CodecServiceProvider interface {
	Initialize(Entity)
	GetEncoding() EncodingFormat
	GetRootSchemaNode(Entity) RootSchemaNode
}

type Protocol int

const (
	Restconf Protocol = iota
	Netconf
)

type DataNode struct {
	Private interface{}
}

type RootSchemaNode struct {
	Private interface{}
}

type CServiceProvider struct {
	Private interface{}
}

type Repository struct {
	Path    string
	Private interface{}
}

//////////////////////////////////////////////////////////////////////////
// Exported utility functions
//////////////////////////////////////////////////////////////////////////

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

	if ancestor == nil {
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

	parents = sort.Reverse(parents).(EntitySlice)

	p = nil
	for _, p1 := range parents {
		if p != nil {
			path_buffer += "/"
		} else {
			p = p1
		}
		path_buffer += p1.GetSegmentPath()
	}
	if p != nil {
		path_buffer += "/"
	}
	path_buffer += current_node.GetSegmentPath()
	return path_buffer

}

func IsSet(Filter YFilter) bool {
	return Filter != NotSet
}

func sortValuePaths(v []NameLeafData) []NameLeafData {
	ret := make([]NameLeafData, 0)
	for _, v := range v {
		ret = append(ret, v)
	}
	sort.Sort(nameLeafDataList(ret))
	return ret
}

func nameValuesEqual(e1, e2 Entity) bool {
	value_path1 := e1.GetEntityPath(e1.GetParent()).ValuePaths
	value_path2 := e2.GetEntityPath(e2.GetParent()).ValuePaths
	path1 := sortValuePaths(value_path1)
	path2 := sortValuePaths(value_path2)

	if len(path1) != len(path2) {
		return false
	}

	ret := true
	for k := range path1 {
		name1 := path1[k].Name
		value1 := path1[k].Data
		name2 := path2[k].Name
		value2 := path2[k].Data

		if name1 != name2 || !reflect.DeepEqual(value1, value2) {
			ret = false
			break
		}
	}
	return ret
}

func deepValueEqual(e1, e2 Entity) bool {
	if e1 == nil && e2 == nil {
		return false
	}
	children1 := e1.GetChildren()
	children2 := e2.GetChildren()

	marker := make(map[string]bool)

	ret := true
	for k, c1 := range children1 {
		marker[k] = true
		if c1.HasDataOrFilter() {
			c2, ok := children2[k]
			if ok && deepValueEqual(c1, c2) {
				ret = ret && nameValuesEqual(c1, c2)
			} else {
				ret = false
				break
			}
		}
	}

	for k := range children2 {
		_, ok := marker[k]
		if !ok {
			ret = false
			break
		}
	}

	return ret
}

func EntityEqual(x, y Entity) bool {
	if x == nil && y == nil {
		return x == y
	}
	return deepValueEqual(x, y)
}
