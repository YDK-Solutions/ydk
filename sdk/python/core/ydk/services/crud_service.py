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

from ydk.ext.services import CrudService as _CrudService
from ydk.errors.error_handler import handle_runtime_error as _handle_error


class CrudService(_CrudService):
    """ Python wrapper for CrudService
    """
    def __init__(self):
        self._crud = _CrudService()

    def create(self, provider, entity):
        with _handle_error():
            return self._crud.create(provider, entity)

    def read(self, provider, filter):
        with _handle_error():
            return self._crud.read(provider, filter)

    def read_config(self, provider, filter):
        with _handle_error():
            return self._crud.read_config(provider, filter)

    def update(self, provider, entity):
        with _handle_error():
            return self._crud.update(provider, entity)

    def delete(self, provider, entity):
        with _handle_error():
            return self._crud.delete(provider, entity)
