FROM public.ecr.aws/spacelift/runner-terraform

USER root
WORKDIR /home/spacelift

RUN apk -U upgrade && apk add --no-cache \
    ansible \
    py3-pip

RUN pip install ansible-runner

USER spacelift
