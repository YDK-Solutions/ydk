YCPPError
=========

.. cpp:class:: YCPPError : public std::exception

  Base class for YDK Exceptions.

    .. cpp:member:: std::string err_msg

    .. cpp:function:: YCPPError(const std::string& msg)

    .. cpp:function:: char* what() const noexcept
