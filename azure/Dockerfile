ARG REPOSTORY_BASE_PATH

FROM ${REPOSITORY_BASE_PATH}:latest

# Capture default user. We need to elevate to install packages, but need to de-elevate after
RUN id -u -> /tmp/uid

USER root

RUN apk add gcc python3-dev musl-dev linux-headers --break-system-packages

RUN pip install azure-cli

# Switch back to the default user
USER $(cat /tmp/uid)
