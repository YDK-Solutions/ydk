Introduction
=============
.. contents:: Table of Contents

YDK consists of two main components: core library, which consists of services and providers, and Go model API, which are APIs generated based on YANG models and packaged as packages.

Core library consists of the below:

    * **Service:** Provides simple API interface to be used with the bindings and providers
    * **ServiceProvider:** Provides concrete implementation that abstracts underlying protocol details (e.g. :go:struct:`NetconfServiceProvider<ydk/providers/NetconfServiceProvider>`, which is based on the NETCONF protocol)

Applications can be written using the Go model API in conjunction with a service and a provider.

Writing an app
--------------

In this example, we set some BGP configuration using the OpenConfig model, the CRUD (Create/Read/Update/Delete) service and the NETCONF service provider. The example in this document is a simplified version of the more complete sample that is available in ``samples/bgp.go``. The more complete sample can be run with the below steps:

.. code-block:: sh
    
    (ydk-go)ydk-go$ cd core/samples

What happens underneath
-----------------------
YDK performs the below actions when running this application:

1. Establish a session with the device
2. Encode go data objects to the protocol format (e.g. netconf XML payload)
3. Perform transport operation with the device and collect the response (e.g. netconf reply)
4. Decode response as go object and return the result to app
5. Raise go exceptions for any errors that occurred


Service Providers
-----------------
The first step in any application is to create a service provider instance. In this case, the NETCONF service provider (defined in :go:struct:`NetconfServiceProvider<ydk/providers/NetconfServiceProvider>`) is responsible for mapping between the CRUD service API and the underlying manageability protocol (NETCONF RPCs).

We instantiate an instance of the service provider that creates a NETCONF session to the machine with address 10.0.0.1:

.. code-block:: go
    :linenos:

    todo

Using the model APIs
--------------------
After establishing the connection, we instantiate the entities and set some data.

.. code-block:: go
    :linenos:
    :lineno-start: 8

    todo

Next, create a :go:struct:`Bgp<ydk/models/openconfig/openconfig_bgp/Bgp>` configuration object and set the attributes:

.. code-block:: go
    :linenos:
    :lineno-start: 10

    # create BGP object

    # set the Global AS

    # Create an AFI SAFI config

    # Add the AFI SAFI config to the global AFI SAFI list

Invoking the CRUD Service
-------------------------
The CRUD service provides methods to create, read, update and delete entities on a device making use of the session provided by a service provider (NETCONF in this case).  In order to use the CRUD service, we need to import the :go:struct:`CRUDService<ydk/services/CRUDService>` struct:

.. code-block:: go
    :linenos:
    :lineno-start: 24

    todo

Next, we instantiate the CRUD service:

.. code-block:: go
    :linenos:
    :lineno-start: 25

    todo

Finally, we invoke the create method of the in this case).  In order to use the CRUD service, we need to import the :go:struct:`CRUDService<ydk/services/CRUDService>` struct passing in the
service provider instance and our entity (``bgp_cfg``):

.. code-block:: go
    :linenos:
    :lineno-start: 26

    todo

Note if there were any errors the above API will raise an exception.

.. _howto-logging:

Logging
-------
todo

.. code-block:: go
    :linenos:

    todo



