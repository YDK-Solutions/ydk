cmake_minimum_required(VERSION 2.8.9)
cmake_policy(SET CMP0048 NEW)
project(ydk VERSION 0.5.4 LANGUAGES C CXX)

set(YDK_DESCRIPTION "YANG Development Kit Library. The library for  YDK API.")

# include custom Modules
set(CMAKE_MODULE_PATH ${CMAKE_SOURCE_DIR}/)

# set default build type if not specified by user
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE debug)
endif()

set(CMAKE_CXX_FLAGS         "${CMAKE_CXX_FLAGS} -Wall -Wextra")
set(CMAKE_CXX_FLAGS_RELEASE "-O2 -DNDEBUG")
set(CMAKE_CXX_FLAGS_DEBUG   "-g3 -O0 -fprofile-arcs -ftest-coverage")

# Set LibYDK_c variables
set(libydk_c_src golang.cpp)

add_library(ydk_c STATIC ${libydk_c_src})

set_property(TARGET ydk_c PROPERTY CXX_STANDARD 11)
set_property(TARGET ydk_c PROPERTY CXX_STANDARD_REQUIRED ON)

set(go_src path_sample.go)

add_custom_command(OUTPUT ${CMAKE_BINARY_DIR}/path_sample
  COMMAND go build -o ${CMAKE_BINARY_DIR}/path_sample ${go_src}
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  DEPENDS ydk_c ${go_src}
  COMMENT "Go build ${go_src}"
  VERBATIM
)

add_custom_target(path ALL DEPENDS ${CMAKE_BINARY_DIR}/path_sample)
