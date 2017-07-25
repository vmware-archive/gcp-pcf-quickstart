# OMG! Ops Manager on Google

[Design Doc](https://docs.google.com/document/d/1HNZ_rV59DGCyuZqz_gUMvccbBPawQHBh7kqFq1phaLY/edit#heading=h.jgubdjc8el47)

## Checking out omg
```
git clone [..]
git submodules init
git submodules update
```

## TODO: Prerequisites

- Setup DNS
- Quota increases
- Enable APIs

## Deploying and Accessing PCF
```bash
# Deploy
DNS_SUFFIX="..." ./deploy_pcf

# Access jumpbox
ssh -i src/omg-tf/keys/jumpbox_ssh omg@$(terraform output -state src/omg-tf/terraform.tfstate jumpbox_public_ip)

# VPN to internal network
sshuttle -e "ssh -i src/omg-tf/keys/jumpbox_ssh -l omg" -r $(terraform output -state src/omg-tf/terraform.tfstate jumpbox_public_ip) 10.0.0.0/16
```