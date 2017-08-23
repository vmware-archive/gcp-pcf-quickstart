#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
workspace_dir="${release_dir}"
omg_tf_dir="${release_dir}/src/omg-tf"

export GOPATH=${workspace_dir}
export PATH=${GOPATH}/bin:${PATH}

pushd ${release_dir} > /dev/null
	source ci/tasks/utils.sh
popd > /dev/null

check_param ${google_region}
set_gcloud_config

export ENV_DIR="${workspace_dir}/env"
mkdir -p ${ENV_DIR}

pushd ${omg_tf_dir}
	./init.sh
popd

env_file="${workspace_dir}/omg-env-out/${env_file_name}"
pushd ${ENV_DIR}
	tar czvf ${env_file} .
popd