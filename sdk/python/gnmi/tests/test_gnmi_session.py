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

"""test_gnmi_session.py
Sanity tests for gNMISession and gNMIServiceProvider
"""

from __future__ import absolute_import
import unittest
import logging

from test_utils import enable_logging, get_local_repo_dir

from ydk.path  import Codec
from ydk.path  import Repository
from ydk.types import EncodingFormat

from ydk.gnmi.path import gNMISession

class SanityGnmiSession(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.codec = Codec()
        self.repo = Repository(get_local_repo_dir())
        self.session = gNMISession( self.repo, "127.0.0.1", 50051, "admin", "admin")
        self.schema = self.session.get_root_schema()
        self.bgp_json = ''
        self.int_json = ''

    def test_gnmi_rpc_caps(self):
        cap_rpc = self.schema.create_rpc("ydk:gnmi-caps")
        caps = self.session.invoke(cap_rpc)
 
        json = self.codec.encode(caps, EncodingFormat.JSON, True)
        cap_update = '''
      {
        "name": "openconfig-bgp",
        "organization": "OpenConfig working group",
        "version": "2016-06-21"
      },'''
        self.assertEqual((cap_update in json), True)
        #print("Server capabilities:\n%s" % json)

    def _build_configuration(self):
        # Configure BGP
        bgp = self.schema.create_datanode("openconfig-bgp:bgp")
        bgp.create_datanode("global/config/as", "65172")
        neighbor = bgp.create_datanode("neighbors/neighbor[neighbor-address='172.16.255.2']")
        neighbor_address = neighbor.create_datanode("config/neighbor-address", "172.16.255.2")
        peer_as = neighbor.create_datanode("config/peer-as","65172")

        # Configure interface
        ifc = self.schema.create_datanode("openconfig-interfaces:interfaces")
        lo10_config = ifc.create_datanode("interface[name='Loopback10']")
        lo10_config.create_datanode("config/name", 'Loopback10')
        lo10_config.create_datanode("config/description", 'Test')

        # Add data-nodes to RPC
        self.bgp_json = self.codec.encode(bgp, EncodingFormat.JSON, False)
        self.int_json = self.codec.encode(ifc, EncodingFormat.JSON, False)

        set_rpc = self.schema.create_rpc("ydk:gnmi-set")
        set_rpc.get_input_node().create_datanode("replace[alias='bgp']/entity", self.bgp_json)
        set_rpc.get_input_node().create_datanode("replace[alias='int']/entity", self.int_json)
        res = self.session.invoke(set_rpc)
        
    def _delete_configuration(self):
        del_rpc = self.schema.create_rpc("ydk:gnmi-set")
        del_rpc.get_input_node().create_datanode("delete[alias='bgp']/entity", self.bgp_json)
        del_rpc.get_input_node().create_datanode("delete[alias='int']/entity", self.int_json)
        res = self.session.invoke(del_rpc)
        
    def test_gnmi_rpc_set_get(self):
        self._build_configuration()
         
 	    # Read configuration
        bgp_read = self.schema.create_datanode("openconfig-bgp:bgp", "")
        json_bgp = self.codec.encode(bgp_read, EncodingFormat.JSON, False)
 
        int_read = self.schema.create_datanode("openconfig-interfaces:interfaces", "")
        json_int = self.codec.encode(int_read, EncodingFormat.JSON, False)
 
        read_rpc = self.schema.create_rpc("ydk:gnmi-get")
        read_rpc.get_input_node().create_datanode("type", "CONFIG")
        read_rpc.get_input_node().create_datanode("request[alias='bgp']/entity", json_bgp)
        read_rpc.get_input_node().create_datanode("request[alias='int']/entity", json_int)
        read_result = self.session.invoke(read_rpc)
 
        # Delete configuration
        self._delete_configuration()

    def test_gnmi_rpc_subscribe(self):
        self._build_configuration()

        rpc = self.schema.create_rpc("ydk:gnmi-subscribe")
        subscription = rpc.get_input_node().create_datanode("subscription", "")
        subscription.create_datanode("mode", "ONCE")
        subscription.create_datanode("qos", "10")
        subscription.create_datanode("encoding", "JSON_IETF")

        bgp_read = self.schema.create_datanode("openconfig-bgp:bgp", "")
        bgp_json = self.codec.encode(bgp_read, EncodingFormat.JSON, False)

        int_read = self.schema.create_datanode("openconfig-interfaces:interfaces", "")
        int_json = self.codec.encode(int_read, EncodingFormat.JSON, False)

        int_subscription = subscription.create_datanode("subscription-list[alias='int']", "")
        int_subscription.create_datanode("entity", int_json)
        int_subscription.create_datanode("subscription-mode", "ON_CHANGE")
        int_subscription.create_datanode("sample-interval", "10000000")
        int_subscription.create_datanode("suppress-redundant", "true")
        int_subscription.create_datanode("heartbeat-interval", "1000000000")

        bgp_subscription = subscription.create_datanode("subscription-list[alias='bgp']", "")
        bgp_subscription.create_datanode("entity", bgp_json)
        bgp_subscription.create_datanode("sample-interval", "20000000")

        self.session.subscribe(rpc, gnmi_service_subscribe_callback)
 
        # Delete configuration
        self._delete_configuration()


int_update = '''val {
      json_ietf_val: "{\\"interface\\":[{\\"name\\":\\"Loopback10\\",\\"config\\":{\\"name\\":\\"Loopback10\\",\\"description\\":\\"Test\\"}}]}"
    }'''

bgp_update = '''val {
      json_ietf_val: "{\\"global\\":{\\"config\\":{\\"as\\":65172} },\\"neighbors\\":{\\"neighbor\\":[{\\"neighbor-address\\":\\"172.16.255.2\\",\\"config\\":{\\"neighbor-address\\":\\"172.16.255.2\\",\\"peer-as\\":65172}}]}}"
    }'''

def gnmi_service_subscribe_callback(response):
    #print(response)
    if (int_update not in response) or (bgp_update not in response):
        print("ERROR in response!")
    
if __name__ == '__main__':
    #device = get_device_info()
    import sys
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityGnmiSession)
    enable_logging(logging.ERROR)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
