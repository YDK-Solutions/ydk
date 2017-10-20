### YDK GO

To build, first install [C++ core](https://github.com/CiscoDevNet/ydk-gen#second-step-generate--install-the-core). Then execute the below steps to install the ydk package (until we have the YDK go repository on github, this is a little unusal way to install go packages).

```
$ export GOPATH=/your/path/to/install/go/packages
$ mkdir -p $GOPATH/src/github.com/CiscoDevNet/ydk-go/ydk
$ cp -r ydk/* $GOPATH/src/github.com/CiscoDevNet/ydk-go/ydk
$ cd ../../..
$ ./generate.py --bundle profiles/test/ydktest-cpp.json --go
$ cp -r gen-api/go/ydktest-bundle/ydk/models/*  $GOPATH/src/github.com/CiscoDevNet/ydk-go/ydk/models
$ cd -
```

Then execute the below to run the samples
```
$ go run samples/cgo_path/cgo_path.go
```
```
$ go run samples/bgp_create/bgp_create.go
```
```
$ go run samples/bgp_read/bgp_read.go
```

To run tests:
```bash
$ cd ../../..
$ ./generate.py --bundle profiles/test/ydktest-cpp.json --go
$ cp -r gen-api/go/ydktest-bundle/ydk/models/* $GOPATH/src/github.com/CiscoDevNet/ydk-go/ydk/models
$ go get gopkg.in/stretchr/testify.v1
$ cd sdk/go/core/tests
$ go test
```
