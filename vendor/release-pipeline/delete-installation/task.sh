#!/bin/bash -exu

function main() {
  local cwd
  cwd="${1}"

  local opsman_dns
  opsman_dns="$(cat "${cwd}/environment/name")"

  if [[ -n $OPSMAN_URL_SUFFIX ]]; then
    if [[ -n $USE_OPTIONAL_OPSMAN ]]; then
      opsman_dns="pcf-optional.$opsman_dns.$OPSMAN_URL_SUFFIX"
    else
      opsman_dns="pcf.$opsman_dns.$OPSMAN_URL_SUFFIX"
    fi
  fi

  om --target "https://${opsman_dns}" \
     --username "${OPSMAN_USERNAME}" \
     --password "${OPSMAN_PASSWORD}" \
     --skip-ssl-validation \
     delete-installation
}

main "${PWD}"
