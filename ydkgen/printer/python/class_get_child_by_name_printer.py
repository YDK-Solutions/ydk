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

 prints python classes

"""


class ClassGetChildByNamePrinter(object):
    def __init__(self, ctx):
        self.ctx = ctx

    def print_class_get_child_by_name(self, clazz, children):
        self._print_class_get_child_header(clazz)
        self._print_class_get_child_body(children)
        self._print_class_get_child_trailer(clazz)

    def _print_class_get_child_header(self, clazz):
        self.ctx.writeln('def get_child_by_name(self, child_yang_name, segment_path):')
        self.ctx.lvl_inc()

    def _print_class_get_child_body(self, children):
        self._print_class_get_child_common()
        for child in children:
            self._print_class_get_child(child)
            self.ctx.bline()

    def _print_class_get_child_common(self):
        self.ctx.writeln('c1 = self.children.get(child_yang_name)')
        self.ctx.writeln('c2 = self.children.get(segment_path)')
        self.ctx.bline()
        self.ctx.writeln('if (c1 is not None):')
        self.ctx.lvl_inc()
        self.ctx.writeln('return c1')
        self.ctx.lvl_dec()
        self.ctx.writeln('elif (c2 is not None):')
        self.ctx.lvl_inc()
        self.ctx.writeln('return c2')
        self.ctx.lvl_dec()
        self.ctx.bline()

    def _print_class_get_child(self, child):
        self.ctx.writeln('if (child_yang_name == "%s"):' % (child.stmt.arg))
        self.ctx.lvl_inc()
        if child.is_many:
            self._print_class_get_child_many(child)
        else:
            self._print_class_get_child_unique(child)
        self.ctx.lvl_dec()

    def _print_class_get_child_many(self, child):
        self.ctx.writeln('for c in %s:' % child.name)
        self.ctx.lvl_inc()
        self.ctx.writeln('segment = c.get_segment_path();')
        self.ctx.writeln('if (segment_path == segment):')
        self.ctx.lvl_inc()
        self.ctx.writeln('children[segment_path] = c.get();')
        self.ctx.writeln('return children.get(segment_path);')
        self.ctx.lvl_dec()
        self.ctx.lvl_dec()
        self.ctx.writeln('c = %s()' % (child.property_type.qn()))
        self.ctx.writeln('c.parent = self')
        self.ctx.writeln('self.children[segment_path] = self.%s.get()' % child.name)
        self.ctx.writeln('return children.get(segment_path)')

    def _print_class_get_child_unique(self, child):
        self.ctx.writeln('if (self.%s is None):' % child.name)
        self.ctx.lvl_inc()
        self.ctx.writeln('%s = %s()' % (child.name, child.property_type.qn()))
        self.ctx.writeln('self.%s.parent = self' % child.name)
        self.ctx.lvl_dec()
        self.ctx.writeln('children["%s"] = %s.get()' % (child.stmt.arg, child.name))
        self.ctx.writeln('return children.get("%s")' % child.stmt.arg)

    def _print_class_get_child_trailer(self, clazz):
        self.ctx.writeln('return None')
        self.ctx.lvl_dec()
        self.ctx.bline()
