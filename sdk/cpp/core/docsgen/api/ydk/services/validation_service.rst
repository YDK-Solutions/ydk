.. _ref-validationservice:

ValidationService
=================


.. cpp:class:: ydk::ValidationService

    Validation Service class for validating C++ model API objects of type :cpp:class:`Entity<Entity>`.

    .. cpp:enum-class:: Option

        All validation is performed in the context of some operation. These options capture the context of use.

        .. cpp:enumerator:: DATASTORE

            Datastore validation.

        .. cpp:enumerator:: GET_CONFIG

            Get config validation. Checks to see if only config nodes are references.

        .. cpp:enumerator:: GET

            Get validation.

        .. cpp:enumerator:: EDIT_CONFIG

            Edit validation. Checks on the values of leafs etc.

    .. cpp:function:: void validate(const ydk::path::ServiceProvider& provider, Entity& entity, ValidationService::Option option)

        Validates an entity based on the option.

        :param provider: An instance of :cpp:class:`ServiceProvider<ydk::path::ServiceProvider>`
        :param entity: An instance of :cpp:class:`Entity<Entity>` class defined under a bundle
        :param option: An instance of type :cpp:class:`Option<Option>`
        :raises YError: If validation errors were detected

    .. cpp:function:: virtual ~ValidationService()
