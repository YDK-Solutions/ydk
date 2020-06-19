..
  #  YDK-YANG Development Kit
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

Executor Service
================

.. cpp:class:: ydk::ExecutorService

    Executor Service class for supporting execution of RPCs.

    .. cpp:function:: ExecutorService()

    .. cpp:function:: std::shared_ptr<Entity> execute_rpc(ydk::ServiceProvider & provider, Entity & rpc_entity, std::shared_ptr<Entity> top_entity = nullptr)

        Create the rpc entity.

        :param provider: An instance of :cpp:class:`ServiceProvider<ServiceProvider>`.
        :param rpc_entity: An rpc instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :param top_entity: Optional arg that should be provided when expecting data to be returned.
        :return: The requested data as :cpp:class:`Entity<ydk::Entity>` or nullptr if N/A.
        :raises: :cpp:class:`YServiceProviderError<YServiceProviderError>`, if an error has occurred.
