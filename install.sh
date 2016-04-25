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
# -*-sh-*-
# Prepare virtual environment for ydk-gen,
# Requirement: device with pip and virtualenv installed

export YDKGEN_HOME=`pwd`
MY_ENV=mypython
OUT_DIR=$YDKGEN_HOME/gen-api/python
PROFILE_FILE=$YDKGEN_HOME/profiles/ydk/ydk_0_4_0.json

# Install virtual environment
echo "Installing virtual environment under ${PWD}"
virtualenv $MY_ENV
source $MY_ENV/bin/activate

# Install pip dependency packages
echo "Installing python dependency"
pip install -r requirements.txt

echo "To exit current environment:"
echo "  $ deactivate"