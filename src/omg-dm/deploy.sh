#!/bin/bash
set -uex

envName=$(date +%s)
projectOwner=$(gcloud config get-value account)

echo "export envName=${envName}" >> .envrc
direnv allow || true

gcloud deployment-manager deployments create omg-${envName} --config=omg.jinja --properties=projectOwner:${projectOwner}
