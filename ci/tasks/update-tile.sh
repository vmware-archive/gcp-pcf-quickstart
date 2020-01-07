#!/bin/bash -e

set -x

if [[ ! -z ${OPS_FILE} ]]; then
  opsfile="$(pwd)/repo/${OPS_FILE}"
  opsfilecommand="--ops-file ${opsfile}"

fi

export RECIPE="$(pwd)/repo/src/omg-cli/templates/assets/deployment.yml ${opsfilecommand}"



get_tile_var () {
    bosh int ${RECIPE} --path /tiles/name=${TILE_NAME}/${1}
}

export PIVNET_PRODUCT_SLUG=$(get_tile_var product/product_slug)
export PIVNET_PRODUCT_GLOB=$(get_tile_var product/file_glob)
export PIVNET_PRODUCT_VERSION=$(jq -r '.Release.Version' tile/metadata.json)
export STEMCELL_VERSION=$(jq -r '.Dependencies | map(select(.Release.Product.Slug | contains("stemcell")))[0].Release.Version'  tile/metadata.json)


pushd repo

echo "Updating tile templates for ${TILE_NAME}/${PIVNET_PRODUCT_VERSION} in ${TILE_BASE_DIR}"
tile-config-generator generate \
                      --include-errands \
                      --do-not-include-product-version \
                      --base-directory=${TILE_BASE_DIR}


echo "Update stemcell to: ${STEMCELL_VERSION} for: ${TILE_NAME}/${PIVNET_PRODUCT_VERSION}"
if [[ ! -z ${OPS_FILE} ]]; then
  # ~1 is a special charater for / https://bosh.io/docs/cli-ops-files/#escaping
  bosh int ${opsfile} -o <(echo -e "
  - type: replace
    path: /path=~1tiles~1name=${TILE_NAME}/value/version
    value: ${PIVNET_PRODUCT_VERSION}
  - type: replace
    path: /path=~1tiles~1name=${TILE_NAME}/value/product/release_version
    value: ${PIVNET_PRODUCT_VERSION}
  - type: replace
    path: /path=~1tiles~1name=${TILE_NAME}/value/stemcell/release_version
    value: '${STEMCELL_VERSION}'
  ") > ${OPS_FILE}.tmp
  mv ${OPS_FILE}{.tmp,}
else
  bosh int ${RECIPE} -o <(echo -e "
  - type: replace
    path: /tiles/name=${TILE_NAME}/version
    value: ${PIVNET_PRODUCT_VERSION}
  - type: replace
    path: /tiles/name=${TILE_NAME}/product/release_version
    value: ${PIVNET_PRODUCT_VERSION}
  - type: replace
    path: /tiles/name=${TILE_NAME}/stemcell/release_version
    value: '${STEMCELL_VERSION}'
  ") > src/omg-cli/templates/assets/deployment.yml.tmp
  mv src/omg-cli/templates/assets/deployment.yml{.tmp,}
fi

git --no-pager diff

echo "Embed updated template files"
go generate src/omg-cli/templates/templates.go
pushd src/omg-cli
UPDATE_FIXTURES=true ginkgo -skipPackage=certification -r ./...
popd

git config --global user.email "ci@starkandwayne.com"
git config --global user.name "CI Bot"

git add -A
git commit -m "Bump tile: ${TILE_NAME}/${PIVNET_PRODUCT_VERSION} stemcell: ${STEMCELL_VERSION}"
