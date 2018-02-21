.. _ydk-errors:

Errors
======

.. go:package:: ydk/errors
    :synopsis: YDK Go Errors

.. code-block:: sh

   import "github.com/CiscoDevNet/ydk-go/ydk/errors"

.. attribute:: Y_ERROR_TYPE

    Represents YDK Go error types with underlying type ``int``

    .. attribute Y_ERROR_TYPE_NONE

        Represented by a value of 0

    .. attribute Y_ERROR_TYPE_ERROR

        Represented by a value of 1

    .. attribute Y_ERROR_TYPE_CLIENT_ERROR

        Represented by a value of 2

    .. attribute Y_ERROR_TYPE_SERVICE_PROVIDER_ERROR

        Represented by a value of 3

    .. attribute Y_ERROR_TYPE_SERVICE_ERROR

        Represented by a value of 4

    .. attribute Y_ERROR_TYPE_ILLEGAL_STATE_ERROR

        Represented by a value of 5

    .. attribute Y_ERROR_TYPE_INVALID_ARGUMENT_ERROR

        Represented by a value of 6

    .. attribute Y_ERROR_TYPE_OPERATION_NOTSUPPORTED_ERROR

        Represented by a value of 7

    .. attribute Y_ERROR_TYPE_MODEL_ERROR

        Represented by a value of 8

.. go:struct:: State

    .. attribute:: Private

        Type is ``interface{}``

.. go:struct:: CState

    .. attribute:: Private

        Type is ``interface{}``

.. _errors-cerror:

.. object:: CError

    CError is an interface type that represents a basic error in Go.

    .. function:: Error()

        :rtype: ``string``

.. go:struct:: YError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YClientError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YClientError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YServiceProviderError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YServiceProviderError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YServiceError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YServiceError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YIllegalStateError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YIllegalStateError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YInvalidArgumentError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YInvalidArgumentError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YOperationNotSupportedError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YOperationNotSupportedError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YModelError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YModelError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YCoreError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YCoreError) Error()

    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``

.. go:struct:: YCodecError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YCodecError) Error()
    
    Satisfies the error interface

    :return: Msg of given YError
    :rtype: A Go ``string``
