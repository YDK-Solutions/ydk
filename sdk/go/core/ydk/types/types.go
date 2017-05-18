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

import (
    "fmt"
    "sort"
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

type AugmentCapabilitiesFunction func() map[string]string

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

	leaf_type  YType
	bits_value Bits

    Value string
    IsSet      bool
    Operation  EditOperation
}

func (y *YLeaf) GetNameLeafdata() NameLeafData {
	return NameLeafData{y.name, LeafData{y.Value, y.Operation, y.IsSet}}
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
		result = append(result, NameLeafData{y.values[i].name, LeafData{y.values[i].Value, y.values[i].Operation, y.values[i].IsSet}})
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

type ServiceProvider interface {
    GetPrivate() interface{}
    Connect()
    Disconnect()
}

type CodecServiceProvider struct {
    Encoding EncodingFormat
}

type Protocol int

const (
	Restconf Protocol = iota
	Netconf
)

type DataNode struct {
    Private interface{}
}

type CServiceProvider struct {
    Private interface{}
}

type Repository struct {
    Path string
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

