---
# Source
#
# Parameters for a git repository with the OMG source code,
# optionally secured TTP basic auth.
#
# Using Google Cloud Source Repositories:
#  source_uri: https://source.developers.google.com/p/<project>/r/<repo name>
#  source_branch: master
#
# For username/password go to https://source.developers.google.com/auth/start?scopes=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fcloud-platform
# Use the fields labeled 'This is your Git username', 'This is your Git password'
source_uri:
source_branch:
source_username:
source_password:

# CI data
#
# JSON key for a GCP Service Account with Storage Object Admin access to the `ci_bucket_name` GCS bucket
ci_json_key_date: |
   {
   ... key here ..
   }
# GCS bucket for CI to store data. Bucket must have Object Versioning enabled. To enable: gsutil versioning set on gs://[bucket_name]
ci_bucket_name: replace-me

# Target environment to deploy OMG
google_project:
# JSON key for a GCP Service Account with Owner access to the `google_project`
google_json_key_data:

# Change not currently supported
google_region: us-east1

# Optionally provide unique names to deploy multiple envs to the same project
env_file_name: env.tgz
env_name: omg-ci

# API Token to network.pivotal.io account (more info: https://network.pivotal.io/docs/api#how-to-authenticate)
PIVNET_API_TOKEN:
# "yes" to accept all network.pivotal.io EULAs
PIVNET_ACCEPT_EULA:
# Existing Cloud DNS Zone to use for PCF deployment
DNS_ZONE_NAME: omg-zone
# URL to base image of Ops Manager
BASE_IMAGE_URL: https://storage.cloud.google.com/ops-manager-us/pcf-gcp-1.11.4.tar.gz
