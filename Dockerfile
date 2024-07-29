FROM python:3.12-alpine AS ansible
ARG ANSIBLE_VERSION=10.0
RUN apk -U upgrade --available &&\
    # Required to install ansible pip package, bear in mind to remove those build deps at the end of this RUN directive
    apk add --virtual=build --no-cache --update gcc musl-dev libffi-dev openssl-dev &&\
    # Add here package mandatory to be able to run ansible
    apk add --no-cache openssh-client ca-certificates&&\
    pip install --no-cache-dir --upgrade pip &&\
    pip install --no-cache-dir ansible==${ANSIBLE_VERSION}.* ansible-runner~=2.4 &&\
    # Cleanup package manager cache and remove build deps
    apk del build &&\
    rm -rf /var/cache/apk/*
ENV HOME=/ansible
RUN mkdir -p /ansible && chown 1983:1983 /ansible

FROM ansible AS base
USER 1983

FROM ansible AS gcp
RUN pip install --no-cache-dir requests==2.* google-auth==2.*
USER 1983

FROM ansible AS aws
RUN pip install --no-cache-dir boto3==1.*
USER 1983
