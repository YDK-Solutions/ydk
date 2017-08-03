Statement
=========


.. cpp:class:: ydk::path::Statement

    Represents the YANG :cpp:class:`Statement<Statement>`.

    .. cpp:member:: std::string keyword

        YANG keyword corresponding to the :cpp:class:`Statement<Statement>`.

    .. cpp:member:: std::string  arg

        The arg if any.

    .. cpp:member:: std::string name_sapce

        The namespace if any.

    .. cpp:function:: Statement()

    .. cpp:function:: Statement(const std::string& mkeyword, const std::string& marg)

    .. cpp:function:: Statement(const Statement& stmt)

    .. cpp:function:: Statement(Statement&& stmt)

    .. cpp:function:: ~Statement()

    .. cpp:function:: Statement& operator=(const Statement& stmt)

    .. cpp:function:: Statement& operator=(Statement&& stmt)
