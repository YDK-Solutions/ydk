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
# dependencies_osx.sh
# Script for running ydk CI on docker via travis-ci.org
#
# ------------------------------------------------------------------

function print_msg {
    echo -e "${MSG_COLOR}*** $(date) *** dependencies_osx.sh | $@ ${NOCOLOR}"
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

function install_confd {
    print_msg "Installing confd"

    wget https://github.com/CiscoDevNet/ydk-gen/files/562559/confd-basic-6.2.darwin.x86_64.zip &> /dev/null
    unzip confd-basic-6.2.darwin.x86_64.zip
    ./confd-basic-6.2.darwin.x86_64.installer.bin ../confd
}

function install_fpm {
    print_msg "Installing fpm"
    brew install gnu-tar > /dev/null
    gem install --no-ri --no-rdoc fpm
}

function install_golang {
    print_msg "Installing Go1.9.2"
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

function check_python_installation {
  print_msg "Checking python3 and pip3 installation"
  #brew rm -f --ignore-dependencies python python3
  python3 -V &> /dev/null
  status=$?
  if [ $status -ne 0 ]; then
    print_msg "Installing python3"
    brew install python@3
  fi
  pip3 -V &> /dev/null
  status=$?
  if [ $status -ne 0 ]; then
    print_msg "Installing pip3"
    sudo easy_install pip3
  fi
  sudo pip3 install virtualenv
}

########################## EXECUTION STARTS HERE #############################

# Terminal colors
NOCOLOR="\033[0m"
YELLOW='\033[1;33m'
MSG_COLOR=$YELLOW

install_libssh
install_confd
install_golang

brew install pybind11
check_python_installation

#install_fpm
