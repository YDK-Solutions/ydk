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

#include "entity_data_node_walker.hpp"
#include "errors.hpp"
#include "ietf_parser.hpp"
#include "gnmi_client.hpp"
#include "gnmi_provider.hpp"
#include "types.hpp"
#include "ydk_yang.hpp"
#include "logger.hpp"

using namespace std;
using namespace ydk;

namespace ydk
{

    gNMIServiceProvider::gNMIServiceProvider(const string& address)
        : session{address}
    {
        YLOG_INFO("Connected to {}", address);
    }

    gNMIServiceProvider::gNMIServiceProvider(const string& address, bool is_secure)
        : session{address, is_secure}
    {
        YLOG_INFO("Connected to {} via Secure Channel", address);
    }


    gNMIServiceProvider::gNMIServiceProvider(path::Repository & repo, const string& address)
        : session{repo, address}
    {
        YLOG_INFO("Connected to {}", address);
    }

    gNMIServiceProvider::gNMIServiceProvider(path::Repository & repo, const string& address, bool is_secure)
        : session{repo, address, is_secure}
    {
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
