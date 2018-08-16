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

from test_utils import enable_logging, get_local_repo_dir, print_entity

from ydk.path  import Repository
from ydk.types import EncodingFormat
from ydk.services  import CRUDService
from ydk.services  import CodecService
from ydk.providers import CodecServiceProvider

from ydk.gnmi.providers import gNMIServiceProvider

from ydk.models.ydktest import openconfig_bgp, openconfig_interfaces

class SanityGnmiCrud(unittest.TestCase):

    @classmethod
    def setUpClass(self):
        self.codec_provider = CodecServiceProvider()
        self.codec = CodecService()
        self.repo = Repository(get_local_repo_dir())
        self.provider = gNMIServiceProvider( repo=self.repo, address="127.0.0.1", port=50051)
        self.schema = self.provider.get_session().get_root_schema()
        self.crud = CRUDService()

    def test_gnmi_crud_all_operations(self):
        # Configure interface
        ifc = openconfig_interfaces.Interfaces()
        lo10 = ifc.Interface()
        lo10.name = 'Loopback10'
        lo10.config.name = 'Loopback10'
        ifc.interface.append(lo10)

        # Configure BGP
        bgp = openconfig_bgp.Bgp()
        bgp.global_.config.as_ = 65172
        neighbor = bgp.Neighbors.Neighbor()
        neighbor.neighbor_address = '172.16.255.2'
        neighbor.config.neighbor_address = '172.16.255.2'
        neighbor.config.peer_as = 65172
        bgp.neighbors.neighbor.append(neighbor)

        res = self.crud.create(self.provider, [ifc, bgp])

        # Update configuration
        lo10.config.description = 'Test'
        res = self.crud.update(self.provider, [ifc, bgp])
 
        # Read config
        read_list = self.crud.read_config(self.provider, [openconfig_interfaces.Interfaces(), openconfig_bgp.Bgp()])
        #for entity in read_list:
        #    print_entity(entity, self.schema)
 
        # Read all
        read_list = self.crud.read(self.provider, [openconfig_interfaces.Interfaces(), openconfig_bgp.Bgp()])
 
        # Delete configuration
        res = self.crud.delete(self.provider, [ifc, bgp])

if __name__ == '__main__':
    import sys
    enable_logging(logging.ERROR)
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityGnmiCrud)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
