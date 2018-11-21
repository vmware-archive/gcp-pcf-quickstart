#!/usr/bin/env bash

export GO111MODULE=on # manually active module mode

check_param() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [[ "$value" == 'replace-me' ]] || [[ "$value" == '' ]]; then
    echo "environment variable $name must be set"
    exit 1
  fi
}

print_git_state() {
  echo "--> last commit..."
  TERM=xterm-256color git log -1
  echo "---"
  echo "--> local changes (e.g., from 'fly execute')..."
  TERM=xterm-256color git status --verbose
  echo "---"
}

declare -a on_exit_items
on_exit_items=()

function on_exit {
  echo "Running ${#on_exit_items[@]} on_exit items..."
  for i in "${on_exit_items[@]}"
  do
    for try in $(seq 0 9); do
      sleep $try
      echo "Running cleanup command $i (try: ${try})"
        eval $i || continue
      break
    done
  done
}

function add_on_exit {
  local n=${#on_exit_items[@]}
  on_exit_items=("${on_exit_items[@]}" "$*")
  if [[ $n -eq 0 ]]; then
    trap on_exit EXIT
  fi
}

set_gcloud_config() {
	check_param 'google_project'
	check_param 'google_json_key_data'

	gcloud config set project $google_project

	key=$(mktemp)
	echo $google_json_key_data > ${key}
    export GOOGLE_APPLICATION_CREDENTIALS=${key}
	gcloud auth activate-service-account --key-file=${key}
}

build_go() {
  my_dir=$(dirname "$(readlink -f "$0")")
  release_dir="$(realpath ${my_dir}/../..)"
  omg_dir="${release_dir}/src/omg-cli"

  pushd ${omg_dir}
    go build -o $release_dir/bin/omg-cli
    export PATH=$release_dir/bin:$PATH
  popd
}

set_resource_dirs() {
    my_dir=$(dirname "$(readlink -f "$0")")

    export release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
    export workspace_dir="$( cd ${release_dir} && cd .. && pwd )"
    export omg_dir="${release_dir}/src/omg-cli"
    export env_file="${workspace_dir}/omg-env-in/${env_file_name}"
    export env_dir="${workspace_dir}/env"
    export env_output_dir="${workspace_dir}/omg-env-out"
    export env_output_file="${env_output_dir}/${env_file_name}"
    export terraform_output="${env_dir}/terraform_output.json"
    export terraform_config="${env_dir}/terraform.tfvars"
    export terraform_state="${env_dir}/terraform.tfstate"
}

extract_env() {
    check_param 'env_file_name'

    # This task may run as part of a failed job. In that case
    # the env_dir will already exist and contain the state.
    if [[ ! -d ${env_dir} ]]; then
        mkdir -p ${env_dir}
        pushd ${env_dir}
            tar zxvf ${env_file}
        popd
    fi
}

generate_env_config() {
    check_param 'env_config'

    mkdir -p ${env_dir}
    echo "${env_config}" > "${env_dir}/config.json"
}

configure_terraform_backend() {
    check_param "terraform_state_bucket"
    check_param "env_file_name"
    cat << EOF > state.tf
terraform {
  backend "gcs" {
    bucket = "${terraform_state_bucket}"
    prefix  = "terraform/${env_file_name%.*}"
  }
}
EOF
}

function save_terraform_state {
	pushd "${env_dir}"
		tar czvf ${env_output_file} .
	popd
}
