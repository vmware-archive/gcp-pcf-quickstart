#!/usr/bin/env bash

set -ex

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
  set_gcloud_config
  generate_env_config
popd > /dev/null

go install omg-cli
set -o allexport
eval $(omg-cli source-config --env-dir="${env_dir}")
set +o allexport

trap save_terraform_state EXIT
pushd "${release_dir}/src/omg-tf"
	ENV_DIR=${env_dir} ./init.sh
popd