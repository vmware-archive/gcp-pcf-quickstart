#!/bin/bash -e

set -x

export GITHUB_USER=${GITHUB_USER}
export GITHUB_PASSWORD=${GITHUB_PASSWORD}


pushd omg-src-in

hub pull-request --base cf-platform-eng/gcp-pcf-quickstart:master --head starkandwayne/gcp-pcf-quickstart:develop -m "Automated Pull Request from the starkandwayne CI"

popd
