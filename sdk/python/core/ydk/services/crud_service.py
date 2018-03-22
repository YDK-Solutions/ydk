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
from ydk.ext.services import CRUDService as _CrudService
from ydk.errors.error_handler import handle_runtime_error as _handle_error
from ydk.errors.error_handler import check_argument as _check_argument

from ydk.types import EntityCollection, Config

class CRUDService(_CrudService):
    """
    Python wrapper for CrudService

    Provides API to Create, Read, Update, and Delete configuration on device.

    Args:
        provider (ydk.ServiceProvider): NetconfServiceProvider.
        entity_holder (ydk.types.Entity or a list(ydk.types.Entity)): Encoding target(s).
        pretty (bool, optional): Pretty formatting. Default value is True.

    Returns:
        Functions read and read-config return, ydk.types.Entity or a list(ydk.types.Entity), depending on provided input.
        Other functions return boolean value: true - for success; false - failure.

    Raises:
        Instance of YPYError in case of operation failure fails.
    """
    def __init__(self):
        self._crud = _CrudService()

    @_check_argument
    def create(self, provider, entity):
        if isinstance(entity, EntityCollection):
            entity = entity.get_entities()
        with _handle_error():
            return self._crud.create(provider, entity)

    @_check_argument
    def read(self, provider, read_filter):
        filters = read_filter
        if isinstance(read_filter, EntityCollection):
            filters = read_filter.get_entities()
        with _handle_error():
            read_entity = self._crud.read(provider, filters)
        if isinstance(read_filter, EntityCollection):
            read_entity = Config(read_entity)
        return read_entity

    @_check_argument
    def read_config(self, provider, read_filter):
        filters = read_filter
        if isinstance(read_filter, EntityCollection):
            filters = read_filter.get_entities()
        with _handle_error():
            read_entity = self._crud.read_config(provider, filters)
        if isinstance(read_filter, EntityCollection):
            read_entity = Config(read_entity)
        return read_entity

    @_check_argument
    def update(self, provider, entity):
        if isinstance(entity, EntityCollection):
            entity = entity.get_entities()
        with _handle_error():
            return self._crud.update(provider, entity)

    @_check_argument
    def delete(self, provider, entity):
        if isinstance(entity, EntityCollection):
            entity = entity.get_entities()
        with _handle_error():
            return self._crud.delete(provider, entity)
