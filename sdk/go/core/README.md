### YDK GO Core and Bundle Installation

The YDK requires Go version 1.9 or higher. Make sure that corresponding software is installed and environment variables GOROOT and GOPATH are properly set before the YDK installation. Follow System Requirements [here](https://github.com/CiscoDevNet/ydk-gen/tree/master/sdk/go#system-requirements).

##gNMI Requirements

In order to enable YDK support for gNMI protocol, which is optional, the following third party software must be installed prior to gNMI YDK component installation.

###Install protobuf and protoc

```
    wget https://github.com/google/protobuf/releases/download/v3.5.0/protobuf-cpp-3.5.0.zip
    unzip protobuf-cpp-3.5.0.zip
    cd protobuf-3.5.0
    ./configure
    make
    make check
    sudo make install
    sudo ldconfig
```

###Install gRPC

```
    git clone -b v1.9.1 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    make
    sudo make install
    sudo ldconfig
    cd -
```

###Run-time environment

There is an open issue with gRPC on Centos/Fedora, which requires an extra step before running any YDK gNMI application. See this issue on `GRPC GitHub <https://github.com/grpc/grpc/issues/10942#issuecomment-312565041>`_ 
for details. As a workaround, the YDK based application runtime environment must include setting of `LD_LIBRARY_PATH` variable:

```
    PROTO="/Your-Protobuf-and-Grpc-installation-directory"
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PROTO/grpc/libs/opt:$PROTO/protobuf-3.5.0/src/.libs:/usr/local/lib64
```

First, install [C++ core](https://github.com/CiscoDevNet/ydk-gen#second-step-generate--install-the-core). Then execute the below steps to install the ydk go core and bundle packages.
```
$ go get gopkg.in/stretchr/testify.v1
$ cd /your/path/to/ydk-gen
$ export GOPATH=/your/path/to/install/go/packages
$ ./generate.py --core --go
$ ./generate.py --bundle profiles/test/ydktest-cpp.json --go
```

Optionally you can install gNMI service:
```
$ ./generate.py --service profiles/services/gnmi-0.4.0.json --go
```

Then execute the below to run sample tests:
```
$ go run samples/cgo_path/cgo_path.go
$ go run samples/bgp_create/bgp_create.go
$ go run samples/bgp_read/bgp_read.go
```

To run YDK tests, please follow these commands:
```
$ cd /your/path/to/ydk-gen
$ cd sdk/cpp/core/tests/confd/ydktest
$ make clean
$ make all start
$ cd -
$ cd sdk/go/core/tests
$ go test
```
##gNMI Requirements

In order to enable YDK support for gNMI protocol, which is optional, the following third party software must be installed prior to gNMI YDK component installation.

###Install protobuf and protoc

```
    wget https://github.com/google/protobuf/releases/download/v3.5.0/protobuf-cpp-3.5.0.zip
    unzip protobuf-cpp-3.5.0.zip
    cd protobuf-3.5.0
    ./configure
    make
    make check
    sudo make install
    sudo ldconfig
```

###Install gRPC

```
    git clone -b v1.9.1 https://github.com/grpc/grpc
    cd grpc
    git submodule update --init
    make
    sudo make install
    sudo ldconfig
    cd -
```

###Run-time environment

There is an open issue with gRPC on Centos/Fedora, which requires an extra step before running any YDK gNMI application. See this issue on `GRPC GitHub <https://github.com/grpc/grpc/issues/10942#issuecomment-312565041>`_ 
for details. As a workaround, the YDK based application runtime environment must include setting of `LD_LIBRARY_PATH` variable:

```
    PROTO="/Your-Protobuf-and-Grpc-installation-directory"
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PROTO/grpc/libs/opt:$PROTO/protobuf-3.5.0/src/.libs:/usr/local/lib64
```

