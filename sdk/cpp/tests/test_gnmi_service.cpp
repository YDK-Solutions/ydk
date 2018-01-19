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
#include <ydk/gnmi_provider.hpp>
#include <ydk/crud_service.hpp>
#include <ydk/gnmi_service.hpp>
#include <ydk_ydktest/openconfig_bgp.hpp>
#include <ydk_ydktest/openconfig_interfaces.hpp>
#include <ydk/filters.hpp>

using namespace std;
using namespace ydk;
using namespace path;

using namespace std;
using namespace ydk;
using namespace ydktest;

void config_bgp(openconfig_bgp::Bgp bgp)
{
    bgp.global->config->as = 65172;
    
    auto neighbor = make_unique<openconfig_bgp::Bgp::Neighbors::Neighbor>();
    neighbor->neighbor_address = "172.16.255.2";
    neighbor->config->neighbor_address = "172.16.255.2";
    neighbor->config->peer_as = 65172;
    
    neighbor->parent = bgp.neighbors.get();
    bgp.neighbors->neighbor.push_back(move(neighbor));
}

void read_sub(const string & s)
{
    cout<<s<<endl;
}

TEST_CASE("gnmi_service_subscribe")
{
    // session
    path::Repository repo{"/Users/abhirame/.ydk/pavarotti:830"};
    string address = "127.0.0.1:50051";

    gNMIServiceProvider provider{repo, "pavarotti", "admin", "admin"};
    gNMIService gs{};

    openconfig_interfaces::Interfaces filter = {};
    auto i = make_shared<openconfig_interfaces::Interfaces::Interface>();
    i->yfilter = YFilter::read;
//    i->name = "Loopback0";
//    i->config->name = "Loopback0";
//    i->config->description = "test";
    filter.interface.push_back(i);

    // Get Request 
    gs.subscribe(provider, filter, "ONCE", 10, "ON_CHANGE", 100000, read_sub);

}

TEST_CASE("gnmi_service_create")
{
    // session
    path::Repository repo{"/Users/abhirame/.ydk/pavarotti:830"};
    string address = "127.0.0.1:50051";

    gNMIServiceProvider provider{repo, "pavarotti", "admin", "admin"};
    gNMIService gs{};

    openconfig_interfaces::Interfaces filter = {};
    auto i = make_shared<openconfig_interfaces::Interfaces::Interface>();
    i->name = "Loopback0";
    i->config->name = "Loopback0";
    filter.interface.push_back(i);

    // Get Request 
    auto get_reply = gs.set(provider, filter, "gnmi_create");
    REQUIRE(get_reply);

}

TEST_CASE("gnmi_service_get")
{
	// session
    path::Repository repo{"/Users/abhirame/.ydk/pavarotti:830"};
    string address = "pavarotti:57400";

    gNMIServiceProvider provider{repo, "pavarotti", "admin", "admin"};
    gNMIService gs{};

//    openconfig_bgp::Bgp filter = {};
    openconfig_interfaces::Interfaces filter = {};
    // Set Create Request
    openconfig_bgp::Bgp bgp = {};
//    config_bgp(bgp);
    bgp.global->config->yfilter = YFilter::read;
//    auto create_reply = gs.set(provider, bgp, "gnmi_create");
//    REQUIRE(create_reply);

    // Get Request

//    auto get_after_create_reply = gs.get(provider, bgp, true);

//    REQUIRE(get_after_create_reply);
    CrudService crud{};
    auto r = crud.read_config(provider, bgp);
}

TEST_CASE("gnmi_service_delete")
{
	// session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1:50051";

    gNMIServiceProvider provider{repo, "127.0.0.1", "admin", "admin", 50051};
    gNMIService gs{};

    openconfig_bgp::Bgp filter = {};
    openconfig_bgp::Bgp bgp = {};

    // Set Delete Request
    auto delete_reply = gs.set(provider, bgp, "gnmi_delete");
    REQUIRE(delete_reply);

    // Get Request
    auto get_after_delete_reply = gs.get(provider, filter, false);
    REQUIRE(get_after_delete_reply);
}
