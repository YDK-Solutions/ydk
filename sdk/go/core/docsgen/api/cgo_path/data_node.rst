DataNode
========

.. go:package:: ydk/cgopath
    :synopsis: CGo Path API DataNode


.. go:struct:: DataNode

    Class represents a Path API's DataNode.

    .. note::

        The Go DataNode is wrapper for YDK C++ DataNode implementation. No constructor is defined and the user could not instantiate a ``DataNode`` instance. However, the user could get an instance of ``DataNode`` through :go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>`:

        .. code-block:: go

            import (
                "C"
                "github.com/CiscoDevNet/ydk-go/ydk"
                "github.com/CiscoDevNet/ydk-go/ydk/cgopath"
                "github.com/CiscoDevNet/ydk-go/ydk/providers"
                "testing"
            )

            repo := C.RepositoryInit()
            defer C.RepositoryFree(repo)
            var address *C.char = C.CString("localhost")
            defer C.free(unsafe.Pointer(address))
            var username *C.char = C.CString("admin")
            defer C.free(unsafe.Pointer(username))
            var password *C.char = C.CString("admin")
            defer C.free(unsafe.Pointer(password))

            provider := C.NetconfServiceProviderInitWithRepo(repo, address, username, password, 12022)
            defer C.NetconfServiceProviderFree(provider)

            rootSchema := C.ServiceProviderGetRootSchema(provider)      # <-- root_schema is an instance of RootSchemaNode

            var bgp_path *C.char = C.CString("openconfig-bgp:bgp")
            defer C.free(unsafe.Pointer(bgp_path))
            bgp := C.RootSchemaNodeCreate(root_schema, bgp_path)        # <-- bgp is an instance of DataNode
