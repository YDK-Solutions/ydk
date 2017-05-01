package main

// #cgo CXXFLAGS: -g -std=c++11
// #cgo LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk_c -L./build -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #include "golang.h"
// #include <stdlib.h>
import "C"

import (
    "fmt"
    "flag"
    "unsafe"
)

func main() {
    verbosePtr := flag.Bool("verbose", false, "Enable verbose")
    flag.Parse()

    if *verbosePtr {
        C.EnableLogging()
    }

    var address *C.char = C.CString("localhost")
    defer C.free(unsafe.Pointer(address))
    var username *C.char = C.CString("admin")
	defer C.free(unsafe.Pointer(username))
	var password *C.char = C.CString("admin")
	defer C.free(unsafe.Pointer(password))
	var path *C.char = C.CString("/usr/local/share/ydktest@0.1.0/")
	defer C.free(unsafe.Pointer(path))
	var runner_path *C.char = C.CString("ydktest-sanity:runner")
	defer C.free(unsafe.Pointer(runner_path))
	var number_path *C.char = C.CString("ytypes/built-in-t/number8")
	defer C.free(unsafe.Pointer(number_path))
	var number_value *C.char = C.CString("2")
	defer C.free(unsafe.Pointer(number_value))
	var create_path *C.char = C.CString("ydk:create")
	defer C.free(unsafe.Pointer(create_path))
	var read_path *C.char = C.CString("ydk:read")
	defer C.free(unsafe.Pointer(read_path))
	var entity_path *C.char = C.CString("entity")
	defer C.free(unsafe.Pointer(entity_path))
	var filter_path *C.char = C.CString("filter")
	defer C.free(unsafe.Pointer(filter_path))

	codec := C.CodecServiceInit()
	defer C.CodecServiceFree(codec)
    repo := C.RepositoryInitWithPath(path)
    defer C.RepositoryFree(repo)
	provider := C.NetconfServiceProviderInit(repo, address, username, password, 12022)
	defer C.NetconfServiceProviderFree(provider)
	root_schema := C.NetconfServiceProviderGetRootSchema(provider)

	runner := C.RootSchemaNodeCreate(root_schema, runner_path)

    C.DataNodeCreate(runner, number_path, number_value)
    var create_xml *C.char = C.CodecServiceEncode(codec, runner, C.XML, 0)
    defer C.free(unsafe.Pointer(create_xml))

    create_rpc := C.RootSchemaNodeRpc(root_schema, create_path)
    input := C.RpcInput(create_rpc)
    C.DataNodeCreate(input, entity_path, create_xml)
    C.RpcExecute(create_rpc, provider)

    read_rpc := C.RootSchemaNodeRpc(root_schema, read_path)
    input = C.RpcInput(read_rpc)
	runner_filter := C.RootSchemaNodeCreate(root_schema, runner_path)
	var read_xml *C.char = C.CodecServiceEncode(codec, runner_filter, C.XML, 0)
	defer C.free(unsafe.Pointer(read_xml))

    C.DataNodeCreate(input, filter_path, read_xml)
    read_data := C.RpcExecute(read_rpc, provider)
    if read_data == nil {
        fmt.Println("uhoh");
    }
    var data *C.char = C.CodecServiceEncode(codec, read_data, C.XML, 1)
    defer C.free(unsafe.Pointer(data))
    s := C.GoString(data)
    fmt.Println("Read data:\n", s);
}
