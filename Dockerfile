ARG BASE_IMAGE=python:3.12-alpine
FROM ${BASE_IMAGE} AS ansible
ARG ANSIBLE_VERSION=10.7
ENV ANSIBLE_COLLECTIONS_PATH=/usr/share/ansible/collections
ENV ANSIBLE_ROLES_PATH=/usr/share/ansible/roles
COPY scripts/* /build/
RUN adduser --disabled-password --uid=1983 spacelift &&\
    apk -U upgrade --available &&\
    # Required to install ansible pip package, bear in mind to remove those build deps at the end of this RUN directive
    apk add --virtual=build --no-cache --update gcc musl-dev libffi-dev openssl-dev &&\
    # Add here package mandatory to be able to run ansible
    apk add --no-cache openssh-client ca-certificates &&\
    # Ensure latest setuptools to override any vulnerable system packages
    apk upgrade --available &&\
    pip install --no-cache-dir --upgrade pip &&\
    pip install --no-cache-dir "setuptools>=70.0.0" "PyYAML>=2.2.2" "virtualenv>=20.26.6" "spaceforge>=0.0.2" &&\
    pip install --no-cache-dir ansible==${ANSIBLE_VERSION}.* ansible-runner~=2.4 &&\
    mkdir -p "${ANSIBLE_COLLECTIONS_PATH}" && chown spacelift:spacelift "${ANSIBLE_COLLECTIONS_PATH}" &&\
    mkdir -p "${ANSIBLE_ROLES_PATH}" && chown spacelift:spacelift "${ANSIBLE_ROLES_PATH}" &&\
    # Cleanup package manager cache and remove build deps
    apk del build &&\
    rm -rf /var/cache/apk/*
ENV HOME=/ansible
RUN mkdir -p /ansible && chown spacelift:spacelift /ansible

FROM ansible AS base
USER spacelift

FROM ansible AS gcp
RUN pip install --no-cache-dir requests==2.* google-auth==2.* && \
    /build/install-collection.sh 'google.cloud:>=1.5.1,<2.0.0' &&\
    spaceforge --version
USER spacelift

FROM ansible AS aws
RUN pip install --no-cache-dir boto3==1.* &&\
    /build/install-collection.sh 'amazon.aws:>=9.2.0,<10.0.0' && \
    spaceforge --version
USER spacelift

FROM ansible AS azure
RUN apk add --virtual=build --no-cache gcc musl-dev linux-headers &&\
    # Install azure collection
    /build/install-azure-collection.sh &&\
    pip install --no-cache-dir azure-cli==2.* &&\
    spaceforge --version &&\
    # Cleanup package manager cache and remove build deps
    apk del build &&\
    rm -rf /var/cache/apk/*
USER spacelift
