ARG REPOSITORY_BASE_PATH

FROM alpine AS ssm-builder

ARG TARGETARCH

RUN apk add dpkg curl; \
    if [ "$TARGETARCH" = "arm64" ]; then \
     curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_arm64/session-manager-plugin.deb" -o "session-manager-plugin.deb"; \
    else \
     curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/ubuntu_64bit/session-manager-plugin.deb" -o "session-manager-plugin.deb"; \
    fi; \
    dpkg -x session-manager-plugin.deb session-manager-plugin

FROM ${REPOSITORY_BASE_PATH}:latest

COPY --from=ssm-builder /session-manager-plugin/usr/local/sessionmanagerplugin/bin/session-manager-plugin /usr/local/bin/

ARG TARGETARCH

USER root

# session-manager-plugin requires gcompat on amd64
RUN if [ "$TARGETARCH" = "amd64" ]; then \
     apk add --no-cache gcompat; \
    fi

RUN chmod +x /usr/local/bin/session-manager-plugin

USER spacelift

RUN pip install boto3 --break-system-packages
