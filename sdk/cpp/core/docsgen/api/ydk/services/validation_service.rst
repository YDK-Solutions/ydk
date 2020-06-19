..
  #  YDK - YANG Development Kit
  #  Copyright 2016 Cisco Systems. All rights reserved
  # *************************************************************
  # Licensed to the Apache Software Foundation (ASF) under one
  # or more contributor license agreements.  See the NOTICE file
  # distributed with this work for additional information
  # regarding copyright ownership.  The ASF licenses this file
  # to you under the Apache License, Version 2.0 (the
  # "License"); you may not use this file except in compliance
  # with the License.  You may obtain a copy of the License at
  #
  #   http:#www.apache.org/licenses/LICENSE-2.0
  #
  #  Unless required by applicable law or agreed to in writing,
  # software distributed under the License is distributed on an
  # "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  # KIND, either express or implied.  See the License for the
  # specific language governing permissions and limitations
  # under the License.
  # *************************************************************
  # This file has been modified by Yan Gorelik, YDK Solutions.
  # All modifications in original under CiscoDevNet domain
  # introduced since October 2019 are copyrighted.
  # All rights reserved under Apache License, Version 2.0.
  # *************************************************************

.. _ref-validationservice:

ValidationService
=================


.. cpp:class:: ydk::ValidationService

    Validation Service class for validating C++ model API objects of type :cpp:class:`Entity<Entity>`.

    .. cpp:enum-class:: Option

        All validation is performed in the context of some operation. These options capture the context of use.

        .. cpp:enumerator:: DATASTORE

            Datastore validation.

        .. cpp:enumerator:: GET_CONFIG

            Get config validation. Checks to see if only config nodes are references.

        .. cpp:enumerator:: GET

            Get validation.

        .. cpp:enumerator:: EDIT_CONFIG

            Edit validation. Checks on the values of leafs etc.

    .. cpp:function:: void validate(const ydk::path::ServiceProvider& provider, Entity& entity, ValidationService::Option option)

        Validates an entity based on the option.

        :param provider: An instance of :cpp:class:`ServiceProvider<ydk::path::ServiceProvider>`.
        :param entity: An instance of :cpp:class:`Entity<Entity>` class defined under a bundle.
        :param option: An instance of type :cpp:class:`Option<Option>`.
        :raises: :cpp:class:`YModelError<YModelError>`, if validation error was detected.
