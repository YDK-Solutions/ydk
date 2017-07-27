.. _ref-serviceprovider:

ServiceProvider
===============

.. toctree::
   :maxdepth: 2

   netconf_provider.rst
   restconf_provider.rst
   codec_provider.rst
   opendaylight_provider.rst

.. cpp:namespace:: ydk::path

.. cpp:class:: ServiceProvider

Interface for all ServiceProvider implementations.

Concrete instances of ServiceProviders are expected to extend this interface.

    .. cpp:function:: EncodingFormat get_encoding()

        Returns the type of encoding supported by the service provider.

    .. py:method:: get_session()

        Returns the instance of the :py:class:`Session<ydk.path.Session>` used to connect to the server

        :return: A :py:class:`Session<ydk.path.Session>` instance.


    .. cpp:function:: virtual ~ServiceProvider()
