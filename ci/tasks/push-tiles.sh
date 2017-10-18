#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
	extract_env
popd > /dev/null

go install omg-cli
omg-cli remote --env-dir="${env_dir}" "push-tiles"