#!/bin/bash

set -e
set -x

VERSION=$(cat version/version)

echo "v${VERSION}"           > release/tag
echo "GCP v${VERSION}"       > release/name

export IMAGE_VERSION=$(cat omg-src-develop/src/omg-cli/templates/assets/opsman-image | cut -d- -f5)

path="omg-src-develop/src/omg-cli"

bosh int ${path}/templates/assets/deployment.yml -o ${path}/templates/assets/options/healthwatch.yml | spruce json \
    | jq -r '.tiles | map("- \(.name)/\(.version) (\(.stemcell.product_slug | split("-") | .[-1])/\(.stemcell.release_version))") | .[]' > release/notes.md
echo "- ops-manager/${IMAGE_VERSION}" >> release/notes.md
