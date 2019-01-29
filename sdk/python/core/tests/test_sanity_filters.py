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

"""test_sanity_levels.py
sanity test for ydktest-sanity.yang
"""
from __future__ import absolute_import

import sys
import unittest

from ydk.models.ydktest import ydktest_sanity as ysanity
from ydk.providers import NetconfServiceProvider
from ydk.services import CRUDService
from ydk.filters import YFilter

from test_utils import ParametrizedTestCase
from test_utils import get_device_info


class SanityYang(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.ncc = NetconfServiceProvider(
            cls.hostname,
            cls.username,
            cls.password,
            cls.port,
            cls.protocol,
            cls.on_demand,
            cls.common_cache,
            cls.timeout)
        cls.crud = CRUDService()

    def setUp(self):
        runner = ysanity.Runner()
        self.crud.delete(self.ncc, runner)

    def tearDown(self):
        runner = ysanity.Runner()
        self.crud.delete(self.ncc, runner)

    def test_read_on_ref_class(self):
        one = ysanity.Runner.YdktestSanityOne()
        one.number, one.name = 1, 'runner:one:name'
        self.crud.create(self.ncc, one)

        one_filter = ysanity.Runner.YdktestSanityOne()
        one_read = self.crud.read(self.ncc, one_filter)
        self.assertIsNotNone(one_read)
        self.assertEqual(one, one_read)

    def test_read_on_leaf(self):
        one = ysanity.Runner.YdktestSanityOne()
        one.number, one.name = 1, 'runner:one:name'
        self.crud.create(self.ncc, one)

        one_filter = ysanity.Runner.YdktestSanityOne()
        one_filter.number = YFilter.read
        one_read = self.crud.read(self.ncc, one_filter)
        self.assertIsNotNone(one_read)
        self.assertEqual(one, one_read)

        # no such value, will return empty data
        one_filter = ysanity.Runner.YdktestSanityOne()
        one_filter.number = 2
        one_read = self.crud.read(self.ncc, one_filter)
        self.assertIsNone(one_read)

    def test_read_on_ref_enum_class(self):
        from ydk.models.ydktest.ydktest_sanity import YdkEnumTest
        r_1 = ysanity.Runner.Ytypes.BuiltInT()
        r_1.enum_value = YdkEnumTest.local
        self.crud.create(self.ncc, r_1)

        r_2 = ysanity.Runner().Ytypes.BuiltInT()
        r_2.enum_value = YFilter.read
        runner_read = self.crud.read(self.ncc, r_2)
        self.assertIsNotNone(runner_read)
        self.assertEqual(r_1.enum_value, runner_read.enum_value)

        r_2 = ysanity.Runner().Ytypes.BuiltInT()
        r_2.enum_value = YdkEnumTest.local
        runner_read = self.crud.read(self.ncc, r_2)
        self.assertEqual(r_1.enum_value, runner_read.enum_value)

        # no such value, nothing returned
        r_2 = ysanity.Runner.Ytypes.BuiltInT()
        r_2.enum_value = YdkEnumTest.remote
        runner_read = self.crud.read(self.ncc, r_2)
        self.assertIsNone(runner_read)

    def test_read_on_ref_list(self):
        r_1 = ysanity.Runner()
        l_1, l_2 = ysanity.Runner.OneList.Ldata(), ysanity.Runner.OneList.Ldata()
        l_1.number, l_2.number = 1, 2
        r_1.one_list.ldata.extend([l_1, l_2])
        self.crud.create(self.ncc, r_1)

        r_2 = ysanity.Runner()
        r_2.one_list.ldata.yfilter = YFilter.read
        runner_read = self.crud.read(self.ncc, r_2)

        self.assertEqual(runner_read, r_1)

    def test_read_on_list_with_key(self):
        r_1 = ysanity.Runner.OneList()
        l_1, l_2 = ysanity.Runner.OneList.Ldata(), ysanity.Runner.OneList.Ldata()
        l_1.number, l_2.number = 1, 2
        r_1.ldata.extend([l_1, l_2])
        self.crud.create(self.ncc, r_1)

        r_2 = ysanity.Runner().OneList()
        r_2.ldata.extend([l_1])
        runner_read = self.crud.read(self.ncc, r_2)

        self.assertEqual(runner_read, r_2)

    def test_read_on_leaflist(self):
        r_1 = ysanity.Runner.Ytypes.BuiltInT()
        r_1.llstring.extend(['1', '2', '3'])
        self.crud.create(self.ncc, r_1)

        r_2 = ysanity.Runner()
        r_2.ytypes.built_in_t.llstring = YFilter.read
        runner_read = self.crud.read(self.ncc, r_2)
        self.assertEqual(runner_read.ytypes.built_in_t.llstring, r_1.llstring)

    def test_read_on_identity_ref(self):
        r_1 = ysanity.Runner.Ytypes.BuiltInT()
        r_1.identity_ref_value = ysanity.ChildIdentity()
        self.crud.create(self.ncc, r_1)

        r_2 = ysanity.Runner().Ytypes.BuiltInT()
        r_2.identity_ref_value = YFilter.read
        t_read = self.crud.read(self.ncc, r_2)
        self.assertIsNotNone(t_read)
        self.assertEqual(r_1.identity_ref_value, t_read.identity_ref_value)

    def test_read_only_config(self):
        r_1 = ysanity.Runner()
        r_1.ydktest_sanity_one.number, r_1.ydktest_sanity_one.name = 1, 'runner:one:name'
        self.crud.create(self.ncc, r_1)
        r_2, r_3 = ysanity.Runner(), ysanity.Runner()
        r_2.ydktest_sanity_one.number = YFilter.read
        r_3.ydktest_sanity_one.number = YFilter.read

        r_2 = self.crud.read_config(self.ncc, r_2)
        r_3 = self.crud.read(self.ncc, r_3)
        # ysanity only have config data, ok to compare
        self.assertEqual(r_2.ydktest_sanity_one.number, r_3.ydktest_sanity_one.number)
        self.assertEqual(r_2.ydktest_sanity_one.name, r_3.ydktest_sanity_one.name)

    def test_decoder(self):
        # send payload to device
        element = ysanity.Runner.OneList.Ldata()
        element.number = 5
        element.name = 'five'

        self.crud.create(self.ncc, element)

        runner_filter = ysanity.Runner().OneList.Ldata()
        runner_filter.number = 5
        element_read = self.crud.read(self.ncc, runner_filter)
        self.assertEqual(element, element_read)


if __name__ == '__main__':
    device, non_demand, common_cache, timeout = get_device_info()

    suite = unittest.TestSuite()
    suite.addTest(ParametrizedTestCase.parametrize(
        SanityYang,
        device=device,
        non_demand=non_demand,
        common_cache=common_cache,
        timeout=timeout))
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
