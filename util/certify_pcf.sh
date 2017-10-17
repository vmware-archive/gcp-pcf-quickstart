#!/usr/bin/env bash

#
# Copyright 2017 Google Inc.
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
#

set -u
cd "$(dirname "$0")/../"

if [ -z ${ENV_DIR+X} ]; then
    export ENV_DIR="env/pcf"
    echo "ENV_DIR unset, using: ${ENV_DIR}"
fi

if ! [ -x "$(command -v ginkgo)" ]; then
    echo "Error: ginkgo missing, try: go install github.com/onsi/ginkgo"
    exit 1
fi

# Ensure absolute path
export ENV_DIR=$(readlink -f ${ENV_DIR})

export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin

pushd src/omg-cli/certification
    ginkgo -r .
popd