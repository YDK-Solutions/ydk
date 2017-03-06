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

import os
import json
import logging
import pkgutil
import importlib
import xml.etree.ElementTree

from ydk.entity_utils import get_data_node_from_entity as _get_data_node_from_entity
from ydk.entity_utils import get_entity_from_data_node as _get_entity_from_data_node
from ydk.path import CodecService as _CodecService
from ydk.path import Capability as _Capability
from ydk.errors import YPYServiceProviderError as _YPYServiceProviderError
from ydk.errors.error_handler import handle_runtime_error as _handle_error
from ydk.errors.error_handler import handle_import_error as _handle_import_error
from ydk.errors.error_handler import check_argument as _check_argument
from ydk.types import EncodingFormat


_TRACE_LEVEL_NUM = 5
_ENTITY_ERROR_MSG = "No local YDK object install for {}"
_REPO_ERROR_MSG = "Failed to initialize provider."
_PAYLOAD_ERROR_MSG = "Codec service only supports one entity per payload, please split payload"


class CodecService(object):

    def __init__(self):
        self.logger = logging.getLogger(__name__)

    @_check_argument
    def encode(self, provider, entity_holder, pretty=True):
        if isinstance(entity_holder, dict):
            payload_map = {}
            for key in entity_holder:
                payload_map[key] = self._encode(provider, entity_holder[key], pretty)
            return payload_map
        else:
            return self._encode(provider, entity_holder, pretty)

    def _encode(self, provider, entity, pretty):
        bundle_name = _get_bundle_name(entity)
        provider.initialize(bundle_name, _get_yang_path(entity))
        root_schema = provider.get_root_schema(bundle_name)

        with _handle_error():
            data_node = _get_data_node_from_entity(entity, root_schema)
            codec_service = _CodecService()
            result = codec_service.encode(data_node, provider.m_encoding, pretty)
            self.logger.debug("Performing encode operation, resulting in {}", result)
            return result

        return ''

    @_check_argument
    def decode(self, provider, payload_holder):
        if isinstance(payload_holder, dict):
            entities = {}
            for key in payload_holder:
                entity = self.decode(provider, payload_holder[key])
                entities[key] = entity
            return entities
        else:
            return self._decode(provider, payload_holder)

    def _decode(self, provider, payload):
        entity = _get_top_entity(payload, provider.m_encoding, self.logger)

        bundle_name = _get_bundle_name(entity)
        provider.initialize(bundle_name, _get_yang_path(entity))

        root_schema = provider.get_root_schema(bundle_name)

        self.logger.debug("Performing decode operation on {}", payload)

        codec_service = _CodecService()
        root_data_node = codec_service.decode(root_schema, payload, provider.m_encoding)

        if len(root_data_node.children()) != 1:
            self.logger.debug(_PAYLOAD_ERROR_MSG)
            raise _YPYServiceProviderError(_PAYLOAD_ERROR_MSG)
        else:
            for data_node in root_data_node.children():
                _get_entity_from_data_node(data_node, entity)
        return entity


def _get_string(string):
    # TODO: check in python3 env
    if isinstance(string, unicode):
        return string.encode('utf-8')
    else:
        return string


def _get_ns_ename(payload, encoding):
    ns, ename = None, None
    if encoding == EncodingFormat.XML:
        payload_root = xml.etree.ElementTree.fromstring(payload)
        ns, ename = payload_root.tag.rsplit('}')
        ns = ns.strip('{')
    else:
        ns, ename = json.loads(payload).keys()[0].split(':')
        ns = _get_string(ns)
        ename = _get_string(ename)

    return (ns, ename)

def _get_top_entity(payload, encoding, logger):
    top_entity = None
    ns_ename = _get_ns_ename(payload, encoding)
    ydk_models = importlib.import_module('ydk.models')
    for (_, name, ispkg) in pkgutil.iter_modules(ydk_models.__path__):
        if ispkg:
            yang_ns = importlib.import_module('ydk.models.{}._yang_ns'.format(name))
            entity_lookup = yang_ns.__dict__['ENTITY_LOOKUP']
            if ns_ename in entity_lookup:
                return entity_lookup[ns_ename].clone_ptr()

    logger.debug(_ENTITY_ERROR_MSG, ename)
    raise _YPYServiceProviderError(_ENTITY_ERROR_MSG.format(ename))


def _get_yang_path(entity):
    m = entity.__module__.rsplit('.', 1)[0]
    m = importlib.import_module(m)
    return os.path.join(m.__path__[0], '_yang')


def _get_bundle_name(entity):
    m = entity.__module__.rsplit('.', 1)[0]
    m = importlib.import_module('.'.join([m, '_yang_ns']))
    return m.__dict__['BUNDLE_NAME']
