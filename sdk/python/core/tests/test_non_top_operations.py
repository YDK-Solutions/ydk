#  ----------------------------------------------------------------
# Copyright 2018 Cisco Systems
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

"""test_non_top_operations.py
Sanity test for ydktest-sanity.yang targetted specifically
to test support for non-top level objects in CRUD operations
"""

from __future__ import absolute_import
from __future__ import print_function

import unittest

from ydk.errors    import YModelError, YServiceError
from ydk.providers import NetconfServiceProvider
from ydk.services  import NetconfService
from ydk.services  import CRUDService
from ydk.filters import YFilter

try:
    from ydk.models.ydktest.ydktest_sanity import Runner, ChildIdentity
except ImportError:
    from ydk.models.ydktest.ydktest_sanity.runner.runner import Runner
    from ydk.models.ydktest.ydktest_sanity.ydktest_sanity import ChildIdentity

class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.ncc = NetconfServiceProvider(
            "127.0.0.1",
            "admin",
            "admin",
            12022,
            )
        self.crud = CRUDService()

    def setUp(self):
        runner = Runner()
        self.crud.delete(self.ncc, runner)

    def test_delete_on_list_with_identitykey(self):
        a1 = Runner.OneList.IdentityList()
        a1.config.id = ChildIdentity()
        a1.id_ref =  a1.config.id
        self.crud.create(self.ncc, a1)

        k = Runner.OneList.IdentityList()
        k.config.id = ChildIdentity()
        k.id_ref = k.config.id
        k.yfilter = YFilter.delete
        self.crud.update(self.ncc, k)

        runner_read = self.crud.read(self.ncc, Runner())
        self.assertIsNone(runner_read)

    def test_iden_list(self):
        # CREATE
        il = Runner.OneList.IdentityList()
        il.config.id = ChildIdentity()
        il.id_ref = ChildIdentity()
        self.crud.create(self.ncc, il)

        # READ & VALIDATE
        runner_filter = Runner.OneList()
        read_one = self.crud.read(self.ncc, runner_filter)
        self.assertIsNotNone(read_one)

        read_il = read_one.identity_list.get(ChildIdentity().to_string())
        self.assertIsNotNone(read_il)
        read_il.parent = None
        self.assertEqual(read_il, il)

        # DELETE & VALIDATE
        self.crud.delete(self.ncc, il)
        runner_read = self.crud.read(self.ncc, Runner())
        self.assertIsNone(runner_read)

if __name__ == '__main__':
    suite = unittest.TestSuite()
    testloader = unittest.TestLoader()
    testnames = testloader.getTestCaseNames(SanityTest)
    for name in testnames:
        suite.addTest(SanityTest(name))
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
