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
