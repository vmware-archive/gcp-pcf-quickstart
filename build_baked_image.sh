#!/usr/bin/env bash

set -xeu

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

export DNS_SUFFIX="example.org"

# Setup infrastructure
pushd src/omg-tf
    if [ ! -f terraform.tfvars ]; then
        ./init.sh
    fi
    terraform apply
    terraform output -json > ../omg-cli/env.json
    export jumpbox_ip=$(terraform output jumpbox_ip)
    export opsman_instance_name=$(terraform output ops_manager_instance_name)
    export opsman_instance_zone=$(terraform output ops_manager_instance_zone)
popd

# Ensure ssh keys are setup
gcloud compute config-ssh

# Bake image
pushd src/omg-cli
    go build
    scp -i ~/.ssh/google_compute_engine -oStrictHostKeyChecking=no omg-cli ${jumpbox_ip}:.
    scp -i ~/.ssh/google_compute_engine -oStrictHostKeyChecking=no env.json ${jumpbox_ip}:.
    ssh -i ~/.ssh/google_compute_engine -oStrictHostKeyChecking=no ${jumpbox_ip} "./omg-cli bake-image --pivnet-api-token=${PIVNET_API_TOKEN}"
popd

# Capture image
image_name="baked-opsman-$(date +%s)"
gcloud compute instances stop --zone=${opsman_instance_zone} ${opsman_instance_name}
gcloud compute images create ${image_name} \
  --source-disk ${opsman_instance_name} \
  --source-disk-zone ${opsman_instance_zone} \

echo "Image built as: ${image_name}"