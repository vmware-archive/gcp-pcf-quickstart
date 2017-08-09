# OMG! Ops Manager on Google

TODO: Introduction

## TODO: Prerequisites

- Quota increases
- Enable APIs

## Setup DNS

- Create a new DNS Zone (eg pcf.example.org) in [Cloud DNS](https://console.cloud.google.com/networking/dns/zones)
- Create NS records in parent domain to point to new DNS Zone (eg NS pcf.example.org -> ns-cloud-{a,b,c,d}1.googledomains.com)

## Deploying and Accessing PCF
```bash
# Deploy
DNS_ZONE_NAME="..." ./deploy_pcf

# Access jumpbox
ssh -i env/omg/keys/jumpbox_ssh omg@$(terraform output -state env/omg/terraform.tfstate jumpbox_public_ip)

# VPN to internal network
sshuttle -e "ssh -i src/omg-tf/keys/jumpbox_ssh -l omg" -r $(terraform output -state env/omg/terraform.tfstate jumpbox_public_ip) 10.0.0.0/16
```