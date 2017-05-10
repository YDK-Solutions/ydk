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
//#include <ydk/ydk.h>
import "C"

import (
    "fmt"
    "github.com/CiscoDevNet/ydk-go/models/openconfig"
    "github.com/CiscoDevNet/ydk-go/providers"
    "github.com/CiscoDevNet/ydk-go/services"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()

    C.EnableLogging()

    var provider providers.NetconfServiceProvider = providers.NetconfServiceProvider{
                                                Repo:providers.Repository{"/usr/local/share/ydktest@0.1.0/"},
                                                Address:"localhost",
                                                Username: "admin",
                                                Password: "admin",
                                                Port: 12022}

    provider.Init()

    bgp := openconfigbgp.Bgp{}
    bgp.Global.Config.As = 65001
    bgp.Global.Config.RouterId = "1.2.3.4"

    crud := services.CrudService{}
    result := crud.Create(provider, bgp)

    if result == true {
        fmt.Println("Operation succeeded!")
    } else {
        fmt.Println("Operation failed!")
    }

    provider.Close()
}