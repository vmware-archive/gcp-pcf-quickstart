#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
workspace_dir="${release_dir}"
omg_dir="${release_dir}/src/omg-cli"
env_dir="${workspace_dir}/env"
env_file="${workspace_dir}/omg-env-in/${env_file_name}"
env_output_dir="${workspace_dir}/omg-env-out"

pushd ${release_dir} > /dev/null
	source ci/tasks/utils.sh
popd > /dev/null
check_param ${env_file_name}

mkdir -p ${env_dir}
pushd ${env_dir}
	tar zxvf ${env_file}
popd

export GOPATH=${workspace_dir}
export PATH=${GOPATH}/bin:${PATH}

terraform_output="${env_dir}/env.json"
terraform_config="${env_dir}/terraform.tfvars"
terraform_state="${env_dir}/terraform.tfstate"

terraform init
terraform get
terraform apply --parallelism=100 -state=${terraform_state} -var-file=${terraform_config} || terraform apply --parallelism=100 -state=${terraform_state} -var-file=${terraform_config}
terraform output -json -state=${terraform_state} > ${terraform_output}

env_file="${env_output_dir}/${env_file_name}"
pushd "${env_dir}"
	tar czvf ${env_file} .
popd