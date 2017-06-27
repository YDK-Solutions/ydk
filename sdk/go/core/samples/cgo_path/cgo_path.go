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
package main

// #cgo CXXFLAGS: -g -std=c++11
// #cgo LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #include <ydk/ydk.h>
// #include <stdlib.h>
import "C"

import (
	"flag"
	"fmt"
	"unsafe"
)

func main() {
	verbosePtr := flag.Bool("verbose", false, "Enable verbose")
	flag.Parse()

	if *verbosePtr {
		C.EnableLogging(C.INFO)
	}

	var address *C.char = C.CString("localhost")
	defer C.free(unsafe.Pointer(address))
	var username *C.char = C.CString("admin")
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString("admin")
	defer C.free(unsafe.Pointer(password))
	var path *C.char = C.CString("/usr/local/share/ydktest@0.1.0/")
	defer C.free(unsafe.Pointer(path))
	var runner_path *C.char = C.CString("ydktest-sanity:runner")
	defer C.free(unsafe.Pointer(runner_path))
	var number_path *C.char = C.CString("ytypes/built-in-t/number8")
	defer C.free(unsafe.Pointer(number_path))
	var number_value *C.char = C.CString("2")
	defer C.free(unsafe.Pointer(number_value))
	var create_path *C.char = C.CString("ydk:create")
	defer C.free(unsafe.Pointer(create_path))
	var read_path *C.char = C.CString("ydk:read")
	defer C.free(unsafe.Pointer(read_path))
	var entity_path *C.char = C.CString("entity")
	defer C.free(unsafe.Pointer(entity_path))
	var filter_path *C.char = C.CString("filter")
	defer C.free(unsafe.Pointer(filter_path))

	codec := C.CodecServiceInit()
	defer C.CodecServiceFree(codec)
	repo := C.RepositoryInitWithPath(path)
	defer C.RepositoryFree(repo)
	provider := C.NetconfServiceProviderInitWithRepo(repo, address, username, password, 12022)
	defer C.NetconfServiceProviderFree(provider)
	root_schema := C.ServiceProviderGetRootSchema(provider)

	runner := C.RootSchemaNodeCreate(root_schema, runner_path)

	C.DataNodeCreate(runner, number_path, number_value)
	var create_xml *C.char = C.CodecServiceEncode(codec, runner, C.XML, 0)
	defer C.free(unsafe.Pointer(create_xml))

	create_rpc := C.RootSchemaNodeRpc(root_schema, create_path)
	input := C.RpcInput(create_rpc)
	C.DataNodeCreate(input, entity_path, create_xml)
	C.RpcExecute(create_rpc, provider)

	read_rpc := C.RootSchemaNodeRpc(root_schema, read_path)
	input = C.RpcInput(read_rpc)
	runner_filter := C.RootSchemaNodeCreate(root_schema, runner_path)
	var read_xml *C.char = C.CodecServiceEncode(codec, runner_filter, C.XML, 0)
	defer C.free(unsafe.Pointer(read_xml))

	C.DataNodeCreate(input, filter_path, read_xml)
	read_data := C.RpcExecute(read_rpc, provider)
	if read_data == nil {
		fmt.Println("uhoh")
	}
	var data *C.char = C.CodecServiceEncode(codec, read_data, C.XML, 1)
	defer C.free(unsafe.Pointer(data))
	s := C.GoString(data)
	fmt.Println("Read data:", s)
}
