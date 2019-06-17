#!/usr/bin/env bash

set -ex

pushd omg-src-in/ci/tasks > /dev/null
	source utils.sh
	set_resource_dirs
	build_go
  set_gcloud_config
  generate_env_config
popd > /dev/null

set -o allexport
eval $(omg-cli source-config --env-dir="${env_dir}")
set +o allexport

trap save_terraform_state EXIT
pushd "${release_dir}/src/omg-tf"
	ENV_DIR=${env_dir} ./init.sh
popd
