SchemaNode
==========

.. go:package:: ydk/cgopath
    :synopsis: CGo Path API SchemaNode


.. go:struct:: SchemaNode

    Class represents a Node in the SchemaTree.

    .. note::

        The Go SchemaNode is wrapper for YDK C++ SchemaNode implementation. No constructor is defined and the user could not instantiate a ``SchemaNode`` instance. However, the user could get an instance of ``SchemaNode`` through :go:struct:`DataNode<ydk/cgopath/DataNode>`:

        .. code-block:: go

            >>> from ydk.providers import NetconfServiceProvider
            >>> provider = NetconfServiceProvider('127.0.0.1', 'admin', 'admin', 830)
            >>> root_schema = provider.get_root_schema()                               # <-- root_schema is an instance of RootSchemaNode
            >>> bgp = root_schema.create('openconfig-bgp:bgp')                         # <-- bgp is an instance of DataNode
            >>> schema_node = bgp.schema()                                             # <-- schema node for bgp