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

#include "path_api.hpp"
#include "gnmi_provider.hpp"
#include "config.hpp"
#include "catch.hpp"

const char* gnmi_expected_bgp_output ="{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172},\"afi-safis\":{\"afi-safi\":[{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"config\":{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"enabled\":true}}]}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172},\"afi-safis\":{\"afi-safi\":[{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"config\":{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"enabled\":true}}]}}]}}}";
const char* gnmi_expected_bgp_read ="{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172},\"afi-safis\":{\"afi-safi\":[{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"config\":{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"enabled\":true}}]}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172},\"afi-safis\":{\"afi-safi\":[{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"config\":{\"afi-safi-name\":\"openconfig-bgp-types:L3VPN_IPV4_UNICAST\",\"enabled\":true}}]}}]}}}";

void gnmi_print_tree(ydk::path::DataNode* dn, const std::string& indent)
{
    ydk::path::Statement s = dn->get_schema_node().get_statement();
    if(s.keyword == "leaf" || s.keyword == "leaf-list" || s.keyword == "anyxml") {
        auto val = dn->get_value();
        std::cout << indent << "{\"" << s.arg << "\":" << val << "}" << std::endl;
    } else {
        std::string child_indent{indent};
        child_indent+="  ";
        std::cout << indent << "{\"" << s.arg << "\":" << std::endl;
        for(auto c : dn->get_children())
	    gnmi_print_tree(c.get(), child_indent);
        std::cout << indent << "}" << std::endl;
    }
}

TEST_CASE("gnmi_test_json_payload"  )
{
    ydk::path::Repository repo{TEST_HOME};

    ydk::gNMIServiceProvider sp{repo,"127.0.0.1:50051"};
    ydk::path::RootSchemaNode& root = sp.get_root_schema();

    ydk::path::Codec s{};

    auto bgp = s.decode(root, R"(
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

    auto edit_config = root.create_rpc("ietf-netconf:get-config");
    edit_config->get_input_node().create_datanode("source/candidate");
    edit_config->get_input_node().create_datanode("filter",json);

    std::string rpc_json = s.encode(edit_config->get_input_node(), ydk::EncodingFormat::JSON, false);
    std::cout<< rpc_json<<std::endl;
}

TEST_CASE("bgp_gnmi_create")
{
    ydk::path::Repository repo{TEST_HOME};

    ydk::gNMIServiceProvider sp{repo,"127.0.0.1:50051"};
    ydk::path::RootSchemaNode& schema = sp.get_root_schema();

    ydk::path::Codec s{};

    auto & bgp = schema.create_datanode("openconfig-bgp:bgp", "");
    //first delete
    std::shared_ptr<ydk::path::Rpc> delete_rpc { schema.create_rpc("ydk:delete") };
    auto json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
    delete_rpc->get_input_node().create_datanode("entity", json);
    //call delete
    (*delete_rpc)(sp);


    auto & as = bgp.create_datanode("global/config/as", "65172");
    auto & l3vpn_ipv4_unicast = bgp.create_datanode("global/afi-safis/afi-safi[afi-safi-name='openconfig-bgp-types:L3VPN_IPV4_UNICAST']", "");
    auto & afi_safi_name = l3vpn_ipv4_unicast.create_datanode("config/afi-safi-name", "openconfig-bgp-types:L3VPN_IPV4_UNICAST");

    //set the enable flag
    auto & enable = l3vpn_ipv4_unicast.create_datanode("config/enabled","true");

    //bgp/neighbors/neighbor
    auto & neighbor = bgp.create_datanode("neighbors/neighbor[neighbor-address='172.16.255.2']", "");
    auto & neighbor_address = neighbor.create_datanode("config/neighbor-address", "172.16.255.2");
    auto & peer_as = neighbor.create_datanode("config/peer-as","65172");

    //bgp/neighbors/neighbor/afi-safis/afi-safi
    auto & neighbor_af = neighbor.create_datanode("afi-safis/afi-safi[afi-safi-name='openconfig-bgp-types:L3VPN_IPV4_UNICAST']", "");
    auto & neighbor_afi_safi_name = neighbor_af.create_datanode("config/afi-safi-name" , "openconfig-bgp-types:L3VPN_IPV4_UNICAST");
    auto & neighbor_enabled = neighbor_af.create_datanode("config/enabled","true");
    json = s.encode(bgp, ydk::EncodingFormat::JSON, false);

    CHECK( !json.empty());

    REQUIRE(json == gnmi_expected_bgp_output);

    //call create
    std::shared_ptr<ydk::path::Rpc> create_rpc{schema.create_rpc("ydk:create")};
    create_rpc->get_input_node().create_datanode("entity", json);
    (*create_rpc)(sp);

    //call read
    std::shared_ptr<ydk::path::Rpc> read_rpc{schema.create_rpc("ydk:read")};
    auto & bgp_read = schema.create_datanode("openconfig-bgp:bgp", "");
    json = s.encode(bgp_read, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json.empty() );
    read_rpc->get_input_node().create_datanode("filter", "{\"openconfig-bgp:bgp\":{}}");
    auto read_result = (*read_rpc)(sp); 
    REQUIRE(read_result != nullptr);
    gnmi_print_tree(read_result.get(),"");
    json = s.encode(*read_result, ydk::EncodingFormat::JSON, false);
    REQUIRE(json == gnmi_expected_bgp_read);

    //call update (update equiv. create in gnmi)
    peer_as.set_value("6500");
    json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
    CHECK( !json.empty());
    std::shared_ptr<ydk::path::Rpc> update_rpc { schema.create_rpc("ydk:create") };
    update_rpc->get_input_node().create_datanode("entity", json);
    (*update_rpc)(sp);
}

TEST_CASE("gnmi_core_validate")
{
    ydk::path::Repository repo{};

    ydk::gNMIServiceProvider sp{repo,"127.0.0.1:50051"};
    ydk::path::RootSchemaNode& schema = sp.get_root_schema();

    auto & runner = schema.create_datanode("ietf-netconf:validate", "");

    auto & ysanity = runner.create_datanode("source/candidate", "");

    ydk::path::Codec s{};
    auto json = s.encode(runner, ydk::EncodingFormat::JSON, false);

    CHECK( !json.empty());

    std::cout << json << std::endl;
}


TEST_CASE("gnmi_bgp_xr_openconfig"  )
{
    ydk::path::Repository repo{TEST_HOME};

    ydk::gNMIServiceProvider sp{repo,"127.0.0.1:50051"};
    ydk::path::RootSchemaNode& schema = sp.get_root_schema();

    ydk::path::Codec s{};

    auto & bgp = schema.create_datanode("openconfig-bgp:bgp", "");
    
    //call create
    auto & as = bgp.create_datanode("global/config/as", "65172");
    auto & router_id = bgp.create_datanode("global/config/router-id", "1.2.3.4");
    auto & l3vpn_ipv4_unicast = bgp.create_datanode("global/afi-safis/afi-safi[afi-safi-name='openconfig-bgp-types:L3VPN_IPV4_UNICAST']", "");
    auto & afi_safi_name = l3vpn_ipv4_unicast.create_datanode("config/afi-safi-name", "openconfig-bgp-types:L3VPN_IPV4_UNICAST");
    auto & enable = l3vpn_ipv4_unicast.create_datanode("config/enabled","true");
    //bgp/neighbors/neighbor
    auto & neighbor = bgp.create_datanode("neighbors/neighbor[neighbor-address='172.16.255.2']", "");
    auto & neighbor_address = neighbor.create_datanode("config/neighbor-address", "172.16.255.2");
    auto & peer_as = neighbor.create_datanode("config/peer-as","65172");
    auto & peer_group = neighbor.create_datanode("config/peer-group","IBGP");
    //bgp/peer-groups/peer-group
    auto & ppeer_group = bgp.create_datanode("peer-groups/peer-group[peer-group-name='IBGP']", "");
    auto & peer_group_name = ppeer_group.create_datanode("config/peer-group-name", "IBGP");
    auto & ppeer_as = ppeer_group.create_datanode("config/peer-as","65172");

    std::shared_ptr<ydk::path::Rpc> create_rpc { schema.create_rpc("ydk:create") };
    auto json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json.empty() );
    create_rpc->get_input_node().create_datanode("entity", json);

    auto res = (*create_rpc)(sp);

	//call read
    std::shared_ptr<ydk::path::Rpc> read_rpc { schema.create_rpc("ydk:read") };
    auto & bgp_read = schema.create_datanode("openconfig-bgp:bgp", "");

    json = s.encode(bgp_read, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json.empty() );
    read_rpc->get_input_node().create_datanode("filter", json);
    read_rpc->get_input_node().create_datanode("only-config");

    auto read_result = (*read_rpc)(sp);

    REQUIRE(read_result != nullptr);
}
