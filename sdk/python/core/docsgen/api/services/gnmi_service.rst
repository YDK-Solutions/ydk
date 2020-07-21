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

.. _gnmi_service:

gNMI Service API
================

.. contents:: Table of Contents

.. py:class:: ydk.gnmi.services.gNMIService

    Supports gNMI Set/Get/Subscribe operations on model API entities. It also allows to get gNMI server capabilities.

    .. py:method:: set(provider, entity)

        Create, update, or delete single entity or multiple entities in the server configuration.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.gnmi.providers.gNMIServiceProvider>`) Provider instance.
        :param entity: (:py:class:`Entity<ydk.types.Entity>`) instance, which represents single container in device supported model.
                         Each **Entity** instance must be annotated with :py:class:`YFilter<ydk.filters.YFilter>`, which defines set operation:

                         * YFilter.replace - add new configuration or replace the whole configuration tree

                         * YFilter.update  - update or create configuration in existing tree

                         * YFilter.delete  - delete part or entire configuration tree

                         For multiple containers the :py:class:`Entity<ydk.types.Entity>` instances must be encapsulated into Python ``list`` or YDK ``EntityCollection`` :py:class:`Config<ydk.types.Config>`.
        :return: ``True`` if successful, ``False`` if not.
        :raises: :py:exc:`YServiceProviderError<ydk.errors.YServiceProviderError>` if an error has occurred.

    .. py:method:: get(provider, read_filter, read_mode='CONFIG')

        Read the entity.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.gnmi.providers.gNMIServiceProvider>`) Provider instance.
        :param read_filter: (:py:class:`Entity<ydk.types.Entity>`) instance, which represents single container in device supported model.

                              For multiple containers the :py:class:`Entity<ydk.types.Entity>` instances must be encapsulated into Python ``list`` or YDK ``EntityCollection`` :py:class:`Filter<ydk.types.Filter>`.
        :param read_mode: (``str``) One of the values: ``CONFIG``, ``STATE``, ``OPERATIONAL``, or ``ALL``.
        :return: For single entity filter - an instance of :py:class:`Entity<ydk.types.Entity>` as identified by the **read_filter** or ``None``, if operation fails.

                 For multiple filters - collection of :py:class:`Entity<ydk.types.Entity>` instances encapsulated into Python ``list`` or YDK ``EntityCollection`` :py:class:`Config<ydk.types.Config>` accordingly to the type of **read_filter**.
        :raises: :py:exc:`YServiceProviderError<ydk.errors.YServiceProviderError>` if an error has occurred.

    .. py:method:: subscribe(provider, subscription, qos=0, mode='ONCE, encoding='PROTO', callback_function=None)

        Subscribe to telemetry updates.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.gnmi.providers.gNMIServiceProvider>`) Provider instance.
        :param subscription: (:py:class:`gNMISubscription<ydk.gnmi.services.gNMISubscription>`) Single instance or Python ``list`` of instances of objects, which represent the subscription.
        :param qos: (``long``) QOS indicating the packet marking.
        :param mode: (``str``) Subscription mode: one of ``STREAM``, ``ONCE`` or ``POLL``.
        :param encoding: (``str``) Encoding method for the output: one of ``JSON``, ``BYTES``, ``PROTO``, ``ASCII``, or ``JSON_IETF``.
        :param callback_function: (``func(str)``) Callback function, which is used to process the subscription data.
                                  The subscription data returned to the user as a string representation of protobuf **SubscribeResponse** message.
                                  If not specified, the response is printed to system stdout.
        :raises: :py:exc:`YServiceProviderError<ydk.errors.YServiceProviderError>` if an error has occurred.

    .. py:method:: capabilities(provider)

        Get gNMI server capabilities

        :param provider: (:py:class:`gNMIServiceProvider<ydk.gnmi.providers.gNMIServiceProvider>`) Provider instance.
        :return: (``str``) JSON encoded string, which represents gNMI server capabilities.

.. py:class:: ydk.gnmi.services.gNMISubscription

        Instance of this class defines subscription for a single entity. Members of the class are:

        * entity: (:py:class:`Entity<ydk.types.Entity>`) Instance of the subscription entity. This parameter must be set by the user.
        * subscription_mode: (``str``) Expected one of the following string values: ``TARGET_DEFINED``, ``ON_CHANGE``, or ``SAMPLE``; default value is ``ON_CHANGE``.
        * sample_interval: (``long``) Time interval in nanoseconds between samples in ``STREAM`` mode; default value is 60000000000 (1 minute).
        * suppress_redundant: (``bool``) Indicates whether values that not changed should be sent in a ``STREAM`` subscription; default value is ``False``
        * heartbeat_interval: (``long``) Specifies the maximum allowable silent period in nanoseconds when **suppress_redundant** is True. If not specified, the **heartbeat_interval** is set to 360000000000 (10 minutes) or **sample_interval** whatever is bigger.

gNMI Service Examples
=====================

To enable YDK logging include the following code:

.. code-block:: python
    :linenos:

    import logging
    logger = logging.getLogger("ydk")
    logger.setLevel(logging.INFO)
    handler = logging.StreamHandler()
    formatter = logging.Formatter(("%(asctime)s - %(name)s - %(levelname)s - %(message)s"))
    handler.setFormatter(formatter)
    logger.addHandler(handler)

To enable GRPC trace set environment variables as followed:

.. code-block:: shell
    :linenos:

    export GRPC_VERBOSITY=debug
    export GRPC_TRACE=transport_security

gNMI 'set' example
~~~~~~~~~~~~~~~~~~

Example of instantiating and using objects of ``gNMIServiceProvider`` with ``gNMIService`` is shown below (assuming you have ``openconfig`` bundle installed).

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_bgp
    from ydk.path import Repository
    from ydk.gnmi.providers import gNMIServiceProvider
    from ydk.gnmi.services import gNMIService

    repository = Repository('/Users/test/yang_models_location')
    provider = gNMIServiceProvider(repo=repository, address='10.0.0.1', port=57400, username='admin', password='admin')
    gnmi_service = gNMIService()

    # Create entire BGP configuration
    bgp = openconfig_bgp.Bgp()
    bgp.global_.config.as_ = 65172
    neighbor = bgp.Neighbors.Neighbor()
    neighbor.neighbor_address = '172.16.255.2'
    neighbor.config.neighbor_address = '172.16.255.2'
    neighbor.config.peer_as = 65172
    bgp.neighbors.neighbor.append(neighbor)

    bgp.yfilter = YFilter.replace	# Define set/create operation

    ok = gnmi_service.set(provider, bgp) # Perform create operation

    # Delete one neighbor
    bgp = openconfig_bgp.Bgp()
    neighbor = bgp.Neighbors.Neighbor()
    neighbor.neighbor_address = '172.16.255.2'
    bgp.neighbors.neighbor.append(neighbor)

    bgp.yfilter = YFilter.delete	# Define set/delete operation

    ok = gnmi_service.set(provider, bgp) # Perform delete operation

gNMI 'get' example
~~~~~~~~~~~~~~~~~~

Example of instantiating and using objects of ``gNMIServiceProvider`` with ``gNMIService`` is shown below (assuming you have ``openconfig`` bundle installed).

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_bgp
    from ydk.path import Repository
    from ydk.gnmi.providers import gNMIServiceProvider
    from ydk.gnmi.services import gNMIService

    repository = Repository('/Users/test/yang_models_location')
    provider = gNMIServiceProvider(repo=repository, address='10.0.0.1', port=57400, username='admin', password='admin')
    gnmi_service = gNMIService()

    capabilities = provider.get_capabilities() # Get list of capabilities

    bgp = openconfig_bgp.Bgp()

    bgp_read = gnmi_service.get(provider, bgp) # Perform get operation

gNMI 'subscribe' example
~~~~~~~~~~~~~~~~~~~~~~~~

Example of subscribing to telemetry using ``gNMIServiceProvider`` with ``gNMIService`` is shown below (assuming you have ``openconfig`` bundle installed).

**NOTE:** The ``gNMIService`` class **can be** used with ``multiprocessing.Pool`` for the ``subscribe`` operation as shown below as the subcription is a long-lived connection.

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_interfaces
    from ydk.path import Repository
    from ydk.gnmi.providers import gNMIServiceProvider
    from ydk.gnmi.services import gNMIService, gNMISubscription

    import datetime

    # Callback function to handle telemetry data
    def print_telemetry_data(s):
        if 'update' in s:
            current_dt = datetime.datetime.now()
            print("\n===> Received subscribe response at %s: \n%s" %
                  (current_dt.strftime("%Y-%m-%d %H:%M:%S"), s))

    # Function to subscribe to telemetry data
    def subscribe(func):
        # Initialize gNMI Service and Provider
        gnmi = gNMIService()
        repository = Repository('/home/yan/ydk-workspace/ydk-gen/scripts/repository/10.30.110.84')
        provider = gNMIServiceProvider(repo=repository, address='10.30.110.85', port=57400,
                                       username='admin', password='admin')

        # Create telemetry subscription entity
        interfaces = openconfig_interfaces.Interfaces()
        interface = openconfig_interfaces.Interfaces.Interface()
        interface.name = '"MgmtEth0/RP0/CPU0/0"'
        interfaces.interface.append(interface)

        # Build subscription
        subscription = gNMISubscription()
        subscription.entity = interfaces
        subscription.subscription_mode = 'SAMPLE'
        subscription.sample_interval = 10 * 1000000000
        subscription.suppress_redundant = False
        subscription.heartbeat_interval = 100 * 1000000000

        # Subscribe for updates in STREAM mode.
        gnmi.subscribe(provider, subscription, 10, 'STREAM', "PROTO", func)

    if __name__ == "__main__":
        subscribe(print_telemetry_data)
