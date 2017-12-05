# Deploying Pivotal Cloud Foundry

Now that you have setup your [prerequisites](./prerequisites.md) and
[DNS](./dns.md) you can deploy Pivotal Cloud Foundry.

## Fetch Tools

Login to your deployment machine and fetch a snapshot of this repository:

```bash
git clone https://github.com/cf-platform-eng/gcp-pcf-quickstart.git
cd gcp-pcf-quickstart
```

## Deploy PCF

Kick off the deployment script and sit back. If a failure occurs follow
any provided instructions and re-run the script. It is safe to repeat.

```bash
export PIVNET_API_TOKEN=<Look for 'API TOKEN' at https://network.pivotal.io/users/dashboard/edit-profile>
./deploy_pcf.sh
```

### Deployment Stages

The installation will perform the following steps:

1. Provision infrastructure with [terraform](https://terraform.io) (5-10 minutes)
1. Configure Ops Manager (<5 minutes)
1. Deploy Pivotal Cloud Foundry (2 hours)

The deployment of Pivotal Cloud Foundry is handled by Ops Manager.
Once it's begun (evident by the streaming of BOSH output) it's safe to
exit the command.

## What's Next?
- [Login to Pivotal Cloud Foundry](login-to-pcf.md) and deploy your first app
- [Delete Deployment](./deleting-deployment.md)
