===============
Getting Started
===============
.. contents:: Table of Contents

Overview
========

The YANG Development Kit (YDK) is a Software Development Kit that provides API's that are modeled in YANG. 
The main goal of YDK is to reduce the learning curve of YANG data models by expressing the model semantics in an API 
and abstracting protocol/encoding details.  YDK is composed of a core package that defines services and providers, 
plus one or more module bundles that are based on YANG models.

Docker
======

A `docker image <https://docs.docker.com/engine/reference/run/>`_ is automatically built with the latest ydk-py installed. 
This be used to run ydk-py without installing anything natively on your machine.

To use the docker image, `install docker <https://docs.docker.com/install/>`_ on your system and run the below command. 
See the `docker documentation <https://docs.docker.com/engine/reference/run/>`_ for more details::

  docker run -it ydkdev/ydk-py


System Requirements
===================

Please follow the below instructions to install the system requirements before installing YDK-Py:

Linux
-----

Ubuntu (Debian-based)
~~~~~~~~~~~~~~~~~~~~~

Download and install third party dependency packages and then YDK core library `libydk`. 
You can install the library using prebuilt debian packages for Xenial and Bionic LTS distributions. 
For other Ubuntu distributions it is recommended to build core library from source.

Install Third-party software dependencies::

   $ sudo apt-get install gdebi-core python3-dev python-dev libtool-bin
   $ sudo apt-get install libcurl4-openssl-dev libpcre3-dev libssh-dev libxml2-dev libxslt1-dev cmake

For Xenial (Ubuntu 16.04.4)::

   # Upgrade compiler to gcc 5.*
   $ sudo apt-get install gcc-5 g++-5 -y > /dev/null
   $ sudo ln -sf /usr/bin/gcc-5 /usr/bin/cc
   $ sudo ln -sf /usr/bin/g++-5 /usr/bin/c++

   # Install YDK core library
   $ wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/xenial/libydk_0.8.1-1_amd64.deb
   $ sudo gdebi libydk_0.8.1-1_amd64.deb

For Bionic (Ubuntu 18.04.1)::

   # Install YDK core library
   $ wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/bionic/libydk_0.8.1-1_amd64.deb
   $ sudo gdebi libydk_0.8.1-1_amd64.deb

Centos (Fedora-based)
~~~~~~~~~~~~~~~~~~~~~

**Note:** Currently, only Centos7 and RHEL7 are known to work with YDK.
Please see `this issue on YDK GitHub <https://github.com/CiscoDevNet/ydk-gen/issues/518>`_ for any potential/usage installation on CentOS.

.. code-block:: sh

   # Install Third-party software dependencies
   $ sudo yum install epel-release
   $ sudo yum install libssh-devel gcc-c++ python-devel python3-devel

   # Upgrade compiler to gcc 5.*
   $ yum install centos-release-scl -y > /dev/null
   $ yum install devtoolset-4-gcc* -y > /dev/null
   $ ln -sf /opt/rh/devtoolset-4/root/usr/bin/gcc /usr/bin/cc
   $ ln -sf /opt/rh/devtoolset-4/root/usr/bin/g++ /usr/bin/c++

   # Install YDK core library
   $ sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.1/libydk-0.8.1-1.x86_64.rpm


MacOS
-----

**Note:** Please do not use the homebrew version of python as it causes issues with installation of YDK packages. Please execute ``brew rm python python3`` to remove any homebrew python packages.

It is required to install Xcode command line tools, `homebrew <http://brew.sh>`_ and the following homebrew packages on your system before installing YDK-Py.

.. code-block:: sh

   # Install Third-party software dependencies
   $ xcode-select --install
   $ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
   $ brew install pkg-config libssh xml2 libxml2 curl pcre cmake pybind11

   # Install YDK core library
   $ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1/libydk-0.8.1-Darwin.pkg
   $ sudo installer -pkg libydk-0.8.1-Darwin.pkg -target /

Libssh Installation
-------------------

The libssh-0.8.0 `does not support <http://api.libssh.org/master/libssh_tutor_threads.html>`_ separate threading library, 
which is required for YDK. If after installation of libssh package the `libssh_threads.a` is missing, please downgrade the installation to libssh-0.7.6, 
or upgrade to libssh-0.8.1 or higher.

gNMI Requirements
===================

In order to have YDK support for gNMI protocol, which is optional, the following third party software must be installed prior to gNMI YDK component installation.

Install protobuf
----------------

.. code-block:: sh

    wget https://github.com/google/protobuf/releases/download/v3.5.0/protobuf-cpp-3.5.0.zip
    unzip protobuf-cpp-3.5.0.zip
    cd protobuf-3.5.0
    ./configure
    make
    sudo make install
    sudo ldconfig

Install gRPC
------------

.. code-block:: sh

    git clone -b v1.9.1 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    sudo ldconfig
    make
    sudo make install

Instal YDK gNMI library
-----------------------

Ubuntu
~~~~~~

For Xenial (Ubuntu 16.04.4):

.. code-block:: sh

   wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/xenial/libydk_gnmi_0.4.0-1_amd64.deb
   sudo gdebi libydk_gnmi_0.4.0-1_amd64.deb

For Bionic (Ubuntu 18.04.1):

.. code-block:: sh

   wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.1/bionic/libydk_gnmi_0.4.0-1_amd64.deb
   sudo gdebi libydk_gnmi_0.4.0-1_amd64.deb

CentOS
~~~~~~

.. code-block:: sh

   sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.1/libydk_gnmi_0.4.0-1.x86_64.rpm
   
MacOS
~~~~~

.. code-block:: sh

   curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.1/libydk_gnmi-0.4.0-Darwin.pkg
   sudo installer -pkg libydk_gnmi-0.4.0-Darwin.pkg -target /


Runtime environment
~~~~~~~~~~~~~~~~~~~~

The YDK based application runtime environment must include setting of **LD_LIBRARY_PATH** variable:

.. code-block:: sh

   PROTO="/Your-Protobuf-and-Grpc-installation-directory"
   export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PROTO/grpc/libs/opt:$PROTO/protobuf-3.5.0/src/.libs:/usr/local/lib64

Python Requirements
===================

YDK supports both Python2 and Python3 versions.  At least Python2.7 or Python3.4 must be installed on your system. 

It is also required for Python installation to include corresponding shared library. As example: 

 - python2.7  - /usr/lib/x86_64-linux-gnu/libpython2.7.so
 - python3.5m - /usr/lib/x86_64-linux-gnu/libpython3.5m.so

Please follow `System Requirements` to assure presence of shared Python libraries.

.. _howto-install:

Quick install
=============

Centos (Fedora-based)
---------------------

You can install the latest model packages from the DevHub artifactory and Python package index.  
Note that, in some systems, you need to install the new package as root.

To install the core and model bundles on Centos, please follow the below steps.

``Python2.7``::

.. code-block:: sh

    pip install ydk
    pip install --install-option="--install-purelib=/usr/lib64/python2.7/site-packages" --no-deps ydk-models-ietf
    pip install --install-option="--install-purelib=/usr/lib64/python2.7/site-packages" --no-deps ydk-models-openconfig
    pip install --install-option="--install-purelib=/usr/lib64/python2.7/site-packages" --no-deps ydk-models-cisco-ios-xr
    pip install --install-option="--install-purelib=/usr/lib64/python2.7/site-packages" --no-deps ydk-models-cisco-ios-xe

``Python3.4``::

.. code-block:: sh

    pip install ydk
    pip install --install-option="--install-purelib=/usr/lib64/python3.4/site-packages" --no-deps ydk-models-ietf
    pip install --install-option="--install-purelib=/usr/lib64/python3.4/site-packages" --no-deps ydk-models-openconfig
    pip install --install-option="--install-purelib=/usr/lib64/python3.4/site-packages" --no-deps ydk-models-cisco-ios-xr
    pip install --install-option="--install-purelib=/usr/lib64/python3.4/site-packages" --no-deps ydk-models-cisco-ios-xe

``Python3.6``::

 .. code-block:: sh

    pip install ydk
    pip install --install-option="--install-purelib=/usr/lib64/python3.6/site-packages" --no-deps ydk-models-ietf
    pip install --install-option="--install-purelib=/usr/lib64/python3.6/site-packages" --no-deps ydk-models-openconfig
    pip install --install-option="--install-purelib=/usr/lib64/python3.6/site-packages" --no-deps ydk-models-cisco-ios-xr
    pip install --install-option="--install-purelib=/usr/lib64/python3.6/site-packages" --no-deps ydk-models-cisco-ios-xe


Other platforms
---------------

You can install the latest model packages from the DevHub artifactory and Python package index.  
Note that, in some systems, you need to install the new package as root.  
You get fully operational YDK environment by installing the ``cisco-ios-xr`` and/or ``cisco-ios-xe`` bundle(s) 
(depending on whether you're developing for an IOS XR or IOS XE platform) which automatically installs all other 
YDK-related packages (``YDK``, ``openconfig`` and ``ietf`` packages):

.. code-block:: sh

    $ pip install https://devhub.cisco.com/artifactory/pypi-ydk/0.8.1/ydk-0.8.1.tar.gz
    $ pip install https://devhub.cisco.com/artifactory/pypi-ydk/0.8.1/ydk-service-gnmi-0.4.0.tar.gz
    $ pip install ydk-models-cisco-ios-xr
    $ pip install ydk-models-cisco-ios-xe

Alternatively, you can perform a partial installation.  
If you only prefer to install the ``openconfig`` bundle and its dependencies (``YDK`` and ``ietf`` packages), execute:

.. code-block:: sh

    $ pip install https://devhub.cisco.com/artifactory/pypi-ydk/0.8.1/ydk-0.8.1.tar.gz
    $ pip install ydk-models-openconfig

If you only want to install the ``ietf`` bundle and its dependencies (``YDK`` package), execute:

.. code-block:: sh

    $ pip install https://devhub.cisco.com/artifactory/pypi-ydk/0.8.1/ydk-0.8.1.tar.gz
    $ pip install ydk-models-ietf

Installing from Source
======================

Installing core package
-----------------------

If you prefer not to use the YDK packages in the DevHub artifactory or Python package index, you need to install manually the ``YDK`` core package and then the model bundles that you plan to use.  
The Python core package is dependent on C++ core library `libydk`, which must be installed prior to Python package installation:

.. code-block:: sh

    $ git clone https://github.com/CiscoDevNet/ydk-gen.git
    $ cd ydk-gen/sdk/cpp/core
    core$ mkdir -p build
    core$ cd build
    build$ cmake ..
    build$ make
    build$ sudo make install

To install the ``YDK`` Python core package, execute:

.. code-block:: sh

    $ cd ydk-gen/sdk/python/core
    core$ python setup.py sdist
    core$ pip install dist/ydk*.gz

Installing model bundles
------------------------

Once you have installed the ``YDK`` core package, you can install one or more model bundles. The source code for the model bundles can be generated by running:

.. code-block:: sh

    $ git clone https://github.com/CiscoDevNet/ydk-py.git -b 0.8.1

Note that some bundles have dependencies on other bundles.  Those dependencies are already captured in the bundle package.  Make sure you install the desired bundles in the order below.  To install the ``ietf`` bundle, execute:

.. code-block:: sh

    core$ cd ../ietf
    ietf$ python setup.py sdist
    ietf$ pip install dist/ydk*.gz

To install the ``openconfig`` bundle, execute:

.. code-block:: sh

    ietf$ cd ../openconfig
    openconfig$ python setup.py sdist
    openconfig$ pip install dist/ydk*.gz

To install the ``cisco-ios-xr`` bundle, execute:

.. code-block:: sh

    openconfig$ cd ../cisco-ios-xr
    cisco-ios-xr$ python setup.py sdist
    cisco-ios-xr$ pip install dist/ydk*.gz
    cisco-ios-xr$ cd ..

Installing gNMI package
-----------------------

Optionaly the gNMI package for Python can be installed. The Python gNMI package is dependent on C++ core library `libydk` 
(see installation of Python core package above) and C++ gNMI library `libydk_gnmi`, which must be installed prior to Python package installation:

.. code-block:: sh

    $ cd ydk-gen/sdk/cpp/gnmi
    gnmi$ mkdir -p build
    gnmi$ cd build
    build$ cmake ..
    build$ make
    build$ sudo make install

To install the ``YDK`` Python gNMI package, execute:

.. code-block:: sh

    $ cd ydk-gen/sdk/python/gnmi
    gnmi$ python setup.py sdist
    gnmi$ pip install dist/ydk*.gz

Using Virtual Environment
=========================

You may want to perform the installation under Python virtual environment (`virtualenv <https://pypi.python.org/pypi/virtualenv/>`_/`virtualenvwrapper  <https://pypi.python.org/pypi/virtualenvwrapper>`_).  A virtual environment allows you to install multiple versions of YDK if needed.  In addition, it prevents any potential conflicts between package dependencies in your system.

To install virtual environment on your system, execute:

.. code-block:: sh

    $ pip install virtualenv virtualenvwrapper
    $ source /usr/local/bin/virtualenvwrapper.sh

**Note**  In some systems (e.g. Debian-based Linux), you need to be a root user or use `sudo` access.

Create new virtual environment:

.. code-block:: sh

    $ mkvirtualenv -p python2.7 ydk-py

At this point, you can perform the quick install or the installation from source described above.  Take into account that you must not attempt to install YDK as root under virtual environment.

Samples
=======

To get started using the YDK API, there are sample applications available in the `YDK-Py repository <https://github.com/CiscoDevNet/ydk-py/tree/master/core/samples>`_. For example, to run the ``bgp.py`` sample, execute:

.. code-block:: sh

    (ydk-py)ydk-py$ cd core/samples
    (ydk-py)samples$ ./bgp.py -h
    Usage: bgp.py [-h | --help] [options]

    Options:
    -h, --help            show this help message and exit
    -v VERSION, --version=VERSION
                          force NETCONF version 1.0 or 1.1
    -u USERNAME, --user=USERNAME
                          login user name
    -p PASSWORD, --password=PASSWORD
                          login user password
    --proto=PROTOCOL      Which transport protocol to use, one of ssh or tcp
    --host=HOST           NETCONF agent hostname or IP address
    --port=PORT           NETCONF agent SSH port

    (ydk-py)samples$ ./bgp.py --host <ip-address-of-netconf-server> -u <username> -p <password> --port <port-number>

Documentation and Support
=========================

- Hundreds of samples can be found in the `YDK-Py samples repository <https://github.com/CiscoDevNet/ydk-py-samples>`_
- Join the `YDK community <https://communities.cisco.com/community/developer/ydk>`_ to connect with other users and with the makers of YDK
