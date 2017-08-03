.. _ref-rpc:

Rpc
===


.. cpp:class:: ydk::path::Rpc

    An instance of the YANG schmea rpc node.

    Instances of this class represent a YANG rpc and are modelled as Callables.

    The input data node tree is used to populate the input parameters to the rpc if any. The Callable takes as a parameter the :cpp:class:`Session<Session>` that can execute this rpc as its parameter returning a pointer to a :cpp:class:`DataNode<DataNode>` tree if output is available.

    .. cpp:function:: virtual ~Rpc()

    .. cpp:function:: virtual std::shared_ptr<DataNode> operator()(const Session& session)

        Execute/Invoke the rpc through the given Session.

        :param session: The Session.
        :return: Pointer to the :cpp:class:`DataNode<DataNode>` or ``nullptr`` if none exists.

    .. cpp:function:: virtual DataNode& get_input_node() const

        Get the input data tree.

        :return: Pointer to the input :cpp:class:`DataNode<DataNode>` or ``nullptr`` if the rpc does not have an input element in the schema.

    .. cpp:function:: virtual has_output_node() const = 0

    .. cpp:function:: virtual SchemaNode& get_schema_node() const

        :return: Pointer to the :cpp:class:`SchemaNode<SchemaNode>` associated with this rpc.
