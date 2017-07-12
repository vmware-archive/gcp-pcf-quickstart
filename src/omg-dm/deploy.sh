#!/bin/bash
set -uex

envName=$(date +%s)
projectOwner=$(gcloud config get-value account)
rootDomain=${rootDomain}
region=${region:-us-east1}
zone1=${zone1:-us-east1-b}
zone2=${zone2:-us-east1-c}
zone3=${zone3:-us-east1-d}


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
export region=${region}
export zone1=${zone1}
export zone2=${zone2}
export zone3=${zone3}
" >> .envrc
direnv allow || true

gcloud deployment-manager deployments create omg-${envName} --config=omg.jinja --properties="projectOwner:${projectOwner},rootDomain:${rootDomain},region:${region},zone1:${zone1},zone2:${zone2},zone3:${zone3}"
