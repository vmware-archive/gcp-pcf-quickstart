#!/usr/bin/env bash

set -e

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
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
	yes "yes" | terraform destroy --parallelism=100 -var-file=${terraform_config} && exit 0

	seconds=300
	echo "terraform destroy failed, trying again in ${seconds} seconds"
	sleep ${seconds}
	yes "yes" | terraform destroy --parallelism=100 -var-file=${terraform_config}
popd
