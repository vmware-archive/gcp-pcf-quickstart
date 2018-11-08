#!/bin/bash

get_var() {
    bosh int ${ENV_DIR}/terraform_output.json --path /$1/value
}

while read var; do unset $var; done < <(env | grep BOSH | cut -d'=' -f1)

export OM_TARGET=$(get_var 'ops_manager_dns')
export OM_PASSWORD=$(get_var 'ops_manager_password')
export BOSH_ALL_PROXY="ssh+socks5://omg@$(get_var 'jumpbox_public_ip'):22?private-key=$ENV_DIR/keys/jumpbox_ssh"
BOSH_PRODUCT_GUID=$(om -t $OM_TARGET -k -u admin curl -s -path /api/v0/deployed/products/ | jq -r -c '.[] | select(.type | contains("p-bosh")) | .guid');
export BOSH_ENVIRONMENT=$(om -t $OM_TARGET -k -u admin curl -s -path /api/v0/deployed/products/$BOSH_PRODUCT_GUID/static_ips | jq -r '.[].ips[]');
export BOSH_CA_CERT=$(om -t $OM_TARGET -k -u admin curl -s -path /api/v0/certificate_authorities | jq -r '.certificate_authorities[].cert_pem');
export BOSH_USERNAME=$(om -t $OM_TARGET -k -u admin curl -s -path /api/v0/deployed/director/credentials/director_credentials | jq -r '.credential.value.identity');
export BOSH_PASSWORD=$(om -t $OM_TARGET -k -u admin curl -s -path /api/v0/deployed/director/credentials/director_credentials | jq -r '.credential.value.password');
echo -e "$BOSH_USERNAME\n$BOSH_PASSWORD\n" | bosh log-in
