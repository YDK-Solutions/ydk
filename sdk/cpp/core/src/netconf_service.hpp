//
// @file crud_service.hpp
// @brief The main ydk public header.
//
// YANG Development Kit
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

#ifndef NETCONF_SERVICE_HPP
#define NETCONF_SERVICE_HPP

#include <map>
#include <memory>
#include <string>
#include "path/netconf_session.hpp"
#include "service.hpp"
#include "types.hpp"

namespace ydk
{

namespace core
{
class DataNode;
class Session;
}

class Entity;

enum class DataStore {
    candidate,
    running,
    startup,
    url
};

class NetconfService : public Service
{
    public:
        NetconfService();
        ~NetconfService();

        bool cancel_commit(NetconfSession & session, int persist_id = -1);

        bool close_session(NetconfSession & session);

        bool commit(NetconfSession & session, bool confirmed = false,
            int confirm_timeout = -1, int persist = -1, int persist_id = -1);

        bool copy_config(NetconfSession & session, DataStore target, DataStore source, std::string url = "");

        bool copy_config(NetconfSession & session, DataStore target, Entity& source_config);

        bool delete_config(NetconfSession & session, DataStore target, std::string url = "");

        bool discard_changes(NetconfSession & session);

        bool edit_config(NetconfSession & session, DataStore target, Entity& config,
            std::string default_operation = "", std::string test_option = "", std::string error_option = "");

        std::shared_ptr<Entity> get_config(NetconfSession & session, DataStore source, Entity& filter);

        std::shared_ptr<Entity> get(NetconfSession & session, Entity& filter);

        bool kill_session(NetconfSession & session, int session_id);

        bool lock(NetconfSession & session, DataStore target);

        bool unlock(NetconfSession & session, DataStore target);

        bool validate(NetconfSession & session, DataStore source, std::string url = "");

        bool validate(NetconfSession & session, Entity& source_config);
};

}

#endif /* NETCONF_SERVICE_HPP */
