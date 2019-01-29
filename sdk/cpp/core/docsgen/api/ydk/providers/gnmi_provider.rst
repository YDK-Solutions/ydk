gNMIServiceProvider
======================


.. cpp:class:: ydk::gNMIServiceProvider : public ydk::ServiceProvider

    Implementation of :cpp:class:`ServiceProvider<ydk::ServiceProvider>` for the `gNMI <https://github.com/openconfig/gnmi>`_ protocol.

    .. cpp:function:: gNMIServiceProvider( \
        path::Repository& repo, \
        std::string & address, \
        int port, \
        std::string & username, \
        std::string & password, \
        std::string & server_certificate = "", \
        std::string & private_key = "")

        Constructs an instance of ``gNMIServiceProvider`` using the provided :cpp:class:`repository<path::Repository>`, connects to gNMI server and retrieves server capabilities.

        :param repository: Reference to an instance of :cpp:class:`Repository<ydk::path::Repository>`.
        :param address: IP address of the device supporting gNMI protocol.
        :param port: Device port used to access the gNMI server.
        :param username: Username to log in to the device.
        :param password: Password to log in to the device.
        :param server_certificate: Full path to a file, which contains server certificate of authorization (public key). If not specified, it is assumed non-secure connection to gNMI server.
        :param private_key: Full path to a file, which contains private key of the application host. If not specified and **server_certificate** is defined (secure connection), the GRPC internally defined private key is used.
        :raises: YServiceError, if connection error occured.

    .. cpp:function:: EncodingFormat get_encoding() const

        Returns the type of encoding supported by the service provider. In the case of gNMI service provider, :cpp:enum:`EncodingFormat::JSON<EncodingFormat>` is always returned.

    .. cpp:function:: std::vector<std::string> get_capabilities() const

        Returns gNMI server capabilities as vector of strings.
