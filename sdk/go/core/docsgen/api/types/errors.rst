.. _types-errors:
Errors
======

.. go:package:: ydk/types
    :synopsis: YDK Go Errors

.. attribute:: YGO_ERROR_TYPE

    Represents YDK Go error types with underlying type ``int``

    .. attribute YGO_ERROR_TYPE_NONE

        Represented by a value of 0

    .. attribute YGO_ERROR_TYPE_ERROR

        Represented by a value of 1

    .. attribute YGO_ERROR_TYPE_CLIENT_ERROR

        Represented by a value of 2

    .. attribute YGO_ERROR_TYPE_SERVICE_PROVIDER_ERROR

        Represented by a value of 3

    .. attribute YGO_ERROR_TYPE_SERVICE_ERROR

        Represented by a value of 4

    .. attribute YGO_ERROR_TYPE_ILLEGAL_STATE_ERROR

        Represented by a value of 5

    .. attribute YGO_ERROR_TYPE_INVALID_ARGUMENT_ERROR

        Represented by a value of 6

    .. attribute YGO_ERROR_TYPE_OPERATION_NOTSUPPORTED_ERROR

        Represented by a value of 7

    .. attribute YGO_ERROR_TYPE_MODEL_ERROR

        Represented by a value of 8

.. go:struct:: State

    .. attribute:: Private

        Type is interface{}

.. go:struct:: CState

    .. attribute:: Private

        Type is interface{}

.. interface -- is this necessary?
.. go:struct:: CError

    TODO

.. go:struct:: YGOError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
    :rtype: A Go ``string``

.. go:struct:: YGOClientError

    .. attribute:: Msg

        Represents the error message. Type is ``string``

.. function:: (e *YGOClientError) Error()

    Satisfies the error interface

    :return: Msg of given YGOError
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
