#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../../.. && pwd )"
workspace_dir="${release_dir}"
test_dir="${my_dir}/system_tests"

export GOPATH=${workspace_dir}
export PATH=${GOPATH}/bin:${PATH}

pushd ${test_dir} > /dev/null
	ginkgo -r .
popd > /dev/null