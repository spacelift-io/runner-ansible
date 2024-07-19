FROM public.ecr.aws/spacelift/runner-terraform

USER root
WORKDIR /home/spacelift

# Pinning to 7 for now, will introduce more tags a bit later

RUN apk -U upgrade && apk add --no-cache gcc py3-pip python3-dev musl-dev libffi-dev && \
  pip install ansible==10.* ansible-runner==2.* --break-system-packages && \
  apk del python3-dev gcc musl-dev libffi-dev

USER spacelift
