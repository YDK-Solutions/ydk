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
source_printer.py

prints Go class method

"""
from .function_printer import FunctionPrinter
from ydkgen.common import get_qualified_yang_name


class ClassGetChildrenPrinter(FunctionPrinter):
    def __init__(self, ctx, clazz, leafs, children):
        super(ClassGetChildrenPrinter, self).__init__(ctx, clazz, leafs, children)

    def print_function_header(self):
        self.print_function_header_helper(
            'GetChildren', return_type='map[string]types.ChildStore')

    def print_function_body(self):
        self.ctx.writeln('children := make(map[string]types.ChildStore)')
        for child in self.children:
            path = get_qualified_yang_name(child)
            if child.is_many:
                self.ctx.writeln('children["%s"] = types.ChildStore{"%s", nil}' % (
                    path, child.go_name()))
                self._print_many(child)
            else:
                self.ctx.writeln('children["%s"] = types.ChildStore{"%s", &%s.%s}' % (
                    path, child.go_name(), self.class_alias, child.go_name()))
        self.ctx.writeln('return children')

    def _print_many(self, child):
        child_stmt = '%s.%s' % (self.class_alias, child.go_name())
        self.ctx.writeln('for i := range %s {' % (child_stmt))
        self.ctx.lvl_inc()
        child_stmt = '%s[i]' % child_stmt
        self.ctx.writeln('children[{0}.GetCommonEntityData().SegmentPath] = types.ChildStore{{"{1}", &{0}}}'.format(
            child_stmt, child.go_name()))
        self.ctx.lvl_dec()
        self.ctx.writeln('}')
