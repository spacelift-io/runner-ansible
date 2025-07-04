#!/bin/bash

ECR_REPO="public.ecr.aws/spacelift/runner-ansible"
GITHUB_REPO="ghcr.io/spacelift-io/runner-ansible"


docker pull ${ECR_REPO}:latest
docker pull ${GITHUB_REPO}:latest

ECR_DIGESTS=""
for digest in $(docker manifest inspect ${ECR_REPO}:latest | jq -r '.manifests[].digest'); do
  ECR_DIGESTS+="--amend ${ECR_REPO}@${digest} "
done

docker manifest create ${ECR_REPO}:legacy ${ECR_DIGESTS}
docker manifest push ${ECR_REPO}:legacy

GITHUB_DIGESTS=""
for digest in $(docker manifest inspect ${GITHUB_REPO}:latest | jq -r '.manifests[].digest'); do
  GITHUB_DIGESTS+="--amend ${GITHUB_REPO}@${digest} "
done

docker manifest create ${GITHUB_REPO}:legacy ${GITHUB_DIGESTS}
docker manifest push ${GITHUB_REPO}:legacy
