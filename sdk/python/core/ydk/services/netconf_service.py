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
import inspect
from ydk.ext.services import DataStore, NetconfService as _NetconfService
from ydk.errors import YPYServiceError as _YPYServiceError
from ydk.errors.error_handler import handle_runtime_error as _handle_error


class NetconfService(_NetconfService):
    """ Python wrapper for NetconfService
    """
    def __init__(self):
        self._ns = _NetconfService()

    def cancel_commit(self, session, persist_id=-1):
        if None in (session, persist_id):
            raise _YPYServiceError("session and persist_id cannot be None")

        with _handle_error():
            return self._ns.cancel_commit(session, persist_id)

    def close_session(self, session):
        if session is None:
            raise _YPYServiceError("session cannot be None")

        with _handle_error():
            return self._ns.close_session(session)

    def commit(self, session, confirmed=False, confirm_timeout=-1, persist=-1, persist_id=-1):
        if session is None:
            raise _YPYServiceError("session cannot be None")

        with _handle_error():
            return self._ns.commit(session, confirmed, confirm_timeout, persist, persist_id)

    def copy_config(self, session, target, source=None, url="", source_config=None):
        if None in (session, target) or (source is None and source_config is None):
            raise _YPYServiceError("session, target, and source/source_config cannot be None")

        with _handle_error():
            if type(source) == DataStore:
                return self._ns.copy_config(session, target, source, url)
            elif source_config is not None:
                return self._ns.copy_config(session, target, source_config)
            else:
                return self._ns.copy_config(session, target, source)

    def delete_config(self, session, target, url=""):
        if None in (session, target):
            raise _YPYServiceError("session and target cannot be None")

        with _handle_error():
            return self._ns.delete_config(session, target, url)

    def discard_changes(self, session):
        if session is None:
            raise _YPYServiceError("session cannot be None")

        with _handle_error():
            return self._ns.discard_changes(session)

    def edit_config(self, session, target, config,
        default_operation="", test_option="", error_option=""):

        if None in (session, target, config):
            raise _YPYServiceError("session, target, and config cannot be None")

        with _handle_error():
            return self._ns.edit_config(session, target, config,
                default_operation, test_option, error_option)

    def get_config(self, session, source, filter):
        if None in (session, source, filter):
            raise _YPYServiceError("session, source, and filter cannot be None")

        with _handle_error():
            return self._ns.get_config(session, source, filter)

    def get(self, session, filter):
        if None in (session, filter):
            raise _YPYServiceError("session and filter cannot be None")

        with _handle_error():
            return self._ns.get(session, filter)

    def kill_session(self, session, session_id):
        if None in (session, session_id):
            raise _YPYServiceError("session and session_id cannot be None")

        with _handle_error():
            return self._ns.kill_session(session, session_id)

    def lock(self, session, target):
        if None in (session, target):
            raise _YPYServiceError("session and target cannot be None")

        with _handle_error():
            return self._ns.lock(session, target)

    def unlock(self, session, target):
        if None in (session, target):
            raise _YPYServiceError("session and target cannot be None")

        with _handle_error():
            return self._ns.unlock(session, target)

    def validate(self, session, source=None, url="", source_config=None):
        if session is None or (source is None and source_config is None):
            raise _YPYServiceError("session and source/source_config cannot be None")

        with _handle_error():
            if type(source) == DataStore:
                return self._ns.validate(session, source, url)
            elif source_config is not None:
                return self._ns.validate(session, source_config)
            else:
                return self._ns.validate(session, source)
