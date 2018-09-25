//
// @file gnmi_crud_service.hpp
// @brief Implementation of CRUD service for gNMI.
//
// YANG Development Kit
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

#ifndef GNMI_CRUD_SERVICE_HPP
#define GNMI_CRUD_SERVICE_HPP

#include <ydk/crud_service.hpp>

namespace ydk
{

class gNMICrudService : public CrudService
{
  public:
    bool create(ServiceProvider & provider, Entity & entity);
    bool create(ServiceProvider & provider, std::vector<Entity*> & entity_list);

    bool update(ServiceProvider & provider, Entity & entity);
    bool update(ServiceProvider & provider, std::vector<Entity*> & entity_list);

    bool delete_(ServiceProvider & provider, Entity & entity);
    bool delete_(ServiceProvider & provider, std::vector<Entity*> & entity_list);

    std::shared_ptr<Entity> read(ServiceProvider & provider, Entity & filter);
    std::vector<std::shared_ptr<Entity>>
        read(ServiceProvider & provider, std::vector<Entity*> & filter_list);

    std::shared_ptr<Entity> read_config(ServiceProvider & provider, Entity & filter);
    std::vector<std::shared_ptr<Entity>>
        read_config(ServiceProvider & provider, std::vector<Entity*> & filter_list);
};

}

#endif /* GNMI_CRUD_SERVICE_HPP */
