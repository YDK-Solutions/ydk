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

#include "gnmi_path_api.hpp"

namespace ydk
{
class gNMIServiceProvider;

class gNMIService
{
  public:
    struct Subscription {
        Entity* entity;
        std::string subscription_mode;
        uint64 sample_interval;
        bool suppress_redundant;
        uint64 heartbeat_interval;
    };

    gNMIService();
    ~gNMIService();

    std::shared_ptr<Entity> get(gNMIServiceProvider & provider,
                                Entity& filter, const std::string & operation) const;
    std::vector<std::shared_ptr<Entity>> get(gNMIServiceProvider & provider,
                                std::vector<Entity*> & filter_list, const std::string & operation) const;

    bool set(gNMIServiceProvider & provider, Entity& entity) const;
    bool set(gNMIServiceProvider & provider, std::vector<Entity*> & entity_list) const;

    void subscribe(gNMIServiceProvider& provider,
                   Subscription* sub,
                   uint32 qos, const std::string & mode,
                   std::function<void(const std::string & response)> out_func,
                   std::function<bool(const std::string & response)> poll_func) const;
    void subscribe(gNMIServiceProvider& provider,
                   std::vector<Subscription*> & sub_list,
                   uint32 qos, const std::string & mode,
                   std::function<void(const std::string & response)> out_func,
                   std::function<bool(const std::string & response)> poll_func) const;

    std::string capabilities(gNMIServiceProvider & provider);
};

}
#endif /* GNMI_SERVICE_HPP */

