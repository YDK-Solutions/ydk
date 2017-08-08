#  ----------------------------------------------------------------
# Copyright 2017 Cisco Systems
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
class_enum_printer.py

 prints Go enum class method

"""

from ydkgen.api_model import Enum
from ydkgen.common import get_module_name
from ydkgen.printer.meta_data_util import get_enum_class_docstring

class EnumPrinter(object):

    def __init__(self, ctx):
        self.ctx = ctx

    def print_enum(self, enum_class):
        assert isinstance(enum_class, Enum)
        self._print_enum_header(enum_class)
        self._print_enum_body(enum_class)
        self._print_enum_trailer(enum_class)

    def _print_enum_header(self, enum_class):
        self.ctx.writeln('type %s string' % enum_class.qualified_go_name())
        self.ctx.bline()
        self.ctx.writeln('const (')
        self.ctx.lvl_inc()

    def _print_enum_body(self, enum_class):
        self._print_enum_literals(enum_class.qualified_go_name(), enum_class)

    def _print_enum_literals(self, enum_class_prefix, enum_class):
        for enum_literal in enum_class.literals:
            self._print_enum_literal(enum_class_prefix, enum_literal)

    def _print_enum_literal(self, enum_class_prefix, enum_literal):
        name = enum_literal.name
        value = enum_literal.value
        self.ctx.writeln('%s_%s %s = "%s"' % (enum_class_prefix,
                                            name,
                                            enum_class_prefix,
                                            enum_literal.stmt.arg))

    def _print_enum_trailer(self, enum_class):
        self.ctx.lvl_dec()
        self.ctx.writeln(')')
        self.ctx.bline()
