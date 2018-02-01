Executor Service
================


.. cpp:class:: ydk::ExecutorService

    Executor Service class for supporting execution of RPCs.

    .. cpp:function:: ExecutorService()

    .. cpp:function:: std::shared_ptr<Entity> execute_rpc(ydk::ServiceProvider & provider, Entity & rpc_entity, std::shared_ptr<Entity> top_entity = nullptr)

        Create the rpc entity.

        :param provider: An instance of :cpp:class:`ServiceProvider<ServiceProvider>`.
        :param rpc_entity: An rpc instance of :cpp:class:`Entity<ydk::Entity>` class defined under a bundle.
        :param top_entity: Optional arg that should be provided when expecting data to be returned.
        :return: The requested data as :cpp:class:`Entity<ydk::Entity>` or nullptr if N/A.
        :raises YError: If an error has occurred.
