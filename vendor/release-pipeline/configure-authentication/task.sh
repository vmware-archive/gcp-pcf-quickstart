#!/bin/bash -exu

function main() {
  local cwd
  cwd="${1}"

  local opsman_dns
  opsman_dns="$(cat "${cwd}/environment/name")"

  if [[ -n $OPSMAN_URL_SUFFIX ]]; then
    opsman_dns="pcf.$opsman_dns.$OPSMAN_URL_SUFFIX"
  fi

  om --target "https://${opsman_dns}" \
     --skip-ssl-validation \
     configure-authentication \
     --username "${OPSMAN_USERNAME}" \
     --password "${OPSMAN_PASSWORD}" \
     --decryption-passphrase "${OPSMAN_DECRYPTION_PASSPHRASE}"
}

main "${PWD}"
