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
class_inits_printer.py

 __init__ printer

"""
from pyang.types import UnionTypeSpec, PathTypeSpec
from ydkgen.api_model import Bits, Class, Package, DataType, Enum
from ydkgen.builder import TypesExtractor


class ClassInitsPrinter(object):

    def __init__(self, ctx):
        self.ctx = ctx

    def print_output(self, clazz, leafs, children):
        self._print_class_inits_header(clazz)
        self._print_class_inits_body(clazz, leafs, children)
        self._print_class_inits_trailer(clazz)

    def _print_class_inits_header(self, clazz):
        self.ctx.writeln('def __init__(self):')
        self.ctx.lvl_inc()

    def _print_class_inits_body(self, clazz, leafs, children):
        if clazz.is_identity():
            arg = '"%s:%s"' % (clazz.module.arg, clazz.stmt.arg)
            line = 'super(%s, self).__init__(%s)' % (clazz.name, arg)
            self.ctx.writeln(line)
        else:
            self.ctx.writeln('super(%s, self).__init__()' % clazz.name)
            self._print_class_inits_body_helper(clazz, leafs, children)
            self.ctx.writeln('self.yang_name = "%s"' % clazz.stmt.arg)
            self.ctx.writeln('self.yang_parent_name = "%s"' % clazz.owner.stmt.arg)

    def _print_class_inits_body_helper(self, clazz, leafs, children):
        yleafs = self._get_leafs(clazz)
        yleaf_lists = self._get_leaf_lists(clazz)

        output = []

        for prop in leafs:
            leaf_type = None
            if prop in yleafs:
                leaf_type = 'YLeaf'
            elif prop in yleaf_lists:
                leaf_type = 'YLeafList'

            output.append('self.%s = %s(YType.%s, "%s")' 
                % (prop.name, leaf_type, self._get_type_name(prop.property_type), prop.stmt.arg))

        for child in children:
            if not child.is_many:
                if (child.stmt.search_one('presence') is None):
                    output.append('self.%s = %s()' % (child.name, child.property_type.qn()))
                    output.append('self.%s.parent = self' % child.name)
                else:
                    output.append('self.%s = None' % (child.name))

        if len(output) > 0:
            for line in output:
                self.ctx.writeln(line)

    def _print_class_inits_trailer(self, clazz):
        self.ctx.lvl_dec()
        self.ctx.bline()

    def _get_type_name(self, prop_type):
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

    def _get_leafs(self, clazz):
        leafs = []
        for child in clazz.owned_elements:
            if child.stmt.keyword == 'leaf':
                leafs.append(child)
        return leafs

    def _get_leaf_lists(self, clazz):
        leaf_lists = []
        for child in clazz.owned_elements:
            if child.stmt.keyword == 'leaf-list':
                leaf_lists.append(child)
        return leaf_lists

class ClassSetAttrPrinter(object):

    def __init__(self, ctx):
        self.ctx = ctx

    def print_setattr(self):
        self._print_class_setattr_header()
        self._print_class_setattr_body()
        self._print_class_setattr_trailer()

    def _print_class_setattr_header(self):
        self.ctx.writeln('def __setattr__(self, name, value):')
        self.ctx.lvl_inc()

    def _print_class_setattr_body(self):
        self.ctx.writeln('if name in ("name", "number") and name in self.__dict__:')
        self.ctx.lvl_inc()
        self.ctx.writeln('self.__dict__[name].set(value)')
        self.ctx.lvl_dec()
        self.ctx.writeln('else:')
        self.ctx.lvl_inc()
        self.ctx.writeln('self.__dict__[name] = value')
        self.ctx.lvl_dec()

    def _print_class_setattr_trailer(self):
        self.ctx.lvl_dec()
        self.ctx.bline()
