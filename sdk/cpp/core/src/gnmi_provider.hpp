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
#include "path_api.hpp"
#include "service_provider.hpp"

namespace ydk 
{
    class gNMIClient;

    class gNMIServiceProvider : public ServiceProvider 
    {
    public:
        typedef struct SecureChannelArguments
        {
            std::shared_ptr<grpc::ChannelCredentials> channel_creds;
            grpc::ChannelArguments args;
        } SecureChannelArguments;

        gNMIServiceProvider(path::Repository & repo, const std::string& address);
        gNMIServiceProvider(path::Repository & repo, const std::string& address, bool is_secure);
        gNMIServiceProvider(const std::string& address);
        gNMIServiceProvider(const std::string& address, bool is_secure);
        virtual ~gNMIServiceProvider();

        virtual EncodingFormat get_encoding() const;
        virtual const path::Session& get_session() const;
        std::vector<std::string> get_capabilities() const;

    private:
        const path::gNMISession session;
    };
}
#endif 
/*_GNMI_PROVIDER_H_*/
