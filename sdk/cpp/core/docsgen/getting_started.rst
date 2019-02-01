===============
Getting Started
===============
.. contents:: Table of Contents

Overview
========

The YANG Development Kit (YDK) is a Software Development Kit that provides API's that are modeled in YANG. The main goal of YDK is to reduce the learning curve of YANG data models by expressing the model semantics in an API and abstracting protocol/encoding details.  YDK is composed of a core package that defines services and providers, plus one or more module bundles that are based on YANG models.

How to Install
==============

You can install YDK-Cpp on Linux or MacOS.  It is currently not supported on Windows.

System Requirements
-------------------

Linux
~~~~~

**Ubuntu (Debian-based)**

The following packages must be present in your system before installing YDK-Cpp:

.. code-block:: sh

   # Linux dependencies
   sudo apt-get install libcurl4-openssl-dev libpcre3-dev libssh-dev libxml2-dev libxslt1-dev libtool-bin cmake
   #
   # gcc-5 and g++5 for modern c++
   sudo apt-get install gcc-5 g++-5 -y > /dev/null
   sudo ln -f -s /usr/bin/g++-5 /usr/bin/c++
   sudo ln -f -s /usr/bin/gcc-5 /usr/bin/cc

**Centos (Fedora-based)**

The following packages must be present in your system before installing YDK-Cpp:

.. code-block:: sh

   # Linux dependencies
   sudo yum install epel-release
   sudo yum install libxml2-devel libxslt-devel libssh-devel libtool gcc-c++ pcre-devel cmake3 wget
   #
   # gcc-5 and g++5 for modern c++
   sudo yum install centos-release-scl -y > /dev/null
   sudo yum install devtoolset-4-gcc* -y > /dev/null
   sudo ln -sf /opt/rh/devtoolset-4/root/usr/bin/gcc /usr/bin/cc
   sudo ln -sf /opt/rh/devtoolset-4/root/usr/bin/g++ /usr/bin/c++
   sudo ln -sf /opt/rh/devtoolset-4/root/usr/bin/gcc /usr/bin/gcc
   sudo ln -sf /opt/rh/devtoolset-4/root/usr/bin/g++ /usr/bin/g++

MacOS
~~~~~

It is recommended to install `homebrew <http://brew.sh>`_ and Xcode command line tools on your system before installing YDK-Cpp:

.. code-block:: sh

  /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
  brew install curl libssh pcre xml2 cmake
  xcode-select --install

Libssh Installation
-------------------

The libssh-0.8.0 `does not support <http://api.libssh.org/master/libssh_tutor_threads.html>`_ separate threading library, 
which is required for YDK. If after installation of libssh package the `libssh_threads.a` is missing, please downgrade the installation to libssh-0.7.6, 
or upgrade to libssh-0.8.1 or higher.

gNMI Requirements
~~~~~~~~~~~~~~~~~

In order to enable YDK support for gNMI protocol, which is optional, the following third party software must be installed prior to gNMI YDK component installation.

**Install protobuf and protoc**

.. code-block:: sh

    wget https://github.com/google/protobuf/releases/download/v3.5.0/protobuf-cpp-3.5.0.zip
    unzip protobuf-cpp-3.5.0.zip
    cd protobuf-3.5.0
    ./configure
    make
    sudo make install
    sudo ldconfig
    cd -

**Install gRPC**

.. code-block:: sh

    git clone -b v1.9.1 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    make
    sudo make install
    sudo ldconfig
    cd -

Quick YDK Installation
----------------------

Install prebuilt libraries
~~~~~~~~~~~~~~~~~~~~~~~~~~

**Ubuntu**

Download and install YDK core library - `libydk`. You can install the library using prebuilt debian packages for Xenial and Bionic LTS distributions. 
For other Ubuntu distributions it is recommended to build core library from source.

For Xenial (Ubuntu 16.04.4):

.. code-block:: sh

   wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/xenial/libydk_0.8.1-1_amd64.deb
   sudo gdebi libydk_0.8.1-1_amd64.deb

For Bionic (Ubuntu 18.04.1):

.. code-block:: sh

   wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/bionic/libydk_0.8.1-1_amd64.deb
   sudo gdebi libydk_0.8.1-1_amd64.deb

**CentOS**

Install YDK core library:

.. code-block:: sh

   sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.1/libydk-0.8.1-1.x86_64.rpm

Install YDK gNMI library (optional):

.. code-block:: sh

   sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.1/libydk_gnmi-0.4.0-1.x86_64.rpm

**MacOS**

You can install the latest model packages using `homebrew <http://brew.sh>`_.  This utility will manage the dependencies between YDK packages and all other sytem dependencies.  First, add the third-party repository (homebrew tap) for YDK:

.. code-block:: sh

   $ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1/libydk-0.8.1-Darwin.pkg
   $ sudo installer -pkg libydk-0.8.1-Darwin.pkg -target /
   $
   $ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1/libydk_gnmi-0.4.0-Darwin.pkg
   $ sudo installer -pkg libydk_gnmi-0.4.0-Darwin.pkg -target /

You get a fully operational YDK environment by installing the ``cisco-ios-xr`` bundle, which automatically installs all other YDK-related packages (``ydk``, ``cisco-ios-xr``, ``openconfig`` and ``ietf`` packages):

.. code-block:: sh

  brew install ydk-cisco-ios-xr

Alternatively, you can perform partial installation.  If you only want to install the ``openconfig`` bundle and its dependencies (``ydk`` and ``ietf`` packages), execute:

.. code-block:: sh

   $ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1-beta/libydk-0.8.1-Darwin.pkg
   $ sudo installer -pkg libydk-0.8.1-Darwin.pkg -target /
   $
   $ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1-beta/libydk_gnmi-0.4.0-Darwin.pkg
   $ sudo installer -pkg libydk_gnmi-0.4.0-Darwin.pkg -target /

If you want to install only the ``ietf`` bundle and its dependencies (``ydk`` package), execute:

.. code-block:: sh

  brew install ydk-ietf

gNMI Service Installation
-------------------------

Instal YDK gNMI library
~~~~~~~~~~~~~~~~~~~~~~~

**Ubuntu**

For Xenial (Ubuntu 16.04.4):

.. code-block:: sh

   wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/xenial/libydk_gnmi_0.4.0-1_amd64.deb
   sudo gdebi libydk_gnmi_0.4.0-1_amd64.deb

For Bionic (Ubuntu 18.04.1):

.. code-block:: sh

   wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/bionic/libydk_gnmi_0.4.0-1_amd64.deb
   sudo gdebi libydk_gnmi_0.4.0-1_amd64.deb

**CentOS**

.. code-block:: sh

   sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.1/libydk_gnmi_0.4.0-1.x86_64.rpm

**MacOS**

.. code-block:: sh

   curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1/libydk_gnmi-0.4.0-Darwin.pkg
   sudo installer -pkg libydk_gnmi-0.4.0-Darwin.pkg -target /

Set runtime environment
~~~~~~~~~~~~~~~~~~~~~~~

The YDK based application runtime environment must include setting of **LD_LIBRARY_PATH** variable:

.. code-block:: sh

   PROTO="/Your-Protobuf-and-Grpc-installation-directory"
   export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PROTO/grpc/libs/opt:$PROTO/protobuf-3.5.0/src/.libs:/usr/local/lib64

Installing from source
----------------------

In order to build YDK components from source, download or clone source files from `YDK-Cpp repository <https://github.com/CiscoDevNet/ydk-cpp>`_

Installing YDK core library
~~~~~~~~~~~~~~~~~~~~~~~~~~~

YDK uses ``cmake`` as the build system of choice. To install the ``core`` package, execute:

.. code-block:: sh

  $ cd ydk-cpp/core/ydk
  ydk$ mkdir build && cd build
  build$ cmake .. && make
  build$ sudo make install

Installing model bundles
~~~~~~~~~~~~~~~~~~~~~~~~

Once you have installed the ``core`` package, you can install one or more model bundles.  Note that some bundles have dependencies on other bundles.  Those dependencies are captured in the bundle packages used for quick installation. 

To install the ``ietf`` bundle, execute:

.. code-block:: sh

  $ cd ydk-cpp/ietf
  ietf$ mkdir build && cd build
  build$ cmake .. && make
  build$ sudo make install

To install the ``openconfig`` bundle, execute:

.. code-block:: sh

  $ cd ydk-cpp/openconfig
  openconfig$ mkdir build && cd build
  build$ cmake .. && make
  build$ sudo make install

To install the ``cisco-ios-xr`` bundle, execute:

.. code-block:: sh

  $ cd ydk-cpp/cisco-ios-xr
  cisco-ios-xr$ mkdir build && cd build
  build$ cmake .. && make
  build$ sudo make install

Instaling YDK gNMI library
~~~~~~~~~~~~~~~~~~~~~~~~~~

Optionaly the YDK gNMI Service library can be installed. Prior to this installation the YDK core library must be installed (see above).

.. code-block:: sh

  $ cd ydk-cpp/gnmi
  gnmi$ mkdir -p build
  gnmi$ cd build
  build$ cmake ..
  build$ make
  build$ sudo make install

Samples
=======

To get started using the YDK API, there are sample apps available in the `YDK-Cpp samples repository <https://github.com/CiscoDevNet/ydk-cpp/tree/master/core/ydk/samples>`_. For example, to run the ``bgp_create.cpp`` sample execute:

.. code-block:: sh

  $ ydk-cpp$ cd core/samples
  samples$ mkdir build && cd build
  build$ cmake .. && make
  build$ ./bgp_create ssh://<username>:<password>@<host-address>:<port> [-v]

Documentation and Support
=========================
- Numerous additional samples can be found in the `YDK-Cpp samples repository <https://github.com/CiscoDevNet/ydk-cpp/tree/master/core/ydk/samples>`_
- Join the `YDK community <https://communities.cisco.com/community/developer/ydk>`_ to connect with other users and with the makers of YDK
