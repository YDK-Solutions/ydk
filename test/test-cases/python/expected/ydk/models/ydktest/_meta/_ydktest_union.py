


import re
import collections

from enum import Enum

from ydk._core._dm_meta_info import _MetaInfoClassMember, _MetaInfoClass, _MetaInfoEnum
from ydk.types import Empty, YList, DELETE, Decimal64, FixedBitsDict
from ydk._core._dm_meta_info import ATTRIBUTE, REFERENCE_CLASS, REFERENCE_LIST, REFERENCE_LEAFLIST,     REFERENCE_IDENTITY_CLASS, REFERENCE_ENUM_CLASS, REFERENCE_BITS, REFERENCE_UNION

from ydk.errors import YPYError, YPYDataValidationError
from ydk.models import _yang_ns

_meta_table = {
    'YdkEnumTestEnum' : _MetaInfoEnum('YdkEnumTestEnum', 'ydk.models.ydktest.ydktest_union',
        {
            'not-set':'NOT_SET',
            'none':'NONE',
            'local':'LOCAL',
            'remote':'REMOTE',
        }, 'ydktest-union', _yang_ns._namespaces['ydktest-union']),
    'BuiltInT' : {
        'meta_info' : _MetaInfoClass('BuiltInT',
            False, 
            [
            _MetaInfoClassMember('younion', REFERENCE_UNION, 'str' , None, None, 
                [], [], 
                '''                ''',
                'younion',
                'ydktest-union', False, [
                    _MetaInfoClassMember('younion', REFERENCE_ENUM_CLASS, 'YdkEnumTestEnum' , 'ydk.models.ydktest.ydktest_union', 'YdkEnumTestEnum', 
                        [], [], 
                        '''                        ''',
                        'younion',
                        'ydktest-union', False),
                    _MetaInfoClassMember('younion', ATTRIBUTE, 'int' , None, None, 
                        [(0, 63)], [], 
                        '''                        ''',
                        'younion',
                        'ydktest-union', False),
                ]),
            ],
            'ydktest-union',
            'built-in-t',
            _yang_ns._namespaces['ydktest-union'],
        'ydk.models.ydktest.ydktest_union'
        ),
    },
}
