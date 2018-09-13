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
package providers

import (
//	"fmt"
//	"github.com/CiscoDevNet/ydk-go/ydk"
	"github.com/CiscoDevNet/ydk-go/ydk/errors"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
	"github.com/CiscoDevNet/ydk-go/ydk/path"
)

type GnmiServiceProvider struct {
	Repo        types.Repository
	Address     string
	Username    string
	Password    string
	Port        int
	ServerCert  string
	PrivateKey  string

	Private     types.CServiceProvider
	State       errors.State
}

// Functions to implement ServiceProvider interface
//
func (provider *GnmiServiceProvider) Connect() {
	if provider.Port == 0 {
		provider.Port = 57400
	}

	provider.Private = path.GnmiServiceProviderConnect(
		&provider.State,
		provider.Repo,
		provider.Address, provider.Port, provider.Username, provider.Password,
		provider.ServerCert, provider.PrivateKey)
}

func (provider *GnmiServiceProvider) Disconnect() {
	if provider.Private.Private == nil {
		return
	}
	path.GnmiServiceProviderDisconnect(provider.Private)
	path.CleanUpErrorState(&provider.State)
}

func (provider *GnmiServiceProvider) GetPrivate() interface{} {
	return provider.Private
}

func (provider *GnmiServiceProvider) GetState() *errors.State {
	return &provider.State
}

func (provider *GnmiServiceProvider) GetType() string {
    return "gNMIServiceProvider"
}

func (provider *GnmiServiceProvider) GetSession() *path.GnmiSession {
	session := path.GnmiServiceProviderGetSession(provider.Private)
	return session
}
