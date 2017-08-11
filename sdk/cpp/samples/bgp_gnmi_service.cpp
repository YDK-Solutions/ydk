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

#include <fstream>
#include <iostream>
#include <memory>
#include <spdlog/spdlog.h>

#include "ydk/path_api.hpp"
#include "ydk/gnmi_provider.hpp"
#include "ydk/gnmi_service.hpp"
#include "args_parser.h"
#include "ydk/crud_service.hpp"
#include "ydk_openconfig/openconfig_bgp.hpp"

using namespace std;
using namespace ydk;
using namespace path;

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


int main(int argc, char* argv[]) 
{
    vector<string> args = parse_args(argc, argv);
    if(args.empty()) return 1;
    
    string host, username, password, port, address;
    username = args[0]; password = args[1]; host = args[2]; port = args[3];

    address.append(host);
    address.append(":");
    address.append(port);
    bool verbose = (args[4]=="--verbose");
    if(verbose)
    {
        auto logger = spdlog::stdout_color_mt("ydk");
        logger->set_level(spdlog::level::debug);
    }

    ydk::path::Repository repo{"/usr/local/share/ydk/0.0.0.0\:50051/"};

    bool is_secure = true;	
    gNMIService gs{"127.0.0.1:50051"};
    openconfig::openconfig_bgp::Bgp filter = {};

    if(is_secure)
    {
    	gNMIServiceProvider provider{repo, address, is_secure};
    	// Get Request 
	    gs.get(provider, filter);

	    // Set Create Request
	    openconfig::openconfig_bgp::Bgp bgp = {};
	    config_bgp(bgp);
	    gs.set(provider, bgp, "gnmi_create");

	    // Get Request
	    gs.get(provider, filter);

	    // Set Delete Request
	    gs.set(provider, bgp, "gnmi_delete");

	    // Get Request
	    gs.get(provider, filter);
    }
    else
    {
    	gNMIServiceProvider provider{repo, address};
    	// Get Request 
	    gs.get(provider, filter);

	    // Set Create Request
	    openconfig::openconfig_bgp::Bgp bgp = {};
	    config_bgp(bgp);
	    gs.set(provider, bgp, "gnmi_create");

	    // Get Request
	    gs.get(provider, filter);

	    // Set Delete Request
	    gs.set(provider, bgp, "gnmi_delete");

	    // Get Request
	    gs.get(provider, filter);
    }

    return 0;
}