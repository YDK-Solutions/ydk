#!/bin/bash
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
#
# Script for running ydk CI on travis-ci.org
#
# ------------------------------------------------------------------

# Terminal colors
RED="\033[0;31m"
NOCOLOR="\033[0m"
YELLOW='\033[1;33m'
MSG_COLOR=$YELLOW

PY_GENERATE="python2"
PY_TEST="python3"

######################################################################
# Utility functions
######################################################################

function print_msg {
    echo -e "${MSG_COLOR}*** $(date): tests.sh | $@ ${NOCOLOR}"
}

function run_exec_test {
    $@
    local status=$?
    if [ $status -ne 0 ]; then
        MSG_COLOR=$RED
        print_msg "Exiting '$@' with status=$status"
        exit $status
    fi
    return $status
}

function run_test_no_coverage {
    print_msg "Executing: ${PYTHON_BIN} $@"
    ${PYTHON_BIN} $@
    local status=$?
    if [ $status -ne 0 ]; then
        MSG_COLOR=$RED
        print_msg "Exiting '${PYTHON_BIN} $@' with status=$status"
        exit $status
    fi
    return $status
}

function run_test {
    if [[ $(command -v coverage) && $run_with_coverage ]]; then
        print_msg "Executing with coverage: $@"
        coverage run --omit=/usr/* --branch --parallel-mode $@ > /dev/null
        local status=$?
        print_msg "Returned status is ${status}"
        if [ $status -ne 0 ]; then
            MSG_COLOR=$RED
            print_msg "Exiting 'coverage run $@' with status=$status"
            exit $status
        fi
        return $status
    fi
    run_test_no_coverage $@
    local status=$?
    return $status
}

function pip_check_install {
    if [[ $(uname) == "Linux" && ${os_info} == *"fedora"* ]]
    then
        print_msg "Custom pip install of $@ for CentOS"
        ${PIP_BIN} install --install-option="--install-purelib=/usr/lib64/python${PYTHON_VERSION}/site-packages" --no-deps $@ -U
    else
        ${PIP_BIN} install $@ -U
    fi
}

######################################################################
# Environment setup-teardown functions
######################################################################

function init_confd {
    cd $1
    print_msg "Initializing confd in $(pwd)"
    source $YDKGEN_HOME/../confd/confdrc
    run_exec_test make stop > /dev/null
    run_exec_test make clean > /dev/null
    run_exec_test make all > /dev/null
    run_exec_test make start
    cd -
}

function init_confd_ydktest {
    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/ydktest
}

function init_rest_server {
    print_msg "Starting REST server"
    pkill -f moco-runner
    export REST_SERVER_PID=$(./test/start_rest_server.sh)
    print_msg "REST server started with PID $REST_SERVER_PID"
}

function init_tcp_server {
    print_msg "Starting TCP server"
    pkill -f tcp_proxy_server
    export TCP_SERVER_PID=$(./test/start_tcp_server.sh)
    print_msg "TCP server started with PID: $TCP_SERVER_PID"
}

function stop_tcp_server {
    print_msg "Stopping TCP server with PID: $TCP_SERVER_PID"
    kill $TCP_SERVER_PID
}

function check_python_installation {
  
  if [[ ${os_type} == "Darwin" ]] ; then
    PYTHON_VERSION=3
    PYTHON_BIN=python3
    PIP_BIN=pip3
  else
    PYTHON_BIN=python${PYTHON_VERSION}
    if [[ ${PYTHON_VERSION} = *"2"* ]]; then
      PIP_BIN=pip
    elif [[ ${PYTHON_VERSION} = *"3.5"* ]]; then
      PIP_BIN=pip3
    else
      PIP_BIN=pip${PYTHON_VERSION}
    fi
  fi

  print_msg "Checking installation of ${PYTHON_BIN}"
  ${PYTHON_BIN} --version &> /dev/null
  status=$?
  if [ $status -ne 0 ]; then
    MSG_COLOR=$RED
    print_msg "Could not locate ${PYTHON_BIN}"
    exit $status
  fi
  print_msg "Checking installation of ${PIP_BIN}"
  ${PIP_BIN} -V &> /dev/null
  status=$?
  if [ $status -ne 0 ]; then
    MSG_COLOR=$RED
    print_msg "Could not locate ${PIP_BIN}"
    exit $status
  fi
  print_msg "Python location: $(which ${PYTHON_BIN})"
  print_msg "Pip location: $(which ${PIP_BIN})"
}

function init_py_env {
  check_python_installation
  print_msg "Initializing Python requirements"
  sudo ${PIP_BIN} install -r requirements.txt pybind11==2.2.2
  if [[ $run_with_coverage ]] ; then
    sudo ${PIP_BIN} install coverage
  fi

  #else
    #print_msg "Initializing Python3 virtual environment"
    #virtualenv macos_pyenv -p python3
    #source macos_pyenv/bin/activate
  #fi
}

function init_go_env {
    print_msg "Initializing Go environment"

    if [[ $(uname) == "Darwin" ]]; then
        source /Users/travis/.gvm/scripts/gvm
        gvm use go1.9.2
        print_msg "GOROOT: $GOROOT"
        print_msg "GOPATH: $GOPATH"
    else
        cd $YDKGEN_HOME
        export GOPATH="$(pwd)/golang"
        export GOROOT=/usr/local/go
        export PATH=$GOROOT/bin:$PATH
        print_msg "Setting GOROOT to $GOROOT"
        print_msg "Setting GOPATH to $GOPATH"
    fi
    print_msg "Running $(go version)"

    go get github.com/stretchr/testify
}

######################################################################
# Core install / test functions
######################################################################

function install_test_cpp_core {
    print_msg "Installing / testing C++ core"
    install_cpp_core
    run_cpp_core_test

    if [[ $(uname) == "Linux" ]]; then
        generate_libydk
    fi
}

function install_cpp_core {
    print_msg "Installing C++ core library"

    cd $YDKGEN_HOME
    mkdir -p $YDKGEN_HOME/sdk/cpp/core/build
    cd $YDKGEN_HOME/sdk/cpp/core/build

    if [[ $run_with_coverage ]] ; then
      print_msg "Compiling with coverage"
      run_exec_test ${CMAKE_BIN} -DCOVERAGE=True ..
    else
      run_exec_test ${CMAKE_BIN} ..
    fi
    run_exec_test make > /dev/null
    sudo make install
}

function generate_libydk {
    print_msg "Generating libydk package for later testing"
    cd $YDKGEN_HOME
    run_test generate.py --libydk
    cd gen-api/cpp/ydk/build
    sudo make package || true
    cp libydk*rpm libydk*deb /ydk-gen &> /dev/null
    cd -
}

function run_cpp_core_test {
    print_msg "Running C++ core test"

    #export CTEST_OUTPUT_ON_FAILURE=1
    #make test

    ./tests/ydk_core_test -d yes
    local status=$?
    if [ $status -ne 0 ]; then
        # If the tests fail, try to run them in verbose to get more details for debug
        #./tests/ydk_core_test -d yes
        MSG_COLOR=$RED
        print_msg "Exiting 'run_cpp_core_test' with status=$status"
        exit $status
    fi
    cd $YDKGEN_HOME
}

function install_go_core {
    print_msg "Installing Go YDK core"
    cd $YDKGEN_HOME
    run_test generate.py --core --go --install
}

function install_py_core {
    print_msg "Installing Python YDK core"

    if [[ $run_with_coverage ]] ; then
      print_msg "Building python with coverage"
      export YDK_COVERAGE=
    fi
    cd $YDKGEN_HOME/sdk/python/core
    ${PYTHON_BIN} setup.py sdist
    ${PIP_BIN} install dist/ydk*.tar.gz

#    print_msg "Generating py binaries"
#    sudo ./generate_python_binary.sh

    cd $YDKGEN_HOME
}

######################################################################
# C++ ydktest bundle install and test functions
######################################################################

function run_cpp_bundle_tests {
    print_msg "Generating and testing C++ bundle"

    cpp_sanity_ydktest_gen_install
    cpp_sanity_ydktest_test
    collect_cpp_coverage
}

function generate_install_specified_cpp_bundle {
   bundle_profile=$1
   bundle_name=$2

   run_test generate.py --bundle ${bundle_profile} --cpp --generate-doc &> /dev/null
   cd gen-api/cpp/${bundle_name}/build
   run_exec_test make > /dev/null
   sudo make install
   cd -
}

function cpp_sanity_ydktest_gen_install {
    print_msg "Generating and installing C++ ydktest bundle"
    generate_install_specified_cpp_bundle profiles/test/ydktest-cpp.json ydktest-bundle

    print_msg "Generating and installing new C++ ydktest bundle"
    generate_install_specified_cpp_bundle profiles/test/ydktest-cpp-new.json ydktest_new-bundle
}

function cpp_sanity_ydktest_test {
    print_msg "Running C++ bundle tests"

    print_msg "Initializing ssh keys for key-based authentication"
    sudo mkdir -p /var/confd/homes/admin/.ssh
    sudo touch /var/confd/homes/admin/.ssh/authorized_keys
    cd $YDKGEN_HOME
    sudo sh -c 'cat sdk/cpp/tests/ssh_host_rsa_key.pub >> /var/confd/homes/admin/.ssh/authorized_keys'
    cd -

    print_msg "Building and running C++ bundle tests"
    mkdir -p $YDKGEN_HOME/sdk/cpp/tests/build && cd sdk/cpp/tests/build
    if [[ $run_with_coverage ]] ; then
      print_msg "Compiling with coverage"
      run_exec_test ${CMAKE_BIN} -DCOVERAGE=True  ..
    else
      run_exec_test ${CMAKE_BIN} ..
    fi
    run_exec_test make > /dev/null

    #export CTEST_OUTPUT_ON_FAILURE=1
    #make test

    ./ydk_bundle_test -d yes
    local status=$?
    if [ $status -ne 0 ]; then
        # If the tests fail, try to run them in verbose to get more details for  # debug
        #./ydk_bundle_test -d yes
        MSG_COLOR=$RED
        print_msg "Exiting C++ bundle tests with status=$status"
        exit $status
    fi
}

function cpp_test_gen_test {
    print_msg "cpp_test_gen_test"

    cd $YDKGEN_HOME
    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/testgen/confd
    mkdir -p gen-api/cpp/models_test-bundle/ydk/models/models_test/test/build
    cd gen-api/cpp/models_test-bundle/ydk/models/models_test/test/build
    run_exec_test ${CMAKE_BIN} ..
    run_exec_test make > /dev/null
    ctest --output-on-failure
}

function cpp_test_gen {
    print_msg "cpp_test_gen"

    cd $YDKGEN_HOME
    run_test generate.py --bundle profiles/test/ydk-models-test.json --generate-tests --cpp &> /dev/null
    cd gen-api/cpp/models_test-bundle/build/
    run_exec_test make > /dev/null
    sudo make install

    # cpp_test_gen_test
}

function collect_cpp_coverage {
  if [[ $run_with_coverage ]] ; then
    print_msg "Collecting coverage for C++"
    cd ${YDKGEN_HOME}/sdk/cpp/core/build
    lcov --directory . --capture --output-file coverage.info &> /dev/null # capture coverage info
    lcov --remove coverage.info '/usr/*' '/Applications/*' '/opt/*' '*/json.hpp' '*/catch.hpp' '*/network_topology.cpp' '*/spdlog/*' --output-file coverage.info # filter out system
    lcov --list coverage.info #debug info
    print_msg "Moving C++ coverage to ${YDKGEN_HOME}"
    cp coverage.info ${YDKGEN_HOME}
  fi
}

######################################################################
# Go ydktest bundle install and test functions
######################################################################

function run_go_bundle_tests {
    print_msg "Generating/installing go sanity bundle tests"
    # TODO: go get
    cd $YDKGEN_HOME
    run_test generate.py -i --bundle profiles/test/ydktest-cpp.json --go

    run_go_tests
}

function run_go_tests {
    print_msg "Running go tests"
    run_go_samples
    run_go_sanity_tests
}

function run_go_samples {
    print_msg "Running go samples"

    export CXX=/usr/bin/c++
    export CC=/usr/bin/cc

    print_msg "CC: ${CC}"
    print_msg "CXX: ${CXX}"

    cd $YDKGEN_HOME/sdk/go/core/samples
    run_exec_test go run cgo_path/cgo_path.go
    run_exec_test go run bgp_create/bgp_create.go -device ssh://admin:admin@localhost:12022
    run_exec_test go run bgp_read/bgp_read.go -device ssh://admin:admin@localhost:12022
    run_exec_test go run bgp_delete/bgp_delete.go -device ssh://admin:admin@localhost:12022
    cd -
}

function run_go_sanity_tests {
    print_msg "Running go sanity tests"
    cd $YDKGEN_HOME/sdk/go/core/tests
    if [[ $run_with_coverage ]] ; then
        run_exec_test go test -race -coverpkg="github.com/CiscoDevNet/ydk-go/ydk/providers","github.com/CiscoDevNet/ydk-go/ydk/services","github.com/CiscoDevNet/ydk-go/ydk/types","github.com/CiscoDevNet/ydk-go/ydk/types/datastore","github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format","github.com/CiscoDevNet/ydk-go/ydk/types/protocol","github.com/CiscoDevNet/ydk-go/ydk/types/yfilter","github.com/CiscoDevNet/ydk-go/ydk/types/ytype","github.com/CiscoDevNet/ydk-go/ydk","github.com/CiscoDevNet/ydk-go/ydk/path" -coverprofile=coverage.txt -covermode=atomic
        print_msg "Moving go coverage to ${YDKGEN_HOME}"
        mv coverage.txt ${YDKGEN_HOME}
    else
        run_exec_test go test
    fi
    cd -
}

######################################################################
# Python ydktest bundle install and test functions
######################################################################

function run_python_bundle_tests {
    print_msg "Running python bundle tests"
    py_sanity_ydktest
    py_sanity_deviation
    py_sanity_augmentation
    py_sanity_common_cache
    py_sanity_one_class_per_module
    py_sanity_backward_compatibility
}

#--------------------------
# Python ydktest bundle
#--------------------------

function py_sanity_ydktest {
    print_msg "Generating, installing and testing python ydktest bundle"

    py_sanity_ydktest_gen
    if [[ $run_with_coverage ]]
    then
        py_sanity_ydktest_test
        py_sanity_ydktest_install
    else
        py_sanity_ydktest_install
        run_py_sanity_ydktest_tests
    fi
}

function py_sanity_ydktest_gen {
    print_msg "Generating python ydk core and ydktest bundle"

    cd $YDKGEN_HOME

    print_msg "py_sanity_ydktest_gen: testing bundle and documentation generation"
    run_test generate.py --bundle profiles/test/ydktest-cpp.json --python --generate-doc &> /dev/null

    print_msg "py_sanity_ydktest_gen: testing core and documentation generation"
    run_test generate.py --core
}

function py_sanity_ydktest_install {
    print_msg "py_sanity_ydktest_install"
    print_msg "Installing"
    cd $YDKGEN_HOME
    pip_check_install gen-api/python/ydktest-bundle/dist/ydk*.tar.gz

    print_msg "Running import tests on generated bundle"
    run_test gen-api/python/ydktest-bundle/ydk/models/ydktest/test/import_tests.py

}

function py_sanity_ydktest_test {
    print_msg "Running py_sanity_ydktest_test with coverage"
    cd $YDKGEN_HOME
    cp -r gen-api/python/ydktest-bundle/ydk/models/* sdk/python/core/ydk/models

    print_msg "Uninstall ydk py core from pip for testing with coverage"
    ${PIP_BIN} uninstall ydk -y
    
    export OLDPYTHONPATH=$PYTHONPATH
    print_msg "Setting OLDPYTHONPATH to $OLDPYTHONPATH"

    print_msg "Build & copy cpp-wrapper to sdk directory to gather coverage"
    cd sdk/python/core/ && ${PYTHON_BIN} setup.py build

    print_msg "Set new python path to gather coverage"
    export PYTHONPATH=$PYTHONPATH:$(pwd)
    print_msg "Setting PYTHONPATH to $PYTHONPATH"
    cp build/lib*/*.so .
    cd -

    run_py_sanity_ydktest_tests

    export PYTHONPATH=$OLDPYTHONPATH
    print_msg "Restored PYTHONPATH to $PYTHONPATH"

    cd sdk/python/core/
    rm -f *.so
    print_msg "Restore ydk py core to pip"
    ${PIP_BIN} install dist/ydk*.tar.gz

    cd $YDKGEN_HOME
}

function run_py_sanity_ydktest_tests {
    print_msg "Running run_py_sanity_ydktest_tests"
    run_test sdk/python/core/tests/test_ydk_types.py
    run_test sdk/python/core/tests/test_sanity_codec.py

    py_sanity_ydktest_test_netconf_ssh
    py_sanity_ydktest_test_tcp

    stop_tcp_server
}

function py_sanity_ydktest_test_netconf_ssh {
    print_msg "py_sanity_ydktest_test_netconf_ssh"

    run_test sdk/python/core/tests/test_netconf_operations.py
    run_test sdk/python/core/tests/test_opendaylight.py
    run_test sdk/python/core/tests/test_restconf_provider.py
    run_test sdk/python/core/tests/test_sanity_delete.py
    run_test sdk/python/core/tests/test_sanity_errors.py
    run_test sdk/python/core/tests/test_sanity_filter_read.py
    run_test sdk/python/core/tests/test_sanity_filters.py
    run_test sdk/python/core/tests/test_sanity_levels.py
    run_test sdk/python/core/tests/test_sanity_netconf.py
    run_test sdk/python/core/tests/test_sanity_path.py
    run_test sdk/python/core/tests/test_netconf_provider.py
    run_test sdk/python/core/tests/test_sanity_service_errors.py
    run_test sdk/python/core/tests/test_sanity_type_mismatch_errors.py
    run_test sdk/python/core/tests/test_sanity_types.py
    run_test sdk/python/core/tests/test_non_top_operations.py
#    run_test_no_coverage sdk/python/core/tests/test_sanity_executor_rpc.py

    print_msg "py_sanity_ydktest_test_netconf_ssh no on-demand"
    run_test sdk/python/core/tests/test_netconf_operations.py --non-demand
    run_test sdk/python/core/tests/test_sanity_delete.py --non-demand
    run_test sdk/python/core/tests/test_sanity_errors.py --non-demand
    run_test sdk/python/core/tests/test_sanity_filter_read.py --non-demand
    run_test sdk/python/core/tests/test_sanity_filters.py --non-demand
    run_test sdk/python/core/tests/test_sanity_levels.py --non-demand
    run_test sdk/python/core/tests/test_sanity_netconf.py --non-demand
    run_test sdk/python/core/tests/test_sanity_path.py --non-demand
    run_test sdk/python/core/tests/test_netconf_provider.py --non-demand
    run_test sdk/python/core/tests/test_sanity_service_errors.py --non-demand
    run_test sdk/python/core/tests/test_sanity_type_mismatch_errors.py --non-demand
    run_test sdk/python/core/tests/test_sanity_types.py --non-demand
#    run_test_no_coverage sdk/python/core/tests/test_sanity_executor_rpc.py
#    --non-demand
}

function py_sanity_ydktest_test_tcp {
    print_msg "py_sanity_ydktest_test_tcp"
    run_test sdk/python/core/tests/test_sanity_netconf.py tcp://admin:admin@127.0.0.1:12307
    init_confd_ydktest
    run_test sdk/python/core/tests/test_sanity_netconf.py tcp://admin:admin@127.0.0.1:12307 --non-demand
}

#--------------------------
# Python deviation bundle
#--------------------------

function py_sanity_deviation {
    print_msg "py_sanity_deviation"

    rm -rf $HOME/.ydk/*
    py_sanity_deviation_ydktest_gen
    py_sanity_deviation_ydktest_install
    py_sanity_deviation_ydktest_test

    py_sanity_deviation_bgp_gen
    py_sanity_deviation_bgp_install
    py_sanity_deviation_bgp_test
    rm -rf $HOME/.ydk/*
}

function py_sanity_deviation_ydktest_gen {
    print_msg "py_sanity_deviation_ydktest_gen"

    rm -rf gen-api/python/*
    cd $YDKGEN_HOME
    run_test generate.py --bundle profiles/test/ydktest-cpp.json --python
}

function py_sanity_deviation_ydktest_install {
    print_msg "py_sanity_deviation_ydktest_install"

    ${PIP_BIN} uninstall ydk-models-ydktest -y
    pip_check_install gen-api/python/ydktest-bundle/dist/ydk*.tar.gz
}

function py_sanity_deviation_ydktest_test {
    print_msg "py_sanity_deviation_ydktest_test"

    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/deviation
    run_test sdk/python/core/tests/test_sanity_deviation.py
    run_test sdk/python/core/tests/test_sanity_deviation.py --non-demand
}

function py_sanity_deviation_bgp_gen {
    print_msg "py_sanity_deviation_bgp_gen"

    rm -rf gen-api/python/*
    cd $YDKGEN_HOME
    run_test generate.py --bundle profiles/test/deviation.json --verbose
}

function py_sanity_deviation_bgp_install {
    print_msg "py_sanity_deviation_bgp_install"

    cd $YDKGEN_HOME
    pip_check_install gen-api/python/deviation-bundle/dist/*.tar.gz
}

function py_sanity_deviation_bgp_test {
    print_msg "py_sanity_deviation_bgp_test"

    run_test sdk/python/core/tests/test_sanity_deviation_bgp.py
    run_test sdk/python/core/tests/test_sanity_deviation_bgp.py --non-demand
}

#--------------------------
# Python augmentation bundle
#--------------------------
function py_sanity_augmentation {
    print_msg "py_sanity_augmentation"

    rm -rf $HOME/.ydk/*
    py_sanity_augmentation_gen
    py_sanity_augmentation_install
    py_sanity_augmentation_test
    rm -rf $HOME/.ydk/*
}

function py_sanity_augmentation_gen {
    print_msg "py_sanity_augmentation_gen"

    cd $YDKGEN_HOME && rm -rf gen-api/python/*
    run_test generate.py --core
    run_test generate.py --bundle profiles/test/ydktest-augmentation.json
}

function py_sanity_augmentation_install {
    print_msg "py_sanity_augmentation_install"

    cd $YDKGEN_HOME
    ${PIP_BIN} uninstall ydk -y
    ${PIP_BIN} install gen-api/python/ydk/dist/ydk*.tar.gz
    pip_check_install gen-api/python/augmentation-bundle/dist/*.tar.gz
}

function py_sanity_augmentation_test {
    print_msg "py_sanity_augmentation_test"

    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/augmentation

    run_test sdk/python/core/tests/test_sanity_augmentation.py
    run_test sdk/python/core/tests/test_sanity_augmentation.py --non-demand
    run_test sdk/python/core/tests/test_on_demand.py
}

function py_sanity_common_cache {
    print_msg "py_sanity_common_cache"

    rm -rf $HOME/.ydk/*
    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/deviation
    run_test sdk/python/core/tests/test_sanity_deviation.py --common-cache
    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/augmentation
    run_test sdk/python/core/tests/test_sanity_augmentation.py --common-cache

    run_test sdk/python/core/tests/test_sanity_levels.py --common-cache
    run_test sdk/python/core/tests/test_sanity_types.py --common-cache
    rm -rf $HOME/.ydk/*
}

function py_sanity_run_limited_tests {
    ${PIP_BIN} uninstall ydk-models-ydktest -y
    pip_check_install gen-api/python/ydktest-bundle/dist/ydk*.tar.gz
    print_msg "Running limited test : run_test sdk/python/core/tests/test_sanity_levels.py"
    run_test sdk/python/core/tests/test_sanity_levels.py
    print_msg "Running limited test : sdk/python/core/tests/test_sanity_types.py"
    run_test sdk/python/core/tests/test_sanity_types.py
    print_msg "Running limited test : sdk/python/core/tests/test_sanity_codec.py"
    run_test sdk/python/core/tests/test_sanity_codec.py
}

function py_sanity_one_class_per_module {
    print_msg "**********************************"
    print_msg "Running ONE CLASS PER MODULE TESTS"
    cd $YDKGEN_HOME
    run_test generate.py --bundle profiles/test/ydktest-cpp.json -o > /dev/null
    py_sanity_run_limited_tests
}

function py_sanity_backward_compatibility {
    print_msg "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
    print_msg "Running YDK BUNDLE BACKWARD COMPATIBILITY TESTS"
    cd $YDKGEN_HOME
    CURRENT_GIT_REV=$(git rev-parse HEAD)
    git checkout 8d278d7d51765a310afe43ed88951f2fe5d8783e
    run_test generate.py --bundle profiles/test/ydktest-cpp.json > /dev/null
    py_sanity_run_limited_tests
    git checkout ${CURRENT_GIT_REV}
}

#-------------------------------------
# Python generated model tests bundle
#-------------------------------------

#function test_gen_tests {
#    print_msg "test_gen_tests"
#
#    cd $YDKGEN_HOME
#    git clone https://github.com/psykokwak4/ydk-test-yang.git sdk/cpp/core/tests/confd/testgen
#
#    py_test_gen
#    cpp_test_gen
#}

#function py_test_gen_test {
#    print_msg "py_test_gen_test"
#
#    cd $YDKGEN_HOME
#    init_confd $YDKGEN_HOME/sdk/cpp/core/tests/confd/testgen/confd
#    cd gen-api/python/models_test-bundle/ydk/models/models_test/test/
#    ${PYTHON_BIN} import_tests.py
#    cd models_test/
#    ${PYTHON_BIN} -m unittest discover
#}

function py_test_gen {
    print_msg "py_test_gen"

    cd $YDKGEN_HOME
    run_test generate.py --core --python
    run_test generate.py --bundle profiles/test/ydk-models-test.json  --generate-tests --python &> /dev/null
    ${PIP_BIN} install gen-api/python/ydk/dist/ydk*.tar.gz
    pip_check_install gen-api/python/models_test-bundle/dist/ydk*.tar.gz

    # py_test_gen_test
}

#-------------------------------------
# Meta-data test
#-------------------------------------

function run_py_metadata_test {
    print_msg "Building ydktest bundle with metadata"
    cd $YDKGEN_HOME
    ${PYTHON_BIN} generate.py --bundle profiles/test/ydktest-cpp.json --generate-meta
    pip_check_install gen-api/python/ydktest-bundle/dist/ydk*.tar.gz

    print_msg "Running metadata test"
    ${PYTHON_BIN} sdk/python/core/tests/test_meta.py
}

#-------------------------------------
# Documentation tests
#-------------------------------------

function sanity_doc_gen {
   print_msg "Generating docs"
   print_msg "Removing previously generated code with sudo to avoid permission issues"
   sudo rm -rf gen-api/cpp/ydk
   run_test generate.py --core --cpp --generate-doc &> /dev/null
   run_test generate.py --core --go --generate-doc &> /dev/null
   run_test generate.py --core --python --generate-doc &> /dev/null
}

function sanity_doc_gen_cache {
   print_msg "Creating cache and moving previous docs to cache"
   mkdir -p gen-api/cache
   mv gen-api/python gen-api/cache

   print_msg "Generating bundle"
   run_test generate.py --bundle profiles/test/ydktest-cpp.json

   print_msg "Generating core docs with cache option"
   run_test generate.py --core --python --generate-doc --cached-output-dir --output-directory gen-api
}

########################## EXECUTION STARTS HERE #############################
######################################
# Parse args to get python version
######################################

args=$(getopt p:d $*)
set -- $args
PYTHON_VERSION=${2}

######################################
# Set up env
######################################

os_type=$(uname)
if [[ ${os_type} == "Linux" ]] ; then
    os_info=$(cat /etc/*-release)
else
    os_info=$(sw_vers)
fi
print_msg "Running OS type: $os_type"
print_msg "OS info: $os_info"
if [[ ${os_type} == "Linux" && ${os_info} != *"trusty"* ]] ; then
    run_with_coverage=1
fi

export YDKGEN_HOME="$(pwd)"
print_msg "YDKGEN_HOME is set to: ${YDKGEN_HOME}"

CMAKE_BIN=cmake
which cmake3
status=$?
if [[ ${status} == 0 ]] ; then
    CMAKE_BIN=cmake3
fi

init_py_env
init_confd_ydktest
init_rest_server
init_tcp_server

######################################
# Install and run C++ core tests
######################################
install_test_cpp_core
run_cpp_bundle_tests

######################################
# Install and run Go tests
######################################
init_go_env
install_go_core
run_go_bundle_tests

######################################
# Install and run Python tests
######################################
install_py_core
run_python_bundle_tests
run_py_metadata_test

# test_gen_tests

######################################
# Documentation tests
######################################
sanity_doc_gen
sanity_doc_gen_cache

cd $YDKGEN_HOME
find . -name '*gcda*'|xargs rm -f
find . -name '*gcno*'|xargs rm -f
find . -name '*gcov*'|xargs rm -f

if [[ $run_with_coverage ]] ; then
    print_msg "Combining C++, Python and Go coverage"
    coverage combine > /dev/null || echo "Coverage not combined"
fi
