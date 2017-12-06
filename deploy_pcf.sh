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

if ! terraform version | grep -q "v0.9.11"; then
    echo "Error: Incompatible version of terraform. v0.9.11 required."
    echo ""
    echo "Linux: https://releases.hashicorp.com/terraform/0.9.11/terraform_0.9.11_linux_amd64.zip"
    exit 1
fi

set -ue
cd "$(dirname "$0")"

export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin
go install omg-cli

if [ -z ${ENV_DIR+X} ]; then
    export ENV_DIR="$PWD/env/pcf"
    echo "ENV_DIR unset, using: ${ENV_DIR}"
fi


mkdir -p ${ENV_DIR}
terraform_output="${ENV_DIR}/terraform_output.json"
terraform_config="${ENV_DIR}/terraform.tfvars"
terraform_state="${ENV_DIR}/terraform.tfstate"
env_config="${ENV_DIR}/config.json"

if [ ! -f $env_config ]; then
    if [ -z ${PIVNET_API_TOKEN+x} ]; then
        echo "PIVNET_API_TOKEN environment required (requires registration)."
        echo "Find the value for 'API TOKEN' at https://network.pivotal.io/users/dashboard/edit-profile"
        echo "and run: export PIVNET_API_TOKEN=<value of 'API TOKEN'> before running this command."
        exit 1
    fi

    omg-cli generate-config --env-dir="${ENV_DIR}"
    echo ""
    echo "The following settings are defaults:"
    echo ""
    omg-cli source-config --env-dir="${ENV_DIR}"

    echo ""
    echo "Review the settings above. Modify them by editing the file: ${env_config} and re-running this script"
    echo ""
    read -p "Accept defaults (y/n)? " choice

    case "$choice" in
      y|Y );;
      * ) exit 0;;
    esac
fi

set -o allexport
eval $(omg-cli source-config --env-dir="${ENV_DIR}")
set +o allexport

if [ -z ${PIVNET_ACCEPT_EULA+x} ]; then
    omg-cli review-eulas --env-dir="${ENV_DIR}"
else
    omg-cli review-eulas --env-dir="${ENV_DIR}" --accept-all
fi

pushd src/omg-tf
    # Verify project is ready
    if [ ! -f $terraform_config ]; then
         omg-cli prepare-project --env-dir="${ENV_DIR}"
        ./init.sh
    fi

    # Setup infrastructure
    gcloud config set project ${PROJECT_ID}
    terraform init
    terraform get
    terraform apply --parallelism=100 -state=${terraform_state} -var-file=${terraform_config} || terraform apply --parallelism=100 -state=${terraform_state} -var-file=${terraform_config}
    terraform output -json -state=${terraform_state} > ${terraform_output}
popd

# Deploy PCF
omg-cli remote --env-dir="${ENV_DIR}" "push-tiles"
omg-cli remote --env-dir="${ENV_DIR}" "deploy $@"
