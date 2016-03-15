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

from . import resolve_profile
import logging
import json
import pyang
from pyang import plugin, Repository
from pyang import error
from pyang import util
import re
from . import api_model, common, gen_target
from .helper import get_rst_file_name
from .module_printer import PythonModulePrinter
import optparse
import os
import shutil
from .ydk_py import YDKPythonGen
from optparse import OptionParser
import sys
from common import YdkGenException

logger = logging.getLogger('ydkgen')
logger.addHandler(logging.NullHandler())


def _parse_and_return_modules(resolved_model_dir):
    """ Use pyang to parse the files and get a list of modules.

        :param str resolved_model_dir The directory where all models to be compiled are found.
        :raise YdkGenException If there was a problem parsing the modules
    """
    repos = pyang.FileRepository(resolved_model_dir, False)
    ctx = pyang.Context(repos)

    filenames = []

    #(name, rev, handle)
    # where handle is (format, absfilename)
    for (_, _, (_, filename)) in repos.get_modules_and_revisions(ctx):
        filenames.append(filename)

    modules = []

    r = re.compile(r"^(.*?)(\@(\d{4}-\d{2}-\d{2}))?\.(yang|yin)$")
    for filename in filenames:
        f = filename
        if filename.startswith('file://'):
            f = filename[len('file://') - 1:]
        try:
            fd = open(f)
            text = fd.read()
        except IOError as ex:
            logger.error("error %s: %s\n" % (filename, str(ex)))
            raise YdkGenException(ex)

        m = r.search(filename)
        ctx.yin_module_map = {}
        if m is not None:
            (name, _dummy, rev, _) = m.groups()
            name = os.path.basename(name)
            logger.debug(
                'Parsing file %s format %s name %s revision %s', filename, format, name, rev)
            module = ctx.add_module(filename, text, format, name, rev,
                                    expect_failure_error=False)
        else:
            module = ctx.add_module(filename, text)
        if module is None:
            raise YdkGenException('Could not add module ')
        else:
            modules.append(module)

    # all the module have been added so get the context to validate
    # call prevalidate before this and post validate after
    ctx.validate()

    def keyfun(e):
        if e[0].ref == filenames[0]:
            return 0
        else:
            return 1

    ctx.errors.sort(key=lambda e: (e[0].ref, e[0].line))
    if len(filenames) > 0:
        # first print error for the first filename given
        ctx.errors.sort(key=keyfun)

    error_messages = []
    for (epos, etag, eargs) in ctx.errors:

        elevel = error.err_level(etag)
        if error.is_warning(elevel):
            logger.warning('%s: %s\n' %
                           (str(epos), error.err_to_str(etag, eargs)))
        else:
            err_msg = '%s: %s\n' % (str(epos), error.err_to_str(etag, eargs))
            logger.error(err_msg)
            error_messages.append(err_msg)

    if len(error_messages) > 0:
        err_msg = '\n'.join(error_messages)
        raise YdkGenException(err_msg)

    return [m for m in modules if m.keyword == 'module']


def generate(profile_file, output_directory, nodoc, ydk_root, groupings_as_class=False):
    """
        Generate ydk-py based in the output_directory using the supplied 
        profile_file

        :param str profile_file: The profile file to use
        :param str output_directory: The output directory where the generated ydk python code will be created.
        :param bool nodoc If set to true the documentation is not generated.
        :param str ydk_root: The ydk root directory. Relative file names in the profile file are resolved relative to this.
        :param bool groupings_as_class: If set to true YANG grouping is converted to a class.
        :raise YdkGenException: if an error has occurred
    """

    resolved_model_dir = None

    if profile_file is None:
        logger.error('profile_file is None.')
        raise YdkGenException('profile_file cannot be None.')

    if output_directory is None:
        logger.error('output_directory is None.')
        raise YdkGenException('output_directory cannot be None.')

    if ydk_root is None:
        logger.error('ydk_root is None.')
        raise YdkGenException('YDKGEN_HOME is not set.')

    try:
        with open(profile_file) as json_file:
            profile_data = json.load(json_file)
            resolved_model_dir = resolve_profile.resolve_profile(
                profile_data, ydk_root)
    except IOError as e:
        logger.error('Cannot open profile file (%s)', e.strerror)
        raise YdkGenException(e.strerror)
    except ValueError as e:
        logger.error('Cannot parse profile file (%s)', e.message)

    modules = _parse_and_return_modules(resolved_model_dir)

    py_sdk_root = output_directory + '/python/'
    py_api_doc_gen = py_sdk_root + '/docsgen'
    py_api_doc = output_directory + '/python/docs_expanded'

    if os.path.isdir(py_sdk_root):
        shutil.rmtree(py_sdk_root)
    if os.path.isdir(py_api_doc):
        shutil.rmtree(py_api_doc)
    if os.path.isfile(py_api_doc_gen + '/getting_started.rst'):
        os.remove(py_api_doc_gen + '/getting_started.rst')

    shutil.copytree(ydk_root + '/sdk/python', output_directory + '/python',
                    ignore=shutil.ignore_patterns('.gitignore', 'ncclient'))

    # no errors found now we can begin the transformation

    # begin generation
    ydk_dir = '%sydk' % py_sdk_root
    ydk_doc_dir = '%sdocsgen' % py_sdk_root

    yang_plugin = YDKPythonGen(ydk_dir, ydk_doc_dir, groupings_as_class)
    yang_plugin.emit(modules)
