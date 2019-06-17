#!/usr/bin/env bash

set -eE

pushd omg-src-in > /dev/null
	source ci/tasks/utils.sh
	set_resource_dirs
  set_gcloud_config
  extract_env
popd > /dev/null

trap save_terraform_state EXIT

function rollback {
	pushd "${release_dir}/src/omg-tf"
		yes "yes" | terraform destroy --parallelism=100 -var-file=${terraform_config}
	popd
}
trap rollback ERR

pushd "${release_dir}/src/omg-tf"
        configure_terraform_backend
	terraform init
	terraform get
	terraform apply --auto-approve --parallelism=100 -var-file=${terraform_config}
	terraform output -json > ${terraform_output}
popd
