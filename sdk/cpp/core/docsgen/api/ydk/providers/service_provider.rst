ServiceProvider
===============


.. cpp:class:: ydk::ServiceProvider

Interface for all ServiceProvider implementations.

    Concrete instances of ServiceProviders are expected to extend this interface.

    .. cpp:function:: virtual EncodingFormat get_encoding()

        Returns the type of encoding supported by the service provider.

    .. cpp:function:: virtual const path::Session& get_session()

        Returns reference to the session used to connect to the server

    .. cpp:function:: virtual ~ServiceProvider()
