# Spacelift Ansible Runner Image

This repository contains the Dockerfile for building our default Ansible runner image.

## Quick Start

To use the Ansible runner image, simply update your [stack settings](https://docs.spacelift.io/concepts/stack/stack-settings#runner-image)
to use `public.ecr.aws/spacelift/runner-ansible` as the runner image for the stack.

## Docker Repository

The image is pushed to the `public.ecr.aws/spacelift/runner-ansible` public repository. It is also pushed to the
`ghcr.io/spacelift-io/runner-ansible` repository as a backup in case of issues with ECR.

Altogether we have 3 flavors of the image:

- `runner-ansible:${ANSIBLE_VERSION}` - built on top of `python:3.12-alpine` base image, with `ansible` and `ansible-runner` installed.
- `runner-ansible:${ANSIBLE_VERSION}-aws` - built on top of `runner-ansible`, with `boto3` installed.
- `runner-ansible:${ANSIBLE_VERSION}-gcp` - built on top of `runner-ansible`, with `google-auth` installed.

Every image is available for the following architectures:

- linux/amd64
- linux/arm64

## Tag Model

This repository create a tag for each minor version of ansible.

In case you don't care about locking the minor version, we also create a tag for the major version that is automatically
bumped when a new minor is released.

You can find below is a non-exhaustive list of tags. This may get outdated with time.

- `10`, `10.2` 
- `10.1`
- `9`, `9.8`
- `9.7`
- `...`

All tags are rebuild every sunday at midnight to be able to get latest security fixes.

## Contributing

The only requirement for working on this repo is a [Docker](https://www.docker.com/) installation.

**ℹ️ Please do not open PR to add a new package to those base images because your workflow need it.**

We want to keep the size of those base image as small as possible 🙏 Only package that are a **strong requirement** to run ansible will be accepted in base images. 

If you need a specific package, please maintain your own version using those image as base image with `FROM public.ecr.aws/spacelift/runner-ansible:10`.

We are open to add new image flavors if needed to support common ansible roles and use cases. The rule of thumb is to keep base image small.

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
