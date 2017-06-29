#!/bin/bash -exu

function provider_configuration() {
  local metadata
  metadata="${1}"

  local env_name
  env_name="$(cat "${cwd}/ert_metadata/name")"

  local config

  case "${PROVIDER}" in
    "gcp")
      config="$(jq \
        --arg access_key "${STORAGE_INTEROP_ACCESS_KEY}" \
        --arg secret "${STORAGE_INTEROP_SECRET}" \
        --arg smtp_username "${SMTP_USERNAME}" \
        --arg smtp_password "${SMTP_PASSWORD}" \
        '. + {storage_interop_access_key:$access_key} + {storage_interop_secret_key: $secret} + {smtp_username: $smtp_username} + {smtp_password: $smtp_password}' \
        < "${metadata}")"

      if [[ "${ENABLE_C2C}" == "true" ]]; then
        config="$(jq '. + {"enable_container_networking": true}' < <(echo "${config}"))"
      fi
      ;;

    "azure")
      config="$(jq \
        --arg smtp_username "${SMTP_USERNAME}" \
        --arg smtp_password "${SMTP_PASSWORD}" \
        '. + {smtp_username: $smtp_username} + {smtp_password: $smtp_password}' \
        < "${metadata}")"

      if [[ "${ERT_VERSION}" =~ 1.[6-9]  ]]; then
        config="$(jq \
          'with_entries(select(.key != "cf_storage_account_name" and .key != "cf_storage_account_access_key" and .key != "cf_buildpacks_storage_container" and .key != "cf_droplets_storage_container" and .key != "cf_packages_storage_container" and .key != "cf_resources_storage_container"))' \
          < <(echo "${config}"))"
      fi
      ;;

    "aws")
      config="$(jq \
        --arg smtp_username "${SMTP_USERNAME}" \
        --arg smtp_password "${SMTP_PASSWORD}" \
        --arg network_name "${env_name}-ert-network" \
        '. + {smtp_username: $smtp_username} + {smtp_password: $smtp_password} + {network_name: $network_name}' \
        < "${metadata}")"
      ;;

    *)
      echo "unsupported provider ${PROVIDER}"
      exit 1
      ;;
  esac

  printf "%s" "${config}"
}

function resources_configuration() {
  local config

  case "${PROVIDER}" in
    "gcp")
      config='{
        "consul_server": { "instances": 3 },
        "diego_cell": { "instances": 3 },
        "router": { "instances": 3 },
        "mysql": { "instances": 3 },
        "mysql_proxy": { "instances": 2 }
      }'
      ;;

    "azure")
      config='{
        "consul_server": { "instances": 3 },
        "diego_cell": { "instances": 3 },
        "ha_proxy": { "instances": 3 },
        "router": { "instances": 3 },
        "mysql": { "instances": 3 },
        "mysql_proxy": { "instances": 2 }
      }'
      ;;

    "aws")
      config='{
        "consul_server": { "instances": 3 },
        "diego_cell": { "instances": 3 },
        "router": { "instances": 3 },
        "mysql": { "instances": 3 },
        "mysql_proxy": { "instances": 2 }
      }'
      ;;

    *)
      echo "unsupported provider ${PROVIDER}"
      exit 1
      ;;
  esac

  if [[ "${ERT_VERSION}" = "1.8" ]]; then
    config="$(jq '. + {"etcd_server": {"instances": 3}} + {"diego_database": {"instances": 3}}' < <(echo "${config}"))"
  elif [[ "${ERT_VERSION}" = "1.9" ]]; then
    config="$(jq '. + {"etcd_tls_server": {"instances": 3}} + {"etcd_server": {"instances": 1}} + {"diego_database": {"instances": 3}}' < <(echo "${config}"))"
  elif [[ "${ERT_VERSION}" = "1.10" ]] || [[ "${ERT_VERSION}" = "1.11" ]]; then
    config="$(jq '. + {"etcd_tls_server": {"instances": 3}} + {"clock_global": {"instances": 2}} + {"diego_database": {"instances": 2}}' < <(echo "${config}"))"
  fi

  printf "%s" "${config}"
}

function main() {
  local cwd="$1"
  local ops_manager_domain
  local provider_configuration
  local resources
  local saml_enabled

  if [[ "${PROVIDER}" == "vsphere" ]]; then
    ops_manager_domain="$(cat ert_metadata/metadata.json | jq '.metadata.ops_manager_domain')"
    provider_configuration="$(cat ert_metadata/metadata.json | jq '.metadata.provider_configuration')"
    resources="$(cat ert_metadata/metadata.json | jq '.metadata.resources')"
  else
    ops_manager_domain="$(cat "${cwd}/ert_metadata/name")"
    provider_configuration="$(provider_configuration "${cwd}/ert_metadata/metadata")"
    resources="$(resources_configuration)"
  fi

  if [[ "${ENABLE_SAML_CERT}" == "true" ]]; then
    saml_enabled="--enable-saml-cert"
  else
    saml_enabled=""
  fi

  export GOPATH="${cwd}/go"
  pushd "${cwd}/go/src/github.com/pivotal-cf/pcf-releng-ci/tasks/om/generate-ert-configuration/" > /dev/null
    go run main.go \
      --ops-manager-domain "${ops_manager_domain}" \
      --provider "${PROVIDER}" \
      --provider-configuration "${provider_configuration}" \
      --resources "${resources}" \
      --ssl-cert "${SSL_CERT}" \
      --ssl-private-key "${SSL_PRIVATE_KEY}" > "${cwd}/product_configuration/config.json" \
      "${saml_enabled}"
  popd > /dev/null
}

main "${PWD}"
