#!/usr/bin/env bash

set -e

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
	build_go
    extract_env
popd > /dev/null

pushd ${release_dir}
  omg-cli remote --env-dir="${env_dir}" "deploy"
popd
