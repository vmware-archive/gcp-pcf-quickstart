#!/usr/bin/env bash

set -ue
cd "$(dirname "$0")"

if [ -z ${DNS_ZONE_NAME+x} ]; then
    echo "DNS_ZONE_NAME required"
    exit 1
fi

if [ -z ${PROJECT_ID+x} ]; then
    export PROJECT_ID=${PROJECT_ID-`gcloud config get-value project  2> /dev/null`}
    echo "PROJECT_ID unset, using: ${PROJECT_ID}"
fi

gcloud config set project ${PROJECT_ID}

if [ -z ${BASE_IMAGE_SELFLINK+x} ]; then
    export BASE_IMAGE_SELFLINK="projects/graphite-demo-jjcf/global/images/baked-opsman-1501021113"
    echo "BASE_IMAGE_SELFLINK unset, using: ${BASE_IMAGE_SELFLINK}"
fi

terraform_output=$(mktemp)

# Setup infrastructure
pushd src/omg-tf
    if [ ! -f terraform.tfvars ]; then
        ./init.sh
    fi
    terraform apply
    terraform output -json > ${terraform_output}
popd

# Deploy PCF
go install omg-cli
omg-cli bootstrap-deploy --ssh-key-path src/omg-tf/keys/jumpbox_ssh --username omg --terraform-output-path ${terraform_output}
