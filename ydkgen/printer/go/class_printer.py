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
        self._print_class_create_new_child(clazz, children)
        self._print_class_get_go_name(clazz, leafs, children)
        self._print_class_get_segment_path(clazz)
        self._print_class_get_children(clazz, leafs, children)
        self._print_yang_models_function(clazz)
        self._print_get_capabilities_table(clazz)
        self._print_get_namespace_table(clazz)

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

        data_alias = '%s.EntityData' % fp.class_alias
        fp.ctx.writeln('%s.YFilter = %s.YFilter' % (data_alias, fp.class_alias))
        fp.ctx.writeln('%s.YangName = "%s"' % (data_alias, fp.clazz.stmt.arg))
        fp.ctx.writeln('%s.BundleName = "%s"' % (data_alias, self.bundle_name.lower()))
        fp.ctx.writeln('%s.ParentYangName = "%s"' % (data_alias, clazz.owner.stmt.arg))
        fp.ctx.bline()

        # GoName - returns string
        # TODO

        # Children - returns map[string]Entity
        # self.ctx.writeln('%s.Children = make(map[string]types.Entity)' % data_alias)
        # for child in children:
        #     if child.is_many:
        #         child_stmt = '%s.%s' % (fp.class_alias, child.go_name())
        #         fp.ctx.writeln('for i := range %s {' % (child_stmt))
        #         fp.ctx.lvl_inc()
        #         child_stmt = '%s[i]' % child_stmt
        #         fp.ctx.writeln('{0}.Children[{1}.GetSegmentPath()] = &{1}'.format(data_alias, child_stmt))
        #         fp.ctx.lvl_dec()
        #         fp.ctx.writeln('}')
        #     else:
        #         path = get_qualified_yang_name(child)
        #         fp.ctx.writeln('%s.Children["%s"] = &%s.%s' % (
        #             data_alias, path, fp.class_alias, child.go_name()))
        self.ctx.writeln('%s.Children = %s.GetChildren()' % (data_alias, fp.class_alias))

        # # Leafs
        fp.ctx.writeln('%s.Leafs = make(map[string]interface{})' % data_alias)
        for leaf in leafs:
            fp.ctx.writeln('%s.Leafs["%s"] = %s.%s' % (
                data_alias, leaf.stmt.arg, fp.class_alias, leaf.go_name()))
        # self.ctx.writeln('%s.Leafs = %s.GetLeafs()' % (data_alias, fp.class_alias))


        fp.ctx.writeln('return &(%s)' % data_alias)
        fp.print_function_trailer()

    def _print_class_create_new_child(self, clazz, children):
        fp = FunctionPrinter(self.ctx, clazz)
        if len(children) == 0:
            fp.quick_print('CreateNewChild', args='childYangName string', return_type='types.Entity', return_stmt='nil')
        else:
            fp.print_function_header_helper('CreateNewChild', args='childYangName string', return_type='types.Entity')
            for child in children:
                args = (get_qualified_yang_name(child), child.qualified_go_name())
                fp.ctx.writeln('if "%s" == childYangName { return &%s{} }' % args)
            fp.ctx.writeln('return nil')
            fp.print_function_trailer()

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

    # GetChildren
    def _print_class_get_children(self, clazz, leafs, children):
        fp = ClassGetChildrenPrinter(self.ctx, clazz, leafs, children)
        fp.print_all()

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
