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

//#include <ydk/path_api.hpp>
//#include <ydk/restconf_client.hpp>
//#include <ydk/codec_provider.hpp>
//#include <ydk/codec_service.hpp>
//#include <ydk/crud_service.hpp>
//#include <ydk/entity_data_node_walker.hpp>
//#include <ydk/entity_util.hpp>
//#include <ydk/executor_service.hpp>
//#include <ydk/filters.hpp>
#include <ydk/gnmi_provider.hpp>
#include <ydk/gnmi_service.hpp>
//#include <ydk/logging_callback.hpp>
//#include <ydk/netconf_provider.hpp>
//#include <ydk/netconf_service.hpp>
//#include <ydk/opendaylight_provider.hpp>
//#include <ydk/restconf_provider.hpp>
#include <ydk/types.hpp>
//#include <ydk/xml_subtree_codec.hpp>

//#include <spdlog/spdlog.h>
//#include <spdlog/sinks/null_sink.h>

using namespace pybind11;
using namespace std;

typedef std::map<std::string, std::shared_ptr<ydk::Entity>> ChildrenMap;
PYBIND11_MAKE_OPAQUE(ChildrenMap)


static object log_debug;
static object log_info;
static object log_warning;
static object log_error;
static object log_critical;
static bool added_nullhandler = false;
static bool enabled_logging = false;


static void add_null_handler(object logger)
{
    if (added_nullhandler) { return; }
    object version = module::import("sys").attr("version_info");
    object ge = version.attr("__ge__");
    // NullHandler is introduced after Python 2.7
    // Add Nullhandler to avoid `handler not found for logger` error for Python > 2.7
    object version_27 = pybind11::make_tuple(2,7);
    bool result = ge(version_27).cast<bool>();
    if (result)
    {
        object null_handler = module::import("logging").attr("NullHandler");
        null_handler = null_handler();
        object add_handler = logger.attr("addHandler");
        add_handler(null_handler);
        added_nullhandler = true;
    }
}

void debug(const char* msg) { log_debug(msg); }
void info(const char* msg) { log_info(msg); }
void warning(const char* msg) { log_warning(msg); }
void error(const char* msg) { log_error(msg); }
void critical(const char* msg) { log_critical(msg); }

void setup_logging()
{
    if (enabled_logging == false)
    {
        object get_logger = module::import("logging").attr("getLogger");
        object logger = get_logger("ydk");

        add_null_handler(logger);
        log_debug = logger.attr("debug");
        log_info = logger.attr("info");
        log_warning = logger.attr("warning");
        log_error = logger.attr("error");
        log_critical = logger.attr("critical");

        ydk::set_logging_callback("debug", debug);
        ydk::set_logging_callback("info", info);
        ydk::set_logging_callback("warning", warning);
        ydk::set_logging_callback("error", error);
        ydk::set_logging_callback("critical", critical);
        enabled_logging = true;
    }
}


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


PYBIND11_MODULE(ydk_, ydk)
{
    module providers = ydk.def_submodule("providers", "providers module");
    module services = ydk.def_submodule("services", "services module");
    module filters = ydk.def_submodule("filters", "filters module");
    module types = ydk.def_submodule("types", "types module");
    module path = ydk.def_submodule("path", "path module");
    module entity_utils = ydk.def_submodule("entity_utils", "entity utils module");
    module logging = ydk.def_submodule("logging", "logging");
    module clients = ydk.def_submodule("clients", "clients");

    bind_map<ChildrenMap>(types, "ChildrenMap");

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
        .def("invoke", &ydk::path::gNMISession::invoke, return_value_policy::reference)
        .def("invoke", &ydk::path::gNMISession::invoke, return_value_policy::void);

    class_<ydk::gNMIServiceProvider, ydk::ServiceProvider>(providers, "gNMIServiceProvider")
        .def(init<ydk::path::Repository&, const string&, const string&, const string&, int>(),
            arg("repo"), arg("address"), arg("username"), arg("password"), arg("port")=57400)
        .def("get_encoding", &ydk::gNMIServiceProvider::get_encoding, return_value_policy::reference)
        .def("get_session", &ydk::gNMIServiceProvider::get_session, return_value_policy::reference)
        .def("get_capabilities", &ydk::gNMIServiceProvider::get_capabilities, return_value_policy::reference);

    class_<ydk::gNMIService>(services, "gNMIService")
	    .def(init<>())
	    .def("get", &ydk::gNMIService::get, arg("provider"), arg("filter"), return_value_policy::reference)
    	.def("set", &ydk::gNMIService::set, arg("provider"), arg("entity"), arg("operation"), return_value_policy::reference)
    	.def("subscribe", &ydk::gNMIService::subscribe, arg("provider"),
    	                                                arg("filter"),
    	                                                arg("list_mode"),
    	                                                arg("qos"),
    	                                                arg("mode"),
    	                                                arg("sample_interval"),
    	                                                arg("callback_function"));

    setup_logging();

};

