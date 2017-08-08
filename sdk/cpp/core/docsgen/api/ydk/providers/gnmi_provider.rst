gNMIServiceProvider
======================


.. cpp:class:: ydk::gNMIServiceProvider : public ydk::ServiceProvider

    Implementation of :cpp:class:`ServiceProvider<ydk::ServiceProvider>` for the `gNMI <https://github.com/openconfig/gnmi>`_ protocol.

    .. cpp:function:: gNMIServiceProvider(\
        std::string address, \
        std::string username, \
        std::string password, \
        int port = 50051)

        Constructs an instance of ``gNMIServiceProvider`` connect to a server 

        :param address: IP address of the device supporting a gNMI interface.
        :param username: Username to log in to the device.
        :param password: Password to log in to the device.
        :param port: Device port used to access the gNMI interface. Default value is 50051.

    .. cpp:function:: gNMIServiceProvider(\
        const path::Repository& repo, \
        std::string address, \
        std::string username, \
        std::string password, \
        int port = 50051)

        Constructs an instance of ``gNMIServiceProvider`` using the provided :cpp:class:`repository<path::Repository>`

        :param repository: Reference to an instance of :cpp:class:`Repository<ydk::path::Repository>`
        :param address: IP address of the device supporting a gNMI interface
        :param username: Username to log in to the device
        :param password: Password to log in to the device
        :param port: Device port used to access the gNMI interface. Default value is 50051.

    .. cpp:function:: EncodingFormat get_encoding() const

        Returns the type of encoding supported by the service provider. In the case of gNMI service provider, :cpp:enum:`EncodingFormat::JSON<EncodingFormat>` is returned.

    .. cpp:function:: ~gNMIServiceProvider()
