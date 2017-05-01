cmake_minimum_required(VERSION 2.8.9)
cmake_policy(SET CMP0048 NEW)
set(YDK_UNITTEST_TARGET_NAME "ydk_core_test")

set(core_tests_src golang_test.cpp               
               main.cpp)

set(CMAKE_CXX_FLAGS         "${CMAKE_CXX_FLAGS} -Wall -Wextra")
set(CMAKE_CXX_FLAGS_RELEASE "-O2 -DNDEBUG")
set(CMAKE_CXX_FLAGS_DEBUG   "-g3 -O0 -fprofile-arcs -ftest-coverage")

add_executable(${YDK_UNITTEST_TARGET_NAME} ${core_tests_src})
set_property(TARGET ${YDK_UNITTEST_TARGET_NAME} PROPERTY CXX_STANDARD 11)
set_property(TARGET ${YDK_UNITTEST_TARGET_NAME} PROPERTY CXX_STANDARD_REQUIRED ON)
target_link_libraries(${YDK_UNITTEST_TARGET_NAME} ydk
                                                 ydk_c                                                 
                                                curl
                                                xslt
                                                pcre
                                                ssh_threads
                                                xml2
                                                ssh
                                                python2.7)
add_test(NAME ${YDK_UNITTEST_TARGET_NAME} COMMAND $<TARGET_FILE:${YDK_UNITTEST_TARGET_NAME}>)
