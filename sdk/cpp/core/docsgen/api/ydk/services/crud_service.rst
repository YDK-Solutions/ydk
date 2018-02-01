CRUD Service
============


.. cpp:class:: ydk::CrudService

    CRUD Service class for supporting CRUD operations on entities.

    .. cpp:function:: CrudService()

    .. cpp:function:: bool create(ydk::ServiceProvider & provider, Entity & entity)

        Create the entity.

        :param provider: An instance of :cpp:class:`ServiceProvider<ydk::ServiceProvider>`.
        :param entity: An instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :return: ``true`` if successful, ``false`` if not.
        :raises YError: If an error has occurred.

    .. cpp:function:: std::shared_ptr<ydk::Entity> read(ydk::ServiceProvider & provider, Entity & filter)

        Read the entity.

        :param provider: An instance of :cpp:class:`ServiceProvider<ydk::ServiceProvider>`.
        :param filter: An instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :return: A pointer to an instance of :cpp:class:`Entity<ydk::Entity>` as identified by the ``filter`` if successful, ``nullptr`` if not.
        :raises YError: If an error has occurred.

    .. cpp:function:: std::shared_ptr<ydk::Entity> read_config(ydk::ServiceProvider & provider, Entity & filter)

        Read only config.

        :param provider: An instance of :cpp:class:`ServiceProvider<ydk::ServiceProvider>`.
        :param filter: An instance of :cpp:class:`entity<ydk::Entity>` class defined under a bundle.
        :return: A pointer to an instance of :cpp:class:`Entity<ydk::Entity>` as identified by the ``filter`` if successful, ``nullptr`` if not.
        :raises YError: If an error has occurred.

    .. cpp:function:: bool update(ydk::ServiceProvider & provider, Entity & entity)

        Update the entity.

        :param provider: An instance of :cpp:class:`ServiceProvider<ydk::ServiceProvider>`.
        :param entity: An instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :return: ``true`` if successful, ``false`` if not.
        :raises YError: If an error has occurred.

    .. cpp:function:: bool delete\_(ydk::ServiceProvider & provider, Entity & entity)

        Delete the entity.

        :param provider: An instance of :cpp:class:`ServiceProvider<ydk::ServiceProvider>`.
        :param entity: An instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :return: ``true`` if successful, ``false`` if not.
        :raises YError: If an error has occurred.
