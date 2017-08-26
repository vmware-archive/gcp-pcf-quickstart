# ci/docker-image

Builds the docker image used by the CI environment to run OMG.

## Configuration

See `credentials.yml.tpl`

## DockerHub Repositories

Two DockerHub repositories are required to run the pipeline. The first (docker_hub_test_repo) is used to test the image. The second (docker_hub_repo) is the final published image.

## Running

```bash
fly -t (your concourse) -p docker-image set-pipeline -c pipeline.yml -l credentials.yml
```

## Validation

The test suite in [system_tests](./system_tests) is used to verify the image contains the needed tools.