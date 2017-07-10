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
generated_entity_lookup_printer.py

 prints generated_entity_lookup.go file

"""
from ydkgen.api_model import Class, Enum, Bits
from .class_printer import ClassPrinter
# from .class_enum_printer import EnumPrinter

from ydkgen.printer.file_printer import FilePrinter
from ydkgen.common import convert_to_reStructuredText


class GeneratedEntityLookupPrinter(object):

    def __init__(self, ctx, bundle_name, packages):
        self.ctx = ctx
        self.bundle_name = bundle_name
        self.packages = packages

        self.packages.sort(lambda x, y: cmp(x.stmt.arg, y.stmt.arg))

    def print_output(self):
        self.print_header()
        self.print_lookup_table_function()

    def print_header(self):
        self.ctx.writeln('package %s' % self.bundle_name)
        self.ctx.bline()
        bundle_name_tuple = (self.bundle_name.title(), self.bundle_name)
        self.ctx.writeln('const Ydk%sModelsPath = "/usr/local/share/%s@0.1.1"' % bundle_name_tuple)
        self.ctx.bline()

    def print_lookup_table_function(self):
        self.ctx.writeln('func %sAugmentLookupTables() map[string]string {' % self.bundle_name.title())
        self.ctx.lvl_inc()
        self.ctx.writeln('caps := make(map[string]string)')
        for package in self.packages:
            self.ctx.writeln('caps["%s"] = "%s"' % (package.stmt.arg, package.revision))
        self.ctx.lvl_dec()
        self.ctx.writeln('}')

    def _print_packages(self):
        for package in self.packages:
            print package.stmt.arg