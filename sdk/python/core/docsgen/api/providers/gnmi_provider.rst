gNMI Service Provider
========================


.. py:class:: ydk.gnmi.providers.gNMIServiceProvider(repo, address, port=57400, username, password, server_certificate="", private_key="")

    A service provider, which allows YDK connect to a `gNMI <https://github.com/openconfig/gnmi>`_ server. By default, the provider works in non-secure mode (tls is off). In order to enable secure mode connection the user must provide the gNMI server certificate of authorization (public key) and optionally the client (YDK application host) private key.

    :param repo: Instance of :py:class:`Repository<ydk.path.Repository>` with path to local directory containing the the `ydk yang model <https://raw.githubusercontent.com/CiscoDevNet/ydk-gen/1344b3f22d746764f17536ac4e666836de4ba84d/sdk/cpp/core/tests/models/ydk%402016-02-26.yang>`_ along with all the yang models supported on the gNMI server.
    :param address: (``str``) Host address of the device supporting a gNMI interface
    :param port: (``int``)Port on which the gNMI interface can be accessed on the device. If not specified, the default value of ``57400`` is assigned.
    :param username: (``str``) Username.
    :param password: (``str``) Password.
    :param server_certificate: (``str``) Full path to a file, which contains server certificate of authorization (public key). If not specified, it is assumed non-secure connection to gNMI server.
    :param private_key: (``str``) Full path to a file, which contains private key of the application host. If not specified and **server_certificate** is defined (secure connection), the GRPC internally defined private key is used.

Examples
--------

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

To get server certificate, you need to retrive the **ems.pem** file from the IOS XR device (after enabling gRPC/TLS)
and put it to YDK application host in a folder, which is specific to the device. 
You can find the **ems.pem** file on the router in either **/misc/config/grpc/** or **/var/xr/config/grpc** folder.

Example of instantiating and using objects of ``gNMIServiceProvider`` is shown below (assuming you have ``openconfig`` bundle installed).

.. code-block:: python
    :linenos:

    from ydk.models.openconfig import openconfig_bgp
    from ydk.path import Repository
    from ydk.gnmi.providers import gNMIServiceProvider
    from ydk.services import CRUDService

    # Create repository with location of directory where yang models from the server are downloaded
    repository = Repository('/Users/test/yang_models_location')
    ca_file_path = '/Users/test/device-ip/ems.pem'

    # Instantiate provider with non-secure connection to gNMI server
    provider = gNMIServiceProvider(repo=repository, address='10.0.0.1', port=57400,
                                   username='admin', password='admin', server_certificate=ca_file_path)

    # Define filter
    bgp = openconfig_bgp.Bgp()

    # Run CRUD read opearaion
    crud = CRUDService()
    bgp_read = crud.read(provider, bgp) # Perform read operation
