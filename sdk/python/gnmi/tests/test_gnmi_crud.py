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
Sanity tests for gNMIServiceProvider and CRUDServices
"""

from __future__ import absolute_import
import unittest
import logging

from test_utils import enable_logging, get_local_repo_dir, print_entity, entity_to_string

from ydk.path  import Repository
from ydk.types import EncodingFormat
from ydk.services  import CRUDService
from ydk.services  import CodecService
from ydk.providers import CodecServiceProvider
from ydk.filters import YFilter

from ydk.gnmi.providers import gNMIServiceProvider

from ydk.models.ydktest import openconfig_bgp, openconfig_interfaces

class SanityGnmiCrud(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.codec_provider = CodecServiceProvider()
        self.codec = CodecService()
        self.repo = Repository(get_local_repo_dir())
        self.provider = gNMIServiceProvider( self.repo, "127.0.0.1", 50051, "admin", "admin")
        self.schema = self.provider.get_session().get_root_schema()
        self.crud = CRUDService()

    def test_gnmi_crud_all_operations(self):
        # Configure interface
        lo10 = openconfig_interfaces.Interfaces.Interface()
        lo10.name = 'Loopback10'
        lo10.config.name = 'Loopback10'

        # Configure BGP
        bgp_global_config = openconfig_bgp.Bgp.Global.Config()
        bgp_global_config.as_ = 65172
        neighbor = openconfig_bgp.Bgp.Neighbors.Neighbor()
        neighbor.neighbor_address = '172.16.255.2'
        neighbor.config.neighbor_address = '172.16.255.2'
        neighbor.config.peer_as = 65172

        res = self.crud.create(self.provider, [lo10, bgp_global_config, neighbor])

        # Update configuration
        lo10.config.description = 'Test'
        res = self.crud.update(self.provider, lo10)
        self.assertTrue(res)
 
        # Read all
        read_list = self.crud.read(self.provider, [openconfig_interfaces.Interfaces(), openconfig_bgp.Bgp()])

        # Read config
        ifc_filter = openconfig_interfaces.Interfaces.Interface()
        ifc_filter.name = 'Loopback10'
        bgp_neighbor_filter = openconfig_bgp.Bgp.Neighbors.Neighbor()
        bgp_neighbor_filter.neighbor_address = '172.16.255.2'

        read_list = self.crud.read_config(self.provider, [ifc_filter, bgp_neighbor_filter])
        self.assertIsNotNone(read_list)
        self.assertEqual(isinstance(read_list, list), True)
        self.assertEqual(len(read_list), 2)
        #for entity in read_list:
        #    print_entity(entity, self.schema)
 
        # Read single container
        lo10 = openconfig_interfaces.Interfaces.Interface()
        lo10.name = 'Loopback10'
        lo10.config = YFilter.read
        ifc_config = self.crud.read(self.provider, lo10)
        #print_entity(ifc_config, self.schema)
        expected = '''<interface>
  <name>Loopback10</name>
  <config>
    <name>Loopback10</name>
    <description>Test</description>
  </config>
</interface>
'''
        self.assertEqual( entity_to_string(ifc_config, self.schema), expected)
 
        # Read single leaf
        lo10 = openconfig_interfaces.Interfaces.Interface()
        lo10.name = 'Loopback10'
        lo10.config.description = YFilter.read
        read_descr = self.crud.read(self.provider, lo10)
        #print_entity(read_descr, self.schema)
        expected = '''<interface>
  <name>Loopback10</name>
  <config>
    <name>Loopback10</name>
    <description>Test</description>
  </config>
</interface>
'''
        self.assertEqual( entity_to_string(read_descr, self.schema), expected)
         
        # Delete configuration
        ifc = openconfig_interfaces.Interfaces.Interface()
        ifc.name = 'Loopback10'
        bgp = openconfig_bgp.Bgp()
        res = self.crud.delete(self.provider, [ifc, bgp])

if __name__ == '__main__':
    import sys
    enable_logging(logging.ERROR)
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityGnmiCrud)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
