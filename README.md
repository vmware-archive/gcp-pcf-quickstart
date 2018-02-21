# Pivotal Cloud Foundry on Google Cloud Platform Quickstart 

This repository contains tools and instructions for deploying Pivotal Cloud Foundry (PCF) on Google Cloud Platform (GCP).
The installation includes [PCF 2.0](https://pivotal.io/platform) (Full size or [Small Footprint](https://docs.pivotal.io/pivotalcf/1-12/customizing/small-footprint.html)), [GCP Service Broker](https://docs.pivotal.io/partners/gcp-sb/index.html),
and the [Stackdriver Nozzle](https://docs.pivotal.io/partners/gcp-sdn/index.html).

This project aims to make installing PCF on GCP straightforward and is currently in **beta**. It is **not** an official Google product.

# Getting Started
1. Verify and setup [prerequisites](./docs/prerequisites.md)
1. Setup [DNS](./docs/dns.md) records for your deployment
1. [Deploy PCF](./docs/quick-deployment.md)
1. [Access your Deployment](./docs/login-to-pcf.md)
1. [Deploy an app using Google Cloud Storage and Vision APIs](./docs/deploy-awwvision.md)
1. [Delete your Deployment](./docs/deleting-deployment.md)

# <a name="resources"></a>Related Resources
1. [Pivotal Cloud Foundry on Google Cloud Platform Webseries](https://www.youtube.com/watch?v=TBsc7kiog5Q&list=PLIivdWyY5sqKJ48ycao632rEDuVbFm8yJ)
1. [Pivotal Cloud Foundry on Google Cloud Platform Solution](https://cloud.google.com/solutions/cloud-foundry-on-gcp)
1. [Official PCF Installation Guide](https://docs.pivotal.io/pivotalcf/1-12/customizing/gcp.html)
1. [terraforming-gcp](https://github.com/pivotal-cf/terraforming-gcp)

# Contributing

For details on how to contribute to this project - including filing bug reports and contributing code changes - please see [CONTRIBUTING.md](./CONTRIBUTING.md).
