"""test_gnmi_provider.py
gNMIServiceProvider test
"""
from __future__ import absolute_import

import sys
import unittest

from ydk.providers import gNMIServiceProvider
from test_utils import ParametrizedTestCase
from test_utils import get_device_info


class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.ncc = gNMIServiceProvider(cls.address, cls.common_cache)

    def test_get_encoding(self):
        encoding = self.ncc.get_encoding()
        self.assertEqual(encoding is not None, True)

    def test_get_capabilities(self):
        capabilities = self.ncc.get_capabilities()
        self.assertEqual(capabilities is not None, True)


if __name__ == '__main__':
    device, common_cache = get_device_info()

    suite = unittest.TestSuite()
    suite.addTest(ParametrizedTestCase.parametrize(SanityTest, device=device, common_cache=common_cache))
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)