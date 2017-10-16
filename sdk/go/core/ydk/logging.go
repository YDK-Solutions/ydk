/*
 * ------------------------------------------------------------------
 * YANG Development Kit
 * Copyright 2017 Cisco Systems. All rights reserved
 *
 *----------------------------------------------
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http:www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *----------------------------------------------
 */
package ydk

// #cgo CXXFLAGS: -g -std=c++11
// #cgo darwin LDFLAGS:  -fprofile-arcs -ftest-coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lpython -lc++
// #cgo linux LDFLAGS:  -fprofile-arcs -ftest-coverage --coverage -lydk -lxml2 -lxslt -lpcre -lssh -lssh_threads -lcurl -lstdc++ -lpython2.7 -ldl
//#include <ydk/ydk.h>
import "C"
import "fmt"

// LogLevel indicates what level of logging to output
type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warning
	Error
)

// EnableLogging enables logging on the code that runs in C++
func EnableLogging(level LogLevel) {
	switch level {
	case Debug:
		C.EnableLogging(C.DEBUG)

	case Info:
		C.EnableLogging(C.INFO)

	case Warning:
		C.EnableLogging(C.WARNING)

	case Error:
		C.EnableLogging(C.ERROR)
	}
}

func YLogInfo(msg string){
	msg = fmt.Sprintf("[Go] %s", msg)
    C.YLogInfo(C.CString(msg));
}

func YLogDebug(msg string){
	msg = fmt.Sprintf("[Go] %s", msg)
    C.YLogDebug(C.CString(msg));
}

func YLogWarn(msg string){
	msg = fmt.Sprintf("[Go] %s", msg)
    C.YLogWarn(C.CString(msg));
}

func YLogError(msg string){
	msg = fmt.Sprintf("[Go] %s", msg)
    C.YLogError(C.CString(msg));
}
