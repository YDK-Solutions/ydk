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
#include <thread>
#include <chrono>

#include <ydk/gnmi_provider.hpp>
#include <ydk/gnmi_service.hpp>
#include <ydk/codec_provider.hpp>
#include <ydk/codec_service.hpp>
#include <ydk/entity_data_node_walker.hpp>

#include <ydk/filters.hpp>

#include "test_utils.hpp"

#include "../../core/src/catch.hpp"
#include "../../core/tests/config.hpp"

using namespace std;
using namespace ydk;
using namespace path;
using namespace ydktest;

static string cap_update = R"(
      {
        "name": "openconfig-bgp",
        "organization": "OpenConfig working group",
        "version": "2016-06-21"
      },)";

TEST_CASE("gnmi_service_capabilities")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
    gNMIService gs{};

    string json_caps = gs.capabilities(provider);
    //std::cout << json_caps << std::endl;
    REQUIRE(json_caps.find(cap_update) != string::npos);
}

static string int_update = R"([{\"name\":\"Loopback10\",\"config\":{\"name\":\"Loopback10\",\"description\":\"Test\"}}])";

static string bgp_update = R"([{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}])";

void read_sub(const char * subscribe_response)
{
    cout << endl << subscribe_response << endl;
}

void gnmi_service_subscribe_callback(const char * subscribe_response)
{
	//read_sub(subscribe_response);
	string response = subscribe_response;
	REQUIRE(response.find(int_update) != string::npos);
}

void gnmi_service_subscribe_multiples_callback(const char * subscribe_response)
{
	//read_sub(subscribe_response);
	string response = subscribe_response;
	REQUIRE(response.find(int_update) != string::npos);
    REQUIRE(response.find(bgp_update) != string::npos);
}

TEST_CASE("gnmi_service_subscribe")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
    gNMIService gs{};

    build_int_config(provider);

    // Build subscription
    auto i = make_shared<openconfig_interfaces::Interfaces::Interface>();
    i->yfilter = YFilter::read;

    gNMISubscription subscription{};
    subscription.entity = i.get();
    //subscription.subscription_mode = "ON_CHANGE";
    subscription.sample_interval = 10000000000;
    //subscription.suppress_redundant = true;
    subscription.heartbeat_interval = 100000000000;

    gs.subscribe(provider, subscription, 10, "ONCE", "JSON_IETF", gnmi_service_subscribe_callback);
}

TEST_CASE("gnmi_service_subscribe_multiples")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
    gNMIService gs{};

    // Build configuration on the server
    build_int_config(provider);
    build_bgp_config(provider);

    // Build subscription
    auto ifc = make_shared<openconfig_interfaces::Interfaces::Interface>();
    ifc->name = "*";

    auto neighbor = make_shared<openconfig_bgp::Bgp::Neighbors::Neighbor>();
    neighbor->neighbor_address = "*";

    // Build subscriptions
    gNMISubscription ifc_subscription{};
    ifc_subscription.entity = ifc.get();
    //ifc_subscription.subscription_mode = "ON_CHANGE";
    ifc_subscription.sample_interval = 10000000000;
    //ifc_subscription.suppress_redundant = true;
    ifc_subscription.heartbeat_interval = 100000000000;

    gNMISubscription bgp_subscription{};
    bgp_subscription.entity = neighbor.get();
    bgp_subscription.subscription_mode = "TARGET_DEFINED";
    bgp_subscription.sample_interval = 20000000000;
    bgp_subscription.suppress_redundant = false;
    bgp_subscription.heartbeat_interval = 200000000000;

    vector<gNMISubscription*> subscription_list{};
    subscription_list.push_back(&ifc_subscription);
    subscription_list.push_back(&bgp_subscription);

    // Subscribe
    gs.subscribe(provider, subscription_list, 10, "ONCE", "JSON_IETF", gnmi_service_subscribe_multiples_callback);
}

static int counter;

void set_counter(int max)
{
	counter = max;
}

bool counter_poll_request(const char * subscribe_response)
{
	std::this_thread::sleep_for(std::chrono::seconds(2));
	return --counter >= 0;
}

TEST_CASE("gnmi_service_poll_subscribe")
{
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
    gNMIService gs{};

    build_int_config(provider);

    // Build subscription
    openconfig_interfaces::Interfaces filter = {};
    auto i = make_shared<openconfig_interfaces::Interfaces::Interface>();
    i->name = "*";
    filter.interface.append(i);

    gNMISubscription subscription{};
    subscription.entity = &filter;

    set_counter(2);
    gs.subscribe(provider, subscription, 10, "POLL", "JSON_IETF", gnmi_service_subscribe_callback, counter_poll_request);
}

TEST_CASE("gnmi_service_stream_subscribe")
{
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
    gNMIService gs{};

    build_int_config(provider);

    // Build subscription
    openconfig_interfaces::Interfaces filter = {};
    auto i = make_shared<openconfig_interfaces::Interfaces::Interface>();
    i->name = "*";
    filter.interface.append(i);

    gNMISubscription subscription{};
    subscription.entity = &filter;
    subscription.sample_interval = 2000000000;
    subscription.heartbeat_interval = 10000000000;

    gs.subscribe(provider, subscription, 10, "STREAM", "JSON_IETF", gnmi_service_subscribe_callback);
}

TEST_CASE("gnmi_service_create")
{
    // session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
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
    //print_entity(get_reply, provider.get_session().get_root_schema());
    string expected = R"( <interfaces>
   <interface>
     <name>Loopback10</name>
     <config>
       <name>Loopback10</name>
       <description>Test</description>
     </config>
   </interface>
 </interfaces>
)";
    REQUIRE(entity2string(get_reply, provider.get_session().get_root_schema()) == expected);
}

TEST_CASE("gnmi_service_get_multiple")
{
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;
    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
    gNMIService gs{};

    // Set Create Request
    openconfig_bgp::Bgp bgp{};
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
}

TEST_CASE("gnmi_service_delete")
{
	// session
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;

    gNMIServiceProvider provider{repo, address, port, "admin", "admin"};
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

TEST_CASE("gnmi_service_get_list_element")
{
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;
    string username = "admin"; string password = "admin";

    gNMIServiceProvider provider{repo, address, port, username, password};
    gNMIService gs{};

    build_int_config(provider);

    auto ifc = make_shared<openconfig_interfaces::Interfaces::Interface>();
    ifc->name = "Loopback10";

    openconfig_interfaces::Interfaces ifcs{};
    ifcs.interface.append(ifc);

    auto get_reply = gs.get(provider, ifcs, "CONFIG");
    REQUIRE(get_reply != nullptr);
    //print_entity(get_reply, provider.get_session().get_root_schema());
    string expected = R"( <interfaces>
   <interface>
     <name>Loopback10</name>
     <config>
       <name>Loopback10</name>
       <description>Test</description>
     </config>
   </interface>
 </interfaces>
)";
    REQUIRE(entity2string(get_reply, provider.get_session().get_root_schema()) == expected);
}

TEST_CASE("gnmi_service_get_leaf")
{
    path::Repository repo{TEST_HOME};
    string address = "127.0.0.1"; int port = 50051;
    string username = "admin"; string password = "admin";

    gNMIServiceProvider provider{repo, address, port, username, password};
    gNMIService gs{};

    build_int_config(provider);

    auto ifc = make_shared<openconfig_interfaces::Interfaces::Interface>();
    ifc->name = "Loopback10";
    ifc->config->description.yfilter = YFilter::read;

    openconfig_interfaces::Interfaces ifcs{};
    ifcs.interface.append(ifc);

    auto get_reply = gs.get(provider, ifcs, "CONFIG");
    REQUIRE(get_reply != nullptr);
    string expected = R"( <interfaces>
   <interface>
     <name>Loopback10</name>
     <config>
       <description>Test</description>
     </config>
   </interface>
 </interfaces>
)";
    REQUIRE(entity2string(get_reply, provider.get_session().get_root_schema()) == expected);
}
