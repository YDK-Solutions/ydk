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

Opendaylight Service Provider
=============================

.. py:class:: ydk.providers.OpendaylightServiceProvider(repo, address, username, password, port, encoding)

    A service provider to be used to communicate with an OpenDaylight instance.

    :param repo: (:py:class:`Repository<ydk.path.Repository>`) User provided repository stores cached models.
    :param address: (``str``) IP address of the ODL instance
    :param username: (``str``) Username to log in to the instance
    :param password: (``str``) Password to log in to the instance
    :param port: (``int``) Device port used to access the ODL instance.
    :param encoding: (:py:class:`EncodingFormat<ydk.types.EncodingFormat>`) Type of encoding to be used for the payload. Default is :py:attr:`JSON<ydk.types.EncodingFormat.JSON>`

    .. py:method:: get_node_provider(node_id)

        Returns the ServiceProvider instance corresponding to the device being controlled by the OpenDaylight instance, indicated by ``node_id``

        :param node_id: (``str``) The name of the device being controlled by the OpenDaylight instance.
        :return: One of supported service provider instance.
        :raises: :py:exc:`YServiceProviderError<ydk.errors.YServiceProviderError>` if no such service provider could be found.

    .. py:method:: get_node_ids()

        Returns a list of node ID’s of the devices being controlled by this OpenDaylight instance.

        :return: List of node ID’s of the devices being controlled by this OpenDaylight instance.
