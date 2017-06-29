#!/bin/bash
set -uex

envName=$(date +%s)

gcloud deployment-manager deployments create omg-${envName} --config=omg.jinja