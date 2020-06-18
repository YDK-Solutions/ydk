Codec Service Provider
======================


.. py:class:: ydk.providers.CodecServiceProvider(type=EncodingFormat.XML, repo=None)

    A provider to be used with :py:class:`CodecService<ydk.services.CodecService>` for performing encoding and decoding.

    :param type: An argument specifies encoding format, could be a Python string (``xml`` or ``json``) or an instance of :py:class:`EncodingFormat<ydk.types.EncodingFormat>`.
    :type type: ``string`` or :py:class:`EncodingFormat<ydk.types.EncodingFormat>`
    :param repo: User provided repository stores cached models.
    :type repo: :py:class:`Repository<ydk.path.Repository>`

    .. py:method:: get_root_schema(bundle_name)

        Return root schema node for bundle_name defined in the 'initialize' method.

        :param bundle_name: (``str``) Bundle name.
        :return: :py:class:`RootSchemaNode<ydk.path.RootSchemaNode>` for this bundle.

    .. py:method:: initialize(bundle_name, models_path)

        Initializes root schema in the user defined location of YANG models.

        :param bundle_name: (``str``): user defined bundle name.
        :param models_path: (``str``): location (directory) for local YANG models.
