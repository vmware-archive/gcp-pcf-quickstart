# omg-tf

## Prerequisites

`gcloud` and `terraform` **0.9.11+** is must be installed on your machine.

You will need to enable the following Google Cloud APIs:
- [Identity and Access Management](https://console.developers.google.com/apis/api/iam.googleapis.com)
- [Cloud Resource Manager](https://console.developers.google.com/apis/api/cloudresourcemanager.googleapis.com/)
- [Cloud DNS](https://console.developers.google.com/apis/api/dns/overview)
- [Cloud SQL API](https://console.developers.google.com/apis/api/sqladmin/overview)
- [Google Compute Engine](https://console.developers.google.com/apis/api/compute_component/overview)

## Prepare Environment

### Deploying

```bash
export PROJECT_ID="your-gcp-project"
export DNS_SUFFIX="cf.your-domain.example.org"
export BASE_IMAGE="https://storage.cloud.google.com/ops-manager-us/pcf-gcp-1.11.4.tar.gz"
export ENV_DIR="$HOME/omg-env/"
mkdir -p ${ENV_DIR}
./init.sh
terraform apply --state=${ENV_DIR}
```

### Destroying
```bash
terraform destroy
rm terraform.tfvars
```

## Connecting to the environment

```bash
sshuttle -e "ssh -i ${ENV_DIR}/keys/jumpbox_ssh" -r $(terraform output jumpbox_public_ip --state=${ENV_DIR}) 10.0.0.0/16
```

## Configuration for omg-cli
```bash
terraform output -json > ${ENV_DIR}/env.json
```

## Appendix

### Var Details
- project: **(required)** ID for your GCP project.
- service_account_key: **(required)** Contents of your service account key file generated using the `gcloud iam service-accounts keys create` command.
- dns_suffix: **(required)** Domain to add environment subdomain to (e.g. foo.example.com)
- ssl_cert: **(required)** SSL certificate for HTTP load balancer configuration. Can be either trusted or self-signed.
- ssl_cert_private_key:  **(required)** Private key for above SSL certificate.
- env_name: *(optional)* An arbitrary unique name for namespacing resources.
- region: *(optional)* Region in which to create resources (e.g. us-central1)
- zones: *(optional)* Zones in which to create resources. Must be within the given region. Currently you must specify exactly 3 Zones for this terraform configuration to work. (e.g. [us-central1-a, us-central1-b, us-central1-c])
- opsman_image_url *(optional)* Source URL of the Ops Manager image you want to boot.
- opsman_storage_bucket_count: *(optional)* Google Storage Bucket for BOSH's Blobstore.

## DNS Records
- pcf.*$env_name*.*$dns_suffix*: Points at the Ops Manager VM's public IP address.
- \*.sys.*$env_name*.*$dns_suffix*: Points at the HTTP/S load balancer in front of the Router.
- doppler.sys.*$env_name*.*$dns_suffix*: Points at the TCP load balancer in front of the Router. This address is used to send websocket traffic to the Doppler server.
- loggregator.sys.*$env_name*.*$dns_suffix*: Points at the TCP load balancer in front of the Router. This address is used to send websocket traffic to the Loggregator Trafficcontroller.
- \*.apps.*$env_name*.*$dns_suffix*: Points at the HTTP/S load balancer in front of the Router.
- \*.ws.*$env_name*.*$dns_suffix*: Points at the TCP load balancer in front of the Router. This address can be used for application websocket traffic.
- ssh.sys.*$env_name*.*$dns_suffix*: Points at the TCP load balancer in front of the Diego brain.
- tcp.*$env_name*.*$dns_suffix*: Points at the TCP load balancer in front of the TCP router.

## Isolation Segments (optional)
- isolation_segment *(optional)* When set to "true" creates HTTP load-balancer across 3 zones for isolation segments.
- iso_seg_ssl_cert: *(optional)* SSL certificate for HTTP load balancer configuration. Can be either trusted or self-signed.
- iso_seg_ssl_cert_private_key:  *(optional)* Private key for above SSL certificate.

## Cloud SQL Configuration (optional)
- external_database: *(optional)* When set to "true", a cloud SQL instance will be deployed for the Ops Manager and ERT.
