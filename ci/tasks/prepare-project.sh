#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
workspace_dir="${release_dir}"
omg_dir="${release_dir}/src/omg-cli"

export GOPATH=${workspace_dir}
export PATH=${GOPATH}/bin:${PATH}
go install omg-cli

pushd ${release_dir} > /dev/null
	source ci/tasks/utils.sh
popd > /dev/null

check_param 'env_config'
set_gcloud_config


mkdir -p env/
echo "${env_config}" > env/config.json

omg-cli prepare-project --env-dir=$PWD/env