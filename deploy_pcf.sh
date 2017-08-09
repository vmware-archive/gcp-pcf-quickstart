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

set -ue
cd "$(dirname "$0")"

if [ -z ${DNS_ZONE_NAME+x} ]; then
    export DNS_ZONE_NAME=omg-zone
    echo "DNS_ZONE_NAME unset, using: ${DNS_ZONE_NAME}"
fi

if [ -z ${PROJECT_ID+x} ]; then
    export PROJECT_ID=${PROJECT_ID-`gcloud config get-value project  2> /dev/null`}
    echo "PROJECT_ID unset, using: ${PROJECT_ID}"
fi

gcloud config set project ${PROJECT_ID}

if [ -z ${BASE_IMAGE_SELFLINK+x} ]; then
    export BASE_IMAGE_SELFLINK="baked-opsman-1501261888"
    echo "BASE_IMAGE_SELFLINK unset, using: ${BASE_IMAGE_SELFLINK}"
fi

if [ -z ${ENV_NAME+X} ]; then
    export ENV_NAME="omg"
    echo "ENV_NAME unset, using: ${ENV_NAME}"
fi

if [ -z ${ENV_DIR+X} ]; then
    export ENV_DIR="$PWD/env/${ENV_NAME}"
    echo "ENV_DIR unset, using: ${ENV_DIR}"
fi

mkdir -p ${ENV_DIR}
terraform_output="${ENV_DIR}/env.json"
terraform_config="${ENV_DIR}/terraform.tfvars"
terraform_state="${ENV_DIR}/terraform.tfstate"

# Setup infrastructure
pushd src/omg-tf
    if [ ! -f $terraform_config ]; then
        ./init.sh
    fi
    terraform init
    terraform get
    terraform apply --parallelism=100 -state=${terraform_state} -var-file=${terraform_config}
    terraform output -json -state=${terraform_state} > ${terraform_output}
popd

# Deploy PCF
export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin
go install omg-cli
omg-cli bootstrap-deploy --ssh-key-path "${ENV_DIR}/keys/jumpbox_ssh" --username omg --terraform-output-path ${terraform_output} $@
