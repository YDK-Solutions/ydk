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

#include <ydk/entity_util.hpp>

#include <ydk_ydktest/ydktest_sanity.hpp>

#include "config.hpp"
#include "catch.hpp"
#include "test_utils.hpp"

using namespace std;
using namespace ydk;
using namespace ydktest;

string check_empty_str_value(string & v)
{
    if (v.empty())
        v = "exists";
    return v;
}

void print_dictionary(const string & legend, map<string,string> & ent_dict)
{
    cout << "\n------> DICTIONARY" << legend << endl;
    for (map<string,string>::iterator it = ent_dict.begin(); it != ent_dict.end(); ++it)
    {
        cout << it->first << ": " << check_empty_str_value(it->second) << endl;
    }
}

void print_diffs(map<string, pair<string,string>> & diff, Entity& ent_left, Entity& ent_right)
{
	cout << "\n------> DIFFS:" << endl;
    for (auto const & entry : diff)
    {
        auto value_pair = entry.second;
        cout << entry.first << ": " << check_empty_str_value(value_pair.first) << " vs "
                                    << check_empty_str_value(value_pair.second) << endl;
        auto & ent = (value_pair.first == "None") ? ent_right : ent_left;
        string path = entry.first;
        REQUIRE(path_to_entity(ent, path) != nullptr);
    }
}

TEST_CASE( "test_entity_diff_two_key" )
{
    auto runner1 = ydktest_sanity::Runner{};
    auto l_1 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_1->first = "f1";
    l_1->second = 11;
    l_1->property = "82";
    auto l_2 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_2->first = "f2";
    l_2->second = 22;
    l_2->property = "83";
    runner1.two_key_list.extend({l_1, l_2});

    auto ent_dict1 = entity_to_dict(runner1);
    REQUIRE(ent_dict1.size() == 4);
    print_dictionary("-LEFT", ent_dict1);

    auto runner2 = ydktest_sanity::Runner{};
    l_1 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_1->first = "f1";
    l_1->second = 11;
    l_1->property = "82";
    l_2 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_2->first = "f2";
    l_2->second = 22;
    l_2->property = "83";
    runner2.two_key_list.extend({l_1, l_2});

    auto diff = entity_diff(runner1, runner2);
    REQUIRE(diff.size() == 0);

    auto runner_clone = runner1.clone();
    diff = entity_diff(runner1, *runner_clone);
    REQUIRE(diff.size() == 0);

    l_1->property = "83";
    auto ent_dict2 = entity_to_dict(runner2);
    print_dictionary("-RIGHT", ent_dict2);
    diff = entity_diff(runner1, runner2);
    REQUIRE(diff.size() == 1);
    print_diffs(diff, runner1, runner2);
}

TEST_CASE( "test_entity_diff_two_key_not_equal" )
{
    auto runner1 = ydktest_sanity::Runner{};
    auto l_1 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_1->first = "f1";
    l_1->second = 11;
    l_1->property = "82";
    auto l_2 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_2->first = "f2";
    l_2->second = 22;
    l_2->property = "83";
    runner1.two_key_list.extend({l_1, l_2});

    auto ent_dict1 = entity_to_dict(runner1);
    REQUIRE(ent_dict1.size() == 4);
    print_dictionary("-LEFT", ent_dict1);

    auto runner2 = ydktest_sanity::Runner{};
    l_1 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_1->first = "f1";
    l_1->second = 11;
    l_1->property = "82";
    l_2 = make_shared<ydktest_sanity::Runner::TwoKeyList>();
    l_2->first = "f3";
    l_2->second = 22;
    l_2->property = "83";
    runner2.two_key_list.extend({l_1, l_2});

    auto ent_dict2 = entity_to_dict(runner2);
    print_dictionary("-RIGHT", ent_dict2);

    auto diff = entity_diff(runner1, runner2);
    REQUIRE(diff.size() == 2);
    print_diffs(diff, runner1, runner2);
}

TEST_CASE( "test_entity_to_dict_aug_onelist" )
{
    auto runner = ydktest_sanity::Runner{};
    auto e_1 = make_shared<ydktest_sanity::Runner::OneList::OneAugList::Ldata>();
    auto e_2 = make_shared<ydktest_sanity::Runner::OneList::OneAugList::Ldata>();
    e_1->number = 1;
    e_1->name = "e_1.name";
    e_2->number = 2;
    e_2->name = "e_2.name";
    runner.one_list->one_aug_list->ldata.extend({e_1, e_2});
    runner.one_list->one_aug_list->enabled = true;

    auto ent_dict = entity_to_dict(runner);
    REQUIRE(ent_dict.size() == 5);
    print_dictionary("", ent_dict);

    auto runner_clone = runner.clone();
    auto diff = entity_diff(runner, *runner_clone);
    REQUIRE(diff.size() == 0);
}

TEST_CASE( "test_entity_to_dict_enum_leaflist" )
{
    auto runner = ydktest_sanity::Runner{};
    runner.ytypes->built_in_t->enum_llist.append(ydktest_sanity::YdkEnumTest::local);
    runner.ytypes->built_in_t->enum_llist.append(ydktest_sanity::YdkEnumTest::remote);
    runner.ytypes->built_in_t->identity_llist.append(ydktest_sanity::ChildIdentity());
    runner.ytypes->built_in_t->identity_llist.append(ydktest_sanity::ChildChildIdentity());

    auto ent_dict = entity_to_dict(runner);
    REQUIRE(ent_dict.size() == 4);
    print_dictionary("", ent_dict);

    auto runner_clone = runner.clone();
    auto diff = entity_diff(runner, *runner_clone);
    REQUIRE(diff.size() == 0);
}

TEST_CASE( "test_entity_diff_presence" )
{
    auto runner = ydktest_sanity::Runner();
    runner.runner_2 = make_shared<ydktest_sanity::Runner::Runner2>();
    runner.runner_2->some_leaf = "some-leaf";
    runner.runner_2->parent = &runner;

    auto ent_dict = entity_to_dict(runner);
    REQUIRE(ent_dict.size() == 2);
    print_dictionary("-LEFT", ent_dict);

    auto empty_runner = ydktest_sanity::Runner();
    ent_dict = entity_to_dict(empty_runner);
    REQUIRE(ent_dict.size() == 0);
    print_dictionary("-RIGHT", ent_dict);

    auto diff = entity_diff(runner, empty_runner);
    REQUIRE(diff.size() == 1);
    print_diffs(diff, runner, empty_runner);

    auto runner_clone = runner.clone();
    diff = entity_diff(runner, *runner_clone);
    REQUIRE(diff.size() == 0);
}

TEST_CASE( "test_entity_diff_two_list_pos" )
{
    auto r_1 = ydktest_sanity::Runner();
    auto e_1 = make_shared<ydktest_sanity::Runner::TwoList::Ldata>();
    auto e_2 = make_shared<ydktest_sanity::Runner::TwoList::Ldata>();
    auto e_11 = make_shared<ydktest_sanity::Runner::TwoList::Ldata::Subl1>();
    auto e_12 = make_shared<ydktest_sanity::Runner::TwoList::Ldata::Subl1>();
    e_1->number = 21;
    e_1->name = "runner:twolist:ldata[21]:name";
    e_11->number = 211;
    e_11->name = "runner:twolist:ldata[21]:subl1[211]:name";
    e_12->number = 212;
    e_12->name = "runner:twolist:ldata[21]:subl1[212]:name";
    e_1->subl1.extend({e_11, e_12});
    auto e_21 = make_shared<ydktest_sanity::Runner::TwoList::Ldata::Subl1>();
    auto e_22 = make_shared<ydktest_sanity::Runner::TwoList::Ldata::Subl1>();
    e_2->number = 22;
    e_2->name = "runner:twolist:ldata[22]:name";
    e_21->number = 221;
    e_21->name = "runner:twolist:ldata[22]:subl1[221]:name";
    e_22->number = 222;
    e_22->name = "runner:twolist:ldata[22]:subl1[222]:name";
    e_2->subl1.extend({e_21, e_22});
    r_1.two_list->ldata.extend({e_1, e_2});

    auto ent_dict = entity_to_dict(r_1);
    REQUIRE(ent_dict.size() == 12);
    print_dictionary("-LEFT", ent_dict);

    auto r_2 = ydktest_sanity::Runner();
    e_1 = make_shared<ydktest_sanity::Runner::TwoList::Ldata>();
    e_2 = make_shared<ydktest_sanity::Runner::TwoList::Ldata>();
    e_11 = make_shared<ydktest_sanity::Runner::TwoList::Ldata::Subl1>();
    e_12 = make_shared<ydktest_sanity::Runner::TwoList::Ldata::Subl1>();
    e_1->number = 21;
    e_1->name = "runner:twolist:ldata[21]:name";
    e_11->number = 211;
    e_11->name = "runner:twolist:ldata[21]:subl1[211]:name";
    e_12->number = 212;
    e_12->name = "runner:twolist:ldata[21]:subl1[212]:name";
    e_1->subl1.extend({e_11, e_12});
    r_2.two_list->ldata.append(e_1);

    ent_dict = entity_to_dict(r_2);
    REQUIRE(ent_dict.size() == 6);
    print_dictionary("-RIGHT", ent_dict);

    auto diff = entity_diff(r_1, r_2);
    REQUIRE(diff.size() == 1);
    print_diffs(diff, r_1, r_2);

    auto runner_clone = r_1.clone();
    diff = entity_diff(r_1, *runner_clone);
    REQUIRE(diff.size() == 0);
}

TEST_CASE("test_entity_diff_no_keys") {
    auto r_1 = ydktest_sanity::Runner();
    auto t1 = make_shared<ydktest_sanity::Runner::NoKeyList> (); t1->test = "t1";
    auto t2 = make_shared<ydktest_sanity::Runner::NoKeyList> (); t2->test = "t2";
    r_1.no_key_list.extend({t1, t2});

    auto ent_dict = entity_to_dict(r_1);
    REQUIRE(ent_dict.size() == 4);
    print_dictionary("-LEFT", ent_dict);

    auto r_2 = ydktest_sanity::Runner();
    t1 = make_shared<ydktest_sanity::Runner::NoKeyList> (); t1->test = "t1";
    t2 = make_shared<ydktest_sanity::Runner::NoKeyList> (); t2->test = "tt";
    auto t3 = make_shared<ydktest_sanity::Runner::NoKeyList> (); t3->test = "t3";
    r_2.no_key_list.extend({t1, t2, t3});

    ent_dict = entity_to_dict(r_2);
    REQUIRE(ent_dict.size() == 6);
    print_dictionary("-RIGHT", ent_dict);

    auto diff = entity_diff(r_1, r_2);
    REQUIRE(diff.size() == 2);
    print_diffs(diff, r_1, r_2);

    auto runner_clone = r_2.clone();
    diff = entity_diff(r_2, *runner_clone);
    REQUIRE(diff.size() == 0);
}
