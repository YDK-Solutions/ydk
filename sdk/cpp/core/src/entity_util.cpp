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

#include "entity_util.hpp"
#include "errors.hpp"
#include "logger.hpp"

using namespace std;

namespace ydk
{
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// gNMI path utils
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
static void parse_entity_children(map<string, shared_ptr<Entity>> & children, vector<PathElem> & path_container);
static void parse_entity(Entity& entity, vector<PathElem> & path_container);

PathKey::PathKey(const std::string & name, const std::string & value)
        : name(name), value(value)
{
}

PathElem::PathElem(const std::string & path, std::vector<PathKey> keys)
        : path(path), keys(keys)
{
}

static void parse_entity(Entity& entity, vector<PathElem> & path_container)
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
    auto c = entity.get_children();
    parse_entity_children(c, path_container);
}

static void parse_entity_children(map<string, shared_ptr<Entity>> & children, vector<PathElem> & path_container)
{
    YLOG_DEBUG("Children count: {}", children.size());
    for(auto const& child : children)
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

void parse_entity_to_prefix_and_paths(Entity& entity, pair<string, string> & prefix, vector<PathElem> & path_container)
{
    EntityPath root_path = get_entity_path(entity, nullptr);
    auto s = root_path.path;
    YLOG_DEBUG("Got root path: {}", s);
    auto p = s.find(":");
    if(p != std::string::npos)
    {
        auto mod = s.substr(0, p);
        auto con = s.substr(p+1);
        prefix = make_pair(mod, con);
        YLOG_DEBUG("Got prefix: {}, {}", prefix.first, prefix.second);
    }
    auto c = entity.get_children();
    parse_entity_children(c, path_container);
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Entity utils
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

std::string get_relative_entity_path(const Entity* current_node, const Entity* ancestor, const std::string & path)
{
    std::ostringstream path_buffer;
    path_buffer << path;
    if(ancestor == nullptr)
    {
        throw(YCPPInvalidArgumentError{"ancestor should not be null."});
    }
    auto p = current_node->parent;
    std::vector<Entity*> parents {};
    while (p != nullptr && p != ancestor) {
        parents.push_back(p);
        p = p->parent;
    }

    if (p == nullptr) {
        throw(YCPPInvalidArgumentError{"parent is not in the ancestor hierarchy."});
    }

    std::reverse(parents.begin(), parents.end());

    p = nullptr;
    for (auto p1 : parents) {
        if (p) {
            path_buffer << "/";
        } else {
             p = p1;
        }
        path_buffer << p1->get_segment_path();
    }
    if(p)
        path_buffer << "/";
    path_buffer<<current_node->get_segment_path();
    return path_buffer.str();

}

bool is_set(const YFilter & yfilter)
{
    return yfilter != YFilter::not_set;
}

static const EntityPath get_entity_path(const Entity & entity, const string & path_buffer)
{
    std::vector<std::pair<std::string, LeafData> > leaf_name_data  = entity.get_name_leaf_data();

    EntityPath entity_path {path_buffer, leaf_name_data};
    return entity_path;
}

static bool is_absolute_path(Entity* ancestor)
{
    return ancestor == nullptr;
}

//
// @brief Get the EntityPath relative to the parent passed in
//
// Returns the EntityPath relative to the ancestor passed in.
// The ancestor must either be null, in which case the absolute path
// from the root is returned, or some other ancestor of this Entity.
//
// @param[in] parent The ancestor relative to which the path is calculated or nullptr
// @return EntityPath
// @throws YCPPInvalidArgumentError if the parent is invalid

const EntityPath get_entity_path(const Entity & entity, Entity* ancestor)
{
    std::ostringstream path_buffer;
    if (is_absolute_path(ancestor))
    {
        if(entity.has_list_ancestor)
        {
            throw(YCPPInvalidArgumentError{"ancestor for entity cannot be nullptr as one of the ancestors is a list. Path: "+entity.get_segment_path()});
        }
        auto a = entity.get_absolute_path();
        if(a.size() == 0)
        {
            path_buffer << entity.get_segment_path();
        }
        else
        {
            path_buffer << a;
        }
    }
    else
    {
        if(entity.is_top_level_class)
        {
            throw(YCPPInvalidArgumentError{"ancestor has to be nullptr for top-level node. Path: "+entity.get_segment_path()});
        }
        path_buffer << get_relative_entity_path(&entity, ancestor, path_buffer.str());
    }
    return get_entity_path(entity, path_buffer.str());
}

}
