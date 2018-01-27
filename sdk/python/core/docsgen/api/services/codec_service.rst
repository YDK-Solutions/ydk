Codec Service
=============


.. py:class:: ydk.services.CodecService()

    Supports encoding and decoding Python model API objects of type :py:class:`Entity<ydk.types.Entity>`.

    .. py:method:: encode(provider, entity, pretty=True, subtree=False)

        Encodes :py:class:`Entity<ydk.types.Entity>` in ``entity`` and returns the payload.

        :param provider: (:py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`) Provider instance.
        :param entity: (:py:class:`Entity<ydk.types.Entity>` or dict of ``str`` and :py:class:`Entity<ydk.types.Entity>`) Either a single entity or a dictionary of entities.
        :param pretty: (:py:class:`EncodingFormat<ydk.types.EncodingFormat>`)Formats the payload in a readable way with idnentation
        :param subtree: (``bool``) Can be used to encode XML subtree filters to use with netconf ``get``/``get-config`` operations
        :return: Either a single ``str`` for a single encoded payload or a dictionary of ``str``
        :raises: :py:exc:`YError<ydk.errors.YError>` if error has occurred.

    .. py:method:: decode(provider, payload, subtree=False)

        Decodes payload in ``payload`` and returns :py:class:`Entity<ydk.types.Entity>`.

        :param provider: (:py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`) Provider instance.
        :param payload: (``str`` or dict of ``str`` and ``str``) Either a single encoded payload or a dictionary of payloads.
        :param subtree: (``bool``) Can be used to decode XML subtree filters to use with netconf ``get``/``get-config`` operations
        :return: Either a single decoded instance of :py:class:`Entity<ydk.types.Entity>` or a dictionary of decoded :py:class:`Entity<ydk.types.Entity>`
        :raises: :py:exc:`YError<ydk.errors.YError>` if error has occurred.
