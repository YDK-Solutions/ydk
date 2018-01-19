gNMI Service
============

.. contents:: Table of Contents

YDK gNMIService provides Set/Get/Subscribe functionalities.

.. py:class:: ydk.services.gNMIService

    Supports gNMI Set/Get operations on model API entities.

    .. py:method:: set(provider, entity, operation)

        Create the entity.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.providers.gNMIServiceProvider>`) Provider instance.
        :param entity: (:py:class:`Entity<ydk.types.Entity>`) Entity instance.
        :param operation: (``str``) Supported gNMI operations include: ``gnmi_create``, ``gnmi_delete``.
        :return: ``True`` if successful, ``False`` if not.
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if an error has occurred.

    .. py:method:: get(provider, read_filter)

        Read the entity.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.providers.gNMIServiceProvider>`) Provider instance.
        :param read_filter: (:py:class:`Entity<ydk.types.Entity>`) Read filter entity instance.
        :return: An instance of :py:class:`Entity<ydk.types.Entity>` as identified by the ``read_filter`` if successful, ``None`` if not.
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if an error has occurred.

    .. py:method:: subscribe(provider, subscribe_filter, list_mode, qos, mode, sample_interval, callback_function)

        Subscribe to telemetry updates.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.providers.gNMIServiceProvider>`) Provider instance.
        :param subscribe_filter: (:py:class:`Entity<ydk.types.Entity>`) Entity instance representing the path for which subscription is reqeusted.
        :param list_mode: (``str``) Subsciption list mode. One of ``STREAM``, ``ONCE`` or ``POLL``.
        :param qos: (``long``) QOS indicating the packet marking.
        :param mode: (``str``) Subscription mode. One of ``SAMPLE``, ``ON_CHANGE`` or ``TARGET_DEFINED``.
        :param sample_interval: (``int``) Sample interval for ``SAMPLE`` mode.
        :param callback_function: (``function``) Callback function with signature ``func(str)`` used to process the subscription data. The data is sent in the form of
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if an error has occurred.

Examples
--------

YDK Logging can be enabled per below:

.. code-block:: python
    :linenos:

    import logging
    logger = logging.getLogger("ydk")
    logger.setLevel(logging.INFO)
    handler = logging.StreamHandler()
    formatter = logging.Formatter(("%(asctime)s - %(name)s - %(levelname)s - %(message)s"))
    handler.setFormatter(formatter)
    logger.addHandler(handler)

GRPC trace can be enabled per below:

.. code-block:: shell
    :linenos:

    export GRPC_VERBOSITY=debug
    export GRPC_TRACE=transport_security

gNMI get example
~~~~~~~~~~~~~~~~

Example of instantiating and using objects of ``gNMIServiceProvider`` with ``gNMIService`` is shown below (assuming you have ``openconfig`` bundle installed). This assumes you have the certificate and key from the server (with a name like ``ems.pem`` and ``ems.key``) copied to the location of the below python script:

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_bgp
    from ydk.path import Repository
    from ydk.providers import gNMIServiceProvider
    from ydk.services import gNMIService

    gnmi_service = gNMIService()
    repository = Repository('/Users/test/yang_models_location')
    provider = gNMIServiceProvider(repo=repository, address='10.0.0.1', username='admin', password='admin')

    capabilities = provider.get_capabilities() # Get list of capabilities

    bgp = openconfig_bgp.Bgp()

    bgp_read = gnmi_service.get(provider, bgp) # Perform get operation

gNMI subscribe example
~~~~~~~~~~~~~~~~~~~~~~

Example of subscribing to telemetry using ``gNMIServiceProvider`` with ``gNMIService`` is shown below (assuming you have ``openconfig`` bundle installed). This assumes you have the certificate and key from the server (with a name like ``ems.pem`` and ``ems.key``) copied to the location of the below python script.

**NOTE:** The ``gNMIService`` class **has to be** used with ``multiprocessing.Pool`` for the ``subscribe`` operation as shown below as the subcription is a long-lived connection.

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_interfaces
    from ydk.path import Repository
    from ydk.providers import gNMIServiceProvider
    from ydk.services import gNMIService
    from ydk.filters import YFilter

    # Import the Pool class from multiprocessing module
    from multiprocessing import Pool

    # Callback function to handle telemetry data
    def print_telemetry_data(s):
        print(s)

    # Function to subscribe to telemetry data
    def subscribe(func):
        gnmi = gNMIService()
        repository = Repository('/Users/test/yang_models_location')
        provider = gNMIServiceProvider(repo=repository, address='10.0.0.1', username='admin', password='admin')

        # The below will create a telemetry subscription path 'openconfig-interfaces:interfaces/interface'
        interfaces = openconfig_interfaces.Interfaces()
        interface = openconfig_interfaces.Interfaces.Interface()
        interface.yfilter = YFilter.read
        interfaces.interface.append(interface)

        # Subscribe to updates in POLL mode. Polling subscriptions are used
        # for on-demand retrieval of data items via long-lived RPCs
        gnmi.subscribe(provider, interfaces, "POLL", 10, "SAMPLE", 100000, func)


    pool = Pool()
    pool.map(subscribe, [print_telemetry_data])
