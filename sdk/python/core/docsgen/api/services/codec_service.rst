Codec Service
=============

.. module:: ydk.services
    :synopsis: YDK Codec service

.. py:class:: CodecService()

    Supports encoding and decoding Python model API objects of type :py:class:`Entity<ydk.types.Entity>`.

    .. py:method:: encode(provider, entity, pretty=True, subtree=False)

        Encodes :py:class:`Entity<ydk.types.Entity>` in ``entity`` and returns the payload.

        :param provider: An instance of :py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`.
        :param entity: Either an instance of :py:class:`Entity<ydk.types.Entity>` or a dictionary of :py:class:`Entity<ydk.types.Entity>`
        :param pretty: Formats the payload in a readable way with idnentation
        :param subtree: Can be used to encode XML subtree filters to use with netconf ``get``/``get-config`` operations
        :return: Either a single ``str`` for a single encoded payload or a dictionary of ``str``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if error has occurred.

    .. py:method:: decode(provider, payload, subtree=False)

        Decodes payload in ``payload`` and returns :py:class:`Entity<ydk.types.Entity>`.

        :param provider: An instance of :py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`.
        :param payload: Either a single ``str`` for a single encoded payload or a dictionary of ``str``
        :param subtree: Can be used to decode XML subtree filters to use with netconf ``get``/``get-config`` operations
        :return: Either a single decoded instance of :py:class:`Entity<ydk.types.Entity>` or a dictionary of decoded :py:class:`Entity<ydk.types.Entity>`
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if error has occurred.
