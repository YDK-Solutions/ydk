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

from __future__ import absolute_import
import unittest

from ydk.services import CrudService
from ydk.models.ydktest import ydktest_sanity as ysanity
from ydk.models.ydktest import ydktest_sanity_types as ysanity_types
from ydk.providers import NetconfServiceProvider
from ydk.types import Empty, Decimal64
from ydk.errors import YPYError, YPYModelError, YPYServiceError
from ydk.models.ydktest.ydktest_sanity import YdkEnumTestEnum, YdkEnumIntTestEnum


class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.ncc = NetconfServiceProvider('127.0.0.1', 'admin', 'admin', 12022)
        self.crud = CrudService()

    @classmethod
    def tearDownClass(self):
        pass

    def setUp(self):
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

    @unittest.skip('segfault')
    def test_int8_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.number8 = 8.5
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "8.5" in "number8" element. Path: '/ydktest-sanity:runner/ytypes/built-in-t/number8'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    def test_int16_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.number16 = {}
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """set(): incompatible function arguments. The following argument types are supported:
    1. (self: ydk_.types.YLeaf, arg0: int) -> None
    2. (self: ydk_.types.YLeaf, arg0: int) -> None
    3. (self: ydk_.types.YLeaf, arg0: int) -> None
    4. (self: ydk_.types.YLeaf, arg0: int) -> None
    5. (self: ydk_.types.YLeaf, arg0: int) -> None
    6. (self: ydk_.types.YLeaf, arg0: int) -> None
    7. (self: ydk_.types.YLeaf, arg0: float) -> None
    8. (self: ydk_.types.YLeaf, arg0: ydk_.types.Empty) -> None
    9. (self: ydk_.types.YLeaf, arg0: ydk_.types.Identity) -> None
    10. (self: ydk_.types.YLeaf, arg0: ydk_.types.Bits) -> None
    11. (self: ydk_.types.YLeaf, arg0: unicode) -> None
    12. (self: ydk_.types.YLeaf, arg0: ydk_.types.YLeaf) -> None
    13. (self: ydk_.types.YLeaf, arg0: ydk_.types.Decimal64) -> None

Invoked with: , {}"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    def test_int32_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.number32 = []
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """set(): incompatible function arguments. The following argument types are supported:
    1. (self: ydk_.types.YLeaf, arg0: int) -> None
    2. (self: ydk_.types.YLeaf, arg0: int) -> None
    3. (self: ydk_.types.YLeaf, arg0: int) -> None
    4. (self: ydk_.types.YLeaf, arg0: int) -> None
    5. (self: ydk_.types.YLeaf, arg0: int) -> None
    6. (self: ydk_.types.YLeaf, arg0: int) -> None
    7. (self: ydk_.types.YLeaf, arg0: float) -> None
    8. (self: ydk_.types.YLeaf, arg0: ydk_.types.Empty) -> None
    9. (self: ydk_.types.YLeaf, arg0: ydk_.types.Identity) -> None
    10. (self: ydk_.types.YLeaf, arg0: ydk_.types.Bits) -> None
    11. (self: ydk_.types.YLeaf, arg0: unicode) -> None
    12. (self: ydk_.types.YLeaf, arg0: ydk_.types.YLeaf) -> None
    13. (self: ydk_.types.YLeaf, arg0: ydk_.types.Decimal64) -> None

Invoked with: , []"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    def test_int64_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.number64 = 9223372036854775808
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "9223372036854775808" in "number64" element. Path: '/ydktest-sanity:runner/ytypes/built-in-t/number64'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')
        # runner = self._create_runner()
        # runner.ytypes.built_in_t.number64 = 9223372036854775808
        # self.crud.create(self.ncc, runner)

    @unittest.skip('segfault')
    def test_uint8_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.u_number8 = -1
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "-1" in "u_number8" element. Path: '/ydktest-sanity:runner/ytypes/built-in-t/u_number8'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    @unittest.skip('segfault')
    def test_uint16_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.u_number16 = 'not an uint'
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "not an uint" in "u_number16" element. Path: \'/ydktest-sanity:runner/ytypes/built-in-t/u_number16\'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    @unittest.skip('segfault')
    def test_uint32_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.u_number32 = 4294967296
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "4294967296" in "u_number32" element. Path: '/ydktest-sanity:runner/ytypes/built-in-t/u_number32'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    @unittest.skip('error, overflow, exception not raised')
    def test_uint64_invalid_1(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.u_number64 = -1
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = "Runner.Ytypes.BuiltInT.u_number64: (INVALID_VALUE, Value is invalid: -1 not in range (0, 18446744073709551615))"
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    @unittest.skip('segfault')
    def test_uint64_invalid_2(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.u_number64 = 18446744073709551616
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "1.84467e+19" in "u_number64" element. Path: '/ydktest-sanity:runner/ytypes/built-in-t/u_number64'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    def test_string_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.name = ['name_str']
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """set(): incompatible function arguments. The following argument types are supported:
    1. (self: ydk_.types.YLeaf, arg0: int) -> None
    2. (self: ydk_.types.YLeaf, arg0: int) -> None
    3. (self: ydk_.types.YLeaf, arg0: int) -> None
    4. (self: ydk_.types.YLeaf, arg0: int) -> None
    5. (self: ydk_.types.YLeaf, arg0: int) -> None
    6. (self: ydk_.types.YLeaf, arg0: int) -> None
    7. (self: ydk_.types.YLeaf, arg0: float) -> None
    8. (self: ydk_.types.YLeaf, arg0: ydk_.types.Empty) -> None
    9. (self: ydk_.types.YLeaf, arg0: ydk_.types.Identity) -> None
    10. (self: ydk_.types.YLeaf, arg0: ydk_.types.Bits) -> None
    11. (self: ydk_.types.YLeaf, arg0: unicode) -> None
    12. (self: ydk_.types.YLeaf, arg0: ydk_.types.YLeaf) -> None
    13. (self: ydk_.types.YLeaf, arg0: ydk_.types.Decimal64) -> None

Invoked with: , ['name_str']"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    @unittest.skip('segfault')
    def test_empty_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.emptee = '0'
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "0" in "emptee" element. Path: '/ydktest-sanity:runner/ytypes/built-in-t/emptee'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    @unittest.skip('error, empty string implicitly converted to false, exception not raised')
    def test_boolean_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.bool_value = ''
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = "Runner.Ytypes.BuiltInT.bool_value: (INVALID_TYPE, Invalid type: 'str'. Expected type: 'bool')"
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    @unittest.skip('segfault')
    def test_enum_invalid(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.enum_value = 'not an enum'
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value "not an enum" in "enum-value" element. Path: '/ydktest-sanity:runner/ytypes/built-in-t/enum-value'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    def test_yleaflist_assignment(self):
        try:
            runner = self._create_runner()
            runner.ytypes.built_in_t.llstring = ['invalid', 'leaf-list', 'assignment']
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = """Invalid value 'llstring' in '['invalid', 'leaf-list', 'assignment']'"""
            self.assertEqual(err.message.strip(), expected_msg)
        else:
            raise Exception('YPYModelError not raised')

    def test_ylist_assignment(self):
        try:
            runner = self._create_runner()
            elems, n = [], 10
            for i in range(n):
                l = ysanity.Runner.OneList.Ldata()
                l.number = i
                l.name = str(i)
                elems.append(l)
            runner.one_list.ldata = elems
            self.crud.create(self.ncc, runner)
        except YPYModelError as err:
            expected_msg = "Attempt to assign object of type list to YList ldata. Please use list append or extend method."
            self.assertEqual(err.message.strip(), expected_msg)

if __name__ == '__main__':
    import sys
    if len(sys.argv) > 1:
        SanityTest.PROVIDER_TYPE = sys.argv.pop()
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityTest)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
