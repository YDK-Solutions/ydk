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

#include "config.hpp"
#include "catch.hpp"
#include <ydk/gnmi_service.hpp>
#include "ydk_ydktest/openconfig_bgp.hpp"
//#include <ydk_ydktest/ydktest_sanity.hpp>
//#include <ydk/types.hpp>

using namespace std;
using namespace ydk;
using namespace path;

using namespace std;
using namespace ydk;
//using namespace ydktest;

void config_bgp(openconfig::openconfig_bgp::Bgp bgp)
{
    bgp.global->config->as = 65172;
    
    auto neighbor = make_unique<openconfig::openconfig_bgp::Bgp::Neighbors::Neighbor>();
    neighbor->neighbor_address = "172.16.255.2";
    neighbor->config->neighbor_address = "172.16.255.2";
    neighbor->config->peer_as = 65172;
    
    neighbor->parent = bgp.neighbors.get();
    bgp.neighbors->neighbor.push_back(move(neighbor));
}

TEST_CASE("gnmi_service_get")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1:50051";

    gNMIServiceProvider provider{repo, address};
    gNMIService gs{address};

    openconfig::openconfig_bgp::Bgp filter = {};

    // Get Request 
    auto get_reply = gs.get(provider, filter);
    REQUIRE(get_reply);

}

TEST_CASE("gnmi_service_create")
{
	// session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1:50051";

    gNMIServiceProvider provider{repo, address};
    gNMIService gs{address};

    openconfig::openconfig_bgp::Bgp filter = {};

    // Set Create Request
    openconfig::openconfig_bgp::Bgp bgp = {};
    config_bgp(bgp);
    auto create_reply = gs.set(provider, bgp, "gnmi_create");
    REQUIRE(create_reply);

    // Get Request
    auto get_after_create_reply = gs.get(provider, filter);
    REQUIRE(get_after_create_reply);
}

TEST_CASE("gnmi_service_delete")
{
	// session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1:50051";

    gNMIServiceProvider provider{repo, address};
    gNMIService gs{address};

    openconfig::openconfig_bgp::Bgp filter = {};
    openconfig::openconfig_bgp::Bgp bgp = {};

    // Set Delete Request
    auto delete_reply = gs.set(provider, bgp, "gnmi_delete");
    REQUIRE(delete_reply);

    // Get Request
    auto get_after_delete_reply = gs.get(provider, filter);
    REQUIRE(get_after_delete_reply);
}