# omg/ci

Builds, tests, deploys PCF with OMG. This pipeline creates a full PCF deployment from an empty project. It requires a few out of band resources but is mostly self sufficient. The deployment is destroyed on successful creation but can serve as a guide for automating PCF creation.

## Configuration

See `credentials.yml.tpl`

## Running

```bash
fly -t (your concourse) -p omg set-pipeline -c pipeline.yml -l credentials.yml
```