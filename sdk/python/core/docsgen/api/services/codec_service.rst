Codec service
=============

.. module:: ydk.services
    :synopsis: YDK Codec service

YDK Codec service provides encode/decode functionality.

.. py:class:: CodecService()

    Codec Service class for supporting encoding and decoding Python model API objects of type :py:class:`Entity<ydk.types.Entity>`.

    .. py:method:: encode(provider, entity)

        Encodes :py:class:`Entity<ydk.types.Entity>` in ``entity`` and returns the payload.

        :param provider: (:py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`) An instance of :py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`.
        :param entity: (:py:class:`Entity<ydk.types.Entity>` or ``dict`` of :py:class:`Entity<ydk.types.Entity>`)Could either be an instance of :py:class:`Entity<ydk.types.Entity>` or a dictionary of :py:class:`Entity<ydk.types.Entity>`.
        :return: Could be payload of dictionary of payload.
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if error has occurred.

    .. py:method:: decode(provider, payload)

        Decodes payload in ``payload`` and returns :py:class:`Entity<ydk.types.Entity>`.

        :param provider: (:py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`) An instance of :py:class:`CodecServiceProvider<ydk.providers.CodecServiceProvider>`.
        :param payload: (``str`` or ``dict`` of ``str``) Could be payload of a dictionary of payload.
        :return: Could be decoded :py:class:`Entity<ydk.types.Entity>`, or a dictionary of decoded :py:class:`Entity<ydk.types.Entity>`.
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` if error has occurred.
