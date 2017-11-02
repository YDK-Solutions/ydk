gNMI Service Provider
========================


.. py:class:: ydk.providers.gNMIServiceProvider(repo, address, is_secure)

    Constructs an instance of the ``gNMIServiceProvider`` to connect to a `gNMI <https://github.com/openconfig/gnmi>`_ server. This has clean up methods implemented in its destructor. So, the user does not need to worry about closing the gNMI session.

    :param repo: Instance of :py:class:`Repository<ydk.path.Repository>` with path to local directory containing all the yang models supported on the gNMI server
    :param address: IP address & port of the device supporting a gNMI interface in the format ``ip-address:port``
    :param is_secure: Indicates usage of gNMI secure channel. **Note:** When using ``is_secure=True``, the certificate file from the server (with a name like ``ems.pem``) needs to be copied to the location of your python script

    .. py:method:: get_capabilities()

        Returns list of capabilities supported by the device


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

Example of instantiating and using objects of ``gNMIServiceProvider`` is shown below (assuming you have ``openconfig`` bundle installed). This assumes you have the certificate file from the server (with a name like ``ems.pem``) copied to the location of the below python script:

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_bgp
    from ydk.path import Repository
    from ydk.providers import gNMIServiceProvider
    from ydk.services import CRUDService

    crud = CRUDService()
    repository = Repository('/Users/test/yang_models_location') # Location of directory where yang models from the server are downloaded
    provider = gNMIServiceProvider(repo=repository, address='10.0.0.1:57400', is_secure=True)

    capabilities = provider.get_capabilities() # Get list of capabilities

    bgp = openconfig_bgp.Bgp()

    bgp_read = crud.read(provider, bgp) # Perform read operation
