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

Codec Service Provider
======================

.. py:class:: ydk.providers.CodecServiceProvider(type=EncodingFormat.XML, repo=None)

    A provider to be used with :py:class:`CodecService<ydk.services.CodecService>` for performing encoding and decoding.

    :param type: An argument specifies encoding format, could be a Python string (``xml`` or ``json``) or an instance of :py:class:`EncodingFormat<ydk.types.EncodingFormat>`.
    :type type: ``string`` or :py:class:`EncodingFormat<ydk.types.EncodingFormat>`
    :param repo: User provided repository stores cached models.
    :type repo: :py:class:`Repository<ydk.path.Repository>`

    .. py:method:: get_root_schema(bundle_name)

        Return root schema node for bundle_name defined in the 'initialize' method.

        :param bundle_name: (``str``) Bundle name.
        :return: :py:class:`RootSchemaNode<ydk.path.RootSchemaNode>` for this bundle.

    .. py:method:: initialize(bundle_name, models_path)

        Initializes root schema in the user defined location of YANG models.

        :param bundle_name: (``str``): user defined bundle name.
        :param models_path: (``str``): location (directory) for local YANG models.
