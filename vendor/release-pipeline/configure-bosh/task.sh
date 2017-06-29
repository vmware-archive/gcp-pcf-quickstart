#!/bin/bash -exu

function main() {
  local cwd="$1"

  local opsman_dns
  if [[ -n $USE_OPTIONAL_OPSMAN ]]; then
    opsman_dns="$(jq -r '.optional_ops_manager_dns | rtrimstr(".")' terraform_output/metadata)"
  else
    opsman_dns="$(jq -r '.ops_manager_dns | rtrimstr(".")' terraform_output/metadata)"
  fi

  local bosh_config_flags

  local iaas_configuration
  iaas_configuration="$(jq -r '.iaas_configuration' bosh_configuration/config.json)"
  if [[ "${iaas_configuration}" != "null" ]]; then
    bosh_config_flags+=(--iaas-configuration "${iaas_configuration}")
  fi

  local director_configuration
  director_configuration="$(jq -r '.director_configuration' bosh_configuration/config.json)"
  if [[ "${director_configuration}" != "null" ]]; then
    bosh_config_flags+=(--director-configuration "${director_configuration}")
  fi

  local az_configuration
  az_configuration="$(jq -r '.az_configuration' bosh_configuration/config.json)"
  if [[ "${az_configuration}" != "null" ]]; then
    bosh_config_flags+=(--az-configuration "${az_configuration}")
  fi

  local networks_configuration
  networks_configuration="$(jq -r '.networks_configuration' bosh_configuration/config.json)"
  if [[ "${networks_configuration}" != "null" ]]; then
    bosh_config_flags+=(--networks-configuration "${networks_configuration}")
  fi

  local network_assignment
  network_assignment="$(jq -r '.network_assignment' bosh_configuration/config.json)"
  if [[ "${network_assignment}" != "null" ]]; then
    bosh_config_flags+=(--network-assignment "${network_assignment}")
  fi

  local resource_configuration
  resource_configuration="$(jq -r '.resource_configuration' bosh_configuration/config.json)"
  if [[ "${resource_configuration}" != "null" ]]; then
    bosh_config_flags+=(--resource-configuration "${resource_configuration}")
  fi

  om --target "https://${opsman_dns}" \
     --skip-ssl-validation \
     --username "${OPSMAN_USERNAME}" \
     --password "${OPSMAN_PASSWORD}" \
     configure-bosh \
     "${bosh_config_flags[@]}"
}

main "${PWD}"
