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
  api_model.py

 The meta-model of the generated API.
 Translation process converts the YANG model to classes defined in this module.
"""

from pyang import statements
from pyang.types import EnumerationTypeSpec, BitsTypeSpec, UnionTypeSpec, PathTypeSpec, \
    IdentityrefTypeSpec
from helper import camel_case, iskeyword
from helper import snake_case, escape_name
from common import YdkGenException



class Element(object):

    """
        The Element class.

        This is the super class of all modelled elements in the API.

        :attribute:: owned_elements
                    list of `Element` owned by this element

        :attribute:: owner
                    The owner of this `Element`.

        :attribute:: comment
                    The comments associated with this element.
    """

    def __init__(self):
        self.owned_elements = []
        self.owner = None
        self.comment = None

class Deviation(Element):
    def __init__(self):
        Element.__init__(self)
        self._stmt = None
        self.d_type = None
        self.d_target = None

    @property
    def stmt(self):
        return self._stmt

    @stmt.setter
    def stmt(self, stmt):
        self._stmt = stmt

    def qn(self):
        names = []
        stmt = self.d_target
        while stmt.parent:
            if stmt.keyword not in ('container', 'list', 'rpc'):
                names.append(self.convert_prop_name(stmt))
            else:
                names.append(self.convert_owner_name(stmt))
            stmt = stmt.parent

        return '.'.join(reversed(names))


    def convert_prop_name(self, stmt):
        name = snake_case(stmt.arg)
        if iskeyword(name):
            name = '%s_' % name

        if name.startswith('_'):
            name = '%s%s' % ('y', name)
        return name

    def convert_owner_name(self, stmt):
        name = escape_name(stmt.arg)
        if stmt.keyword == 'grouping':
            name = '%s_Grouping' % camel_case(name)
        elif stmt.keyword == 'identity':
            name = '%s_Identity' % camel_case(name)
        elif stmt.keyword == 'rpc':
            name = camel_case(name) + 'Rpc'
        else:
            name = camel_case(name)
        if iskeyword(name):
            name = '%s_' % name

        if name.startswith('_'):
            name = '%s%s' % ('Y', name)
        return name

    def get_package(self):
        if self.owner is None:
            return None
        if isinstance(self.owner, Package):
            return self.owner
        else:
            if hasattr(self.owner, 'get_package'):
                return self.owner.get_package()

class NamedElement(Element):

    ''' 

        An abstract element that may have a name
        The name is used for identification of the named element
        within the namespace that is defined or accessible

       :attribute:: name
                   The name of the Element

    '''

    def __init__(self):
        ''' The name of the named element'''
        super(NamedElement, self).__init__()
        self.name = None

    def get_py_mod_name(self):
        """
            Get the python module name that contains this
            NamedElement.
        """
        pkg = self
        while pkg is not None and not isinstance(pkg, Package):
            pkg = pkg.owner
        if pkg is None:
            return ''
        t = pkg.name.split('_')
        t = [n for n in t if n.lower() not in ['cisco', 'ios', 'xr']]
        if t:
            sub = t[0].lower()
            if iskeyword(sub):
                sub = '%s_' % sub
            sub = '.%s' % sub
        else:
            sub = ''
        py_mod_name = 'ydk.models%s.%s' % (sub, pkg.name)
        return py_mod_name

    def get_meta_py_mod_name(self):
        """
            Get the python meta module that contains the meta model
            information about this NamedElement.

        """
        pkg = self
        while pkg is not None and not isinstance(pkg, Package):
            pkg = pkg.owner
        if pkg is None:
            return ''
        t = pkg.name.split('_')
        t = [n for n in t if n.lower() not in ['cisco', 'ios', 'xr']]
        if t:
            sub = t[0].lower()
            if iskeyword(sub):
                sub = '%s_' % sub
            sub = '.%s' % sub
        else:
            sub = ''
        py_meta_mod_name = 'ydk.models%s._meta' % sub
        return py_meta_mod_name

    def fqn(self):
        ''' get the Fully Qualified Name '''
        names = []
        element = self
        while element is not None:
            if isinstance(element, Deviation):
                element = element.owner
            names.append(element.name)
            element = element.owner
        return '.'.join(reversed(names))

    def qn(self):
        ''' get the qualified name , name sans
        package name '''
        names = []
        element = self
        while element is not None and not isinstance(element, Package):
            if isinstance(element, Deviation):
                element = element.owner
            names.append(element.name)
            element = element.owner
        return '.'.join(reversed(names))


class Package(NamedElement):

    """
        Represents a Package in the API
    """

    def __init__(self):
        super(Package, self).__init__()
        self._stmt = None

    def qn(self):
        """ Return the qualified name """
        return self.name

    @property
    def stmt(self):
        """ Return the `pyang.statements.Statement` associated
            with this package. This is usually a module statement.
        """
        return self._stmt

    @stmt.setter
    def stmt(self, stmt):
        name = stmt.arg.replace('-', '_')
        if iskeyword(name):
            name = '%s_' % name
        if name[0] == '_':
            name = 'y%s' % name

        self.name = name
        self._stmt = stmt
        desc = stmt.search_one('description')
        if desc is not None:
            self.comment = desc.arg

    def imported_types(self):
        """
            Returns a list of all types imported by elements in
            this package.
        """
        imported_types = []
        for clazz in [c for c in self.owned_elements if isinstance(c, Class)]:
            imported_types.extend(clazz.imported_types())
        return imported_types

    def __eq__(self, other):
        if isinstance(other, self.__class__):
            return self.name == other.name
        else:
            return False


class DataType(NamedElement):

    """
        Represents a DataType
    """

    def __init__(self):
        super(DataType, self).__init__()


class Class(NamedElement):

    """
       Represents a Class in the api. 
    """

    def __init__(self):
        super(Class, self).__init__()
        self._stmt = None
        self._extends = []
        self._module = None

    def is_grouping_contribution(self):
        ''' Returns true if this Class is either a grouping class or is defined
        within a grouping class '''
        if self.is_grouping():
            return True

        owner = self.owner
        while owner is not None and isinstance(owner, Class):
            if owner.is_grouping():
                return True
            owner = owner.owner
        return False

    @property
    def extends(self):
        """ Returns the immediate super classes of this class. """
        if self.is_identity():
            base = []
            base_stmt = self.stmt.search_one('base')
            if base_stmt is not None and hasattr(base_stmt, 'i_identity'):
                base_identity = base_stmt.i_identity
                if hasattr(base_identity, 'i_class'):
                    base.append(base_identity.i_class)
            return base
        else:
            return self._extends

    def is_identity(self):
        """ Returns True if this is a class for a YANG identity. """
        return self._stmt.keyword == 'identity'

    def is_grouping(self):
        """ Returns True if this is a class for a YANG grouping. """
        return self._stmt.keyword == 'grouping'

    def is_rpc(self):
        return self._stmt.keyword == 'rpc'

    def all_owned_elements(self):
        """ Returns all the owned_element of this class and its super classes."""
        all_owned_elements = []
        for super_class in self.extends:
            all_owned_elements.extend(super_class.all_owned_elements())
        all_owned_elements.extend(self.owned_elements)
        return all_owned_elements

    def properties(self):
        """ Returns the properties defined by this class. """
        return get_properties(self.owned_elements)

    def get_package(self):
        """ Returns the Package that contains this Class. """
        if self.owner is None:
            return None
        if isinstance(self.owner, Package):
            return self.owner
        else:
            if hasattr(self.owner, 'get_package'):
                return self.owner.get_package()

    def imported_types(self):
        """ Returns all types that are referenced in this Class that are not
        from the same package as this Class."""

        imported_types = []
        package = self.get_package()

        # look at the super classes
        for super_class in self.extends:
            if super_class.get_package() != package:
                if super_class not in imported_types:
                    imported_types.append(super_class)

        for p in self.properties():
            prop_type = p.property_type
            if isinstance(prop_type, Class) or isinstance(prop_type, Enum) or isinstance(prop_type, Bits):
                if prop_type.get_package() != package:
                    if prop_type not in imported_types:
                        imported_types.append(prop_type)

        # do this for nested classes too
        for nested_class in [clazz for clazz in self.owned_elements if isinstance(clazz, Class)]:
            imported_types.extend(
                [c for c in nested_class.imported_types() if not c in imported_types])

        return imported_types

    def get_dependent_siblings(self):
        ''' This will return all types that are referenced by this Class
        or nested Classes that are at the same level as this type within the package and are
        used as super types .

        This is useful to determine which type needs to be printed
        before declaring this type in languages that do not support
        forward referencing like Python '''
        classes_at_same_level = []
        classes_at_same_level.extend(
            [c for c in self.owner.owned_elements if isinstance(c, Class) and c is not self])
        dependent_siblings = []
        package = self.get_package()

        def _walk_supers(clazz):
            for super_class in clazz.extends:
                if super_class.get_package() == package and super_class in classes_at_same_level:
                    if super_class not in dependent_siblings:
                        dependent_siblings.append(super_class)
                _walk_supers(super_class)

        def _walk_nested_classes(clazz):
            for nested_class in [c for c in clazz.owned_elements if isinstance(c, Class)]:
                _walk_supers(nested_class)
                _walk_nested_classes(nested_class)

        _walk_supers(self)
        _walk_nested_classes(self)

        return dependent_siblings

    def is_config(self):
        """
            Returns True if an instance of this Class represents config data.
        """
        if hasattr(self.stmt, 'i_config'):
            return self.stmt.i_config
        elif isinstance(self.owner, Class):
            return self.owner.is_config
        else:
            return True

    @property
    def stmt(self):
        """
            Returns the `pyang.statements.Statement` instance associated with this Class.
        """
        return self._stmt

    @property
    def module(self):
        """
            Returns the module `pyang.statements.Statement` that this Class was derived from.
        """
        return self._module

    @stmt.setter
    def stmt(self, stmt):
        name = escape_name(stmt.arg)
        if stmt.keyword == 'grouping':
            name = '%s_Grouping' % camel_case(name)
        elif stmt.keyword == 'identity':
            name = '%s_Identity' % camel_case(name)
        elif stmt.keyword == 'rpc':
            name = camel_case(name) + 'Rpc'
        else:
            name = camel_case(name)
        if iskeyword(name):
            name = '%s_' % name
        self.name = name

        if self.name.startswith('_'):
            self.name = '%s%s' % ('Y', name)

        self._stmt = stmt
        desc = stmt.search_one('description')
        if desc is not None:
            self.comment = desc.arg

        if stmt.keyword == 'module':
            self._module = stmt
        else:
            self._module = stmt.i_module

    def get_key_props(self):
        """ Returns a list of the properties of this class that are keys
        of a YANG list. """

        key_props = []

        if self.stmt.keyword == 'list':
            if hasattr(self.stmt, 'i_key'):
                key_stmts = self.stmt.i_key
                # do not use #properties here
                for prop in [p for p in self.owned_elements if isinstance(p, Property)]:
                    if prop.stmt in key_stmts:
                        key_props.append(prop)
        return sorted(key_props, key=lambda k: k.name)

    def __eq__(self, other):
        if isinstance(other, self.__class__):
            return self.fqn() == other.fqn()
        else:
            return False


class AnyXml(NamedElement):

    """
        Represents an AnyXml element.
    """

    def __init__(self):
        super(AnyXml, self).__init__()
        self._stmt = None

    @property
    def stmt(self):
        """ Returns the `pyang.statements.Statement` instance associated with this AnyXml instance."""
        return self._stmt

    @stmt.setter
    def stmt(self, stmt):
        self._stmt = stmt
        desc = stmt.search_one('description')
        if desc is not None:
            self.comment = desc.arg

    def properties(self):
        return get_properties(self.owned_elements)


class Bits(DataType):

    """
        A DataType representing the bits type in YANG.
    """

    def __init__(self):
        super(DataType, self).__init__()
        self._stmt = None
        self._dictionary = None
        self._pos_map = None

    def __eq__(self, other):
        if isinstance(other, self.__class__):
            return self.fqn() == other.fqn()
        else:
            return False

    def get_package(self):
        """ Returns the Package for this DataType. """
        if self.owner is None:
            return None
        if isinstance(self.owner, Package):
            return self.owner
        else:
            return self.owner.get_package()

    @property
    def stmt(self):
        return self._stmt

    @stmt.setter
    def stmt(self, stmt):
        self._stmt = stmt
        self._dictionary = {}
        self._pos_map = {}
        # the name of the enumeration is derived either from the typedef
        # or the leaf under which it is defined
        leaf_or_typedef = stmt
        while leaf_or_typedef.parent is not None and not leaf_or_typedef.keyword in ('leaf', 'leaf-list', 'typedef'):
            leaf_or_typedef = leaf_or_typedef.parent

        name = '%s_Bits' % camel_case(leaf_or_typedef.arg)
        if iskeyword(name):
            name = '%s_' % name
        self.name = name

        desc = stmt.search_one('description')
        if desc is not None:
            self.comment = desc.arg
        else:
            desc = leaf_or_typedef.search_one('description')
            if desc is not None:
                self.comment = desc.arg
        for bit_stmt in stmt.search('bit'):
            self._dictionary[bit_stmt.arg] = False
            pos_stmt = bit_stmt.search_one('position')
            if pos_stmt is not None:
                self._pos_map[bit_stmt.arg] = pos_stmt.arg


class Property(NamedElement):

    """ Represents an attribute or reference of a Class.
    """

    def __init__(self):
        super(Property, self).__init__()
        self._stmt = None
        self.is_static = False
        self.featuring_classifiers = []
        self.read_only = False
        self.is_many = False
        # self.property_type = None
        self.default_value = None
        self.visibility = '+'
        self.id = False
        self.ordered = False
        self.unique = False
        self._property_type = None
        self.max_elements = None
        self.min_elements = None

    def is_key(self):
        """ Returns True if this property represents a key of a YANG list."""
        if isinstance(self.owner, Class):
            return self in self.owner.get_key_props()
        return False

    @property
    def stmt(self):
        return self._stmt

    @stmt.setter
    def stmt(self, stmt):
        self._stmt = stmt
        name = snake_case(stmt.arg)
        if iskeyword(name):
            name = '%s_' % name
        self.name = name

        if self.name.startswith('_'):
            self.name = '%s%s' % ('y', name)
        if stmt.keyword in ['leaf-list', 'list']:
            self.is_many = True
        desc = stmt.search_one('description')
        if desc is not None:
            self.comment = desc.arg

        max_elem_stmt = stmt.search_one('max-elements')
        min_elem_stmt = stmt.search_one('min-elements')
        if max_elem_stmt:
            self.max_elements = max_elem_stmt.arg
        if min_elem_stmt:
            self.min_elements = min_elem_stmt.arg

    @property
    def property_type(self):
        """ Returns the type of this property. """
        if self._property_type is not None:
            return self._property_type

        if self._stmt is None:
            return None

        if self._stmt.keyword in ['leaf', 'leaf-list']:
            type_stmt = self._stmt.search_one('type')
            return type_stmt.i_type_spec
        else:
            return None

    @property_type.setter
    def property_type(self, property_type):
        self._property_type = property_type


class Enum(DataType):

    """ Represents an enumeration. """

    def __init__(self):
        super(Enum, self).__init__()
        self._stmt = None
        self.literals = []

    def get_package(self):
        """ Returns the Package that this enum is found in. """
        if self.owner is None:
            return None
        if isinstance(self.owner, Package):
            return self.owner
        else:
            return self.owner.get_package()

    @property
    def stmt(self):
        return self._stmt

    @stmt.setter
    def stmt(self, stmt):
        self._stmt = stmt
        # the name of the numeration is derived either from the typedef
        # or the leaf under which it is defined
        leaf_or_typedef = stmt
        while leaf_or_typedef.parent is not None and not leaf_or_typedef.keyword in ('leaf', 'leaf-list', 'typedef'):
            leaf_or_typedef = leaf_or_typedef.parent

        name = '%sEnum' % camel_case(escape_name(leaf_or_typedef.arg))
        if iskeyword(name):
            name = '%s_' % name

        if name[0] == '_':
            name = 'Y%s' % name

        self.name = name

        desc = stmt.search_one('description')
        if desc is not None:
            self.comment = desc.arg
        else:
            desc = leaf_or_typedef.search_one('description')
            if desc is not None:
                self.comment = desc.arg
        for enum_stmt in stmt.search('enum'):
            literal = EnumLiteral()
            literal.stmt = enum_stmt
            self.literals.append(literal)


class EnumLiteral(NamedElement):

    """ Represents an enumeration literal. """

    def __init__(self):
        super(EnumLiteral, self).__init__()
        self._stmt = None
        self.value = None

    @property
    def stmt(self):
        return self._stmt

    @stmt.setter
    def stmt(self, stmt):
        self._stmt = stmt
        self.name = stmt.arg.upper().replace('-', '_')
        self.name = self.name.replace('+', '__PLUS__')
        self.name = self.name.replace('/', '__FWD_SLASH__')
        self.name = self.name.replace('\\', '__BACK_SLASH__')
        self.name = self.name.replace('.', '__DOT__')
        self.name = self.name.replace('*', '__STAR__')
        self.name = self.name.replace('$', '__DOLLAR__')
        self.name = self.name.replace('@', '__AT__')
        self.name = self.name.replace('#', '__POUND__')
        self.name = self.name.replace('^', '__CARET__')
        self.name = self.name.replace('&', '__AMPERSAND__')
        self.name = self.name.replace('(', '__LPAREN__')
        self.name = self.name.replace(')', '__RPAREN__')
        self.name = self.name.replace('=', '__EQUALS__')
        self.name = self.name.replace('{', '__LCURLY__')
        self.name = self.name.replace('}', '__RCURLY__')
        self.name = self.name.replace("'", '__SQUOTE__')
        self.name = self.name.replace('"', '__DQUOTE__')
        self.name = self.name.replace('<', '__GREATER_THAN__')
        self.name = self.name.replace('>', '__LESS_THAN__')
        self.name = self.name.replace(',', '__COMMA__')
        self.name = self.name.replace(':', '__COLON__')
        self.name = self.name.replace('?', '__QUESTION__')
        self.name = self.name.replace('!', '__BANG__')
        self.name = self.name.replace(';', '__SEMICOLON__')

        if self.name[0:1].isdigit():
            self.name = 'Y_%s' % self.name

        if self.name[0] == '_':
            self.name = 'Y%s' % self.name

        self.value = stmt.i_value
        desc = stmt.search_one('description')
        if desc is not None:
            self.comment = desc.arg


def _get_type_stmt(stmt, typeSpec):
    if stmt.keyword == 'type':
        type_stmt = stmt
    else:
        type_stmt = stmt.search_one('type')

    if hasattr(type_stmt, 'i_typedef') and type_stmt.i_typedef is not None:
        typedef_stmt = type_stmt.i_typedef
        return _get_type_stmt(typedef_stmt, typeSpec)
    elif hasattr(type_stmt, 'i_type_spec'):
        type_spec = type_stmt.i_type_spec
        while isinstance(type_spec, PathTypeSpec):
            type_stmt = type_spec.i_target_node.search_one('type')
            type_spec = type_stmt.i_type_spec
            if hasattr(type_stmt, 'i_typedef') and type_stmt.i_typedef is not None:
                typedef_stmt = type_stmt.i_typedef
                return _get_type_stmt(typedef_stmt, typeSpec)

        if isinstance(type_spec, typeSpec):
            return type_stmt
        else:
            return None
    else:
        return None

_get_enum_type_stmt = lambda stmt: _get_type_stmt(stmt, EnumerationTypeSpec)
_get_identity_ref_type_stmt = lambda stmt: _get_type_stmt(
    stmt, IdentityrefTypeSpec)
_get_bits_type_stmt = lambda stmt: _get_type_stmt(stmt, BitsTypeSpec)
_get_union_type_stmt = lambda stmt: _get_type_stmt(stmt, UnionTypeSpec)


def _add_enums_and_bits(s, pe):
    enum_type_stmt = _get_enum_type_stmt(s)
    bits_type_stmt = _get_bits_type_stmt(s)
    union_type_stmt = _get_union_type_stmt(s)

    if enum_type_stmt is not None:
        enum_class = Enum()
        enum_class.stmt = enum_type_stmt
        enum_type_stmt.parent.i_enum = enum_class
        enum_type_stmt.i_enum = enum_class
        pe.owned_elements.append(enum_class)
        enum_class.owner = pe

    if bits_type_stmt is not None:
        bits_class = Bits()
        bits_class.stmt = bits_type_stmt
        bits_type_stmt.parent.i_bits = bits_class
        bits_type_stmt.i_bits = bits_class
        pe.owned_elements.append(bits_class)
        bits_class.owner = pe

    if union_type_stmt is not None:
        # need to process the type stmts under the union
        for contained_type in union_type_stmt.i_type_spec.types:
            _add_enums_and_bits(contained_type, pe)


def resolve_expanded_cross_references(element):
    """ Resolves cross references in the api_model Elements.

       Prerequisite before calling this function is that 
       the 'pyang.statements.Statement' tree for the list of modules
       must have their associated api_model Elements set in the i_class
       i_package, i_enum variables. These variables are used to resolve
       the cross references.

      :param `NamedElement` element :- The element whose references need to be resolved.
      :raise `common.YdkGenException` If cross resolution failed. 
    """
    if isinstance(element, Property):
        enum_type = _get_enum_type_stmt(element.stmt)
        if enum_type is not None and not isinstance(element.property_type, Enum):
            if not hasattr(enum_type.parent, 'i_enum'):
                # case where the type is a leafref pointing to a leaf with an
                # embedded enum definition
                if hasattr(enum_type.parent, 'i_property') and isinstance(enum_type.parent.i_property.property_type, Enum):
                    element.property_type = enum_type.parent.i_property.property_type
                else:
                    raise YdkGenException(
                        'Cross resolution of enum failed for ' + element.fqn())
            else:
                element.property_type = enum_type.parent.i_enum

        else:
            # check for identity_ref's
            identity_ref_type = _get_identity_ref_type_stmt(element.stmt)
            if identity_ref_type is not None:
                if not hasattr(identity_ref_type.i_type_spec.base.i_identity, 'i_class'):
                    raise YdkGenException(
                        'Cross resolution of identity class failed for ' + element.fqn())
                element.property_type = identity_ref_type.i_type_spec.base.i_identity.i_class
            else:
                # check for bits
                bits_ref_type = _get_bits_type_stmt(element.stmt)
                if bits_ref_type is not None and not isinstance(element.property_type, Bits):
                    if not hasattr(bits_ref_type.parent, 'i_bits'):
                        raise YdkGenException(
                            'Cross resolution of bits failed for ' + element.fqn())

                    element.property_type = bits_ref_type.parent.i_bits

    if hasattr(element, 'owned_elements'):
        for owned_element in element.owned_elements:
            resolve_expanded_cross_references(owned_element)

def add_to_deviation_package(target, parent_element, deviation_packages):
    i_deviation = target.i_deviation
    for d_type in i_deviation:
        d_obj = Deviation()
        d_obj.d_type = d_type
        d_obj.d_target = target
        d_obj.owner = parent_element
        d_obj.d_stmts = []
        for (d_module, d_stmt) in i_deviation[d_type]:
            d_module_name = d_module.arg
            target_package = [p for p in deviation_packages if p.stmt.arg == d_module_name][0]
            d_obj.d_stmts.append(d_stmt)
            if d_stmt.keyword == 'type':
                d_obj.d_target = target.copy()
                d_target = d_obj.d_target
                idx = d_target.substmts.index(d_target.search_one('type'))
                d_target.substmts[idx] = d_stmt
                add_leaf_leaflist_prop(d_target, d_obj)
            if d_obj not in target_package.owned_elements:
                target_package.owned_elements.append(d_obj)

def add_leaf_leaflist_prop(stmt, parent_element):
    prop = Property()
    stmt.i_property = prop
    prop.stmt = stmt
    parent_element.owned_elements.append(prop)
    prop.owner = parent_element
    # for inlined enum types where leaf { type enumeration {
    enum_type = _get_enum_type_stmt(stmt)
    bits_type = _get_bits_type_stmt(stmt)
    union_type = _get_union_type_stmt(stmt)
    # if the type statement is totally self contained
    # then we need to extract this type
    if enum_type is not None and enum_type == stmt.search_one('type'):
            # we have to create the enum
            enum_class = Enum()
            enum_class.stmt = enum_type
            parent_element.owned_elements.append(enum_class)
            enum_class.owner = parent_element
            prop.property_type = enum_class
    elif bits_type is not None and bits_type == stmt.search_one('type'):
            # we have to create the specific subclass of FixedBitsDict
            bits_class = Bits()
            bits_class.stmt = bits_type
            parent_element.owned_elements.append(bits_class)
            bits_class.owner = parent_element
            prop.property_type = bits_class
    elif union_type is not None and union_type == stmt.search_one('type'):
        def _add_union_type(union_type_stmt, parent_element):
            for contained_type in union_type_stmt.i_type_spec.types:
                contained_enum_type = _get_enum_type_stmt(contained_type)
                contained_bits_type = _get_bits_type_stmt(contained_type)
                contained_union_type = _get_union_type_stmt(contained_type)

                if contained_enum_type is not None and contained_enum_type == contained_type:
                    enum_class = Enum()
                    enum_class.stmt = contained_enum_type
                    parent_element.owned_elements.append(enum_class)
                    enum_class.owner = parent_element
                    contained_enum_type.i_enum = enum_class

                if contained_bits_type is not None and contained_bits_type == contained_type:
                    bits_class = Bits()
                    bits_class.stmt = contained_bits_type
                    parent_element.owned_elements.append(bits_class)
                    bits_class.owner = parent_element
                    contained_bits_type.i_bits = bits_class

                if contained_union_type is not None and contained_union_type == contained_type:
                    _add_union_type(contained_union_type, parent_element)

        # is this embedded ?
        if union_type == stmt.search_one('type'):
            # we need to check for the types under the union to see if
            # any of them need to be handled differently
            _add_union_type(union_type, parent_element)


def create_expanded_api_model(stmt, parent_element, deviation_packages):
    """
        Converts the stmt to an Element in the api_model according
        to the expanded code generation algorithm.

        The expanded code generation algorithm uses the tree view of the YANG models
        to generate the API. For each data node in the YANG model a corresponding Class
        element will be created.

        In the first pass Elements that encapsulate the references are created
        this is done for all the stmts we are interested
        after this is done, resolve cross references is called on all the elements
        to resolve all the cross references (examples include
        super classes extends field in a class)

        :param `pyang.statements.Statement` stmt The statement to convert
        :param  `Element` The parent element.
        :param list of 'Package' The deviation packages.
    """

    # process typedefs first so that they are resolved
    # when we have to use them
    element = parent_element

    # identities
    if hasattr(stmt, 'i_identities'):
        for identity_stmt in stmt.i_identities.values():
            identity_class = Class()
            identity_class.stmt = identity_stmt
            identity_class.owner = parent_element
            parent_element.owned_elements.append(identity_class)
            identity_stmt.i_class = identity_class

    if hasattr(stmt, 'i_typedefs'):
        for typedef_stmt_name in stmt.i_typedefs:
            typedef_stmt = stmt.i_typedefs[typedef_stmt_name]
            _add_enums_and_bits(typedef_stmt, parent_element)

    if stmt.keyword == 'module':
        pass

    elif stmt.keyword == 'container' or stmt.keyword == 'list' or stmt.keyword == 'rpc' or stmt.keyword == 'input' or stmt.keyword == 'output':
        if (stmt.keyword == 'input' or stmt.keyword == 'output') and len(stmt.substmts) == 0:
            pass
        else:
            clazz = Class()
            stmt.i_class = clazz
            clazz.stmt = stmt
            parent_element.owned_elements.append(clazz)
            clazz.owner = parent_element
            element = clazz

            if not isinstance(parent_element, Package):
                # create a property along with the class
                prop = Property()
                stmt.i_property = prop
                prop.stmt = stmt
                prop.property_type = clazz
                parent_element.owned_elements.append(prop)
                prop.owner = parent_element

    elif stmt.keyword == 'leaf' or stmt.keyword == 'leaf-list' or stmt.keyword == 'anyxml':
        add_leaf_leaflist_prop(stmt, parent_element)

    if hasattr(stmt, 'i_deviation'):
        add_to_deviation_package(stmt, parent_element, deviation_packages)
    
    # walk the children
    if hasattr(stmt, 'i_children'):
        sanitize_namespace(stmt)

        chs = [child for child in stmt.i_children
               if child.keyword in statements.data_definition_keywords
               or child.keyword == 'case'
               or child.keyword == 'rpc'
               or child.keyword == 'input'
               or child.keyword == 'output']
        for child_stmt in chs:
            create_expanded_api_model(child_stmt, element, deviation_packages)


def sanitize_namespace(stmt):
    """
        Detects if there is a name clash for a statement
        under a given node. 

        Note in the expanded tree algorithm
        augments show up under the parent node and their namespaces
        are not taken into account when the corresponding api_model.Element
        is generated. By sanitizing the namespace we avoid name collisions
        that can occur because of this situations

        :param `pyang.statements.Statement` The stmt to sanitize.
    """
    def _get_num_clashes(i_children):
        stmts = [stmt for stmt in i_children]
        clashes = []

        while len(stmts) > 0:
            stmt = stmts.pop()
            if stmt.arg in [s.arg for s in stmts] and stmt.arg not in clashes:
                clashes.append(stmt.arg)

        return clashes

    def _kill_clash(clash, i_children):
        for stmt in i_children:
            if clash == stmt.arg:
                old_arg = stmt.arg
                new_arg = '%s_%s' % (stmt.top.arg, old_arg)
                stmt.arg = new_arg

    clashes = _get_num_clashes(stmt.i_children)

    if len(clashes) > 0:
        for clash in clashes:
            _kill_clash(clash, stmt.i_children)


def generate_expanded_api_model(modules):
    """
        Generates and returns the packages for the list of modules supplied.

        This is the function that converts the list of pyang modules to 
        Packages in the api_model

        :param list of `pyang.statements.Statement`

    """

    d_modules = [module for module in modules if hasattr(module, 'is_deviation_module')]
    modules = [module for module in modules if not hasattr(module, 'is_deviation_module')]
    only_modules = [module for module in modules if module.keyword == 'module']
    packages = []
    deviation_packages = []

    for module in d_modules:
        package = Package()
        module.i_package = package
        package.stmt = module
        package.is_deviation = True
        deviation_packages.append(package)

    for module in only_modules:
        package = Package()
        module.i_package = package
        package.stmt = module
        create_expanded_api_model(module, package, deviation_packages)
        packages.append(package)

    packages.extend(deviation_packages)
    # resolve cross references
    for package in packages:
        resolve_expanded_cross_references(package)
    return packages


def resolve_grouping_class_cross_references(element):
    """
        Resolve cross references in the grouping as class code generation
        strategy.

        :param `api_model.Element` The model element whose cross references have to be
                resolved.

        :raise `common.YdkGenException' if cross resolution of references failed.

    """
    if isinstance(element, Class) and not element.is_identity():
        uses_stmts = element.stmt.search('uses')

        # set the super classes or the extend property
        groupings_used = []
        for uses in uses_stmts:
            groupings_used.append(uses.i_grouping)

        extends = []
        for grouping_stmt in groupings_used:
            if grouping_stmt.i_class is None:
                raise YdkGenException(
                    'Unresolved grouping class for ' + element.fqn())
            extends.append(grouping_stmt.i_class)
        element._extends = extends
    if isinstance(element, Property):
        enum_type = _get_enum_type_stmt(element.stmt)
        if enum_type is not None and not isinstance(element.property_type, Enum):
            if not hasattr(enum_type.parent, 'i_enum'):
                raise YdkGenException(
                    'Cross resolution of enum failed for ' + element.fqn())

            element.property_type = enum_type.parent.i_enum

        else:
            # check for identity_ref's
            identity_ref_type = _get_identity_ref_type_stmt(element.stmt)
            if identity_ref_type is not None:
                element.property_type = identity_ref_type.i_type_spec.base.i_identity.i_class
            else:
                # check for bits
                bits_ref_type = _get_bits_type_stmt(element.stmt)
                if bits_ref_type is not None and not isinstance(element.property_type, Bits):
                    if not hasattr(bits_ref_type.parent, 'i_bits'):
                        raise YdkGenException(
                            'Cross resolution of bits failed for ' + element.fqn())

                    element.property_type = bits_ref_type.i_bits

    if hasattr(element, 'owned_elements'):
        for owned_element in element.owned_elements:
            resolve_grouping_class_cross_references(owned_element)


def create_grouping_class_api_model(stmt, parent_element):
    """
        Converts the stmt to an Element in the api_model according
        to the grouping as class algorithm.

        In the grouping as class code generations strategy a grouping in YANG
        is converted to a class. Every class that represents a container or a list
        or a grouping that has a uses statement in it , inherits from the grouping class that
        corresponds to the grouping in the uses statement.

        for example

        grouping abc {                        class Abc_Grouping(object):...
        ....
        }

        container A {                         class A(Abc_Grouping): ...
            uses abc;
        }

        In the first pass Elements that encapsulate the references are created
        this is done for all the stmts we are interested
        after this is done, resolve cross references is called on all the elements
        to resolve all the cross references (examples include
        super classes extends field in a class)

        :param `pyang.statements.Statement` stmt The statement to convert
        :param  `Element` The parent element.
    """
    # process typedefs first so that they are resolved
    # when we have to use them
    element = parent_element

    # identities
    if hasattr(stmt, 'i_identities'):
        for identity_stmt in stmt.i_identities.values():
            identity_class = Class()
            identity_class.stmt = identity_stmt
            identity_class.owner = parent_element
            parent_element.owned_elements.append(identity_class)
            identity_stmt.i_class = identity_class

    if hasattr(stmt, 'i_typedefs'):
        for typedef_stmt_name in stmt.i_typedefs:
            typedef_stmt = stmt.i_typedefs[typedef_stmt_name]
            _add_enums_and_bits(typedef_stmt, parent_element)

    # walk the groupings first
    if hasattr(stmt, 'i_groupings'):
        for grouping_name in stmt.i_groupings:
            create_grouping_class_api_model(
                stmt.i_groupings[grouping_name], element)

    if stmt.keyword == 'grouping':
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

        if not isinstance(parent_element, Package):
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
        # for inlined enum types where leaf { type enumeration {
        enum_type = _get_enum_type_stmt(stmt)
        bits_type = _get_bits_type_stmt(stmt)
        if enum_type is not None:
            if enum_type == stmt.search_one('type'):
                # we have to create the enum
                enum_class = Enum()
                enum_class.stmt = enum_type
                parent_element.owned_elements.append(enum_class)
                enum_class.owner = parent_element
                prop.property_type = enum_class
                enum_type.parent.i_enum = enum_class
        elif bits_type is not None:
            if bits_type == stmt.search_one('type'):
                # we have to create the specific subclass of FixedBitsDict
                bits_class = Bits()
                bits_class.stmt = bits_type
                parent_element.owned_elements.append(bits_class)
                bits_class.owner = parent_element
                prop.property_type = bits_class

    # walk the children
    if hasattr(stmt, 'i_children'):
        grouping_stmt_names = []

        if stmt.keyword != 'module':
            uses_stmts = stmt.search('uses')
            groupings_used = []
            for uses_stmt in uses_stmts:
                groupings_used.append(uses_stmt.i_grouping)

            for grouping in groupings_used:
                grouping_stmt_names.extend(
                    [s.arg for s in grouping.i_children])

        chs = [ch for ch in stmt.i_children
               if ch.keyword in statements.data_definition_keywords and ch.arg not in grouping_stmt_names]
        for child_stmt in chs:
            create_grouping_class_api_model(child_stmt, element)


def generate_grouping_class_api_model(modules):
    """
        Generates and returns the packages for the list of modules supplied.

        This is the function that converts the list of pyang modules to 
        Packages in the api_model for the grouping as classes strategy.

        :param list of `pyang.statements.Statement`

    """
    only_modules = [m for m in modules if m.keyword == 'module']
    packages = []
    for m in only_modules:
        p = Package()
        m.i_package = p
        p.stmt = m
        create_grouping_class_api_model(m, p)
        packages.append(p)

    # resolve cross references
    for p in packages:
        resolve_grouping_class_cross_references(p)
    return packages


def get_properties(owned_elements):
    """ get all properties from the owned_elements. """
    props = []
    all_props = []
    all_props.extend([p for p in owned_elements if isinstance(p, Property)])

    # first get the key properties
    key_props = [p for p in all_props if p.is_key()]
    props.extend(sorted(key_props, key=lambda p: p.name))

    non_key_props = [p for p in all_props if not p.is_key()]
    props.extend(sorted(non_key_props, key=lambda p: p.name))

    return props
