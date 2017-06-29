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

PY_GENERATE="python2"
PY_TEST="python3"

function print_msg {
    echo -e "${RED}*** $(date) *** tests.sh | $1${NOCOLOR}"
}

function run_exec_test {
    $@
    local status=$?
    if [ $status -ne 0 ]; then
        exit $status
    fi
    return $status
}

function run_test_no_coverage {
    python $@
    local status=$?
    if [ $status -ne 0 ]; then
        exit $status
    fi
    return $status
}

function run_test {
    coverage run --source=ydkgen,sdk,generate --branch --parallel-mode $@ > /dev/null
    local status=$?
    if [ $status -ne 0 ]; then
        exit $status
    fi
    return $status
}

function init_env {
    print_msg "init_env"

    PY_GENERATE="$1"
    PY_TEST="$2"

    YDK_GEN_ENV=`which $PY_GENERATE`
    YDK_TEST_ENV=`which $PY_TEST`

    print_msg "init_env: Generating interpreter $YDK_GEN_ENV"
    print_msg "init_env: Testing interpreter $YDK_TEST_ENV"

    virtualenv -p $PY_GENERATE gen_env
    virtualenv -p $PY_TEST test_env

    source test_env/bin/activate
    pip install -r requirements.txt coverage > /dev/null

    source gen_env/bin/activate
    pip install -r requirements.txt coverage > /dev/null

    init_rest_server
}

function init_go_env {
    export PATH=$PATH:$GOPATH/bin
    export PATH=$PATH:$GOROOT/bin
    go get github.com/stretchr/testify
}

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

function init_rest_server {
    cd $YDKGEN_HOME/test
    print_msg "starting rest server"
    rest_server_id=$(./start_rest_server.sh)
    cd -
}

function teardown_env {
    print_msg "teardown_env"
    deactivate
    cd $YDKGEN_HOME && rm -rf gen_env test_env
}

function test_gen_tests {
    print_msg "test_gen_tests"

    init_env "python" "python"
    cd $YDKGEN_HOME && source gen_env/bin/activate
    git clone https://github.com/abhikeshav/ydk-test-yang.git sdk/cpp/core/tests/confd/testgen

    py_test_gen
    cpp_test_gen
}

########################## EXECUTION STARTS HERE #############################

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..

export YDKGEN_HOME=`pwd` && cd $YDKGEN_HOME
echo $YDKGEN_HOME

source $YDKGEN_HOME/test/test_python.sh
source $YDKGEN_HOME/test/test_cpp.sh
source $YDKGEN_HOME/test/test_go.sh

if [[ -z "$GOPATH" ]]; then
    export GOPATH=$YDKGEN_HOME/golang
fi

run_exec_test py_tests
run_exec_test cpp_tests

init_go_env
run_exec_test go_tests

# test_gen_tests

cd $YDKGEN_HOME
print_msg "gathering cpp coverage"
print_msg "combining python coverage"
coverage combine
