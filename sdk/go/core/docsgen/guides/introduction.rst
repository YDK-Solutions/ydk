Introduction
=============
.. contents:: Table of Contents

YDK consists of two main components: the core library, which consists of services and providers, and Go model API, which are APIs generated based on YANG models and packaged as packages.

The core library consists of the below:

    * **Service:** Provides simple API interface to be used with the bindings and providers
    * **ServiceProvider:** Provides concrete implementation that abstracts underlying protocol details (e.g. :go:struct:`NetconfServiceProvider<ydk/providers/NetconfServiceProvider>`, which is based on the NETCONF protocol)

Applications can be written using the Go model API in conjunction with a service and a provider.

Writing an app
--------------

In this example, we set some BGP configuration using the OpenConfig model, the CRUD (Create/Read/Update/Delete) service and the NETCONF service provider. The example in this document is a simplified version of the more complete sample that is available in ``samples/bgp.go``. The more complete sample can be run with the below steps:

.. code-block:: sh
    
    $ cd core/samples/bgp_create
    $ go run bgp_create.go


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
The first step in any application is to create a service provider instance. In this case, the NETCONF service provider (defined in :go:struct:`NetconfServiceProvider<ydk/providers/NetconfServiceProvider>`) is responsible for mapping between the CRUD service API and the underlying manageability protocol (NETCONF RPCs). To access the providers, we must include the following import statement at the top of our app.

.. code-block:: c

    import "github.com/CiscoDevNet/ydk-go/ydk/providers"

We instantiate an instance of the service provider to the machine with address 10.0.0.1:

.. code-block:: c
    :linenos:

    func main() {
        var provider providers.NetconfServiceProvider = providers.NetconfServiceProvider{
            Address:  "127.0.0.1",
            Username: "admin",
            Password: "admin",
            Port:     12022}
        provider.Connect()


Using the model APIs
--------------------
After establishing the connection, we instantiate the entities and set some data. In this example, we'll be using the OpenConfig BGP package. To use this package, we must include the following import statments at the top of our app.

.. code-block:: c

    import oc_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp"
    import oc_bgp_types "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp_types"

In our main function, we'll create a :go:struct:`Bgp<ydk/models/openconfig/openconfig_bgp/Bgp>` configuration object and set the attributes:

.. code-block:: c
    :linenos:
    :lineno-start: 8

        // create BGP object
        bgp := oc_bgp.Bgp{}

        // set the Global AS
        bgp.Global.Config.As = 65172

        // Create an AFI SAFI config
        ipv6_afisafi := oc_bgp.Bgp_Global_AfiSafis_AfiSafi{}
        ipv6_afisafi.AfiSafiName = &oc_bgp_types.IPV6UNICAST{}
        ipv6_afisafi.Config.AfiSafiName = &oc_bgp_types.IPV6UNICAST{}
        ipv6_afisafi.Config.Enabled = true

        // Add the AFI SAFI config to the global AFI SAFI list
        bgp.Global.AfiSafis.AfiSafi = append(bgp.Global.AfiSafis.AfiSafi, ipv6_afisafi)

Invoking the CRUD Service
-------------------------

The CRUD service provides methods to create, read, update and delete entities on a device making use of the session provided by a service provider (NETCONF in this case). To use the :go:struct:`CrudService<ydk/services/CrudService>` service, we need to include the import statment:

.. code-block:: c

    import "github.com/CiscoDevNet/ydk-go/ydk/services"

In our main function, we instantiate the CRUD service:

.. code-block:: c
    :linenos:
    :lineno-start: 24

        crud := services.CrudService{}

Finally, we invoke the create method of the in this case).  In order to use the CRUD service, we need to import the :go:struct:`CrudService<ydk/services/Crudervice>` struct passing in the
service provider instance and our entity (``bgp_cfg``):

.. code-block:: c
    :linenos:
    :lineno-start: 25

        crud.Create(&provider, &bgp)
    }

Note if there were any errors the above API will raise an exception.

Using non-top level objects
---------------------------

In the example above you noticed that we started building model from top-level object - :go:struct:`Bgp<ydk/models/openconfig/openconfig_bgp/Bgp>` and then built the object tree down the hierarchy. 
However in certain conditions we can build independently non-top level objects and still be able to do all the CRUD operations.

Top level object vs. non-top
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
The top level object represents top-level container in the Yang model. Examples of top-level objects:

 * oc_bgp.Bgp
 * oc_int.Interfaces

The non-top level object represents a container in the Yang model, which is located under top level container. A member of a non-top level list can also be considered as non-top level object.
Examples of non-top level objects:

 * oc_bgp.Bgp_Global_AfiSafis_AfiSafi
 * oc_bgp.Bgp_Neighbors
 * oc_bgp.Bgp_Neighbors_Neighbor
 * oc_bgp.Bgp_Neighbors_Neighbor_Config
 * oc_int.Interfaces_Interface

How to use non-top level objects
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You should be able to work with non-top level objects similarly as with top level. 
Your program will look more simple and straight to the point.
The above example will look now like this:

.. code-block:: c
 :linenos:

    import oc_bgp "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp"
    import oc_bgp_types "github.com/CiscoDevNet/ydk-go/ydk/models/ydktest/openconfig_bgp_types"
    import "github.com/CiscoDevNet/ydk-go/ydk/services"
    import "github.com/CiscoDevNet/ydk-go/ydk/providers"
 
    func createBgpAfisafiConfig(provider NetconfServiceProvider) {
        // Create single AFI SAFI configuration
        ipv6_afisafi := oc_bgp.Bgp_Global_AfiSafis_AfiSafi{}
        ipv6_afisafi.AfiSafiName = &oc_bgp_types.IPV6UNICAST{}
        ipv6_afisafi.Config.AfiSafiName = &oc_bgp_types.IPV6UNICAST{}
        ipv6_afisafi.Config.Enabled = true

        crud := services.CrudService{}
        crud.Create(&provider, &ipv6_afisafi)
    }
    
    func readBgpAfisafiConfig(provider NetconfServiceProvider) {
        // Read single AFI SAFI configuration
        afisafiFilter := oc_bgp.Bgp_Global_AfiSafis_AfiSafi{}
        afisafiFilter.AfiSafiName = &oc_bgp_types.IPV6UNICAST{}
        crud := services.CrudService{}
        afisafiEntity := crud.ReadConfig(&provider, &afisafiFilter)
        afisafi := afisafiEntity.(*oc_bgp.Bgp_Global_AfiSafis_AfiSafi)
    }

Limitations
~~~~~~~~~~~

Not all non-top level objects can be used independently. Here is the rule:

  When building non-top level object, we have to define all the list keys on the way up to the top level object. 
  In the example above the object `ipv6_afisafi` is a member of the list. We can use it as long as its key `AfiSafiName` is defined. 
  
Under the hood
~~~~~~~~~~~~~~

The programmability protocols like Netconf, gNMI, etc. are always working with top level model objects. 
When non-top level object is presented to `CrudService` or `NetconfService`, the YDK creates corresponding top-level object and perform the requested operation.
In case of read/get operation the protocol returns always top-level objects. 
When specified filter is a non-top level object, the YDK traverses the response object tree and finds corresponding non-top level object.

.. _howto-logging:

Logging
-------
Go logging relies on wrapper functions around cpp logging. All modules are based on the ydk log. The below code snippet shows how to enable basic logging with the INFO level, which is useful for most users of YDK. Using the DEBUG level will produces a lot more detailed logs, which may be useful for developers working on YDK.

.. code-block:: c
    :linenos:
    
    package main

    import "github.com/CiscoDevNet/ydk-go/ydk"

    func main() {
        ydk.EnableLogging(ydk.Info)
    }

