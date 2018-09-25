/// YANG Development Kit
// Copyright 2018 Cisco Systems. All rights reserved
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

#include <ydk/logger.hpp>
#include <ydk/common_utilities.hpp>

#include "gnmi_crud_service.hpp"
#include "gnmi_service.hpp"
#include "gnmi_provider.hpp"

using namespace std;

namespace ydk {

bool gNMICrudService::create(ServiceProvider & provider, Entity & entity)
{
    YLOG_INFO("Executing CRUD create operation on [{}]", entity.get_segment_path());
    entity.yfilter = YFilter::update;
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.set(gsp, entity);
}

bool gNMICrudService::create(ServiceProvider & provider, vector<Entity*> & entity_list)
{
    YLOG_INFO("Executing CRUD create operation on {}", entity_vector_to_string(entity_list));
    for (auto entity : entity_list)
        entity->yfilter = YFilter::update;
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.set(gsp, entity_list);
}

bool gNMICrudService::update(ServiceProvider & provider, Entity & entity)
{
    YLOG_INFO("Executing CRUD update operation on [{}]", entity.get_segment_path());
    entity.yfilter = YFilter::update;
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.set(gsp, entity);
}

bool gNMICrudService::update(ServiceProvider & provider, vector<Entity*> & entity_list)
{
    YLOG_INFO("Executing CRUD create operation on {}", entity_vector_to_string(entity_list));
    for (auto entity : entity_list)
        entity->yfilter = YFilter::update;
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.set(gsp, entity_list);
}

bool gNMICrudService::delete_(ServiceProvider & provider, Entity & entity)
{
    YLOG_INFO("Executing CRUD delete operation on [{}]", entity.get_segment_path());
    entity.yfilter = YFilter::delete_;
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.set(gsp, entity);
}

bool gNMICrudService::delete_(ServiceProvider & provider, vector<Entity*> & entity_list)
{
    YLOG_INFO("Executing CRUD delete operation on {}", entity_vector_to_string(entity_list));
    for (auto entity : entity_list)
        entity->yfilter = YFilter::delete_;
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.set(gsp, entity_list);
}

shared_ptr<Entity>
gNMICrudService::read(ServiceProvider & provider, Entity & filter)
{
    YLOG_INFO("Executing CRUD read operation on [{}]", filter.get_segment_path());
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.get(gsp, filter, "ALL");
}

vector<shared_ptr<Entity>>
gNMICrudService::read(ServiceProvider & provider, vector<Entity*> & filter_list)
{
    YLOG_INFO("Executing CRUD read operation on {}", entity_vector_to_string(filter_list));
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.get(gsp, filter_list, "ALL");
}

shared_ptr<Entity>
gNMICrudService::read_config(ServiceProvider & provider, Entity & filter)
{
    YLOG_INFO("Executing CRUD read_config operation on [{}]", filter.get_segment_path());
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.get(gsp, filter, "CONFIG");
}

vector<shared_ptr<Entity>>
gNMICrudService::read_config(ServiceProvider & provider, vector<Entity*> & filter_list)
{
    YLOG_INFO("Executing CRUD read operation on {}", entity_vector_to_string(filter_list));
    gNMIService gs;
    gNMIServiceProvider& gsp = dynamic_cast<gNMIServiceProvider&> (provider);
    return gs.get(gsp, filter_list, "CONFIG");
}

}
