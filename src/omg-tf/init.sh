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

if [ -z ${ENV_DIR+X} ]; then
    echo "ENV_DIR required"
    exit 1
fi

if [ -z ${ENV_NAME+X} ]; then
    echo "ENV_NAME required"
    exit 1
fi

cd ${ENV_DIR}

if [ -z ${DNS_ZONE_NAME+x} ]; then
    echo "DNS_ZONE_NAME required"
    exit 1
fi

if [ -z ${PROJECT_ID+x} ]; then
    export PROJECT_ID=${PROJECT_ID-`gcloud config get-value project  2> /dev/null`}
    echo "PROJECT_ID unset, using: ${PROJECT_ID}"
fi

dns_suffix=$(gcloud dns managed-zones describe ${DNS_ZONE_NAME} --project ${PROJECT_ID} --format="value(dnsName)" 2> /dev/null)

if [ $? != 0 ]; then
    echo "ERROR: Expected to find Cloud DNS managed zone ${DNS_ZONE_NAME} in project ${PROJECT_ID}"
    exit 1
fi

# trim trailing '.' from response
export DNS_SUFFIX=${dns_suffix%.}
echo "DNS_SUFFIX set to: ${DNS_SUFFIX}"

if [ `dig ${DNS_SUFFIX} NS +short | wc -l` == "0" ]; then
    echo "Failed to resolve NS records for ${DNS_SUFFIX}"
    exit 1
fi

if [ -z ${BASE_IMAGE_URL+x} ] && [ -z ${BASE_IMAGE_SELFLINK+x} ]; then
    echo "BASE_IMAGE_URL or BASE_IMAGE_SELFLINK is required"
    exit 1
fi

set -e
seed=$(date +%s)

#
# Generate SSL/SSH Keys
#

mkdir -p keys
pushd keys
  openssl genrsa -passout pass:x -out server.pass.key 2048
  openssl rsa -passin pass:x -in server.pass.key -out server.key
  openssl req -new -key server.key -out server.csr \
  -subj "/C=US/ST=Washington/L=Seattle/CN=${ENV_NAME}.${DNS_SUFFIX}/subjectAltName=*.${DNS_SUFFIX}"
  openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

  rm -f jumpbox_ssh jumpbox_ssh.pub
  ssh-keygen -b 2048 -t rsa -f jumpbox_ssh -q -N ""
popd

#
# Create environment config
#

cat << VARS_FILE > terraform.tfvars
env_name = "${ENV_NAME}"
project = "${PROJECT_ID}"
dns_suffix = "${DNS_SUFFIX}"
dns_zone_name = "${DNS_ZONE_NAME}"
opsman_image_url = "${BASE_IMAGE_URL:-}"
opsman_image_selflink = "${BASE_IMAGE_SELFLINK:-}"
opsman_external_ip = "true"
ops_manager_skip_ssl_verify = "true"
region = "${REGION}"
zones = ["${ZONE1}", "${ZONE2}", "${ZONE3}"]
external_database = "true"

ssl_cert = <<SSL_CERT
$(cat keys/server.crt)
SSL_CERT

ssl_cert_private_key = <<SSL_KEY
$(cat keys/server.key)
SSL_KEY

service_account_key = ""

ssh_public_key = <<SSH_PUBLIC_KEY
$(cat keys/jumpbox_ssh.pub)
SSH_PUBLIC_KEY
VARS_FILE
