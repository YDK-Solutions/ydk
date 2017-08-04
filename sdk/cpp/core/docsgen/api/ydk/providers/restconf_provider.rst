RestconfServiceProvider
========================


.. cpp:class:: ydk::RestconfServiceProvider : public ServiceProvider

    Implementation of :cpp:class:`ServiceProvider<ServiceProvider>` for the `restconf <https://tools.ietf.org/html/draft-ietf-netconf-restconf-18>`_ protocol.

    .. cpp:function::RestconfServiceProvider(\
        path::Repository & repo, \
        const std::string & address, \
        const std::string & username, \
        const std::string & password, \
        int port = 80, \
        EncodingFormat encoding = EncodingFormat::JSON, \
        const std::string & config_url_root = "/data", \
        const std::string & state_url_root = "/data")

        Constructs an instance of the ``RestconfServiceProvider`` to connect to a restconf server

        :param repository: Reference to an instance of :cpp:class:`path::Repository<ydk::path::Repository>`
        :param address: IP address of the device supporting a restconf interface
        :param username: Username to log in to the device
        :param password: Password to log in to the device
        :param port: Device port used to access the restconf interface. Default value is 80
        :param encoding: Type of encoding to be used for the payload. Default is :cpp:enumerator:`JSON<EncodingFormat::JSON>`
        :param config_url_root: To provider backwards compatibility with older drafts of restconf RFC, this can be "/config" or "/data" (which is the default)
        :param state_url_root: To provider backwards compatibility with older drafts of restconf RFC, this can be "/operational" or "/data" (which is the default)

    .. cpp:function::RestconfServiceProvider(\
        std::unique_ptr<RestconfClient> client, \
        std::shared_ptr<ydk::path::RootSchemaNode> root_schema, \
        const std::string & edit_method, \
        const std::string & config_url_root, \
        const std::string & state_url_root, \
        EncodingFormat encoding)

    .. cpp:function:: EncodingFormat get_encoding() const

        Returns the type of encoding supported by the service provider.

    .. cpp:function:: const RestconfSession& get_session()

        Returns a reference of the :cpp:class:`RestconfSession<path::RestconfSession>` used to connect to the restconf server.

    .. cpp:function:: ~RestconfServiceProvider()
