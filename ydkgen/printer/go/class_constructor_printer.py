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

 prints Go class constructor

"""

from ydkgen.api_model import Bits, Class, DataType, Enum
from .function_printer import FunctionPrinter

class ClassConstructorPrinter(FunctionPrinter):
    def __init__(self, ctx, clazz, leafs):
        super(ClassConstructorPrinter, self).__init__(ctx, clazz, leafs)

    def print_function_header(self):
        self.ctx.writeln('//////////////////////////////////////////////////////////////////////////')
        self.ctx.writeln('// %s' % self.clazz.qualified_go_name())
        self.ctx.writeln('//////////////////////////////////////////////////////////////////////////')
        self.ctx.writeln('type %s struct {' % self.clazz.qualified_go_name())
        self.ctx.lvl_inc()

    def print_function_body(self):
        self.ctx.writeln('parent types.Entity')
        self.ctx.writeln('Filter types.YFilter')
        self.ctx.bline()
        self._print_inits()

    def _print_inits(self):
        self._print_leaf_inits()
        self._print_children_inits()

    def _print_leaf_inits(self):
        index = 0
        while index < len(self.leafs):
            prop = self.leafs[index]
            index += 1
            leaf_name = ''
            if prop.stmt.i_module.arg != self.clazz.stmt.i_module.arg:
                leaf_name = prop.stmt.i_module.arg + ':' + prop.stmt.arg
            else:
                leaf_name = prop.go_name()
            type_name = get_type_name(prop.property_type)

            if prop.is_many:
                line = '%s []%s' %(leaf_name, prop.qualified_go_name())
            else:
                line = '%s interface{} // %s' % (leaf_name, type_name)
            self.ctx.writeln(line)
    
    def _print_children_inits(self):
        for prop in self.clazz.properties():
            if not prop.is_many:
                self._print_child_inits_unique(prop)
            elif prop.stmt.keyword != 'anyxml':
                self._print_child_inits_many(prop)

    def _print_child_inits_unique(self, prop):
        if (isinstance(prop.property_type, Class)
            and not prop.property_type.is_identity()):
            presence_stmt = ''
            if prop.property_type.stmt.search_one('presence') is not None:
                presence_stmt = ' // presence node'
            self.ctx.writeln('%s %s%s' % (
                prop.go_name(), prop.qualified_go_name(), presence_stmt))

    def _print_child_inits_many(self, prop):
        if (prop.is_many and isinstance(prop.property_type, Class)
            and not prop.property_type.is_identity()):
            self.ctx.writeln('%s []%s' % (
                prop.go_name(), prop.qualified_go_name()))

def get_type_name(prop_type):
    if prop_type.name == 'string':
        return 'str'
    elif prop_type.name == 'leafref':
        return 'str'
    elif prop_type.name == 'decimal64':
        return 'str'
    elif prop_type.name == 'union':
        return 'str'
    elif prop_type.name == 'binary':
        return 'str'
    elif prop_type.name == 'instance-identifier':
        return 'str'
    elif isinstance(prop_type, Bits):
        return 'bits'
    elif isinstance(prop_type, Class) and prop_type.is_identity():
        return 'identityref'
    elif isinstance(prop_type, Enum):
        return 'enumeration'
    elif isinstance(prop_type, DataType):
        return 'str'
    return prop_type.name
