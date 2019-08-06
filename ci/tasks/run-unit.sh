#!/usr/bin/env bash

set -e

pushd omg-src-in/ci/tasks > /dev/null
	source utils.sh
	set_resource_dirs
popd > /dev/null

pushd ${omg_dir}
	ginkgo -skipPackage=certification -r .
popd
