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
#include <ydk/codec_provider.hpp>
#include <ydk/codec_service.hpp>

using namespace std;
using namespace ydk;
using namespace path;

using namespace std;
using namespace ydk;
using namespace ydktest;

void print_entity(shared_ptr<ydk::Entity> entity, ydk::path::RootSchemaNode& root);

void config_bgp(openconfig_bgp::Bgp bgp)
{
    bgp.global->config->as = 65172;
    
    auto neighbor = make_shared<openconfig_bgp::Bgp::Neighbors::Neighbor>();
    neighbor->neighbor_address = "172.16.255.2";
    neighbor->config->neighbor_address = "172.16.255.2";
    neighbor->config->peer_as = 65172;
    
    neighbor->parent = bgp.neighbors.get();
    bgp.neighbors->neighbor.append(neighbor);
}

void read_sub(const string & s)
{
    cout<<s<<endl;
}

TEST_CASE("gnmi_service_capabilities")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port};
    gNMIService gs{};

    string json_caps = gs.capabilities(provider);
    read_sub(json_caps);
}

TEST_CASE("gnmi_service_subscribe")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port};
    gNMIService gs{};

    openconfig_interfaces::Interfaces filter = {};
    auto i = make_shared<openconfig_interfaces::Interfaces::Interface>();
    i->name = "*";
    filter.interface.append(i);

    // Build subscription
    gNMIService::Subscription subscription{};
    subscription.entity = &filter;
    subscription.subscription_mode = "ON_CHANGE";
    subscription.sample_interval = 10000000;
    subscription.suppress_redundant = true;
    subscription.heartbeat_interval = 1000000000;

    gs.subscribe(provider, &subscription, 10, "ONCE", read_sub);
}

TEST_CASE("gnmi_service_subscribe_multiples")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port};
    gNMIService gs{};

    openconfig_interfaces::Interfaces ifcs = {};
    auto ifc = make_shared<openconfig_interfaces::Interfaces::Interface>();
    ifc->name = "*";
    ifcs.interface.append(ifc);

    openconfig_bgp::Bgp bgp = {};
    auto neighbor = make_shared<openconfig_bgp::Bgp::Neighbors::Neighbor>();
    neighbor->neighbor_address = "*";
    bgp.neighbors->neighbor.append(neighbor);

    // Build subscriptions
    gNMIService::Subscription ifc_subscription{};
    ifc_subscription.entity = &ifcs;
    ifc_subscription.subscription_mode = "ON_CHANGE";
    ifc_subscription.sample_interval = 10000000;
    ifc_subscription.suppress_redundant = true;
    ifc_subscription.heartbeat_interval = 1000000000;

    gNMIService::Subscription bgp_subscription{};
    bgp_subscription.entity = &bgp;
    bgp_subscription.subscription_mode = "TARGET_DEFINED";
    bgp_subscription.sample_interval = 20000000;
    bgp_subscription.suppress_redundant = true;
    bgp_subscription.heartbeat_interval = 2000000000;

    vector<gNMIService::Subscription*> subscription_list{};
    subscription_list.push_back(&ifc_subscription);
    subscription_list.push_back(&bgp_subscription);

    // Get Request
    gs.subscribe(provider, subscription_list, 10, "ONCE", read_sub);
}

TEST_CASE("gnmi_service_create")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port};
    gNMIService gs{};
    CodecServiceProvider codec_provider{EncodingFormat::JSON};
    CodecService codec_service{};

    auto ifc = make_shared<openconfig_interfaces::Interfaces::Interface>();
    ifc->name = "Loopback10";
    ifc->config->name = "Loopback10";
    ifc->config->description = "Test";

    openconfig_interfaces::Interfaces ifcs{};
    ifcs.interface.append(ifc);
    ifcs.yfilter = YFilter::replace;

    // Set-replace Request
    auto set_reply = gs.set(provider, ifcs);

    openconfig_interfaces::Interfaces filter{};
    auto get_reply = gs.get(provider, filter, "CONFIG");
    REQUIRE(get_reply != nullptr);
    print_entity(get_reply, provider.get_session().get_root_schema());
}

TEST_CASE("gnmi_service_get_multiple")
{
	// session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;
    gNMIServiceProvider provider{repo, address, port};

//    string address = "10.30.110.86"; int port = 57400;
//    gNMIServiceProvider provider{repo, address, "admin", "admin", port};

    gNMIService gs{};

    // Set Create Request
    openconfig_bgp::Bgp bgp = {};
    bgp.yfilter = YFilter::replace;
    config_bgp(bgp);

    auto ifc = make_shared<openconfig_interfaces::Interfaces::Interface>();
    ifc->name = "Loopback10";
    ifc->config->name = "Loopback10";
    ifc->config->description = "Test";

    openconfig_interfaces::Interfaces ifcs{};
    ifcs.yfilter = YFilter::replace;
    ifcs.interface.append(ifc);

    vector<Entity*> set_entities;
    set_entities.push_back(&bgp);
    set_entities.push_back(&ifcs);

    auto set_reply = gs.set(provider, set_entities);

    // Get Request
    openconfig_bgp::Bgp bgp_filter{};
    openconfig_interfaces::Interfaces int_filter{};
    vector<Entity*> filter;
    filter.push_back(&bgp_filter);
    filter.push_back(&int_filter);

    auto get_reply = gs.get(provider, filter, "CONFIG");
    REQUIRE(get_reply.size() == 2);
//    for (auto entity : get_reply) {
//        print_entity(entity, provider.get_session().get_root_schema());
//    }
}

TEST_CASE("gnmi_service_delete")
{
	// session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port};
    gNMIService gs{};

    openconfig_bgp::Bgp filter = {};
    openconfig_bgp::Bgp bgp = {};

    // Set Delete Request
    bgp.yfilter = YFilter::delete_;
    auto delete_reply = gs.set(provider, bgp);
    REQUIRE(delete_reply);

    // Get Request
    auto get_after_delete_reply = gs.get(provider, filter, "STATE");
    REQUIRE(get_after_delete_reply);
}
