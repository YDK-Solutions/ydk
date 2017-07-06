#include <iostream>
#include <fstream>
#include <memory>
#include <string>
#include <sstream>
#include <vector>
#include <iterator>

#include "errors.hpp"
#include "gnmi_client.hpp"
#include "types.hpp"
#include "logger.hpp"
#include "path_api.hpp"
#include "json.hpp"

#include <libyang/libyang.h>
#include <grpc++/grpc++.h>
#include <grpc++/client_context.h>
#include <grpc++/create_channel.h>
#include <libgnmi/gnmi.grpc.pb.h>
#include <libgnmi/gnmi.pb.h>
#include "entity_data_node_walker.hpp"
#include "ietf_parser.hpp"
#include "ydk_yang.hpp"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using gnmi::gNMI;
using gnmi::CapabilityRequest;
using gnmi::CapabilityResponse;
using gnmi::GetRequest;
using gnmi::GetResponse;

using namespace std;
using namespace ydk;
using json = nlohmann::json;

typedef vector<string> StringVec;
std::ifstream capfile("/usr/local/share/ydk/0.0.0.0\:50051/capabilities.txt");

gNMIClient::gNMIClient(shared_ptr<Channel> channel)
	: stub_(gNMI::NewStub(channel)){}

StringVec gNMIClient::Capabilities() 
{
	CapabilityRequest request;
	CapabilityResponse response;
	gNMIClient::get_capabilities(request, &response);
	return capabilities;
}

int gNMIClient::connect() 
{
    grpc::CreateChannel("0.0.0.0:50051", grpc::InsecureChannelCredentials());
    init_capabilities();
    return EXIT_SUCCESS;
}

void gNMIClient::init_capabilities() 
{
	std::string cap;
	while (std::getline(capfile, cap)) 
	{
		capabilities.push_back(cap);
		//std::cout << cap;
	}
}

bool gNMIClient::get_capabilities(const CapabilityRequest& request, CapabilityResponse* response) 
{
	grpc::ClientContext context;
	grpc::Status status = stub_->Capabilities(&context, request, response);
	if (!status.ok()) 
	{
	  std::cout << "Capabilities RPC Failed." << std::endl;
	  return false;
	}
	if (!response->gnmi_version().size() > 0) 
	{
	  std::cout << "Server has not returned gnmi version" << std::endl;
	  return false;
	} else 
	{
		std::cout << "Receive capability response with gNMI version: " << response->gnmi_version() << std::endl;
		for (unsigned int i = 0, n = response->supported_models_size(); i < n; i++) 
		{
			std::cout << "supported_models size" << response->supported_models_size() << std::endl;
			std::cout << "supported_models";
			response->supported_models(i);
		}
		for (int i = 0, n = response->supported_encodings_size(); i < n; i++) 
		{
			std::cout << "supported_models size" << response->supported_encodings_size() << std::endl;
			std::cout << "supported_models";
			response->supported_encodings(i);
		}
	}
    return true;
}

void gNMIClient::Get() {}

GetRequest gNMIClient::ParseGetRequest(GetRequest request, string payload) 
{	
	::gnmi::Path* prefix = new ::gnmi::Path();
	::gnmi::Path* path;

	// create parsable payload 
	std::string payload_recieved = "{";
	payload_recieved.append(payload);
	payload_recieved.append("}");
	std::cout << payload << endl;
	auto payload_to_parse = json::parse(payload_recieved);
	
	// capture filter value
	std::string payload_filter = payload_to_parse.value("/rpc/ietf-netconf:get-config/filter"_json_pointer, "Empty Filter");

	// get prefix and path values
	std::istringstream filter_value(payload_filter);
	std::string token;
	std::vector<std::string> prefix_container;
	std::vector<std::string> path_container;
	while(std::getline(filter_value, token, '"')) 
	{
		if(token.find("{")==std::string::npos && token.find("}")==std::string::npos) 
		{
			std::istringstream pref(token);
			std::string subtoken;
			while(std::getline(pref, subtoken, ':')) 
			{
				prefix_container.push_back(subtoken);
			}
		}
	}

	// populate values into request payload
	// prefix - {"openconfig-bgp", "bgp"}
	for(int i = 0; i < prefix_container.size(); ++i) 
	{
		prefix->add_element(prefix_container[i]);
	}
	request.set_allocated_prefix(prefix);

	// path - {"global", "config", "as"}
	path = request.add_path();
	path->add_element("*");

	// type (CONFIG/STATE/OPERATIONAL) - "CONFIG"
	request.set_type(::gnmi::GetRequest::CONFIG);

	// encoding - default JSON
	request.set_encoding(::gnmi::Encoding::JSON);
	return request;
}

std::string gNMIClient::ParseGetResponse(GetResponse* response) {
	/* 
	Notification Message:
	Timestamp - Timestamp values MUST be represented as the number of nanoseconds since the Unix epoch 
		(January 1st 1970 00:00:00 UTC). The value MUST be encoded as a signed 64-bit integer (int64)[^1].
 	Prefix - {"openconfig-bgp", "bgp"}
 	Update 
		Path - {"global", "config", "as"}
		Value - bytes ("65172"), type (0)
	*/
	::gnmi::Notification notification;
	::gnmi::Path prefix;
	::gnmi::Path path;
	::gnmi::Update update;
	::gnmi::Value value;
	std::string element;
	std::string reply_to_parse;
	int element_size;
	std::string prefix_to_prepend;
	std::string path_to_prepend;
	int notification_size = response->notification_size();

	// read from all response fields and create reply
	for(int i = 0; i < notification_size; ++i) 
	{
		notification = response->notification(i);

		if(notification.has_prefix()) {
			prefix = notification.prefix();
			element_size = prefix.element_size();
			for(int j = 0; j < element_size; ++j)
            {
              	element.append("\"");
		        element.append(prefix.element(j));
		        element.append("\" ");
		        prefix_to_prepend.append(prefix.element(j));
		        if(j != (element_size - 1))
		        {
		        	prefix_to_prepend.append(":");
		        }
            }
            prefix_to_prepend.append("\": {");
		}

		int update_size = notification.update_size();
		
		if(notification.update_size() != 0) {
	        for(int k = 0; k < update_size; ++k) 
	        {
	            update = notification.update(k);
	            
	            if(update.has_path()) 
	            {
	            	path = update.path();
	           		element_size = path.element_size();
	             	element.clear();
	              	for(int l = 0; l < element_size; ++l) 
	              	{
		                element.append("\"");
		                element.append(path.element(l));
		                element.append("\" ");
		                path_to_prepend.append(path.element(l));
		                if(l != element_size-1)
		                {
		                	path_to_prepend.append(":");
		                }
	              	}
	              path_to_prepend.append("\": {");
	            }

	            if(update.has_value()) 
	            {
	              	value = update.value();
	              	reply_to_parse.append("\"data\": {\"");
	              	reply_to_parse.append(prefix_to_prepend); 
	              	reply_to_parse.append(value.value());
	              	reply_to_parse.append(" } } }");
            	}
           	}   
        }
	}
	return reply_to_parse;
}

std::string gNMIClient::execute_wrapper(const string & payload)
{
	GetRequest request;
	GetResponse response;
	request = gNMIClient::ParseGetRequest(request, payload);
	std::cout << "===============Get Request Sent================" << std::endl;
	std::cout << request.DebugString() << std::endl;
	std::string reply = execute_payload(payload, request, &response);
	return reply;
}

std::string gNMIClient::execute_payload(const string & payload, const GetRequest& request, GetResponse* response)
{
	grpc::ClientContext context;
	grpc::Status status;

	status = stub_->Get(&context, request, response);
	std::cout << "=============Get Response Received=============" << std::endl;
	std::cout << response->DebugString() << std::endl;

	if (!status.ok()) 
	{
		return status.error_message();
	} 
	else 
	{
		std::cout << "===================RPC Passed===================" << std::endl;
		std::string reply = gNMIClient::ParseGetResponse(response);	
		return reply;
	}
}

gNMIClient::~gNMIClient() {}
