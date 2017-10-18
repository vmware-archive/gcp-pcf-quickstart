#!/usr/bin/env bash

set -e

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
    extract_env
popd > /dev/null

trap save_terraform_state EXIT

pushd "${release_dir}/src/omg-tf"
	terraform init
	terraform get
	yes "yes" | terraform destroy --parallelism=100 -state=${terraform_state} -var-file=${terraform_config}
popd
