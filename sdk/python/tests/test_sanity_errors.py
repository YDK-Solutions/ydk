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

import ydk.types as ytypes
import unittest

from ydk.services import CRUDService
from ydk.models.ydktest import ydktest_sanity as ysanity
from ydk.models.ydktest import ydktest_sanity_types as ysanity_types
from ydk.providers import NetconfServiceProvider
from ydk.types import Empty, DELETE, Decimal64
from tests.compare import is_equal
from ydk.errors import YPYError, YPYModelError

from ydk.models.ydktest.ydktest_sanity import YdkEnumTestEnum, YdkEnumIntTestEnum

class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.ncc = NetconfServiceProvider(address='127.0.0.1',
            username='admin', password='admin', port=12022)
        self.crud = CRUDService()

    @classmethod
    def tearDownClass(self):
        self.ncc.close()

    def setUp(self):
        print "\nIn method", self._testMethodName
        runner = ysanity.Runner()
        self.crud.delete(self.ncc, runner)

    def tearDown(self):
        runner = ysanity.Runner()
        self.crud.delete(self.ncc, runner)

    def _create_runner(self):
        runner = ysanity.Runner()
        runner.ytypes = runner.Ytypes()
        runner.ytypes.built_in_t = runner.ytypes.BuiltInT()

        return runner

    def test_int8_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.number8 = 8.5
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_int16_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.number16 = {}
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_int32_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.number32 = []
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_int64_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.number64 = 9223372036854775808
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_uint8_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.u_number8 = -1
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_uint16_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.u_number16 = 'not an uint'
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_uint16_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.u_number16 = 'not an uint'
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_uint32_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.u_number32 = 4294967296
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_uint64_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.u_number64 = 18446744073709551616
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYError,
            self.crud.create, self.ncc, runner)

    # but invalid input is sent to the device
    def test_string_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.name = ['name_str']
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)
        # self.crud.create(self.ncc, runner)

    def test_empty_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.emptee = '0'
        
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_boolean_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.bool_value = ''
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)

    def test_enum_invalid(self):
        runner = self._create_runner()
        runner.ytypes.built_in_t.enum_value = 'not an enum'
        # self.crud.create(self.ncc, runner)
        self.assertRaises(YPYModelError,
            self.crud.create, self.ncc, runner)


