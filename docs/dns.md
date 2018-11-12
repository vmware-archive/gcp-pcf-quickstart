# Setting up DNS

A registered domain name is required to deploy PCF. If you don't already
have a domain name, you can create and register a new domain name at
[Google Domains](https://domains.google), or you can use a third-party
domain name registrar.

The installation requires a child zone to be hosted on [Cloud DNS](https://cloud.google.com/dns/docs/)

## Create a Cloud DNS Zone

### Using the UI

1. Open [Cloud DNS](https://console.cloud.google.com/net-services/dns/zones) for your GCP Project
1. Click 'Create Zone' and populate the following fields
   - **Zone Name**: `pcf-zone`
   - **DNS Name**: `pcf.<your-domain-name.com>`
1. Click 'Create'

### Or Using `gcloud`

Run the following, replacing `<your-domain-name.com>`:

```bash
gcloud dns managed-zones create pcf-zone --description="PCF Zone" --dns-name="pcf.<your-domain-name.com>"
```

## Setup NS records

This process will differ by DNS host and you should consult the documentation
for your provider to create the record.

The follow record must be created:

- **DNS Name**: `pcf.<your-domain-name.com>`
- **Record Type**: `NS`
- **Name Server**: Populate with the rrdatas from: `gcloud dns record-sets list --zone=pcf-zone --format=flattened --filter="type=NS"`.
  They will be in the format of `ns-cloud-{a,b,c,..}{1,2,3..}.googledomains.com`

## Verifying NS records

Verify your NS records are properly configured by running dig:

```bash
dig pcf.<your-domain-name.com> NS +short
```

If the command returns your list on Cloud DNS servers (`ns-cloud-..`)
then you have successfuly setup your zone. If it returns an empty result
then you may need to wait for the record to propagate or the NS record may
be incorrectly configured.

**> After verifying your DNS records you can now move on to [Deploying Pivotal Cloud Foundry](./quick-deployment.md)**
