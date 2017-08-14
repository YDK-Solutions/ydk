gNMI Service Provider
========================


.. py:class:: ydk.providers.gNMIServiceProvider(address, username, password, port=50051, repo=None)

    Constructs an instance of the ``gNMIServiceProvider`` to connect to a server which **has** to support model download. Since the class is a Python wrapper for C++ ``gNMIServiceProvider`` class, which has clean up methods implemented in its destructor. The user does not need to worry about close gNMI session.

    :param address: (``str``) IP address of the device supporting a gNMI interface
    :param port: (``int``) The port to use, defaults to 50051
    :param username: (``str``) Username to log in to the device
    :param password: (``str``) Password to log in to the device
    :param repo: User provided repository stores cached models
    :type repo: :py:class:`Repository<ydk.path.Repository>`

    .. py:method:: get_encoding()

        Returns the type of encoding supported by the service provider.
