#!/bin/bash -exu

function main() {
  local cwd="$1"

  local opsman_dns
  if [[ -n $USE_OPTIONAL_OPSMAN ]]; then
    opsman_dns="$(jq -r '.optional_ops_manager_dns | rtrimstr(".")' terraform_output/metadata)"
  else
    opsman_dns="$(jq -r '.ops_manager_dns | rtrimstr(".")' terraform_output/metadata)"
  fi

  local product_network
  product_network="$(jq -r '.product.network' product_configuration/config.json)"

  local product_properties
  product_properties="$(jq -r '.product.properties' product_configuration/config.json)"

  local product_config_flags
  product_config_flags=(--product-network "${product_network}")

  if [[ "${product_properties}" != "null" ]]; then
    product_config_flags+=(--product-properties "${product_properties}")
  fi

  local product_resources
  product_resources="$(jq -r '.product.resources' product_configuration/config.json)"

  if [[ "${product_resources}" != "null" ]]; then
    product_config_flags+=(--product-resources "${product_resources}")
  fi

  om --target "https://${opsman_dns}" \
     --skip-ssl-validation \
     --username "${OPSMAN_USERNAME}" \
     --password "${OPSMAN_PASSWORD}" \
     configure-product \
     --product-name "${PRODUCT}" \
     "${product_config_flags[@]}"
}

main "${PWD}"
