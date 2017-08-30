YDK Exceptions
==============

.. go:package:: ydk/errors
    :synopsis: YDK Exceptions

This package contains YDK Go exceptions classes.

.. go:exception:: YPYError

    Bases: :go:exc:`exceptions/Exception`

    Base class for YPY Errors. The subclasses give a specialized view of the error that has occurred.

.. go:exception:: YPYClientError

    Bases: :go:exc:`ydk/errors/YPYError`

    Exception for Client Side Validation

.. go:exception:: YPYIllegalStateError

    Bases: :go:exc:`ydk/errors/YPYError`

    Illegal State Error. Thrown when an operation/service is invoked on an object that is not in the right state. Use the error_msg for the error.

.. go:exception:: YPYInvalidArgumentError

    Bases: :go:exc:`ydk/errors/YPYError`

    Invalid Argument. Use the error_msg for the error.

.. go:exception:: YPYModelError

    Bases: :go:exc:`ydk/errors/YPYError`

    Model Error. Thrown when a model constraint is violated.

.. go:exception:: YPYOperationNotSupportedError

    Bases: :go:exc:`ydk/errors/YPYError`

    Operation Not Supported Error. Thrown when an operation is not supported.

.. go:exception:: YPYServiceError

    Bases: :go:exc:`ydk/errors/YPYError`

    Exception for Service Side Validation

.. go:exception:: YPYServiceProviderError

    Bases: :go:exc:`ydk/errors/YPYError`

    Exception for Provider Side Validation

