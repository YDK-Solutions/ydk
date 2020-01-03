ValidationService
=================


.. cpp:class:: ydk::path::ValidationService

    Instances of this class can validate the DataNode Tree based on the option supplied.

    .. cpp:function:: validate(const ydk::path::DataNode& dn, ydk::ValidationService::Option option)

        Validates data node based on the option.

        :param provider: An instance of :cpp:class:`DataNode<ydk::path::DataNode>`.
        :param option: An instance of type :cpp:class:`Option<Option>`.
        :raises: :cpp:class:`YModelError<YModelError>`, if validation error was detected.
