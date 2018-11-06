#!/bin/bash

slug=$1
version=$(pivnet releases -p $1 --format json | jq -r '.[0].version')

get_pivnet_meta() {
    pivnet product-files --product-slug $1 -r $2 --format json | \
        jq --arg slug "$1" --arg filter "$3" -r 'map(select(.aws_object_key | contains($filter))) |
           map(.release_id = (.["_links"].download.href | capture("/releases/(?<r>[^/]+)").r)) |
           map("\(.name):\nconfig.PivnetMetadata{\n\"\($slug)\",\n\(.release_id),\n\(.id),\n\"\(.sha256)\",\n},") | .[]'
}

echo "$1/$version"
get_pivnet_meta $slug $version ".pivotal"

read -r stemcell_slug stemcell_version <<<$(
    pivnet release-dependencies --product-slug $slug -r $version --format json | \
    jq -r 'map(select(.release.product.slug | contains("stemcells"))) |
        sort_by(- .release.id)[0].release | "\(.product.slug) \(.version)"')

get_pivnet_meta $stemcell_slug $stemcell_version "google"

pivnet product-files --product-slug $stemcell_slug -r $stemcell_version --format json | \
    jq -r 'map(select(.aws_object_key | contains("google")))[0].aws_object_key |
           capture("/(?<n>[^/]+).tgz").n | "\"\(.)\","'
