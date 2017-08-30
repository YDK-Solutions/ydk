Repository
==========

.. go:package:: ydk/cgopath
    :synopsis: CGo Path API Repository

.. go:struct:: Repository(*args)

    Repository is used to create a :go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>` given a set of Capabilities. Behind the scenes the repository is responsible for loading and parsing the YANG modules and creating the :go:struct:`SchemaNode<ydk/cgopath/SchemaNode>` tree. Service provider is expected to use the function :func:`CreateRootSchema` to generate the :go:struct:`RootSchemaNode<ydk/cgopath/RootSchemaNode>`.

    :param args: Search directory of type ``str`` or nothing.

    .. note::

        If a directory path of type ``str`` is provided, model search path will be located in this path, otherwise, default location ``~/.cache/ydk`` will be used. Example usage:

        .. code-block:: go

            import (
                "C"
                "github.com/CiscoDevNet/ydk-go/ydk"
            )
            
            repo := C.RepositoryInit()
            defer C.RepositoryFree(repo)