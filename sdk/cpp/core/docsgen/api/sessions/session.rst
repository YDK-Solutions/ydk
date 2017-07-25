Session
=======

.. cpp:namespace:: ydk::path

.. cpp:class:: Session

Class represent the base class :cpp:class:`Session<Session>` to be inherited for implementation.

    .. cpp:function:: virtual RootSchemaNode* get_root_schema() const

        Returns the :cpp:class:`RootSchemaNode<RootSchemaNode>` tree supported by this instance of the ``Session``.

        :return: Pointer to the :cpp:class:`RootSchemaNode<RootSchemaNode>` or ``nullptr`` if one could not be created.

    .. cpp:function:: virtual path::DataNode* invoke(path::Rpc* rpc) const

        Invokes or executes the given rpc and returns a :cpp:class:`DataNode<DataNode>` pointer if the Rpc has an output modeled in YANG.

        :param rpc: Pointer to the :cpp:class:`Rpc<Rpc>` node.
        :return: Pointer to the :cpp:class:`DataNode<DataNode>` representing the output.

    .. cpp:function:: ~Session()
