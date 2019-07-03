#!/bin/bash -e

set -x

# pushd repo
url=$(yq r opsman-tile/*.yml us)
fullurl=https://storage.cloud.google.com/${url}

file=repo/src/omg-cli/templates/opsman-image

echo ${fullurl} > ${file}

go generate repo/src/omg-cli/templates/templates.go

cd repo
git status
