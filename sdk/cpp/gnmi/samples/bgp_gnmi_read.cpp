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
#include <iostream>
#include <memory>
#include "ydk/path_api.hpp"
#include "ydk/gnmi_provider.hpp"
#include "args_parser.h"
#include "ydk/crud_service.hpp"
#include "ydk_openconfig/openconfig_bgp.hpp"

#include <spdlog/spdlog.h>
#include <fstream>

using namespace std;
using namespace ydk;
using namespace path;

void print_paths(ydk::path::SchemaNode & sn)
{
    std::cout << sn.get_path() << std::endl;
    for(auto const& p : sn.get_children()) print_paths(*p);
}

int main(int argc, char* argv[]) 
{

    vector<string> args = parse_args(argc, argv);
    if(args.empty()) return 1;
    
    string host, username, password, port, address;
    username = args[0]; password = args[1]; host = args[2]; port = args[3];

    address.append(host);
    address.append(":");
    address.append(port);
    bool verbose = (args[4]=="--verbose");
    std::cout << "verbose: " << verbose << std::endl;
    if(true)
    {
        auto logger = spdlog::stdout_color_mt("ydk");
        logger->set_level(spdlog::level::debug);
    }

    ydk::path::Repository repo{"/usr/local/share/ydk/0.0.0.0\:50051/"};

    gNMIServiceProvider provider{repo, address};

    CrudService crud{};

    auto bgp_filter = make_unique<openconfig::openconfig_bgp::Bgp>();
    crud.read_config(provider, *bgp_filter);

    auto bgp_read = crud.read_config(provider, *bgp_filter);
    if(bgp_read == nullptr)
    {
        cout << "=================================================="<<endl;
        cout << "No entries found"<<endl<<endl;
        cout << "=================================================="<<endl;
        return 0;
    }
    openconfig::openconfig_bgp::Bgp * bgp_read_ptr = dynamic_cast<openconfig::openconfig_bgp::Bgp*>(bgp_read.get());

    cout << "=================================================="<<endl;
    cout << "BGP configuration: " << endl<<endl;
    cout << "AS: " << bgp_read_ptr->global->config->as << endl;
    cout << "Router ID: " << bgp_read_ptr->global->config->router_id << endl<<endl;

    for(size_t index=0; index < bgp_read_ptr->neighbors->neighbor.size(); index++)
    {
        openconfig::openconfig_bgp::Bgp::Neighbors::Neighbor & neighbor = *(bgp_read_ptr->neighbors->neighbor[index]);

        cout << "Neighbor address: " << neighbor.neighbor_address <<endl;
        cout << "Neighbor local AS: " <<  neighbor.config->local_as << endl;
        cout << "Neighbor peer group: " <<  neighbor.config->peer_group << endl;
        cout << "Neighbor peer type: " <<  neighbor.config->peer_type << endl<<endl;
    }

    for(size_t index=0; index < bgp_read_ptr->global->afi_safis->afi_safi.size(); index++)
    {
        openconfig::openconfig_bgp::Bgp::Global::AfiSafis::AfiSafi & afi_safi = *(bgp_read_ptr->global->afi_safis->afi_safi[index]);
        cout << "AFI-SAFI name: " << afi_safi.afi_safi_name <<endl;
        cout << "AFI-SAFI config name: " <<  afi_safi.config->afi_safi_name <<endl;
        cout << "AFI-SAFI enabled: " <<  afi_safi.config->enabled <<endl<<endl;
    }

    cout << "=================================================="<<endl<<endl;

    return 0;
}
