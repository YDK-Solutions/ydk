YError
=========

.. cpp:class:: ydk::YError : public std::exception

    Base class for YDK exceptions.

    .. cpp:member:: std::string err_msg

    .. cpp:function:: YError(const std::string& msg)

    Class constructor

    .. cpp:function:: char* what() const noexcept

    Function to access error message

.. cpp:class:: ydk::YClientError : public ydk::YError

    YDK exception thrown when error occurred in protocol client application.

.. cpp:class:: ydk::YServiceProviderError : public ydk::YError

    YDK exception thrown when service provider errors occurred.

.. cpp:class:: ydk::YServiceError : public ydk::YError

    YDK exception thrown when service errors occurred.

.. cpp:class:: ydk::YIllegalStateError : public ydk::YError

    YDK exception thrown when an operation/service is invoked on an object that is not in the right state.

.. cpp:class:: ydk::YInvalidArgumentError : public ydk::YError

    YDK exception thrown when a function given parameter(s) with wrong values.

.. cpp:class:: ydk::YOperationNotSupportedError : public ydk::YError

    YDK exception thrown when specified yfilter is not supported by protocol.

.. cpp:class:: ydk::YModelError : public ydk::YError

    YDK exception thrown when a model constraint is violated.
