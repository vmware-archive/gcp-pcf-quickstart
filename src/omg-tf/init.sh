#!/usr/bin/env bash

if [ -z ${DNS_SUFFIX+x} ]; then
    echo "DNS_SUFFIX required"
    exit 1
fi

if [ -z ${PROJECT_ID+x} ]; then
    export PROJECT_ID=${PROJECT_ID-`gcloud config get-value project  2> /dev/null`}
    echo "PROJECT_ID unset, using: ${PROJECT_ID}"
fi

# TODO(jrjohnson): Once a baked OpsMan image is ready, default to using it here
if [ -z ${BASE_IMAGE_URL+x} ] && [ -z ${BASE_IMAGE_SELFLINK+x} ]; then
    echo "BASE_IMAGE_URL or BASE_IMAGE_SELFLINK is required"
    exit 1
fi

if [ -z ${ENV_NAME+X} ]; then
    export ENV_NAME="omg"
    echo "ENV_NAME unset, using: ${ENV_NAME}"
fi


service_account_email=omg-terraform@${PROJECT_ID}.iam.gserviceaccount.com
service_account_file=$(mktemp)

gcloud iam service-accounts create omg-terraform --display-name terraform  2> /dev/null
gcloud iam service-accounts keys create ${service_account_file} --iam-account ${service_account_email}
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member serviceAccount:${service_account_email} \
  --role roles/owner

mkdir -p keys
pushd keys
  openssl genrsa -des3 -passout pass:x -out server.pass.key 2048
  openssl rsa -passin pass:x -in server.pass.key -out server.key
  openssl req -new -key server.key -out server.csr \
  -subj "/C=US/ST=Washington/L=Seattle/CN=${ENV_NAME}.${DNS_SUFFIX}/subjectAltName=*.${ENV_NAME}.${DNS_SUFFIX}"
  openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

  rm -f jumpbox_ssh jumpbox_ssh.pub
  ssh-keygen -b 2048 -t rsa -f jumpbox_ssh -q -N ""
popd

cat << VARS_FILE > terraform.tfvars
env_name = "${ENV_NAME}"
project = "${PROJECT_ID}"
dns_suffix = "${DNS_SUFFIX}"
opsman_image_url = "${BASE_IMAGE_URL}"
opsman_image_selflink = "${BASE_IMAGE_SELFLINK}"

ssl_cert = <<SSL_CERT
$(cat keys/server.crt)
SSL_CERT

ssl_cert_private_key = <<SSL_KEY
$(cat keys/server.key)
SSL_KEY

service_account_key = <<SERVICE_ACCOUNT_KEY
$(cat ${service_account_file})
SERVICE_ACCOUNT_KEY

ssh_public_key = <<SSH_PUBLIC_KEY
$(cat keys/jumpbox_ssh.pub)
SSH_PUBLIC_KEY

VARS_FILE
