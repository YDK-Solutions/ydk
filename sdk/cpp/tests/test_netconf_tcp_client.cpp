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

/*
    This is a test file target XR 6.3.1 release, modified from
    test_netconf_client.cpp It test against XR image loaded in Vagrant with
    NETCONF TCP port forwarded to 41631.
*/

#define TEST_MODULE NetconfTCPClientTest
#include <string.h>
#include "../core/src/netconf_tcp_client.hpp"
#include "../core/src/errors.hpp"
#include <iostream>
#include <sys/time.h>
#include "catch.hpp"

using namespace ydk;
using namespace std;
#define NC_VERB_VERBOSE 2

TEST_CASE("tcp_xr")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

    struct timeval t1, t2;
    double elapsedTime;

    // start timer
    gettimeofday(&t1, NULL);

    string reply = client.execute_payload(
     "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
     "<get-config>"
     "<source><running/></source>"
     "</get-config>"
     "</rpc>");


    // stop timer
    gettimeofday(&t2, NULL);

    // compute and print the elapsed time in millisec
    elapsedTime = (t2.tv_sec - t1.tv_sec) * 1000.0;      // sec to ms
    elapsedTime += (t2.tv_usec - t1.tv_usec) / 1000.0;   // us to ms
    cout << elapsedTime << " ms.\n";

    REQUIRE(result == OK);
}


TEST_CASE("tcp_create")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);


    REQUIRE(result == OK);
}

TEST_CASE("tcp_edit_get_config")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

    string reply = client.execute_payload(
     "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
     "<edit-config>"
     "<target><candidate/></target>"
     "<config>"
     "<bgp xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XR-ipv4-bgp-cfg\">"
     "<instance>"
     "<instance-name>default</instance-name>"
     "<instance-as>"
     "<as>0</as>"
     "<four-byte-as>"
     "<as>65001</as>"
     "<bgp-running></bgp-running>"
     "<default-vrf>"
     "<bgp-entity>"
     "<neighbor-groups>"
     "<neighbor-group>"
     "<neighbor-group-name>IBGP</neighbor-group-name>"
     "<create></create>"
     "<neighbor-group-afs>"
     "<neighbor-group-af>"
     "<af-name>ipv4-unicast</af-name>"
     "<activate></activate>"
     "</neighbor-group-af>"
     "</neighbor-group-afs>"
     "<remote-as>"
     "<as-xx>0</as-xx>"
     "<as-yy>65001</as-yy>"
     "</remote-as>"
     "<update-source-interface>Loopback0</update-source-interface>"
     "</neighbor-group>"
     "</neighbor-groups>"
     "<neighbors>"
     "<neighbor>"
     "<neighbor-address>172.16.255.2</neighbor-address>"
     "<neighbor-group-add-member>IBGP</neighbor-group-add-member>"
     "</neighbor>"
     "</neighbors>"
     "</bgp-entity>"
     "<global>"
     "<global-afs>"
     "<global-af>"
     "<af-name>ipv4-unicast</af-name>"
     "<enable></enable>"
     "</global-af>"
     "</global-afs>"
     "</global>"
     "</default-vrf>"
     "</four-byte-as>"
     "</instance-as>"
     "</instance>"
     "</bgp>"
     "</config>"
     "</edit-config>"
     "</rpc>");
    REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));

    // ?
    // reply = client.execute_payload(
    //  "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
    //  "<commit/>"
    //  "</rpc>"
    // );
    // REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));

    // ?
    // reply = client.execute_payload(
    //  "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
    //  "<get-config>"
    //  "<source><running/></source>"
    //  "<filter>"
    //  "<bgp xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XR-ipv4-bgp-cfg\"/>"
    //  "</filter>"
    //  "</get-config>"
    //  "</rpc>");
    // REQUIRE(NULL != strstr(reply.c_str(), "<neighbor-address>172.16.255.2</neighbor-address>"));

    reply = client.execute_payload(
     "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
     "<edit-config>"
     "<target><candidate/></target>"
     "<config xmlns:xc=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
     "<bgp xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XR-ipv4-bgp-cfg\" xc:operation=\"delete\"/>"
     "</config>"
     "</edit-config>"
     "</rpc>");
    REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));

    // ?
    // reply = client.execute_payload(
    //  "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
    //  "<commit/>"
    //  "</rpc>"
    // );
    // REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));

    REQUIRE(result == OK);
}

TEST_CASE("tcp_validate")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

    string reply = client.execute_payload(
     "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
     "<validate>"
     "<source><candidate/></source>"
     "</validate>"
     "</rpc>");

    REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));


    REQUIRE(result == OK);
}

TEST_CASE("tcp_lock_unlock")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

    string reply = client.execute_payload(
                                          "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
                                          "<discard-changes/>"
                                          "</rpc>");

    REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));


    reply = client.execute_payload(
         "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
         "<lock>"
         "<target><candidate/></target>"
         "</lock>"
         "</rpc>");

    REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));

    reply = client.execute_payload(
         "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
         "<unlock>"
         "<target><candidate/></target>"
         "</unlock>"
         "</rpc>");

    REQUIRE(NULL != strstr(reply.c_str(), "<ok/>"));

    try
    {
        // throw exception for every rpc error
        reply = client.execute_payload(
             "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
             "<unlock>"
             "<target><candidate/></target>"
             "</unlock>"
             "</rpc>");
    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
    }

    REQUIRE(result == OK);
}

TEST_CASE("tcp_rpc_error")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);
    try
    {
        string reply = client.execute_payload(
         "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
         "<edit-config>"
         "<target><candidate/></target>"
         "<config>"
         "<bgp xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XR-ipv4-bgp-cfg\">"
         "<instance>"
         "<instance-name>default</instance-name>"
         "<instance-as>"
         "<as>0</as>"
         "<four-byte-as>"
         "<as>65001</as>"
         "<bgp-running></bgp-running>"
         "<default-vrf>"
         "<bgp-entity>"
         "<neighbor-groups>"
         "<neighbor-group>"
         "<neighbor-group-name>IBGP</neighbor-group-name>"
         "<create></create>"
         "<neighbor-group-afs>"
         "<neighbor-group-af>"
         "<af-name>ipv4-unicast</af-name>"
         "<activate></activate>"
         "</neighbor-group-af>"
         "</neighbor-group-afs>"
         "<remote-as>"
         "<as-xx>0</as-xx>"
         "<as-yy>WRONG STRING</as-yy>"
         "</remote-as>"
         "<update-source-interface>Loopback0</update-source-interface>"
         "</neighbor-group>"
         "</neighbor-groups>"
         "<neighbors>"
         "<neighbor>"
         "<neighbor-address>172.16.255.2</neighbor-address>"
         "<neighbor-group-add-member>IBGP</neighbor-group-add-member>"
         "</neighbor>"
         "</neighbors>"
         "</bgp-entity>"
         "<global>"
         "<global-afs>"
         "<global-af>"
         "<af-name>ipv4-unicast</af-name>"
         "<enable></enable>"
         "</global-af>"
         "</global-afs>"
         "</global>"
         "</default-vrf>"
         "</four-byte-as>"
         "</instance-as>"
         "</instance>"
         "</bgp>"
         "</config>"
         "</edit-config>"
         "</rpc>");
        REQUIRE(NULL != strstr(reply.c_str(), "<rpc-error>"));
    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
    }

}

TEST_CASE("tcp_device_not_connected_execute")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    try
    {
        string s = client.execute_payload(
                 "<rpc xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">"
                 "<edit-config>"
                 "<target><candidate/></target>"
                 "<config>"
                 "<runner xmlns=\"http://cisco.com/ns/yang/ydktest-sanity\"><ytypes><built-in-t><number8>aaa</number8></built-in-t></ytypes></runner>"
                 "</config>"
                 "</edit-config>"
                 "</rpc>");
        REQUIRE(s== "");
    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not execute payload. Not connected to 127.0.0.1");
    }

}


TEST_CASE("tcp_rpc_invalid")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

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


    REQUIRE(result == OK);
}

TEST_CASE("tcp_wrong_xml")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

    try
    {
        string reply = client.execute_payload(
                "<testing>"
         );
        REQUIRE(reply== "");
    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
    }


    REQUIRE(result == OK);
}

TEST_CASE("tcp_correct_xml_wrong_rpc")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

    try
    {
        string reply = client.execute_payload(
                "<testing/>"
        );
        REQUIRE(reply== "");
    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
    }



    REQUIRE(result == OK);
}

TEST_CASE("tcp_empty_rpc")
{
    NetconfTCPClient client{"vagrant", "vagrant", "127.0.0.1", 41631};
    int OK = 0;

    int result = client.connect();
    REQUIRE(result == OK);

    try
    {
        string reply = client.execute_payload("");
        REQUIRE(reply== "");
    }
    catch (YCPPError & e)
    {
        REQUIRE(e.err_msg=="YCPPClientError: Could not build payload");
    }


    REQUIRE(result == OK);
}

// Timeout error
// TEST_CASE("tcp_MultipleClients")
// {
//     NetconfTCPClient client1{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client2{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client3{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client4{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client5{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client6{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client7{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client8{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client9{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client10{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client11{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client12{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client13{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client14{"vagrant", "vagrant", "127.0.0.1", 41631};
//     NetconfTCPClient client15{"vagrant", "vagrant", "127.0.0.1", 41631};

//     int result = client1.connect() && client2.connect() && client3.connect() && client4.connect() && client5.connect()
//          && client6.connect() && client7.connect() && client8.connect() && client9.connect() && client10.connect()
//          && client11.connect() && client12.connect() && client13.connect() && client14.connect() && client15.connect();
//     REQUIRE(result == 0);

// }

