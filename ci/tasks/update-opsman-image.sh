#!/bin/bash -e

set -x

version=$(cat opsman-tile/version | cut -d# -f1)
url=$(yq r opsman-tile/*.yml us)
fullurl=https://storage.cloud.google.com/${url}

pushd repo
file=src/omg-cli/templates/assets/opsman-image

echo ${fullurl} > ${file}

go generate src/omg-cli/templates/templates.go

git config --global user.email "ci@starkandwayne.com"
git config --global user.name "CI Bot"

git add -A
git commit -m "Bump opsman-image: ${version}"
