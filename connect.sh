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

if [ -z ${ENV_NAME+X} ]; then
    export ENV_NAME="omg"
    echo "ENV_NAME unset, using: ${ENV_NAME}"
fi

if [ -z ${ENV_DIR+X} ]; then
    export ENV_DIR="$PWD/env/${ENV_NAME}"
    echo "ENV_DIR unset, using: ${ENV_DIR}"
fi

terraform_state="${ENV_DIR}/terraform.tfstate"
ssh_key="${ENV_DIR}/keys/jumpbox_ssh"

sshuttle -e "ssh -i ${ssh_key} -l omg" -r $(terraform output -state ${terraform_state} jumpbox_public_ip) 10.0.0.0/16
