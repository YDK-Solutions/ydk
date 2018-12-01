# - Try to find LibYDK
# Once done this will define
#
#  LIBYDK_GNMI_FOUND - system has LibYDK
#  LIBYDK_GNMI_INCLUDE_DIRS - the LibYDK include directory
#  LIBYDK_GNMI_LIBRARIES - Link these to use LibYANG LibPCRE Lib
#
# ####################################################################
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
######################################################################

if (LIBYDK_GNMI_LIBRARIES AND LIBYDK_GNMI_INCLUDE_DIRS)
  # in cache already
  set(LIBYDK_GNMI_FOUND TRUE)
else (LIBYDK_GNMI_LIBRARIES AND LIBYDK_GNMI_INCLUDE_DIRS)

  find_path(LIBYDK_GNMI_INCLUDE_DIR
    NAMES
      ydk/gnmi_provider.hpp
      ydk/gnmi_client.hpp
      ydk/gnmi_service.hpp
      ydk/gnmi_util.hpp
      ydk/gnmi_path_api.hpp
    PATHS
      /usr/include
      /usr/local/include
      /opt/local/include
      /sw/include
      ${CMAKE_INCLUDE_PATH}
      ${CMAKE_INSTALL_PREFIX}/include
  )
  
  find_library(LIBYDK_GNMI_LIBRARY
    NAMES
      libydk_gnmi
    PATHS
      /usr/lib
      /usr/lib64
      /usr/local/lib
      /usr/local/lib64
      /opt/local/lib
      /sw/lib
      ${CMAKE_LIBRARY_PATH}
      ${CMAKE_INSTALL_PREFIX}/lib
  )

  if (LIBYDK_GNMI_INCLUDE_DIR AND LIBYDK_GNMI_LIBRARY)
    set(LIBYDK_GNMI_FOUND TRUE)
  else (LIBYDK_GNMI_INCLUDE_DIR AND LIBYDK_GNMI_LIBRARY)
    set(LIBYDK_GNMI_FOUND FALSE)
  endif (LIBYDK_GNMI_INCLUDE_DIR AND LIBYDK_GNMI_LIBRARY)

  set(LIBYDK_GNMI_INCLUDE_DIRS ${LIBYDK_INCLUDE_DIR})
  set(LIBYDK_GNMI_LIBRARIES ${LIBYDK_LIBRARY})

  # show the LIBYDK_INCLUDE_DIRS and LIBYDK_LIBRARIES variables only in the advanced view
  mark_as_advanced(LIBYDK_GNMI_INCLUDE_DIRS LIBYDK_GNMI_LIBRARIES)

endif (LIBYDK_GNMI_LIBRARIES AND LIBYDK_GNMI_INCLUDE_DIRS)

