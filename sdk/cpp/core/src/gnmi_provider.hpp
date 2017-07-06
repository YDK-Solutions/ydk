
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

#include <memory>
#include <string>
#include <libgnmi/gnmi.grpc.pb.h>
#include <libgnmi/gnmi.pb.h>
#include "gnmi_client.hpp"
#include "path_api.hpp"

namespace ydk {

	class gNMIClient;

	class gNMIServiceProvider : public path::ServiceProvider {
	public:
			gNMIServiceProvider(path::Repository & repo,
								   std::string address,
								   std::string username,
								   std::string password);
			gNMIServiceProvider(std::string address,
								   std::string username,
								   std::string password);
			~gNMIServiceProvider();
			path::RootSchemaNode& get_root_schema() const;
			std::shared_ptr<path::DataNode> invoke(path::Rpc& rpc) const;
			EncodingFormat get_encoding() const;
			void print_paths(ydk::path::SchemaNode& sn) const;
			void print_data_paths(ydk::path::DataNode& dn) const;
			void print_root_paths(ydk::path::RootSchemaNode& rsn) const;

	private:
		    std::shared_ptr<path::DataNode> handle_read(path::Rpc& rpc) const;
			void initialize(path::Repository& repo);
			std::string execute_payload(const std::string & payload) const;

	private:
			std::unique_ptr<gNMIClient> client;
			//std::unique_ptr<path::ModelProvider> model_provider;
			std::shared_ptr<ydk::path::RootSchemaNode> root_schema;
			std::vector<std::string> server_capabilities;

	};
}
#endif /*_GNMI_PROVIDER_H_*/
