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

#ifndef _RESTCONF_SESSION_H_
#define _RESTCONF_SESSION_H_

#include <map>
#include <memory>
#include <string>

#include "../path_api.hpp"

namespace ydk {
static const std::string default_capabilities_url = "/ietf-restconf-monitoring:restconf-state/capabilities";

class CapabilitiesParser;
struct OpenDaylightNode;
class RestconfClient;

class RestconfSession : public path::Session {
public:
    RestconfSession(
        path::Repository & repo,
        const std::string & address,
        const std::string & username,
        const std::string & password,
        int port = 80,
        EncodingFormat encoding = EncodingFormat::JSON,
        const std::string & state_url_root = "/data"
    );

    RestconfSession(
        std::unique_ptr<RestconfClient> client,
        std::shared_ptr<ydk::path::RootSchemaNode> root_schema,
        const std::string & edit_method,
        const std::string & state_url_root,
        EncodingFormat encoding
    );

    ~RestconfSession();

    path::RootSchemaNode& get_root_schema() const;
    std::shared_ptr<path::DataNode> invoke(path::Rpc& rpc) const;

private:
    void initialize(path::Repository & repo);
    std::shared_ptr<path::DataNode> handle_edit(path::Rpc& rpc, const std::string & yfilter) const;
    std::shared_ptr<path::DataNode> handle_read(path::Rpc& rpc) const;
private:
        std::unique_ptr<RestconfClient> client;
        std::shared_ptr<ydk::path::RootSchemaNode> root_schema;
        std::vector<std::string> server_capabilities;

        EncodingFormat encoding;
        std::string edit_method;
        std::string config_url_root;
        std::string state_url_root;
};
}

#endif /*_RESTCONF_SESSION_H_*/
