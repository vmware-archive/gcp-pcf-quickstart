#!/usr/bin/env bash

set -e
my_dir="$( cd $(dirname $0) && pwd )"
pushd ${my_dir} > /dev/null
	source utils.sh
	set_resource_dirs
    extract_env
popd > /dev/null

go install omg-cli

set +e
n=0
until [ $n -ge 5 ]
do
  omg-cli remote --env-dir="${env_dir}" "delete-installation" && break
  n=$[$n+1]
  delay=180
  echo "delete failed, trying again in ${delay} seconds"
  sleep ${delay}
done
