---
# Source
source_uri:
source_branch:
source_username:
source_password:

# CI 
ci_json_key_date: |             # service account with access to the omg source and GCS storage
ci_bucket_name:	                # GCS bucket for CI to store data

# Target project for deploying OMG
google_project:
google_json_key_data:
google_region:
env_file_name: env.tgz
env_name: omg-ci
PIVNET_API_TOKEN:
PIVNET_ACCEPT_EULA:               # "yes" to accept all network.pivotal.io EULAs
DNS_ZONE_NAME:                    # existing cloud DNS zone
BASE_IMAGE_URL:                   # URL to base image of Ops Manager
