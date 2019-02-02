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

"""test_gnmi_service.py
Sanity tests for gNMIService
"""

from __future__ import absolute_import
import unittest
import logging

from multiprocessing import Pool

from test_utils import enable_logging, get_local_repo_dir, print_entity, entity_to_string, EmptyTest

from ydk.filters import YFilter
from ydk.path  import Repository
from ydk.types import EncodingFormat, Filter, Config
from ydk.services  import CodecService
from ydk.providers import CodecServiceProvider

from ydk.gnmi.providers import gNMIServiceProvider
from ydk.gnmi.services import gNMIService, gNMISubscription

from ydk.models.ydktest import openconfig_bgp, openconfig_interfaces

int_update = '''json_ietf_val: "{\\"interface\\":[{\\"name\\":\\"Loopback10\\",\\"config\\":{\\"name\\":\\"Loopback10\\",\\"description\\":\\"Test\\"}}]}"'''

bgp_update = '''json_ietf_val: "{\\"global\\":{\\"config\\":{\\"as\\":65172} },\\"neighbors\\":{\\"neighbor\\":[{\\"neighbor-address\\":\\"172.16.255.2\\",\\"config\\":{\\"neighbor-address\\":\\"172.16.255.2\\",\\"peer-as\\":65172}}]}}"'''

def build_bgp_config():
    #Build BGP configuration on server
    bgp = openconfig_bgp.Bgp()
    bgp.global_.config.as_ = 65172
    neighbor = bgp.Neighbors.Neighbor()
    neighbor.neighbor_address = '172.16.255.2'
    neighbor.config.neighbor_address = '172.16.255.2'
    neighbor.config.peer_as = 65172
    bgp.neighbors.neighbor.append(neighbor)
    return bgp

def build_int_config():
    lo10 = openconfig_interfaces.Interfaces.Interface()
    lo10.name = 'Loopback10'
    lo10.config.name = 'Loopback10'
    lo10.config.description = 'Test'
    return lo10

class SanityGnmiService(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.codec_provider = CodecServiceProvider()
        self.codec = CodecService()
        self.repo = Repository(get_local_repo_dir())
        self.provider = gNMIServiceProvider( self.repo, "127.0.0.1", 50051, "admin", "admin")
        self.schema = self.provider.get_session().get_root_schema()
        self.gs = gNMIService()

    def test_gnmi_service_capabilities(self):
        json_caps = self.gs.capabilities(self.provider);
        cap_update = '''
      {
        "name": "openconfig-bgp",
        "organization": "OpenConfig working group",
        "version": "2016-06-21"
      },'''
        self.assertEqual((cap_update in json_caps), True)
 
    def test_gnmi_service_set_get(self):
        # Create interface configuration
        ifc = build_int_config()
        ifc.yfilter = YFilter.replace
        reply = self.gs.set(self.provider, ifc)
 
        # Update interface configuration
        ifc.config.description = 'Test'
        ifc.yfilter = YFilter.update
        reply = self.gs.set(self.provider, ifc)
 
        # Get interface configuration
        ifc_filter = openconfig_interfaces.Interfaces.Interface()
        ifc_filter.name = 'Loopback10'
        response = self.gs.get(self.provider, ifc_filter, "CONFIG")
        self.assertIsNotNone(response)
        self.assertTrue(hasattr(response, 'name'))
        self.assertEqual(response.name, 'Loopback10')
 
        # Get interface description only
        ifc_filter = openconfig_interfaces.Interfaces.Interface()
        ifc_filter.name = 'Loopback10'
        ifc_filter.config.description = YFilter.read

        response = self.gs.get(self.provider, ifc_filter, "CONFIG")
        self.assertIsNotNone(response)
        #print_entity(ifc_config, self.schema)
        expected = '''<interface>
  <name>Loopback10</name>
  <config>
    <name>Loopback10</name>
    <description>Test</description>
  </config>
</interface>
'''
        self.assertEqual( entity_to_string(response, self.schema), expected)
 
        # Delete interface configuration
        ifc_delete = openconfig_interfaces.Interfaces.Interface()
        ifc_delete.name = 'Loopback10'
        ifc_delete.yfilter = YFilter.delete
        reply = self.gs.set(self.provider, ifc_delete)
 
    def test_gnmi_service_set_get_multiple(self):
        # Create interface and BGP configuration
        ifc = build_int_config()
        ifc.yfilter = YFilter.replace
        bgp = build_bgp_config()
        bgp.yfilter = YFilter.replace
        reply = self.gs.set(self.provider, Config(ifc, bgp))
 
        # Get and print interface and BGP configuration
        ifc_filter = openconfig_interfaces.Interfaces()
        bgp_filter = openconfig_bgp.Bgp()
        response = self.gs.get(self.provider, Filter(ifc_filter, bgp_filter), "CONFIG")
        self.assertIsNotNone(response)
        self.assertEqual(response.__len__(), 2)
        #for entity in response:
        #    print_entity(entity, self.schema)
 
        # Delete interface and BGP configuration
        ifc.yfilter = YFilter.delete
        bgp.yfilter = YFilter.delete
        reply = self.gs.set(self.provider, [ifc, bgp])
  
    def test_gnmi_service_subscribe_once(self):
        # Create BGP configuration
        bgp = build_bgp_config()
        bgp.yfilter = YFilter.replace
        reply = self.gs.set(self.provider, bgp)
 
        subscription = gNMISubscription()
        bgp_filter = openconfig_bgp.Bgp()
        subscription.subscription_mode = "ON_CHANGE"
        subscription.sample_interval = 10 * 1000000000
        subscription.suppress_redundant = True
        subscription.heartbeat_interval = 100 * 1000000000
        subscription.entity = bgp_filter
 
        self.gs.subscribe(self.provider, subscription, 10, "ONCE", "JSON_IETF", gnmi_subscribe_callback);
 
        # Delete BGP configuration
        bgp.yfilter = YFilter.delete
        reply = self.gs.set(self.provider, bgp)
 
    def test_gnmi_service_subscribe_multiple(self):
        # Create interface and BGP configurations
        bgp = build_bgp_config()
        bgp.yfilter = YFilter.replace
 
        ifc = build_int_config()
        ifc.yfilter = YFilter.replace
        reply = self.gs.set(self.provider, [bgp, ifc])
         
        bgp_filter = openconfig_bgp.Bgp()
        bgp_subscription = gNMISubscription()
        bgp_subscription.entity = bgp_filter
 
        int_filter = openconfig_interfaces.Interfaces()
        int_subscription = gNMISubscription()
        int_subscription.entity = int_filter
 
        self.gs.subscribe(self.provider, [int_subscription, bgp_subscription], 10, "ONCE", "JSON_IETF", gnmi_subscribe_multiples_callback)
 
        # Delete BGP configuration
        bgp.yfilter = YFilter.delete
        ifc.yfilter = YFilter.delete
        reply = self.gs.set(self.provider, [bgp, ifc])
 
    def test_gnmi_service_subscribe_stream(self):
        # Create BGP configuration
        bgp = build_bgp_config()
        bgp.yfilter = YFilter.replace
        reply = self.gs.set(self.provider, bgp)

        subscription = gNMISubscription()
        bgp_filter = openconfig_bgp.Bgp()
        subscription.subscription_mode = "SAMPLE"
        subscription.sample_interval = 2 * 1000000000
        subscription.suppress_redundant = False
        subscription.heartbeat_interval = 8 * 1000000000
        subscription.entity = bgp_filter

        self.gs.subscribe(self.provider, subscription, 10, "STREAM", "JSON_IETF", gnmi_subscribe_callback)

        # Delete BGP configuration
        bgp.yfilter = YFilter.delete
        reply = self.gs.set(self.provider, bgp)

    def test_gnmi_service_subscribe_poll(self):
        # Create BGP configuration
        bgp = build_bgp_config()
        bgp.yfilter = YFilter.replace
        reply = self.gs.set(self.provider, bgp)

        subscription = gNMISubscription()
        bgp_filter = openconfig_bgp.Bgp()
        subscription.entity = bgp_filter

        self.gs.subscribe(self.provider, subscription, 10, "POLL", "JSON_IETF", gnmi_subscribe_callback)

        # Delete BGP configuration
        bgp.yfilter = YFilter.delete
        reply = self.gs.set(self.provider, bgp)

def read_sub(response):
    print("===> In response callback. Received subscribe response:")
    print(response)

def gnmi_subscribe_callback(response):
    #read_sub(response)
    test = EmptyTest()
    test.assertTrue(bgp_update in response)

def gnmi_subscribe_multiples_callback(response):
    test = EmptyTest()
    test.assertTrue(int_update in response)
    test.assertTrue(bgp_update in response)

if __name__ == '__main__':
    import sys
    enable_logging(logging.ERROR)
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityGnmiService)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
