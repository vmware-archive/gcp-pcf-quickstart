#!/bin/bash
set -uex

envName=$(date +%s)
projectOwner=$(gcloud config get-value account)

gcloud deployment-manager deployments create omg-${envName} --config=omg.jinja --properties=projectOwner:${projectOwner}
