ModelProvider
=============


.. cpp:class:: ydk::path::ModelProvider

  Interface for a YANG model provider which can download YANG models from a server.

    .. cpp:enum-class:: Format

       Format of model to be downloaded

       .. cpp:enumerator:: YANG

       .. cpp:enumerator:: YIN

    .. cpp:function:: virtual std::string get_model(const std::string& name, const std::string& version, Format format) = 0

        Returns the model identified by the name and version.

        :param name: Model name.
        :param version: Model revision
        :param format: Model download format.
        :return: Context of downloaded model. Empty if the model cannot be provided.

    .. cpp:function:: virtual std::string get_hostname_port() = 0

       Return the hostname and port of the connected server.
