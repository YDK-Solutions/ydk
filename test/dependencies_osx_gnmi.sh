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
# dependencies_osx_gnmi.sh
# Script for running ydk CI on docker via travis-ci.org
#
# ------------------------------------------------------------------

function print_msg {
    echo -e "${MSG_COLOR}*** $(date) *** dependencies_osx_gnmi.sh | $@ ${NOCOLOR}"
}

function install_protobuf {
    print_msg "Downloading protobuf and protoc"
    wget https://github.com/google/protobuf/releases/download/v3.5.0/protobuf-cpp-3.5.0.zip > /dev/null
    unzip protobuf-cpp-3.5.0.zip > /dev/null
    cd protobuf-3.5.0
    print_msg "Configuring protobuf and protoc"
    ./configure > /dev/null
    print_msg "Compiling protobuf and protoc"
    make > /dev/null
    print_msg "Installing protobuf and protoc"
    sudo make install
    cd $curr_dir
}

function install_grpc {
    print_msg "Installing grpc"

    #LIBTOOL=glibtool LIBTOOLIZE=glibtoolize make
    git clone -b v1.9.1 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    make > /dev/null
    sudo make install
    cd $curr_dir
}

########################## EXECUTION STARTS HERE #############################

# Terminal colors
NOCOLOR="\033[0m"
YELLOW='\033[1;33m'
MSG_COLOR=$YELLOW

curr_dir="$(pwd)"
install_protobuf
install_grpc

