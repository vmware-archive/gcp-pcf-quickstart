#!/usr/bin/env bash

set -e


pushd omg-src-in > /dev/null
	source utils.sh
	set_resource_dirs
	build_go
	extract_env
popd > /dev/null

pushd ${release_dir}
  omg-cli review-eulas --env-dir="${env_dir}" --accept-all
popd
