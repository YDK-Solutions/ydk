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

import sys
import unittest

from ydk.errors import YPYInvalidArgumentError
from ydk.types  import EntityCollection, Filter, Config

from ydk.models.ydktest import ydktest_sanity as ysanity

from test_utils import assert_with_error

test_add_unsupported_pattern = '''Argument <(type|class) 'str'> is not supported by EntityCollection class; data ignored'''

class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        pass

    def test_entity_collection_empty(self):
        anydata = EntityCollection()
        self.assertEqual(anydata.is_empty(), True)
        self.assertEqual(anydata.len(), 0)
        self.assertEqual(format(anydata), "Entities in EntityCollection: []")

    def test_config_add(self):
        runner = ysanity.Runner()
        native = ysanity.Native()

        config = Config()
        config.add(runner)
        self.assertEqual(config.is_empty(), False)
        self.assertEqual(config.len(), 1)
        self.assertEqual(format(config), "Entities in Config: ['ydk.models.ydktest.ydktest_sanity.Runner']")

        config.add(native)
        self.assertEqual(config.len(), 2)
        self.assertEqual(format(config), "Entities in Config: ['ydk.models.ydktest.ydktest_sanity.Runner', 'ydk.models.ydktest.ydktest_sanity.Native']")

        config.clear()

    def test_filter_list(self):
        runner = ysanity.Runner()
        native = ysanity.Native()

        read_filter = Filter([runner, native])

        self.assertEqual(read_filter.get_keys(), ['ydktest-sanity:runner', 'ydktest-sanity:native'])
        entities = read_filter.get_entities()
        self.assertEqual(format(entities[0]), 'ydk.models.ydktest.ydktest_sanity.Runner')
        self.assertEqual(format(entities[1]), 'ydk.models.ydktest.ydktest_sanity.Native')
        read_filter.clear()

    def test_access_config_by_key(self):
        runner = ysanity.Runner()
        native = ysanity.Native()

        config = Config([runner, native])
        for key in config.get_keys():
            if key=='ydktest-sanity:runner':
                self.assertEqual(config[key].get_segment_path(), 'ydktest-sanity:runner')
            else:
                self.assertEqual(config[key].get_segment_path(), 'ydktest-sanity:native')
        config.clear()

    def test_access_config_by_item_no(self):
        runner = ysanity.Runner()
        native = ysanity.Native()

        config = Config([runner, native])
        for item in range(0, config.len()):
            if item==0:
                self.assertEqual(format(config[item]), "ydk.models.ydktest.ydktest_sanity.Runner")
            else:
                self.assertEqual(format(config[item]), "ydk.models.ydktest.ydktest_sanity.Native")
        config.clear()

    def test_access_config_by_iter(self):
        runner = ysanity.Runner()
        native = ysanity.Native()

        config = Config([runner, native])
        for entity in config:
            print(entity)

        self.assertEqual(config.len(), 2)
        config.clear()

    @assert_with_error(test_add_unsupported_pattern, YPYInvalidArgumentError)
    def test_add_unsupported(self):
        anydata = EntityCollection()
        anydata.add('native')
        self.assertEqual(anydata.is_empty(), True)

    def test_init_delete(self):
        runner = ysanity.Runner()
        native = ysanity.Native()

        anydata = Config(runner)
        deleted = anydata.delete(0)
        self.assertEqual(anydata.is_empty(), True)
        self.assertEqual(deleted.__class__.__name__, "Runner")

        anydata = Config([runner, native])
        deleted = anydata.delete('ydktest-sanity:native')
        self.assertEqual(anydata.len(), 1)
        self.assertEqual(deleted.__class__.__name__, "Native")

        deleted = anydata.delete('ydktest-sanity:native')
        self.assertEqual(deleted, None)

        deleted = anydata.delete(runner)
        self.assertEqual(anydata.is_empty(), True)
        self.assertEqual(deleted.__class__.__name__, "Runner")

if __name__ == '__main__':
    suite = unittest.TestSuite()
    testloader = unittest.TestLoader()
    testnames = testloader.getTestCaseNames(SanityTest)
    for name in testnames:
        suite.addTest(SanityTest(name))
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
