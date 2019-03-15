## YDK GO Core and Bundle Installation

The YDK requires Go version 1.9 or higher. Make sure that corresponding software is installed and environment variables GOROOT and GOPATH are properly set before the YDK installation. Follow System Requirements [here](https://github.com/CiscoDevNet/ydk-gen/tree/master/sdk/go#system-requirements).

### YDK core and model bundle installation

First, install [C++ core](https://github.com/CiscoDevNet/ydk-gen#second-step-generate--install-the-core). Then execute the below steps to install the ydk go core and bundle packages.

```
$ go get gopkg.in/stretchr/testify.v1
$ cd /your-path-to-ydk-gen
$ export GOPATH=/your-path-to-go-packages-installation-directory
$ ./generate.py -i --core --go
$ ./generate.py -i --bundle profiles/test/ydktest-cpp.json --go
```

You can test your installation now. Execute the below commands to run sample tests.

```
$ go run samples/cgo_path/cgo_path.go
$ go run samples/bgp_create/bgp_create.go
$ go run samples/bgp_read/bgp_read.go
```

To run YDK sanity tests execute:

```
$ cd /your-path-to-ydk-gen
$ cd sdk/go/core/tests
$ go test
```
