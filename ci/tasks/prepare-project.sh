#!/usr/bin/env bash

set -e

pushd omg-src-in/ci/tasks > /dev/null
	source utils.sh
	set_resource_dirs
	build_go
  set_gcloud_config
  generate_env_config
popd > /dev/null

omg-cli prepare-project --env-dir=${env_dir}
