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
import logging
import pkgutil
import importlib

from ydk import models
from ydk.types import EncodingFormat
from ydk.errors import YPYServiceProviderError
from ydk.path import Capability as _Capability
from ydk.path import Repository as _Repository


_TRACE_LEVEL_NUM = 5
_USER_PROVIDED_REPO = "ydk-user-provider-repo"


class CodecServiceProvider(object):
    """Python CodecServiceProvider wrapper.
    """

    def __init__(self, encoding, repo=None):
        self.logger = logging.getLogger(__name__)
        self.m_encoding = encoding
        self._m_root_schema_table = {}
        if repo is None:
            self._user_provided_repo = False
            self._m_repo = None
        else:
            self._user_provided_repo = True
            self._m_repo = repo

    def initialize(self, bundle_name, models_path):
        if self._user_provided_repo:
            self._initialize_root_schema(_USER_PROVIDED_REPO, self._m_repo, bundle_name)

        if bundle_name in self._m_root_schema_table:
            return

        self.logger.log(_TRACE_LEVEL_NUM, "Creating repo in path {}".format(models_path))
        repo = _Repository(models_path)
        self._initialize_root_schema(bundle_name, repo)

    def get_root_schema(self, bundle_name):
        if self._user_provided_repo:
            return self._m_root_schema_table[_USER_PROVIDED_REPO]

        if bundle_name not in self._m_root_schema_table:
            self.logger.error("Root schema not created")
            raise YPYServiceProviderError("Root schema not created")

        return self._m_root_schema_table[bundle_name]

    def _initialize_root_schema(self, repo_name, repo, bundle_name=''):
        bundle_name = repo_name if bundle_name == '' else bundle_name
        self.logger.log(_TRACE_LEVEL_NUM, "Initializing root schema for {}".format(repo_name))
        # TODO: turn on and off libyang logging
        capabilities = self._get_bundle_capabilities(bundle_name)
        self._m_root_schema_table[repo_name] = repo.create_root_schema(capabilities)

    def _get_bundle_capabilities(self, bundle_name):
        capabilities = []
        capability_map = {}
        for (_, name, ispkg) in pkgutil.iter_modules(models.__path__):
            if ispkg and name == bundle_name:
                try:
                    mod_yang_ns = importlib.import_module('ydk.models.{}._yang_ns'.format(name))
                    capability_map = mod_yang_ns.__dict__['CAPABILITIES']
                except ImportError:
                    continue

        for name in capability_map:
            capabilities.append(_Capability(name, capability_map[name]))

        return capabilities
