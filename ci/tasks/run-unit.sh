#!/usr/bin/env bash

set -e

pushd omg-src-ci > /dev/null
	source ci/tasks/utils.sh
	set_resource_dirs
popd > /dev/null

pushd omg-src-in
	ginkgo -skipPackage=certification -r .
popd
