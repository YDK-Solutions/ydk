.. _types-errors:
Errors
======

.. go:package:: ydk/types
    :synopsis: YDK Go Errors

.. code-block:: sh

   import "github.com/CiscoDevNet/ydk-go/ydk/types"

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

.. go:struct:: YGOServiceProviderError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOServiceProviderError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOServiceError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOServiceError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOIllegalStateError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOIllegalStateError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOInvalidArgumentError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOInvalidArgumentError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOOperationNotSupportedError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOOperationNotSupportedError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOModelError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOModelError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOCoreError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOCoreError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOCodecError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOCodecError) Error()
    
    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``
