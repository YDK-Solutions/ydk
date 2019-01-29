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

"""test_meta.py
Tests how _meta_info is generated,
getting Enum and Union info
"""
from __future__ import absolute_import

import unittest

try:
    from ydk.models.ydktest import openconfig_bgp_types
    from ydk.models.ydktest.ydktest_sanity import Runner
    from ydk.models.ydktest.openconfig_bgp import Bgp
except ImportError:
    from ydk.models.ydktest.ydktest_sanity.runner.runner import Runner
    from ydk.models.ydktest.openconfig_bgp.bgp.bgp import Bgp

from ydk._core._dm_meta_info import _MetaInfoEnum, REFERENCE_UNION

class MetaSanityTest(unittest.TestCase):

    def test_runner(self):
        runner = Runner()
        runner_meta = runner._meta_info()
        self.assertEqual(runner_meta.module_name, "ydktest-sanity")
        self.assertEqual(runner_meta.name, "Runner")
        self.assertEqual(runner_meta.yang_name, "runner")
        print(runner_meta.name, runner_meta.doc)
        for member_meta in runner_meta.meta_info_class_members:
            print("{}  --  {}".format(member_meta.name, member_meta._doc.strip()))

        one_list_meta = runner_meta.member('one-list')
        self.assertIsNotNone(one_list_meta)
        self.assertEqual(one_list_meta.name, 'one-list')

        built_in_t_meta = runner.ytypes.built_in_t._meta_info()
        embeded_enum_meta = built_in_t_meta.member('embeded-enum')
        self.assertIsNotNone(embeded_enum_meta)

        embeded_enum_map = embeded_enum_meta.enum_dict()
        print("\nEnum dictionary:\n    %s" % embeded_enum_map)
        self.assertTrue(len(embeded_enum_map) > 0)
        self.assertEqual(embeded_enum_map['seven'].value, 7L)


    def test_enum_union_meta(self):
        nbr_ipv6 = Bgp.Neighbors.Neighbor()
        nbr_ipv6.neighbor_address = '2001:db8:fff1::1'
        nbr_ipv6.config.neighbor_address = '2001:db8:fff1::1'
        nbr_ipv6.config.peer_as = 65002
        nbr_ipv6.config.peer_type = openconfig_bgp_types.PeerType.INTERNAL

        # Print neighbor address Union
        config_meta = nbr_ipv6.config._meta_info()
        print("{}  --  {}".format(config_meta.name, config_meta.doc))
        for member in config_meta.meta_info_class_members:
            if member.mtype == REFERENCE_UNION:
                print("\nUnion list of tuples:")
                for item in member.union_list():
                    print("    {}".format(item))

        # Print PeerType enum
        peer_type_meta = openconfig_bgp_types.PeerType()._meta_info()
        if isinstance(peer_type_meta, _MetaInfoEnum):
            peer_type_meta_enum_dict = peer_type_meta.enum_dict()
            self.assertTrue(len(peer_type_meta_enum_dict) > 0)
            print("\nEnum dictionary:\n    %s" % peer_type_meta_enum_dict)


if __name__ == '__main__':
    import sys
    suite = unittest.TestLoader().loadTestsFromTestCase(MetaSanityTest)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
