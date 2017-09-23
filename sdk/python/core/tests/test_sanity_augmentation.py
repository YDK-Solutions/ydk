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
"""
    test_sanity_bundle_augmentation.py
        Unittest for bundle augmentation.
"""
from __future__ import absolute_import

import sys
import unittest

from ydk.services import CRUDService
from ydk.providers import NetconfServiceProvider
from ydk.models.augmentation import ietf_aug_base_1
from ydk.models.augmentation import ietf_aug_base_2

from test_utils import assert_with_error
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
        self.crud.delete(self.ncc, ietf_aug_base_1.Cpython())
        self.crud.delete(self.ncc, ietf_aug_base_2.Cpython())

    def tearDown(self):
        self.crud.delete(self.ncc, ietf_aug_base_1.Cpython())
        self.crud.delete(self.ncc, ietf_aug_base_2.Cpython())

    def test_aug_base_1(self):
        cpython = ietf_aug_base_1.Cpython()
        cpython.doc.ydktest_aug_1.aug_one = 'aug one'
        cpython.doc.ydktest_aug_2.aug_two = 'aug two'
        cpython.doc.ydktest_aug_4.aug_four = 'aug four'
        cpython.lib.ydktest_aug_1.ydktest_aug_nested_1.aug_one = 'aug one'
        cpython.lib.ydktest_aug_2.ydktest_aug_nested_2.aug_two = 'aug two'
        cpython.lib.ydktest_aug_4.ydktest_aug_nested_4.aug_four = 'aug four'
        cpython.doc.disutils.four_aug_list.enabled = True

        item1 = cpython.doc.disutils.four_aug_list.Ldata()
        item2 = cpython.doc.disutils.four_aug_list.Ldata()

        item1.name, item1.number = 'one', 1
        item2.name, item1.number = 'two', 2

        self.crud.create(self.ncc, cpython)
        cpython_read = self.crud.read(self.ncc, ietf_aug_base_1.Cpython())

        self.assertEqual(cpython, cpython_read)

    def test_aug_base_2(self):
        cpython = ietf_aug_base_2.Cpython()
        cpython.tools.aug_four = 'aug four'

        self.crud.create(self.ncc, cpython)
        cpython_read = self.crud.read(self.ncc, ietf_aug_base_2.Cpython())

        self.assertEqual(cpython, cpython_read)


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
