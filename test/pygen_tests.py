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

'''
Python generation tests

Copyright (c) 2015 by Cisco Systems, Inc.
All rights reserved.

'''
import sys
import os
import optparse
import re
import subprocess
import shutil
import unittest
import filecmp
import difflib
from shutil import rmtree
from filecmp import dircmp

yang_mod_path = []

def we_are_frozen():
    #All of the modules are built-in to the interpreter e.g by p2e
    return hasattr(sys, "frozen")

def module_path():
    encoding = sys.getfilesystemencoding()
    if we_are_frozen():
        return os.path.dirname(unicode(sys.executable, encoding))
    return os.path.dirname(unicode(__file__, encoding))    

def suite():
    
    class PyGenTest(unittest.TestCase):
        def __init__(self):
            setattr(self, "Python Gen", self.translate_and_check)
            unittest.TestCase.__init__(self, "Python Gen")

        def id(self):
            return 'Python Gen'
    
        def compare(self, src, dest):
            if os.path.isfile(src):
                if src.split('/')[-1].startswith('.'):
                    return
                same = filecmp.cmp(src, dest, False)
                if not same:
                    dest_file = open(dest, 'r')
                    src_file = open(src, 'r')
                    diff = difflib.context_diff(dest_file.readlines(), src_file.readlines())
                    delta = ''.join(diff)
                    self.assertEquals(same, True, 'File comparison failed for test case %s\n\
                       expected %s \n \
                       generated %s \n \
                       if generated is correct use command \n \
                       cp %s %s \n \
                         diff %s'%(self.test_case_name, dest, src, src, os.path.dirname(dest), delta))
                    return
        
    
        def translate_and_check(self):
            self.assertTrue(os.environ.has_key('YDKGEN_HOME'), 'Need to have YDKGEN_HOME set!')
            ydk_root = os.environ['YDKGEN_HOME']
            if not os.environ.has_key('PYANG_PLUGINPATH'):
                os.environ['PYANG_PLUGINPATH'] = ydk_root + '/src/pyang-plugins'
            else:
                os.environ['PYANG_PLUGINPATH'] += ':' + ydk_root + '/src/pyang-plugins'
            
            pythonpath_additions = ydk_root + ':' + ydk_root + '/src/translators'
            if not os.environ.has_key('PYTHONPATH'):
                os.environ['PYTHONPATH'] = pythonpath_additions
            else:
                os.environ['PYTHONPATH'] += ':' +pythonpath_additions
            
            os.chdir(ydk_root)
            
            yang_models = []
            ydk_model_paths = []
            
            # Create a list of fully expanded models to work with (yang_models)
            # and a list of the directories within which models reside
            # (yang_model_paths)
            #
            
            for path, sub, files in os.walk(os.getenv('YANG_MODELS', ydk_root + '/yang/ydktest')):
                for f in files:
                    if f.endswith('.yang'):
                        with open(os.path.join(path, f)) as fd:
                            subm = False
                            for line in fd.readlines():
                                # check first line for module or submodule
                                if 'submodule' in line:
                                    subm = True
                                    break
                                elif 'module' in line:
                                    break
                            if subm:
                                continue
                            if path not in ydk_model_paths:
                                ydk_model_paths.append(path)
                            yang_models.append(os.path.join(path, f))
                                
                
            # Convert the list of model paths to a mroe standard string in a unix
            # path format for use later as a parameter to pyang
            #
            ydk_model_paths = ':'.join(ydk_model_paths)
            
            current_env = os.environ.copy()
            pyang, stderr = subprocess.Popen(['which', 'pyang'], stdout=subprocess.PIPE, env=current_env).communicate()
            self.assertTrue(len(pyang) > 0, "Cannot find pyang, please enable in your environment!")
            pyang = pyang[:-1]
            py_sdk_root = ydk_root + '/test/test-cases/python/actual/ydk'
            expected_py_sdk_root = ydk_root + '/test/test-cases/python/expected/ydk'
            
            if os.path.isdir(py_sdk_root):
                rmtree(py_sdk_root)
            
            os.mkdir(py_sdk_root)
            os.mkdir('%s/models'%py_sdk_root)
            
            args = [sys.executable, pyang, '-p', ydk_model_paths,
                            '-f', 'ydk-py', '--ydk-dir', py_sdk_root]
            
            args.extend(yang_models)
            
            print ' '.join(args)
            exit_code = subprocess.call(args, env=os.environ.copy(), cwd=os.environ['YDKGEN_HOME'])
            if exit_code != 0:
                print 'Code generation failed !!!'
            
            def check_diff_files(dcmp, diff_files):
                
                for name in dcmp.diff_files:
                    diff_files.append('File %s/%s does not match'%(dcmp.left,name))
                    #dest_file = open('%s/%s'%(dcmp.left,name), 'r')
                    #src_file = open('%s/%s'%(dcmp.right,name), 'r')
                    #diff = difflib.context_diff(dest_file.readlines(), src_file.readlines())
                    #delta = ''.join(diff)
                    #diff_files.append('File %s at %s does not match diff:-\n %s \n' % (name, dcmp.left, delta))
                for sub_dcmp in dcmp.subdirs.values():
                    check_diff_files(sub_dcmp, diff_files)
            
            result = dircmp(py_sdk_root, expected_py_sdk_root)
            diff_files = []
            check_diff_files(result, diff_files)
            self.assertTrue(len(diff_files) == 0, 'Total number of files that differ %s.\nList of files that differ:- %s'%(len(diff_files),'\n'.join(diff_files)))
            
            #print result.report()
            
            
        
        def setUp(self):
            pass

        def tearDown(self):
            pass


        def testName(self):
            self.yang_file
    
    suite = unittest.TestSuite()

    suite.addTest(PyGenTest())
    
    return suite


if __name__ == "__main__":
    ret = unittest.TextTestRunner(verbosity=2).run(suite()).wasSuccessful()
    if ret:
        sys.exit(0)
    else:
        sys.exit(1)

def load_tests(loader , tests, pattern):
    return suite()
