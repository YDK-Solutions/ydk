/// YANG Development Kit
// Copyright 2018 Cisco Systems. All rights reserved
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
#ifndef GNMI_TEST_UTIL_HPP
#define GNMI_TEST_UTIL_HPP

#include <ydk_ydktest/openconfig_bgp.hpp>
#include <ydk_ydktest/openconfig_interfaces.hpp>

#include <ydk/path_api.hpp>
#include <ydk/gnmi_provider.hpp>

void print_data_node(std::shared_ptr<ydk::path::DataNode> dn);
void print_entity(std::shared_ptr<ydk::Entity> entity, ydk::path::RootSchemaNode& root);

std::string entity2string(std::shared_ptr<ydk::Entity> entity, ydk::path::RootSchemaNode& root);
std::string print_tree(ydk::path::DataNode* dn, const std::string& indent);

void config_bgp(ydktest::openconfig_bgp::Bgp& bgp);
void build_bgp_config(ydk::gNMIServiceProvider& provider);
void build_int_config(ydk::gNMIServiceProvider& provider);

void delete_bgp_config(ydk::gNMIServiceProvider& provider);
void delete_int_config(ydk::gNMIServiceProvider& provider);

#endif /* GNMI_TEST_UTIL_HPP */
