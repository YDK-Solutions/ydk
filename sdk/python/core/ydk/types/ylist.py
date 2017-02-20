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

""" types.py

    Contains type definitions for:
        - YList
"""

class YList(list):
    """ Represents a list with support for hanging a parent

        All YANG based entity classes that have lists in them use YList
        to represent the list.

        The "list" statement is used to define an interior data node in the
        schema tree.  A list node may exist in multiple instances in the data
        tree.  Each such instance is known as a list entry.  The "list"
        statement takes one argument, which is an identifier, followed by a
        block of substatements that holds detailed list information.

        A list entry is uniquely identified by the values of the list's keys,
        if defined.

    """
    def __init__(self, parent):
        super(YList, self).__init__()
        self.parent = parent

    def append(self, item):
        item.parent = self.parent
        super(YList, self).append(item)

    def extend(self, items):
       for item in items:
           self.append(item)
