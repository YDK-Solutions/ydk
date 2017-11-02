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

#include "catch.hpp"
#include "config.hpp"
#include "path_api.hpp"

const char* gnmi_expected_bgp_output ="{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}}}";
const char* gnmi_expected_bgp_read ="{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}}}";

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

    ydk::path::gNMISession session{repo,"127.0.0.1:50051"};
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


TEST_CASE("gnmi_bgp_create")
{
    ydk::path::Repository repo{"/Users/abhirame/.ydk/pavarotti:830"};

    ydk::path::gNMISession session{repo,"pavarotti:57400", true};
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

TEST_CASE("gnmi_core_validate")
{
    ydk::path::Repository repo{};

    ydk::path::gNMISession session{repo,"127.0.0.1:50051"};
    ydk::path::RootSchemaNode& schema = session.get_root_schema();

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

    ydk::path::gNMISession session{repo,"127.0.0.1:50051"};
    ydk::path::RootSchemaNode& schema = session.get_root_schema();

    ydk::path::Codec s{};

    auto & bgp = schema.create_datanode("openconfig-bgp:bgp", "");
    
    //call create
    auto & as = bgp.create_datanode("global/config/as", "65172");
    //bgp/neighbors/neighbor
    auto & neighbor = bgp.create_datanode("neighbors/neighbor[neighbor-address='172.16.255.2']", "");
    auto & neighbor_address = neighbor.create_datanode("config/neighbor-address", "172.16.255.2");
    auto & peer_as = neighbor.create_datanode("config/peer-as","65172");    

    std::shared_ptr<ydk::path::Rpc> create_rpc { schema.create_rpc("ydk:create") };
    auto json = s.encode(bgp, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json.empty() );
    create_rpc->get_input_node().create_datanode("entity", json);
    auto res = (*create_rpc)(session);

	//call read
    std::shared_ptr<ydk::path::Rpc> read_rpc { schema.create_rpc("ydk:read") };
    auto & bgp_read = schema.create_datanode("openconfig-bgp:bgp", "");
    json = s.encode(bgp_read, ydk::EncodingFormat::JSON, false);
    REQUIRE( !json.empty() );
    read_rpc->get_input_node().create_datanode("filter", json);
    read_rpc->get_input_node().create_datanode("only-config");
    auto read_result = (*read_rpc)(session);

    REQUIRE(read_result != nullptr);
}
