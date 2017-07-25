//
// @file path_api.hpp
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

#ifndef _SERVICE_PROVIDER_H_
#define _SERVICE_PROVIDER_H_

#include "path_api.hpp"

namespace ydk {

///
/// @brief Interface for all ServiceProvider implementations
///
/// Concretes instances of ServiceProviders are expected to extend this interface.
///
class ServiceProvider
{
public:
    ///
    /// @brief return the SchemaTree supported by this instance of the ServiceProvider
    ///
    /// @return pointer to the RootSchemaNode or nullptr if one could not be created
    ///

    virtual ~ServiceProvider();

    virtual EncodingFormat get_encoding() const = 0;

    virtual Session get_session() const = 0;

};
}
#endif /*_SERVICE_PROVIDER_H_*/
