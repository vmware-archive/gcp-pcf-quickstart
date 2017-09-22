#!/usr/bin/env bash

set -e

my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
workspace_dir="$( cd ${release_dir} && cd .. && pwd )"

export GOPATH=${release_dir}
export PATH=${GOPATH}/bin:${PATH}

pushd ${release_dir} > /dev/null
	source ci/tasks/utils.sh
popd > /dev/null

# Version info
semver_version=`cat release-version-semver/number`
echo $semver_version > promoted/semver_version

today=$(date +%Y-%m-%d)
cp -r omg-src-in promoted/repo

pushd promoted/repo
  # generate versioned file
  SEMVER=${semver_version} go generate src/omg-cli/version/version.go

  git diff | cat
  git add .

  git config --global user.email cf-bosh-eng@pivotal.io
  git config --global user.name CI
  git commit -m "tag: release v${semver_version}"
popd
