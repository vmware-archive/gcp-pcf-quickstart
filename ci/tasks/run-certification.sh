#!/usr/bin/env bash

set -e

pushd omg-src-in/ci/tasks > /dev/null
	source utils.sh
	set_resource_dirs
	extract_env
popd > /dev/null

pushd omg-src-in
  [[ $(cat ${env_dir}/config.json | jq .SmallFootprint) == "true" ]] && echo "Skipping certification test on Small Footprint" && exit 0
	ENV_DIR=${env_dir} ginkgo -r certification
popd
