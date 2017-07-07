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

"""
module_printer.py

 prints Go classes

"""
from ydkgen.api_model import Class, Enum, Bits
from .class_printer import ClassPrinter
# from .class_enum_printer import EnumPrinter

from ydkgen.printer.file_printer import FilePrinter


class ModulePrinter(FilePrinter):

    def __init__(self, ctx, bundle_name, sort_clazz):
        super(ModulePrinter, self).__init__(ctx)
        self.bundle_name = bundle_name
        self.sort_clazz = sort_clazz

    def print_header(self, package):
        self.ctx.writeln('package %s' % package.name)
        self.ctx.bline()
        self._print_imports(package)

    def print_body(self, package):
        classes = [clazz for clazz in package.owned_elements if isinstance(clazz, Class)]
        for clazz in classes:
            self._print_class(clazz)

    def print_extra(self, package):
        # self._print_enums(package, multi_file.class_list)
        pass

    def print_trailer(self, package):
        pass

    def _print_class(self, clazz):
        cp = ClassPrinter(self.ctx, self.bundle_name, self.sort_clazz)
        cp.print_output(clazz)

    def _print_enums(self, package, classes):
        # self.enum_printer.print_enum_to_string_funcs(package, classes)
        pass

    def _print_imports(self, package):
        self.ctx.writeln('import (')
        self.ctx.lvl_inc()
        self._print_static_imports()
        self._print_derived_imports(package)
        self.ctx.lvl_dec()
        self.ctx.writeln(')')
        self.ctx.bline()

    def _print_static_imports(self):
        self.ctx.writeln('"fmt"')
        self.ctx.writeln('"github.com/CiscoDevNet/ydk-go/ydk/types"')

    def _print_derived_imports(self, package):
        derived_imports = ['"github.com/CiscoDevNet/ydk-go/ydk/%s"' % self.bundle_name]
        for imported_type in package.imported_types():
            if( self.is_derived_identity(package, imported_type) ):
                stmt = '"github.com/CiscoDevNet/ydk-go/ydk/models/%s"' % (
                    imported_type.get_py_mod_name().split('.')[2])
                if stmt not in derived_imports:
                    derived_imports.append(stmt)

        self.ctx.writelns(derived_imports)
        self.ctx.bline()
