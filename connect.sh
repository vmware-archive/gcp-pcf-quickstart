#!/usr/bin/env bash

if [ -z ${ENV_DIR+X} ]; then
    export ENV_DIR="$PWD/env/omg"
    echo "ENV_DIR unset, using: ${ENV_DIR}"
fi

terraform_state="${ENV_DIR}/terraform.tfstate"
ssh_key="${ENV_DIR}/keys/jumpbox_ssh"

sshuttle -e "ssh -i ${ssh_key} -l omg" -r $(terraform output -state ${terraform_state} jumpbox_public_ip) 10.0.0.0/16
