.. _types-ydk:

YDK Types
=========

.. go:package:: ydk/types
    :synopsis: YDK Types

.. attribute:: EncodingFormat

    representing the encoding format with underlying type ``int``:

    .. attribute:: XML

        Represented by a value of 0

    .. attribute:: JSON

        Represented by a value of 1

.. interface
.. go:struct:: ServiceProvider

    TODO

.. interface
.. go:struct:: CodecServiceProvider

    TODO

.. attribute:: Protocol

    Represents the protocol to use using underlying type ``int``

    .. attribute Restconf

        Represented by a value of 0

    .. attribute Netconf

        Represented by a value of 1

.. go:struct:: DataNode

    .. attribute:: Private

        Type is interface{}

.. go:struct:: RootSchemaNode
    
    .. attribute:: Private

        Type is interface{}

.. go:struct:: CServiceProvider

    .. attribute:: Private

        Type is interface{}

.. go:struct:: COpenDaylightServiceProvider

    .. attribute:: Private

        Type is interface{}

.. go:struct:: Repository

    .. attribute:: Path

        Type is ``string``

    .. attribute:: Private

        Type is interface{}
