:mod:`ydk.services` --- YDK Services
====================================

.. module:: ydk.services
    :synopsis: YDK Services

.. contents:: Table of Contents

This module contains YDK Python supported services and datatype used for NETCONF service.

YDK Services
------------

.. toctree::
   :maxdepth: 2

   crud_service.rst
   netconf_service.rst
   codec_service.rst
   executor_service.rst

YDK Services datatype
---------------------

As specified in `NETCONF RFC 6241 <https://tools.ietf.org/html/rfc6241#section-5.1>`_, YDK provides datatype for configuration datastores.

.. py:class:: DataStore

    Netconf datastore type

    .. data:: candidate

        Candidate configuration datastore

    .. data:: running

        Running configuration datastore

    .. data:: startup

        Startup configuration datastore

    .. data:: url

        URL
