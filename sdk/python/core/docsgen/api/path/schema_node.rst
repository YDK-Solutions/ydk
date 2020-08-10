..
  #  YDK-YANG Development Kit
  #  Copyright 2016 Cisco Systems. All rights reserved
  # *************************************************************
  # Licensed to the Apache Software Foundation (ASF) under one
  # or more contributor license agreements.  See the NOTICE file
  # distributed with this work for additional information
  # regarding copyright ownership.  The ASF licenses this file
  # to you under the Apache License, Version 2.0 (the
  # "License"); you may not use this file except in compliance
  # with the License.  You may obtain a copy of the License at
  #
  #   http:#www.apache.org/licenses/LICENSE-2.0
  #
  #  Unless required by applicable law or agreed to in writing,
  # software distributed under the License is distributed on an
  # "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  # KIND, either express or implied.  See the License for the
  # specific language governing permissions and limitations
  # under the License.
  # *************************************************************
  # This file has been modified by Yan Gorelik, YDK Solutions.
  # All modifications in original under CiscoDevNet domain
  # introduced since October 2019 are copyrighted.
  # All rights reserved under Apache License, Version 2.0.
  # *************************************************************

SchemaNode
==========

.. module:: ydk.path
    :synopsis: Path API' SchemaNode


.. py:class:: SchemaNode

    Class represents a Node in the SchemaTree.

    .. note::

        The Python SchemaNode is wrapper for YDK C++ SchemaNode implementation. No constructor is defined and the user could not instantiate a ``SchemaNode`` instance. However, the user could get an instance of ``SchemaNode`` through :py:class:`DataNode<ydk.path.DataNode>`:

        .. code-block:: python

            >>> from ydk.path import NetconfSession
            >>> session = NetconfSession('127.0.0.1', 'admin', 'admin', 830)
            >>> root_schema = session.get_root_schema()                               # <-- root_schema is an instance of RootSchemaNode
            >>> bgp = root_schema.create_datanode('openconfig-bgp:bgp')               # <-- bgp is an instance of DataNode
            >>> schema_node = bgp.get_schema_node()                                   # <-- schema node for bgp

    .. py:method:: find(path)

        Finds descendant nodes that match the given xpath expression. This API finds descendant nodes in the schema node tree that satisfy the given path expression. See :ref:`howto-path`.

        :param path: (``str``) The path expression.
        :return: List of schema node satisfies the criterion.
        :rtype: list of :py:class:`SchemaNode`
        :raises: YPathError, if the path expression in invalid, see error code for details.
        :raises: YInvalidArgumentError, if the argument is invalid.

   .. py:method:: get_keys()

        Returns vector of YANG :py:class:`Statement<Statement>`, which represents keys for the list node.

        :return: List of statements of type :cpp:class:`Statement<Statement>`, which represents list node keys.

   .. py:method:: get_parent()

        Get the parent node of this schema node in the tree.

        :return: Parent schema node.
        :rtype: :py:class:`SchemaNode`

    .. py:method:: get_path()

        Get the path expression representing this schema node in in the schema node tree.

        :return: Path to this schema node.
        :rtype: A Python string

    .. py:method:: get_root()

        Get the root schema node of current schema node.

        :return: Root schema node of current schema node.
        :rtype: :py:class:`SchemaNode`

    .. py:method:: get_statement()

        Get current schema node's YANG statement.

        :return: YANG statements for this schema node.
        :rtype: :py:class:`Statement`
