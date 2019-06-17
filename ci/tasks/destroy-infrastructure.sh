#!/usr/bin/env bash

set -e

pushd omg-src-in/ci/tasks > /dev/null
	source utils.sh
	set_resource_dirs
  set_gcloud_config
  extract_env
popd > /dev/null

trap save_terraform_state EXIT

pushd "${release_dir}/src/omg-tf"
        configure_terraform_backend
	terraform init
	terraform get
	terraform destroy --auto-approve --parallelism=100 -var-file=${terraform_config} && exit 0

	seconds=300
	echo "terraform destroy failed, trying again in ${seconds} seconds"
	sleep ${seconds}
	terraform destroy --auto-approve --parallelism=100 -var-file=${terraform_config}
popd
