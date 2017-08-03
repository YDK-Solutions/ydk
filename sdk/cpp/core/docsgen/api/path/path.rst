YDK Path API
============
YDK C++ provides a new interface in the form of Path API, which can be used to write apps using a generic API, using xpath-like path expression to create and access YANG data nodes. Internally, the nodes created using the YDK model API are converted to Path API data nodes for validation and encoding to respective protocol payloads.

Path API Nodes
~~~~~~~~~~~~~~

.. toctree::
   :maxdepth: 2

   nodes/annotation.rst
   nodes/statement.rst
   nodes/data_node.rst
   nodes/schema_node.rst
   nodes/root_schema_node.rst

Path API Sessions
~~~~~~~~~~~~~~~~~

.. toctree::
   :maxdepth: 2

   sessions/session.rst
   sessions/netconf_session.rst
   sessions/restconf_session.rst

Path API Errors
~~~~~~~~~~~~~~~

.. toctree::
   :maxdepth: 2

   errors/codec_error.rst
   errors/core_error.rst
   errors/data_validation_error.rst
   errors/path_error.rst

Path API
~~~~~~~~

.. toctree::
   :maxdepth: 2

   capability.rst
   codec.rst
   model_provider.rst
   repository.rst
   rpc.rst
   types.rst
   validation_service.rst
