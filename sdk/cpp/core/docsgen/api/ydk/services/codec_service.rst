Codec Service
=============

.. cpp:class:: ydk::CodecService

    Codec Service class for supporting encoding and decoding C++ model API objects of type :cpp:class:`Entity<Entity>`.

    .. cpp:function:: CodecService()

        Constructs an instance of CodecService

    .. cpp:function:: std::string encode(CodecServiceProvider & provider, Entity & entity, bool pretty=false, bool subtree=false)

        Perform encoding.

        :param provider: An instance of :cpp:class:`CodecServiceProvider<ydk::CodecServiceProvider>`.
        :param entity: An instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :param pretty: Optionally produce formatted output.
        :param subtree: Subtree filter.
        :return: Encoded payload.
        :raises YError: If an error has occurred

    .. cpp:function:: std::map<std::string, std::string> encode(CodecServiceProvider & provider, std::map<std::string, std::shared_ptr<Entity>> & entity, bool pretty=false)

        Perform encoding

        :param provider: An instance of :cpp:class:`CodecServiceProvider<ydk::CodecServiceProvider>`
        :param entity: A map of `Entity` class defined under same bundle
        :param pretty: Optionally produce formatted output
        :return: A map of encodec payload.
        :raises YError: If an error has occurred

    .. cpp:function:: std::shared_ptr<ydk::Entity> decode(CodecServiceProvider & provider, const std::string & payload, bool subtree=false)

        Decode the payload to produce an instance of `Entity`.

        :param provider: An instance of :cpp:class:`CodecServiceProvider<ydk::CodecServiceProvider>`.
        :param payload: Payload to be decoded.
        :param subtree: Subtree filter.
        :return: Pointer to the decoded `Entity`.
        :raises YError: If an error has occurred

    .. cpp:function:: std::map<std::string, std::shared_ptr<Entity>> decode(CodecServiceProvider & provider, std::map<std::string, std::string> & payload_map, std::map<std::string, std::shared_ptr<Entity>> entity_map)

        Decode map of payload to map of `Entity`.

        :param provider: An instance of :cpp:class:`CodecServiceProvider<ydk::CodecServiceProvider>`.
        :param payload_map: Module name payload map.
        :param entity_map: Module name entity map.
        :return: A ``std::map`` of the decoded `Entity`.
        :raises YError: If an error has occurred.

XmlSubtreeCodec
===============

.. cpp:class:: ydk::XmlSubtreeCodec

    XmlSubtreeCodec class designed to provide encoding and decoding C++ model API objects of type :cpp:class:`Entity<Entity>` to/from XML encoded string.
    Compared to :cpp:class:`CodecService<ydk::CodecService>` the class does not validate encoded data for their types and values.

    .. cpp:function:: XmlSubtreeCodec()

        Constructs an instance of `XmlSubtreeCodec` class.

    .. cpp:function:: std::string encode(Entity & entity, path::RootSchemaNode & root_schema)

        Performs encoding of C++ model API objects of type :cpp:class:`Entity<Entity>` to well formatted XML encoded string.

        :param entity: An instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :param root_schema: An instance of :cpp:class:`RootSchemaNode<path::RootSchemaNode>`, which includes the model bundle.
        :return: A `std::string`. Encoded well formatted multi-line XML payload.
        :raises YServiceProviderError: If an error has occurred; usually appears when model is not present in the bundle.

    .. cpp:function:: std::shared_ptr<ydk::Entity> decode(const std::string & payload, std::shared_ptr<Entity> entity)

        Decodes the XML encoded string to produce corresponding instance of :cpp:class:`Entity<Entity>`.

        :param payload: `std::string`, XML encoded string to be decoded.
        :param entity: `std::shared_ptr<Entity>`, instance of shared pointer to expected top level `Entity` class.
        :return: `std::shared_ptr<Entity>`, shared pointer to the decoded `Entity`.
        :raises YServiceProviderError: If an error has occurred; usually appears when payload does not correspond to `Entity` model.

JsonSubtreeCodec
================

.. cpp:class:: ydk::JsonSubtreeCodec

    JsonSubtreeCodec class designed to provide encoding and decoding C++ model API objects of type :cpp:class:`Entity<Entity>` to/from JSON encoded string.
    Compared to :cpp:class:`CodecService<ydk::CodecService>` the class does not validate encoded data for their types and values.

    .. cpp:function:: JsonSubtreeCodec()

        Constructs an instance of `JsonSubtreeCodec` class.

    .. cpp:function:: std::string encode(Entity & entity, path::RootSchemaNode & root_schema, bool pretty=true)

        Performs encoding of C++ model API objects of type :cpp:class:`Entity<Entity>` to JSON encoded string.

        :param entity: An instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :param root_schema: An instance of :cpp:class:`RootSchemaNode<path::RootSchemaNode>`, which includes the model bundle.
        :param pretty: `bool`. If set to `true`, the function produces well formatted multi-line JSON string. If set to `false` - one line string.
        :return: A `std::string`. Encoded JSON payload.
        :raises YServiceProviderError: If an error has occurred; usually appears when model is not present in the bundle.

    .. cpp:function:: std::shared_ptr<ydk::Entity> decode(const std::string & payload, std::shared_ptr<Entity> entity)

        Decodes the JSON encoded string to produce corresponding instance of :cpp:class:`Entity<Entity>`.

        :param payload: `std::string`, JSON encoded string to be decoded.
        :param entity: `std::shared_ptr<Entity>`, instance of shared pointer to expected top level `Entity` class.
        :return: `std::shared_ptr<Entity>`, shared pointer to the decoded `Entity`.
        :raises YServiceProviderError: If an error has occurred; usually appears when payload does not correspond to `Entity` model.

Example of JsonSubtreeCodec usage
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

In this example we use :cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>` and
:cpp:class:`CRUDService<ydk::CRUDService>` to get interface configuration from IOS XR device and
then print it using :cpp:class:`JsonSubtreeCodec<ydk::JsonSubtreeCodec>`::

    #include <iostream>

    #include <ydk/crud_service.hpp>
    #include <ydk/gnmi_provider.hpp>
    #include <ydk/json_subtree_codec.hpp>

    #include <ydk_cisco_ios_xr/Cisco_IOS_XR_ifmgr_cfg.hpp>

    using namespace ydk;
    using namespace std;
    using namespace cisco_ios_xr;

    int main()
    {
        auto repo = path::Repository("/home/yan/ydk-workspace/ydk-gen/scripts/repository/10.30.110.84");
        gNMIServiceProvider provider{repo, "10.30.110.84", 57400, "admin", "admin"};
        CrudService crud{};

        // Build filter to retrieve interface configuration
        auto ifcs_config = Cisco_IOS_XR_ifmgr_cfg::InterfaceConfigurations();
        auto ifc_config = make_shared<Cisco_IOS_XR_ifmgr_cfg::InterfaceConfigurations::InterfaceConfiguration>();
        ifc_config->active = "\"act\"";
        ifc_config->interface_name = "\"Loopback0\"";
        ifcs_config.interface_configuration.append(ifc_config);

        // Read interface configuration
        auto ifc_read = crud.read(provider, ifcs_config);

        // Print interface configuration
        if (ifc_read) {
            JsonSubtreeCodec jcodec{};
            auto json_payload = jcodec.encode(*ifc_read, provider.get_session().get_root_schema(), true);
            cout << "Interface Configuration:" << endl << json_payload << endl;
        }
    }
