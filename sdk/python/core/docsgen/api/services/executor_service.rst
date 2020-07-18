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

.. py:class:: ydk.services.ExecutorService

    Provides the functionality to execute RPCs

    .. py:method:: execute_rpc(self, provider, rpc_entity, top_entity=None)

        Create the entity

        :param provider: (:py:class:`ServiceProvider<ydk.path.ServiceProvider>`.) Provider instance.
        :param rpc_entity: (:py:class:`Entity<ydk.types.Entity>`) An instance of an RPC class defined under the ydk.models package or any of its subpackages.
        :param top_entity: (:py:class:`Entity<ydk.types.Entity>` optional)  Provide an instance of :py:class:`Entity<ydk.types.Entity>` only when expecting data to be returned.

        :return: An instance of :py:class:`Entity<ydk.types.Entity>` when provided top_entity or None otherwise
        :raises: :py:exc:`YServiceProviderError<ydk.errors.YServiceProviderError>`, if an error has occurred.

        Possible Errors:

            * a server side error
            * there isn't enough information in the entity to prepare the message (eg. missing keys)
