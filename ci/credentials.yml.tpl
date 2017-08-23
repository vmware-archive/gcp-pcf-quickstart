---
# Source
source_uri:
source_branch:
source_username:
source_password:

# CI data
ci_json_key_date: |             # service account with access to the omg source and GCS storage
ci_bucket_name:	                # GCS bucket for CI to store data

# Target project for deploying OMG
google_project:                   # Project to deploy to
google_json_key_data:             # JSON key to a service account with owner access to the google_project
google_region: us-east1           # Change not currently supported
env_file_name: env.tgz
env_name: omg-ci
PIVNET_API_TOKEN:                 # API Token to network.pivotal.io account (more info: https://network.pivotal.io/docs/api#how-to-authenticate)
PIVNET_ACCEPT_EULA:               # "yes" to accept all network.pivotal.io EULAs
DNS_ZONE_NAME:                    # existing cloud DNS zone
BASE_IMAGE_URL: https://storage.cloud.google.com/ops-manager-us/pcf-gcp-1.11.4.tar.gz # URL to base image of Ops Manager
