Using OpenDaylight with YDK
============================
.. contents::

YDK makes it easy to interact with OpenDaylight programmatically using the YANG model APIs.

Applications can be written using the Go model API in conjunction with a service and a provider.

Writing the app
---------------

In this example, we set some BGP configuration using the Cisco IOS XR model, the :go:struct:`CRUD (Create/Read/Update/Delete) service<ydk/services/CRUDService>` and the :go:struct:`OpenDaylightServiceProvider<ydk/providers/OpendaylightServiceProvider>`. The example in this document is a simplified version of the more complete sample that is available in ``core/samples/bgp_xr_opendaylight.go``. Assuming you have performed the ``core`` and ``cisco-ios-xr`` bundle installations first(see :ref:`howto-install`), more complete sample can be run with the below steps:

.. code-block:: bash
    
    todo

What happens underneath
-----------------------
YDK performs the below actions when running this application:

 1. Establish a session with the OpenDaylight instance and fetch the details of the nodes mounted on the southbound
 2. Encode Go data objects to the protocol format (e.g. restconf JSON payload)
 3. For a chosen node on the southbound, perform transport operation with the device and collect the response (e.g. restconf reply)
 4. Decode response as Go object and return the result to app
 5. Raise Go exceptions for any errors that occurred

Import libraries
----------------
In our example YDK application, first, let us import the necessary libraries

.. code-block:: go
    :linenos:

    todo

OpenDaylight service provider
-----------------------------
The first step in any application is to create a service provider instance. In this case, the OpenDaylight service provider is responsible for mapping between the CRUD service API and the underlying manageability protocol (Restconf).

We first instantiate a :go:struct:`Repository<ydk/cgopath/Repository>` using the location of the schema cache of the OpenDaylight instance. We instantiate an instance of the service provider that can communicate using Restconf with an OpenDaylight instance running at host address: ``127.0.0.1`` and port: ``8181``

.. code-block:: go
    :linenos:
    :lineno-start: 17

    todo


Using the model APIs
--------------------
After establishing the connection, let's instantiate the entities and set some data. Now, create a Cisco IOS XR :go:struct:`Bgp<ydk/models/cisco_ios_xr/Cisco_IOS_XR_ipv4_bgp_cfg/Bgp>` configuration object and set the attributes

.. code-block:: go
    :linenos:
    :lineno-start: 19

    # Create BGP object
    # BGP instance
    # global address family
    # add the instance to the parent BGP object


Invoking the CRUD Service
-------------------------
The CRUD service provides methods to create, read, update and delete entities on a device making use of the session provided by a service provider.  In order to use the CRUD service, we need to instantiate the :go:struct:`CRUDService<ydk/services/CRUDService>` class

.. code-block:: go
    :linenos:
    :lineno-start: 41

    todo

At this point we can explore the southbound device node-IDs using the function call: :go:func:`OpenDaylightServiceProviderGetNodeIDByIndex<OpenDaylightServiceProviderGetNodeIDByIndex>`. Let us assume there is a XR device mounted with the node ID ``xr``. We can obtain the :go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>` instance corresponding to this node using the : :go:func:`OpenDaylightServiceProviderGetNodeProvider<OpenDaylightServiceProviderGetNodeProvider>`.

Finally, we invoke the create method of the :go:struct:`CRUDService<ydk/services/CRUDService>` class passing in the service provider instance and our entity, ``bgp``

.. code-block:: go
    :linenos:
    :lineno-start: 42

    todo


Note if there were any errors the above API will raise an exception.

Logging
-------
todo
