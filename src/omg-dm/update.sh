#!/bin/bash
set -uex

projectOwner=$(gcloud config get-value account)
envName=${envName}
rootDomain=${rootDomain}

gcloud deployment-manager deployments update omg-${envName} --config=omg.jinja --properties="projectOwner:${projectOwner},rootDomain:${rootDomain},region:${region},zone1:${zone1},zone2:${zone2},zone3:${zone3}"
