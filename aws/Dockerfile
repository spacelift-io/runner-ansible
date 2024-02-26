ARG REPOSITORY_BASE_PATH

# Build SSM
FROM golang:1.15.3-alpine as ssm-builder

ARG VERSION=1.2.279.0

RUN set -ex && apk add --no-cache make git gcc libc-dev curl bash zip && \
    curl -sLO https://github.com/aws/session-manager-plugin/archive/${VERSION}.tar.gz && \
    mkdir -p /go/src/github.com && \
    tar xzf ${VERSION}.tar.gz && \
    mv session-manager-plugin-${VERSION} /go/src/github.com/session-manager-plugin && \
    cd /go/src/github.com/session-manager-plugin && \
    make release

# Add AWS specifics
FROM ${REPOSITORY_BASE_PATH}:latest

USER root

RUN apk add --no-cache gcompat

COPY --from=ssm-builder /go/src/github.com/session-manager-plugin/bin/linux_amd64_plugin/session-manager-plugin /usr/local/bin/

USER spacelift

RUN pip install boto3 botocore --break-system-packages
