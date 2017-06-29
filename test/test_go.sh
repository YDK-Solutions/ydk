#!/bin/bash
#  ----------------------------------------------------------------
# Copyright 2017 Cisco Systems
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
# Script for running ydk CI on travis-ci.org
#
# ------------------------------------------------------------------

function go_samples {
    print_msg "CC: ${CC}"
    print_msg "CXX: ${CXX}"
    export CXX=/usr/bin/clang++
    export CC=/usr/bin/clang
    run_exec_test cd $YDKGEN_HOME/sdk/go/core/samples/cgo_path && go run cgo_path*
    run_exec_test cd $YDKGEN_HOME/sdk/go/core/samples/bgp_create && go run bgp_create*
    run_exec_test cd $YDKGEN_HOME/sdk/go/core/samples/bgp_read && go run bgp_read*
    run_exec_test cd $YDKGEN_HOME/sdk/go/core/samples/bgp_delete && go run bgp_delete*
}

function go_behavioral_tests {
    run_exec_test cd $YDKGEN_HOME/sdk/go/core/tests && go test
}

function go_tests {
    go_samples
    go_behavioral_tests
}
