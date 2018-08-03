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
#include <pybind11/operators.h>
#include <pybind11/pybind11.h>
#include <pybind11/functional.h>
#include <pybind11/stl.h>
#include <pybind11/stl_bind.h>

#include <ydk/gnmi_path_api.hpp>
#include <ydk/gnmi_provider.hpp>
#include <ydk/gnmi_service.hpp>
#include <ydk/types.hpp>

using namespace pybind11;
using namespace std;

using ListCasterBase = detail::list_caster<std::vector<ydk::path::SchemaNode *>, ydk::path::SchemaNode *>;
namespace pybind11{ namespace detail {
template<> struct type_caster<std::vector<ydk::path::SchemaNode *>> : ListCasterBase {
    static handle cast(const std::vector<ydk::path::SchemaNode *> &src, return_value_policy, handle parent) {
        return ListCasterBase::cast(src, return_value_policy::reference, parent);
    }
    static handle cast(const std::vector<ydk::path::SchemaNode *> *src, return_value_policy pol, handle parent) {
        return cast(*src, pol, parent);
    }
};
}}

PYBIND11_MODULE(ydk_gnmi_, ydk_gnmi)
{
    module providers = ydk_gnmi.def_submodule("providers", "providers module");
    module services  = ydk_gnmi.def_submodule("services", "services module");
    module path      = ydk_gnmi.def_submodule("path", "path module");

    class_<ydk::path::gNMISession, ydk::path::Session>(path, "gNMISession")
        .def(init<ydk::path::Repository&, const std::string&, const std::string&, const std::string&, int>(),
             arg("repo"),
             arg("address"),
             arg("username"),
             arg("password"),
             arg("port")=57400)
        .def(init<ydk::path::Repository&, const std::string&, int>(),
             arg("repo"),
             arg("address"),
             arg("port")=57400)
        .def(init<const std::string&, const std::string&, const std::string&, int>(),
             arg("address"),
             arg("username"),
             arg("password"),
             arg("port")=57400)
        .def("get_root_schema", &ydk::path::gNMISession::get_root_schema, return_value_policy::reference)
        .def("invoke", (std::shared_ptr<ydk::path::DataNode> (ydk::path::gNMISession::*)(ydk::path::Rpc&) const)
             &ydk::path::gNMISession::invoke, arg("rpc"), return_value_policy::reference)
        .def("invoke", (void (ydk::path::gNMISession::*)(ydk::path::Rpc& rpc,
                                                         std::function<void(const std::string & response)> out_func,
                                                         std::function<bool(const std::string & response)> poll_func) const)
             &ydk::path::gNMISession::invoke, arg("rpc"),
                                              arg("output_callback_function")=nullptr,
                                              arg("poll_callback_function")=nullptr);

    class_<ydk::gNMIServiceProvider, ydk::ServiceProvider>(providers, "gNMIServiceProvider")
        .def(init<ydk::path::Repository&, const string&, const string&, const string&, int>(),
            arg("repo"), arg("address"), arg("username"), arg("password"), arg("port")=57400)
        .def(init<ydk::path::Repository&, const string&, int>(),
            arg("repo"), arg("address"), arg("port")=57400)
        .def("get_encoding", &ydk::gNMIServiceProvider::get_encoding, return_value_policy::reference)
        .def("get_session", &ydk::gNMIServiceProvider::get_session, return_value_policy::reference)
        .def("get_capabilities", &ydk::gNMIServiceProvider::get_capabilities, return_value_policy::reference);

    class_<ydk::gNMIService>(services, "gNMIService")
	    .def(init<>())
        .def("get", (shared_ptr<ydk::Entity> (ydk::gNMIService::*)
                (ydk::gNMIServiceProvider & provider, ydk::Entity& filter, const string & operation) const)
                &ydk::gNMIService::get, arg("provider"), arg("filter"), arg ("operation"), return_value_policy::reference)
        .def("get", (vector<shared_ptr<ydk::Entity>> (ydk::gNMIService::*)
                (ydk::gNMIServiceProvider & provider, vector<ydk::Entity*> & filter, const string & operation) const)
                &ydk::gNMIService::get, arg("provider"), arg("filter"), arg ("operation"), return_value_policy::reference)
        .def("set", (bool (ydk::gNMIService::*)(ydk::gNMIServiceProvider & provider, ydk::Entity& entity) const)
                &ydk::gNMIService::set, arg("provider"), arg("entity"), return_value_policy::reference)
        .def("set", (bool (ydk::gNMIService::*)(ydk::gNMIServiceProvider & provider, vector<ydk::Entity*> & entity_list) const)
                &ydk::gNMIService::set, arg("provider"), arg("entity"), return_value_policy::reference)

        .def("subscribe", (void (ydk::gNMIService::*)(ydk::gNMIServiceProvider& provider,
                                                      ydk::gNMIService::Subscription* subscription,
                                                      ydk::uint32 qos, const string & mode,
                                                      std::function<void(const string & response)> out_func,
                                                      std::function<bool(const string & response)> poll_func) const)
                &ydk::gNMIService::subscribe, arg("provider"),
                                              arg("subscription"),
                                              arg("qos"),
                                              arg("mode"),
                                              arg("output_callback_function"),
                                              arg("poll_callback_function"))

        .def("subscribe", (void (ydk::gNMIService::*)(ydk::gNMIServiceProvider& provider,
                                                      vector<ydk::gNMIService::Subscription*> & sub_list,
                                                      ydk::uint32 qos, const string & mode,
                                                      std::function<void(const string & response)> out_func,
                                                      std::function<bool(const string & response)> poll_func) const)
               &ydk::gNMIService::subscribe, arg("provider"),
                                             arg("subscription_list"),
                                             arg("qos"),
                                             arg("mode"),
                                             arg("output_callback_function"),
                                             arg("poll_callback_function"));

};

