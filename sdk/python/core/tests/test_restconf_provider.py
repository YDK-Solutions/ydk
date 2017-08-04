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

"""test_restconf_provider.py
RestconfServiceProvider test
"""
from __future__ import absolute_import

import os
import unittest

from ydk.providers import RestconfServiceProvider
from ydk.types import EncodingFormat
from ydk.path import Repository


class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        # Need to keep a local reference for repo to keep it alive
        repo_path = os.path.dirname(__file__)
        repo_path = os.path.join(repo_path, '..', '..', '..', 'cpp', 'core', 'tests', 'models')
        self.repo = Repository(repo_path)
        self.restconf_provider = RestconfServiceProvider(self.repo, 'localhost', 'admin', 'admin', 12306, EncodingFormat.JSON)

    def test_get_session(self):
        session = self.restconf_provider.get_session()
        self.assertEqual(session is not None, True)

    def test_get_encoding(self):
        encoding = self.restconf_provider.get_encoding()
        self.assertEqual(encoding is not None, True)


if __name__ == '__main__':
    import sys
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityTest)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
