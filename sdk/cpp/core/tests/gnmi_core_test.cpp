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

#include "../src/path_api.hpp"
#include "../src/path/path_private.hpp"
#include "config.hpp"
#include "catch.hpp"
#include "gnmi_util.hpp"
#include "mock_data.hpp"

TEST_CASE( "test_gnmi_datanode_to_path" )
{
    ydk::path::Repository repo{TEST_HOME};
    mock::MockSession sp{TEST_HOME, test_openconfig};
    auto & schema = sp.get_root_schema();

    auto & bgp = schema.create_datanode("openconfig-bgp:bgp");
    auto & neighbor = bgp.create_datanode("neighbors/neighbor[neighbor-address='172.16.255.2']");
    auto & neighbor_address = neighbor.create_datanode("config/neighbor-address", "172.16.255.2");
    auto & peer_as = neighbor.create_datanode("config/peer-as","65172");

    gnmi::Path* path = new gnmi::Path();
    ydk::parse_datanode_to_path(&peer_as, path);

    std::string expected = R"(origin: "openconfig-bgp"
elem {
  name: "bgp"
}
elem {
  name: "neighbors"
}
elem {
  name: "neighbor"
  key {
    key: "neighbor-address"
    value: "172.16.255.2"
  }
}
elem {
  name: "config"
}
elem {
  name: "peer-as"
}
)";
    REQUIRE(path->DebugString() == expected);
}

TEST_CASE( "test_gnmi_datanode_to_path_2" )
{
    ydk::path::Repository repo{TEST_HOME};
    mock::MockSession sp{TEST_HOME, test_openconfig};
    auto & schema = sp.get_root_schema();

    auto & runner = schema.create_datanode("ydktest-sanity:runner");
    auto & two_pr  = runner.create_datanode("two-key-list[first='GigabitEthernet0/0/0/0'][second='222']");

    gnmi::Path* path = new gnmi::Path();
    ydk::parse_datanode_to_path(&two_pr, path);

    std::string expected = R"(origin: "ydktest-sanity"
elem {
  name: "runner"
}
elem {
  name: "two-key-list"
  key {
    key: "first"
    value: "GigabitEthernet0/0/0/0"
  }
  key {
    key: "second"
    value: "222"
  }
}
)";
    REQUIRE(path->DebugString() == expected);
}
