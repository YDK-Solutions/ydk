//////////////////////////////////////////////////////////////////
// @file golang.cpp
//
// YANG Development Kit
// Copyright 2017 Cisco Systems. All rights reserved
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

#include <ydk/ydk.h>
#include "catch.hpp"
#include <iostream>

using namespace std;

const char* test_string="<runner xmlns=\"http://cisco.com/ns/yang/ydktest-sanity\"><ytypes><built-in-t><number8>2</number8></built-in-t></ytypes></runner>";

TEST_CASE( "codec_encode"  )
{
    CodecService c = CodecServiceInit();
    Repository repo = RepositoryInitWithPath("/usr/local/share/ydktest@0.1.0/");
    ServiceProvider provider = NetconfServiceProviderInitWithRepo(repo, "localhost", "admin", "admin", 12022);

    RootSchemaNode root_schema = ServiceProviderGetRootSchema(provider);

    DataNode runner = RootSchemaNodeCreate(root_schema, "ydktest-sanity:runner");

    DataNodeCreate(runner, "ytypes/built-in-t/number8", "2");

    string s { CodecServiceEncode(c, runner, XML, 0)};
    REQUIRE(s == test_string);

    NetconfServiceProviderFree(provider);
    RepositoryFree(repo);
    CodecServiceFree(c);
}

TEST_CASE( "codec_decode"  )
{
    CodecService c = CodecServiceInit();
    Repository repo = RepositoryInitWithPath("/usr/local/share/ydktest@0.1.0/");
    ServiceProvider provider = NetconfServiceProviderInitWithRepo(repo, "localhost", "admin", "admin", 12022);

    RootSchemaNode root_schema = ServiceProviderGetRootSchema(provider);

    DataNode runner = CodecServiceDecode(c, root_schema, test_string, XML);

    REQUIRE(runner!=NULL);

    string s { CodecServiceEncode(c, runner, XML, 0)};
    REQUIRE(s == test_string);

    NetconfServiceProviderFree(provider);
    RepositoryFree(repo);
    CodecServiceFree(c);
}

TEST_CASE( "provider_withpath"  )
{
    Repository repo = RepositoryInitWithPath("/usr/local/share/ydktest@0.1.0/");
    ServiceProvider provider = NetconfServiceProviderInitWithRepo(repo, "localhost", "admin", "admin", 12022);

    REQUIRE(repo!=NULL);
    REQUIRE(provider!=NULL);

    NetconfServiceProviderFree(provider);
    RepositoryFree(repo);
}

TEST_CASE( "provider"  )
{
    Repository repo = RepositoryInit();
    ServiceProvider provider = NetconfServiceProviderInitWithRepo(repo, "localhost", "admin", "admin", 12022);

    REQUIRE(repo!=NULL);
    REQUIRE(provider!=NULL);

    NetconfServiceProviderFree(provider);
    RepositoryFree(repo);
}

TEST_CASE( "rpc" )
{
    CodecService c = CodecServiceInit();

    Repository repo = RepositoryInitWithPath("/usr/local/share/ydktest@0.1.0/");

    ServiceProvider provider = NetconfServiceProviderInitWithRepo(repo, "localhost", "admin", "admin", 12022);
    RootSchemaNode root_schema = ServiceProviderGetRootSchema(provider);

    DataNode runner = RootSchemaNodeCreate(root_schema, "ydktest-sanity:runner");

    DataNodeCreate(runner, "ytypes/built-in-t/number8", "2");
    const char* create_xml = CodecServiceEncode(c, runner, XML, 0);

    Rpc create_rpc = RootSchemaNodeRpc(root_schema, "ydk:create");
    DataNode input = RpcInput(create_rpc);
    DataNodeCreate(input, "entity", create_xml);
    RpcExecute(create_rpc, provider);

    Rpc read_rpc = RootSchemaNodeRpc(root_schema, "ydk:read");
    input = RpcInput(read_rpc);
    DataNode runner_filter = RootSchemaNodeCreate(root_schema, "ydktest-sanity:runner");
    const char* read_xml = CodecServiceEncode(c, runner_filter, XML, 0);

    DataNodeCreate(input, "filter", read_xml);
    DataNode read_data = RpcExecute(read_rpc, provider);

    delete create_xml;
    delete read_xml;
    NetconfServiceProviderFree(provider);
    RepositoryFree(repo);
    CodecServiceFree(c);
}
