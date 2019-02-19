/// YANG Development Kit
// Copyright 2019 Cisco Systems. All rights reserved
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

#include "config.hpp"
#include "catch.hpp"
#include "test_utils.hpp"

#include <ydk/codec_provider.hpp>
#include <ydk/codec_service.hpp>

#include "ydk_ydktest_oc_nis/openconfig_network_instance.hpp"

using namespace std;
using namespace ydk;
using namespace ydktest_oc_nis;

const string nis_xml = R"(
<network-instances xmlns="http://openconfig.net/yang/network-instance">
  <network-instance>
    <name>default</name>
    <protocols>
      <protocol>
        <identifier xmlns:policy-types="http://openconfig.net/yang/policy-types">policy-types:ISIS</identifier>
        <name>DEFAULT</name>
        <config>
          <identifier xmlns:policy-types="http://openconfig.net/yang/policy-types">policy-types:ISIS</identifier>
          <name>DEFAULT</name>
        </config>
        <isis>
          <global>
            <afi-safi>
              <af>
                <afi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:IPV4</afi-name>
                <safi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:UNICAST</safi-name>
                <config>
                  <afi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:IPV4</afi-name>
                  <safi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:UNICAST</safi-name>
                </config>
              </af>
            </afi-safi>
          </global>
        </isis>
      </protocol>
    </protocols>
  </network-instance>
</network-instances>
)";


TEST_CASE( "decode_nis_part" )
{
    CodecServiceProvider codec_provider{EncodingFormat::XML};
    CodecService codec{};

    auto nis_top = make_shared<ydktest_oc_nis::openconfig_network_instance::NetworkInstances>();
    auto entity = codec.decode(codec_provider, nis_xml, nis_top);
    REQUIRE(entity != nullptr);
}
