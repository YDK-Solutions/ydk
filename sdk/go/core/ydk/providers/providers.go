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

import (
	"github.com/CiscoDevNet/ydk-go/ydk/path"
	"github.com/CiscoDevNet/ydk-go/ydk/types"
)

type NetconfServiceProvider struct {
	Repo     types.Repository
	Address  string
	Username string
	Password string
	Port     int

	Private types.CServiceProvider
}

type RestconfServiceProvider struct {
	Repo     types.Repository
	Address  string
	Username string
	Password string
	Port     int
	Encoding types.EncodingFormat

	Private types.CServiceProvider
}

func (provider *NetconfServiceProvider) GetPrivate() interface{} {
	return provider.Private
}

func (provider *NetconfServiceProvider) Connect() {
	provider.Private = path.ConnectToProvider(provider.Repo, provider.Address, provider.Username, provider.Password, provider.Port)
}

func (provider *NetconfServiceProvider) Disconnect() {
	if provider.Private.Private == nil {
		return
	}
	path.DisconnectFromProvider(provider.Private)
}
