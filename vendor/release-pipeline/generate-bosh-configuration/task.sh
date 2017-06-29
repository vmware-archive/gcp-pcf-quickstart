#!/bin/bash -exu

function main() {
  local cwd="$1"

  local env_name
  env_name="$(cat "${cwd}/terraform_output/name")"

  local generate_config_flags
  generate_config_flags+=(--provider "${PROVIDER}")
  generate_config_flags+=(--provider-configuration "$(cat "${cwd}/terraform_output/metadata")")
  generate_config_flags+=(--env-name "${env_name}")

  if [[ -n "${COMPILATION_VM_TYPE}" ]]; then
    generate_config_flags+=(--compilation-vm-type ${COMPILATION_VM_TYPE})
  fi

  export GOPATH="${cwd}/go"
  pushd "${cwd}/go/src/github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-bosh-configuration/" > /dev/null
    go run main.go "${generate_config_flags[@]}" > "${cwd}/bosh_configuration/config.json"
  popd > /dev/null
}

main "${PWD}"
