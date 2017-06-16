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
   YDK GO converter

'''

from __future__ import print_function

import os

from ydkgen.api_model import Class, Enum
from ydkgen.common import get_rst_file_name
from ydkgen.printer.language_bindings_printer import LanguageBindingsPrinter, _EmitArgs

from .module_printer import ModulePrinter
# from ..doc import DocPrinter
# from ..tests import TestPrinter

class GoBindingsPrinter(LanguageBindingsPrinter):

    def __init__(self, ydk_root_dir, bundle, generate_tests, sort_clazz):
        super(GoBindingsPrinter, self).__init__(ydk_root_dir, bundle, generate_tests, sort_clazz)

    def print_files(self):
        only_modules = [package.stmt for package in self.packages]
        size = len(only_modules)

        for index, package in enumerate(self.packages):
            self._print_module(index, package, size)

        # RST documentation
        # self._print_go_rst_toc()

        # if self.generate_tests:
        #     self._print_cmake_file(self.packages, self.bundle_name, self.test_dir)
        # return self.source_files

    def _print_module(self, index, package, size):
        print('Processing %d of %d %s' % (index + 1, size, package.stmt.pos.ref))
        # Skip generating module for empty modules
        if len(package.owned_elements) == 0:
            return

        # Set up module directory
        module_dir = self.initialize_output_directory(self.models_dir)

        # Generate go module
        self._print_go_module(package, self.models_dir)

        # RST documentation
        # self._print_cpp_rst_doc(package)
        
        # if self.generate_tests:
        #     self._print_tests(package, self.test_dir)

    def _print_go_module(self, package, path):
        go_module_file_name = '%s/%s.go' % (path, package.name)
        with open(go_module_file_name, 'w+') as file_descriptor:
            self.ypy_ctx.fd = file_descriptor
            mp = ModulePrinter(self.ypy_ctx, self.bundle_name, self.sort_clazz)
            mp.print_output(package)

