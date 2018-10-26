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
# Script for running ydk CI on docker via travis-ci.org
#
# ------------------------------------------------------------------

RED="\033[0;31m"
NOCOLOR="\033[0m"

function print_msg {
    echo -e "${RED}*** $(date) *** dependencies_osx.sh | $@ ${NOCOLOR}"
}

function install_dependencies {
    print_msg "Installing dependencies"

    brew install curl doxygen pcre xml2 lcov
    brew rm -f --ignore-dependencies python python3
    wget https://www.python.org/ftp/python/3.6.3/python-3.6.3-macosx10.6.pkg
    sudo installer -pkg python-3.6.3-macosx10.6.pkg  -target /
}

function install_confd {
    print_msg "Installing confd"

    wget https://github.com/CiscoDevNet/ydk-gen/files/562559/confd-basic-6.2.darwin.x86_64.zip &> /dev/null
    unzip confd-basic-6.2.darwin.x86_64.zip
    ./confd-basic-6.2.darwin.x86_64.installer.bin ../confd
}

function install_libssh {
    print_msg "Installing libssh-0.7.6"
    brew reinstall openssl
    export OPENSSL_ROOT_DIR=/usr/local/opt/openssl
    wget https://git.libssh.org/projects/libssh.git/snapshot/libssh-0.7.6.tar.gz
    tar zxf libssh-0.7.6.tar.gz && rm -f libssh-0.7.6.tar.gz
    mkdir libssh-0.7.6/build && cd libssh-0.7.6/build
    cmake ..
    sudo make install
    cd -
}

function install_golang {
    print_msg "Installing Go1.9.2"
    print_msg "Before installation GO version: $(go version)"
    bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
    source /Users/travis/.gvm/scripts/gvm
    print_msg "GO version before installation: $(go version)"
    gvm install go1.4 -B
    gvm use go1.4
    export GOROOT_BOOTSTRAP=$GOROOT
    gvm install go1.9.2
    gvm use go1.9.2
    print_msg "GOROOT: $GOROOT"
    print_msg "GOPATH: $GOPATH"
    print_msg "GO version: $(go version)"
    print_msg " "
}

########################## EXECUTION STARTS HERE #############################

install_dependencies
install_confd
install_libssh
install_golang

sudo pip install virtualenv
