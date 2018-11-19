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

set -eu
cd "$(dirname "$0")/../"
root=$(pwd)

if [ -z ${ENV_DIR+X} ]; then
    export ENV_DIR="${PWD}/env/pcf"
    echo "ENV_DIR unset, using: ${ENV_DIR}"
fi

pushd src/omg-cli
go build -o $root/bin/omg-cli
popd
export PATH=$root/bin:$PATH

omg-cli remote --env-dir="${ENV_DIR}" "$@"
