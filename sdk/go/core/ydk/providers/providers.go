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
package providers

// #cgo CXXFLAGS: -g -std=c++11
// #cgo LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #include <ydk/ydk.h>
// #include <stdlib.h>
import "C"

import (
    "fmt"
    "unsafe"
)

type Repository struct {
	Path string
}

type ServiceProvider interface {
    GetPrivate() C.ServiceProvider
	Init()
    Close()
}

type NetconfServiceProvider struct {
	Repo     Repository
	Address  string
	Username string
	Password string
	Port     int

	private C.ServiceProvider
}

type RestconfServiceProvider struct {
}

func (provider *NetconfServiceProvider) GetPrivate() C.ServiceProvider {
    return provider.private
}

func (provider *NetconfServiceProvider) Init() {
	var address *C.char = C.CString(provider.Address)
	defer C.free(unsafe.Pointer(address))
    var username *C.char = C.CString(provider.Username)
    defer C.free(unsafe.Pointer(username))
    var password *C.char = C.CString(provider.Password)
    defer C.free(unsafe.Pointer(password))
    var port C.int = C.int(provider.Port)

	if len(provider.Repo.Path) > 0 {
        var path *C.char = C.CString(provider.Repo.Path)
        repo := C.RepositoryInitWithPath(path)
        provider.private = C.NetconfServiceProviderInitWithRepo(repo, address, username, password, port)
	} else {
        provider.private = C.NetconfServiceProviderInit(address, username, password, port)
	}
    if provider.private == nil {
        fmt.Println("Could not connect")
    }
}

func (provider *NetconfServiceProvider) Close() {
	C.NetconfServiceProviderFree(provider.private)
}
