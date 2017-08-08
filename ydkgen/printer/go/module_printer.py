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
from .class_identity_printer import IdentityPrinter
# from .class_enum_printer import EnumPrinter

from ydkgen.printer.file_printer import FilePrinter
from ydkgen.common import convert_to_reStructuredText


class ModulePrinter(FilePrinter):

    def __init__(self, ctx, bundle_name, sort_clazz, identity_subclasses):
        super(ModulePrinter, self).__init__(ctx)
        self.bundle_name = bundle_name
        self.sort_clazz = sort_clazz
        self.identity_subclasses = identity_subclasses
        self.cp = ClassPrinter(ctx, bundle_name, sort_clazz, identity_subclasses)
        self.ip = IdentityPrinter(ctx, bundle_name, sort_clazz, identity_subclasses)

    def print_header(self, package):
        self._print_package_description(package)
        self.ctx.writeln('package %s' % package.name)
        self.ctx.bline()
        self._print_imports(package)
        self._print_init(package)

    def print_body(self, package):
        classes = [clazz for clazz in package.owned_elements if isinstance(clazz, Class)]
        for clazz in classes:
            if clazz.is_identity():
                self.ip.print_identity(clazz)
            else:
                self.cp.print_output(clazz)

    def print_extra(self, package):
        # self._print_enums(package, multi_file.class_list)
        pass

    def print_trailer(self, package):
        pass

    def _print_package_description(self, package):
        comment = package.stmt.search_one('description')
        if comment is not None:
            comment = comment.arg
            for line in comment.split('\n'):
                self.ctx.writeln("// %s" % convert_to_reStructuredText(line))

    def _print_class(self, clazz):
        self.cp.print_output(clazz)
        self.ip.print_identity(clazz)

    def _print_enums(self, package, classes):
        # self.enum_printer.print_enum_to_string_funcs(package, classes)
        pass

    def _print_imports(self, package):
        self.ctx.writeln('import (')
        self.ctx.lvl_inc()
        self._print_static_imports(package)
        self._print_derived_imports(package)
        self.ctx.lvl_dec()
        self.ctx.writeln(')')
        self.ctx.bline()

    def _print_init(self, package):
        self.ctx.writeln('func init() {')
        self.ctx.lvl_inc()
        self.ctx.writeln('fmt.Println("Registering top level entities for package {}")'.format(package.name))
        for e in package.owned_elements:
            ns = package.stmt.search_one('namespace')
            if ns is not None and isinstance(e, Class) and not e.is_identity():
                self.ctx.writeln('ydk.RegisterEntity("{{{} {}}}", reflect.TypeOf({}{{}}))'.format(ns.arg, e.stmt.arg, e.go_name()))
                self.ctx.writeln('ydk.RegisterEntity("{}:{}", reflect.TypeOf({}{{}}))'.format(package.stmt.arg, e.stmt.arg, e.go_name()))
        self.ctx.lvl_dec()
        self.ctx.writeln('}')
        self.ctx.bline()

    def _print_static_imports(self, package):
        self.ctx.writeln('"fmt"')
        has_top_entity = False
        for c in package.owned_elements:
            if isinstance(c, Class) and not c.is_identity():
                has_top_entity = True
                break;
        if has_top_entity:
            self.ctx.writeln('"github.com/CiscoDevNet/ydk-go/ydk"')
            self.ctx.writeln('"github.com/CiscoDevNet/ydk-go/ydk/types"')
            self.ctx.writeln('"github.com/CiscoDevNet/ydk-go/ydk/models/{}"'.format(self.bundle_name))
            self.ctx.writeln('"reflect"')

    def _print_derived_imports(self, package):
        # derived_imports = ['"github.com/CiscoDevNet/ydk-go/ydk/models/%s"' % self.bundle_name]
        derived_imports = []
        for imported_type in package.imported_types():
            if( self.is_derived_identity(package, imported_type) ):
                stmt = '"github.com/CiscoDevNet/ydk-go/ydk/models/%s"' % (
                    imported_type.get_py_mod_name().split('.')[2])
                if stmt not in derived_imports:
                    derived_imports.append(stmt)

        if len(derived_imports) > 0:
            self.ctx.writelns(derived_imports)
        self.ctx.bline()
