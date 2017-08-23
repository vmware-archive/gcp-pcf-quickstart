#!/usr/bin/env bash

set -e

my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../../.. && pwd )"
workspace_dir="${release_dir}"
test_dir="${release_dir}/ci/docker-image/system_tests"

pushd ${test_dir} > /dev/null
	go get github.com/onsi/gomega
	go get github.com/onsi/ginkgo
	ginkgo -r .
popd > /dev/null