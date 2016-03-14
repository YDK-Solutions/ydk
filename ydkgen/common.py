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
 common.py 
 
 YANG model driven API, common definitions.
"""

#  ----------------------------------------------------------------
#  Generic lookups
# -----------------------------------------------------------------
yang_int = {
    'int8',
    'int16',
    'int32',
    'int64',
    'uint8',
    'uint16',
    'uint32',
    'uint64',
}

yang_int_ranges = {
    'int8': (-128, 127),
    'int16': (-32768, 32767),
    'int32': (-2147483648, 2147483647),
    'int64': (-9223372036854775808, 9223372036854775807),
    'uint8': (0, 255),
    'uint16': (0, 65535),
    'uint32': (0, 4294967295),
    'uint64': (0, 18446744073709551615),
}

yang_base_types = {
    'binary',
    'bits',
    'boolean',
    'decimal64',
    'empty',
    'identityref',
    'instance-identifier',
    'int8',
    'int16',
    'int32',
    'int64',
    'leafref',
    'string',
    'uint8',
    'uint16',
    'uint32',
    'uint64',
    # union, separate handling
    # enumeration, separate handling
}

container_nodes = {
    'module',
    'container',
    'choice',
    'case',
    'list',
    'augment',
    #    'grouping',
    'uses',
    'rpc',
    'input',
    'output',
}


class YdkGenException(Exception):

    """Exception raised when there is a problem in the generation.

        .. attribute:: msg
                      The message describing the error.

    """

    def __init__(self, msg):
        self.msg = msg


class PrintCtx(object):

    """
        Print Context.

        Used to encapsulate information needed by the printers.
    """

    def __init__(self):
        self.fd = None
        self.lvl = 0
        # internal
        self.all_classes = []
        self.aug_stmt = None
        self.augment_path = ''
        self.class_list = []
        self.class_name = ''
        self.class_stack = []
        self.comment = False
        self.contact = ''
        self.depth = 0
        self.env = None
        self.first = True
        self.group_list = []
        self.groupings = {}
        self.idx = 0
        self.idx_stack = []
        self.import_enum = []
        self.imports = {}
        self.loader = None
        self.local_group_list = []
        self.meta = True
        self.module = None
        self.module_name = ''
        self.namespace = ''
        self.ns = []
        self.organization = ''
        self.prefix = ''
        self.printer = None
        self.revision = ''
        self.rpc = []
        self.tab_size = 4
        self.target = ''
        self.templates = None
        self.types = []
        self.uses = []

    def str(self, msg):
        self.fd.write(msg)

    def tab(self, lvl=None):
        if lvl is None:
            lvl = self.lvl
        if lvl > 0:
            fmt = '%%%ds' % (lvl * self.tab_size)
        else:
            return ''
        return fmt % ' '

    def write(self, msg):
        if self.lvl > 0:
            fmt = '%%%ds' % (self.lvl * self.tab_size)
            self.fd.write(fmt % ' ')
        self.fd.write(msg)

    def writeln(self, msg, tab=0):
        if self.lvl + tab > 0:
            fmt = '%%%ds' % ((self.lvl + tab) * self.tab_size)
            self.fd.write(fmt % ' ')
        self.fd.write(msg)
        self.fd.write('\n')

    def bline(self):
        self.fd.write('\n')

    def lvl_inc(self, tab=1):
        self.lvl += tab

    def lvl_dec(self, tab=1):
        self.lvl -= tab

    def push_idx(self):
        self.idx_stack.append(self.idx)

    def pop_idx(self):
        self.idx = self.idx_stack.pop()

    def push_class(self):
        self.class_stack.append(self.class_list)
        self.class_list = []

    def pop_class(self):
        self.class_list = self.class_stack.pop()


def yang_id(stmt):
    if hasattr(stmt, 'arg') and stmt.arg is not None:
        return stmt.arg.replace(':', '_')
    else:
        return None
