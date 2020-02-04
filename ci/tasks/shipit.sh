#!/bin/bash

set -e
set -x

VERSION=$(cat version/version)

echo "v${VERSION}"           > release/tag
echo "GCP v${VERSION}"       > release/name

export IMAGE_VERSION=$(cat omg-src-develop/src/omg-cli/templates/assets/opsman-image | cut -d- -f5)

spruce json omg-src-develop/src/omg-cli/templates/assets/deployment.yml \
    | jq -r '.tiles | map("- \(.name)/\(.version) (\(.stemcell.product_slug | split("-") | .[-1])/\(.stemcell.release_version))") | .[]' > release/notes.md
echo "- ops-manager/${IMAGE_VERSION}" >> release/notes.md
