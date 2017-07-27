
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

#ifndef _GNMI_PROVIDER_H_
#define _GNMI_PROVIDER_H_

#include <fstream>
#include <libyang/libyang.h>

#include "gnmi_client.hpp"
#include "ietf_parser.hpp"
#include "path_api.hpp"
#include "ydk_yang.hpp"

namespace ydk 
{
    class gNMIClient;

    class gNMIServiceProvider : public path::ServiceProvider 
    {
    public:
        typedef struct SecureChannelArguments
        {
        std::shared_ptr<grpc::ChannelCredentials> channel_creds;
        grpc::ChannelArguments args;
            
        } SecureChannelArguments;

        gNMIServiceProvider(path::Repository & repo, std::string address);
        gNMIServiceProvider(std::string address);
        ~gNMIServiceProvider();

        std::shared_ptr<path::DataNode> invoke(path::Rpc& rpc) const;
        path::RootSchemaNode& get_root_schema() const;
        EncodingFormat get_encoding() const;

    private:
        std::shared_ptr<path::DataNode> handle_edit(path::Rpc& ydk_rpc, std::string operation) const;
        std::shared_ptr<path::DataNode> handle_read(path::Rpc& rpc, std::string operation) const;
        void initialize(path::Repository& repo, std::string address);
        std::string execute_payload(const std::string & payload, std::string operation) const;
        void print_root_paths(ydk::path::RootSchemaNode& rsn) const;
        void print_paths(ydk::path::SchemaNode& sn) const;

    private:
        std::unique_ptr<gNMIClient> client;
        std::shared_ptr<ydk::path::RootSchemaNode> root_schema;
        std::vector<std::string> server_capabilities;

    };
}
#endif /*_GNMI_PROVIDER_H_*/
