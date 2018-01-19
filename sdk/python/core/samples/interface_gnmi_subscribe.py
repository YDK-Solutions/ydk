#!/usr/bin/env python
#
# Copyright 2016 Cisco Systems, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

from argparse import ArgumentParser
from urllib.parse import urlparse

from multiprocessing import Pool
import logging

from ydk.providers import gNMIServiceProvider
from ydk.path import Repository
from ydk.services import gNMIService
from ydk.models.ydktest import openconfig_interfaces
from ydk.filters import YFilter


def print_telemetry_data(s):
    print(s)


def subscribe(args):
    func = args[0]
    device = args[1]
    mode = args[2]
    gnmi = gNMIService()
    repository = Repository('/Users/abhirame/.ydk/pavarotti:830')
    provider = gNMIServiceProvider(repo=repository,
                                      address=device.hostname,
                                      username=device.username,
                                      password=device.password)

    inf = openconfig_interfaces.Interfaces()
    i = openconfig_interfaces.Interfaces.Interface()
    i.yfilter = YFilter.read
    inf.interface.append(i)
    gnmi.subscribe(provider, inf, mode, 10, "ON_CHANGE", 100000, func)


if __name__ == "__main__":
    """Execute main program."""
    parser = ArgumentParser()
    parser.add_argument("-v", "--verbose", help="print debugging messages",
                        action="store_true")
    parser.add_argument("device",
                        help="gNMI device (ssh://user:password@host:port)")
    parser.add_argument("-m", "--mode", help="Subscription mode. One of 'POLL', 'ONCE', 'STREAM'", dest='mode', default='STREAM')
    args = parser.parse_args()
    device = urlparse(args.device)

    # log debug messages if verbose argument specified
    if args.verbose:
        logger = logging.getLogger("ydk")
        logger.setLevel(logging.INFO)
        handler = logging.StreamHandler()
        formatter = logging.Formatter(("%(asctime)s - %(name)s - "
                                      "%(levelname)s - %(message)s"))
        handler.setFormatter(formatter)
        logger.addHandler(handler)
 
    pool = Pool()
    pool.map(subscribe, [(print_telemetry_data, device, args.mode)])

