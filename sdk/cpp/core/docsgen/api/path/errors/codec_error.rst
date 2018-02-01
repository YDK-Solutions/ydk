YCodecError
==============


.. cpp:class:: ydk::path::YCodecError : public YCoreError

    Exception that encapsulates the validation errors for YDK :cpp:class:`CodecService`.

    .. cpp:member: Error err

    .. cpp:function:: YCodecError(YCodecError::Error merror)

    .. cpp:enum-class:: Error

        .. cpp:enumerator:: SUCCESS

            No error.

        .. cpp:enumerator:: XML_MISS

            Missing XML object.

        .. cpp:enumerator:: XML_INVAL

            Invalid XML object.

        .. cpp:enumerator:: XML_INCHAR

            Invalid XML character.

        .. cpp:enumerator:: EOF_ERR

            Unexpected end of input data.
