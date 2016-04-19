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
  gen_target.py 
  YANG model driven API, python emitter.
  
"""
from api_model import Bits
from api_model import Class
from api_model import Enum
from common import yang_id
from helper import convert_to_reStructuredText, get_sphinx_ref_label, is_config_stmt
from meta_data_util import get_meta_info_data


def emit_comment(ctx, comment):
    ctx.printer.comment(lines=comment.split('\n'))


def _is_entity_clazz(c):
    return isinstance(c, Class) and not (c.is_identity() or c.is_grouping())


def _tree_prefix(indent):
    return ''.join([' ' for _ in range(1, indent)])


def _construct_prop_tree(ctx, prop, indent):
    """ Print the prop tree

    """
    property_type = prop.property_type
    node_prefix = convert_to_reStructuredText('%s' % _tree_prefix(indent))

    node_name = convert_to_reStructuredText(prop.name)
    # ctx.writeln('| %s\|'%node_prefix)
    config_or_oper = 'configuration'
    if not is_config_stmt(prop.stmt):
        config_or_oper = 'operational'
    meta_info_data = get_meta_info_data(
        prop, prop.property_type, prop.stmt.search_one('type'))

    type_display = meta_info_data.doc_link

    ctx.writeln('| %s\|\-\- **%s** (%s, %s)' % (node_prefix,
                                                node_name,
                                                type_display,
                                                config_or_oper))

    if _is_entity_clazz(property_type):
        new_indent = indent + 4
        for p in property_type.properties():
            _construct_prop_tree(ctx, p, new_indent)


def _construct_tree(ctx, package):
    """ Construct the class hierarchy

    """
    top_classes = [c for c in package.owned_elements if _is_entity_clazz(c)]

    if len(top_classes) > 0:
        ctx.printer.rst_printer.writeln(
            '.. topic:: Entity Tree defined by %s.py' % convert_to_reStructuredText(package.name))
        ctx.printer.rst_printer.lvl_inc()
        ctx.printer.rst_printer.bline()

        for c in top_classes:
            config_or_oper = 'configuration'
            if not is_config_stmt(c.stmt):
                config_or_oper = 'operational'

            clazz_name = convert_to_reStructuredText(c.name)

            # :ref:`Link title <label-name>`
            ctx.printer.rst_printer.writeln(
                '| :ref:`%s <%s>` (%s)' % (clazz_name, get_sphinx_ref_label(c), config_or_oper))
            # ctx.printer.rst_printer.writeln('| :class: `%s.%s`'%(c.get_py_mod_name(), c.qn()))
            # ctx.printer.rst_printer.writeln('| :class: `%s`'%(c.qn()))
            for p in c.properties():
                _construct_prop_tree(ctx.printer.rst_printer, p, 4)
        ctx.printer.rst_printer.writeln('| ')
        ctx.printer.rst_printer.writeln('| ')
        ctx.printer.rst_printer.bline()
        ctx.printer.rst_printer.lvl_dec()
        ctx.printer.rst_printer.bline()


def emit_module_header(ctx, package, mheader=None, is_meta=False):
    # ::::::::::::::::::::::::::::::::::::::::
    # Print the header
    # ::::::::::::::::::::::::::::::::::::::::
    s = package.stmt
    if is_meta:
        rpcs = [idx for idx in package.owned_elements if isinstance(idx, Class) and idx.is_rpc()]
        anyxml_import = ''
        if len(rpcs) > 0:
            anyxml_import = ', ANYXML_CLASS'
        ctx.printer.meta_header(anyxml_import)
    else:
        comment = s.search_one('description')
        ctx.writeln('""" %s ' % package.name)
        ctx.bline()

        if comment is not None and not is_meta:
            ctx.comment = comment.arg
            for line in ctx.comment.split('\n'):
                ctx.writeln(convert_to_reStructuredText(line))
        ctx.bline()
        ctx.writeln('"""')
        ctx.printer.header(mheader)
        ctx.printer.imports(package)
    ctx.root = s
    ctx.augment_path = ''
    ctx.aug_stmt = None
    ctx.module_name = yang_id(s)
    ctx.module = s
    # get the yang meta information.
    prefix = s.search_one('prefix')
    if prefix is not None:
        ctx.prefix = yang_id(prefix)
    namespace = s.search_one('namespace')
    if namespace is not None:
        ctx.namespace = yang_id(namespace)
    org = s.search_one('organization')
    if org is not None:
        ctx.organization = yang_id(org)
    contact = s.search_one('contact')
    if contact is not None:
        ctx.contact = yang_id(contact)
    revision = s.search_one('revision')
    if revision is not None:
        ctx.revision = yang_id(revision)

    ctx.ns += [(yang_id(s), ctx.namespace)]


def emit_yang_ns(ctx, packages):

    ctx.printer.print_yang_ns_header()
    ns_list = []
    namespace_map = {}
    for m in [p.stmt for p in packages]:

        ns = m.search_one('namespace')
        if ns is not None:
            ns_list.append((m.arg.replace('-', '_'), ns.arg, yang_id(m)))

    for package in packages:
        ns = package.stmt.search_one('namespace')
        for ele in package.owned_elements:
            if hasattr(ele, 'stmt') and ele.stmt is not None and ele.stmt.keyword == 'container':                
                namespace_map[(ns.arg, ele.stmt.arg)] = (package.get_py_mod_name(), ele.name)


    ctx.printer.print_namespaces(ns_list)
    ctx.printer.print_identity_map(packages)
    ctx.printer.print_namespaces_map(namespace_map)


def emit_importests(ctx, packages):
    ctx.printer.print_import_tests(packages)


def emit_python_rst(ctx, named_element):
    ctx.printer.print_python_rst(named_element)


def emit_ydk_models_rst(ctx, packages):
    ctx.printer.print_ydk_models_rst(packages)


def emit_module(ctx, package, mheader):

    emit_module_header(ctx, package, mheader=mheader)

    if package is not None:
        emit_module_enums(ctx, package)
        emit_module_bits(ctx, package)
        emit_module_classes(ctx, package)
        ctx.bline()


def emit_module_enums(ctx, package):
    enumz = []
    enumz.extend(
        [element for element in package.owned_elements if isinstance(element, Enum)])
    for nested_enumz in sorted(enumz, key=lambda element: element.name):
        ctx.printer.print_enum(nested_enumz)


def emit_module_bits(ctx, package):
    bits = []
    bits.extend(
        [bit for bit in package.owned_elements if isinstance(bit, Bits)])
    for bit in sorted(bits, key=lambda bit: bit.name):
        ctx.printer.print_bits(bit)


def emit_module_classes(ctx, package):
    ctx.printer.print_classes_at_same_level(
        [clazz for clazz in package.owned_elements if isinstance(clazz, Class)])


def emit_test_module(ctx, package):
    ctx.printer.print_testcases(package)


def emit_meta(ctx, package):
    ctx.print_meta = True
    emit_module_header(ctx, package, is_meta=True)
    if package is not None:
        emit_meta_table_open(ctx)
        for nested_enumz in [e for e in package.owned_elements if isinstance(e, Enum)]:
            ctx.printer.print_enum_meta(nested_enumz)
        ctx.printer.print_classes_meta([c for c in package.owned_elements if isinstance(c, Class)])
        emit_meta_table_close(ctx)
        ctx.printer.print_classes_meta_parents(
            [c for c in package.owned_elements if isinstance(c, Class)])


def emit_meta_table_open(ctx):
    ctx.writeln('_meta_table = {')
    ctx.lvl_inc()


def emit_meta_table_close(ctx):
    ctx.lvl_dec()
    ctx.writeln('}')

def emit_deviation(ctx, package):
    ctx.printer.print_deviation(package)
