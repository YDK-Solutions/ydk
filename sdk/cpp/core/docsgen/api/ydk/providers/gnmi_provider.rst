..
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

gNMIServiceProvider
======================


.. cpp:class:: ydk::gNMIServiceProvider : public ydk::ServiceProvider

    Implementation of :cpp:class:`ServiceProvider<ydk::ServiceProvider>` for the `gNMI <https://github.com/openconfig/gnmi>`_ protocol.

    .. cpp:function:: gNMIServiceProvider( \
        path::Repository& repo, \
        std::string & address, \
        int port, \
        std::string & username, \
        std::string & password, \
        std::string & server_certificate = "", \
        std::string & private_key = "")

        Constructs an instance of ``gNMIServiceProvider`` using the provided :cpp:class:`repository<path::Repository>`, connects to gNMI server and retrieves server capabilities.

        :param repository: Reference to an instance of :cpp:class:`Repository<ydk::path::Repository>`.
        :param address: IP address of the device supporting gNMI protocol.
        :param port: Device port used to access the gNMI server.
        :param username: Username to log in to the device.
        :param password: Password to log in to the device.
        :param server_certificate: Full path to a file, which contains server certificate of authorization (public key). If not specified, it is assumed non-secure connection to gNMI server.
        :param private_key: Full path to a file, which contains private key of the application host. If not specified and **server_certificate** is defined (secure connection), the GRPC internally defined private key is used.
        :raises: :cpp:class:`YServiceProviderError<YServiceProviderError>`, if connection error occurred.

    .. cpp:function:: EncodingFormat get_encoding() const

        Returns the type of encoding supported by the service provider. In the case of gNMI service provider, :cpp:enum:`EncodingFormat::JSON<EncodingFormat>` is always returned.

    .. cpp:function:: std::vector<std::string> get_capabilities() const

        Returns gNMI server capabilities as vector of strings.
