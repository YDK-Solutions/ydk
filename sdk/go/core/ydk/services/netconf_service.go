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
package services

// #cgo CXXFLAGS: -g -std=c++11
// #cgo darwin LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #cgo linux LDFLAGS:  -fprofile-arcs -ftest-coverage --coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lstdc++ -lpython2.7 -ldl
// #include <ydk/ydk.h>
// #include <stdlib.h>
import "C"

import (
    "encoding/json"
    "encoding/xml"
    "strconv"
    "fmt"
    "errors"
    "log"
    "github.com/CiscoDevNet/ydk-go/ydk"
    "github.com/CiscoDevNet/ydk-go/ydk/path"
    "github.com/CiscoDevNet/ydk-go/ydk/types"
)

const (
    Candidate int = iota
    Running
    StartUp
    Url
)

type NetconfService struct {
}

func (ns *NetconfService) CancelCommit(provider NetconfServiceProvider, persistId int) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    cancelCommitRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:cancel-commit")
    if (persistId > 1) {
        input := C.RpcInput(cstate, cancelCommitRpc)
        persistIdStr := strconv.Itoa(persistId)
        C.DataNodeCreate(cstate, input, "persist-id", C.CString(persistIdStr))
    }

    readData := C.RpcExecute(cstate, cancelCommitRpc, provider)
    return readData == nil
}

func (ns *NetconfService) CloseSession(provider NetconfServiceProvider) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    closeSessionRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:close-session")

    readData := C.RpcExecute(cstate, closeSessionRpc, provider)
    return readData == nil
}

func (ns *NetconfService) Commit(provider NetconfServiceProvider, confirmed bool, confirmTimeOut, persist, persistId int) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    commitRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:commit")

    if (confirmed) {
        input := C.RpcInput(cstate, commitRpc)
        C.DataNodeCreate(cstate, input, "confirmed", "")
    }

    if (confirmTimeOut > -1) {
        input := C.RpcInput(cstate, commitRpc)
        confirmTimeOutStr := strconv.Itoa(confirmTimeOut)
        C.DataNodeCreate(cstate, input, "confirm-timeout", confirmTimeOutStr)
    }

    if (persist > -1) {
        input := C.RpcInput(cstate, commitRpc)
        persistStr := strconv.Itoa(persist)
        C.DataNodeCreate(cstate, input, "persist", persistStr)
    }

    if (persistId > -1) {
        input := C.RpcInput(cstate, commitRpc)
        persistIdStr := strconv.Itoa(persistId)
        C.DataNodeCreate(cstate, input, "persist-id", persistIdStr)
    }

    readData := C.RpcExecute(cstate, commitRpc, provider)
    return readData == nil
}

// TODO: figure out optional arguments
// func (ns *NetconfService) CopyConfig(provider NetconfServiceProvider, target, source int, url string) bool {
// }

// func (ns *NetconfService) CopyConfig(provider NetconfServiceProvider, target int, source types.Entity) bool {   
// }

func (ns *NetconfService) DeleteConfig(provider NetconfServiceProvider, target int, url string) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    deleteConfigRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:delete-config")
    input := C.RpcInput(cstate, deleteConfigRpc)

    // target options: startup | url
    createInputLeaf(cstate, input, target, "target", "")

    readData := C.RpcExecute(cstate, deleteConfigRpc, provider)
    return readData == nil
}

func (ns *NetconfService) DiscardChanges(provider NetconfServiceProvider) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    discardChangesRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:discard-changes")
    readData := C.RpcExecute(cstate, discardChangesRpc, provider)
    return readData == nil
}

// TODO
func (ns *NetconfService) EditConfig(
    provider NetconfServiceProvider, target int, config types.Entity, defaultOper, testOp, errorOp string) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    editConfigRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:edit-config")
    // TODO
}

// incomplete
func (ns *NetconfService) GetConfig(provider NetconfServiceProvider, source int, filter types.Entity) types.Entity {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    getConfigRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:get-config")
    input := C.RpcInput(cstate, getConfigRpc)

    // source options: candidate | running | startup
    createInputLeaf(cstate, input, source, "source")

    payload := getXmlSubtreeFilterPayload(cstate, entity, provider)
    C.DataNodeCreate(cstate, input, "filter", payload)

    readData := C.RpcExecute(cstate, getConfigRpc, provider)
    if (readData == nil) {
        return nil
    }

    topEntity := getTopEntityFromFilter(filter)

    // TODO: how to handle DataNodeChildren struct go?
    node := C.DataNodeGetChildren(readData)[0]
    getEntityFromDataNode(node, topEntity)

    return topEntity
}

// incomplete
func (ns *NetconfService) Get(provider NetconfServiceProvider, filter types.Entity) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    getRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:get")
    input := C.RpcInput(cstate, getRpc)

    payload := getXmlSubtreeFilterPayload(cstate, entity, provider)
    C.DataNodeCreate(cstate, input, "filter", payload)

    readData := C.RpcExecute(cstate, getConfigRpc, provider)
    if (readData == nil) {
        return nil
    }

    topEntity := getTopEntityFromFilter(filter)

    // TODO: how to handle DataNodeChildren struct go?
    node := C.DataNodeGetChildren(readData)[0]
    getEntityFromDataNode(node, topEntity)

    return topEntity
}

func (ns *NetconfService) KillSession(provider NetconfServiceProvider, sessionId int) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    killSessionRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:kill-session")
    input := C.RpcInput(cstate, killSessionRpc)
    sessionIdStr := strconv.Itoa(sessionId)
    C.DataNodeCreate(cstate, input, "session-id", sessionIdStr)

    readData := C.RpcExecute(cstate, killSessionRpc, provider)
    return readData == nil
}

func (ns *NetconfService) Lock(provider NetconfServiceProvider, target int) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    lockRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:lock")
    input := C.RpcInput(cstate, lockRpc)

    // target options: candidate | running | startup
    createInputLeaf(cstate, input, target, "target", "")

    readData := C.RpcExecute(cstate, lockRpc, provider)
    return readData == nil
}

func (ns *NetconfService) Unlock(provider NetconfServiceProvider, target DataStore) bool {
    var cstate C.YDKStatePtr = C.YDKStateCreate()
    defer C.YDKStateFree(cstate)

    unlockRpc := getRpcFromProvider(cstate, provider, "ietf-netconf:unlock")
    input := C.RpcInput(cstate, unlockRpc)

    // target options: candidate | running | startup
    createInputLeaf(cstate, input, target, "target", "")

    readData := C.RpcExecute(cstate, unlockRpc, provider)
    return readData == nil
}

// TODO: figure out optional arguments
// func (ns *NetconfService) Validate(provider NetconfServiceProvider, source DataStore, url string) bool {
// }

// func (ns *NetconfService) Validate(provider NetconfServiceProvider, source DataStore) bool {
// }

func getRpcFromProvider(cstate C.YDKStatePtr, provider NetconfServiceProvider, path string) C.Rpc {
    rootSchemaNode := C.ServiceProviderGetRootSchema(cstate, provider)
    rpc := C.RootSchemaNodeRpc(cstate, rootSchemaNode, C.CString(path))
    return rpc
}

func createInputLeaf(cstate C.YDKStatePtr, inputDataNode C.DataNode, dataStore int, dataStoreStr string, url string) {
    var path bytes.Buffer
    path.WriteString(dataStoreStr)

    switch dataStore {
    case Candidate:
        path.WriteString("/candidate")
    case Running:
        path.WriteString("/running")
    case StartUp:
        path.WriteString("/startup")
    case Url:
        if (url == "") {
            err := errors.New("URL needs to be specified")
            log.Fatal(err)
        }
        path.WriteString("/url")
    }

    C.DataNodeCreate(cstate, inputDataNode, path.String(), "")
}

func getXmlSubtreeFilterPayload(cstate C.YDKStatePtr, entity types.Entity, provider NetconfServiceProvider) string {
    // XmlSubTreeCodec* not imported or implemented!!
    var xmlSubtreeCodec C.XmlSubtreeCodec = C.XmlSubtreeCodecCreate()
    defer C.XmlSubtreeCodec(xmlSubtreeCodec)

    payload := xmlSubtreeCodecEncode(cstate, xmlSubtreeCodec, entity, provider)
    return payload
}

func getTopEntityFromFilter(filter types.Entity) types.Entity {
    parent := filter.GetParent()
    if (parent == nil) {
        return filter.GetFilter()
    }
    return getTopEntityFromFilter(parent)
}

// TODO
func getEntityFromDataNode(node C.DataNode, entity types.Entity) {
    if (node == nil || entity == nil) {
        return
    }

    // not sure how this is supposed to work...
    children := C.DataNodeGetChildren(node)
}
