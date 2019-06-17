#!/usr/bin/env bash

set -e

pushd omg-src-in > /dev/null
	source utils.sh
	set_resource_dirs
popd > /dev/null

pushd omg-src-in
	ginkgo -skipPackage=certification -r .
popd
