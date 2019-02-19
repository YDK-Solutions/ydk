#  ----------------------------------------------------------------
# Copyright 2019 Cisco Systems
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

"""test_oc_nis.py
Tests for openconfig-network-instance model
"""
from __future__ import absolute_import

import unittest

from ydk.providers import CodecServiceProvider
from ydk.services import CodecService

from ydk.models.ydktest_oc_nis import openconfig_network_instance as oc_nis

nis_xml = '''
<network-instances xmlns="http://openconfig.net/yang/network-instance">
  <network-instance>
    <name>default</name>
    <protocols>
      <protocol>
        <identifier xmlns:policy-types="http://openconfig.net/yang/policy-types">policy-types:ISIS</identifier>
        <name>DEFAULT</name>
        <config>
          <identifier xmlns:policy-types="http://openconfig.net/yang/policy-types">policy-types:ISIS</identifier>
          <name>DEFAULT</name>
        </config>
        <isis>
          <global>
            <afi-safi>
              <af>
                <afi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:IPV4</afi-name>
                <safi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:UNICAST</safi-name>
                <config>
                  <afi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:IPV4</afi-name>
                  <safi-name xmlns:isis-types="http://openconfig.net/yang/isis-types">isis-types:UNICAST</safi-name>
                </config>
              </af>
            </afi-safi>
          </global>
        </isis>
      </protocol>
    </protocols>
  </network-instance>
</network-instances>
'''

class SanityTest(unittest.TestCase):

    def test_oc_nis_decode(self):
        codec_provider = CodecServiceProvider(type='xml')
        codec = CodecService()

        nis_top = oc_nis.NetworkInstances()
        entity = codec.decode(codec_provider, nis_xml, nis_top);
        self.assertIsNotNone(entity)

if __name__ == '__main__':
    import sys
    suite = unittest.TestLoader().loadTestsFromTestCase(SanityTest)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
