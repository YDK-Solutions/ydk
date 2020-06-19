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

YError
=========

.. cpp:class:: ydk::YError : public std::exception

    Base class for YDK exceptions.

    .. cpp:member:: std::string err_msg

    .. cpp:function:: YError(const std::string& msg)

    Class constructor

    .. cpp:function:: char* what() const noexcept

    Function to access error message

.. cpp:class:: ydk::YClientError : public ydk::YError

    YDK exception thrown when error occurred in protocol client application.

.. cpp:class:: ydk::YServiceProviderError : public ydk::YError

    YDK exception thrown when service provider errors occurred.

.. cpp:class:: ydk::YServiceError : public ydk::YError

    YDK exception thrown when service errors occurred.

.. cpp:class:: ydk::YIllegalStateError : public ydk::YError

    YDK exception thrown when an operation/service is invoked on an object that is not in the right state.

.. cpp:class:: ydk::YInvalidArgumentError : public ydk::YError

    YDK exception thrown when a function given parameter(s) with wrong values.

.. cpp:class:: ydk::YOperationNotSupportedError : public ydk::YError

    YDK exception thrown when specified yfilter is not supported by protocol.

.. cpp:class:: ydk::YModelError : public ydk::YError

    YDK exception thrown when a model constraint is violated.
