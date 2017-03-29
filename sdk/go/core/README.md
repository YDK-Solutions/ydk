### YDK GO

To build, first install [C++ core](https://github.com/CiscoDevNet/ydk-gen#second-step-generate--install-the-core). Then execute the below steps to build the go file.

```
$ mkdir build && cd build
$ cmake -DCMAKE_C_COMPILER=/usr/bin/clang -DCMAKE_CXX_COMPILER=/usr/bin/clang++ ..
$ make
```

Then execute the below to run the path_sample
```
$ ./path_sample [-verbose]
```