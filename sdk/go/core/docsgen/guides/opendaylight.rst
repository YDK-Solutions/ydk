How do I use OpenDaylight with YDK?
===================================
.. contents::

YDK makes it easy to interact with OpenDaylight programmatically using the YANG model APIs.

Applications can be written using the Go model API in conjunction with a service and a provider.

Writing the app
---------------

In this example, we set some BGP configuration using the Cisco IOS XR model, the :go:struct:`CRUD (Create/Read/Update/Delete) service<ydk/services/CrudService>` and the :go:struct:`OpenDaylightServiceProvider<ydk/providers/OpendaylightServiceProvider>`. The example in this document is a simplified version of the more complete sample that is available in ``core/samples/bgp_xr_opendaylight.go``. Assuming you have performed the ``core`` and ``cisco-ios-xr`` bundle installations first(see :ref:`howto-install`), more complete sample can be run with the below steps:

.. code-block:: bash
    
    $ cd core/samples
    $ go run bgp_xr_opendaylight.py -device http://<username>:<password>@<host-address>:<port>

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

.. code-block:: c
    :linenos:

    package main

    import (
        xr_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/openconfig/ipv4_bgp_cfg"
        xr_bgp_types "github.com/CiscoDevNet/ydk-go/ydk/models/openconfig/ipv4_bgp_datatypes"
        "github.com/CiscoDevNet/ydk-go/ydk"
        "github.com/CiscoDevNet/ydk-go/ydk/providers"
        "github.com/CiscoDevNet/ydk-go/ydk/services"
        "github.com/CiscoDevNet/ydk-go/ydk/types"
        "github.com/CiscoDevNet/ydk-go/ydk/types/protocol"
        encoding "github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format"
    )

OpenDaylight service provider
-----------------------------
The first step in any application is to create a service provider instance. In this case, the OpenDaylight service provider is responsible for mapping between the CRUD service API and the underlying manageability protocol (Restconf).

We instantiate an instance of the service provider that can communicate using Restconf with an OpenDaylight instance running at host address: ``127.0.0.1`` and port: ``8181``

.. code-block:: c
    :linenos:
    :lineno-start: 13

    func main() {
        provider := providers.OpenDaylightServiceProvider{
            // In this case, we have a ODL boron instance with this schema cache location
            Path:           "/Users/home/distribution-karaf-0.5.2-Boron-SR2/cache/schema",
            Address:        "127.0.0.1",
            Username:       "admin",
            Password:       "admin",
            Port:           8181,
            EncodingFormat: encoding.XML,
            Protocol:       protocol.Restconf}
        suite.Provider.Connect()


Using the model APIs
--------------------
After establishing the connection, let's instantiate the entities and set some data. Now, create a Cisco IOS XR :go:struct:`Bgp<ydk/models/cisco_ios_xr/ipv4_bgp_cfg/Bgp>` configuration object and set the attributes

.. code-block:: c
    :linenos:
    :lineno-start: 24

        // Create BGP object
        bgp := xr_bgp.Bgp{}

        // BGP instance
        bgp.Instance = xr_bgp.Bgp_Instance{}
        bgp.Instance.InstanceName = "test"

        instanceAs := xr_bgp.Bgp_Instance_InstanceAs{}
        instanceAs.As = 65001
        
        fourByteAs := xr_bgp.Bgp_Instance_As_FourByteAs{}
        fourByteAs.As = 65001
        fourByteAs.BgpRunning = types.Empty{}

        // global address family
        fourByteAs.DefaultVrf = xr_bgp.Bgp_Instance_InstanceAs_FourByteAs_DefaultVrf{}
        fourByteAs.DefaultVrf.Global = xr_bgp.Bgp_Instance_InstanceAs_FourByteAs_DefaultVrf_Global{}
        fourByteAs.DefaultVrf.Global.GlobalAfs = xr_bgp.Bgp_Instance_InstanceAs_FourByteAs_DefaultVrf_Global_GlobalAfs{}

        globalAf := xr_bgp.Bgp_Instance_InstanceAs_FourByteAs_DefaultVrf_Global_GlobalAfs_GlobalAf{}
        globalAf.AfName = xr_bgp_types.BgpAddressFamily_ipv4_unicast
        globalAf.Enable = types.Empty{}

        fourByteAs.DefaultVrf.Global.GlobalAfs.GlobalAf = append(fourByteAs.DefaultVrf.Global.GlobalAfs.GlobalAf, &globalAf)

        // add the instance to the parent bgp object
        instanceAs.FourByteAs = append(instanceAs.FourByteAs, &fourByteAs)
        bgp.Instance.InstanceAs = append(bgp.Instance.InstanceAs, &instanceAs)


Invoking the CRUD Service
-------------------------
The CRUD service provides methods to create, read, update and delete entities on a device making use of the session provided by a service provider.  In order to use the CRUD service, we need to instantiate the :go:struct:`CrudService<ydk/services/CrudService>` class

.. code-block:: go
    :linenos:
    :lineno-start: 55

        crud := services.CrudService{}

At this point we can explore the southbound device node-IDs using the function call: :go:func:`OpenDaylightServiceProviderGetNodeIDByIndex<OpenDaylightServiceProviderGetNodeIDByIndex>`. Let us assume there is a XR device mounted with the node ID ``xr``. We can obtain the :go:struct:`ServiceProvider<ydk/types/ServiceProvider>` instance corresponding to this node using the : :go:func:`OpenDaylightServiceProviderGetNodeProvider<OpenDaylightServiceProviderGetNodeProvider>`.

Finally, we invoke the create method of the :go:struct:`CrudService<ydk/services/CrudService>` class passing in the service provider instance and our entity, ``bgp``

.. code-block:: go
    :linenos:
    :lineno-start: 56

        p := provider.GetNodeProvider("xr")
        crud.Create(p, &bgp)
    }


Note if there were any errors the above API will raise an exception.

Logging
-------
See :ref:`howto-logging`.
