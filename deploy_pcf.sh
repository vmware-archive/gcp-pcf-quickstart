#!/usr/bin/env bash

set -uex

if [ -z ${PROJECT_ID+x} ]; then
    export PROJECT_ID=${PROJECT_ID-`gcloud config get-value project  2> /dev/null`}
    echo "PROJECT_ID unset, using: ${PROJECT_ID}"
fi

gcloud config set project ${PROJECT_ID}

if [ -z ${BASE_IMAGE_SELFLINK+x} ]; then
    echo "BASE_IMAGE_SELFLINK required"
    exit 1
fi

if [ -z ${DNS_SUFFIX+x} ]; then
    echo "DNS_SUFFIX required"
    exit 1
fi

function wait_for_host {
    host=$1

    set +e
    while ! ssh -oStrictHostKeyChecking=no -i ~/.ssh/google_compute_engine ${host} "exit 0"
    do
        sleep 2
        echo "Waiting for host to wake up.."
    done
    set -e
}

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

# Deploy PCF
pushd src/omg-cli
    go build
    wait_for_host $jumpbox_ip
    scp -i ~/.ssh/google_compute_engine -oStrictHostKeyChecking=no omg-cli ${jumpbox_ip}:.
    scp -i ~/.ssh/google_compute_engine -oStrictHostKeyChecking=no env.json ${jumpbox_ip}:.
    ssh -i ~/.ssh/google_compute_engine -oStrictHostKeyChecking=no ${jumpbox_ip} "./omg-cli --mode ConfigureOpsManager"
popd
