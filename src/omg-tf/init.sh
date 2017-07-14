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
if [ -z ${BASE_IMAGE+x} ]; then
    echo "BASE_IMAGE required"
    exit 1
fi


service_account_email=omg-terraform@${PROJECT_ID}.iam.gserviceaccount.com
service_account_file=$(mktemp)

gcloud iam service-accounts create omg-terraform --display-name terraform  2> /dev/null
gcloud iam service-accounts keys create ${service_account_file} --iam-account ${service_account_email}
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member serviceAccount:${service_account_email} \
  --role roles/owner

mkdir -p ssl
pushd ssl
  openssl genrsa -des3 -passout pass:x -out server.pass.key 2048
  openssl rsa -passin pass:x -in server.pass.key -out server.key
  openssl req -new -key server.key -out server.csr \
  -subj "/C=US/ST=Washington/L=Seattle/CN=*.${DNS_SUFFIX}"
  openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
popd

cat << VARS_FILE > terraform.tfvars
project = "${PROJECT_ID}"
dns_suffix = "${DNS_SUFFIX}"
opsman_image_url = "${BASE_IMAGE}"

ssl_cert = <<SSL_CERT
$(cat ssl/server.crt)
SSL_CERT

ssl_cert_private_key = <<SSL_KEY
$(cat ssl/server.key)
SSL_KEY

service_account_key = <<SERVICE_ACCOUNT_KEY
$(cat ${service_account_file})
SERVICE_ACCOUNT_KEY

VARS_FILE
