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

"""test_sanity_rpc.py
sanity test for netconf
"""
from __future__ import absolute_import

import unittest

from ydk.errors import YPYModelError, YPYError, YPYServiceError
from ydk.models.ydktest import ydktest_sanity as ysanity
try:
    from ydk.models.ydktest import ietf_netconf
except:
    pass
from ydk.providers import NetconfServiceProvider
from ydk.services import ExecutorService
from ydk.types import Empty


class SanityTest(unittest.TestCase):
    PROVIDER_TYPE = "non-native"

    @classmethod
    def setUpClass(self):
        self.ncc = NetconfServiceProvider('127.0.0.1', 'admin', 'admin', 12022)
        self.es = ExecutorService()

    @classmethod
    def tearDownClass(self):
        pass

    def setUp(self):
        from ydk.services import CRUDService
        crud = CRUDService()
        runner = ysanity.Runner()
        crud.delete(self.ncc, runner)

    def tearDown(self):
        pass

    def test_zclose_session_rpc(self):
        rpc = ietf_netconf.CloseSessionRpc()

        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

    # persist-id is broken?
    def test_commit_rpc(self):
        rpc = ietf_netconf.CommitRpc()
        rpc.input.confirmed = Empty()
        rpc.input.confirm_timeout = 5
        rpc.input.persist = '0'
        # rpc.input.persist_id = '0'

        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

    def test_copy_config_rpc(self):
        rpc = ietf_netconf.CopyConfigRpc()
        rpc.input.target.candidate = Empty()
        rpc.input.source.running = Empty()

        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

    @unittest.skip('Issues in netsim')
    def test_delete_config_rpc(self):
        rpc = ietf_netconf.DeleteConfigRpc()
        rpc.input.target.url = "http://test"

        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

    def test_discard_changes_rpc(self):
        rpc = ietf_netconf.DiscardChangesRpc()

        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

    @unittest.skip('Missing config/edit_content field in ietf_netconf.EditConfigRpc.Input')
    # filter in ietf_netconf.GetConfigRpc.Input ???
    # Missing filter field in ietf_netconf.GetRpc.Input
    def test_edit_config_rpc(self):
        runner = ysanity.Runner()
        runner.one.number = 1
        runner.one.name = 'runner:one:name'

        edit_rpc = ietf_netconf.EditConfigRpc()
        edit_rpc.input.target.candidate = Empty()
        # edit_rpc.input.config = runner
        # edit_rpc.input.edit_content.config = runner
        reply = self.es.execute_rpc(self.ncc, edit_rpc)
        self.assertIsNone(reply)

        get_config_rpc = ietf_netconf.GetConfigRpc()
        get_config_rpc.input.source.candidate = Empty()
        reply = self.es.execute_rpc(self.ncc, get_config_rpc)
        self.assertIsNotNone(reply)

        get_rpc = ietf_netconf.GetRpc()
        get_rpc.input.filter = ysanity.Runner()
        reply = self.es.execute_rpc(self.ncc, get_rpc)
        self.assertIsNotNone(reply)

    @unittest.skip('YCPPServiceProviderError')
    def test_kill_session(self):
        rpc = ietf_netconf.KillSessionRpc()
        rpc.input.session_id = 3
        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

    def test_lock_rpc(self):
        lock_rpc = ietf_netconf.LockRpc()
        lock_rpc.input.target.candidate = Empty()
        from pdb import set_trace; set_trace()
        reply = self.es.execute_rpc(self.ncc, lock_rpc)
        self.assertIsNone(reply)

        unlock_rpc = ietf_netconf.UnlockRpc()
        unlock_rpc.input.target.candidate = Empty()
        reply = self.es.execute_rpc(self.ncc, unlock_rpc)
        self.assertIsNone(reply)

    def test_validate_rpc_1(self):
        rpc = ietf_netconf.ValidateRpc()
        rpc.input.source.candidate = Empty()
        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

    @unittest.skip('Missing config field in ietf_netconf.ValidateRpc.Input.Source')
    def test_validate_rpc_2(self):
        runner = ysanity.Runner()
        runner.one.number = 1
        runner.one.name = 'runner:one:name'

        rpc = ietf_netconf.ValidateRpc()
        rpc.input.source.config = runner
        reply = self.es.execute_rpc(self.ncc, rpc)
        self.assertIsNone(reply)

if __name__ == '__main__':
    import sys
    if len(sys.argv) > 1:
        SanityNetconf.PROVIDER_TYPE = sys.argv.pop()

    suite = unittest.TestLoader().loadTestsFromTestCase(SanityNetconf)
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)