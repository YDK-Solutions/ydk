#  ----------------------------------------------------------------
# Copyright 2016 Cisco Systems
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------------

'''
   YDK PY converter
   
'''


import api_model
import common
import gen_target
import optparse
import os
import shutil

from pyang import error
from pyang import plugin
from module_printer import PythonModulePrinter
from helper import get_rst_file_name
from api_model import Class, Enum


class _EmitArgs:

    def __init__(self, ctx, package, extra_args=None):
        self.ctx = ctx
        self.package = package
        self.extra_args = extra_args


# ============================================================
# Plugin class definitions
# ============================================================
class YDKPythonGen():

    def __init__(self, ydk_dir, ydk_doc_dir, groupings_as_classes):
        self.ydk_dir = ydk_dir
        self.groupings_as_classes = groupings_as_classes
        self.ydk_doc_dir = ydk_doc_dir

    def emit(self, modules):
        self.ypy_ctx = None
        self.packages = []
        self.models_dir = ''
        self.test_dir = ''

        self.parse_modules(modules)
        self.initialize_print_environment()
        self.print_files()

    def initialize_print_environment(self):
        self.initialize_top_level_directories()
        self.initialize_printer_context()

    def parse_modules(self, modules):
        if not self.groupings_as_classes:
            self.packages = api_model.generate_expanded_api_model(modules)
        else:
            self.packages = api_model.generate_grouping_class_api_model(
                modules)
        self.packages = sorted(self.packages, key=lambda package: package.name)

    def initialize_top_level_directories(self):
        self.models_dir = self.initialize_output_directory(
            self.ydk_dir + '/models', True)
        self.test_dir = self.initialize_output_directory(
            self.ydk_dir + '/tests', True)

    def initialize_printer_context(self):
        self.ypy_ctx = common.PrintCtx()
        self.ypy_ctx.meta = True
        self.ypy_ctx.tab_size = 4
        self.ypy_ctx.printer = PythonModulePrinter(self.ypy_ctx)

    def print_files(self):
        self.print_modules()
        self.print_init_file(self.models_dir)
        self.print_yang_ns_file()
        self.print_import_tests_file()

        # RST Documentation
        if self.ydk_doc_dir is not None:
            self.print_python_rst_ydk_models()

    def print_modules(self):
        only_modules = [package.stmt for package in self.packages]
        size = len(only_modules)

        for index, package in enumerate(self.packages):
            self.print_module(index, package, size)

    def print_module(self, index, package, size):
        print 'Processing %d of %d %s' % (index + 1, size, package.stmt.pos.ref)

        py_mod_name = package.get_py_mod_name()
        sub = py_mod_name[len('ydk.models.'): py_mod_name.rfind('.')]

        module_dir = self.initialize_output_directory(
            '%s/%s' % (self.models_dir, sub))
        meta_dir = self.initialize_output_directory(module_dir + '/_meta')
        test_output_dir = self.initialize_output_directory(
            '%s/%s' % (self.test_dir, sub))

        # RST Documentation
        self.print_python_module(package, index, module_dir, size, sub)
        self.print_meta_module(package, meta_dir)
        self.print_test_module(package, test_output_dir)
        if self.ydk_doc_dir is not None:
            self.print_python_rst_module(package)

    def print_python_rst_module(self, package):
        if self.ydk_doc_dir is None:
            return
        # Skip generating documentation for empty modules
        if len(package.owned_elements) == 0:
            return

        def _walk_n_print(named_element, p):
            self.print_file(get_python_rst_file_name(p, named_element),
                            gen_target.emit_python_rst,
                            _EmitArgs(self.ypy_ctx, named_element))

            for owned_element in named_element.owned_elements:
                if isinstance(owned_element, Class) or isinstance(owned_element, Enum):
                    _walk_n_print(owned_element, p)

        _walk_n_print(package, self.ydk_doc_dir)

    def print_python_rst_ydk_models(self):
        if self.ydk_doc_dir is None:
            return
        packages = [p for p in self.packages if len(p.owned_elements) > 0]

        self.print_file(get_python_rst_ydk_models_file_name(self.ydk_doc_dir),
                        gen_target.emit_ydk_models_rst,
                        _EmitArgs(self.ypy_ctx, packages))

    def print_python_module(self, package, index, path, size, sub):
        self.print_init_file(path)

        package.parent_pkg_name = sub
        self.print_file(get_python_module_file_name(path, package),
                        gen_target.emit_module,
                        _EmitArgs(self.ypy_ctx, package, (sub, package.name)))

    def print_meta_module(self, package, path):
        self.print_init_file(path)
        self.print_file(get_meta_module_file_name(path, package),
                        gen_target.emit_meta,
                        _EmitArgs(self.ypy_ctx, package))

    def print_test_module(self, package, path):
        self.print_file(get_test_module_file_name(path, package),
                        gen_target.emit_test_module,
                        _EmitArgs(self.ypy_ctx, package))

    def print_yang_ns_file(self):
        self.print_file(get_yang_ns_file_name(self.models_dir),
                        gen_target.emit_yang_ns,
                        _EmitArgs(self.ypy_ctx, self.packages))

    def print_import_tests_file(self):
        self.print_file(get_import_test_file_name(self.test_dir),
                        gen_target.emit_importests,
                        _EmitArgs(self.ypy_ctx, self.packages))

    def print_init_file(self, path):
        file_name = get_init_file_name(path)
        if not os.path.isfile(file_name):
            self.print_file(file_name)

    def print_file(self, path, emit_func=None, emit_args=None):
        with open(path, 'w+') as file_descriptor:
            if emit_func is not None and emit_args is not None:
                emit_args.ctx.fd = file_descriptor
                if emit_args.extra_args is None:
                    emit_func(emit_args.ctx, emit_args.package)
                else:
                    emit_func(
                        emit_args.ctx, emit_args.package, emit_args.extra_args)

    def initialize_output_directory(self, path, delete_if_exists=False):
        if delete_if_exists:
            if os.path.exists(path):
                shutil.rmtree(path)
        if not os.path.isdir(path):
            os.mkdir(path)
        return path


def get_init_file_name(path):
    return path + '/__init__.py'


def get_yang_ns_file_name(path):
    return path + '/_yang_ns.py'


def get_import_test_file_name(path):
    return path + '/import_tests.py'


def get_python_rst_file_name(path, named_element):
    return '%s/%s.rst' % (path, get_rst_file_name(named_element))


def get_python_rst_ydk_models_file_name(path):
    return '%s/ydk.models.rst' % path


def get_python_module_file_name(path, package):
    return '%s/%s.py' % (path, package.name)


def get_meta_module_file_name(path, package):
    return '%s/_%s.py' % (path, package.name)


def get_test_module_file_name(path, package):
    return '%s/%sTest.py' % (path, package.stmt.arg.replace('-', '_'))
