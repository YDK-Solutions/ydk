NetconfServiceProvider
======================

.. cpp:namespace:: ydk

.. cpp:class:: NetconfServiceProvider : public ydk::ServiceProvider

Implementation of :cpp:class:`ServiceProvider<ydk::ServiceProvider>` for the `netconf <https://tools.ietf.org/html/rfc6241>`_ protocol.

    .. cpp:function:: NetconfServiceProvider(
        std::string address,\
        std::string username,\
        std::string password, int port = 830)

        Constructs an instance of :cpp:class::`session::<path::NetconfSession>` to connect to a server which **has** to support model download

        :param address: IP address of the device supporting a netconf interface
        :param username: Username to log in to the device
        :param password: Password to log in to the device
        :param port: Device port used to access the netconf interface. Default value is 830

    .. cpp:function:: NetconfServiceProvider(
        const path::Repository& repo,\
        std::string address,\
        std::string username,\
        std::string password, int port = 830)

        Constructs an instance of :cpp:class::`session::<path::NetconfSession>` using the provided :cpp:class:`repository<path::Repository>`

        :param repository: Reference to an instance of :cpp:class:`path::Repository<ydk::path::Repository>`
        :param address: IP address of the device supporting a netconf interface
        :param username: Username to log in to the device
        :param password: Password to log in to the device
        :param port: Device port used to access the netconf interface. Default value is 830

    .. cpp:function:: EncodingFormat get_encoding()

        Returns the type of encoding supported by the service provider. In the case of netconf service provider, :cpp:enum:`EncodingFormat::XML<EncodingFormat>` is returned.

    .. cpp:function:: NetconfSession get_session()

        Returns the instance of the :cpp:class:`NetconfSession<path::NetconfSession>` used to connect to the netconf server

    .. cpp:function:: ~NetconfServiceProvider()
