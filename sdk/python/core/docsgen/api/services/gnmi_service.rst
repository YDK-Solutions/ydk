gNMI Service
============


YDK gNMIService provides Set/Get functionalities.

.. py:class:: ydk.services.gNMIService

    Supports gNMI Set/Get operations on model API entities.

    .. py:method:: set(provider, entity, operation)

        Create the entity.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.providers.gNMIServiceProvider>`) Provider instance.
        :param entity: (:py:class:`Entity<ydk.types.Entity>`) Entity instance.
        :param operation: (``str``) Supported gNMI operations include: ``create``, ``update`` and ``delete``.
        :return: ``True`` if successful, ``False`` if not.
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if an error has occurred.

    .. py:method:: get(provider, read_filter)

        Read the entity.

        :param provider: (:py:class:`gNMIServiceProvider<ydk.providers.gNMIServiceProvider>`) Provider instance.
        :param read_filter: (:py:class:`Entity<ydk.types.Entity>`) Read filter entity instance.
        :return: An instance of :py:class:`Entity<ydk.types.Entity>` as identified by the ``read_filter`` if successful, ``None`` if not.
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

Example of instantiating and using objects of ``gNMIServiceProvider`` with ``gNMIService`` is shown below (assuming you have ``openconfig`` bundle installed). This assumes you have the certificate file from the server (with a name like ``ems.pem``) copied to the location of the below python script:

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_bgp
    from ydk.path import Repository
    from ydk.providers import gNMIServiceProvider
    from ydk.services import gNMIService

    gnmi_service = gNMIService()
    repository = Repository('/Users/test/yang_models_location')
    provider = gNMIServiceProvider(repo=repository, address='10.0.0.1:57400', is_secure=True)

    capabilities = provider.get_capabilities() # Get list of capabilities

    bgp = openconfig_bgp.Bgp()

    bgp_read = gnmi_service.get(provider, bgp) # Perform get operation
