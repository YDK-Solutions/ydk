.. _ref-repository:

Repository
==========


.. cpp:class:: ydk::path::Repository

    The Repository is responsible for loading module locally or downloading module remotely if certain model providers are provided; It holds a variable pointing to the location for model downloads and could be used as a factory class to create :cpp:class:`RootSchemaNode<RootSchemaNode>` based on a list of :cpp:class:`Capability<Capability>` passed in or list of capability lookup tables and a list of :cpp:class:`Capability<Capability>` passed in.

    .. cpp:member:: std::string path

        Location where YANG models are present and/or downloaded to

    .. cpp:function:: Repository(ModelCachingOption caching_option = ModelCachingOption::PER_DEVICE)

        Create an instance of ``Repository``, the model downloading path is set to ``~/.ydk``.

        If the caching option is set to :cpp:enum:`PER_DEVICE<ModelCachingOption::PER_DEVICE>`, directories with format of ``~/.ydk/hostname:port`` (``~/.ydk/127.0.0.1:1220`` for example) are used for model downloading. Otherwise, a common directory(``~/.ydk/common_cache``) is used for all devices.

        :param caching_option: Model caching option.

        :raises: :cpp:class:`YIllegalStateError<YIllegalStateError>` If the model downloads directory could not be created.

    .. cpp:function:: Repository(const std::string& search_dir, ModelCachingOption caching_option = ModelCachingOption::PER_DEVICE)

        Create an instance of ``Repository``, the model downloading path is set to ``search_dir``.

        :param search_dir: The path in the filesystem where yang files can be found.
        :param caching_option: Model caching option.

        :raises: :cpp:class:`YIllegalStateError<YIllegalStateError>` If the model downloads directory could not be created.

    .. cpp:function:: std::shared_ptr<RootSchemaNode> create_root_schema(const std::vector<Capability>& capabilities)

        Creates a pointer of :cpp:class:`RootSchemaNode<RootSchemaNode>` based on the list of capabilities passed in.

        :param capabilities: List of capabilities available.
        :return: Pointer to the :cpp:class:`RootSchemaNode<RootSchemaNode>`.

    .. cpp:function:: std::shared_ptr<RootSchemaNode> create_root_schema(\
            const std::vector<std::unordered_map<std::string, path::Capability>>& lookup_tables,\
            const std::vector<path::Capability>& caps_to_load)

        Creates a pointer of :cpp:class:`RootSchemaNode<RootSchemaNode>` based on the list of capabilities passed in and the list of lookup tables passed in.

        :param lookup_tables: List of capabilities lookup table available.
        :param caps_to_load: List of capabilities available.
        :return: Pointer to the :cpp:class:`RootSchemaNode<RootSchemaNode>`.

    .. cpp:function:: void add_model_provider(ModelProvider* model_provider)

        The model provider is responsible for downloading the model from a remote device. It is meant to be invoked if the repository cannot load the model from the local directory.

        :param model_provider: The model provider to add

    .. cpp:function:: void remove_model_provider(ModelProvider* model_provider)

        Removes the given model provider from this repository

        :param model_provider: The model provider to remove

    .. cpp:function:: std::vector<ModelProvider*> get_model_providers() const

        Gets all model providers registered with this repository.

        :return: List of model providers associated with this repository
