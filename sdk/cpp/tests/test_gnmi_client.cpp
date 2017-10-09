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
#define TEST_MODULE gNMIClientTest
#include <string.h>
#include <ydk/gnmi_client.hpp>
#include <ydk/errors.hpp>
#include <iostream>
#include <sys/time.h>
#include "catch.hpp"

using namespace ydk;
using namespace std;

string address = "127.0.0.1:50051";

TEST_CASE("gnmi_xr")
{
    gNMIClient client(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));

    int ok = 0;

    int result = client.connect(address);

    std::cout << "result: " << result << std::endl;
    REQUIRE(result == ok);

    struct timeval t1, t2;
    double elapsedTime;

    // start timer
    gettimeofday(&t1, NULL);

    string reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:get-config":{"source":{"running":[null]},"filter":"{\"openconfig-bgp:bgp\":{}}"}})", "read");
    
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

    int result = client.connect(address);
    REQUIRE(result == ok);
}

TEST_CASE("gnmi_edit_get_config")
{
    gNMIClient client(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));

    int ok = 0;

    int result = client.connect(address);
    REQUIRE(result == ok);

    string reply = client.execute_wrapper(
     R"("rpc":{"ietf-netconf:edit-config":{"target":{"running":[null]},"error-option":"rollback-on-error","config":"{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"172.16.255.2\",\"config\":{\"neighbor-address\":\"172.16.255.2\",\"peer-as\":65172}}]}}}"}})", "create");
    REQUIRE(NULL != strstr(reply.c_str(), "Success"));

    reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:get-config":{"source":{"running":[null]},"filter":"{\"openconfig-bgp:bgp\":{}}"}})", "read");
    
    REQUIRE(NULL != strstr(reply.c_str(), R"({"as":65172})"));

    reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:edit-config":{"target":{"running":[null]},"error-option":"rollback-on-error","config":"{\"openconfig-bgp:bgp\":{}}"}})", "delete");
    REQUIRE(NULL != strstr(reply.c_str(), "Success"));

    REQUIRE(result == ok);
}


TEST_CASE("gnmi_device_not_connected_execute")
{
    gNMIClient client(grpc::CreateChannel("", grpc::InsecureChannelCredentials()));

    try
    {
        string s = client.execute_wrapper( R"("rpc":{"ietf-netconf:edit-config":{"target":{"running":[null]},"error-option":"rollback-on-error","config":"{\"openconfig-bgp:bgp\":{\"global\":{\"config\":{\"as\":65172,\"router-id\":\"1.2.3.4\"}},\"neighbors\":{\"neighbor\":[{\"neighbor-address\":\"6.7.8.9\",\"config\":{\"local-as\":65001,\"neighbor-address\":\"6.7.8.9\",\"peer-as\":65001,\"peer-group\":\"IBGP\"}}]},\"peer-groups\":{\"peer-group\":[{\"peer-group-name\":\"IBGP\",\"config\":{\"description\":\"test description\",\"local-as\":65001,\"peer-as\":65001,\"peer-group-name\":\"IBGP\"}}]}}}"}})", "create");
        std::cout << "s: " << s << std::endl;
        REQUIRE(s == "");
    }
    catch(YCPPError & e)
    {
        REQUIRE(e.err_msg=="Connect Failed");
    }

}

/* To Do
TEST_CASE("RpcInvalid")
{
    gNMIClient::gNMIClient(shared_ptr<Channel> channel)
    : stub_(gNMI::NewStub(channel)){}
    int ok = 0;

    int result = client.connect(address);
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

    int result = client.connect(address);
    REQUIRE(result == ok);

    try
    {
        string reply = client.execute_wrapper(R"("rpc":{"ietf-netconf:get-config":{"source":{"running":[null]},"filter":"{\"testing}"}})", "read");
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
//    int result = client.connect(address);
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

    int result = client.connect(address);
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
    gNMIClient client1(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client2(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client3(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client4(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client5(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client6(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client7(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client8(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client9(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client10(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client11(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client12(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client13(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client14(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));
    gNMIClient client15(grpc::CreateChannel(address, grpc::InsecureChannelCredentials()));

    int result = client1.connect(address) && client2.connect(address) && client3.connect(address) && client4.connect(address) && client5.connect(address)
         && client6.connect(address) && client7.connect(address) && client8.connect(address) && client9.connect(address) && client10.connect(address)
         && client11.connect(address) && client12.connect(address) && client13.connect(address) && client14.connect(address) && client15.connect(address);
    REQUIRE(result == 0);

}

