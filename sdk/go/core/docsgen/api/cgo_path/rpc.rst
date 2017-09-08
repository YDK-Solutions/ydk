Rpc
===

.. go:package:: ydk/cgopath
    :synopsis: CGo Path API Rpc


.. go:struct:: Rpc

    Instances of this class represent a YANG rpc and are modelled as Callables. The input data node tree is used to populate the input parameters to the rpc if any. The Callable takes as a parameter the :go:struct:`ServiceProvider<ServiceProvider>` that can execute this rpc as its parameter returning a :go:struct:`DataNode<DataNode>` instance if output is available.