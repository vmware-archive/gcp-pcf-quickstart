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

if [ -z ${ZONE_1+x} ]; then
    echo "ZONE_1 required"
    exit 1
fi

if [ -z ${ZONE_2+x} ]; then
    echo "ZONE_2 required"
    exit 1
fi

if [ -z ${ZONE_3+x} ]; then
    echo "ZONE_3 required"
    exit 1
fi

if [ -z ${REGION+x} ]; then
    echo "REGION required"
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

nat_instance_count=3
nat_machine_type="n1-standard-1"
opsman_machine_type="n1-standard-2"
jumpbox_machine_type="n1-standard-1"
external_database="true"
if [ "${SMALL_FOOTPRINT}" == "true" ]; then
  nat_instance_count=1
  nat_machine_type="g1-small"
  opsman_machine_type="n1-standard-1"
  jumpbox_machine_type="g1-small"
  external_database="false"
fi

set -e
seed=$(date +%s)

#
# Generate SSL/SSH Keys
#

mkdir -p keys
pushd keys
  SYS_DOMAIN=sys.${DNS_SUFFIX}
  APPS_DOMAIN=apps.${DNS_SUFFIX}

  SSL_FILE=sslconf-pas.conf
  ROOT_CERT=rootca
  PAS_CERT=server

  #Generate SSL Config with SANs
if [ ! -f $SSL_FILE ]; then
cat > $SSL_FILE <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
[req_distinguished_name]
countryName_default = US
stateOrProvinceName_default = CA
localityName_default = SF
organizationalUnitName_default = Pivotal
[ v3_req ]
# Extensions to add to a certificate request
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = *.${SYS_DOMAIN}
DNS.2 = *.${APPS_DOMAIN}
DNS.3 = *.login.${SYS_DOMAIN}
DNS.4 = *.uaa.${SYS_DOMAIN}
EOF
fi

  openssl genrsa -out ${ROOT_CERT}.key 2048
  openssl req -x509 -new -nodes -key ${ROOT_CERT}.key -sha256 -days 1024 -subj "/C=US/ST=CA/O=Pivotal/L=SF/OU=PA/CN=pivotal.io" -out ${ROOT_CERT}.crt

  openssl genrsa -out ${PAS_CERT}.key 2048
  openssl req -new -out ${PAS_CERT}.csr -subj "/C=US/ST=CA/O=Pivotal/L=SF/OU=PA/CN=pivotal.io" -key ${PAS_CERT}.key -config ${SSL_FILE}
  openssl x509 -req -days 3650 -sha256 -in ${PAS_CERT}.csr -CA ${ROOT_CERT}.crt -CAkey ${ROOT_CERT}.key -CAcreateserial -out ${PAS_CERT}.crt -extensions v3_req -extfile ${SSL_FILE}

  rm ${ROOT_CERT}.key
  rm ${ROOT_CERT}.srl
  rm ${PAS_CERT}.csr

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
opsman_machine_type = "${opsman_machine_type}"
nat_machine_type="${nat_machine_type}"
jumpbox_machine_type="${jumpbox_machine_type}"
ops_manager_password = "${OPSMAN_ADMIN_PASSWORD:-}"
ops_manager_skip_ssl_verify = "true"
region = "${REGION}"
zones = ["${ZONE_1}", "${ZONE_2}", "${ZONE_3}"]
external_database = "${external_database}"

ssl_root_ca = <<SSL_ROOT_CA
$(cat keys/rootca.crt)
SSL_ROOT_CA

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

nat_instance_count = ${nat_instance_count}
VARS_FILE
