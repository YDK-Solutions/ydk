// Package types provides built-in types specified in
// YANG RFC 6020 and types used in YDK Go APIs.
//
// YANG Development Kit Copyright 2017 Cisco Systems. All rights reserved.
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package types

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	encoding "github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format"
	"github.com/CiscoDevNet/ydk-go/ydk/errors"
	"github.com/CiscoDevNet/ydk-go/ydk/types/yfilter"
	"github.com/CiscoDevNet/ydk-go/ydk/types/ytype"
)

// Empty represents a YANG built-in Empty type
type Empty struct {
}

// String returns the representation of the Empty type as a string (an empty string)
func (e *Empty) String() string {
	return ""
}

// LeafData represents the data contained in a YANG leaf
type LeafData struct {
	Value  string
	Filter yfilter.YFilter
	IsSet  bool
}

// NameLeafData represents a YANG leaf to which a name and data can be assigned
type NameLeafData struct {
	Name string
	Data LeafData
}

// nameLeafDataList represents a YANG leaf-list
type nameLeafDataList []NameLeafData

// Len returns the length (int) of a given nameLeafDataList
func (p nameLeafDataList) Len() int {
	return len(p)
}

// Swap swaps the NameLeafData at indices i and j of the given nameLeafDataList
func (p nameLeafDataList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less returns whether the name of the NameLeafData at index i is less than the one at index j of a given nameLeafDataList
func (p nameLeafDataList) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

// EntityPath
type EntityPath struct {
	Path       string
	ValuePaths []NameLeafData
}

type CommonEntityData struct {
	YangName 		string
	BundleName 		string
	YFilter 		yfilter.YFilter
	Parent 			Entity
	ParentYangName 	string
	GoName 			string
	Children 		map[string]Entity
	Leafs 			map[string]interface{}
}

// Entity is a basic type that represents containers in YANG
type Entity interface {
	GetCommonEntityData()			*CommonEntityData

	GetGoName(string)				string
	GetSegmentPath() 				string
	GetChildByName(string, string) 	Entity

	GetChildren() 					map[string]Entity

	GetCapabilitiesTable() 			map[string]string
	GetNamespaceTable() 			map[string]string
	GetBundleYangModelsLocation() 	string

	GetFilter() 					yfilter.YFilter
}

// Bits is a basic type that represents the YANG bits type
type Bits map[string]bool

type BitsList struct {
	Value []map[string]bool
}


/////////////////////////////////////
// CommonEntityData Utility Functions -- not really needed
/////////////////////////////////////

// GetYangName returns the given entity's YANG name
func GetYangName(entity Entity) string {
	data := entity.GetCommonEntityData()
	return data.YangName
}

// GetBundleName returns the name of the bundle that the given entity belongs to
func GetBundleName(entity Entity) string {
	data := entity.GetCommonEntityData()
	return data.BundleName
}

// GetParentYangName returns the YANG name of the parent of the given entity
func GetParentYangName(entity Entity) string {
	data := entity.GetCommonEntityData()
	return data.ParentYangName
}

// GetLeafs returns a map of the leafs contained in the given entity
func GetLeafs(entity Entity) map[string]interface{} {
	data := entity.GetCommonEntityData()
	return data.Leafs
}


/////////////////////////////////////
// Entity Utility Functions
/////////////////////////////////////

// HasDataOrFilter returns a bool representing whether the entity or any of its children have their data/filter set
func HasDataOrFilter(entity Entity) bool {
	if (entity.GetFilter() != yfilter.NotSet) {
		return true
	}

	children := entity.GetChildren()
	leafs := GetLeafs(entity)

	// children
	for _, child := range children {
		if (child.GetFilter() != yfilter.NotSet || HasDataOrFilter(child)) {
			return true
		}
	}

	v := reflect.ValueOf(entity).Elem()

	// checking leafs
	for name, leaf := range leafs {
		goName := entity.GetGoName(name)
		field := v.FieldByName(goName)

		if field.Kind() != reflect.Slice {
			if leaf != nil { return true }
		} else {
			for _, l := range field.Interface().([]interface{}) {
				if l != nil { return true }
			}
		}
	}

	return false
}

// GetEntityPath returns an EntityPath struct for the given entity
func GetEntityPath(entity Entity) EntityPath {
	entityPath := EntityPath{Path: entity.GetSegmentPath()}
	leafs := GetLeafs(entity)
	v := reflect.ValueOf(entity).Elem()

	// leafs
	var leafData LeafData
	for name, leaf := range leafs {
		goName := entity.GetGoName(name)
		field := v.FieldByName(goName)

		if leaf != nil && field.Kind() != reflect.Slice {
			switch leaf.(type) {
			case yfilter.YFilter:
				// yfilter
				leafData = LeafData{
					IsSet: true, Filter: leaf.(yfilter.YFilter)}
			case map[string]bool:
				// bits
				var used_bits []string
				for bit, enabled := range(leaf.(map[string]bool)) {
					if enabled {
						used_bits = append(used_bits, bit)
					}
				}
				v := strings.Join(used_bits, " ")
				leafData = LeafData{IsSet: true, Value: v}
			default:
				var v string
				if reflect.TypeOf(leaf) != reflect.TypeOf(Empty{}) {
					v = fmt.Sprintf("%v", leafs[name])
				}
				leafData = LeafData{
					IsSet: true, Value: v}
			}
			entityPath.ValuePaths = append(
				entityPath.ValuePaths,
				NameLeafData{Name: name, Data: leafData})
		}
	}

	return entityPath
}

// todo: need to instantiate child if ok == false
// GetChildByName takes an Entity and returns the child Entity described by the given childYangName and segmentPath
func GetChildByName(
	entity Entity,
	childYangName string,
	segmentPath string) Entity {

	children := entity.GetChildren()

	// child
	if child, ok := children[childYangName]; ok {
		return child
	}

	// todo: instantiate the child
	// need a way to retrieve type from childYangName and/or segmentPath
	// (ie bake in a map going from string to type in apis)
	return nil
}

// SetValue sets the leaf value for given entity, valuePath, and value args
func SetValue(entity Entity, valuePath string, value interface{}) {
	goName := entity.GetGoName(valuePath)

	s := reflect.ValueOf(entity).Elem()
	v := s.FieldByName(goName)
	if v.IsValid() {
		if v.Type() == reflect.TypeOf(make(map[string]bool)) {
			bits := v.Interface().(map[string]bool)
			bits[value.(string)] = true

			v.Set(reflect.ValueOf(bits))
		} else if v.Type() == reflect.TypeOf(BitsList{}) {
			bitsValue := make(map[string]bool)
			bitsValue[value.(string)] = true
			
			bitslist := v.Interface().(BitsList)
			bitslist.Value = append(bitslist.Value, bitsValue)

			v.Set(reflect.ValueOf(bitslist))
		} else if v.Kind() == reflect.Slice {
			v.Set(reflect.Append(v, reflect.ValueOf(value)))
		} else {
			v.Set(reflect.ValueOf(value))
		}
	}
}


// Decimal64 represents a YANG built-in Decimal64 type
type Decimal64 struct {
	value string
}

// EnumYLeaf represents variable data
type EnumYLeaf struct {
	value int
	name  string
}

// Enum represents a YANG built-in enum type, a base type for all YDK enums.
type Enum struct {
	EnumYLeaf
}

// YLeaf represents a YANG leaf to which data can be assigned.
type YLeaf struct {
	name string

	leafType  	ytype.YType
	bitsValue 	Bits

	Value  	string
	IsSet  	bool
	Filter 	yfilter.YFilter
}

// GetNameLeafdata instantiates and returns NameLeafData type for this leaf
func (y *YLeaf) GetNameLeafdata() NameLeafData {
	return NameLeafData{y.name, LeafData{y.Value, y.Filter, y.IsSet}}
}

// YLeafList represents a YANG leaf-list to which multiple instances of data can be appended
type YLeafList struct {
	name   	string
	values 	[]YLeaf

	Filter    	yfilter.YFilter
	leafType 	ytype.YType
}

// GetYLeafs is a getter function for YLeafList values
func (y *YLeafList) GetYLeafs() []YLeaf {
	return y.values
}

// GetNameLeafdata instantiates and returns name NameLeafData for this YLeafList
func (y *YLeafList) GetNameLeafdata() [](NameLeafData) {
	result := make([]NameLeafData, len(y.values))
	for i := 0; i < len(y.values); i++ {
		result = append(result, NameLeafData{y.values[i].name, LeafData{y.values[i].Value, y.values[i].Filter, y.values[i].IsSet}})
	}
	return result
}

// ServiceProvider
type ServiceProvider interface {
	GetPrivate() interface{}
	Connect()
	Disconnect()
	GetState() *errors.State
}

// CodecServiceProvider
type CodecServiceProvider interface {
	Initialize(Entity)
	GetEncoding() encoding.EncodingFormat
	GetRootSchemaNode(Entity) RootSchemaNode
	GetState() *errors.State
}

// DataNode represents a containment hierarchy
type DataNode struct {
	Private interface{}
}

// RootSchemaNode represents the root of the SchemaTree. 
// It can be used to instantiate a DataNode tree or an Rpc object. 
// The children of the RootSchemaNode represent the top level SchemaNode in the YANG module submodules.
type RootSchemaNode struct {
	Private interface{}
}

// CServiceProvider
type CServiceProvider struct {
	Private interface{}
}

// COpenDaylightServiceProvider is a service provider to be used to communicate with an OpenDaylight instance.
type COpenDaylightServiceProvider struct {
	Private interface{}
}

// Repository represents the Repository of YANG models.
// A instance of the Repository will be used to create a RootSchemaNode given a set of ©pabilities.
// Behind the scenes the repository is responsible for loading and parsing the YANG modules and creating the SchemaNode tree. 
type Repository struct {
	Path    string
	Private interface{}
}

//////////////////////////////////////////////////////////////////////////
// Exported utility functions
//////////////////////////////////////////////////////////////////////////

// EntitySlice is a slice of entities
type EntitySlice []Entity

// Len returns the length of given EntitySlice
func (s EntitySlice) Len() int {
	return len(s)
}

// Less returns whether the Entity at index i is less than the one at index j of the given EntitySlice
func (s EntitySlice) Less(i, j int) bool {
	return s[i].GetSegmentPath() < s[j].GetSegmentPath()
}

// Swap swaps the Entities at indices i and j of the given EntitySlice
func (s EntitySlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// GetRelativeEntityPath returns the relative entity path (string)
func GetRelativeEntityPath(current_node Entity, ancestor Entity, path string) string {
	path_buffer := path

	if ancestor == nil {
		return ""
	}
	p := current_node.GetCommonEntityData().Parent
	parents := EntitySlice{}
	for p != nil && p != ancestor {
		//append(parents, p)
		p = p.GetCommonEntityData().Parent
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

// IsSet returns whether the given filter is set or not
func IsSet(Filter yfilter.YFilter) bool {
	return Filter != yfilter.NotSet
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
	valuePath1 := GetEntityPath(e1).ValuePaths
	valuePath2 := GetEntityPath(e2).ValuePaths
	path1 := sortValuePaths(valuePath1)
	path2 := sortValuePaths(valuePath2)

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
		if HasDataOrFilter(c1) {
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

// EntityEqual returns whether the entities x and y and their children are equal in value
func EntityEqual(x, y Entity) bool {
	if x == nil && y == nil {
		return x == y
	}
	return deepValueEqual(x, y)
}
