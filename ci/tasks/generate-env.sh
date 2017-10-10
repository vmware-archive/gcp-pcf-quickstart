#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
workspace_dir="$( cd ${release_dir} && cd .. && pwd )"
omg_tf_dir="${release_dir}/src/omg-tf"
env_output_dir="${workspace_dir}/omg-env-out"

export GOPATH=${release_dir}
export PATH=${GOPATH}/bin:${PATH}

pushd ${release_dir} > /dev/null
	source ci/tasks/utils.sh
popd > /dev/null

check_param 'google_project'
check_param 'google_json_key_data'
check_param 'env_config'

check_param 'PIVNET_API_TOKEN'
check_param 'PIVNET_ACCEPT_EULA'

set_gcloud_config

export ENV_DIR="${workspace_dir}/env"

mkdir -p ${ENV_DIR}
echo "${env_config}" > "${ENV_DIR}/config.json"

go install omg-cli
set -o allexport
eval $(omg-cli source-config --env-dir="${ENV_DIR}")
set +o allexport

pushd ${omg_tf_dir}
	./init.sh
popd

env_file="${env_output_dir}/${env_file_name}"
pushd "${ENV_DIR}"
	tar czvf ${env_file} .
popd