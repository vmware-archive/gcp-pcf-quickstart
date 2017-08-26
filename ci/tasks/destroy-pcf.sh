#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
workspace_dir="$( cd ${release_dir} && cd .. && pwd )"
omg_dir="${release_dir}/src/omg-cli"
env_file="${workspace_dir}/omg-env/${env_file_name}"
env_dir="${workspace_dir}/env"

pushd ${release_dir} > /dev/null
	source ci/tasks/utils.sh
popd > /dev/null
check_param ${env_file_name}


if [ ! -d ${env_dir} ]; then
	mkdir -p ${env_dir}
	pushd ${env_dir}
		tar zxvf ${env_file}
	popd
fi

export GOPATH=${release_dir}
export PATH=${GOPATH}/bin:${PATH}

go install omg-cli
omg-cli remote --env-dir="${env_dir}" "delete-installation"