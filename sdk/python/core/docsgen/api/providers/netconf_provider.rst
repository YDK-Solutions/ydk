NETCONF Service Provider
========================


.. py:class:: ydk.providers.NetconfServiceProvider(address, username, password=None, port=830, protocol='ssh', timeout=-1, repo=None, private_key_path="", public_key_path="")

    Constructs an instance of the `NetconfServiceProvider` class and connects to Netconf server via :py:class:`NetconfSession<ydk.path.NetconfSession>`.

    :param address: (``str``) IP address or DNS name of device, which supports Netconf server; required parameter
    :param username: (``str``) Username to log in to the device; required parameter
    :param password: (``str``) Password to log in to the device
    :param port: (``int``) The Netconf server access port; defaults to 830
    :param protocol: (``str``) Currently supported `ssh` for secure connection and `tcp` for insecure connection; default - `ssh`
    :param timeout: (``int``) The timeout in microseconds, -1 for infinite timeout, 0 for non-blocking
    :param repo: (:py:class:`Repository<ydk.path.Repository>`) User provided repository - directory, which stores cached Yang models
    :param private_key_path: (``str``) Path to private key file. Requires public_key_path field.
    :param public_key_path: (``str``) Path to public key file. Does not allow password field.

    .. py:method:: get_encoding()

        Returns the type of encoding supported by the service provider.

    .. py:method:: get_session()

        Returns the instance of the :py:class:`NetconfSession<ydk.path.NetconfSession>` used to connect to the Netconf server

        :return: A :py:class:`NetconfSession<ydk.path.NetconfSession>` instance.

    .. py:method:: get_capabilities()

        Returns a list of capabilities of the Netconf server

        :returns: A list of ``str`` representing the client's/server's capabilities