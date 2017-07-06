#!/bin/bash
set -uex

projectOwner=$(gcloud config get-value account)

gcloud deployment-manager deployments update omg-${envName} --config=omg.jinja --properties=projectOwner:${projectOwner}
