.. _ref-codecservice:

Codec
=====

Codec class is part of YDK path API, which operates on generic path-based YANG data nodes.

.. cpp:class:: ydk::path::Codec

    .. cpp:function:: virtual std::string encode(const std::unique_ptr<DataNode> datanode, EncodingFormat format, bool pretty)

        Encode the given DataNode Tree.

        :param datanode: The :cpp:class:`DataNode<DataNode>` to encode.
        :param format: format to encode to, either :cpp:enumerator:`JSON<EncodingFormat::JSON>` or :cpp:enumerator:`XML<EncodingFormat::XML>`.
        :param pretty: The output is indented for human consumption, if pretty is ``true``.
        :return: The encoded string.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>`, if errors appear during encoding.

    .. cpp:function:: virtual std::string encode(std::vector<DataNode*> data_nodes, EncodingFormat format, bool pretty)

        Encode multiple DataNodes.

        :param data_nodes: An instance of **std::vector<DataNode\*>** class, which contains one or more data nodes.
        :param format: format to encode to, either :cpp:enumerator:`JSON<EncodingFormat::JSON>` or :cpp:enumerator:`XML<EncodingFormat::XML>`.
        :param pretty: The output is indented for human consumption. if **pretty** is ``true``.
        :return: The encoded string.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>`, errors appear during encoding.

    .. cpp:function:: std::shared_ptr<DataNode> decode(const RootSchemaNode& root_schema, const std::string& payload, Format format)

        Decode the buffer to return a DataNode.

        :param root_schema: The root schema to use.
        :param payload: The string representation of the :cpp:class:`DataNode<DataNode>`.
        :param format: Format of the **payload**: either :cpp:enumerator:`JSON<EncodingFormat::JSON>` or :cpp:enumerator:`XML<EncodingFormat::XML>`.
        :return: The :cpp:class:`DataNode<DataNode>` instance or ``nullptr`` in case of error.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>`, if errors appear during decoding.

    .. cpp:function:: std::shared_ptr<DataNode> decode_rpc_output(RootSchemaNode & root_schema, const std::string& payload, const std:: string & rpc_path, EncodingFormat format)

        Decode the rpc output to return a DataNode.

        :param root_schema: The root schema to use.
        :param payload: String representation of one or multiple :cpp:class:`DataNode<DataNode>`.
        :param format: Format of the **payload**: either :cpp:enumerator:`JSON<EncodingFormat::JSON>` or :cpp:enumerator:`XML<EncodingFormat::XML>`.
        :return: The :cpp:class:`DataNode<DataNode>` instance or ``nullptr`` in case of error.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>`, if errors appear during decoding.
