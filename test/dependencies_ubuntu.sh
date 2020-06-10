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
# dependencies_ubuntu (Ubuntu 16.04, 18.04)
# ------------------------------------------------------------------

function print_msg {
    echo -e "${MSG_COLOR}*** $(date) *** dependencies_ubuntu.sh | $@ ${NOCOLOR}"
}

function run_cmd {
    local cmd=$@
    print_msg "Running: $cmd"
    $@
    local status=$?
    if [ $status -ne 0 ]; then
        MSG_COLOR=$RED
        print_msg "Exiting '$@' with status=$status"
        exit $status
    fi
    return $status
}

function install_dependencies {
    print_msg "Installing dependencies"

    apt update -y
    apt install sudo -y
    run_cmd sudo apt-get install -y --no-install-recommends apt-utils
    run_cmd sudo apt-get update -y > /dev/null
    run_cmd sudo apt-get install libtool-bin -y > /dev/null
    local status=$?
    if [[ ${status} != 0 ]]; then
        run_cmd sudo apt-get install libtool -y > /dev/null
    fi
    run_cmd sudo apt-get install -y bison curl doxygen flex git unzip wget cmake gdebi-core lcov vim > /dev/null
    run_cmd sudo apt-get install -y libcmocka0 libcurl4-openssl-dev libpcre3-dev libpcre++-dev > /dev/null
    run_cmd sudo apt-get install -y libssh-dev libxml2-dev libxslt1-dev > /dev/null
    run_cmd sudo apt-get install -y python3-dev python3-lxml python3-pip python3-venv > /dev/null
    run_cmd sudo apt-get install -y pkg-config software-properties-common zlib1g-dev openjdk-8-jre > /dev/null
    run_cmd sudo apt-get install -y valgrind > /dev/null
}

function check_install_gcc {
  which gcc
  local status=$?
  if [[ $status == 0 ]]
  then
    gcc_version=$(echo $(gcc --version) | awk '{ print $3 }' | cut -d '-' -f 1)
    print_msg "Current gcc/g++ version is $gcc_version"
  else
    print_msg "The gcc/g++ not installed"
    gcc_version="4.0"
  fi
  gcc_version=$(echo `gcc --version` | awk '{ print $3 }' | cut -d '-' -f 1)
  print_msg "Current gcc/g++ version is $gcc_version"
  if [[ $(echo $gcc_version | cut -d '.' -f 1) < 5 ]]
  then
    print_msg "Upgrading gcc/g++ to version 5"
    sudo add-apt-repository ppa:ubuntu-toolchain-r/test -y
    sudo apt-get update > /dev/null
    sudo apt-get install gcc-5 g++-5 -y > /dev/null
    sudo ln -fs /usr/bin/g++-5 /usr/bin/c++
    sudo ln -fs /usr/bin/gcc-5 /usr/bin/cc
    gcc_version=$(echo $(gcc --version) | awk '{ print $3 }' | cut -d '-' -f 1)
    print_msg "Installed gcc/g++ version is $gcc_version"
  fi
}

function check_install_go {
  go_exec=$(which go)
  if [[ -z ${go_exec} && -d /usr/local/go ]]; then
    go_exec=/usr/local/go/bin/go
  fi
  if [[ -x ${go_exec} ]]
  then
    go_version=$(echo `${go_exec} version` | awk '{ print $3 }' | cut -d 'o' -f 2)
    print_msg "Current Go version is $go_version"
    minor=$(echo $go_version | cut -d '.' -f 2)
  else
    print_msg "The Go is not installed"
    minor=0
  fi
  if (( $minor < 9 )); then
#    if (( $minor > 0 )); then
#      print_msg "Removing pre-installed Golang"
#      sudo apt-get remove golang -y
#    fi
    print_msg "Installing Golang version 1.9.2 in /usr/local/go"
    run_cmd sudo wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz &> /dev/null
    sudo tar -zxf  go1.9.2.linux-amd64.tar.gz -C /usr/local/
    rm -f go1.9.2.linux-amd64.tar.gz
  fi
}

function install_confd {
  if [[ ! -s $HOME/confd/bin/confd ]]; then
    print_msg "Installing confd"
    run_cmd wget https://github.com/CiscoDevNet/ydk-gen/files/562538/confd-basic-6.2.linux.x86_64.zip &> /dev/null
    unzip confd-basic-6.2.linux.x86_64.zip
    run_cmd ./confd-basic-6.2.linux.x86_64.installer.bin $HOME/confd
    rm -f confd-basic-6.2.* ConfD*
  fi
}

function install_fpm {
    print_msg "Installing fpm"
    apt-get install ruby ruby-dev rubygems build-essential -y > /dev/null
    gem install --no-ri --no-rdoc fpm
}

########################## EXECUTION STARTS HERE #############################

# Terminal colors
NOCOLOR="\033[0m"
YELLOW='\033[1;33m'
MSG_COLOR=$YELLOW

install_dependencies
check_install_gcc
check_install_go

install_confd

#install_fpm
