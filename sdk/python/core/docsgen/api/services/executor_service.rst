Executor Service
================

.. module:: ydk.services
    :synopsis: YDK Executor service

YDK Executor service provides functionality to execute RPCs

.. py:class:: ExecutorService

        Bases: :class:`ydk.services.Service`

        Executor Service class for supporting execution of RPCs.

        .. py:method:: execute_rpc(self, provider, rpc_entity, top_entity=None):

                Create the entity

                :param provider: An instance of ydk.providers.ServiceProvider
                :param rpc_entity: An instance of an RPC class defined under the ydk.models package or subpackages
                :param top_entity: Optional arg that should be provided when expecting data to be returned

                :return: None

                :raises ydk.errors.YPYDataValidationError: if validation error has occurred
                :raises ydk.errors.YPYError: if other error has occurred

                Possible Errors:

                * a server side error
                * there isn't enough information in the entity to prepare the message (eg. missing keys)