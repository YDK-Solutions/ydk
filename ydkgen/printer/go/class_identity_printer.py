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
identity_printer.py

  print yang Identity classes
"""

class IdentityPrinter(object):
    def __init__(self, ctx, bundle_name, identity_subclasses):
        self.ctx = ctx
        self.bundle_name = bundle_name
        self.identity_subclasses = identity_subclasses

    def print_identity(self, clazz):
        self._print_identity_class_body(clazz)
        self._print_identity_class_string(clazz)

    def _print_identity_class_body(self, clazz):
        self.ctx.writeln('type {} struct {{'.format(clazz.qualified_go_name()))
        self.ctx.writeln('}')
        self.ctx.bline()

    def _print_identity_class_string(self, clazz):
        self.ctx.writeln('func (id {}) String() string {{'.format(clazz.qualified_go_name()))
        self.ctx.writeln('\treturn "{}:{}"'.format(clazz.module.arg, clazz.stmt.arg))
        self.ctx.writeln('}')
        self.ctx.bline()
