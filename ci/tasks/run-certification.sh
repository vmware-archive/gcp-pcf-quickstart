#!/usr/bin/env bash

set -e

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
	extract_env
popd > /dev/null

pushd ${omg_dir}
  [[ $(cat ${env_dir}/config.json | jq .SmallFootprint) == "true" ]] && echo "Skipping certification test on Small Footprint" && exit 0
	ENV_DIR=${env_dir} ginkgo -r certification
popd