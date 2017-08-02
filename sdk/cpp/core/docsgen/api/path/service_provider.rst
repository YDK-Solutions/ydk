.. _ref-serviceprovider:

ServiceProvider
===============

.. cpp:namespace:: ydk::path

.. cpp:class:: ServiceProvider

    Interface for all ServiceProvider implementations.

    Concrete instances of ServiceProviders are expected to extend this interface.

    .. cpp:function:: EncodingFormat get_encoding()

        Returns the type of encoding supported by the service provider.

    .. py:method:: const path::Session& get_session()

        Returns the instance of the :py:class:`Session<ydk::path::Session>` used to connect to the server

        :return: A :py:class:`Session<ydk::path::Session>` instance.


    .. cpp:function:: virtual ~ServiceProvider()
