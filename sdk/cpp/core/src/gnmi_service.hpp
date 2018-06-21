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

#include "path_api.hpp"
#include "types.hpp"

namespace ydk
{
class gNMIServiceProvider;

class gNMIService
{
public:
    gNMIService();
    ~gNMIService();
    std::shared_ptr<Entity> get(gNMIServiceProvider& provider, Entity& filter, bool only_config) const;
    bool set(gNMIServiceProvider& provider, Entity& filter) const;

    void subscribe(gNMIServiceProvider& provider, Entity& filter, const std::string & list_mode, long long qos,
                    const std::string & mode, int sample_interval, std::function<void(const std::string &)> func) const;
};

}
#endif /* GNMI_SERVICE_HPP */

