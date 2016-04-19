""" ydktest_union 

This module contains a collection of YANG definitions
for unions
    

Copyright (c) 2013\-2014 by Cisco Systems, Inc.
All rights reserved.

"""


import re
import collections

from enum import Enum

from ydk.types import Empty, YList, DELETE, Decimal64, FixedBitsDict

from ydk.errors import YPYError, YPYDataValidationError



class YdkEnumTestEnum(Enum):
    """
    YdkEnumTestEnum

    YDK Enum test

    .. data:: NOT_SET = 0

    	Not Set

    .. data:: NONE = 1

    	None

    .. data:: LOCAL = 2

    	Local

    .. data:: REMOTE = 3

    	Remote

    """

    NOT_SET = 0

    NONE = 1

    LOCAL = 2

    REMOTE = 3


    @staticmethod
    def _meta_info():
        from ydk.models.ydktest._meta import _ydktest_union as meta
        return meta._meta_table['YdkEnumTestEnum']



class BuiltInT(object):
    """
    config for built\-in types
    
    .. attribute:: younion
    
    	
    	**type**\: one of { :py:class:`YdkEnumTestEnum <ydk.models.ydktest.ydktest_union.YdkEnumTestEnum>` | int }
    
    

    """

    _prefix = 'ydkut'
    _revision = '2015-11-17'

    def __init__(self):
        self.younion = None

    @property
    def _common_path(self):

        return '/ydktest-union:built-in-t'

    def is_config(self):
        ''' Returns True if this instance represents config data else returns False '''
        return True

    def _has_data(self):
        if not self.is_config():
            return False
        if self.younion is not None:
            return True

        return False

    @staticmethod
    def _meta_info():
        from ydk.models.ydktest._meta import _ydktest_union as meta
        return meta._meta_table['BuiltInT']['meta_info']


