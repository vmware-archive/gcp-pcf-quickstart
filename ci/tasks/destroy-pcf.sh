#!/usr/bin/env bash

set -e
my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
	set_gcloud_config
  extract_env
	build_go
popd > /dev/null

omg-cli remote --env-dir="${env_dir}" "delete-installation" && exit 0
echo "delete failed, cleaning project instead"
omg-cli cleanup-project --env-dir="${env_dir}" --no-dry-run
