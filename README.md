[![Codacy Badge](https://api.codacy.com/project/badge/Grade/6e111527081b48e1b2252c3562e08a3b)](https://www.codacy.com/app/ydk/ydk-gen?utm_source=github.com&utm_medium=referral&utm_content=CiscoDevNet/ydk-gen&utm_campaign=badger)
[![License](https://cloud.githubusercontent.com/assets/17089095/19458582/dd626d2c-9481-11e6-8019-8227c5c66a06.png)](https://github.com/CiscoDevNet/ydk-gen/blob/master/LICENSE) [![Build Status](https://travis-ci.org/CiscoDevNet/ydk-gen.svg?branch=master)](https://travis-ci.org/CiscoDevNet/ydk-gen)
[![codecov](https://codecov.io/gh/CiscoDevNet/ydk-gen/branch/master/graph/badge.svg)](https://codecov.io/gh/CiscoDevNet/ydk-gen)
[![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/ydkdev/ydk-gen/)

![ydk-logo-128](https://cloud.githubusercontent.com/assets/16885441/24175899/2010f51e-0e56-11e7-8fb7-30a9f70fbb86.png)

YANG Development Kit (Generator)
================================

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Overview](#overview)
- [Backward compatibility](#backward-compatibility)
- [Docker](#docker)
- [System requirements](#system-requirements)
  - [Linux](#linux)
  - [MacOS](#macos)
  - [Windows](#windows)
- [Installation](#installation)
  - [Setting up your environment](#setting-up-your-environment)
  - [Clone ydk-gen and install the requirements](#clone-ydk-gen-and-install-the-requirements)
- [Generate YDK components](#generate-ydk-components)
  - [First step: choose model bundle profile](#first-step-choose-model-bundle-profile)
  - [Second step. Generate and install the core](#second-step-generate-and-install-the-core)
  - [Third step. Generate and install your bundle](#third-step-generate-and-install-model-bundle)
  - [Fourth step: Writing your first app](#fourth-step-writing-your-first-app)
  - [Documentation](#documentation)
- [Generating an "Adhoc" YDK-Py Bundle](#generating-an-adhoc-ydk-py-bundle)
- [Notes](#notes)
  - [Python Requirements](#python-requirements)
  - [Directory structure](#directory-structure)
  - [Troubleshooting](#troubleshooting)
- [Running Unit Tests](#running-unit-tests)
  - [Python](#python)
  - [C++](#c)
  - [Go](#go)
- [Support](#support)
- [Release Notes](#release-notes)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Overview

**YDK** is a developer tool that can generate API's that are modeled in YANG. Currently, it generates language binding for Python, Go and C++ with planned support for other language bindings in the future.

Other tools and libraries are used to deliver `YDK` functionality:

* YANG model analysis and code generation is implemented using APIs from the [pyang](https://github.com/mbj4668/pyang) library
* Documentation is generated using [Sphinx](http://www.sphinx-doc.org/en/stable/)
* Run time yang model analysis is done using [libyang](https://github.com/CESNET/libyang)
* C++ to python bindings are created using [pybind11](https://github.com/pybind/pybind11)
* C++ uses [catch](https://github.com/catchorg/Catch2) and [spdlog](https://github.com/gabime/spdlog) for tests and logging respectively

Of course, many other libraries are used as an integral part of ydk-gen and it's dependencies, too many to mention!

The output of ydk-gen is either a core package, that defines services and providers, or a module bundle, consisting of APIs based on YANG models. Each module bundle is generated using a bundle profile and the ydk-gen tool. Developers can either use pre-packaged generated bundles (e.g. [ydk-py](http://cs.co/ydk-py)), or they can define their own bundle, consisting of a set of YANG models, using a bundle profile (e.g. [```ietf_0_1_1.json```](profiles/bundles/ietf_0_1_1.json)). This gives a developer the ability to customize the scope of their bundle based on their requirements.


# Backward compatibility

The YDK-0.8.1 core is bacward compatible with generated in YDK-0.7.3 model bundle code. It is not compatible with YDK-0.7.2 and earlier bundle packages due to changes in modeling and handling of YList objects.

# Docker

A [docker image](https://docs.docker.com/engine/reference/run/) is automatically built with the latest ydk-gen installed. This be used to run ydk-gen without installing anything natively on your machine.

To use the docker image, [install docker](https://docs.docker.com/install/) on your system and run the below command. See the [docker documentation](https://docs.docker.com/engine/reference/run/) for more details.

```
docker run -it ydkdev/ydk-gen
```

# System requirements

Please follow the below instructions to install the system requirements before installing YDK-Py/YDK-Cpp/YDK-Go. 
**Please note**. If you are using the latest ydk-gen master branch code, you may not be able to use prebuilt libraries and packages. In this case you need to build all the components [from source](#second-step-generate-and-install-the-core) after installing the below requirements:

## Linux
### Ubuntu (Debian-based)

#### Install OS dependency packages

```
   $ sudo apt-get install gdebi-core python3-dev python-dev libtool-bin
   $ sudo apt-get install libcurl4-openssl-dev libpcre3-dev libssh-dev libxml2-dev libxslt1-dev cmake

   $ # Upgrade compiler to gcc 5.*
   $ sudo apt-get install gcc-5 g++-5 -y > /dev/null
   $ sudo ln -sf /usr/bin/g++-5 /usr/bin/c++
   $ sudo ln -sf /usr/bin/gcc-5 /usr/bin/cc
```

#### Install libydk library

You can install the latest `libydk` core library using prebuilt binaries for Xenial and Bionic. For other Ubuntu distributions it is recommended to build `libydk` library from source.

For Xenial (Ubuntu 16.04.4):

```
$ wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/xenial/libydk_0.8.1-1_amd64.deb
$ sudo gdebi libydk_0.8.1-1_amd64.deb
```

For Bionic (Ubuntu 18.04.1):

```
$ wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/bionic/libydk_0.8.1-1_amd64.deb
$ sudo gdebi libydk_0.8.1-1_amd64.deb
```

### Centos (Fedora-based)

#### Install OS dependency packages

```
   $ sudo yum install epel-release
   $ sudo yum install libssh-devel gcc-c++

   # Upgrade compiler to gcc 5.*
   $ sudo yum install centos-release-scl -y > /dev/null
   $ sudo yum install devtoolset-4-gcc* -y > /dev/null
   $ sudo ln -sf /opt/r h/devtoolset-4/root/usr/bin/gcc /usr/bin/cc
   $ sudo ln -sf /opt/rh/devtoolset-4/root/usr/bin/g++ /usr/bin/c++
```

#### Install prebuilt libydk binary

```  
   $ sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.1/libydk-0.8.1-1.x86_64.rpm
```

### Build from source

Install dependencies OS dependencies then generate and install YDK C++ libraries:

```
   # Generate and install libydk library
   $ ./generate --cpp --core
   $ cd gen-api/cpp/ydk/build/
   $ make
   $ sudo make install
   
   # Generate and install libydk_gnmi library (optional)
   $ ./generate --cpp --service profiles/services/gnmi-0.4.0.json
   $ cd gen-api/cpp/ydk-service-gnmi/build/
   $ make
   $ sudo make install
```

## MacOS

It is recommended to install [homebrew](http://brew.sh) and Xcode command line tools on your system before installing YDK-Py/YDK-Cpp/YDK-Go.

You can download the latest Python package from [here](https://www.python.org/downloads/). Please do not use the homebrew version of Python as it causes issues with installation of YDK packages. Please execute `brew rm python python3` to remove any homebrew Python packages.

#### Install prebuilt libydk and optionally libydk_gnmi libraries

```
   $ xcode-select --install
   $ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
   $ brew install pkg-config libssh xml2 curl pcre cmake libxml2 pybind11
   $
   $ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1/libydk-0.8.1-Darwin.pkg
   $ sudo installer -pkg libydk-0.8.1-Darwin.pkg -target /
```

### Build from source

Install dependencies OS dependencies then generate and install YDK C++ libraries:

```
   # Generate and install libydk library
   $ ./generate --cpp --core
   $ cd gen-api/cpp/ydk/build/
   $ make
   $ sudo make install
   
   # Generate and install libydk_gnmi library (optional)
   $ ./generate --cpp --service profiles/services/gnmi-0.4.0.json
   $ cd gen-api/cpp/ydk-service-gnmi/build/
   $ make
   $ sudo make install
```

## Libssh installation

Please note that libssh-0.8.0 `does not support <http://api.libssh.org/master/libssh_tutor_threads.html>`_ separate threading library, 
which is required for YDK. Therefore, if after installation of libssh package you find that the `libssh_threads.a` library is missing, 
please downgrade the installation of libssh to version 0.7.6, or upgrade to 0.8.1 or higher.

```
$ wget https://git.libssh.org/projects/libssh.git/snapshot/libssh-0.7.6.tar.gz
$ tar zxf libssh-0.7.6.tar.gz && rm -f libssh-0.7.6.tar.gz
$ mkdir libssh-0.7.6/build && cd libssh-0.7.6/build
$ cmake ..
$ sudo make install
```

## Windows

Currently, ``YDK-Py`` and ``YDK-Cpp`` from release ``0.6.0`` onwards is not supported on Windows.

## gNMI Requirements

In order to enable YDK support for gNMI protocol, which is optional, the following third party software must be installed prior to gNMI YDK component installation.

### Install protobuf and protoc

```
    wget https://github.com/google/protobuf/releases/download/v3.5.0/protobuf-cpp-3.5.0.zip
    unzip protobuf-cpp-3.5.0.zip
    cd protobuf-3.5.0
    ./configure
    make
    sudo make install
    sudo ldconfig
```

### Install gRPC

```
    git clone -b v1.9.1 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    make
    sudo make install
    sudo ldconfig
    cd -
```

### Instal YDK gNMI library

#### Ubuntu

For Xenial (Ubuntu 16.04.4):

```
$ wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/xenial/libydk_gnmi_0.4.0-1_amd64.deb
$ sudo gdebi libydk_gnmi_0.4.0-1_amd64.deb
```

For Bionic (Ubuntu 18.04.1)

```
$ wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/bionic/libydk_gnmi_0.4.0-1_amd64.deb
$ sudo gdebi libydk_gnmi_0.4.0-1_amd64.deb
```

For MacOS:

```
$ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1/libydk_gnmi-0.4.0-Darwin.pkg
$ sudo installer -pkg libydk_gnmi-0.4.0-Darwin.pkg -target /
```

#### CentOS

```
   sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.1/libydk_gnmi_0.4.0-1.x86_64.rpm
```

### Run-time environment

There is an open issue with gRPC on Centos/Fedora, which requires an extra step before running any YDK gNMI application. 
See this issue on [GRPC GitHub](https://github.com/grpc/grpc/issues/10942#issuecomment-312565041) for details. 
As a workaround, the YDK based application runtime environment must include setting of `LD_LIBRARY_PATH` variable:

```
    PROTO="/Your-Protobuf-and-Grpc-installation-directory"
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PROTO/grpc/libs/opt:$PROTO/protobuf-3.5.0/src/.libs:/usr/local/lib64
```

# Installation
## Setting up your environment

We recommend that you run ydk-gen under Python virtual environment (``virtualenv``/``virtualenvwrapper``).  To install Python virtual environment in your system, execute:

```
  $ pip install virtualenv virtualenvwrapper
  $ source /usr/local/bin/virtualenvwrapper.sh
```

In some systems (e.g. Debian-based Linux), you may need to install support for Python virtual environments as root

Create new virtual environment:

```
  $ mkvirtualenv -p python2.7 py2
```

## Clone ydk-gen and install the requirements

```
$ git clone https://github.com/CiscoDevNet/ydk-gen.git
$ cd ydk-gen
$ pip install -r requirements.txt
```

## Generate YDK components
All the YDK components/packages can be generated by using Python script `generate.py`.

```
$ ./generate.py --help
usage: generate.py [-h] [-l] [--core] [--service SERVICE] [--bundle BUNDLE]
                   [--adhoc-bundle-name ADHOC_BUNDLE_NAME]
                   [--adhoc-bundle ADHOC_BUNDLE [ADHOC_BUNDLE ...]]
                   [--generate-meta] [--generate-doc] [--generate-tests]
                   [--output-directory OUTPUT_DIRECTORY] [--cached-output-dir]
                   [-p] [-c] [-g] [-v] [-o]

Generate YDK artifacts:

optional arguments:
  -h, --help            show this help message and exit
  -l, --libydk          Generate libydk core package
  --core                Generate and/or install core library
  --service SERVICE     Location of service profile JSON file
  --bundle BUNDLE       Location of bundle profile JSON file
  --adhoc-bundle-name ADHOC_BUNDLE_NAME
                        Name of the adhoc bundle
  --adhoc-bundle ADHOC_BUNDLE [ADHOC_BUNDLE ...]
                        Generate an SDK from a specified list of files
  --generate-meta       Generate meta-data for Python bundle
  --generate-doc        Generate documentation
  --generate-tests      Generate tests
  --output-directory OUTPUT_DIRECTORY
                        The output directory where the sdk will get created.
  --cached-output-dir   The output directory specified with --output-directory
                        includes a cache of previously generated gen-
                        api/<language> files under a directory called 'cache'.
                        To be used to generate docs for --core
  -p, --python          Generate Python SDK
  -c, --cpp             Generate C++ SDK
  -g, --go              Generate Go SDK
  -v, --verbose         Verbose mode
  -o, --one-class-per-module
                        Generate separate modules for each python class
                        corresponding to yang containers or lists.
```
The below steps specify how to use `generate.py` to generate YDK core, model bundle, and service packages. All these packages are available for Python, Go and C++ in corresponding github repositories: [ydk-py](https://github.com/CiscoDevNet/ydk-py),  [ydk-go](https://github.com/CiscoDevNet/ydk-go) and [ydk-cpp](https://github.com/CiscoDevNet/ydk-cpp). 

The script [create_ydk_sdk_for_github.sh](create_ydk_sdk_for_github.sh) can be used to generate the `ydk-py`, `ydk-cpp` and `ydk-go` repositories after having generated all the bundles and core packages using `generate.py`.

## First step: choose model bundle profile

The first step in using ydk-gen is either using one of the already built [bundle profiles](https://github.com/CiscoDevNet/ydk-gen/tree/master/profiles/bundles) or constructing your own bundle profile, consisting of the YANG models you are interested to include into the bundle:

Construct a bundle profile file, such as [```ietf_0_1_1.json```](profiles/bundles/ietf_0_1_1.json) and specify its dependencies.

A sample bundle profile file is described below. The file is in a JSON format. Specify the `name` of your bundle, the `version` of the bundle and the `ydk_version`, which refers to [the version](https://github.com/CiscoDevNet/ydk-gen/releases) of the ydk core package you want to use with this bundle. The `name` of the bundle here is especially important as this will form part of the installation path of the bundle.

```
{
    "name":"cisco-ios-xr",
    "version": "0.1.0",
    "ydk_version": "0.8.1",
    "Author": "Cisco",
    "Copyright": "Cisco",
    "Description": "Cisco IOS-XR Native Models From Git",
```

The "models" section of the file describes where to source models from. There are 3 sources:

- Directories
- Specific files
- Git, within which specific relative directories and files may be referenced

The sample below shows the use of git sources only.

```
    "models": {
        "git": [
```

We have a list of git sources. Each source must specify a URL. This URL should be one that allows the repository to be cloned without requiring user intervention, so please use a public URL such as the example below. There are three further options that can be specified:

- ```commitid``` - Optional specification of a commit in string form. The files identified will be copied from the context of this commit.
- ```dir``` - List of **relative** directory paths within git repository. All .yang files in this directory **and any sub-directories** will be pulled into the generated bundle.
- ```file```- List of **relative** file paths within the git repository.

Only directory examples are shown below.

```
            {
                "url": "https://github.com/YangModels/yang.git",
                "dir": [
                    "vendor/cisco/xr/532"
                ]
            },
            {
                "url": "https://github.com/YangModels/yang.git",
                "commitid": "f6b4e2d59d4eedf31ae8b2fa3119468e4c38259c",
                "dir": [
                    "experimental/openconfig/bgp",
                    "experimental/openconfig/policy"
                ]
            }
        ]
    },
```

## Second step: Generate and install the core

Some model bundles have bin packaged and published in [Pypi](https://pypi.org) repository. These bundles can be installed with `pip` utility. For example, when executing `pip install ydk-models-cisco-ios-xr`, you will install the latest released in PyPi IOS XR device package.

**Note:** 
There usually would have been changes on the master branch since the last [released version](https://github.com/CiscoDevNet/ydk-py/releases). To install the latest code at your own risk, you need to follow the below steps in the exact order.

First generate and install ``libydk``

```
$ ./generate.py --libydk
$ cd gen-api/cpp/ydk/build
$ make
$ [sudo] make install
```

To create the libydk binary package to use for later installation, run the below command

```
$ [sudo] make package
```

For Python:

```
$ ./generate.py --python --core
$ pip install gen-api/python/ydk/dist/ydk*.tar.gz
```

For Go:

```
$ export $GOPATH=/your-go-path-installation-directory
$ ./generate.py --go --core
```

## Third step. Generate and install model bundle

Generate model bundle using a bundle profile and install it.

For Python:

```
$ ./generate.py --python --bundle profiles/bundles/<name-of-profile>.json
$ [sudo] pip install gen-api/python/<name-of-bundle>-bundle/dist/ydk*.tar.gz
```

Now, the `pip list | grep ydk` should show the `ydk` (referring to the core package) and `ydk-<name-of-bundle>` packages installed:

```
$ pip list | grep ydk
...

ydk (0.8.1)
ydk-models-<name-of-bundle> (0.5.1)
...
```

For Go:

```
$ export $GOPATH=/your-go-path-installation-directory
$ ./generate.py --go --bundle profiles/bundles/<name-of-profile>.json
```

For C++:

```
$ ./generate.py --cpp --bundle profiles/bundles/<name-of-profile>.json
$ cd gen-api/cpp/<name-of-bundle>-bundle/build
$ make
$ [sudo] make install
```

## Fourth step: Writing your first app

Now, you can start creating apps based on the models in your bundle. Assuming you generated a python bundle, the models will be available for importing in your app under `ydk.models.<name-of-your-bundle>`. For examples, see [ydk-py-samples](https://github.com/CiscoDevNet/ydk-py-samples#a-hello-world-app) and [C++ samples](sdk/cpp/samples). Also refer to the [documentation for python](http://ydk.cisco.com/py/docs/developer_guide.html), [Go](http://ydk.cisco.com/go/docs/developer_guide.html) and [for C++](http://ydk.cisco.com/cpp/docs/developer_guide.html).

## Documentation

When generating the YDK documentation for several bundles and the core, it is recommended to generate the bundles without the `--generate-doc` option. After generating all the bundles, the combined documentation for all the bundles and the core can be generated using the `--core --generate-doc` option. For example, the below sequence of commands will generate the documentation for the three python bundles and the python core (for C++, use `--cpp`; for Go, use `--go`).

Note that the below process could take a few hours due to the size of the `cisco_ios_xr` bundle.

```
./generate.py --python --bundle profiles/bundles/ietf_0_1_1.json
./generate.py --python --bundle profiles/bundles/openconfig_0_1_1.json
./generate.py --python --bundle profiles/bundles/cisco_ios_xr_6_1_1.json
./generate.py --python --core --generate-doc
```

If you have previously generated documentation, using the `--cached-output-dir --output-directory <dir>` option can be used to reduce document generation time. Taking Python as an example:

```
mkdir gen-api/cache
mv gen-api/python gen-api/cache

./generate.py --python --bundle profiles/bundles/ietf_0_1_5.json
./generate.py --python --bundle profiles/bundles/openconfig_0_1_5.json
./generate.py --python --bundle profiles/bundles/cisco_ios_xr_6_3_2.json
./generate.py --python --core --generate-doc --output-directory gen-api --cached-output-dir -v
```

Pre-generated documentation for [ydk-py](http://ydk.cisco.com/py/docs/) and [ydk-cpp](http://ydk.cisco.com/cpp/docs/) are available.

# Generating an "Adhoc" YDK-Py Bundle

The ability to generate an adhoc bundle directly from the command line and without creating a bundle file can be done something like this:

```
$ ./generate.py --adhoc-bundle-name test --adhoc-bundle \
    /opt/git-repos/clean-yang/vendor/cisco/xr/621/Cisco-IOS-XR-ipv4-bgp-oper*.yang \
    /opt/git-repos/clean-yang/vendor/cisco/xr/621/Cisco-IOS-XR-types.yang
    /opt/git-repos/clean-yang/vendor/cisco/xr/621/Cisco-IOS-XR-ipv4-bgp-datatypes.yang
```

When run in this way, we will generate a bundle that only contains the files specified with the `--adhoc-bundle` option, creating a `pip` package name by the `--adhoc-bundle-name`, with a version `0.1.0` and a dependency on the base IETF bundle. Note that **all** dependencies for the bundle must be listed, and the expectation is that this option will typically be used for generating point YDK-Py bundles for specific testing, the `--verbose` option is automatically enabled to quickly and easily let a user see if dependencies have been satisfied.

# Notes

## Python Requirements

YDK supports both Python2 and Python3 versions.  At least Python2.7 or Python3.4 along with corresponding utilities pip and pip3 must be installed on your system. 

It is also required for Python installation to include corresponding shared library. As example: 

 - python2.7  - /usr/lib/x86_64-linux-gnu/libpython2.7.so
 - python3.5m - /usr/lib/x86_64-linux-gnu/libpython3.5m.so

Please follow [System requirements](#system-requirements) to assure presence of shared Python libraries.

## Directory structure

```
README          - install and usage notes
gen-api         - generated bundle/core
                    - python (Python SDK)
                    - go (Go SDK)
                    - cpp (C++ SDK)

generate.py     - script used to generate SDK for yang models
profiles        - profile files used during generation
yang            - some yang models used for testing
requirements.txt- python dependencies used during installation
sdk             - sdk core and stubs for python and cpp
test            - test code
```

## Troubleshooting

Sometimes, developers using ydk-gen may run across errors when generating a YDK bundle using generate.py with some yang models. If there are issues with the .json profile file being used, such errors will be easily evident. Other times, when the problem is not so evident, it is recommended to try running with the `[--verbose|-v]` flag, which may reveal syntax problems with the yang models being used. For example,

```
./generate.py --python --bundle profiles/bundles/ietf_0_1_1.json --verbose
```

Also, it may be a good idea to obtain a local copy of the yang models and compile them using `pyang` to ensure the validity of the models,
```
cd /path/to/yang/models
pyang *.yang
```

# Running Unit Tests

## Python

First, generate and install the [test](profiles/test/ydktest.json) bundle and core package

```
$ ./generate.py --core
$ pip install gen-api/python/core/dist/ydk*.tar.gz

$ ./generate.py --bundle profiles/test/ydktest.json
$ pip install gen-api/python/ydktest-bundle/dist/ydk*.tar.gz
```

To run the sanity tests, do the following:

```
$ cd ydk-gen/sdk/python
$ python test/test_sanity_types.py
$ python test/test_sanity_levels.py
$ python test/test_sanity_filters.py
```

## C++

First, install the core and [test](profiles/test/ydktest-cpp.json) bundle package.

```
$ ./generate.py --core --cpp
$ ./generate.py --bundle profiles/test/ydktest-cpp.json --cpp
$ cd gen-api/cpp/ydktest-bundle/build
$ make && make install
$ cd -
```

To run the core and bundle tests, do the following:

```
$ cd ydk-gen/sdk/cpp/ydk
$ mkdir build && cd build
$ cmake .. && make && sudo make install test

$ cd ydk-gen/sdk/cpp/tests
$ mkdir build && cd build
$ cmake .. && make all test
```

## Go

Please refer [here](https://github.com/CiscoDevNet/ydk-gen/blob/master/sdk/go/core/README.md).

Support
=======
Join the [YDK community](https://communities.cisco.com/community/developer/ydk) to connect with other users and with the makers of YDK.

Release Notes
===============
The current YDK release version is 0.8.1. The version of the latest YDK-Gen master branch is 0.8.1. YDK-Gen is licensed under the Apache 2.0 License.
