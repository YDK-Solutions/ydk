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

function cpp_sanity_core_gen_install {
    print_msg "cpp_sanity_core_gen_install"

    cd $YDKGEN_HOME && source gen_env/bin/activate
    cd $YDKGEN_HOME/sdk/cpp/core
    mkdir -p build && cd build
    run_exec_test cmake -DCMAKE_C_COMPILER=/usr/bin/clang -DCMAKE_CXX_COMPILER=/usr/bin/clang++ ..
    run_exec_test make install
    cd $YDKGEN_HOME
}

function cpp_sanity_core_test {
    print_msg "Running cpp core test"

    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/ydktest
    cd $YDKGEN_HOME/sdk/cpp/core/build
    make test
    local status=$?
    if [ $status -ne 0 ]; then
    # If the tests fail, try to run them in verbose to get more details for  # debug
        ./tests/ydk_core_test -d yes
        exit $status
    fi
    cd $YDKGEN_HOME
}

function cpp_sanity_ydktest {
    print_msg "Generating and testing bundle"

    cpp_sanity_ydktest_gen_install
    cpp_sanity_ydktest_test
}

function generate_install_cpp_bundle {
   bundle_profile=$1
   bundle_name=$2
   cd $YDKGEN_HOME && source gen_env/bin/activate
    run_test generate.py --bundle $bundle_profile --cpp --generate-doc
    cd gen-api/cpp/$2/build
    run_exec_test make install
    cd -
}

function cpp_sanity_ydktest_gen_install {
    print_msg "Generating and installing ydktest bundle"
    generate_install_cpp_bundle profiles/test/ydktest-cpp.json ydktest-bundle

    print_msg "Generating and installing new ydktest bundle"
    generate_install_cpp_bundle profiles/test/ydktest-cpp-new.json ydktest_new-bundle
}

function cpp_sanity_ydktest_test {
    print_msg "Running cpp bundle tests"

    mkdir -p $YDKGEN_HOME/sdk/cpp/tests/build && cd sdk/cpp/tests/build
    run_exec_test cmake -DCMAKE_C_COMPILER=/usr/bin/clang -DCMAKE_CXX_COMPILER=/usr/bin/clang++ ..
    run_exec_test make
    make test
    local status=$?
    if [ $status -ne 0 ]; then
    # If the tests fail, try to run them in verbose to get more details for  # debug
        ./ydk_bundle_test -d yes
        exit $status
    fi
}

function cpp_test_gen_test {
    print_msg "cpp_test_gen_test"

    cd $YDKGEN_HOME
    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/testgen/confd
    mkdir -p gen-api/cpp/models_test-bundle/ydk/models/models_test/test/build
    cd gen-api/cpp/models_test-bundle/ydk/models/models_test/test/build
    run_exec_test cmake -DCMAKE_C_COMPILER=/usr/bin/clang -DCMAKE_CXX_COMPILER=/usr/bin/clang++ ..
    run_exec_test make
    ctest --output-on-failure
}

function cpp_test_gen {
    print_msg "cpp_test_gen"

    cd $YDKGEN_HOME
    cpp_sanity_core_gen_install
    run_test generate.py --bundle profiles/test/ydk-models-test.json --generate-tests --cpp
    cd gen-api/cpp/models_test-bundle/build/
    run_exec_test make install

    cpp_test_gen_test
}

function cpp_tests {
    init_env "python" "python"
    cpp_sanity_core_test
    cpp_sanity_ydktest
    teardown_env
}
