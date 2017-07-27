cmake_minimum_required(VERSION 2.8.9)
cmake_policy(SET CMP0048 NEW)

set(YDK_CORE_SAMPLES_DESCRIPTION "Samples")

set(samples bgp_gnmi_create 
            bgp_gnmi_read)

find_library(ydk_openconfig_location ydk_openconfig)
find_library(ydk_cisco_ios_xr_location ydk_cisco_ios_xr)
find_library(ydk_location ydk)

link_libraries(${CMAKE_THREAD_LIBS_INIT})
foreach(sample IN LISTS samples)
    add_executable(${sample} ${sample}.cpp args_parser.cpp)
    set_property(TARGET ${sample} PROPERTY CXX_STANDARD 11)
    set_property(TARGET ${sample} PROPERTY CXX_STANDARD_REQUIRED ON)
    target_link_libraries(${sample} ydk ${ydk_openconfig_location} ${ydk_cisco_ios_xr_location})
endforeach(sample)

