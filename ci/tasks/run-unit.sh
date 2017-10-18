#!/usr/bin/env bash

set -e

my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
popd > /dev/null

pushd ${omg_dir}
	ginkgo -skipPackage=certification -r .
popd