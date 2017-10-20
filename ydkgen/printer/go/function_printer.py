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

 prints Go classes

"""
class FunctionPrinter(object):
    def __init__(self, ctx, clazz, leafs=None, children=None):
        self.ctx = ctx
        self.clazz = clazz
        self.leafs = leafs
        self.children = children

        self.class_alias = clazz.go_name(case='lowerCamel')

    def print_all(self):
        self.print_function_header()
        self.print_function_body()
        self.print_function_trailer()

    def print_function_header_helper(self, name, args='', return_type=' '):
        if return_type != ' ':
            return_type = ' %s ' % return_type

        self.ctx.writeln('func (%s *%s) %s(%s)%s{' % (
            self.class_alias, self.clazz.qualified_go_name(),
            name, args, return_type))
        self.ctx.lvl_inc()

    def print_function_trailer(self):
        self.ctx.lvl_dec()
        self.ctx.writeln('}')
        self.ctx.bline()
