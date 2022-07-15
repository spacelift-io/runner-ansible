# Spacelift Ansible Runner Image

This repository contains the Dockerfile for building our default Ansible runner image.

## Quick Start

To use the Ansible runner image, simply update your [stack settings](https://docs.spacelift.io/concepts/stack/stack-settings#runner-image)
to use `public.ecr.aws/spacelift/runner-ansible` as the runner image for the stack.

## Docker Repository

The image is pushed to the `public.ecr.aws/spacelift/runner-ansible` public repository. It is also pushed to the
`ghcr.io/spacelift-io/runner-ansible` repository as a backup in case of issues with ECR.

## Branch Model

This repository uses two main branches:

- `main` - contains the production version of the runner image.
- `future` - used to test development changes.

Pushes to main deploy to the latest tag, whereas pushes to future deploy to the future tag. This
means that to use the development version you can use the `public.ecr.aws/spacelift/runner-ansible:future` image.

## Development

The only requirement for working on this repo is a [Docker](https://www.docker.com/) installation.

### Testing a new Image

The easiest way to test a new image before opening a pull request is to push it to your own
Docker repository and then update a test stack to use your custom image. The following steps
explain the process using Docker Hub, but any other public container registry can be used.

First, sign-up for an account at [Docker Hub](https://hub.docker.com/), and login via `docker login`:

```shell
docker login
```

Next, build the image using your Docker Hub username to tag it. For example, if your username
is `abc123`, you would use the following command to build the image:

```shell
docker build -t abc123/runner-ansible:latest .
```

Once the build has completed, push your changes:

```shell
docker push abc123/runner-ansible:latest
```

Congratulations! You can now update your stack to use `abc123/runner-ansible:latest` as its
runner image.
