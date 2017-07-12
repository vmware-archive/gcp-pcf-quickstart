#!/bin/bash
set -uex

envName=$(date +%s)
projectOwner=$(gcloud config get-value account)
rootDomain=${rootDomain}

mkdir -p ssl
pushd ssl
  openssl genrsa -des3 -passout pass:x -out server.pass.key 2048
  openssl rsa -passin pass:x -in server.pass.key -out server.key
  openssl req -new -key server.key -out server.csr \
  -subj "/C=US/ST=Washington/L=Seattle/CN=*.${rootDomain}"
  openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
popd

echo "
export envName=${envName}
export rootDomain=${rootDomain}
" >> .envrc
direnv allow || true

gcloud deployment-manager deployments create omg-${envName} --config=omg.jinja --properties="projectOwner:${projectOwner},rootDomain:${rootDomain}"
