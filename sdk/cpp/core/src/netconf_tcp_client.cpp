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

#include <algorithm>
#include <cstdio>
#include <map>
#include <regex>
#include <string.h>
#include <sstream>
#include <string>
#include <unistd.h>
#include <vector>

#include <curl/curl.h>
#include <libxml/parser.h>
#include <libxml/tree.h>

#include "capabilities_parser.hpp"
#include "errors.hpp"
#include "ietf_parser.hpp"
#include "logger.hpp"
#include "netconf_tcp_client.hpp"


namespace ydk
{
static const long TIMEOUT = 600L;
static const int EIGHT_K = 8196;

static const std::string RPC_ERROR_PATH("/rpc-reply/rpc-error");

static const char FRAMING_11[] = "]]>]]>";
static const char FRAMING_10[] = "\n##\n";
static const char LF_HASH[] = "\n#";
static const char LF[] = "\n";

static const char *HELLO_11 = "\n<hello xmlns=\"urn:ietf:params:xml:ns:netconf:base:1.0\">\n"
                                   "<capabilities>\n"
                                   "<capability>urn:ietf:params:netconf:base:1.1</capability>\n"
                                   "</capabilities>\n"
                                   "</hello>\n"
                                   "]]>]]>";

static const std::regex QUOT_P("&quot;");
static const std::regex AMP_P("&amp");
static const std::regex APOS_P("&apos");
static const std::regex LT_P("&lt");
static const std::regex GT_P("&gt");
static const std::regex CHUNK_START_P("\n#[0-9]*\n");
static const std::regex FRAMING_10_P(FRAMING_10);

static const std::string QUOT("\"");
static const std::string AMP("&");
static const std::string APOS("'");
static const std::string LT("<");
static const std::string GT(">");
static const std::string EMPTY("");

static const std::map<const std::string, const std::regex> PATTERN_MAP{{QUOT, QUOT_P},
                                                                       {AMP, AMP_P},
                                                                       {APOS, APOS_P},
                                                                       {LT, LT_P},
                                                                       {GT, GT_P},
                                                                       {EMPTY, CHUNK_START_P}};

static std::string trim_reply(std::string& str);
static int wait_on_socket(curl_socket_t sockfd, int for_recv, long timeout_ms);
static bool ends_with_framing(const char* buf, size_t nread, int version);
static xmlDocPtr get_xml_doc(const std::string &payload);
static bool check_xml_doc(xmlNodePtr root, std::string cur_path, int max_lvl);


NetconfTCPClient::NetconfTCPClient(std::string username, std::string password,
                                   std::string address, int port)
    : NetconfClientBase(),
      username(username), hostname(address), password(password), port(port), msgid(0)
{
    initialize(address, port);
    YLOG_INFO("Ready to communicate with {} via TCP", address);
}

NetconfTCPClient::~NetconfTCPClient()
{
    curl_easy_cleanup(curl);
    curl_global_cleanup();
}

void NetconfTCPClient::initialize(std::string address, int port)
{
    initialize_curl(address, port);
}

void NetconfTCPClient::initialize_curl(std::string address, int port)
{
    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (!curl)
    {
        throw(YCPPClientError{"Unable to create curl environment."});
    }

    // curl_easy_setopt(curl, CURLOPT_VERBOSE, 1L); // uncomment for debugging
    curl_easy_setopt(curl, CURLOPT_URL, address.c_str());
    curl_easy_setopt(curl, CURLOPT_PORT, port);
    curl_easy_setopt(curl, CURLOPT_CONNECT_ONLY, 1L);

    CURLcode res;
    res = curl_easy_perform(curl);
    check_ok(res, "Connection failed: {}");

    // get socket descriptor
    res = curl_easy_getinfo(curl, CURLINFO_ACTIVESOCKET, &sockfd);
    check_ok(res, "Unable to retrieve sockfd: {}");
}

int NetconfTCPClient::connect()
{

    char buf[EIGHT_K];
    size_t nread = 0;
    std::stringstream hello_ss;
    CURLcode res;

    for(;;)
    {
        bzero(&buf, sizeof(buf));
        nread = 0;
        res = curl_easy_recv(curl, buf, sizeof(buf)-1, &nread);

        check_timeout(res, 1, "Handshake failed, TCP client connect session timeout: {}");

        if (strncmp(buf, "Username: ", 10) == 0 || strncmp(buf, "Password: ", 10) == 0)
        {
            std::string value = strncmp(buf, "Username: ", 10) ? username : password;
            size_t value_len = value.size();
            send_value(value.c_str(), value_len);
        }
        else if(nread)
        {
            buf[nread] = '\0';
            hello_ss << buf;

            if (ends_with_framing(buf, nread, 11))
            {
                // required, send hello request back to device
                send_value(HELLO_11, strlen(HELLO_11));
                auto found = hello_ss.str().find(FRAMING_11);
                hello_msg = hello_ss.str().substr(0, found);
                YLOG_DEBUG("Received hello message from device:\n{}", hello_msg);
                break;
            }
        }
    }
    YLOG_INFO("Finished TCP connection.");
    init_capabilities();

    connected = true;

    return EXIT_SUCCESS;
}

void NetconfTCPClient::init_capabilities()
{
    IetfCapabilitiesXmlParser xml_parser{};
    IetfCapabilitiesParser capabilities_parser{};
    server_capabilities = xml_parser.parse(hello_msg);
}

std::string NetconfTCPClient::execute_payload(const std::string & payload)
{
    if(!connected)
    {
        auto err_msg = "Could not execute payload. Not connected to " + hostname;
        throw(YCPPClientError{err_msg});
    }
    send(payload);
    return recv();
}

std::string NetconfTCPClient::recv()
{
    // TODO: check and parse every reply?
    // TODO: raise exception for every rpc-error?
    auto reply = recv_value();

    if (reply.empty())
    {
        return reply;
    }

    auto rpc_reply = trim_reply(reply);
    if(!check_rpc_reply(rpc_reply))
    {
        return rpc_reply;
    }
    else
    {
        YLOG_ERROR("Could not build payload");
        throw(YCPPClientError{"Could not build payload"});
    }

}

bool NetconfTCPClient::check_rpc_reply(std::string &reply)
{
    bool ret = false;
    auto doc = get_xml_doc(reply);
    auto root = xmlDocGetRootElement(doc);

    // check doc
    std::string root_path("");
    ret = check_xml_doc(root, root_path, 2);
    xmlFreeDoc(doc);
    return ret;
}

static bool check_xml_doc(xmlNodePtr root, std::string cur_path, int max_lvl)
{

    bool ret = false;
    if (max_lvl <= 0)
    {
        return ret;
    }

    int i;
    std::string path;
    xmlNodePtr node = NULL;

    for(node = root, i = 0; node; node = node->next, i = i + 1)
    {
        path = cur_path;
        path += "/";
        path += reinterpret_cast<const char*>(node->name);
        if (path == RPC_ERROR_PATH)
        {
            return true;
        }
        ret |= check_xml_doc(node->children, path, max_lvl-1);
    }
    return ret;
}

static xmlDocPtr get_xml_doc(const std::string &payload)
{
    xmlDocPtr doc;
    YLOG_INFO("get_xml_doc: payload = {}", payload);
    doc = xmlReadMemory(payload.c_str(), strlen(payload.c_str()), "noname.xml", NULL, 0);
    if (doc == NULL) {
        YLOG_ERROR("Could not build payload");
        throw(YCPPClientError{"Could not build payload"});
    }
    return doc;
}

static void xml_to_string(xmlDocPtr doc, xmlNodePtr root, std::string &out)
{
    xmlBufferPtr buf = xmlBufferCreate();
    if (buf != NULL)
    {
        xmlNodeDump(buf, doc, root, 0, 1);
        std::string tmp{reinterpret_cast<char*>(buf->content)};
        out = tmp;
        xmlBufferFree(buf);
    }
}

std::string NetconfTCPClient::add_message_id(const std::string &payload)
{
    msgid++;
    char msg_id_str[16];
    sprintf(msg_id_str, "%llu", msgid);

    auto doc = get_xml_doc(payload);
    auto cur = xmlDocGetRootElement(doc);
    xmlNewProp(cur, (const xmlChar *) "message-id", (const xmlChar *) msg_id_str);

    std::string new_payload;
    xml_to_string(doc, cur, new_payload);
    xmlFreeDoc(doc);

    return new_payload;
}

void NetconfTCPClient::send(const std::string &payload)
{
    // add message id property to payload
    auto new_payload = add_message_id(payload);
    YLOG_DEBUG("TCP client send payload:\n{}", new_payload);
    // add chunk size
    std::ostringstream ss;
    ss << LF_HASH << new_payload.size() << LF << new_payload << FRAMING_10;
    send_value(ss.str().c_str(), ss.str().size());
}

std::vector<std::string> NetconfTCPClient::get_capabilities()
{
    return server_capabilities;
}

std::string NetconfTCPClient::get_hostname_port()
{
    std::ostringstream os;
    os << hostname << ":" << port;
    return os.str();
}

void NetconfTCPClient::send_value(const char* value, size_t value_len)
{
    CURLcode res;
    size_t nsent = 0;
    size_t nsent_total = 0;
    do
    {
        nsent = 0;
        do
        {
            res = curl_easy_send(curl, value + nsent_total, value_len - nsent_total, &nsent);
            nsent_total += nsent;
            check_timeout(res, 0, "TCP client send_value session timeout: {}");
        } while (res == CURLE_AGAIN);

        check_ok(res, "TCP client error: {}");
        YLOG_DEBUG("Sent {} bytes.\n", (curl_off_t)nsent);

    } while (nsent_total < value_len &&wait_on_socket(sockfd, 0, TIMEOUT));
    YLOG_DEBUG("TCP client sent total {} bytes:\n{}", nsent_total, value);
}

std::string NetconfTCPClient::recv_value()
{
    // wait on socket
    unsigned int usecs = 100000;
    usleep(usecs);

    std::stringstream ss;
    char buf[EIGHT_K];
    size_t nread;
    CURLcode res;
    bool found;

    for(;;)
    {
        do
        {
            found = false;
            bzero(&buf, sizeof(buf));
            nread = 0;
            res = curl_easy_recv(curl, buf, sizeof(buf)-1, &nread);
            buf[nread] = '\0';
            ss << buf;
            // curl_easy_recv is not the basic socket recv block,
            // buf might contain multiple messages, for example
            // a <ok> message and a <rpc-reply>
            found = ends_with_framing(buf, nread, 10);

            if(found)
            {
                break;
            }

        } while(!wait_on_socket(sockfd, 1, TIMEOUT));

        if (found || !nread)
        {
            break;
        }

        check_ok(res, "TCP client error: {}");
    } // timeout?

    YLOG_DEBUG("TCP client received {} bytes:\n{}", ss.str().length(), ss.str());
    return ss.str();
}

void NetconfTCPClient::check_ok(CURLcode res, const char* fmt)
{
    if (res != CURLE_OK)
    {
        YLOG_ERROR(fmt, curl_easy_strerror(res));
        throw(YCPPClientError{curl_easy_strerror(res)});
    }
}

void NetconfTCPClient::check_timeout(CURLcode res, int for_recv, const char* fmt)
{
    if (res == CURLE_AGAIN && !wait_on_socket(sockfd, for_recv, TIMEOUT))
    {
        YLOG_ERROR(fmt, curl_easy_strerror(res));
        throw(YCPPClientError(curl_easy_strerror(res)));
    }
}

// from libcurl examples
/* Auxiliary function that waits on the socket. */
static int wait_on_socket(curl_socket_t sockfd, int for_recv, long timeout_ms)
{
    struct timeval tv;
    fd_set infd, outfd, errfd;
    int res;

    tv.tv_sec = timeout_ms / 1000;
    tv.tv_usec= (timeout_ms % 1000) * 1000;

    FD_ZERO(&infd);
    FD_ZERO(&outfd);
    FD_ZERO(&errfd);

    FD_SET(sockfd, &errfd); /* always check for error */

    if(for_recv) {
        FD_SET(sockfd, &infd);
    }
    else {
        FD_SET(sockfd, &outfd);
    }
    /* select() returns the number of signalled sockets or -1 */
    res = select((int)sockfd + 1, &infd, &outfd, &errfd, &tv);
    return res;
}

static std::string trim_reply(std::string& str)
{
    // drop <ok> message if necessary
    int cnt = 0;
    std::size_t pos = str.find(FRAMING_10, 0);
    while (pos != std::string::npos)
    {
        cnt ++;
        pos = str.find(FRAMING_10, pos+1);
    }

    if(cnt == 2)
    {
        YLOG_DEBUG("TCP client dropped <ok> message, raw msg:\n {}", str);
        str = str.substr(str.find(FRAMING_10, 0));
    }

    for (auto val: PATTERN_MAP)
    {
        str = std::regex_replace(str, val.second, val.first);
    }
    str = std::regex_replace(str, FRAMING_10_P, EMPTY);

    YLOG_DEBUG("TCP client received reply:\n{}", str);
    return str;
}

static bool ends_with_framing(const char* buf, size_t nread, int version)
{
    const char* last_several = NULL;
    const char* framing = NULL;
    if (version == 10)
    {
        last_several = &buf[nread-4];
        framing = FRAMING_10;
    }
    else if(version == 11)
    {
        last_several = &buf[nread-6];
        framing = FRAMING_11;
    }
    return strcmp(last_several, framing) == 0;
}

}
