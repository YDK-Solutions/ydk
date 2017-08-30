Path
====

.. go:package:: ydk/path
    :synopsis: Path API

.. function:: ExecuteRpc(provider ServiceProvider, entity Entity, filter string, dataTag string, setConfigFlag bool)

    This function executes rpc.

    :param provider: (:go:struct:`ServiceProvider<ydk/types/ServiceProvider>`).
    :param entity: (:go:struct:`Entity<ydk/types/Entity>`).
    :param filter: (``string``) A Go string.
    :param dataTag: (``string``) A Go string.
    :param setConfigFlag: (``bool``) A Go bool.
    :return: A data node representing the result of the executed rpc.
    :rtype: :go:struct:`DataNode<ydk/types/DataNode>`

.. function:: ReadDatanode(filter Entity, readDataNode DataNode)

    This function reads the top level entity from a given data node.

    :param filter: (:go:struct:`Entity<ydk/types/Entity>`)
    :param readDataNode: (:go:struct:`DataNode<ydk/types/DataNode>`)
    :return: The top entity from readDataNode.
    :rtype: :go:struct:`Entity<ydk/types/Entity>`

.. function:: ConnectToProvider(repo Repository, address, username, password string, port int)
    
    This function creates a connection to a given provider using given address, username, password, and port.

    :param repo: (:go:struct:`Repository<ydk/types/Repository>`).
    :param address: (``string``) A Go string.
    :param username: (``string``) A Go string.
    :param password: (``string``) A Go string.
    :param port: (``int``) An integer.
    :return: The service provider.
    :rtype: :go:struct:`CServiceProvider<ydk/types/CServiceProvider>`

.. function:: DisconnectFromProvider(provider CServiceProvider)

    This function disconnects from a given service provider.

    :param: provider: (:go:struct:`CServiceProvider<ydk/cgopath/CServiceProvider>`) A service provider instance.

.. function:: ConnectToRestconfProvider(path, address, username, password string, port int)
    
    This function creates a connection to a given restconf provider using given path, address, username, password, and port.

    :param path: (``string``) A Go string.
    :param address: (``string``) A Go string.
    :param username: (``string``) A Go string.
    :param password: (``string``) A Go string.
    :param port: (``int``) An integer.
    :return: The service provider.
    :rtype: :go:struct:`CServiceProvider<ydk/types/CServiceProvider>`

.. function:: DisconnectFromRestconfProvider(provider CServiceProvider)

    This function disconnects from a given restconf provider.

    :param: provider: (:go:struct:`CServiceProvider<ydk/cgopath/CServiceProvider>`) A service provider instance.

.. function:: InitCodecServiceProvider(entity Entity, repo Repository)

    This function initializes a codec service provider
    
    :param entity: :go:struct:`Entity<ydk/types/Entity>`
    :param repo: (:go:struct:`Repository<ydk/types/Repository>`).
    :return: Root data node of derrived from repo.
    :rtype: :go:struct:`RootSchemaNode<ydk/types/RootSchemaNode>`

.. function:: CodecServiceEncode(entity Entity, rootSchema RootSchemaNode, encoding EncodingFormat)

    This function encodes for codec service.

    :param entity: (:go:struct:`Entity<ydk/types/Entity>`).
    :param rootSchema: (:go:struct:`RootSchemaNode<ydk/types/RootSchemaNode>`).
    :param encoding: (:go:struct:`EncodingFormat<ydk/types/EncodingFormat>`).
    :return: The resulting payload.
    :rtype: (``string``) A Go string.

.. function:: CodecServiceDecode(rootSchema RootSchemaNode, payload string, encoding EncodingFormat, topEntity Entity)

    This function decodes payload for codec service.

    :param rootSchema: (:go:struct:`RootSchemaNode<ydk/types/RootSchemaNode>`).
    :param payload: (``string``) A Go string.
    :param encoding: (:go:struct:`EncodingFormat<ydk/types/EncodingFormat>`).
    :param topEntity: (:go:struct:`Entity<ydk/types/Entity>`)
    :return: The top level entity from resulting data node.
    :rtype: :go:struct:`Entity<ydk/types/Entity>`

.. function:: ConnectToOpenDaylightProvider(path, address, username, password string, port int, encoding EncodingFormat, protocol Protocol)

    This function connects the opendaylight provider.

    :param path: (``string``) A Go string.
    :param address: (``string``) A Go string.
    :param username: (``string``) A Go string.
    :param password: (``string``) A Go string.
    :param port: (``int``) An integer.
    :param encoding: (:go:struct:`EncodingFormat<ydk/types/EncodingFormat>`).
    :param protocol: (:go:struct:`Protocol<ydk/types/Protocol>`).
    :return: Returns the provider to which this function has connected.
    :rtype: :go:struct:`COpenDaylightServiceProvider<ydk/types/COpenDaylightServiceProvider>`

.. function:: DisconnectFromOpenDaylightProvider(provider COpenDaylightServiceProvider)

    This function disconnects the opendaylight provider.

    :param provider: (:go:struct:`COpenDaylightServiceProvider<ydk/types/COpenDaylightServiceProvider>`).


.. function:: OpenDaylightServiceProviderGetNodeIDs(provider COpenDaylightServiceProvider)

    This is a getter function for the node ids given the opendaylight service provider.

    :param provider: (:go:struct:`COpenDaylightServiceProvider<ydk/types/COpenDaylightServiceProvider>`).
    :returns: A slice of Go strings representing node ids.
    :rtype: ``[]string``

.. function:: OpenDaylightServiceProviderGetNodeProvider(provider COpenDaylightServiceProvider, nodeID string)

    This is a getter function for the node provider given the opendaylight service provider and node id.

    :param provider: (:go:struct:`COpenDaylightServiceProvider<ydk/types/COpenDaylightServiceProvider>`).
    :param nodeID: (``string``) A Go string.
    :return: The service provider.
    :rtype: :go:struct:`CServiceProvider<ydk/types/CServiceProvider>`
