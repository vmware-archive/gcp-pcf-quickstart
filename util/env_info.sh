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
    echo ""
fi

if [ -z ${1+X} ] || [ "${1}" == "opsman" ]; then
    echo "Pivotal Operations Manager Credentials ================ "
    printf "URL: https://$(util/terraform_output.sh ops_manager_dns)\nusername: $(util/terraform_output.sh ops_manager_username)\npassword: $(util/terraform_output.sh ops_manager_password)\n"
fi

if [ -z ${1+X} ] || [ "${1}" == "cf" ]; then
    echo "Cloud Foundry API Credentials ========================= "
    echo "URL: https://api.$(util/terraform_output.sh sys_domain)"
    bin/omg-cli remote --env-dir=${ENV_DIR} "get-credential --app-name=cf --credential=.uaa.admin_credentials"
fi