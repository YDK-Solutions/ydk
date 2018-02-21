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

"""test_sanity_netconf.py
sanity tests for netconf
"""
from __future__ import absolute_import

import sys
import unittest
import logging

from ydk.errors import YPYModelError, YPYError, YPYServiceError
from ydk.models.ydktest import ydktest_sanity as ysanity
from ydk.providers import NetconfServiceProvider
from ydk.services  import NetconfService, Datastore
from ydk.services  import CRUDService
from ydk.services  import CodecService
from ydk.providers import CodecServiceProvider
from ydk.ext.types import EncodingFormat

from test_utils import ParametrizedTestCase
from test_utils import get_device_info


class SanityNetconf(ParametrizedTestCase):

    @classmethod
    def setUpClass(cls):
        cls.ncc = NetconfServiceProvider(
            cls.hostname,
            cls.username,
            cls.password,
            cls.port,
            cls.protocol,
            not cls.on_demand,
            cls.common_cache,
            cls.timeout)
        cls.netconf_service = NetconfService()

    def setUp(self):
        from ydk.services import CRUDService
        crud = CRUDService()
        runner = ysanity.Runner()
        crud.delete(self.ncc, runner)

    def tearDown(self):
        pass

    def test_edit_commit_get(self):
        runner = ysanity.Runner()
        runner.ydktest_sanity_one.number = 1
        runner.ydktest_sanity_one.name = 'runner:one:name'

        get_filter = ysanity.Runner()

        op = self.netconf_service.edit_config(self.ncc, Datastore.candidate, runner)
        self.assertEqual(True, op)

        result = self.netconf_service.get_config(self.ncc, Datastore.candidate, get_filter)
        self.assertEqual(runner, result)

        op = self.netconf_service.commit(self.ncc)
        self.assertEqual(True, op)

        result = self.netconf_service.get(self.ncc, get_filter)
        self.assertEqual(runner, result)

    def test_lock_unlock(self):
        op = self.netconf_service.lock(self.ncc, Datastore.running)
        self.assertEqual(True, op)

        op = self.netconf_service.unlock(self.ncc, Datastore.running)
        self.assertEqual(True, op)

    # Failing - NetconfService glue code needed
    def test_lock_unlock_fail(self):
        op = self.netconf_service.lock(self.ncc, Datastore.candidate)
        self.assertEqual(True, op)

        try:
            op = self.netconf_service.unlock(self.ncc, Datastore.running)
        except Exception as e:
            self.assertIsInstance(e, YPYError)

    def test_validate(self):
        op = self.netconf_service.validate(self.ncc, source=Datastore.candidate)
        self.assertEqual(True, op)

        runner = ysanity.Runner()
        runner.ydktest_sanity_one.number = 1
        runner.ydktest_sanity_one.name = 'runner:one:name'
        op = self.netconf_service.validate(self.ncc, source=runner)
        self.assertEqual(True, op)

    def test_validate_fail(self):
        # should have been handled by YDK local validation
        pass

    def test_commit_discard(self):
        runner = ysanity.Runner()
        runner.two.number = 2
        runner.two.name = 'runner:two:name'
        get_filter = ysanity.Runner()

        op = self.netconf_service.edit_config(self.ncc, Datastore.candidate, runner)
        self.assertEqual(True, op)

        op = self.netconf_service.discard_changes(self.ncc)
        self.assertEqual(True, op)

        op = self.netconf_service.edit_config(self.ncc, Datastore.candidate, runner)
        self.assertEqual(True, op)

        op = self.netconf_service.commit(self.ncc)
        self.assertEqual(True, op)

        result = self.netconf_service.get(self.ncc, get_filter)
        self.assertEqual(runner, result)

    @unittest.skip('No message id in cancel commit payload')
    def test_confirmed_commit(self):
        runner = ysanity.Runner()
        runner.two.number = 2
        runner.two.name = 'runner:two:name'
        get_filter = ysanity.Runner()

        op = self.netconf_service.edit_config(self.ncc, Datastore.candidate, runner)
        self.assertEqual(True, op)

        op = self.netconf_service.commit(self.ncc, confirmed=True, confirm_timeout=120)
        self.assertEqual(True, op)

        result = self.netconf_service.get(self.ncc, get_filter)
        self.assertEqual(runner, result)

        op = self.netconf_service.cancel_commit(self.ncc)
        self.assertEqual(True, op)

    def test_copy_config(self):
        op = self.netconf_service.copy_config(self.ncc, Datastore.candidate, Datastore.running)
        self.assertEqual(True, op)

        runner = ysanity.Runner()
        runner.two.number = 2
        runner.two.name = 'runner:two:name'
        get_filter = ysanity.Runner()

        op = self.netconf_service.edit_config(self.ncc, Datastore.candidate, runner)
        self.assertEqual(True, op)

        op = self.netconf_service.copy_config(self.ncc, Datastore.running, Datastore.candidate)
        self.assertEqual(True, op)

        result = self.netconf_service.get_config(self.ncc, Datastore.running, get_filter)
        self.assertEqual(runner, result)

        runner.two.name = '%smodified' % runner.two.name

        op = self.netconf_service.copy_config(self.ncc, Datastore.running, runner)
        self.assertEqual(True, op)

        result = self.netconf_service.get_config(self.ncc, Datastore.running, get_filter)
        self.assertEqual(runner, result)

    def test_delete_config(self):
        pass
        # startup and candidate cannot be both enabled in ConfD
        # op = self.netconf_service.delete_config(self.ncc, Datastore.startup)
        # self.assertEqual(True, op)

    # Error not thrown by TCP client, YPYError is populated instead
    def test_delete_config_fail(self):
        found = False
        try:
            self.netconf_service.delete_config(self.ncc, Datastore.running)
        except (YPYError, YPYModelError):
            found = True
        self.assertEqual(found, True)

    # Failing - NetconfService glue code needed
    def test_copy_config_fail(self):
        self.assertRaises(YPYServiceError,
                          self.netconf_service.copy_config,
                          self.ncc,
                          target=123,
                          source=456)

    # Failing - NetconfService glue code needed
    def test_edit_config_fail(self):
        self.assertRaises(YPYServiceError,
                          self.netconf_service.edit_config,
                          self.ncc,
                          Datastore.startup,
                          Datastore.candidate)

    # Failing - NetconfService glue code needed
    def test_get_config_fail(self):
        runner = ysanity.Runner()
        self.assertRaises(YPYServiceError,
                          self.netconf_service.get_config,
                          self.ncc,
                          "invalid-input",
                          runner)

    # Failing - NetconfService glue code needed
    def test_lock_fail(self):
        self.assertRaises(YPYServiceError,
                          self.netconf_service.lock,
                          self.ncc,
                          "invalid-input")

    # Failing - NetconfService glue code needed
    def test_unlock_fail(self):
        self.assertRaises(YPYServiceError,
                          self.netconf_service.unlock,
                          self.ncc,
                          "invalid-input")

    def test_sanity_crud_read_interface(self):
        enable_logging(logging.ERROR)

        address = ysanity.Native.Interface.Loopback.Ipv4.Address();
        address.ip = "2.2.2.2"
        address.netmask = "255.255.255.255"

        loopback = ysanity.Native.Interface.Loopback()
        loopback.name = 2222
        loopback.ipv4.address.append(address)

        native = ysanity.Native()
        native.interface.loopback.append(loopback)

        crud = CRUDService()
        result = crud.create(self.ncc, native)

        native_read = ysanity.Native()
        interfaces = crud.read(self.ncc, native_read)

        codec_service = CodecService()
        codec_provider = CodecServiceProvider()
        codec_provider.encoding = EncodingFormat.XML
        xml_encode = codec_service.encode(codec_provider, interfaces)
        # print(xml_encode)
        # enable_logging(logging.ERROR)

def enable_logging(level):
    log = logging.getLogger('ydk')
    log.setLevel(level)
    handler = logging.StreamHandler()
    formatter = logging.Formatter(("%(asctime)s - %(name)s - %(levelname)s - %(message)s"))
    handler.setFormatter(formatter)
    log.addHandler(handler)

if __name__ == '__main__':
    device, non_demand, common_cache, timeout = get_device_info()

    suite = unittest.TestSuite()
    suite.addTest(ParametrizedTestCase.parametrize(
        SanityNetconf,
        device=device,
        non_demand=non_demand,
        common_cache=common_cache,
        timeout=timeout))
    ret = not unittest.TextTestRunner(verbosity=2).run(suite).wasSuccessful()
    sys.exit(ret)
