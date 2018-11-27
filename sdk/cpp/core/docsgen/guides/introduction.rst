Introduction
============
.. contents:: Table of Contents

YDK consists of two main components: core library, which consists of services and providers, and C++ model API, which are APIs generated based on YANG models and packaged as bundles.

Core library consists of the below:

 * **Service:** Provides simple API interface to be used with the bindings and providers
 * **ServiceProvider:** Provides concrete implementation that abstracts underlying protocol details (e.g. :cpp:class:`NetconfServiceProvider<ydk::NetconfServiceProvider>`, which is based on the NETCONF protocol)

Applications can be written using the C++ model API in conjunction with a service and a provider.

Writing an app
---------------

In this example, we set some BGP configuration using the OpenConfig model, the CRUD (Create/Read/Update/Delete) service and the NETCONF service provider. The example in this document is a simplified version of the more complete sample that is available in ``core/samples/bgp_create.cpp``. Assuming you have performed the core and bundle installations first, that more complete sample can be run with the below steps:

.. code-block:: sh

  ydk-cpp$ cd core/samples
  samples$ mkdir build && cd build
  build$ cmake .. && make
  build$ ./bgp_create ssh://<username>:<password>@<host-address>:<port> [-v]

What happens underneath
~~~~~~~~~~~~~~~~~~~~~~~~
YDK performs the below actions when running this application:

 1. Establish a session with the device and, optionally, fetch any data models which are advertised by the device
 2. Encode C++ data objects to the protocol format (e.g. netconf XML payload)
 3. Perform transport operation with the device and collect the response (e.g. netconf reply)
 4. Decode response as C++ object and return the result to app
 5. Raise C++ exceptions for any errors that occurred

Header includes
---------------
In our example YDK application, first, let us include the necessary header files

.. code-block:: c++
 :linenos:

 #include <iostream>
 #include <spdlog/spdlog.h>

 #include "ydk/crud_service.hpp"
 #include "ydk/netconf_provider.hpp"

 #include "ydk_openconfig/openconfig_bgp.hpp"

 using namespace std;
 using namespace ydk;

Service Providers
----------------------
The first step in any application is to create a service provider instance. In this case, the NETCONF service provider (defined in ``ydk/netconf_provider.hpp``) is responsible for mapping between the CRUD service API and the underlying manageability protocol (NETCONF RPCs).

We instantiate an instance of the service provider that creates a NETCONF session to the machine with address 10.0.0.1

.. code-block:: c++

 NetconfServiceProvider provider{"10.0.0.1", "test", "test", 830};


Using the model APIs
----------------------
After establishing the connection, we instantiate the entities and set some data. Now, create an :cpp:class:`openconfig BGP<ydk::openconfig_bgp::Bgp>` configuration object and set the attributes

.. code-block:: c++
 :linenos:

 // Create BGP object
 auto bgp = make_shared<openconfig_bgp::Bgp>();

 // Set the Global AS
 bgp->global->config->as = 65001;
 bgp->global->config->router_id = "1.2.3.4";

 // Create a neighbor
 auto neighbor = make_shared<openconfig_bgp::Bgp::Neighbors::Neighbor>();
 neighbor->neighbor_address = "6.7.8.9";
 neighbor->config->neighbor_address = "6.7.8.9";
 neighbor->config->peer_as = 65001;
 neighbor->config->local_as = 65001;
 neighbor->config->peer_group = "IBGP";

 // Add the neighbor config to the BGP neighbors list
 bgp->neighbors->neighbor.append(neighbor);


Invoking the CRUD Service
---------------------------
The CRUD service provides methods to create, read, update and delete entities on a device making use of the session provided by a service provider (NETCONF in this case).  In order to use the CRUD service, we need to instantiate the :cpp:class:`CrudService<ydk::CrudService>` class

.. code-block:: c++

 CrudService crud_service{};

Finally, we invoke the create method of the :cpp:class:`CrudService<ydk::CrudService>` class passing in the service provider instance and our entity, ``bgp``

.. code-block:: c++
 :linenos:

 try
 {
   crud_service.create(provider, *bgp);
 }
 catch(YError & e)
 {
   cerr << "Error details: " << e.what() << endl;
 }

Note if there were any errors the above API will raise an exception with the base type :cpp:class:`YError<ydk::YError>`

Using non-top level objects
---------------------------
In the example above you noticed that we started building model from top-level object `openconfig_bgp::Bgp` and build the object tree down the hierarchy. 
However in certain conditions we can build independently non-top level objects and still be able to do all CRUD and Netconf operations.

Top level object vs. non-top
~~~~~~~~~~~~~~~~~~~~~~~~~~~~
The top level object represents top-level container in the Yang model. Examples of top-level objects:

 * openconfig_bgp::Bgp
 * openconfig_interfaces::Interfaces

The non-top level object represents a container in the Yang model, which is located under top level container. A member of a non-top level list can also be considered as non-top level object.
Examples of non-top level objects:

 * openconfig_bgp::Bgp::Neighbors
 * openconfig_bgp::Bgp::Neighbors::Neighbor
 * openconfig_bgp::Bgp::Neighbors::Neighbor::Config
 * openconfig_Interfaces::Interfaces
 * openconfig_Interfaces::Interfaces::Interface

How to use non-top level objects
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You should be able to work with non-top level objects similarly as with top level. Your program will look more simple and straight to the point.

.. code-block:: c++
 :linenos:

 // Create BGP neighbor
 auto neighbor = make_unique<openconfig_bgp::Bgp::Neighbors::Neighbor>();
 neighbor->neighbor_address = "6.7.8.9";    // This is the list object key
 neighbor->config->neighbor_address = "6.7.8.9";
 neighbor->config->peer_as = 65001;
 neighbor->config->local_as = 65001;
 neighbor->config->peer_group = "IBGP";

 crud_service.create(provider, *neighbor);
 
 // Read all configuration of specific BGP neighbor
 auto bgp = make_unique<openconfig_bgp::Bgp>();
 auto filter = make_shared<openconfig_bgp::Bgp::Neighbors::Neighbor>();
 filter->neighbor_address = "6.7.8.9";    // This is the list object key
 bgp->neighbors->neighbor.append(filter);
 
 // Get the neighbor configuration
 auto bgp_entity = crud_service.read_config(provider, *bgp);
 auto bgp_object = dynamic_cast<openconfig_bgp::Bgp*> (bgp_entity.get());
 auto neighbor_entity = bgp_object->neighbors->neighbor.get("6.7.8.9");
 auto neighbor_object = dynamic_cast<openconfig_bgp::Bgp::Neighbors::Neighbor*> (neighbor_entity.get());

Limitations
~~~~~~~~~~~

Not all non-top level objects can be used independently. Here is the rule:

  When building non-top level object, we have to define all the list keys on the way up to the top level object. 
  In the example above the object `Neighbor` is a member of the list. We can use it as long as its key `neighbor_address` is defined. 
  Other words - the absolute Yang model path of the object must be non-ambiguous. 
  In the example above the absolute path would be: `/openconfig-bgp/bgp/neighbors/neighbor[neighbor-address='6.7.8.9']`
  
Current YDK implementation in C++ does not allow to instantiate Entity class object based on its name. 
Therefore all read/get operations still require to specify top level filter object in read and get operations (see example above).
This feature is to be developed in future YDK releases.


Logging
-------

YDK uses the `spdlog` logging library. The logging can be enabled as follows by creating a logger called "ydk". For other options like logging the "ydk" log to a file, see the `spdlog reference <https://github.com/gabime/spdlog#usage-example>`_.

.. code-block:: c++
 :linenos:

 if(verbose)
 {
   auto console = spdlog::stdout_color_mt("ydk");
 }

