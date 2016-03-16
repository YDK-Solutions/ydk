
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

_global_inherit_nsp = 'http://cisco.com/ns/yang/inherit'
_global_main_nsp = 'http://cisco.com/ns/yang/main'
_global_main_aug1_nsp = 'http://cisco.com/ns/yang/main-aug1'
_global_main_aug2_nsp = 'http://cisco.com/ns/yang/main-aug2'
_global_main_aug3_nsp = 'http://cisco.com/ns/yang/main-aug3'
_global_oc_pattern_nsp = 'http://cisco.com/ns/yang/oc-pattern'
_global_ydktest_filterread_nsp = 'http://cisco.com/ns/yang/ydk-filter'
_global_ydktest_sanity_nsp = 'http://cisco.com/ns/yang/ydktest-sanity'
_global_ydktest_sanity_augm_nsp = 'http://cisco.com/ns/yang/ydktest-sanity-augm'
_global_ydktest_union_nsp = 'http://cisco.com/ns/yang/ydktest-union'
_namespaces = { \
    'inherit' : 'http://cisco.com/ns/yang/inherit', 
    'main' : 'http://cisco.com/ns/yang/main', 
    'main-aug1' : 'http://cisco.com/ns/yang/main-aug1', 
    'main-aug2' : 'http://cisco.com/ns/yang/main-aug2', 
    'main-aug3' : 'http://cisco.com/ns/yang/main-aug3', 
    'oc-pattern' : 'http://cisco.com/ns/yang/oc-pattern', 
    'ydktest-filterread' : 'http://cisco.com/ns/yang/ydk-filter', 
    'ydktest-sanity' : 'http://cisco.com/ns/yang/ydktest-sanity', 
    'ydktest-sanity-augm' : 'http://cisco.com/ns/yang/ydktest-sanity-augm', 
    'ydktest-union' : 'http://cisco.com/ns/yang/ydktest-union', 
}

_identity_map = { \
    ('ydktest-sanity', 'base-identity'):('ydk.models.ydktest.ydktest_sanity', 'BaseIdentity_Identity'),
    ('ydktest-sanity', 'child-child-identity'):('ydk.models.ydktest.ydktest_sanity', 'ChildChildIdentity_Identity'),
    ('ydktest-sanity', 'child-identity'):('ydk.models.ydktest.ydktest_sanity', 'ChildIdentity_Identity'),
}

