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
    echo -e "${RED}*** $(date) *** dependencies_linux.sh | $1${NOCOLOR}"
}

function install_confd {
    print_msg "Installing confd"

    wget https://github.com/CiscoDevNet/ydk-gen/files/562538/confd-basic-6.2.linux.x86_64.zip &> /dev/null
    unzip confd-basic-6.2.linux.x86_64.zip
    ./confd-basic-6.2.linux.x86_64.installer.bin ../confd
}

function install_fpm {
    print_msg "Installing fpm"
    apt-get install ruby ruby-dev rubygems build-essential -y > /dev/null
    gem install --no-ri --no-rdoc fpm
}

function install_protobuf {
    print_msg "Installing protobuf and protoc"

    wget https://github.com/google/protobuf/releases/download/v3.3.0/protobuf-cpp-3.3.0.zip
    unzip protobuf-cpp-3.3.0.zip
    cd protobuf-3.3.0
    ./configure
    make
    make check
    sudo make install
    sudo ldconfig
    cd -
}

function install_grpc {
    print_msg "Installing grpc"

    git clone -b 1.4.5 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    sudo ldconfig
    make
    sudo make install
    cd -
}

########################## EXECUTION STARTS HERE #############################

./test/dependencies_ubuntu_basic.sh
install_confd

#install_fpm
install_protobuf
install_grpc
