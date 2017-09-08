YDK CGo Path API
================

.. go:package:: ydk/path
    :synopsis: YDK Path API

This package contains YDK CGo Path API. These are the functions that were defined in C/CPP and ported to Go.

YDK Go provides a new interface in the form of Path API, which can be used to write apps using a generic API, using xpath-like path expression to create and access YANG data nodes. Internally, the nodes created using the YDK model API are converted to Path API data nodes for validation and encoding to respective protocol payloads.


.. toctree::
    :maxdepth: 2

    capability.rst
    codec_service.rst
    data_node.rst
    repository.rst
    root_schema_node.rst
    rpc.rst
    schema_node.rst
    service_provider.rst
    cgo_path.rst