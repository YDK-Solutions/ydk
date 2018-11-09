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

"""test_metadata.py
Tests how metadata is generated
"""
from __future__ import absolute_import

import os
import unittest

try:
    from ydk.models.ydktest.ydktest_sanity import Runner
except:
    from ydk.models.ydktest.ydktest_sanity.runner.runner import Runner

class MetaSanityTest(unittest.TestCase):

    def test_runner(self):
        runner = Runner()
        meta = runner._meta_info()
        self.assertEqual(meta.module_name, "ydktest-sanity")
        self.assertEqual(meta.name, "Runner")
        self.assertEqual(meta.yang_name, "runner")
        for member_meta in meta.meta_info_class_members:
            name = member_meta.name
            print(name)

if __name__ == '__main__':
    import sys
    suite = unittest.TestLoader().loadTestsFromTestCase(MetaSanityTest)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
