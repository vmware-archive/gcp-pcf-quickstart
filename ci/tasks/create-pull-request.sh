#!/bin/bash -e

set -x

pushd omg-src-in

hub pull-request --base cf-platform-eng/gcp-pcf-quickstart:master -m "Automated Pull Request from the starkandwayne CI"

popd
