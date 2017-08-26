#!/usr/bin/env bash

set -e


my_dir="$( cd $(dirname $0) && pwd )"
release_dir="$( cd ${my_dir} && cd ../.. && pwd )"
workspace_dir="${release_dir}"
omg_dir="${release_dir}/src/omg-cli"

export GOPATH=${workspace_dir}
export PATH=${GOPATH}/bin:${PATH}

pushd ${release_dir} > /dev/null
	source ci/tasks/utils.sh
popd > /dev/null

check_param ${google_region}
set_gcloud_config

go install omg-cli
omg-cli prepare-project --project-id=${google_project} --region=${google_region}