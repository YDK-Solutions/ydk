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
class_path_printer.py

Prints Go class methods

"""
from ydkgen.api_model import Bits, Enum, Package
from ydkgen.common import is_empty_prop
from .function_printer import FunctionPrinter

class GetSegmentPathPrinter(FunctionPrinter):
    def __init__(self, ctx, clazz):
        super(GetSegmentPathPrinter, self).__init__(ctx, clazz)

    def print_function_header(self):
        self.print_function_header_helper('GetSegmentPath', return_type='string')

    def print_function_body(self):
        path = ['"']
        prefix = ''
        if self.clazz.owner is not None:
            if isinstance(self.clazz.owner, Package):
                prefix = '%s:' % self.clazz.owner.stmt.arg
            elif self.clazz.owner.stmt.i_module.arg != self.clazz.stmt.i_module.arg:
                prefix = '%s:' % self.clazz.stmt.i_module.arg
        path.append('%s%s"' % (prefix, self.clazz.stmt.arg))

        key_props = self.clazz.get_key_props()
        for key_prop in key_props:
            prefix = ''
            if key_prop.stmt.i_module.arg != self.clazz.stmt.i_module.arg:
                prefix = '%s:' % (key_prop.stmt.i_module.arg)

            predicate = '''{0}"[%s%s='"{0}fmt.Sprintf("%%v", %s.%s){0}"']"'''

            predicate = predicate.format(' + ') % (
                prefix, key_prop.stmt.arg, self.class_alias, key_prop.go_name())
            path.append(predicate)

        self.ctx.writeln('return %s' % (''.join(path)))
