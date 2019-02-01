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
# dependencies_ubuntu (Ubuntu 16.04)
# ------------------------------------------------------------------

function print_msg {
    echo -e "${MSG_COLOR}*** $(date) *** dependencies_ubuntu.sh | $@ ${NOCOLOR}"
}

function install_dependencies {
    print_msg "Installing dependencies"

    apt update -y > /dev/null
    apt install sudo -y > /dev/null
    sudo apt-get update > /dev/null
    sudo apt-get install libtool-bin -y > /dev/null
    local status=$?
    if [[ ${status} != 0 ]]; then
        sudo apt-get install libtool -y > /dev/null
    fi
    sudo apt-get install -y bison \
                            curl \
                            doxygen \
                            flex \
                            git \
                            libcmocka0 \
                            libcurl4-openssl-dev \
                            libpcre3-dev \
                            libpcre++-dev \
                            libssh-dev \
                            libxml2-dev \
                            libxslt1-dev \
                            pkg-config \
                            python-dev \
                            python-pip \
                            python3-dev \
                            python-lxml \
                            python3-lxml \
                            python3-pip \
                            python-virtualenv \
                            software-properties-common \
                            unzip \
                            wget \
                            zlib1g-dev\
                            cmake \
                            openjdk-8-jre \
                            gdebi-core\
                            lcov > /dev/null
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

function install_go {
    print_msg "Removing pre-installed Golang"
    sudo apt-get remove golang -y
    print_msg "Installing Golang version 1.9.2"
    sudo wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz &> /dev/null
    sudo tar -zxf  go1.9.2.linux-amd64.tar.gz -C /usr/local/
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

########################## EXECUTION STARTS HERE #############################

# Terminal colors
NOCOLOR="\033[0m"
YELLOW='\033[1;33m'
MSG_COLOR=$YELLOW

install_dependencies
check_install_gcc
install_go

install_confd

#install_fpm
