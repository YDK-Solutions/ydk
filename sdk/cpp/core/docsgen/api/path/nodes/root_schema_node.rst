RootSchemaNode
==============


Instances of this class represent the Root of the SchemaTree.

.. cpp:class:: ydk::path::RootSchemaNode : public SchemaNode

    Instances of this class represent the Root of the :cpp:class:`SchemaNode<SchemaNode>` Tree. A :cpp:class:`RootSchemaNode<RootSchemaNode>` can be used to instantiate a :cpp:class:`DataNode<DataNode>` tree or an :cpp:class:`Rpc<Rpc>` object. The children of the :cpp:class:`RootSchemaNode<RootSchemaNode>` represent the top level :cpp:class:`SchemaNode<SchemaNode>` in the YANG module submodules.

    .. cpp:function:: virtual ~RootSchemaNode()

        Destructor for the :cpp:class:`RootSchemaNode<RootSchemaNode>`

    .. cpp:function:: std::string get_path() const

        Get the path expression representing this Node in in the NodeTree.

        :return: ``std::string`` representing the path to this Node.

    .. cpp:function:: virtual std::vector<SchemaNode*> find(const std::string& path) const

        Finds descendant nodes that match the given xpath expression.

        This API finds descendant nodes in the :cpp:class:`SchemaNode<SchemaNode>` tree that satisfy the given path expression. See :ref:`how to path <ref-howtopath>`.

        :param path: The path expression.
        :return: Vector of :cpp:class:`SchemaNode<SchemaNode>`  that satisfies the criterion.
        :raises: :cpp:class:`YPathError<YPathError>` if the path expression in invalid, See error code for details.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>` if the argument is invalid.

    .. cpp:function:: virtual SchemaNode* get_parent() const noexcept

        Get the parent node of this :cpp:class:`SchemaNode<SchemaNode>` in the tree.

        :return: Pointer to the parent node or ``nullptr`` in case this is the root.

    .. cpp:function:: virtual std::vector<SchemaNode*> get_children() const

        Get the children of this :cpp:class:`SchemaNode<SchemaNode>` in the NodeTree.

        :return: The children of this node.

    .. cpp:function:: virtual const SchemaNode& get_root() const noexcept

        Get the root of NodeTree this node is part of.

        :return: Reference to the root.

    .. cpp:function:: virtual DataNode& create_datanode(const std::string& path, const std::string& value) const

        Create a :cpp:class:`DataNode<DataNode>` corresponding to the path and set its value.

        This methods creates a :cpp:class:`DataNode<DataNode>` tree based on the path passed in. The path expression must identify a single node. If the last node created is of schema type ``list``, ``leaf-list`` or ``anyxml`` that value is also set in the node.

        The returned :cpp:class:`DataNode<DataNode>` is the last node created (the terminal part of the path).

        The user is responsible for managing the memory of this returned tree. Use :cpp:func:`root` to get the root element of the this tree and use that pointer to dispose of the entire tree.

        Note in the case of List nodes the keys must be present in the path expression in the form of predicates.

        :param path: The XPath expression identifying the node relative to the root of the schema tree.
        :param value: The string representation of the value to set.
        :return: Reference to :cpp:class:`DataNode<DataNode>` created.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>` In case the argument is invalid.
        :raises: :cpp:class:`YPathError<YPathError>` In case the path is invalid.

    .. cpp:function:: virtual DataNode& create_datanode(const std::string& path) const

        Create a :cpp:class:`DataNode<DataNode>` corresponding to the path.

        This methods creates a DataNode tree based on the path passed in. The path expression must identify a single node. If the last node created is of schema type ``list``, ``leaf-list`` or ``anyxml`` that value is also set in the node.

        The returned :cpp:class:`DataNode<DataNode>` is the last node created (the terminal part of the path).

        The user is responsible for managing the memory of this returned tree. Use :cpp:func:`root` to get the root element of the this tree and use that pointer to dispose of the entire tree.

        Note in the case of List nodes the keys must be present in the path expression in the form of predicates.

        :param path: The XPath expression identifying the node.
        :return: Reference to :cpp:class:`DataNode<DataNode>` created.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>` In case the argument is invalid.
        :raises: :cpp:class:`YPathError<YPathError>` In case the path is invalid.

    .. cpp:function:: virtual Statement get_statement() const

        Return the :cpp:class:`Statement<Statement>` representing this :cpp:class:`SchemaNode<SchemaNode>`.

        Note the :cpp:class:`RootSchemaNode<RootSchemaNode>` has no YANG statement representing it.

        So this method returns an empty statement.

        :return: An empty statement.

    .. cpp:function:: virtual std::vector<Statement> get_keys() const

            Returns vector of YANG statement corresponding the the keys.

            :return: Vector of :cpp:class:`Statement` that represent keys.

    .. cpp:function:: virtual std::shared_ptr<Rpc> create_rpc(const std::string& path) const

        Create rpc and return pointer to this rpc.

        The path expression should point to a :cpp:class:`SchemaNode<SchemaNode>` that represents the :cpp:class:`Rpc<Rpc>`.

        :param path: The path to the rpc schema node
        :return: Pointer to :cpp:class:`Rpc<Rpc>` or ``nullptr``.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>` if the argument is invalid.
        :raises: :cpp:class:`YPathError<YPathError>` if the path is invalid.
