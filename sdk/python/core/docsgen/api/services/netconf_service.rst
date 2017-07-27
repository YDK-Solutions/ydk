NETCONF service
===============

.. module:: ydk.services
    :synopsis: NETCONF service

    NETCONF service provides NETCONF execution functionalities.

.. py:class:: NetconfService
\
    .. py:method:: cancel_commit(session, persist_id=-1)

        Cancels an ongoing confirmed commit. If the ``persist_id`` parameter is not given, the operation **MUST** be issued on the same session that issued the confirmed commit.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param persist_id: An ``int`` that cancels a persistent confirmed commit
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: close_session(session)

        Request graceful termination of a NETCONF session

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: commit(session, confirmed=False, confirm_timeout=-1, persist=-1, persist_id=-1)

        Instructs the device to implement the configuration data contained in the candidate configuration.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param confirmed: A ``bool`` that signals a confirmed commit operation
        :param confirm_timeout: An ``int`` representing the timeout interval for a confirmed commit
        :param persist: An ``int`` that makes the confirmed commit persistent
        :param persist_id: An ``int`` that is given in order to commit a persistent confirmed commit
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: copy_config(session, target, source, url='')

        Create or replace an entire configuration DataStore with the contents of another complete configuration DataStore. If the target DataStore exists, it is overwritten. Otherwise, a new one is created, if allowed.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param target: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being used as the destination
        :param source: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being used as the source
        :param url: A ``str`` representing the configuration url
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: delete_config(session, target, url="")

        Delete a configuration DataStore. The RUNNING configuration DataStore cannot be deleted

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param target: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration to be deleted
        :param url: Optional arg: A ``str`` representing the configuration url required only when target is set to ``url``
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: discard_changes(session)

        Used to revert the candidate configuration to the current running configuration.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: edit_config(session, target, config, default_operation='', error_option='', test_option='')

        Loads all or part of a specified configuration to the specified target configuration datastore. Allows the new configuration to be expressed using a local file, a remote file, or inline. If the target configuration datastore does not exist, it will be created.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param target: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being edited
        :param config: An instance of :py:class:`Entity<ydk.types.Entity>` that is a hierarchy configuration of data as defined by one of the device’s data models
        :param default_operation: A ``str`` that changes the default from ``merge`` to either ``merge``, ``replace``, or ``none``
        :param error_option: Optional arg: A ``str`` that can be set to ``test-then-set``, ``set``, or ``test-only`` if the device advertises the :validate:1.1 capability
        :param test_option: Optional arg: A ``str`` that can be set to ``stop-on-error``, ``continue-on-error``, or ``rollback-on-error``
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred


    .. py:method:: get_config(session, source, filter)

        Retrieve all or part of a specified configuration datastore

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param filter: An instance of :py:class:`Entity<ydk.types.Entity>`
        :param source: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration being queried
        :return: The requested data in :py:class:`Entity<ydk.types.Entity>` instance
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: get(session, filter)

        Retrieve running configuration and device state information.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param filter: An instance of :py:class:`Entity<ydk.types.Entity>` which specifies the portion of the system configuration and state data to retrieve
        :return: The requested data in :py:class:`Entity<ydk.types.Entity>` instance
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: kill_session(session, session_id)

        Force the termination of a NETCONF session.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param session_id: An ``int`` that is the session identifier of the NETCONF session to be terminated
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: lock(session, target)

        Allows the client to lock the entire configuration datastore system of a device.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param target: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration to lock
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: unlock(session, target)

        Used to release a configuration lock, previously obtained with the LOCK operation.

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param target: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration to unlock
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred

    .. py:method:: validate(session, source=None, url='', source_config=None)

        Execute a validate operation to validate the contents of the specified configuration

        :param session: An instance of :py:class:`NetconfSession<ydk.path.NetconfSession>`
        :param source: An instance of :py:class:`DataStore<ydk.services.DataStore>` representing the configuration datastore to validate
        :param url: A ``str`` representing the configuration url
        :param source_config: An instance of :py:class:`Entity<ydk.types.Entity>` representing the configuration to validate
        :return: ``True`` if the operation succeeds, else ``False``
        :raises: :py:exc:`YPYError<ydk.errors.YPYError>` If error has occurred
