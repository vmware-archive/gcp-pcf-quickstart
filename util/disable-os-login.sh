#!/usr/bin/env sh

which gcloud &>/dev/null || (echo error: gcloud is not on your PATH. Install the gcloud cli to continue; echo https://cloud.google.com/sdk/gcloud/; exit 1)

project=$(gcloud config list 2>/dev/null)
if [ $? -ne 0 ]; then
    echo error: cannot read gcloud config
    exit 1
fi

project=$(echo ${project} | grep -o -E 'project = .*' | sed 's/project = //')

echo The following project will be used to disable OS Login:
echo ${project}
read -n 1 -s -r -p "Press any key to continue, or ^C to cancel"
echo

set -x
gcloud compute project-info add-metadata --metadata enable-oslogin=FALSE

