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
#include <spdlog/spdlog.h>

#include <ydk/gnmi_path_api.hpp>
#include <ydk/gnmi_util.hpp>

#include "../../core/src/path/path_private.hpp"
#include "../../core/src/catch.hpp"
#include "../../core/tests/config.hpp"
#include "../../core/tests/mock_data.hpp"

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
    ydk::path::parse_datanode_to_path(&peer_as, path);

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
    ydk::path::parse_datanode_to_path(&two_pr, path);

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
elem {
  name: "second"
}
)";
    REQUIRE(path->DebugString() == expected);
}

TEST_CASE("gnmi_test_json_payload"  )
{
    ydk::path::Repository repo{TEST_HOME};

    ydk::path::gNMISession session{repo, "127.0.0.1", 50051};
    ydk::path::RootSchemaNode& schema = session.get_root_schema();

    ydk::path::Codec s{};

    auto bgp = s.decode(schema, R"(
    {"openconfig-bgp:bgp": 
      {
         "@": {
            "ietf-netconf:operation":"merge"
        },
        "global": {
          "config": {
            "as": 65172,
            "router-id": "1.2.3.4"
          }
        },
        "neighbors": {
          "neighbor": [
            {
              "neighbor-address": "6.7.8.9",
              "config": {
                "local-as": 65001,
                "neighbor-address": "6.7.8.9",
                "peer-as": 65001,
                "peer-group": "IBGP"
              }
            }
          ]
        },
        "peer-groups": {
          "peer-group": [
            {
              "peer-group-name": "IBGP",
              "config": {
                "description": "test description",
                "local-as": 65001,
                "peer-as": 65001,
                "peer-group-name": "IBGP"
              }
            }
          ]
        }
      }
    })", ydk::EncodingFormat::JSON);

    std::string json = s.encode(*bgp, ydk::EncodingFormat::JSON, false);
    std::cout<< json<<std::endl;

    auto edit_config = schema.create_rpc("ietf-netconf:get-config");
    edit_config->get_input_node().create_datanode("source/candidate");
    edit_config->get_input_node().create_datanode("filter",json);

    std::string rpc_json = s.encode(edit_config->get_input_node(), ydk::EncodingFormat::JSON, false);
    std::cout<< rpc_json<<std::endl;
}

/* NOT READY FOR THIS TEST!!
 *
TEST_CASE("gnmi_bgp_create")
{
    ydk::path::Repository repo{"/Users/abhirame/.ydk/pavarotti:830"};

    ydk::path::gNMISession session{repo,"pavarotti:57400", "admin", "admin"};
    ydk::path::RootSchemaNode& schema = session.get_root_schema();

    ydk::path::Codec s{};

//    auto & bgp = schema.create_datanode("openconfig-bgp:bgp", "");
////    //first delete
//    std::shared_ptr<ydk::path::Rpc> delete_rpc { schema.create_rpc("ydk:delete") };
//    auto json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
//    delete_rpc->get_input_node().create_datanode("entity", json);
//    //call delete
//    (*delete_rpc)(session);
//
//    auto & as = bgp.create_datanode("global/config/as", "65172");
//
//    //bgp/neighbors/neighbor
//    auto & neighbor = bgp.create_datanode("neighbors/neighbor[neighbor-address='172.16.255.2']", "");
//    auto & neighbor_address = neighbor.create_datanode("config/neighbor-address", "172.16.255.2");
//    auto & peer_as = neighbor.create_datanode("config/peer-as","65172");
//
//    auto json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
//
//    CHECK( !json.empty());
//
//    REQUIRE(json == gnmi_expected_bgp_output);
//
//    //call create
//    std::shared_ptr<ydk::path::Rpc> create_rpc{schema.create_rpc("ydk:create")};
//    create_rpc->get_input_node().create_datanode("entity", json);
//    (*create_rpc)(session);

    //call read
    std::shared_ptr<ydk::path::Rpc> read_rpc{schema.create_rpc("ydk:read")};
    auto & bgp_read = schema.create_datanode("openconfig-bgp:bgp", "");
    auto json = s.encode(bgp_read, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json.empty() );
    read_rpc->get_input_node().create_datanode("filter", "{\"openconfig-bgp:bgp\":{}}");
    auto read_result = (*read_rpc)(session);
    REQUIRE(read_result != nullptr);
    gnmi_print_tree(read_result.get(),"");
    json = s.encode(*read_result, ydk::EncodingFormat::JSON, false);
//    REQUIRE(json == gnmi_expected_bgp_read);
//
//    //call update (update equiv. create in gnmi)
//    peer_as.set_value("6500");
//    json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
//    CHECK( !json.empty());
//    std::shared_ptr<ydk::path::Rpc> update_rpc { schema.create_rpc("ydk:create") };
//    update_rpc->get_input_node().create_datanode("entity", json);
//    (*update_rpc)(session);
}
*/

void read_sub(const std::string & s);

TEST_CASE("gnmi_rpc_subscribe")
{
    ydk::path::Repository repo{TEST_HOME};
    ydk::path::gNMISession session{repo, "127.0.0.1", 50051};
    ydk::path::RootSchemaNode& schema = session.get_root_schema();

    ydk::path::Codec codec{};

    std::shared_ptr<ydk::path::Rpc> rpc { schema.create_rpc("ydk:gnmi-subscribe") };
    auto & subscription = rpc->get_input_node().create_datanode("subscription", "");
    subscription.create_datanode("mode", "ONCE");
    subscription.create_datanode("qos", "10");

    auto & ifs = schema.create_datanode("openconfig-interfaces:interfaces", "");
    ifs.create_datanode("interface[name='*']/state");
    auto int_json = codec.encode(ifs, ydk::EncodingFormat::JSON, false);

    auto & bgp = schema.create_datanode("openconfig-bgp:bgp", "");
    bgp.create_datanode("neighbors/neighbor[neighbor-address='0.0.0.0']/state");
    auto bgp_json = codec.encode(bgp, ydk::EncodingFormat::JSON, false);

    auto & int_subscription = subscription.create_datanode("subscription-list[alias='int']", "");
    int_subscription.create_datanode("entity", int_json);
    int_subscription.create_datanode("subscription-mode", "ON_CHANGE");
    int_subscription.create_datanode("sample-interval", "10000000");
    int_subscription.create_datanode("suppress-redundant", "true");
    int_subscription.create_datanode("heartbeat-interval", "1000000000");

    auto & bgp_subscription = subscription.create_datanode("subscription-list[alias='bgp']", "");
    bgp_subscription.create_datanode("entity", bgp_json);
    // bgp_subscription.create_datanode("subscription-mode", "ON_CHANGE");	// can be skipped as default
    bgp_subscription.create_datanode("sample-interval", "20000000");
    // bgp_subscription.create_datanode("suppress-redundant", "true");  	// can be skipped as default
    // bgp_subscription.create_datanode("heartbeat-interval", "2000000000");// can be skipped as it is calculated from sample-interval

    // auto json = codec.encode(rpc->get_input_node(), ydk::EncodingFormat::JSON, true);
    // std::cout << "Built RPC:" << std::endl << json << std::endl;

    session.invoke(*rpc, read_sub);
}

TEST_CASE("gnmi_rpc_caps")
{
    ydk::path::Repository repo{TEST_HOME};
    ydk::path::gNMISession session{repo, "127.0.0.1", 50051};
    ydk::path::RootSchemaNode& schema = session.get_root_schema();

    std::shared_ptr<ydk::path::Rpc> cap_rpc { schema.create_rpc("ydk:gnmi-caps") };
    auto caps = (*cap_rpc)(session);

    ydk::path::Codec codec{};
    auto json = codec.encode(*caps, ydk::EncodingFormat::JSON, true);
    //std::cout << "Server capabilities:" << std::endl << json << std::endl;
}

TEST_CASE("gnmi_rpc_set_get_bgp")
{
    ydk::path::Repository repo{TEST_HOME};
    ydk::path::gNMISession session{repo, "127.0.0.1", 50051};
    ydk::path::RootSchemaNode& schema = session.get_root_schema();

    // Configure BGP
    auto & bgp = schema.create_datanode("openconfig-bgp:bgp");
    auto & as = bgp.create_datanode("global/config/as", "65172");
    auto & neighbor = bgp.create_datanode("neighbors/neighbor[neighbor-address='172.16.255.2']");
    auto & neighbor_address = neighbor.create_datanode("config/neighbor-address", "172.16.255.2");
    auto & peer_as = neighbor.create_datanode("config/peer-as","65172");    

    // Configure interface
    auto & ifc = schema.create_datanode("openconfig-interfaces:interfaces");
    auto & lo10_config = ifc.create_datanode("interface[name='Loopback10']");

    // Add data-nodes to RPC
    ydk::path::Codec s{};
    auto bgp_json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
    auto int_json = s.encode(ifc, ydk::EncodingFormat::JSON, false);

    std::shared_ptr<ydk::path::Rpc> set_rpc { schema.create_rpc("ydk:gnmi-set") };
    set_rpc->get_input_node().create_datanode("replace[alias='bgp']/entity", bgp_json);
    set_rpc->get_input_node().create_datanode("replace[alias='int']/entity", int_json);

    auto res = (*set_rpc)(session);

	//call read
    auto & bgp_read = schema.create_datanode("openconfig-bgp:bgp", "");
    auto json_bgp = s.encode(bgp_read, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json_bgp.empty() );

    auto & int_read = schema.create_datanode("openconfig-interfaces:interfaces", "");
    auto json_int = s.encode(int_read, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json_int.empty() );

    std::shared_ptr<ydk::path::Rpc> read_rpc { schema.create_rpc("ydk:gnmi-get") };
    read_rpc->get_input_node().create_datanode("type", "CONFIG");

    read_rpc->get_input_node().create_datanode("request[alias='bgp']/entity", json_bgp);
    read_rpc->get_input_node().create_datanode("request[alias='int']/entity", json_int);

    auto read_result = (*read_rpc)(session);

    REQUIRE(read_result != nullptr);
}
