#  ----------------------------------------------------------------
# Copyright 2016 Cisco Systems
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------------
import contextlib
from ydk.errors import YPYError as _YPYError
from ydk.errors import YPYClientError as _YPYClientError
from ydk.errors import YPYIllegalStateError as _YPYIllegalStateError
from ydk.errors import YPYInvalidArgumentError as _YPYInvalidArgumentError
from ydk.errors import YPYModelError as _YPYModelError
from ydk.errors import YPYOperationNotSupportedError as _YPYOperationNotSupportedError
from ydk.errors import YPYServiceError as _YPYServiceError
from ydk.errors import YPYServiceProviderError as _YPYServiceProviderError


_ERRORS = { "YCPPError": _YPYError,
            "YCPPClientError": _YPYClientError,
            "YCPPIllegalStateError": _YPYIllegalStateError,
            "YCPPInvalidArgumentError": _YPYInvalidArgumentError,
            "YCPPModelError": _YPYModelError,
            "YCPPOperationNotSupportedError": _YPYOperationNotSupportedError,
            "YCPPServiceError": _YPYServiceError,
            "YCPPServiceProviderError": _YPYServiceProviderError,
}


@contextlib.contextmanager
def handle_runtime_error():
    try:
        yield
    except RuntimeError, err:
        etype_str, msg = err.message.split(':', 1)
        etype = _ERRORS.get(etype_str)
        raise etype(msg)


@contextlib.contextmanager
def handle_type_error():
    """Rethrow TypeError as YPYModelError"""
    try:
        yield
    except TypeError, err:
        raise _YPYModelError(err.message)
