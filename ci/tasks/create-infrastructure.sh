#!/usr/bin/env bash

set -eE

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
    extract_env
popd > /dev/null

trap save_terraform_state EXIT

function rollback {
	pushd "${release_dir}/src/omg-tf"
		yes "yes" | terraform destroy --parallelism=100 -state=${terraform_state} -var-file=${terraform_config}
	popd
}
trap rollback ERR

pushd "${release_dir}/src/omg-tf"
	terraform init
	terraform get
	terraform apply --parallelism=100 -state=${terraform_state} -var-file=${terraform_config} || terraform apply --parallelism=100 -state=${terraform_state} -var-file=${terraform_config}
	terraform output -json -state=${terraform_state} > ${terraform_output}
popd

