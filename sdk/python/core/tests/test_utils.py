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

"""test_utils.py

Utility function for test cases.
"""
import re
import unittest

def assert_with_error(pattern, ErrorClass):
    def assert_with_pattern(func):
        def helper(self, *args, **kwargs):
            try:
                func(self)
            except ErrorClass as error:
                res = re.match(pattern, error.message.strip())
                self.assertEqual(res is not None, True)
        return helper
    return assert_with_pattern


class ParametrizedTestCase(unittest.TestCase):
    """ TestCase classes that want to be parametrized should
        inherit from this class.
    """
    def __init__(self, methodName='runTest', port=12022, protocol='ssh'):
        super(ParametrizedTestCase, self).__init__(methodName)
        self.port = port
        self.protocol = protocol

    @staticmethod
    def parametrize(testcase_klass, port=12022, protocol='ssh'):
        """ Create a suite containing all tests taken from the given
            subclass, passing them the parameter 'param'.
        """
        testloader = unittest.TestLoader()
        testcase_klass.port = port
        testcase_klass.protocol = protocol
        testnames = testloader.getTestCaseNames(testcase_klass)
        suite = unittest.TestSuite()
        for name in testnames:
            suite.addTest(testcase_klass(name))
        return suite
