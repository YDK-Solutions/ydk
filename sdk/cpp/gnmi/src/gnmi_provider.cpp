/// YANG Development Kit
// Copyright 2016 Cisco Systems. All rights reserved
//
////////////////////////////////////////////////////////////////
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
//////////////////////////////////////////////////////////////////

#include <iostream>
#include <sstream>
#include <memory>
#include <libyang/libyang.h>

#include <ydk/errors.hpp>
#include <ydk/entity_data_node_walker.hpp>
#include <ydk/ietf_parser.hpp>
#include <ydk/types.hpp>
#include <ydk/logger.hpp>

#include "gnmi_client.hpp"
#include "gnmi_provider.hpp"

using namespace std;

namespace ydk
{
    gNMIServiceProvider::gNMIServiceProvider(path::Repository & repo,
                   const std::string& address,
                   const std::string& username,
                   const std::string& password,
                   int port)
        : session{repo, address, username, password, port}
    {
    	is_secure = true;
        YLOG_INFO("Connected to {} via Secure Channel", address);
    }

    gNMIServiceProvider::gNMIServiceProvider(path::Repository & repo,
                   const std::string& address,
                   int port)
        : session{repo, address, port}
    {
    	is_secure = false;
        YLOG_INFO("Connected to {} over insecure connection", address);
    }

    gNMIServiceProvider::gNMIServiceProvider(const std::string& address,
                   const std::string& username,
                   const std::string& password,
                   int port)
        : session{address, username, password, port}
    {
    	is_secure = true;
        YLOG_INFO("Connected to {} via Secure Channel", address);
    }

    gNMIServiceProvider::~gNMIServiceProvider()
    {
        YLOG_INFO("Disconnected from device");
    }

    EncodingFormat gNMIServiceProvider::get_encoding() const
    {
        return EncodingFormat::JSON;
    }

    const path::Session& gNMIServiceProvider::get_session() const
    {
        return session;
    }

    std::vector<std::string> gNMIServiceProvider::get_capabilities() const
    {
        return session.get_capabilities();
    }

}
