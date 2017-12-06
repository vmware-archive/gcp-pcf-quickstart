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
    export ENV_DIR="${PWD}/env/pcf"
    echo "ENV_DIR unset, using: ${ENV_DIR}"
fi

read -p "Delete deployment in ${ENV_DIR} (y/n)? " choice
case "$choice" in
  y|Y ) echo "begin delete";;
  * ) exit 0;;
esac

terraform_output="${ENV_DIR}/terraform_output.json"
terraform_config="${ENV_DIR}/terraform.tfvars"
terraform_state="${ENV_DIR}/terraform.tfstate"

export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin
go install omg-cli
omg-cli remote --env-dir ${ENV_DIR} "delete-installation"
omg-cli cleanup-project --env-dir ${ENV_DIR} --no-dry-run

pushd src/omg-tf
    yes "yes" | terraform destroy --parallelism=100 -state=${terraform_state} -var-file=${terraform_config}
popd
