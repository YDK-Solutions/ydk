Introduction
============
.. contents:: Table of Contents

YDK consists of two main components: core library, which consists of services and providers, and python model API, which are APIs generated based on YANG models and packaged as bundles.

Core library consists of the below:

 * **Service:** Provides simple API interface to be used with the bindings and providers
 * **ServiceProvider:** Provides concrete implementation that abstracts underlying protocol details (e.g. :py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`, which is based on the NETCONF protocol)

Applications can be written using the python model API in conjunction with a service and a provider.

Writing an app
--------------

In this example, we set some BGP configuration using the OpenConfig model, the CRUD (Create/Read/Update/Delete) service and the NETCONF service provider. The example in this document is a simplified version of the more complete sample that is available in ``samples/bgp.py``. The more complete sample can be run with the below steps:

.. code-block:: sh

    (ydk-py)ydk-py$ cd core/samples
    (ydk-py)samples$ ./bgp.py -h
    Usage: bgp.py [-h | --help] [options]

    Options:
    -h, --help            show this help message and exit
    -v VERSION, --version=VERSION
                        force NETCONF version 1.0 or 1.1
    -u USERNAME, --user=USERNAME
    -p PASSWORD, --password=PASSWORD
                        password
    --proto=PROTO         Which transport protocol to use, one of ssh or tcp
    --host=HOST           NETCONF agent hostname
    --port=PORT           NETCONF agent SSH port

    (ydk-py)samples$ ./bgp.py --host <ip-address-of-netconf-server> -u <username> -p <password> --port <port-number>

What happens underneath
-----------------------
YDK performs the below actions when running this application:

 1. Establish a session with the device
 2. Encode python data objects to the protocol format (e.g. netconf XML payload)
 3. Perform transport operation with the device and collect the response (e.g. netconf reply)
 4. Decode response as python object and return the result to app
 5. Raise python exceptions for any errors that occurred


Service providers
-----------------
The first step in any application is to create a service provider instance. In this case, the NETCONF service provider (defined in :py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) is responsible for mapping between the CRUD service API and the underlying manageability protocol (NETCONF RPCs).

We instantiate an instance of the service provider that creates a NETCONF session to the machine with address 10.0.0.1:

.. code-block:: python
    :linenos:

    from ydk.providers import NetconfServiceProvider

    sp_instance = NetconfServiceProvider(address='10.0.0.1',
                                         port=830,
                                         username='test',
                                         password='test',
                                         protocol='ssh')

Using the model APIs
--------------------
After establishing the connection, we instantiate the entities and set some data. First, we import the types from the OpenConfig BGP module:

.. code-block:: python
    :linenos:
    :lineno-start: 8

    from ydk.models.openconfig import openconfig_bgp
    from ydk.models.openconfig import openconfig_bgp_types

Next, create a :py:class:`Bgp<ydk.models.openconfig.openconfig_bgp.Bgp>` configuration object and set the attributes:

.. code-block:: python
    :linenos:
    :lineno-start: 10

    # create BGP object
    bgp_cfg = openconfig_bgp.Bgp()

    # set the Global AS
    bgp_cfg.global_.config.as_ = 65001

    # Create an AFI SAFI config
    ipv4_afsf = bgp_cfg.global_.afi_safis.AfiSafi()
    ipv4_afsf.afi_safi_name = openconfig_bgp_types.Ipv4Unicast()
    ipv4_afsf.config.afi_safi_name = openconfig_bgp_types.Ipv4Unicast()
    ipv4_afsf.config.enabled = True

    # Add the AFI SAFI config to the global AFI SAFI list
    bgp_cfg.global_.afi_safis.afi_safi.append(ipv4_afsf)

Invoking the CRUD Service
-------------------------
The CRUD service provides methods to create, read, update and delete entities on a device making use of the session provided by a service provider (NETCONF in this case).  In order to use the CRUD service, we need to import the :py:class:`CRUDService<ydk.services.CRUDService>` class:

.. code-block:: python
    :linenos:
    :lineno-start: 24

    from ydk.services import CRUDService

Next, we instantiate the CRUD service:

.. code-block:: python
    :linenos:
    :lineno-start: 25

    crud_service = CRUDService()

Finally, we invoke the create method of the in this case).  In order to use the CRUD service, we need to import the :py:class:`CRUDService<ydk.services.CRUDService>` class passing in the
service provider instance and our entity (``bgp_cfg``):

.. code-block:: python
    :linenos:
    :lineno-start: 26

    try:
        crud_service.create(sp_instance, bgp_cfg)
    except YError:

Note if there were any errors the above API will raise a YError exception.

Using non-top level objects
---------------------------

In the example above you noticed that we started building model from top-level object - :py:class:`Bgp<ydk.models.openconfig.openconfig_bgp.Bgp>` and then built the object tree down the hierarchy. 
However in certain conditions we can build independently non-top level objects and still be able to do all the CRUD operations.

Top level object vs. non-top
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
The top level object represents top-level container in the Yang model. Examples of top-level objects:

 * openconfig_bgp.Bgp
 * openconfig_interfaces.Interfaces

The non-top level object represents a container in the Yang model, which is located under top level container. A member of a non-top level list can also be considered as non-top level object.
Examples of non-top level objects:

 * openconfig_bgp.Bgp.Global\_.AfiSafis.AfiSafi
 * openconfig_bgp.Bgp.Neighbors
 * openconfig_bgp.Bgp.Neighbors.Neighbor
 * openconfig_bgp.Bgp.Neighbors.Neighbor.Config
 * openconfig_interfaces.Interfaces.Interface

How to use non-top level objects
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You should be able to work with non-top level objects similarly as with top level. 
Your program will look more simple and straight to the point.
The above example will look now like this:

.. code-block:: python
 :linenos:

    from ydk.models.openconfig import openconfig_bgp
    from ydk.models.openconfig import openconfig_bgp_types
 
    # Create an AFI SAFI config
    ipv4_afsf = openconfig_bgp.Bgp.Global_.AfiSafis.AfiSafi()
    ipv4_afsf.afi_safi_name = openconfig_bgp_types.Ipv4Unicast()
    ipv4_afsf.config.afi_safi_name = openconfig_bgp_types.Ipv4Unicast()
    ipv4_afsf.config.enabled = True

    crud_service = CRUDService()
    crud_service.create(sp_instance, ipv4_afsf)
    
    # Read single AFI SAFI config
    afisafiFilter = openconfig_bgp.Bgp.Global_.AfiSafis.AfiSafi()
    afisafiFilter.afi_safi_name = openconfig_bgp_types.IPV6UNICAST{}

    afisafi = crud.ReadConfig(sp_instance, afisafiFilter)

Limitations
~~~~~~~~~~~

Not all non-top level objects can be used independently. Here is the rule:

  When building non-top level object, we have to define all the list keys on the way up to the top level object. 
  In the example above the object `ipv4_afsf` is a member of the list. We can use it as long as its key `afi_safi_name` is defined. 
  
Under the hood
~~~~~~~~~~~~~~

The programmability protocols like Netconf, gNMI, etc. are always working with top level model objects. 
When non-top level object is presented to `CRUDService` or `NetconfService`, the YDK creates corresponding top-level object and perform the requested operation.
In case of read/get operation the protocol returns always top-level objects. 
When specified filter is a non-top level object, the YDK traverses the response object tree and finds corresponding non-top level object.

.. _howto-logging:

Logging
-------
YDK uses common Python logging.  All modules are based on the ``ydk`` log. The below code snippet shows how to enable basic logging with the ``INFO`` level, which is useful for most `users` of YDK. Using the ``DEBUG`` level will produces a lot more detailed logs, which may be useful for `developers` working on YDK.

.. code-block:: python
    :linenos:

    import logging
    log = logging.getLogger('ydk')
    log.setLevel(logging.INFO)
    handler = logging.StreamHandler()
    log.addHandler(handler)

To see time stamps and logging levels, please see the below code snippet.

.. code-block:: python
   :linenos:

   import logging
   log = logging.getLogger('ydk')
   log.setLevel(logging.INFO)
   handler = logging.StreamHandler()
   formatter = logging.Formatter(("%(asctime)s - %(name)s - %(levelname)s - %(message)s"))
   handler.setFormatter(formatter)
   log.addHandler(handler)
