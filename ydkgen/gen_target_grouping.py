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
from collections import OrderedDict
from pyang.types import EnumerationTypeSpec
from pyang import statements
from pyang.error import EmitError
"""YANG to Python emitter"""

import sys
from types import *
import optparse
import copy
import re

#from errors import YPYError, YPYDataValidationError
from common import *
from pyang.statements import Statement
from api_model import Class
from api_model import Package
from api_model import Property
from api_model import Enum
from helper import snake_case


def emit_comment(ctx, comment):
    ctx.printer.comment(lines=comment.split('\n'))


def module_xlat(ctx, s):
    # ::::::::::::::::::::::::::::::::::::::::
    # Print the header
    # ::::::::::::::::::::::::::::::::::::::::
    ctx.printer.header()
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
    comment = s.search_one('description')
    if comment is not None:
        ctx.comment = yang_id(comment)
        emit_comment(ctx, yang_id(comment))
    ctx.ns += [(yang_id(s), ctx.namespace)]
    typdef = s.search('typedef')
    ctx.types = []
    for typ in typdef:
        ctx.types.append(typ)
    if s.i_typedefs is not None:
        for key in s.i_typedefs:
            typ = s.i_typedefs[key]
            if not typ in ctx.types and belongs_to(s, typ) == True:
                ctx.types.append(typ)


def emit_yang_ns(ctx):

    ctx.printer.yang_ns_header()
    ns_list = []

    sorted_keys = sorted(ctx.ctx.modules.keys())

#     for key in ctx.ctx.modules:
    for key in sorted_keys:
        m = ctx.ctx.modules[key]
        ns = m.search_one('namespace')
        if ns is not None:
            ns_list.append((key[0].replace('-', '_'), ns.arg, yang_id(m)))

    ctx.printer.namespace(ns_list)


def emit_module(ctx, stmt):
    module_xlat(ctx, stmt)

    p = Package()
    stmt.i_package = p
    p.stmt = stmt

    emit_walk(ctx, stmt, p)

    if p is not None:
        for element in p.owned_elements:
            if isinstance(element, Class):
                ctx.printer.print_class(element)
            elif isinstance(element, Enum):
                ctx.printer.print_enum(element)

    for imported_type in p.imported_types():
        ctx.printer.print_import(imported_type)


def emit_walk(ctx, stmt, parent_element):
    # process typedefs first so that they are resolved
    # when we have to use them

    element = parent_element

    # walk the groupings first
    if hasattr(stmt, 'i_groupings'):
        for grouping_name in stmt.i_groupings:
            emit_walk(ctx, stmt.i_groupings[grouping_name], element)

    # get the groupings from the uses stmt also
    for uses_stmt in stmt.search('uses'):
        grouping_stmt = uses_stmt.i_grouping
        # add a check in here to ensure that
        # the grouping belongs to the same module
        if not hasattr(grouping_stmt, 'i_class'):
            emit_walk(ctx, grouping_stmt, element)

    def _get_enum_type_stmt(stmt):
        type_stmt = stmt.search_one('type')
        if hasattr(type_stmt, 'i_typedef') and type_stmt.i_typedef is not None:
            typedef_stmt = type_stmt.i_typedef
            return _get_enum_type_stmt(typedef_stmt)
        elif hasattr(type_stmt, 'i_type_spec') and isinstance(type_stmt.i_type_spec, EnumerationTypeSpec):
            return type_stmt
        else:
            return None

    if hasattr(stmt, 'i_typedefs'):
        for typedef_stmt_name in stmt.i_typedefs:
            typedef_stmt = stmt.i_typedefs[typedef_stmt_name]
            enum_type_stmt = _get_enum_type_stmt(typedef_stmt)
            if enum_type_stmt is not None:
                enum_class = Enum()
                enum_class.stmt = enum_type_stmt
                parent_element.owned_elements.append(enum_class)
                enum_class.owner = parent_element

    if stmt.keyword == 'module':
        clazz = Class()
        clazz.stmt = stmt
        stmt.i_class = clazz
        parent_element.owned_elements.append(clazz)
        clazz.owner = parent_element
        parent_element = clazz
        element = parent_element
    elif stmt.keyword == 'grouping':
        clazz = Class()
        stmt.i_class = clazz
        clazz.stmt = stmt
        parent_element.owned_elements.append(clazz)
        clazz.owner = parent_element
        element = clazz

    elif stmt.keyword == 'container' or stmt.keyword == 'list':
        clazz = Class()
        stmt.i_class = clazz
        clazz.stmt = stmt
        parent_element.owned_elements.append(clazz)
        clazz.owner = parent_element
        element = clazz

        # create a property along with the class
        prop = Property()
        stmt.i_property = prop
        prop.stmt = stmt
        prop.property_type = clazz
        parent_element.owned_elements.append(prop)
        prop.owner = parent_element

    elif stmt.keyword == 'leaf' or stmt.keyword == 'leaf-list':
        prop = Property()
        stmt.i_property = prop
        prop.stmt = stmt
        parent_element.owned_elements.append(prop)
        prop.owner = parent_element

        enum_type = _get_enum_type_stmt(stmt)
        if enum_type is not None:
            if not hasattr(enum_type, 'i_enum'):
                # we have to create the enum
                enum_class = Enum()
                enum_class.stmt = enum_type
                enum_type.i_enum = enum_class

                if enum_type.i_module == stmt.i_module:
                    # we can add to the same package
                    parent_element.owned_elements.append(enum_class)
                    enum_class.owner = parent_element
            prop.property_type = enum_type.i_enum

    if hasattr(stmt, 'i_children'):
        for uses_stmt in stmt.search('uses'):
            grouping_stmt = uses_stmt.i_grouping
            element.extends.append(grouping_stmt.i_class)
        chs = [ch for ch in stmt.i_children
               if ch.keyword in statements.data_definition_keywords]
        if stmt.keyword != 'grouping':
            chs = [ch for ch in chs if ch.parent.keyword != 'grouping']

        for child_stmt in chs:
            emit_walk(ctx, child_stmt, element)
