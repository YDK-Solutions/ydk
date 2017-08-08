### YDK GO

To build, first install [C++ core](https://github.com/CiscoDevNet/ydk-gen#second-step-generate--install-the-core). Then execute the below steps to install the ydk package (until we have the YDK go repository on github, this is a little unusal way to install go packages).

```
$ export GOPATH=/your/path/to/install/go/packages
$ mkdir -p $GOPATH/src/github.com/CiscoDevNet/ydk-go/ydk
$ cp -r ydk/* $GOPATH/src/github.com/CiscoDevNet/ydk-go/ydk
$ cd ../packages
$ cp -r ydk/*  $GOPATH/src/github.com/CiscoDevNet/ydk-go/ydk
```

Then execute the below to run the samples
```
$ cd samples/cgo_path && go run cgo_path
```
```
$ cd samples/bgp_create && go run bgp_create
```
```
$ cd samples/bgp_read && go run bgp_read
```

To test gen code:
```bash
$ go get gopkg.in/stretchr/testify.v1
$ cd tests
$ go test gen_code_test.go
```
