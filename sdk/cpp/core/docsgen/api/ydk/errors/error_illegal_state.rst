YIllegalStateError
=====================


.. cpp:class:: ydk::YIllegalStateError : public YError

    Thrown when an operation/service is invoked on an object that is not in the right state. Use the ``msg`` for the error.

    .. cpp:function:: YIllegalStateError(const std::string& msg)
