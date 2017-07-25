ServiceProvider
===============

.. cpp:namespace:: ydk

.. cpp:class:: ServiceProvider

Class represents the base class :cpp:class:`ServiceProvider<ydk::ServiceProvider>` to be inherited for implementation.

    .. cpp:function:: virtual EncodingFormat get_encoding() const

        Returns the type of encoding supported by the service provider.

    .. cpp:function:: virtual Session get_session() const

        Returns the instance of the :cpp:class:`Session<path::Session>` used to connect to the server

    .. cpp:function:: virtual ~ServiceProvider()
