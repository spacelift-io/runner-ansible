FROM python:3.12-alpine AS ansible
ARG ANSIBLE_VERSION=10.0
ENV ANSIBLE_COLLECTIONS_PATH=/usr/share/ansible/collections
ENV ANSIBLE_ROLES_PATH=/usr/share/ansible/roles
RUN apk -U upgrade --available &&\
    # Required to install ansible pip package, bear in mind to remove those build deps at the end of this RUN directive
    apk add --virtual=build --no-cache --update gcc musl-dev libffi-dev openssl-dev &&\
    # Add here package mandatory to be able to run ansible
    apk add --no-cache openssh-client ca-certificates&&\
    pip install --no-cache-dir --upgrade pip &&\
    pip install --no-cache-dir ansible==${ANSIBLE_VERSION}.* ansible-runner~=2.4 &&\
    mkdir -p /usr/share/ansible/collections &&\
    mkdir -p /usr/share/ansible/roles &&\
    # Cleanup package manager cache and remove build deps
    apk del build &&\
    rm -rf /var/cache/apk/*
ENV HOME=/ansible
RUN mkdir -p /ansible && chown 1983:1983 /ansible

FROM ansible AS base
USER 1983

FROM ansible AS gcp
RUN pip install --no-cache-dir requests==2.* google-auth==2.* && \
    ansible-galaxy collection install 'google.cloud:>=1.5.1,<2.0.0' &&\
    # Reset ansible permissions after installing newer collections \
    rm -rf /ansible &&\
    mkdir -p /ansible &&\
    chown 1983:1983 /ansible
USER 1983

FROM ansible AS aws
RUN pip install --no-cache-dir boto3==1.* &&\
    ansible-galaxy collection install 'amazon.aws:>=9.2.0,<10.0.0' &&\
    # Reset ansible permissions after installing newer collections \
    rm -rf /ansible &&\
    mkdir -p /ansible &&\
    chown 1983:1983 /ansible
USER 1983

FROM ansible AS azure
RUN apk add --virtual=build --no-cache gcc musl-dev linux-headers &&\
    # Install azure collection
    ansible-galaxy collection install 'azure.azcollection:>=3.2.0,<4.0.0' &&\
    pip install --no-cache-dir -r /usr/share/ansible/collections/ansible_collections/azure/azcollection/requirements.txt &&\
    pip install --no-cache-dir azure-cli==2.* &&\
    # Reset ansible permissions after installing newer collections \
    rm -rf /ansible &&\
    mkdir -p /ansible &&\
    chown 1983:1983 /ansible &&\
    # Cleanup package manager cache and remove build deps
    apk del build &&\
    rm -rf /var/cache/apk/*
USER 1983
