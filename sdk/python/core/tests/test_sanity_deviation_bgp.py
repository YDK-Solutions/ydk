from __future__ import absolute_import
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
import sys
import unittest

import ydk.types as ytypes
from ydk.services import CRUDService
from ydk.providers import NetconfServiceProvider
from ydk.types import Empty, Decimal64
from ydk.errors import YPYError, YPYModelError
from ydk.models.deviation import openconfig_bgp, openconfig_bgp_types
from ydk.models.deviation.openconfig_routing_policy import DefaultPolicyType, RoutingPolicy

from test_utils import assert_with_error
from test_utils import ParametrizedTestCase
from test_utils import get_device_info


class SanityTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        hostname = getattr(cls, 'hostname', '127.0.0.1')
        username = getattr(cls, 'username', 'admin')
        password = getattr(cls, 'password', 'admin')
        port = 12023
        protocol = getattr(cls, 'protocol', 'ssh')
        on_demand = not getattr(cls, 'non_demand', True)
        common_cache = getattr(cls, "common_cache", False)
        cls.ncc = NetconfServiceProvider(hostname, username, password, port, protocol, on_demand, common_cache)
        cls.crud = CRUDService()

    @classmethod
    def tearDownClass(cls):
        pass

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_bgp(self):
        # Bgp.Global.AfiSafis.AfiSafi.ApplyPolicy is not supported
        bgp_cfg = openconfig_bgp.Bgp()
        ipv4_afsf = bgp_cfg.global_.afi_safis.AfiSafi()
        ipv4_afsf.afi_safi_name = openconfig_bgp_types.AfiSafiType()
        ipv4_afsf.apply_policy.config.default_export_policy = DefaultPolicyType.ACCEPT_ROUTE
        bgp_cfg.global_.afi_safis.afi_safi.append(ipv4_afsf)

        self.assertRaises(YPYModelError, self.crud.create, self.ncc, bgp_cfg)


if __name__ == '__main__':
    device, non_demand, common_cache = get_device_info()

    suite = unittest.TestSuite()
    suite.addTest(ParametrizedTestCase.parametrize(SanityTest, device=device, non_demand=non_demand, common_cache=common_cache))
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
