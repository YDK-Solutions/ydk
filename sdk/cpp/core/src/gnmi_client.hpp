/*  ----------------------------------------------------------------
 Copyright 2016 Cisco Systems

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 ------------------------------------------------------------------*/

#ifndef _YDK_GNMI_CLIENT_H_
#define _YDK_GNMI_CLIENT_H_
#include <map>
#include <string>
#include <vector>
#include <iostream>
#include <sstream>

#include "errors.hpp"

#include <grpc++/grpc++.h>
#include <libgnmi/gnmi.grpc.pb.h>

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using gnmi::gNMI;
using gnmi::CapabilityRequest;
using gnmi::CapabilityResponse;
using gnmi::GetRequest;
using gnmi::GetResponse;

typedef struct capabilities 
{
    capabilities(std::vector<std::string> cas){caps=cas;}
    std::vector<std::string> caps;
} capabilities;

namespace ydk 
{
    class gNMIClient 
    {
    public:
        gNMIClient(std::shared_ptr<Channel> channel);
        ~gNMIClient();

        std::vector<std::string> Capabilities();
        void Get();
        int connect();
        std::string execute_wrapper(const std::string & payload);
        std::string execute_payload(const std::string & payload, const GetRequest& request, GetResponse* response);
        GetRequest ParseGetRequest(GetRequest request, std::string payload);
        std::string ParseGetResponse(GetResponse* response); 
        std::unique_ptr<gNMI::Stub> stub_;

    private:
        void init_capabilities();
        bool get_capabilities(const CapabilityRequest& request, CapabilityResponse* response);
        std::vector<std::string> capabilities;
    };
}

#endif /* _YDK_GNMI_CLIENT_H_ */
