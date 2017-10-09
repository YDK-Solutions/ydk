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

#ifndef GNMI_SERVICE_HPP
#define GNMI_SERVICE_HPP

#include <iostream>
#include <map>
#include <memory>
#include <sstream>
#include <string>

#include "codec_service.hpp"
#include "entity_data_node_walker.hpp"
#include "errors.hpp"
#include "gnmi_provider.hpp"
#include "logger.hpp"
#include "path_api.hpp"
#include "types.hpp"

namespace ydk
{
	class gNMIService
	{
	    public:
	        gNMIService();
	        ~gNMIService();
    		std::shared_ptr<path::DataNode> get(gNMIServiceProvider& provider, Entity& filter);
    		bool set(gNMIServiceProvider& provider, Entity& filter, std::string operation);
	};
}
#endif /* GNMI_SERVICE_HPP */
