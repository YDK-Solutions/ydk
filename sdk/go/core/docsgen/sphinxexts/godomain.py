# -*- coding: utf-8 -*-
"""
    sphinx.domains.go
    ~~~~~~~~~~~~~~~~~~~~~

    The Go domain.

    :copyright: Copyright 2007-2017 by the Sphinx team, see AUTHORS.
    :license: BSD, see LICENSE for details.
"""

import logging
import re

from six import iteritems

from docutils import nodes
from docutils.parsers.rst import Directive, directives

from sphinx import addnodes
from sphinx.roles import XRefRole
from sphinx.locale import l_, _
from sphinx.domains import Domain, ObjType, Index
from sphinx.directives import ObjectDescription
from sphinx.util.nodes import make_refnode
from sphinx.util.docfields import Field, GroupedField, TypedField

logger = logging.getLogger(__name__)


# REs for Go signatures
go_sig_re = re.compile(
    r'''^ ([\w.]*\.)?            # struct name(s)
          (\w+)  \s*             # thing name
          (?: \((.*)\)           # optional: arguments
           (?:\s* -> \s* (.*))?  #           return annotation
          )? $                   # and nothing more
          ''', re.VERBOSE)


def _pseudo_parse_arglist(signode, arglist):
    # type: (addnodes.desc_signature, unicode) -> None
    """"Parse" a list of arguments separated by commas.

    Arguments can have "optional" annotations given by enclosing them in
    brackets.  Currently, this will split at any comma, even if it's inside a
    string literal (e.g. default argument value).
    """
    paramlist = addnodes.desc_parameterlist()
    stack = [paramlist]
    try:
        for argument in arglist.split(','):
            argument = argument.strip()
            ends_open = ends_close = 0
            while argument.startswith('['):
                stack.append(addnodes.desc_optional())
                stack[-2] += stack[-1]
                argument = argument[1:].strip()
            while argument.startswith(']'):
                stack.pop()
                argument = argument[1:].strip()
            while argument.endswith(']') and not argument.endswith('[]'):
                ends_close += 1
                argument = argument[:-1].strip()
            while argument.endswith('['):
                ends_open += 1
                argument = argument[:-1].strip()
            if argument:
                stack[-1] += addnodes.desc_parameter(argument, argument)
            while ends_open:
                stack.append(addnodes.desc_optional())
                stack[-2] += stack[-1]
                ends_open -= 1
            while ends_close:
                stack.pop()
                ends_close -= 1
        if len(stack) != 1:
            raise IndexError
    except IndexError:
        # if there are too few or too many elements on the stack, just give up
        # and treat the whole argument list as one argument, discarding the
        # already partially populated paramlist node
        signode += addnodes.desc_parameterlist()
        signode[-1] += addnodes.desc_parameter(arglist, arglist)
    else:
        signode += paramlist


# This override allows our inline type specifiers to behave like :struct: link
# when it comes to handling "." and "~" prefixes.
class GoXrefMixin(object):
    def make_xref(self,
                  rolename,                  # type: unicode
                  domain,                    # type: unicode
                  target,                    # type: unicode
                  innernode=nodes.emphasis,  # type: nodes.Node
                  contnode=None,             # type: nodes.Node
                  env=None,                  # type: BuildEnvironment
                  ):
        # type: (...) -> nodes.Node

        result = super(GoXrefMixin, self).make_xref(rolename, domain, target,  # type: ignore
                                                    innernode, contnode, env)
        logger.warning(result)

        result['refspecific'] = True
        if target.startswith(('/', '~')):
            prefix, result['reftarget'] = target[0], target[1:]
            if prefix == '/':
                text = target[1:]
            elif prefix == '~':
                text = target.split('/')[-1]
            for node in result.traverse(nodes.Text):
                node.parent[node.parent.index(node)] = nodes.Text(text)
                break
        return result

    def make_xrefs(self,
                   rolename,                  # type: unicode
                   domain,                    # type: unicode
                   target,                    # type: unicode
                   innernode=nodes.emphasis,  # type: nodes.Node
                   contnode=None,             # type: nodes.Node
                   env=None,                  # type: BuildEnvironment
                   ):
        # type: (...) -> List[nodes.Node]
        delims = r'(\s*[\[\]\(\),](?:\s*or\s)?\s*|\s+or\s+)'
        delims_re = re.compile(delims)
        sub_targets = re.split(delims, target)

        split_contnode = bool(contnode and contnode.astext() == target)

        results = []
        for sub_target in filter(None, sub_targets):
            if split_contnode:
                contnode = nodes.Text(sub_target)

            if delims_re.match(sub_target):  # type: ignore
                results.append(contnode or innernode(sub_target, sub_target))
            else:
                results.append(self.make_xref(rolename, domain, sub_target,
                                              innernode, contnode, env))

        return results


class GoField(GoXrefMixin, Field):
    pass


class GoGroupedField(GoXrefMixin, GroupedField):
    pass


class GoTypedField(GoXrefMixin, TypedField):
    pass


class GoObject(ObjectDescription):
    """
    Description of a general Go object.

    :cvar allow_nesting: an object that allows for nested namespaces
    :vartype allow_nesting: bool
    """
    option_spec = {
        'noindex': directives.flag,
        'package': directives.unchanged,
        'annotation': directives.unchanged,
    }

    doc_field_types = [
        TypedField('parameter', label=l_('Parameters'),
                   names=('param', 'parameter', 'arg', 'argument'),
                   typerolename='type', typenames=('type',)),
        Field('returnvalue', label=l_('Returns'), has_arg=False,
              names=('returns', 'return')),
        GoField('returntype', label=l_('Return type'), has_arg=False,
              names=('rtype',)),
    ]

    allow_nesting = False

    @staticmethod
    def get_signature_prefix(sig):
        # type: (unicode) -> unicode
        """May return a prefix to put before the object name in the
        signature.
        """
        return ''

    @staticmethod
    def needs_arglist():
        # type: () -> bool
        """May return true if an empty argument list is to be generated even if
        the document contains none.
        """
        return False

    def handle_signature(self, sig, signode):
        # type: (unicode, addnodes.desc_signature) -> Tuple[unicode, unicode]
        """Transform a Go signature into RST nodes.

        Return (fully qualified name of the thing, structname if any).

        If inside a struct, the current struct name is handled intelligently:
        * it is stripped from the displayed name if present
        * it is added to the full name (return value) if not present
        """
        m = go_sig_re.match(sig)  # type: ignore
        if m is None:
            raise ValueError
        name_prefix, name, arglist, retann = m.groups()

        # determine package and struct name (if applicable), as well as full name
        pkgname = self.options.get(
            'package', self.env.ref_context.get('go:package'))
        structname = self.env.ref_context.get('go:struct')
        if structname:
            add_package = False
            if name_prefix and name_prefix.startswith(structname):
                fullname = name_prefix + name
                # struct name is given again in the signature
                name_prefix = name_prefix[len(structname):].lstrip('/')
            elif name_prefix:
                # struct name is given in the signature, but different
                # (shouldn't happen)
                fullname = structname + '/' + name_prefix + name
            else:
                # struct name is not given in the signature
                fullname = structname + '/' + name
        else:
            add_package = True
            if name_prefix:
                structname = name_prefix.rstrip('/')
                fullname = name_prefix + name
            else:
                structname = ''
                fullname = name

        signode['package'] = pkgname
        signode['struct'] = structname
        signode['fullname'] = fullname

        sig_prefix = self.get_signature_prefix(sig)

        if sig_prefix:
            signode += addnodes.desc_annotation(sig_prefix, sig_prefix)

        if name_prefix:
            signode += addnodes.desc_addname(name_prefix, name_prefix)
        # exceptions are a special case, since they are documented in the
        # 'exceptions' package.
        elif add_package and self.env.config.add_module_names:
            pkgname = self.options.get(
                'package', self.env.ref_context.get('go:package'))
            if pkgname and pkgname != 'exceptions':
                nodetext = pkgname + '/'
                signode += addnodes.desc_addname(nodetext, nodetext)

        anno = self.options.get('annotation')

        signode += addnodes.desc_name(name, name)
        if not arglist:
            if self.needs_arglist():
                # for callables, add an empty parameter list
                signode += addnodes.desc_parameterlist()
            if retann:
                signode += addnodes.desc_returns(retann, retann)
            if anno:
                signode += addnodes.desc_annotation(' ' + anno, ' ' + anno)
            return fullname, name_prefix

        _pseudo_parse_arglist(signode, arglist)
        if retann:
            signode += addnodes.desc_returns(retann, retann)
        if anno:
            signode += addnodes.desc_annotation(' ' + anno, ' ' + anno)
        return fullname, name_prefix

    def get_index_text(self, pkgname, name_struct):
        # type: (unicode, unicode) -> unicode
        """Return the text for the index entry of the object."""
        raise NotImplementedError('must be implemented in inherited struct')

    def add_target_and_index(self, name_struct, sig, signode):
        # type: (unicode, unicode, addnodes.desc_signature) -> None
        pkgname = self.options.get(
            'package', self.env.ref_context.get('go:package'))
        fullname = (pkgname and pkgname + '/' or '') + name_struct[0]
        # note target
        if fullname not in self.state.document.ids:
            signode['names'].append(fullname)
            signode['ids'].append(fullname)
            signode['first'] = (not self.names)
            self.state.document.note_explicit_target(signode)
            objects = self.env.domaindata['go']['objects']
            if fullname in objects:
                self.state_machine.reporter.warning(
                    'duplicate object description of %s, ' % fullname +
                    'other instance in ' +
                    self.env.doc2path(objects[fullname][0]) +
                    ', use :noindex: for one of them',
                    line=self.lineno)
            objects[fullname] = (self.env.docname, self.objtype)

        indextext = self.get_index_text(pkgname, name_struct)
        if indextext:
            self.indexnode['entries'].append(('single', indextext,
                                              fullname, '', None))

    def before_content(self):
        # type: () -> None
        """Handle object nesting before content

        :go:struct:`GoObject` represents Go language constructs. For
        constructs that are nestable, this method will
        build up a stack of the nesting heirarchy so that it can be later
        de-nested correctly, in :py:meth:`after_content`.

        For constructs that aren't nestable, the stack is bypassed, and instead
        only the most recent object is tracked. This object prefix name will be
        removed with :py:meth:`after_content`.
        """
        prefix = None
        if self.names:
            # fullname and name_prefix come from the `handle_signature` method.
            # fullname represents the full object name that is constructed using
            # object nesting and explicit prefixes. `name_prefix` is the
            # explicit prefix given in a signature
            (fullname, name_prefix) = self.names[-1]
            if self.allow_nesting:
                prefix = fullname
            elif name_prefix:
                prefix = name_prefix.strip('/')
        if prefix:
            self.env.ref_context['go:struct'] = prefix
            if self.allow_nesting:
                structs = self.env.ref_context.setdefault('go:structs', [])
                structs.append(prefix)

    def after_content(self):
        # type: () -> None
        """Handle object de-nesting after content

        If the object is nestable, removing the last nested prefix
        ends further nesting in the object.

        If this object is not a nestable object, the list should not
        be altered as we didn't affect the nesting levels in
        :py:meth:`before_content`.
        """
        structs = self.env.ref_context.setdefault('go:structs', [])
        if self.allow_nesting:
            try:
                structs.pop()
            except IndexError:
                pass
        self.env.ref_context['go:struct'] = (structs[-1] if len(structs) > 0
                                            else None)


class GoPackagelevel(GoObject):
    """
    Description of an object on package level (functions, data).
    """

    def needs_arglist(self):
        # type: () -> bool
        return self.objtype == 'function'

    def get_index_text(self, pkgname, name_struct):
        # type: (unicode, unicode) -> unicode
        if self.objtype == 'function':
            if not pkgname:
                return _('%s() (built-in function)') % name_struct[0]
            return _('%s() (in package %s)') % (name_struct[0], pkgname)
        elif self.objtype == 'data':
            if not pkgname:
                return _('%s (built-in variable)') % name_struct[0]
            return _('%s (in package %s)') % (name_struct[0], pkgname)
        else:
            return ''


class GoStructlike(GoObject):
    """
    Description of a struct-like object (structs, interfaces, exceptions).
    """

    allow_nesting = True

    def get_signature_prefix(self, sig):
        # type: (unicode) -> unicode
        return self.objtype + ' '

    def get_index_text(self, pkgname, name_struct):
        # type: (unicode, unicode) -> unicode
        if self.objtype == 'struct':
            if not pkgname:
                return _('%s (built-in struct)') % name_struct[0]
            return _('%s (struct in %s)') % (name_struct[0], pkgname)
        elif self.objtype == 'exception':
            return name_struct[0]
        else:
            return ''


class GoStructmember(GoObject):
    """
    Description of a struct member (functions, attributes).
    """

    def get_index_text(self, pkgname, name_struct):
        # type: (unicode, unicode) -> unicode
        name = name_struct[0]
        add_packages = self.env.config.add_module_names
        if self.objtype in ('function', 'func'):
            try:
                structname, funcname = name.rsplit('/', 1)
            except ValueError:
                if pkgname:
                    return _('%s() (in package %s)') % (name, pkgname)
                else:
                    return '%s()' % name
            if pkgname and add_packages:
                return _('%s() (%s/%s function)') % (funcname, pkgname, structname)
            else:
                return _('%s() (%s function)') % (funcname, structname)
        elif self.objtype == 'attribute':
            try:
                structname, attrname = name.rsplit('/', 1)
            except ValueError:
                if pkgname:
                    return _('%s (in package %s)') % (name, pkgname)
                else:
                    return name
            if pkgname and add_packages:
                return _('%s (%s/%s attribute)') % (attrname, pkgname, structname)
            else:
                return _('%s (%s attribute)') % (attrname, structname)
        else:
            return ''


class GoPackage(Directive):
    """
    Directive to mark description of a new package.
    """

    has_content = False
    required_arguments = 1
    optional_arguments = 0
    final_argument_whitespace = False
    option_spec = {
        'platform': lambda x: x,
        'synopsis': lambda x: x,
        'noindex': directives.flag,
        'deprecated': directives.flag,
    }

    def run(self):
        # type: () -> List[nodes.Node]
        env = self.state.document.settings.env
        pkgname = self.arguments[0].strip()
        noindex = 'noindex' in self.options
        env.ref_context['go:package'] = pkgname
        ret = []
        if not noindex:
            env.domaindata['go']['packages'][pkgname] = \
                (env.docname, self.options.get('synopsis', ''),
                 self.options.get('platform', ''), 'deprecated' in self.options)
            # make a duplicate entry in 'objects' to facilitate searching for
            # the package in GoDomain.find_obj()
            env.domaindata['go']['objects'][pkgname] = (env.docname, 'package')
            targetnode = nodes.target('', '', ids=['package-' + pkgname],
                                      ismod=True)
            self.state.document.note_explicit_target(targetnode)
            # the platform and synopsis aren't printed; in fact, they are only
            # used in the modindex currently
            ret.append(targetnode)
            indextext = _('%s (package)') % pkgname
            inode = addnodes.index(entries=[('single', indextext,
                                             'package-' + pkgname, '', None)])
            ret.append(inode)
        return ret


class GoCurrentPackage(Directive):
    """
    This directive is just to tell Sphinx that we're documenting
    stuff in package foo, but links to package foo won't lead here.
    """

    has_content = False
    required_arguments = 1
    optional_arguments = 0
    final_argument_whitespace = False
    option_spec = {}  # type: Dict

    def run(self):
        # type: () -> List[nodes.Node]
        env = self.state.document.settings.env
        pkgname = self.arguments[0].strip()
        if pkgname == 'None':
            env.ref_context.pop('go:package', None)
        else:
            env.ref_context['go:package'] = pkgname
        return []


class GoXRefRole(XRefRole):

    @staticmethod
    def process_link(env, refnode, has_explicit_title, title, target):
        # type: (BuildEnvironment, nodes.Node, bool, unicode, unicode) -> Tuple[unicode, unicode]  # NOQA
        refnode['go:package'] = env.ref_context.get('go:package')
        refnode['go:struct'] = env.ref_context.get('go:struct')
        if not has_explicit_title:
            title = title.lstrip('/')    # only has a meaning for the target
            target = target.lstrip('~')  # only has a meaning for the title
            # if the first character is a tilde, don't display the package/class
            # parts of the contents
            if title[0:1] == '~':
                title = title[1:]
                slash = title.rfind('/')
                if slash != -1:
                    title = title[slash + 1:]
        # if the first character is a slash, search more specific namespaces first
        # else search builtins first
        if target[0:1] == '/':
            target = target[1:]
            refnode['refspecific'] = True
        return title, target


class GoPackageIndex(Index):
    """
    Index subclass to provide the Go package index.
    """

    name = 'pkgindex'
    localname = l_('Go Package Index')
    shortname = l_('packages')

    def generate(self, docnames=None):
        # type: (Iterable[unicode]) -> Tuple[List[Tuple[unicode, List[List[Union[unicode, int]]]]], bool]  # NOQA
        content = {}  # type: Dict[unicode, List]
        # list of prefixes to ignore
        ignores = None  # type: List[unicode]
        ignores = self.domain.env.config['pkgindex_common_prefix']  # type: ignore
        ignores = sorted(ignores, key=len, reverse=True)
        # list of all packages, sorted by package name
        packages = sorted(iteritems(self.domain.data['packages']),
                         key=lambda x: x[0].lower())
        # sort out collapsable packages
        prev_pkgname = ''
        num_toplevels = 0
        for pkgname, (docname, synopsis, platforms, deprecated) in packages:
            if docnames and docname not in docnames:
                continue

            for ignore in ignores:
                if pkgname.startswith(ignore):
                    pkgname = pkgname[len(ignore):]
                    stripped = ignore
                    break
            else:
                stripped = ''

            # we stripped the whole package name?
            if not pkgname:
                pkgname, stripped = stripped, ''

            entries = content.setdefault(pkgname[0].lower(), [])

            package = pkgname.split('/')[0]
            if package != pkgname:
                # it's a subpackage
                if prev_pkgname == package:
                    # first subpackage - make parent a group head
                    if entries:
                        entries[-1][1] = 1
                elif not prev_pkgname.startswith(package):
                    # subpackage without parent in list, add dummy entry
                    entries.append([stripped + package, 1, '', '', '', '', ''])
                subtype = 2
            else:
                num_toplevels += 1
                subtype = 0

            qualifier = deprecated and _('Deprecated') or ''
            entries.append([stripped + pkgname, subtype, docname,
                            'package-' + stripped + pkgname, platforms,
                            qualifier, synopsis])
            prev_pkgname = pkgname

        # apply heuristics when to collapse modindex at page load:
        # only collapse if number of toplevel packages is larger than
        # number of subpackages
        collapse = len(packages) - num_toplevels < num_toplevels

        # sort by first letter
        sorted_content = sorted(iteritems(content))

        return sorted_content, collapse


class GoDomain(Domain):
    """Go language domain."""
    name = 'go'
    label = 'Golang'
    object_types = {
        'function':     ObjType(l_('function'),      'func', 'obj'),
        'data':         ObjType(l_('data'),          'data', 'obj'),
        'struct':       ObjType(l_('struct'),        'struct', 'exc', 'obj'),
        'exception':    ObjType(l_('exception'),     'exc', 'struct', 'obj'),
        'package':      ObjType(l_('package'),       'pkg', 'obj'),
    }  # type: Dict[unicode, ObjType]

    directives = {
        'function':        GoPackagelevel,
        'data':            GoPackagelevel,
        'struct':          GoStructlike,
        'exception':       GoStructlike,
        # 'function':        GoStructmember,
        # 'attribute':       GoStructmember,
        'package':         GoPackage,
        'currentpackage':  GoCurrentPackage,
    }
    roles = {
        'data':   GoXRefRole(),
        'exc':    GoXRefRole(),
        'func':   GoXRefRole(fix_parens=True),
        'struct': GoXRefRole(),
        # 'class': GoXRefRole(),
        # 'const': GoXRefRole(),
        # 'attr':  GoXRefRole(),
        # 'meth':  GoXRefRole(fix_parens=True),
        'pkg':    GoXRefRole(),
        'type':   GoXRefRole(),
        'obj':    GoXRefRole(),
    }
    initial_data = {
        'objects': {},  # fullname -> docname, objtype
        'packages': {},  # pkgname -> docname, synopsis, platform, deprecated
    }  # type: Dict[unicode, Dict[unicode, Tuple[Any]]]
    indices = [
        GoPackageIndex,
    ]

    def clear_doc(self, docname):
        # type: (unicode) -> None
        for fullname, (fn, _l) in list(self.data['objects'].items()):
            if fn == docname:
                del self.data['objects'][fullname]
        for pkgname, (fn, _x, _x, _x) in list(self.data['packages'].items()):
            if fn == docname:
                del self.data['packages'][pkgname]

    def merge_domaindata(self, docnames, otherdata):
        # type: (List[unicode], Dict) -> None
        # XXX check duplicates?
        for fullname, (fn, objtype) in otherdata['objects'].items():
            if fn in docnames:
                self.data['objects'][fullname] = (fn, objtype)
        for pkgname, data in otherdata['packages'].items():
            if data[0] in docnames:
                self.data['packages'][pkgname] = data

    def find_obj(self, env, pkgname, structname, name, typ, searchmode=0):
        # type: (BuildEnvironment, unicode, unicode, unicode, unicode, int) -> List[Tuple[unicode, Any]]  # NOQA
        """Find a Go object for "name", perhaps using the given package
        and/or structname.  Returns a list of (name, object entry) tuples.
        """
        # skip parens
        if name[-2:] == '()':
            name = name[:-2]

        if not name:
            return []

        objects = self.data['objects']
        matches = []  # type: List[Tuple[unicode, Any]]

        newname = None
        if searchmode == 1:
            if typ is None:
                objtypes = list(self.object_types)
            else:
                objtypes = self.objtypes_for_role(typ)
            if objtypes is not None:
                if pkgname and structname:
                    fullname = pkgname + '/' + structname + '/' + name
                    if fullname in objects and objects[fullname][1] in objtypes:
                        newname = fullname
                if not newname:
                    if pkgname and pkgname + '/' + name in objects and \
                       objects[pkgname + '/' + name][1] in objtypes:
                        newname = pkgname + '/' + name
                    elif name in objects and objects[name][1] in objtypes:
                        newname = name
                    else:
                        # "fuzzy" searching mode
                        searchname = '/' + name
                        matches = [(oname, objects[oname]) for oname in objects
                                   if oname.endswith(searchname) and
                                   objects[oname][1] in objtypes]
        else:
            # NOTE: searching for exact match, object type is not considered
            if name in objects:
                newname = name
            elif typ == 'pkg':
                # only exact matches allowed for packages
                return []
            elif structname and structname + '/' + name in objects:
                newname = structname + '/' + name
            elif pkgname and pkgname + '/' + name in objects:
                newname = pkgname + '/' + name
            elif pkgname and structname and \
                    pkgname + '/' + structname + '/' + name in objects:
                newname = pkgname + '/' + structname + '/' + name
            # special case: builtin exceptions have package "exceptions" set
            elif typ == 'exc' and '/' not in name and \
                    'exceptions.' + name in objects:
                newname = 'exceptions.' + name
            # special case: object methods
            elif typ in ('func', 'meth') and '/' not in name and \
                    'object.' + name in objects:
                newname = 'object.' + name
        if newname is not None:
            matches.append((newname, objects[newname]))
        return matches

    def resolve_xref(self, env, fromdocname, builder, typ, target, node, contnode):
        # type: (BuildEnvironment, unicode, Builder, unicode, unicode, nodes.Node, nodes.Node) -> nodes.Node  # NOQA
        pkgname = node.get('go:package')
        structname = node.get('go:struct')
        searchmode = node.hasattr('refspecific') and 1 or 0
        matches = self.find_obj(env, pkgname, structname, target,
                                typ, searchmode)
        if not matches:
            return None
        elif len(matches) > 1:
            logger.warning('more than one target found for cross-reference %r: %s',
                           target, ', '.join(match[0] for match in matches),
                           location=node)
        name, obj = matches[0]

        if obj[1] == 'package':
            return self._make_package_refnode(builder, fromdocname, name,
                                             contnode)
        else:
            return make_refnode(builder, fromdocname, obj[0], name,
                                contnode, name)

    def resolve_any_xref(self, env, fromdocname, builder, target, node, contnode):
        # type: (BuildEnvironment, unicode, Builder, unicode, nodes.Node, nodes.Node) -> List[Tuple[unicode, nodes.Node]]  # NOQA
        pkgname = node.get('go:package')
        structname = node.get('go:struct')
        results = []  # type: List[Tuple[unicode, nodes.Node]]

        # always search in "refspecific" mode with the :any: role
        matches = self.find_obj(env, pkgname, structname, target, None, 1)
        for name, obj in matches:
            if obj[1] == 'package':
                results.append(('go:pkg',
                                self._make_package_refnode(builder, fromdocname,
                                                          name, contnode)))
            else:
                results.append(('go:' + self.role_for_objtype(obj[1]),
                                make_refnode(builder, fromdocname, obj[0], name,
                                             contnode, name)))
        return results

    def _make_package_refnode(self, builder, fromdocname, name, contnode):
        # type: (Builder, unicode, unicode, nodes.Node) -> nodes.Node
        # get additional info for packages
        docname, synopsis, platform, deprecated = self.data['packages'][name]
        title = name
        if synopsis:
            title += ': ' + synopsis
        if deprecated:
            title += _(' (deprecated)')
        if platform:
            title += ' (' + platform + ')'
        return make_refnode(builder, fromdocname, docname,
                            'package-' + name, contnode, title)

    def get_objects(self):
        # type: () -> Iterator[Tuple[unicode, unicode, unicode, unicode, unicode, int]]
        for pkgname, info in iteritems(self.data['packages']):
            yield (pkgname, pkgname, 'package', info[0], 'package-' + pkgname, 0)
        for refname, (docname, typ) in iteritems(self.data['objects']):
            if typ != 'package':  # packages are already handled
                yield (refname, refname, typ, docname, refname, 1)

    @staticmethod
    def get_full_qualified_name(node):
        # type: (nodes.Node) -> unicode
        pkgname = node.get('go:package')
        structname = node.get('go:struct')
        target = node.get('reftarget')
        if target is None:
            return None
        else:
            return '/'.join(filter(None, [pkgname, structname, target]))


def setup(app):
    # type: (Sphinx) -> Dict[unicode, Any]
    app.add_domain(GoDomain)

    return {
        'version': 'builtin',
        'parallel_read_safe': True,
        'parallel_write_safe': True,
    }
