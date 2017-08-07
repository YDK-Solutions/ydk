RESTCONF Service Provider
=========================


.. py:class:: ydk.providers.RestconfServiceProvider(repo, address, username, password, port, encoding)

    Constructs an instance of the ``RestconfServiceProvider`` to connect to a server which has to support model download. Since the class is a Python wrapper for C++ ``RestconfServiceProvider`` class, which has clean up methods implemented in its destructor. The user does not need to worry about close RESTCONF session.

    :param repo: (:py:class:`Repository<ydk.path.Repository>`) User provided repository stores cached models
    :param address: (``str``) IP address of the device supporting a restconf interface
    :param username: (``str``) Username to log in to the device
    :param password: (``str``) Password to log in to the device
    :param port: (``int``) Device port used to access the restconf interface, the default value being 80
    :param encoding: (:py:class:`EncodingFormat<ydk.types.EncodingFormat>`) Type of encoding to be used for the payload, the default being :py:attr:`JSON<ydk.types.EncodingFormat.JSON>`


    .. py:method:: get_encoding()

        Returns the type of encoding supported by the service provider.

    .. py:method:: get_session()

        Returns the instance of the :py:class:`NetconfSession<ydk.path.NetconfSession>` used to connect to the netconf server

        :return: A :py:class:`NetconfSession<ydk.path.NetconfSession>` instance.
