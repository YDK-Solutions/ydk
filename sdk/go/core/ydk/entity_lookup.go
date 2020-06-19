/*  ----------------------------------------------------------------
 YDK - YANG Development Kit
 Copyright 2016 Cisco Systems. All rights reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 -------------------------------------------------------------------
 This file has been modified by Yan Gorelik, YDK Solutions.
 All modifications in original under CiscoDevNet domain
 introduced since October 2019 are copyrighted.
 All rights reserved under Apache License, Version 2.0.
 ------------------------------------------------------------------*/

package ydk

import (
	"fmt"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"reflect"
)

var topEntityRegistry = make(map[string]reflect.Type)

// RegisterEntity
func RegisterEntity(name string, entity_type reflect.Type) {
	topEntityRegistry[name] = entity_type
}

// GetTopEntity
func GetTopEntity(name string) (types.Entity, bool) {
	_, ok := topEntityRegistry[name]
	if !ok {
		YLogError(fmt.Sprintf("Entity '%s' is not registered. Please import corresponding package to your application.", name))
		return nil, ok
	}
	return reflect.New(topEntityRegistry[name]).Elem().Addr().Interface().(types.Entity), ok
}
