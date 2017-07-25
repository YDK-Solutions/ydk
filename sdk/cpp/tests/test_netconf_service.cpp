/*  ----------------------------------------------------------------
 Copyright 2016 Cisco Systems

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 ------------------------------------------------------------------*/

#include <string.h>
#include <iostream>

#include <ydk/netconf_provider.hpp>
#include <ydk/netconf_service.hpp>
#include <ydk_ydktest/ydktest_sanity.hpp>
#include <ydk_ydktest/openconfig_bgp.hpp>
#include <ydk/types.hpp>
#include <spdlog/spdlog.h>

#include "config.hpp"
#include "catch.hpp"

using namespace ydk;
using namespace ydktest;
using namespace std;

// cancel_commit -- issues in netsim
//TEST_CASE("cancel_commit")
//{
//    // session
//    path::Repository repo{TEST_HOME};
//    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
//    NetconfService ns{};
//
//    auto reply = ns.cancel_commit(session);
//    REQUIRE(reply);
//}

// close_session
TEST_CASE("close_session")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    auto reply = ns.close_session(session);
    REQUIRE(reply);
}

// commit
TEST_CASE("commit")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    auto reply = ns.commit(session);
    REQUIRE(reply);
}

// copy_config
TEST_CASE("copy_config")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    DataStore target = DataStore::candidate;
    DataStore source = DataStore::running;

    auto reply = ns.copy_config(session, target, source);
    REQUIRE(reply);
}

// delete_config -- issues in netsim
TEST_CASE("delete_config")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    DataStore target = DataStore::url;

//    auto reply = ns.delete_config(session, target, "http://test");
    CHECK_THROWS_AS(ns.delete_config(session, target, "http://test"), YCPPError);
}

// discard_changes
TEST_CASE("discard_changes")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    auto reply = ns.discard_changes(session);
    REQUIRE(reply);
}

// edit_config, get_config
TEST_CASE("edit_config")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    DataStore target = DataStore::candidate;
    DataStore source = DataStore::candidate;
    openconfig_bgp::Bgp filter = {};
    openconfig_bgp::Bgp bgp = {};
    bgp.global->config->as = 6500;

    auto reply = ns.edit_config(session, target, bgp);
    REQUIRE(reply);

    auto data = ns.get_config(session, source, filter);
    REQUIRE(data);

    auto data_ptr = dynamic_cast<openconfig_bgp::Bgp*>(data.get());
    REQUIRE(data_ptr != nullptr);
    REQUIRE(data_ptr->global->config->as == bgp.global->config->as);

    reply = ns.discard_changes(session);
    REQUIRE(reply);
}

// get
TEST_CASE("get")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    openconfig_bgp::Bgp filter = {};

    auto reply = ns.get(session, filter);
    REQUIRE(reply);
}

// kill_session
TEST_CASE("kill_session")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    int session_id = 3;

//    auto reply = ns.kill_session(session, session_id);
    CHECK_THROWS_AS(ns.kill_session(session, session_id), YCPPError);
}

// lock, unlock
TEST_CASE("lock")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    DataStore target = DataStore::candidate;

    auto reply = ns.lock(session, target);
    REQUIRE(reply);

    reply = ns.unlock(session, target);
    REQUIRE(reply);
}

// validate
TEST_CASE("validate")
{
    // session
    path::Repository repo{TEST_HOME};
    NetconfSession session{repo, "127.0.0.1", "admin", "admin", 12022};
    NetconfService ns{};

    DataStore source = DataStore::candidate;

    auto reply = ns.validate(session, source);
    REQUIRE(reply);
}
