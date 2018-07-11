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

#include <algorithm>
#include <vector>

#include "gnmi_util.hpp"
#include "errors.hpp"
#include "logger.hpp"

using namespace std;

namespace ydk
{

PathKey::PathKey(const std::string & name, const std::string & value)
        : name(name), value(value)
{
}

PathElem::PathElem(const std::string & path, std::vector<PathKey> keys)
        : path(path), keys(keys)
{
}

static void parse_entity_children(Entity & entity, vector<PathElem> & path_container);

static void parse_entity(Entity & entity, vector<PathElem> & path_container)
{
    EntityPath path = get_entity_path(entity, entity.parent);
    auto s = entity.get_segment_path();
    vector<PathKey> keys;
    YLOG_DEBUG("Got path {}", s);
    auto p = s.find("[");
    if(p != std::string::npos)
    {
        s = s.substr(0, p);
        for(const pair<string, LeafData> & name_value : path.value_paths)
        {
            LeafData leaf_data = name_value.second;
            if(leaf_data.is_set)
            {
                YLOG_DEBUG("Creating key {} with value: '{}'", name_value.first, leaf_data.value);
                PathKey key{name_value.first, leaf_data.value};
                keys.push_back(key);
            }
         }
    }

    path_container.push_back({s, keys});
    parse_entity_children(entity, path_container);
}

static void parse_entity_children(Entity & entity, vector<PathElem> & path_container)
{
    for(auto const& child : entity.get_children())
    {
        if(child.second == nullptr)
            continue;
        YLOG_DEBUG("==================");
        YLOG_DEBUG("Looking at child '{}': {}",child.first, get_entity_path(*(child.second), child.second->parent).path);
        if(child.second->has_operation() || child.second->has_data() || child.second->is_presence_container)
            parse_entity(*(child.second), path_container);
        else
            YLOG_DEBUG("Child has no data and no operations");
    }
}

void parse_entity_prefix(Entity& entity, pair<string, string> & prefix)
{
    EntityPath root_path = get_entity_path(entity, nullptr);
    auto s = root_path.path;
    string mod = s;
    string con = {};
    auto p = s.find(":");
    if (p != std::string::npos)
    {
        mod = s.substr(0, p);
        con = s.substr(p+1);
        YLOG_DEBUG("Got entity prefix: '{}:{}'", mod, con);
    }
    else {
        YLOG_DEBUG("Got unexpected entity root path: '{}'", s);
    }
    prefix = make_pair(mod, con);
}

void parse_entity_to_prefix_and_paths(Entity& entity, pair<string, string> & prefix, vector<PathElem> & path_container)
{
    parse_entity_prefix(entity, prefix);

    parse_entity_children(entity, path_container);
}
static void parse_entity_children(Entity & entity, gnmi::Path* path);

static void parse_entity(Entity & entity, gnmi::Path* path)
{
    EntityPath entity_path = get_entity_path(entity, entity.parent);
    auto s = entity.get_segment_path();
    map<string,string> keys{};

    auto p = s.find("[");
    if (p != std::string::npos) {
        s = s.substr(0, p);
        for (const pair<string, LeafData> & name_value : entity_path.value_paths) {
            LeafData leaf_data = name_value.second;
            if (leaf_data.is_set) {
                YLOG_DEBUG("Creating key '{}' with value: '{}'", name_value.first, leaf_data.value);
                keys[name_value.first] = leaf_data.value;
            }
         }
    }
    gnmi::PathElem* elem = path->add_elem();
    elem->set_name(s);
    for (auto key : keys) {
        auto key_map = elem->mutable_key();
        (*key_map)[key.first] = key.second;
    }

    parse_entity_children(entity, path);
}

static void parse_entity_children(Entity & entity, gnmi::Path* path)
{
    for (auto const & child : entity.get_children())
    {
        if (child.second == nullptr)
            continue;
        YLOG_DEBUG("gnmi_util::parse_entity_children:");
        YLOG_DEBUG("Looking at child '{}': '{}'", child.first, get_entity_path(*(child.second), child.second->parent).path);
        if (child.second->has_operation() || child.second->has_data() || child.second->is_presence_container)
            parse_entity(*(child.second), path);
        else
            YLOG_DEBUG("Child has no data and no operations");
    }
}

void parse_entity_prefix(Entity& entity, gnmi::Path* path)
{
    // Add origin and first container to the path
	EntityPath root_path = get_entity_path(entity, nullptr);
    auto s = root_path.path;
    parse_prefix_to_path(s, path);
}

void parse_entity_to_path(Entity& entity, gnmi::Path* path)
{
    // Add origin and first container to the path
	parse_entity_prefix(entity, path);

    // Add children path
    parse_entity_children(entity, path);
}



static void add_path_elem(gnmi::Path* path, string s)
{
    map<string,string> keys{};

    auto p = s.find("[");
    if (p != std::string::npos) {
        string key_path = s.substr(p);
        s = s.substr(0, p);
        size_t open_bracket_pos = 0;
        while (open_bracket_pos != string::npos) {
            auto equal_pos = key_path.find("=", open_bracket_pos);
        	auto close_bracket_pos = key_path.find("]", equal_pos);
            string key_name = key_path.substr(open_bracket_pos+1, equal_pos-open_bracket_pos-1);
            string key_value = key_path.substr(equal_pos+2, close_bracket_pos-equal_pos-3);
            keys[key_name] = key_value;
            if (close_bracket_pos == key_path.length()-1)
                break;
            else
                open_bracket_pos = key_path.find("[", close_bracket_pos+1);
        }
    }
    gnmi::PathElem* elem = path->add_elem();
    elem->set_name(s);
    for (auto key : keys) {
        auto key_map = elem->mutable_key();
        (*key_map)[key.first] = key.second;
    }
}

void parse_prefix_to_path(const string& prefix, gnmi::Path* path)
{
    // Add origin and first container to the path
    string mod = prefix;
    string con = {};
    auto p = prefix.find(":");
    if (p != string::npos) {
        mod = prefix.substr(0, p);
        con = prefix.substr(p+1);
        YLOG_DEBUG("parse_prefix_to_path: Got data node path prefix: '{}:{}'", mod, con);
        path->set_origin(mod);
        add_path_elem(path, con);
    }
    else {
        YLOG_DEBUG("parse_prefix_to_path: Got unexpected data node path: '{}', missing prefix.", prefix);
        add_path_elem(path, mod);
    }
}

static path::DataNode* get_last_datanode(path::DataNode* dn)
{
    auto children = dn->get_children();
    if (!children.empty()) {
        // Select last child
        return get_last_datanode(children[children.size()-1].get());
    }
    return dn;
}

void parse_datanode_to_path(path::DataNode* dn, gnmi::Path* path)
{
    path::DataNode* last_datanode = get_last_datanode(dn);
    string full_path = last_datanode->get_path();
    std::vector<std::string> segments = path::segmentalize(full_path);

    // Add origin and first container to the path
    auto s = segments[1];
    parse_prefix_to_path(s, path);

    // Add the rest of the segments to the path
    for (size_t i=2; i < segments.size(); i++) {
    	add_path_elem(path, segments[i]);
    }
}

}
