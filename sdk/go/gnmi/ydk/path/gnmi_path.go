// Package providers implements support for Go service providers.
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
package path

// #cgo linux LDFLAGS:  -lydk_gnmi -lgrpc++ -lprotobuf
// #include <ydk/ydk.h>
// #include <ydk/ydk_gnmi.h>
// #include <stdlib.h>
import "C"

import (
//	"fmt"
//	"reflect"
	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/errors"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"unsafe"
)

// GnmiSession declaration and methods
//
type GnmiSession struct {
	Repo        types.Repository
	Address     string
	Username    string
	Password    string
	Port        int
	ServerCert  string
	PrivateKey  string

	Private		interface{}
	State		errors.State
}

// Connect to gNMI Session
func (gs *GnmiSession) Connect() {

	var caddress *C.char = C.CString(gs.Address)
	defer C.free(unsafe.Pointer(caddress))
	if gs.Port == 0 {
		gs.Port = 57400
	}
	var cport C.int = C.int(gs.Port)

	var cserver *C.char = C.CString(gs.ServerCert)
	defer C.free(unsafe.Pointer(cserver))

	var cclient *C.char = C.CString(gs.PrivateKey)
	defer C.free(unsafe.Pointer(cclient))

	AddCState(&gs.State)
	cstate := GetCState(&gs.State)
	
	var repopath *C.char = C.CString(gs.Repo.Path)
	defer C.free(unsafe.Pointer(repopath))
	repo := C.RepositoryInitWithPath( *cstate, repopath)
	PanicOnCStateError(cstate)

	var cusername *C.char = C.CString(gs.Username)
	defer C.free(unsafe.Pointer(cusername))
	var cpassword *C.char = C.CString(gs.Password)
	defer C.free(unsafe.Pointer(cpassword))
	
	gs.Private = C.GnmiSessionInit( *cstate, repo, caddress, cport, cusername, cpassword, cserver, cclient);
	PanicOnCStateError(cstate)
}

// Disconnect from gNMI Session
func (gs *GnmiSession) Disconnect() {
	if gs.Private == nil {
		return
	}
	realSession := gs.Private.(C.GnmiSession)
	C.GnmiSessionFree(realSession)
	CleanUpErrorState(&gs.State)
}

func (gs *GnmiSession) GetRootSchemaNode() types.RootSchemaNode {
	cstate := GetCState(&gs.State)
	
	realSession := gs.Private.(C.GnmiSession)

	var rootSchema C.RootSchemaWrapper = C.GnmiSessionGetRootSchemaNode( *cstate, realSession)
	PanicOnCStateError(cstate)
	//ydkpath.PanicOnStateError(gs.State)
	if rootSchema == nil {
        ydk.YLogError("Root schema is nil!")
		panic(1)
	}

	rsn := types.RootSchemaNode{Private: rootSchema}
	return rsn
}

func (gs *GnmiSession) ExecuteRpc(rpc types.Rpc) types.DataNode {
	cstate := GetCState(&gs.State)

	csession := gs.Private.(C.GnmiSession)
	crpc := rpc.Private.(C.Rpc)

	cdn := C.GnmiSessionExecuteRpc( *cstate, csession, crpc)
	PanicOnCStateError(cstate)
	
	dn := types.DataNode{Private: cdn}
	return dn
}

func (gs *GnmiSession) ExecuteSubscribeRpc(rpc types.Rpc) {
	cstate := GetCState(&gs.State)

	csession := gs.Private.(C.GnmiSession)
	crpc := rpc.Private.(C.Rpc)

	C.GnmiSessionExecuteSubscribeRpc( *cstate, csession, crpc)
	PanicOnCStateError(cstate)
}

// Utility functions
//
func (gs *GnmiSession) SubscribeInProgress() bool {
	realSession := gs.Private.(C.GnmiSession)
	cstate := GetCState(&gs.State)
	
	var cresponse C.boolean = C.GnmiSessionSubscribeInProgress( *cstate, realSession)
	PanicOnCStateError(cstate)
	var response bool = false
	if cresponse == 1 {
		response = true;
	}
	return response
}

func (gs *GnmiSession) GetLastSubscribeResponse(previousResponse string) string {
	realSession := gs.Private.(C.GnmiSession)
	cstate := GetCState(&gs.State)
	
	var cprevious *C.char = C.CString(previousResponse)
	defer C.free(unsafe.Pointer(cprevious))
	
	cresponse := C.GetLastSubscribeResponse( *cstate, realSession, cprevious)
	return C.GoString(cresponse)
}

// gNMI Service Provider C-wrapper functions
//
func GnmiServiceProviderConnect(
	state *errors.State,
	repo types.Repository,
	address string, port int,
	username, password string,
	serverCert, privateKey string) types.CServiceProvider {

	AddCState(state)
	cstate := GetCState(state)

	var repopath *C.char = C.CString(repo.Path)
	crepo := C.RepositoryInitWithPath(*cstate, repopath)
	PanicOnCStateError(cstate)

	var caddress *C.char = C.CString(address)
	defer C.free(unsafe.Pointer(caddress))
	
	var cport C.int = C.int(port)

	var p C.ServiceProvider
	var cusername *C.char = C.CString(username)
	defer C.free(unsafe.Pointer(cusername))
	var cpassword *C.char = C.CString(password)
	defer C.free(unsafe.Pointer(cpassword))

	var cserver *C.char = C.CString(serverCert)
	defer C.free(unsafe.Pointer(cserver))
	var cclient *C.char = C.CString(privateKey)
	defer C.free(unsafe.Pointer(cclient))
	
	p = C.GnmiServiceProviderInit( *cstate, crepo, caddress, cport, cusername, cpassword, cserver, cclient);
	PanicOnCStateError(cstate)

	cprovider := types.CServiceProvider{Private: p}
	return cprovider
}

func GnmiServiceProviderDisconnect(provider types.CServiceProvider) {
	realProvider := provider.Private.(C.ServiceProvider)
	C.GnmiServiceProviderFree(realProvider)
}

func GnmiServiceProviderGetSession(provider types.CServiceProvider) *GnmiSession {

	var state errors.State
	AddCState(&state)
	cstate := GetCState(&state)

	realProvider := provider.Private.(C.ServiceProvider)
	realSession := C.GnmiServiceProviderGetSession( *cstate, realProvider)
	PanicOnCStateError(cstate)
	
	gnmiSession := GnmiSession{Private: realSession, State: state}
	// TODO. Possibly need to populate the rest of structure members

	return &gnmiSession
}

// GnmiService C-wrapper functions
//
func GnmiServiceGetCapabilities(provider types.CServiceProvider) string {
	realProvider := provider.Private.(C.ServiceProvider)
	var state errors.State
	AddCState(&state)
	cstate := GetCState(&state)

	var ccaps *C.char = C.GnmiServiceGetCapabilities( *cstate, realProvider)
	PanicOnCStateError(cstate)

	return C.GoString(ccaps)
}

