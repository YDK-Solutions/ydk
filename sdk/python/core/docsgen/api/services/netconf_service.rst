NETCONF service
===============

.. py:class:: ydk.services.NetconfService

    .. py:method:: cancel_commit(provider, persist_id=None)

        Cancels an ongoing confirmed commit. If the ``persist_id`` parameter is None, the operation **MUST** be issued on the same session that issued the confirmed commit.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param persist_id: (``int``) An ``int`` that cancels a persistent confirmed commit.
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: close_session(provider)

        Request graceful termination of a NETCONF session

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: commit(provider, confirmed=False, confirm_timeout=None, persist=None, persist_id=None)

        Instructs the device to implement the configuration data contained in the candidate configuration.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param confirmed: (``bool``) A ``bool`` that signals a confirmed commit operation
        :param confirm_timeout: (``int``) An ``int`` representing the timeout interval for a confirmed commit
        :param persist: (``int``) An ``int`` that makes the confirmed commit persistent
        :param persist_id: (``int``) An ``int`` that is given in order to commit a persistent confirmed commit
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: copy_config(provider, target, source, url='')

        Create or replace an entire configuration DataStore with the contents of another complete configuration DataStore. If the target DataStore exists, it is overwritten. Otherwise, a new one is created, if allowed.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param target: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being used as the destination
        :param source: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being used as the source
        :param url: (``str``) A ``str`` representing the configuration url
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: delete_config(provider, target, url="")

        Delete a configuration DataStore. The RUNNING configuration DataStore cannot be deleted

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param target: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration to be deleted
        :param url: (``str``, optional) A ``str`` representing the configuration url required only when target is set to ``url``
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: discard_changes(provider)

        Used to revert the candidate configuration to the current running configuration.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: edit_config(provider, target, config, default_operation='', error_option='', test_option='')

        Loads all or part of a specified configuration to the specified target configuration datastore. Allows the new configuration to be expressed using a local file, a remote file, or inline. If the target configuration datastore does not exist, it will be created.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param target: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being edited
        :param config: An instance of :py:class:`Entity<ydk.types.Entity>` that is a hierarchy configuration of data as defined by one of the device’s data models
        :param default_operation: (``str``) A ``str`` that changes the default from ``merge`` to either ``merge``, ``replace``, or ``none``
        :param error_option: (``str``, optional) A ``str`` that can be set to ``test-then-set``, ``set``, or ``test-only`` if the device advertises the :validate:1.1 capability
        :param test_option: (``str``, optional) A ``str`` that can be set to ``stop-on-error``, ``continue-on-error``, or ``rollback-on-error``
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred


    .. py:method:: get_config(provider, source, filter)

        Retrieve all or part of a specified configuration datastore

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param filter: (:py:class:`Entity<ydk.types.Entity>`) An instance of :py:class:`Entity<ydk.types.Entity>`
        :param source: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being queried
        :return: The requested data in :py:class:`Entity<ydk.types.Entity>` instance
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: get(provider, filter)

        Retrieve running configuration and device state information.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param filter: (:py:class:`Entity<ydk.types.Entity>`) An instance of :py:class:`Entity<ydk.types.Entity>` which specifies the portion of the system configuration and state data to retrieve
        :return: The requested data in :py:class:`Entity<ydk.types.Entity>` instance
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: kill_session(provider, session_id)

        Force the termination of a NETCONF session.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param session_id: (``int``) An ``int`` that is the session identifier of the NETCONF session to be terminated
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: lock(provider, target)

        Allows the client to lock the entire configuration datastore system of a device.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param target: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration to lock
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: unlock(provider, target)

        Used to release a configuration lock, previously obtained with the LOCK operation.

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param target: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration to unlock
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: validate(provider, source=None, url='', source_config=None)

        Execute a validate operation to validate the contents of the specified configuration

        :param provider: (:py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`) NETCONF provider instance.
        :param source: (:py:class:`DataStore<ydk.services.DataStore>`) An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration datastore to validate
        :param url: (``str``) A ``str`` representing the configuration url
        :param source_config: (:py:class:`Entity<ydk.types.Entity>`) An instance of :py:class:`Entity<ydk.types.Entity>` representing the configuration to validate
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred
