#!/bin/bash
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
#
# generate-bundles.sh script designed to generate YDK bundles 
# for the lates versions of cisco devices YANG models
#
# ------------------------------------------------------------------

######################################################################
# Utility functions
######################################################################

function print_msg {
    echo -e "${MSG_COLOR}*** $(date): generate-bundles.sh | $@ ${NOCOLOR}"
}

function run_cmd {
    print_msg "Running command: $@"
    $@
    local status=$?
    if [ $status -ne 0 ]; then
        MSG_COLOR=$RED
        print_msg "Exiting '$@' with status=$status"
        exit $status
    fi
    return $status
}

########################## EXECUTION STARTS HERE #############################
#

# Terminal colors
RED="\033[0;31m"
NOCOLOR="\033[0m"
YELLOW='\033[1;33m'
MSG_COLOR=$YELLOW

######################################
# Set up env

os_type=$(uname)
print_msg "Running OS type: $os_type"

export YDKGEN_HOME="$(pwd)"

args=$(getopt p:d $*)
set -- $args
PYTHON_VERSION=${2}
PYTHON_BIN=python${PYTHON_VERSION}

if [[ ${PYTHON_VERSION} = "2"* ]]; then
    PIP_BIN=pip
elif [[ ${PYTHON_VERSION} = "3.5"* ]]; then
    PIP_BIN=pip3
else
    PIP_BIN=pip${PYTHON_VERSION}
fi

print_msg "Checking installation of ${PYTHON_BIN}"
${PYTHON_BIN} -V
status=$?
if [ $status -ne 0 ]; then
    print_msg "Could not locate ${PYTHON_BIN}"
    exit $status
fi
print_msg "Checking installation of ${PIP_BIN}"
${PIP_BIN} -V
status=$?
if [ $status -ne 0 ]; then
    print_msg "Could not locate ${PIP_BIN}"
    exit $status
fi
print_msg "Python location: $(which ${PYTHON_BIN})"
print_msg "Pip location: $(which ${PIP_BIN})"

CMAKE_BIN=cmake
which cmake3
status=$?
if [[ ${status} == 0 ]] ; then
    CMAKE_BIN=cmake3
fi

print_msg "Generating Python core and bundles"
run_cmd ./generate.py --core
run_cmd ./generate.py --service profiles/services/gnmi-0.4.0.json
run_cmd ./generate.py --bundle profiles/bundles/ietf_0_1_5_post2.json
run_cmd ./generate.py --bundle profiles/bundles/openconfig_0_1_6_post1.json
run_cmd ./generate.py --bundle profiles/bundles/cisco-ios-xe_16_9_1.json
run_cmd ./generate.py --bundle profiles/bundles/cisco-ios-xr_6_5_1.json
run_cmd ./generate.py --bundle profiles/bundles/cisco-nx-os-9_2_2.json

print_msg "Generating C++ core and bundles"
run_cmd ./generate.py --core --cpp
run_cmd ./generate.py --service profiles/services/gnmi-0.4.0.json --cpp
run_cmd ./generate.py --bundle profiles/bundles/ietf_0_1_5_post2.json --cpp
run_cmd ./generate.py --bundle profiles/bundles/openconfig_0_1_6_post1.json --cpp
run_cmd ./generate.py --bundle profiles/bundles/cisco-ios-xe_16_9_1.json --cpp
run_cmd ./generate.py --bundle profiles/bundles/cisco-ios-xr_6_5_1.json  --cpp
run_cmd ./generate.py --bundle profiles/bundles/cisco-nx-os-9_2_2.json  --cpp

print_msg "Generating Go core and bundles"
run_cmd ./generate.py --core --go
run_cmd ./generate.py --service profiles/services/gnmi-0.4.0.json --go
run_cmd ./generate.py --bundle profiles/bundles/ietf_0_1_5_post2.json --go
run_cmd ./generate.py --bundle profiles/bundles/openconfig_0_1_6_post1.json --go
run_cmd ./generate.py --bundle profiles/bundles/cisco-ios-xe_16_9_1.json --go
run_cmd ./generate.py --bundle profiles/bundles/cisco-ios-xr_6_5_1.json  --go
run_cmd ./generate.py --bundle profiles/bundles/cisco-nx-os-9_2_2.json --go
