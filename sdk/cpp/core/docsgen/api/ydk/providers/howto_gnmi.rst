Using gNMI with YDK
============================
.. contents::

YDK makes it easy to interact with gNMI programmatically using the YANG model APIs.

Applications can be written using the C++ model API in conjunction with a service and a provider.

Writing the app
----------------

In this example, we set some BGP configuration using the Cisco IOS XR model, the :cpp:class:`CRUD (Create/Read/Update/Delete) service<ydk::CrudService>` and the :cpp:class:`gNMI service provider<ydk::gNMIServiceProvider>`. The example in this document is a simplified version of the more complete sample that is available in ``core/samples/bgp_gnmi_create.cpp``. Assuming you have performed the core and cisco-ios-xr bundle installations first, that more complete sample can be run with the below steps:

.. code-block:: sh

  ydk-cpp$ cd core/samples
  samples$ mkdir build && cd build
  build$ cmake .. && make
  build$ ./bgp_gnmi_create http://<username>:<password>@<host-address>:<port> [-v]

What happens underneath
~~~~~~~~~~~~~~~~~~~~~~~~
YDK performs the below actions when running this application:

 1. Establish a gRPC session with the gNMI server and fetch the details of the nodes
 2. Encode C++ data objects to the protocol format (e.g. JSON/JSON_IETF payload)
 3. For a chosen node, perform transport operation with the device and collect the response (e.g. gNMI reply)
 4. Decode response as C++ object and return the result to app
 5. Raise C++ exceptions for any errors that occurred

Header includes
----------------------
In our example YDK application, first, let us include the necessary header files

.. code-block:: c++
 :linenos:

 #include <iostream>
 #include <spdlog/spdlog.h>

 #include <ydk/crud_service.hpp>
 #include <ydk/path_api.hpp>
 #include <ydk/gnmi_provider.hpp>
 #include <ydk/types.hpp>

 #include <ydk_cisco_ios_xr/Cisco_IOS_XR_ipv4_bgp_cfg.hpp>
 #include <ydk_cisco_ios_xr/Cisco_IOS_XR_ipv4_bgp_datatypes.hpp>

  // indicate the namespaces being used (optional)
  using namespace std;
  using namespace ydk;
  using namespace Cisco_IOS_XR_ipv4_bgp_cfg;
  using namespace Cisco_IOS_XR_ipv4_bgp_datatypes;

gNMI service provider
------------------------------
The first step in any application is to create a service provider instance. In this case, the gNMI service provider is responsible for mapping between the CRUD service API and the underlying manageability protocol (gNMI).

We first instantiate a :cpp:class:`Repository<ydk::path::Repository>` using the location of the schema cache of the gNMI instance. We instantiate an instance of the service provider that can communicate using GRPC with an gNMI instance running at host address: ``127.0.0.1`` and port: ``50051``

.. code-block:: c++

 path::Repository repo{"/Users/home/distribution-karaf-0.5.2-Boron-SR2/cache/schema"}; // In this case, we have a GRPC service with this schema cache location
 gNMIServiceProvider gnmi_provider{repo, "127.0.0.1", "admin", "admin", 50051};


Using the model APIs
----------------------
After establishing the connection, we instantiate the entities and set some data. Now, create an :cpp:class:`Cisco IOS XR BGP<ydk::Cisco_IOS_XR_ipv4_bgp_cfg::Bgp>` configuration object and set the attributes

.. code-block:: c++
 :linenos:

 // Create BGP object
 auto bgp = make_unique<Bgp>();

 // BGP instance
 auto instance = make_unique<Bgp::Instance>();
 instance->instance_name = "test";
 auto instance_as = make_unique<Bgp::Instance::InstanceAs>();
 instance_as->as = 65001;
 auto four_byte_as = make_unique<Bgp::Instance::InstanceAs::FourByteAs>();
 four_byte_as->as = 65001;
 four_byte_as->bgp_running = Empty();

 // global address family
 auto global_af = make_unique<Bgp::Instance::InstanceAs::FourByteAs::DefaultVrf::Global::GlobalAfs::GlobalAf>();
 global_af->af_name = BgpAddressFamilyEnum::ipv4_unicast;
 global_af->enable = Empty();
 global_af->parent = four_byte_as->default_vrf->global->global_afs.get();
 four_byte_as->default_vrf->global->global_afs->global_af.push_back(move(global_af));

 // add the instance to the parent BGP object
 four_byte_as->parent = instance_as.get();
 instance_as->parent = instance.get();
 instance->parent = bgp.get();
 instance_as->four_byte_as.push_back(move(four_byte_as));
 instance->instance_as.push_back(move(instance_as));
 bgp->instance.push_back(move(instance));


Invoking the CRUD Service
---------------------------
The CRUD service provides methods to create, read, update and delete entities on a device making use of the session provided by a service provider.  In order to use the CRUD service, we need to instantiate the :cpp:class:`CrudService<ydk::CrudService>` class

.. code-block:: c++

 CrudService crud_service{};

We can obtain the :cpp:class:`ServiceProvider<ydk::ServiceProvider>` instance corresponding to this node using: ``gNMIServiceProvider provider{repo, address}; ``.

Finally, we invoke the create method of the :cpp:class:`CrudService<ydk::CrudService>` class passing in the service provider instance and our entity, ``bgp``

.. code-block:: c++
 :linenos:

 try
 {
   gNMIServiceProvider provider{repo, "127.0.0.1:50051"};  
   crud_service.create(provider, *bgp);
 }
 catch(YCPPError & e)
 {
   cerr << "Error details: " << e.what() << endl;
 }

Note if there were any errors the above API will raise an exception with the base type :cpp:class:`YCPPError<ydk::YCPPError>`

Logging
----------------------
YDK uses the `spdlog` logging library. The logging can be enabled as follows by creating a logger called "ydk". For other options like logging the "ydk" log to a file, see the `spdlog reference <https://github.com/gabime/spdlog#usage-example>`_.

.. code-block:: c++
 :linenos:

 if(verbose)
 {
   auto console = spdlog::stdout_color_mt("ydk");
 }
