NetconfSession
======================

.. cpp:namespace:: ydk

.. cpp:class:: NetconfSession : public path::Session

Implementation of :cpp:class:`Session<path::Session>` for the `netconf <https://tools.ietf.org/html/rfc6241>`_ protocol.

    .. cpp:function:: NetconfSession(
        std::string address,\
        std::string username,\
        std::string password,\
        int port = 830)

        Constructs an instance of the ``NetconfSession>`` to connect to a server which **has** to support model download

        :param address: IP address of the device supporting a netconf interface
        :param username: Username to log in to the device
        :param password: Password to log in to the device
        :param port: Device port used to access the netconf interface. Default value is 830

    .. cpp:function:: NetconfSession(
        const path::Repository& repo,\
        std::string address,\
        std::string username,\
        std::string password, int port = 830)

        Constructs an instance of the ``NetconfSession>`` using the provided :cpp:class:`repository<path::Repository>`

        :param repository: Reference to an instance of :cpp:class:`path::Repository<ydk::path::Repository>`
        :param address: IP address of the device supporting a netconf interface
        :param username: Username to log in to the device
        :param password: Password to log in to the device
        :param port: Device port used to access the netconf interface. Default value is 830

    .. cpp:function:: path::RootSchemaNode& get_root_schema() const

        Returns the :cpp:class:`RootSchemaNode<path::RootSchemaNode>` tree supported by this instance of the ``NetconfSession``.

        :return: Pointer to the :cpp:class:`RootSchemaNode<path::RootSchemaNode>` or ``nullptr`` if one could not be created.

    .. cpp:function:: std::shared_ptr<path::DataNode> invoke(path::Rpc& rpc) const

        Invokes or executes the given rpc and returns a :cpp:class:`DataNode<DataNode>` pointer if the Rpc has an output modelled in YANG.

        :param rpc: Reference to the :cpp:class:`Rpc<Rpc>` node.
        :return: Shared pointer to the :cpp:class:`DataNode<DataNode>` representing the output.

    .. cpp:function:: ~NetconfSession()
