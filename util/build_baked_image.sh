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

if [ -z ${PROJECT_ID+x} ]; then
    export PROJECT_ID=${PROJECT_ID-`gcloud config get-value project  2> /dev/null`}
    echo "PROJECT_ID unset, using: ${PROJECT_ID}"
fi

gcloud config set project ${PROJECT_ID}

if [ -z ${BASE_IMAGE_URL+x} ]; then
    echo "BASE_IMAGE_URL required"
    exit 1
fi

if [ -z ${PIVNET_API_TOKEN+x} ]; then
    echo "PIVNET_API_TOKEN required"
    exit 1
fi

if [ -z ${DNS_ZONE_NAME+x} ]; then
    echo "DNS_ZONE_NAME required"
    exit 1
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

# Setup infrastructure
pushd src/omg-tf
    if [ ! -f terraform.tfvars ]; then
        ./init.sh
    fi
    terraform apply --parallelism=100 -state=${ENV_DIR}
    terraform output -json -state=${ENV_DIR} > ${terraform_output}
    export opsman_instance_name=$(terraform output -state=${ENV_DIR} ops_manager_instance_name)
    export opsman_instance_zone=$(terraform output -state=${ENV_DIR} ops_manager_instance_zone)
popd

# Hydrate Ops Manager
export GOPATH=`pwd`
export PATH=$PATH:$GOPATH/bin
go install omg-cli
omg-cli remote --env-dir="${ENV_DIR}" "push-tiles --pivnet-api-token=${PIVNET_API_TOKEN}"

# Capture image
image_name="baked-opsman-$(date +%s)"
gcloud compute instances stop --zone=${opsman_instance_zone} ${opsman_instance_name}
gcloud compute images create ${image_name} \
  --source-disk ${opsman_instance_name} \
  --source-disk-zone ${opsman_instance_zone} \

echo "Image built as: ${image_name}"
