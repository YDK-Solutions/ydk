YError
=========

.. cpp:class:: ydk::YError : public std::exception

    Base class for YDK Exceptions.

    .. cpp:member:: std::string err_msg

    .. cpp:function:: YError(const std::string& msg)

    .. cpp:function:: char* what() const noexcept
