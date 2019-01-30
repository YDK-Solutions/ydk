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
#include <pwd.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>

#include <ydk/gnmi_provider.hpp>
#include <ydk/gnmi_service.hpp>
#include <ydk/errors.hpp>

#include "../../core/src/catch.hpp"
#include "../../core/tests/config.hpp"

#include <ydk_ydktest/openconfig_interfaces.hpp>

using namespace ydk;
using namespace std;
using namespace ydktest;

TEST_CASE("GNMICreateWithRepo")
{
	ydk::path::Repository repo{TEST_HOME};
	gNMIServiceProvider provider{repo, "127.0.0.1", 50051, "admin", "admin"};

	CHECK_NOTHROW(provider.get_encoding());
}

TEST_CASE("GNMICreateNoRepo")
{
	ydk::path::Repository repo{};
	gNMIServiceProvider provider{repo, "127.0.0.1", 50051, "admin", "admin"};

	CHECK_NOTHROW(provider.get_encoding());
}

static std::string get_temp_model_path()
{
    const char *homeDir = getenv("HOME");
    if (!homeDir) {
        struct passwd* pwd = getpwuid(getuid());
        if (pwd)
           homeDir = pwd->pw_dir;
    }

    std::ostringstream models_path{};
    models_path << homeDir << "/.ydk/temp";
    auto path = models_path.str();

    struct stat st;
    memset(&st, 0, sizeof(struct stat));
    if (stat(path.c_str(), &st) == -1) {
        mkdir(path.c_str(), 0755);
    }
    return path;
}

TEST_CASE("GNMICreateWithNoFiles")
{
	ydk::path::Repository repo{get_temp_model_path()};
	gNMIServiceProvider provider{repo, "127.0.0.1", 50051, "admin", "admin"};
	gNMIService gs{};

	openconfig_interfaces::Interfaces ifcs{};

	CHECK_THROWS_AS(
		gs.get(provider, ifcs, "ALL"),
		ydk::path::YCoreError
	);
}

