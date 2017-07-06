
/// YANG Development Kit
// Copyright 2016 Cisco Systems. All rights reserved
//
////////////////////////////////////////////////////////////////
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  Yo::u may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
//////////////////////////////////////////////////////////////////

#include <iostream>
#include <sstream>
#include <memory>
#include <libyang/libyang.h>
#include <grpc++/grpc++.h>
#include <libgnmi/gnmi.grpc.pb.h>
#include <libgnmi/gnmi.pb.h>
#include "entity_data_node_walker.hpp"
#include "errors.hpp"
#include "ietf_parser.hpp"
#include "gnmi_client.hpp"
#include "gnmi_provider.hpp"
//#include "gnmi_model_provider.hpp"
#include "types.hpp"
#include "ydk_yang.hpp"
#include "logger.hpp"

using grpc::Channel;

using namespace std;
using namespace ydk;

namespace ydk
{
	//const char* CANDIDATE = "urn:ietf:params:netconf:capability:candidate:1.0";
	static path::SchemaNode* get_schema_for_operation(path::RootSchemaNode& root_schema, string operation);

	static shared_ptr<path::Rpc> create_rpc_instance(path::RootSchemaNode & root_schema, string rpc_name);
	static path::DataNode& create_rpc_input(path::Rpc & gnmi_rpc);
	static void create_input_source(path::DataNode & input, bool config);

	static string get_read_rpc_name(bool config);
	static bool is_config(path::Rpc & rpc);
	static string get_filter_payload(path::Rpc & ydk_rpc);
	static string get_gnmi_payload(path::DataNode & input, string data_tag, string data_value);
	static std::shared_ptr<path::DataNode> handle_read_reply(string reply, path::RootSchemaNode & root_schema);

	gNMIServiceProvider::gNMIServiceProvider(string address, string username, string password)
		: client(make_unique<gNMIClient>(grpc::CreateChannel(address, grpc::InsecureChannelCredentials())))
	{
		path::Repository repo;
		initialize(repo);
		YLOG_INFO("Connected to {} on port {} using ssh", address);
	}

	gNMIServiceProvider::gNMIServiceProvider(path::Repository & repo, string address, string username, string password)
		: client(make_unique<gNMIClient>(grpc::CreateChannel(address, grpc::InsecureChannelCredentials())))
	{
		initialize(repo);
		YLOG_INFO("Connected to {} on port {} using ssh", address);
	}

	gNMIServiceProvider::~gNMIServiceProvider() = default;

	void gNMIServiceProvider::initialize(path::Repository & repo) {
		IetfCapabilitiesParser capabilities_parser{};
		client->connect();
		server_capabilities = client->Capabilities();
		
		root_schema = repo.create_root_schema(capabilities_parser.parse(server_capabilities));

		if(root_schema.get() == nullptr)
		{
			YLOG_ERROR("Root schema cannot be obtained");
			throw(YCPPIllegalStateError{"Root schema cannot be obtained"});
		}
	}

	EncodingFormat gNMIServiceProvider::get_encoding() const
	{
		return EncodingFormat::JSON;
	}

	path::RootSchemaNode& gNMIServiceProvider::get_root_schema() const
	{
		return *root_schema;
	}

	std::shared_ptr<path::DataNode> gNMIServiceProvider::invoke(path::Rpc& rpc) const
	{
		path::SchemaNode* read_schema = get_schema_for_operation(*root_schema, "ydk:read");

		//for now we only support crud rpc's
    	path::SchemaNode* rpc_schema = &(rpc.schema());
		//print_paths(*rpc_schema);

		std::shared_ptr<path::DataNode> datanode = nullptr;
		
		if(rpc_schema == read_schema)
	    {
	    	//print_data_paths(rpc.input());
	        return handle_read(rpc);
	    }
	    else
	    {
	        YLOG_ERROR("rpc is not supported");
	        throw(YCPPOperationNotSupportedError{"rpc is not supported!"});
	    }

	    return datanode;
	}

	void gNMIServiceProvider::print_paths(ydk::path::SchemaNode& sn) const
	{
	    std::cout << sn.path() << std::endl;
	    for(auto const& p : sn.children())
	        print_paths(*p);
	}

	void gNMIServiceProvider::print_data_paths(ydk::path::DataNode& dn) const
	{
	    std::cout << dn.path() << std::endl;
	    for(auto const& p : dn.children())
	        print_data_paths(*p);
	}

	void gNMIServiceProvider::print_root_paths(ydk::path::RootSchemaNode& rsn) const
	{
	    std::cout << rsn.path() << std::endl;
	    for(auto const& p : rsn.children())
	        print_paths(*p);
	}


	std::shared_ptr<path::DataNode> gNMIServiceProvider::handle_read(path::Rpc& ydk_rpc) const
	{
	    //for now we only support crud rpc's
	    //print_data_paths(ydk_rpc.input());
	    bool config = is_config(ydk_rpc);
	    auto gnmi_rpc = create_rpc_instance(*root_schema, get_read_rpc_name(config));
	    auto & input = create_rpc_input(*gnmi_rpc);
	    //print_data_paths(input);
	    create_input_source(input, config);
	    std::string filter_value = get_filter_payload(ydk_rpc);

	    std::string gnmi_payload = get_gnmi_payload(input, "filter", filter_value);
		//print_root_paths(*root_schema);
		//client->execute_wrapper(gnmi_payload);
		std::shared_ptr<path::DataNode> datanode = nullptr;
	    return handle_read_reply(execute_payload(gnmi_payload), *root_schema);
	}

	
	static std::shared_ptr<path::DataNode> handle_read_reply(string reply, path::RootSchemaNode & root_schema)
	{
	    path::CodecService codec_service{};
	    auto empty_data = reply.find("data");
	    if(empty_data == std::string::npos)
	    {
	        YLOG_INFO( "Found empty data tag");
	        return nullptr;
	    }

	    auto data_start = reply.find("\"data\":");
	    if(data_start == std::string::npos)
	    {
	        YLOG_ERROR( "Can't find data tag in reply sent by device {}", reply);
	        throw(YCPPServiceProviderError{reply});
	    }
	    data_start+= sizeof("\"data\":") - 1;
	    auto data_end = reply.find_last_of("}");
	    if(data_end == std::string::npos)
	    {
	        YLOG_ERROR( "No end data tag found in reply sent by device {}", reply);
	        throw(YCPPError{"No end data tag found"});
	    }

	    string data = reply.substr(data_start, data_end-data_start);

		auto datanode = std::shared_ptr<path::DataNode>(codec_service.decode(root_schema, data, EncodingFormat::JSON));

	    if(!datanode)
	    {
	        YLOG_ERROR( "Codec service failed to decode datanode");
	        throw(YCPPError{"Problems deserializing output"});
	    }
	    return datanode;
	}

	static bool is_config(path::Rpc & rpc)
	{
	    if(!rpc.input().find("only-config").empty())
	    {
	        return true;
	    }
	    return false;
	}

	static shared_ptr<path::Rpc> create_rpc_instance(path::RootSchemaNode & root_schema, string rpc_name)
	{
	    auto rpc = shared_ptr<path::Rpc>(root_schema.rpc(rpc_name));
	    if(rpc == nullptr){
	        YLOG_ERROR("Cannot create payload for RPC: {}", rpc_name);
	        throw(YCPPIllegalStateError{"Cannot create payload for RPC: "+ rpc_name});
	    }
	    return rpc;
	}

	static string get_read_rpc_name(bool config)
	{
	    if(config)
	    {
	        return "ietf-netconf:get-config";
	    }
	    return "ietf-netconf:get";
	}

	static path::DataNode& create_rpc_input(path::Rpc & gnmi_rpc)
	{
	    return gnmi_rpc.input();
	}

	static void create_input_source(path::DataNode & input, bool config)
	{
	    if(config)
	    {
	        input.create("source/running");
	    }
	}

	static string get_filter_payload(path::Rpc & ydk_rpc)
	{
	    auto entity = ydk_rpc.input().find("filter");
	    if(entity.empty()){
	        YLOG_ERROR("Failed to get entity node.");
	        throw(YCPPInvalidArgumentError{"Failed to get entity node"});
	    }

	    auto datanode = entity[0];
	    return datanode->get();
	}

	static string get_gnmi_payload(path::DataNode & input, string data_tag, string data_value)
	{
	    path::CodecService codec_service{};
	    input.create(data_tag, data_value);

	    string payload{"\"rpc\":"};
	    payload+=codec_service.encode(input, EncodingFormat::JSON, true);
	    //payload+="</rpc>";
	    YLOG_INFO("===========Generating Target Payload============");
	    YLOG_INFO(payload.c_str());
	    YLOG_INFO("\n");
	    std::cout << "\n=============Generating Target Payload============" << std::endl;
	    return payload;
	}

	std::string gNMIServiceProvider::execute_payload(const std::string & payload) const
	{
	    std::string reply = client->execute_wrapper(payload);
	    YLOG_INFO("=============Reply payload received from device=============");
	    YLOG_INFO(reply.c_str());
	    YLOG_INFO("\n");
	    std::cout << "\n================Payload to Decode================" << std::endl;
	    std::cout << reply << std::endl;
	    return reply;
	}

	static path::SchemaNode* get_schema_for_operation(path::RootSchemaNode & root_schema, string operation)
	{
		auto c = root_schema.find(operation);
		if(c.empty())
		{
			YLOG_ERROR("CRUD read rpc schema not found!");
			throw(YCPPIllegalStateError{"CRUD read rpc schema not found!"});
		}
		return c[0];
	}
}
