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

#include <string.h>
#include <iostream>
#include <sys/time.h>

#include <ydk/gnmi_client.hpp>
#include <ydk/errors.hpp>

#include "../../core/src/catch.hpp"
#include "../../core/tests/config.hpp"

using namespace ydk;
using namespace std;

TEST_CASE("gnmi_secure_server")
{
    string model_dir = TEST_HOME;
    auto pos = model_dir.find("sdk");
    string gnmi_server_sert = model_dir.substr(0, pos) + "test/gnmi_server/keys/ems.pem";

	gNMIClient client("127.0.0.1", 50051, "admin", "admin", gnmi_server_sert);

	vector<string> caps = client.get_capabilities();
	for (auto cap : caps)
	    cout << cap << endl;
}

/* To Do
TEST_CASE("gnmi_xr")	// Not running
{
    gNMIClient client(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()), "admin", "admin");

    int ok = 0;

    int result = client.connect();

    std::cout << "result: " << result << std::endl;
    REQUIRE(result == ok);

    struct timeval t1, t2;
    double elapsedTime;

    // start timer
    gettimeofday(&t1, NULL);

    string reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:get-config":{"source":{"running":[null]},"filter":"{\"openconfig-bgp:bgp\":{}}"}})", "read", true);
    
    // stop timer
    gettimeofday(&t2, NULL);

    // compute and print the elapsed time in millisec
    elapsedTime = (t2.tv_sec - t1.tv_sec) * 1000.0;      // sec to ms
    elapsedTime += (t2.tv_usec - t1.tv_usec) / 1000.0;   // us to ms
    cout << elapsedTime << " ms.\n";

    REQUIRE(result == ok);
}

TEST_CASE("gnmi_create")
{
    gNMIClient client(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));

    int ok = 0;

    int result = client.connect();
    REQUIRE(result == ok);
}

TEST_CASE("gnmi_edit_get_config")	// Not working
{
    gNMIClient client(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));

    int ok = 0;

    int result = client.connect();
    REQUIRE(result == ok);

    string reply = client.execute_wrapper(
     R"("rpc":{"ietf-netconf:edit-config":{"target":{"running":[null]},"error-option":"rollback-on-error","config":"{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}}}"}})", "create", true);
    REQUIRE(NULL != strstr(reply.c_str(), "Success"));

    reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:get-config":{"source":{"running":[null]},"filter":"{\"openconfig-bgp:bgp\":{}}"}})", "read", true);
    
    REQUIRE(NULL != strstr(reply.c_str(), R"({"as":65172})"));

    reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:edit-config":{"target":{"running":[null]},"error-option":"rollback-on-error","config":"{\"openconfig-bgp:bgp\":{}}"}})", "delete", true);
    REQUIRE(NULL != strstr(reply.c_str(), "Success"));

    REQUIRE(result == ok);
}


TEST_CASE("gnmi_device_not_connected_execute")	// Not working
{
    gNMIClient client(grpc::CreateChannel("", grpc::InsecureChannelCredentials()));

    try
    {
        string s = client.execute_wrapper( R"("rpc":{"ietf-netconf:edit-config":{"target":{"running":[null]},"error-option":"rollback-on-error","config":"{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172,\"router-id\":\"1.2.3.4\"}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"6.7.8.9\",\"config\":{\"local-as\":65001,\"neighbor-address\":\"6.7.8.9\",\"peer-as\":65001,\"peer-group\":\"IBGP\"}}]},\"peer-groups\":{\"peer-group\":[{\"peer-group-name\":\"IBGP\",\"config\":{\"description\":\"test description\",\"local-as\":65001,\"peer-as\":65001,\"peer-group-name\":\"IBGP\"}}]}}}"}})", "create", true);
        std::cout << "s: " << s << std::endl;
        REQUIRE(s == "");
    }
    catch(YError & e)
    {
        REQUIRE(e.err_msg=="Connect Failed");
    }

}

TEST_CASE("RpcInvalid")
{
    gNMIClient::gNMIClient(shared_ptr<Channel> channel)
    : stub_(gNMI::NewStub(channel)){}
    int ok = 0;

    int result = client.connect();
    REQUIRE(result == ok);

    try
    {
        string reply = client.execute_payload(
             "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
             "<lock>"
             "<source><candidate/></source>"
             "</lock>"
             "</rpc>");

        REQUIRE(NULL != strstr(reply.c_str(), ""));

    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
    }


    REQUIRE(result == ok);
}*/


/*TEST_CASE("gnmi_wrong_json")
{
    gNMIClient client(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));

    int ok = 0;

    int result = client.connect();
    REQUIRE(result == ok);

    try
    {
        string reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:get-config":{"source":{"running":[null]},"filter":"{\"testing}"}})", "read", true);
        std::cout << "reply: " << reply << std::endl;
        //REQUIRE(NULL != strstr(reply.c_str(), "parse error"));
    }
    catch (YCPPError & e)
    {
        std::cout << e.err_msg << std::endl;
        //REQUIRE(e.err_msg=="parse error - unexpected '{'; expected string literal");
    }
    //REQUIRE(result == ok);
}*/

// Disabled as we want to be able to send any RPC via client
//TEST_CASE("CorrectXmlWrongRpc")
//{
//    gNMIClient::gNMIClient(shared_ptr<Channel> channel)
//        : stub_(gNMI::NewStub(channel)){}
//    std::unique_ptr<gNMIClient> client;;
//    int ok = 0;
//
//    int result = client.connect();
//    REQUIRE(result == ok);
//
//    try
//    {
//        string reply = client.execute_payload(
//                "<testing/>"
//        );
//        REQUIRE(reply== "");
//    }
//    catch (YCPPError & e)
//    {
//        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
//    }
//
//
//
//    REQUIRE(result == ok);
//}

/* To Do
TEST_CASE("EmptyRpc")
{
    gNMIClient::gNMIClient(shared_ptr<Channel> channel)
    : stub_(gNMI::NewStub(channel)){}

    std::unique_ptr<gNMIClient> client;

    int ok = 0;

    int result = client.connect();
    REQUIRE(result == ok);

    try
    {
        string reply = client.execute_payload("");
        REQUIRE(reply== "");
    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
    }


    REQUIRE(result == ok);
}*/

TEST_CASE("gnmi_multiple_clients")
{
    gNMIClient client1("127.0.0.1", 50051, "", "");
    gNMIClient client2("127.0.0.1", 50051, "", "");
    gNMIClient client3("127.0.0.1", 50051, "", "");
    gNMIClient client4("127.0.0.1", 50051, "", "");
    gNMIClient client5("127.0.0.1", 50051, "", "");
    gNMIClient client6("127.0.0.1", 50051, "", "");
    gNMIClient client7("127.0.0.1", 50051, "", "");
    gNMIClient client8("127.0.0.1", 50051, "", "");
    gNMIClient client9("127.0.0.1", 50051, "", "");
    gNMIClient client10("127.0.0.1", 50051, "", "");
    gNMIClient client11("127.0.0.1", 50051, "", "");
    gNMIClient client12("127.0.0.1", 50051, "", "");
    gNMIClient client13("127.0.0.1", 50051, "", "");
    gNMIClient client14("127.0.0.1", 50051, "", "");
    gNMIClient client15("127.0.0.1", 50051, "", "");
}

