#  ----------------------------------------------------------------
# Copyright 2017 Cisco Systems
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

from ydk.providers import NetconfServiceProvider
from ydk.path.sessions import NetconfSession
from ydk.path import Repository


class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        hostname = getattr(cls, 'hostname', '127.0.0.1')
        username = getattr(cls, 'username', 'admin')
        password = getattr(cls, 'password', 'admin')
        port = getattr(cls, 'port', 12022)
        protocol = getattr(cls, 'protocol', 'ssh')
        on_demand = not getattr(cls, 'non_demand', True)
        common_cache = getattr(cls, "common_cache", False)
        cls.ncc = NetconfServiceProvider(hostname, username, password, port, protocol, on_demand, common_cache)

    @classmethod
    def tearDownClass(self):
        pass

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_get_session(self):
        session = self.ncc.get_session()
        self.assertEqual(session is not None, True)

    def test_get_encoding(self):
        encoding = self.ncc.get_encoding()
        self.assertEqual(encoding is not None, True)

    def test_get_capabilities(self):
        capabilities = self.ncc.get_capabilities()
        self.assertEqual(capabilities is not None, True)


if __name__ == '__main__':
    import sys
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityTest)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
