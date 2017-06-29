#!/bin/bash -exu

function main() {
  local cwd
  cwd="${1}"

  local opsman_dns
  opsman_dns="$(cat "${cwd}/environment/name")"


  local errands
  errands=$(echo "${ERRANDS_LIST//[\,]/ }")

  for errand in ${errands}; do
    om --target "https://${OPSMAN_URL_PREFIX}.${opsman_dns}.${OPSMAN_URL_SUFFIX}" \
       --skip-ssl-validation \
       --username "${OPSMAN_USERNAME}" \
       --password "${OPSMAN_PASSWORD}" \
       set-errand-state \
       -p "${PRODUCT}" \
       -e "${errand}" \
       --post-deploy-state "${STATE}"
  done
}

main "${PWD}"
