CGo Path
========

.. go:package:: ydk/path
    :synopsis: CGo Path API

.. function:: CapabilityCreate(cstate C.YDKStatePtr, cstate model CString, revision CString)

    Creates a capability
    
    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param model: (``CString``) Model name.
    :param revision: (``CString``) Model revision.
    :return: Created capability
    :rtype: :go:struct:`Capability<ydk/cgopath/Capability>`

    .. note::

        If no revision is provided, use empty string:

        .. code-block:: go

            import "C"
            var cstate C.YDKStatePtr = C.YDKStateCreate()
            defer C.YDKStateFree(cstate)

            var cap1 C.Capability
            var cap2 C.Capability
            
            cap1 = C.CapabilityCreate(cstate, C.CString("openconfig-bgp"), C.CString(""))
            cap2 = C.CapabilityCreate(cstate, C.CString("openconfig-bgp"), C.CString("2015-10-09"))

            defer C.CapabilityFree(cap1)
            defer C.CapabilityFree(cap2)

.. function:: CapabilityFree(capability Capability)

    Deletes a capability

    :param capability: (:go:struct:`Capability<ydk/cgopath/Capability>`) capability to delete.




.. function:: RepositoryInitWithPath(cstate C.YDKStatePtr, path CString)
    
    Constructs an instance of the ``Repository`` struct.

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param path: (``CString``) The path in the filesystem where yang files can be found.
    :return: The initialized repository
    :rtype: :go:struct:`Repository<ydk/cgopath/Repository>`

.. function:: RepositoryInit()

    Constructs an instance of the ``Repository`` struct. If the server supports model download, the repo will attempt to download all models from the server.

    :return: The initialized repository
    :rtype: :go:struct:`Repository<ydk/cgopath/Repository>`

.. function:: RepositoryCreateRootSchemaWrapper(cstate C.YDKStatePtr, repo Repository, []caps Capability, capSize int)

    Creates a wrapper for root root schema.

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param repo: (:go:struct:`Repository<ydk/cgopath/Repository>`)
    :param []caps: Slice of :go:struct:`Capability<ydk/cgopath/Capability>`)
    :param capSize: (``int``) Size of the slice.
    :return: root schema wrapper
    :rtype: :go:struct:`RootSchemaWrapper<ydk/cgo/RootSchemaWrapper>`

.. function:: RepositoryFree(repo Repository)

    Deletes a repository

    :param repo: (:go:struct:`Repository<ydk/cgopath/Repository>`) repository to delete.




.. function:: NetconfServiceProviderInitWithRepo(cstate C.YDKStatePtr, repo Repository, address, username, password CString, port int)

    Constructs an instance of the ``NetconfServiceProvider`` using the provided :go:struct:`Repository<ydk/cgopath/Repository>`

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param repo: (:go:struct:`Repository<ydk/cgopath/Repository>`) Repository with which to initialize.
    :param address: (``CString``) IP address of the device supporting a netconf interface
    :param username: (``CString``) Username to log in to the device
    :param password: (``CString``) Password to log in to the device
    :param port: (``integer``) Device port used to access the netconf interface.
    :return: The initialized service provider
    :rtype: :go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`

.. function:: NetconfServiceProviderInit(cstate C.YDKStatePtr, address, username, password CString, port int)

    Constructs an instance of the ``NetconfServiceProvider`` to connect to a server which *has* to support model download

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param address: (``CString``) IP address of the device supporting a netconf interface
    :param username: (``CString``) Username to log in to the device
    :param password: (``CString``) Password to log in to the device
    :param port: (``integer``) Device port used to access the netconf interface.
    :return: The initialized service provider
    :rtype: :go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`

.. function:: ServiceProviderGetRootSchema(cstate C.YDKStatePtr, provider ServiceProvider)
    
    Returns the :go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>` tree supported by this instance of the :go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param provider: (:go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`) Service provider from which to get root schema node
    :return: The root schema node tree supported by given service provider
    :rtype: :go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>`

.. function:: ServiceProviderGetEncoding(provider ServiceProvider)

    Returns the type of encoding supported by the service provider.

    :param provider: (:go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`) Service provider from which to get the encoding
    :return: The encoding format supported by the service provider
    :rtype: :go:struct:`EncodingFormat<ydk/types/EncodingFormat>`

.. function:: NetconfServiceProviderFree(provider ServiceProvider)

    Deletes the netconf service provider

    :param repo: (:go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`) provider to delete.



.. function:: RestconfServiceProviderInitWithRepo(cstate C.YDKStatePtr, repo Repository, address, username, password CString, port int)

    Constructs an instance of the ``RestconfServiceProvider`` using the provided :go:struct:`Repository<ydk/cgopath/Repository>`

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param repo: (:go:struct:`Repository<ydk/cgopath/Repository>`) Repository with which to initialize.
    :param address: (``CString``) IP address of the device supporting a restconf interface
    :param username: (``CString``) Username to log in to the device
    :param password: (``CString``) Password to log in to the device
    :param port: (``integer``) Device port used to access the restconf interface.
    :return: The initialized service provider
    :rtype: :go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`

.. function:: RestconfServiceProviderFree(provider ServiceProvider)
    
    Deletes the restconf service provider.

    :param provider: (:go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`) provider to delete.



.. function:: OpenDaylightServiceProviderInitWithRepo(cstate C.YDKStatePtr, repo Repository, address, username, password CString, port int, encoding EncodingFormat, protocol Protocol)

    Constructs an instance of the ``OpenDaylightServiceProvider`` using the provided :go:struct:`Repository<ydk/cgopath/Repository>`

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param repo: (:go:struct:`Repository<ydk/cgopath/Repository>`) Repository with which to initialize.
    :param address: (``CString``) IP address of the device supporting a opendaylight interface
    :param username: (``CString``) Username to log in to the device
    :param password: (``CString``) Password to log in to the device
    :param port: (``integer``) Device port used to access the opendaylight interface.
    :return: The initialized service provider
    :rtype: :go:struct:`OpenDaylightServiceProvider<ydk/cgopath/OpenDaylightServiceProvider>`

.. function:: OpenDaylightServiceProviderFree(provider OpenDaylightServiceProvider)
    
    Deletes the restconf service provider.

    :param provider: (:go:struct:`OpenDaylightServiceProvider<ydk/cgopath/OpenDaylightServiceProvider>`) provider to delete.

.. function:: OpenDaylightServiceProviderGetNodeProvider(cstate C.YDKStatePtr, provider OpenDaylightServiceProvider, nodeId CString)

    Returns service provider given a node id

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param provider: (:go:struct:`OpenDaylightServiceProvider<ydk/cgopath/OpenDaylightServiceProvider>`)
    :param nodeId: (``CString``) Id of the node
    :return: The provider associated with the node
    :rtype: :go:struct:`ServiceProvider<ydk/cgopath/ServiceProvider>`


.. function:: OpenDaylightServiceProviderGetNodeIDByIndex(cstate C.YDKStatePtr, provider OpenDaylightServiceProvider, idx int)

    Returns node id with given index.

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param provider: (:go:struct:`OpenDaylightServiceProvider<ydk/cgopath/OpenDaylightServiceProvider>`)
    :param idx: (``int``) index with which to get node id
    :return: node id
    :rtype: ``CString``



.. function:: CodecInit()
    
    :return: Initialized codec instance.
    :rtype: :go:struct:`Codec<ydk/cgopath/Codec>`

.. function:: CodecFree(codec Codec)

    :param: (:go:struct:`Codec<ydk/cgopath/Codec>`) Codec to free.

.. function:: CodecEncode(cstate C.YDKStatePtr, codec Codec, dataNode DataNode, encoding EncodingFormat, pretty bool)

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param codec: (:go:struct:`Codec<ydk/cgopath/Codec>`) Codec to encode.
    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`) Path ``DataNode`` to encode.
    :param encoding: (:go:struct:`EncodingFormat<ydk/types/EncodingFormat>`) Encoding format.
    :param pretty: (``bool``) Pretty format flag.
    :return: payload
    :rtype: ``CString``

.. function:: CodecDecode(cstate C.YDKStatePtr, codec Codec, rootSchemaNode RootSchemaNode, payload CString, encoding EncodingFormat)

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param codec: (:go:struct:`Codec<ydk/cgopath/Codec>`) Codec to decode.
    :param rootSchemaNode: (:go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>`) A Path ``RootSchemaNode``
    :param payload: (``CString``) Payload to decode.
    :param encoding: (:go:struct:`EncodingFormat<ydk/types/EncodingFormat>`) Encoding format.
    :return: data node
    :rtype: :go:struct:`DataNode<ydk/cgopath/DataNode>`



.. function:: RootSchemaNodeCreate(cstate C.YDKStatePtr, rootSchemaNode RootSchemaNode, path CString)

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param rootSchemaNode: (:go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>`)
    :param path: (``CString``) The XPath expression identifying the node relative to the root of the schema tree.
    :return: data node
    :rtype: :go:struct:`DataNode<ydk/cgopath/DataNode>`

.. function:: RootSchemaNodeRpc(cstate C.YDKStatePtr, rootSchemaNode RootSchemaNode, path CString)

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param rootSchemaNode: (:go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>`)
    :param path: (``CString``) The path to the rpc schema node.
    :return: Rpc instance that is created.
    :rtype: :go:struct:`Rpc`

    Excample usage for creating a ``ydk:create`` rpc:

        .. code-block:: go
            :linenos:

            import (
                "C"
                "github.com/CiscoDevNet/ydk-go/ydk"
                "github.com/CiscoDevNet/ydk-go/ydk/path"
                "github.com/CiscoDevNet/ydk-go/ydk/providers"
            )

            var address *C.char = C.CString("127.0.0.1")
            defer C.free(unsafe.Pointer(address))
            var username *C.char = C.CString("admin")
            defer C.free(unsafe.Pointer(username))
            var password *C.char = C.CString("admin")
            defer C.free(unsafe.Pointer(password))


            var cport C.int = C.int(port)

            provider := C.NetconfServiceProviderInit(address, username, password, cport)
            defer C.NetconfServiceProviderFree(provider)

            rootSchema := C.ServiceProviderGetRootSchema(provider)      # <-- rootSchema is an instance of RootSchemaNode

            var createPath *C.char = C.CString("ydk:create")
            defer C.free(unsafe.Pointer(createPath))
            createRpc := C.RootSchemaNodeRpc(rootSchema, createPath)

.. function:: RootSchemaWrapperUnwrap(wrapper RootSchemaWrapper)

    :param wrapper: (:go:struct:`RootSchemaWrapper<ydk/cgopath/RootSchemaWrapper>`)
    :return: The root schema node contained in the wrapper.
    :rtype: :go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>`



.. function:: RpcInput(cstate C.YDKStatePtr, rpc Rpc)

    Get the input data tree.

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param rpc: (:go:struct:`Rpc<ydk/cgopath/Rpc>`)
    :return: A data node representing the input data tree.
    :rtype: :go:struct:`DataNode<ydk/cgopath/DataNode>`

.. function:: RpcExecute(cstate C.YDKStatePtr, rpc Rpc, provider ServiceProvider)

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param rpc: (:go:struct:`Rpc<ydk/cgopath/Rpc>`)
    :param provider: (:go:struct:`ServiceProvider`) The service provider.
    :return: data node
    :rtype: :go:struct:`DataNode<ydk/cgopath/DataNode>`



.. function:: DataNodeCreate(cstate C.YDKStatePtr, dataNode DataNode, path, value CString)

    Create a DataNode corresponding to the path and set its value.

    :param cstate: (:cpp:class:`YDKStatePtr`) Current state of execution
    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`)
    :param path: (``CString``) The XPath expression identifying the node.
    :param value: (``CString``) The value to be set.
    :return: data node
    :rtype: :go:struct:`DataNode<ydk/cgopath/DataNode>`

.. function:: DataNodeGetArgument(dataNode DataNode)

    Returns a given data node’s argument.
    
    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`)
    :return: Argument of the given data node.
    :rtype: (``CString``)

.. function:: DataNodeGetKeyword(dataNode DataNode)
    
    Returns a given data node’s keyword.
    
    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`)
    :return: Keyword of the given data node.
    :rtype: (``CString``)

.. function:: DataNodeGetPath(dataNode DataNode)

    Returns the path expression representing the given data node in in the NodeTree.

    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`)
    :return: Path to the given data node.
    :rtype: (``CString``)

.. function:: DataNodeGetValue(dataNode DataNode)
    
    Returns the given data node’s value.
    
    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`)
    :return: Value of the given data node.
    :rtype: (``CString``)

.. function:: DataNodeGetParent(dataNode DataNode)

    Returns the given data node’s parent.
    
    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`)
    :return: Parent of the given data node.
    :rtype: :go:struct:`DataNode<ydk/cgopath/DataNode>`

.. function:: DataNodeAddAnnotation(dataNode DataNode, annotation CString)

    This function adds a given annotation to a given data node.

    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`) DataNode to which to add annotation.
    :param annotation: (``CString``) Annotation to be added.

.. function:: DataNodeGetChildren(dataNode DataNode)

    Return list of children for a given data node.

    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`) DataNode from which to get children.
    :return: A struct containing a list of data node children.
    :rtype: :go:struct:`DataNodeChildren<ydk/cgopath/DataNodeChildren>`

.. function:: DataNodeGetSegmentPath(dataNode DataNode)

    Return list of children for a given data node.

    :param dataNode: (:go:struct:`DataNode<ydk/cgopath/DataNode>`) DataNode from which to get children.
    :return: The segment path for the given data node.
    :rtype: (``CString``)




.. function:: EnableLogging(level LogLevel)

    Enables logging

    :param level: (:go:struct:`LogLevel<ydk/cgopath/LogLevel>`)

.. function:: GetLoggingLevel()

    Returns logging level

    :return: The logging level
    :rtype: :go:struct:`LogLevel<ydk/cgopath/LogLevel>`
