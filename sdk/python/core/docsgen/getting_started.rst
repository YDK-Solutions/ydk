===============
Getting Started
===============
.. contents:: Table of Contents

Overview
========

The YANG Development Kit (YDK) is a Software Development Kit that provides API's that are modeled in YANG. The main goal of YDK is to reduce the learning curve of YANG data models by expressing the model semantics in an API and abstracting protocol/encoding details.  YDK is composed of a core package that defines services and providers, plus one or more module bundles that are based on YANG models.

System Requirements
===================
Please follow the below instructions to install the system requirements before installing YDK-Py:

gNMI Requirements
------------------

Install protobuf

.. code-block:: sh

    wget https://github.com/google/protobuf/releases/download/v3.3.0/protobuf-cpp-3.3.0.zip
    unzip protobuf-cpp-3.3.0.zip
    cd protobuf-3.3.0
    ./configure
    make
    make check
    sudo make install
    sudo ldconfig
    cd -

Install gRPC

.. code-block:: sh

    git clone -b v1.4.5 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    sudo ldconfig
    make
    sudo make install
    cd -

**Note:** There is an open issue with gRPC on Centos/Fedora which requires an extra step before running any YDK gNMI app. See `this issue on GRPC GitHub <https://github.com/grpc/grpc/issues/10942#issuecomment-312565041>`_ for details.

.. code-block:: sh

   $ export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/location_where_grpc_is_downloaded/grpc/libs/opt:/location_where_protobuf_is_downloaded/protobuf-3.3.0/src/.libs:/usr/local/lib64


Linux
-----
Ubuntu (Debian-based)

.. code-block:: sh

   $ sudo apt-get install gdebi-core
   $ wget https://devhub.cisco.com/artifactory/debian-ydk/0.8.0-alpha/libydk_0.8.0-1_amd64.deb
   $ sudo gdebi libydk_0.8.0-1_amd64.deb

Centos (Fedora-based)

.. code-block:: sh

   $ sudo yum install epel-release libssh-devel gcc-c++ python-devel
   $ sudo yum install https://devhub.cisco.com/artifactory/rpm-ydk/0.8.0-alpha/libydk-0.8.0-1.x86_64.rpm

Please see `this issue on YDK GitHub <https://github.com/CiscoDevNet/ydk-gen/issues/518>`_ for any potential/usage installation on CentOS.

macOS
-----
You can download the latest python package from here. **Note:** Please do not use the homebrew version of python as it causes issues with installing ydk packages. Please execute ``brew rm python python3`` to remove any homebrew python packages.

It is required to install Xcode command line tools, `homebrew <http://brew.sh>`_ and the following homebrew packages on your system before installing YDK-Py.

.. code-block:: sh

   $ xcode-select --install
   $ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
   $ brew install pkg-config libssh libxml2 curl pcre cmake
   $ curl -O https://devhub.cisco.com/artifactory/osx-ydk/0.8.0-alpha/libydk-0.8.0-Darwin.pkg
   $ sudo installer -pkg libydk-0.8.0-Darwin.pkg -target /

Windows
-------
Currently, ``YDK-Py`` from release ``0.6.0`` onwards is not supported on Windows.

Python Requirements
===================
Both Python 2 and 3 are supported.  At least Python2.7 or Python 3.4 must be installed in your system.

Backwards Compatibility Notes
=============================
When installing and using the ``0.6.0`` and newer releases of ``YDK-Py``, please refer to the :ref:`compatibility`.

.. _howto-install:

How to install
==============
Quick Install
-------------
You can install the latest model packages from the DevHub artifactory and Python package index.  Note that, in some systems, you need to install the new package as root.  You get a fully operational YDK environment by installing the ``cisco-ios-xr`` and/or ``cisco-ios-xe`` bundle(s) (depending on whether you're developing for an IOS XR or IOS XE platform) which automatically installs all other YDK-related packages (``ydk``, ``openconfig`` and ``ietf`` packages):

.. code-block:: sh

    $ pip install https://devhub.cisco.com/artifactory/osx-ydk/0.8.0-alpha/ydk-0.8.0a0.tar.gz
    $ pip install ydk-models-cisco-ios-xr
    $ pip install ydk-models-cisco-ios-xe

Alternatively, you can perform a partial installation.  If you only want to install the ``openconfig`` bundle and its dependencies (``ydk`` and ``ietf`` packages), execute:

.. code-block:: sh

    $ pip install https://devhub.cisco.com/artifactory/osx-ydk/0.8.0-alpha/ydk-0.8.0a0.tar.gz
    $ pip install ydk-models-openconfig

If you only want to install the ``ietf`` bundle and its dependencies (``ydk`` package), execute:

.. code-block:: sh

    $ pip install https://devhub.cisco.com/artifactory/osx-ydk/0.8.0-alpha/ydk-0.8.0a0.tar.gz
    $ pip install ydk-models-ietf

Installing from Source
----------------------
If you prefer not to use the YDK packages in the DevHub artifactory or Python package index, you need to install manually the ``ydk`` core package and then the model bundles you plan to use.  To install the ``ydk`` core package, execute:

.. code-block:: sh

    $ git clone https://github.com/CiscoDevNet/ydk-gen.git -b gnmi
    $ cd ydk-gen/sdk/python/core
    core$ python setup.py sdist
    core$ pip install dist/ydk*.gz

Once you have installed the ``ydk`` core package, you can install one more model bundles. The source code for the model bundles can be got by running:

.. code-block:: sh

    $ git clone https://github.com/CiscoDevNet/ydk-py.git -b 0.6.2

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

Using a Virtual Environment
---------------------------
You may want to perform the installation under a Python virtual environment (`virtualenv <https://pypi.python.org/pypi/virtualenv/>`_/`virtualenvwrapper  <https://pypi.python.org/pypi/virtualenvwrapper>`_).  A virtual environment allows you to install multiple versions of YDK if needed.  In addition, it prevents any potential conflicts between package dependencies in your system.

To install virtual environment support in your system, execute:

.. code-block:: sh

    $ pip install virtualenv virtualenvwrapper
    $ source /usr/local/bin/virtualenvwrapper.sh

In some systems (e.g. Debian-based Linux), you need to install support for Python virtual environments as root:

.. code-block:: sh

    $ sudo pip install virtualenv virtualenvwrapper
    $ source /usr/local/bin/virtualenvwrapper.sh

Create a new virtual environment:

.. code-block:: sh

    $ mkvirtualenv -p python2.7 ydk-py

At this point, you can perform the quick install or the installation from source described above.  Take into account that must not attempt to install YDK as root under a virtual environment.

Samples
=======
To get started with using the YDK API, there are sample apps available in the `YDK-Py repository <https://github.com/CiscoDevNet/ydk-py/tree/master/core/samples>`_. For example, to run the ``bgp.py`` sample, execute:

.. code-block:: sh

    (ydk-py)ydk-py$ cd core/samples
    (ydk-py)samples$ ./bgp.py -h
    Usage: bgp.py [-h | --help] [options]

    Options:
    -h, --help            show this help message and exit
    -v VERSION, --version=VERSION
                        force NETCONF version 1.0 or 1.1
    -u USERNAME, --user=USERNAME
    -p PASSWORD, --password=PASSWORD
                        password
    --proto=PROTO         Which transport protocol to use, one of ssh or tcp
    --host=HOST           NETCONF agent hostname
    --port=PORT           NETCONF agent SSH port

    (ydk-py)samples$ ./bgp.py --host <ip-address-of-netconf-server> -u <username> -p <password> --port <port-number>

Documentation and Support
=========================
- Hundreds of samples can be found in the `YDK-Py samples repository <https://github.com/CiscoDevNet/ydk-py-samples>`_
- Join the `YDK community <https://communities.cisco.com/community/developer/ydk>`_ to connect with other users and with the makers of YDK
