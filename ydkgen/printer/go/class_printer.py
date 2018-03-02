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
from ydkgen.api_model import Class, Enum, snake_case
from ydkgen.common import get_qualified_yang_name, sort_classes_at_same_level

from .function_printer import FunctionPrinter
from .class_constructor_printer import ClassConstructorPrinter
from .class_enum_printer import EnumPrinter
from .class_get_entity_path_printer import GetSegmentPathPrinter
from .class_get_child_printer import ClassGetChildPrinter
from .class_get_children_printer import ClassGetChildrenPrinter


class ClassPrinter(object):
    def __init__(self, ctx, bundle_name, identity_subclasses):
        self.ctx = ctx
        self.bundle_name = bundle_name
        self.identity_subclasses = identity_subclasses
        self.enum_printer = EnumPrinter(ctx)

    def print_output(self, clazz):
        leafs = []
        children = []
        self._get_class_members(clazz, leafs, children)
        self._print_class_constructor(clazz, leafs, children)
        self._print_class_method_definitions(clazz, leafs, children)
        self._print_child_classes(clazz)
        self._print_child_enums(clazz)

    def _print_class_constructor(self, clazz, leafs, children):
        ClassConstructorPrinter(self.ctx, clazz, leafs, self.identity_subclasses).print_all()

    def _print_class_method_definitions(self, clazz, leafs, children):
        self._print_class_get_entity_common_data(clazz, leafs, children)
        self._print_class_get_filter(clazz)
        self._print_class_set_filter(clazz)
        self._print_class_get_go_name(clazz, leafs, children)
        self._print_class_get_segment_path(clazz)
        self._print_class_get_child(clazz, leafs, children)
        self._print_class_get_children(clazz, leafs, children)
        self._print_class_get_leafs(clazz, leafs)
        self._print_yang_models_function(clazz)
        self._print_get_capabilities_table(clazz)
        self._print_get_namespace_table(clazz)
        self._print_set_parent_function(clazz)
        self._print_get_parent_function(clazz)
        self._print_get_parent_yang_name_function(clazz)

    def _get_class_members(self, clazz, leafs, children):
        for prop in clazz.properties():
            ptype = prop.property_type
            if isinstance(ptype, Class) and not ptype.is_identity():
                children.append(prop)
            elif ptype is not None:
                leafs.append(prop)

    def _print_class_get_entity_common_data(self, clazz, leafs, children):
        fp = FunctionPrinter(self.ctx, clazz)
        fp.print_function_header_helper('GetCommonEntityData', return_type='*types.CommonEntityData')

        fp.ctx.writeln('%s.EntityData.YangName = "%s"' % (fp.class_alias, fp.clazz.stmt.arg))
        fp.ctx.writeln('%s.EntityData.BundleName = "%s"' % (fp.class_alias, self.bundle_name.lower()))
        fp.ctx.writeln('%s.EntityData.ParentYangName = "%s"' % (fp.class_alias, clazz.owner.stmt.arg))
        fp.ctx.bline()

        # GoName - returns string
        # TODO

        # Children - returns map[string]Entity
        # TODO

        fp.ctx.writeln('%s.EntityData.Leafs = make(map[string]interface{})' % fp.class_alias)
        for leaf in leafs:
            fp.ctx.writeln('{0}.EntityData.Leafs["{1}"] = {0}.{2}'.format(
                fp.class_alias, leaf.stmt.arg, leaf.go_name()))

        fp.ctx.writeln('return &(%s.EntityData)' % fp.class_alias)
        fp.print_function_trailer()

    # GetFilter
    def _print_class_get_filter(self, clazz):
        fp = FunctionPrinter(self.ctx, clazz)
        rstmt = '%s.YFilter' % fp.class_alias
        fp.quick_print('GetFilter', return_type='yfilter.YFilter', return_stmt=rstmt)

    # SetFilter
    def _print_class_set_filter(self, clazz):
        fp = FunctionPrinter(self.ctx, clazz)
        stmt = '%s.YFilter = yf' % fp.class_alias
        fp.quick_print('SetFilter', args='yf yfilter.YFilter', stmt=stmt)

    # GetGoName
    def _print_class_get_go_name(self, clazz, leafs, children):
        fp = FunctionPrinter(self.ctx, clazz)
        fp.print_function_header_helper('GetGoName', args='yname string', return_type='string')
        for leaf in leafs:
            fp.ctx.writeln('if yname == "%s" { return "%s" }' % (leaf.stmt.arg, leaf.go_name()))
        for child in children:
            fp.ctx.writeln('if yname == "%s" { return "%s" }' % (get_qualified_yang_name(child), child.go_name()))
        fp.ctx.writeln('return ""')
        fp.print_function_trailer()

    # GetSegmentPath
    def _print_class_get_segment_path(self, clazz):
        fp = GetSegmentPathPrinter(self.ctx, clazz)
        fp.print_all()

    # GetChildByName
    def _print_class_get_child(self, clazz, leafs, children):
        fp = ClassGetChildPrinter(self.ctx, clazz, leafs, children)
        fp.print_all()

    # GetChildren
    def _print_class_get_children(self, clazz, leafs, children):
        fp = ClassGetChildrenPrinter(self.ctx, clazz, leafs, children)
        fp.print_all()

    # GetLeafs
    def _print_class_get_leafs(self, clazz, leafs):
        fp = FunctionPrinter(self.ctx, clazz)
        fp.print_function_header_helper('GetLeafs', return_type='map[string]interface{}')
        fp.ctx.writeln('leafs := make(map[string]interface{})')
        for leaf in leafs:
            fp.ctx.writeln('leafs["{1}"] = {0}.{2}'.format(
                fp.class_alias, leaf.stmt.arg, leaf.go_name()))
        fp.ctx.writeln('return leafs')
        fp.print_function_trailer()

    # GetBundleYangModelsLocation
    def _print_yang_models_function(self, clazz):
        bundle_name = snake_case(self.bundle_name)
        fp = FunctionPrinter(self.ctx, clazz)
        rstmt = '%s.GetModelsPath()' % bundle_name
        fp.quick_print('GetBundleYangModelsLocation', return_type='string', return_stmt=rstmt)

    # GetCapabilitiesTable
    def _print_get_capabilities_table(self, clazz):
        bundle_name = snake_case(self.bundle_name)
        fp = FunctionPrinter(self.ctx, clazz)
        fp.print_function_header_helper(
            'GetCapabilitiesTable', return_type='map[string]string')
        fp.ctx.write("return %s.GetCapabilities() " % (bundle_name))
        fp.print_function_trailer()

    # GetNamespaceTable
    def _print_get_namespace_table(self, clazz):
        bundle_name = snake_case(self.bundle_name)
        fp = FunctionPrinter(self.ctx, clazz)
        fp.print_function_header_helper(
            'GetNamespaceTable', return_type='map[string]string')
        fp.ctx.write("return %s.GetNamespaces() " % (bundle_name))
        fp.print_function_trailer()

    # SetParent
    def _print_set_parent_function(self, clazz):
        fp = FunctionPrinter(self.ctx, clazz)
        stmt = '%s.parent = parent' % fp.class_alias
        fp.quick_print('SetParent', args='parent types.Entity', stmt=stmt)

    # GetParent
    def _print_get_parent_function(self, clazz):
        fp = FunctionPrinter(self.ctx, clazz)
        rstmt = '%s.parent' % fp.class_alias
        fp.quick_print('GetParent', return_type='types.Entity', return_stmt=rstmt)

    # GetParentYangName
    def _print_get_parent_yang_name_function(self, clazz):
        fp = FunctionPrinter(self.ctx, clazz)
        rstmt = '"%s"' % clazz.owner.stmt.arg
        fp.quick_print('GetParentYangName', return_type='string', return_stmt=rstmt)

    def _print_child_classes(self, parent):
        unsorted_classes = [nested_class for nested_class in parent.owned_elements if isinstance(nested_class, Class)]
        sorted_classes = sort_classes_at_same_level(unsorted_classes)

        for clazz in sorted_classes:
            cp = ClassPrinter(self.ctx, self.bundle_name, self.identity_subclasses)
            cp.print_output(clazz)

    def _print_child_enums(self, parent):
        enumz = []
        enumz.extend([nested_enum for nested_enum in parent.owned_elements
                        if isinstance(nested_enum, Enum)])

        for nested_enum in sorted(enumz, key=lambda e: e.name):
            self.enum_printer.print_enum(nested_enum)
