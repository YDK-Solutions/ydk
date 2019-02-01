gNMISession
==============

.. module:: ydk.gnmi.path
    :synopsis: gNMISession


.. py:class:: gNMISession(repo, address, port=57400, username, password, server_certificate="", private_key="")

    :param repo: Instance of :py:class:`Repository<ydk.path.Repository>` with path to local directory containing the the `ydk yang model <https://raw.githubusercontent.com/CiscoDevNet/ydk-gen/1344b3f22d746764f17536ac4e666836de4ba84d/sdk/cpp/core/tests/models/ydk%402016-02-26.yang>`_ along with all the yang models supported on the gNMI server.
    :param address: (``str``) Host address of the device supporting a gNMI interface
    :param port: (``int``)Port on which the gNMI interface can be accessed on the device. If not specified, the default value of ``57400`` is assigned.
    :param username: (``str``) Username.
    :param password: (``str``) Password.
    :param server_certificate: (``str``) Full path to a file, which contains server certificate of authorization (public key). If not specified, it is assumed non-secure connection to gNMI server.
    :param private_key: (``str``) Full path to a file, which contains private key of the application host. If not specified and **server_certificate** is defined (secure connection), the GRPC internally defined private key is used.

    .. py:method:: get_root_schema()

        :returns: :py:class:`RootSchemaNode<RootSchemaNode>` for this gNMI session.

    .. py:method:: invoke(rpc)

        Executes gNMI RPCs **gnmi-set**, **gnmi-get**, and **gnmi-caps**.
 
        :param rpc: (:py:class:`Rpc<ydk.path.Rpc>`) RPC to be executed.
        :returns: :py:class:`DataNode<DataNode>` for **gnmi-get** RPC, or ``None`` otherwise.

    .. py:method:: subscribe(rpc, output_callback_function=None)

       Executes **gnmi-subscribe** RPC. Returns subscription data over specified callback.

        :param rpc: (:py:class:`Rpc<ydk.path.Rpc>`) RPC to be executed.
        :param output_callback_function: (``func(str)``) Callback function, which is used to process the subscription data. 
                                         The subscription data returned to the user as a string representation of protobuf **SubscribeResponse** message.
                                         If not specified, the response is printed to system stdout.
        :returns: None.
