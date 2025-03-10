#!/usr/bin/env ash

# This script will only install collections if the version of ansible is greater than 7.7
# Old versions of ansible use `cert_file` in HTTP calls which was removed in Python 3.12.
# Ansible galaxy collection installs will fail on these versions. So we will skip this step.
# https://github.com/ansible/ansible/pull/80751
collection=$1
if [ "$(python -c "print($ANSIBLE_VERSION > 7.7)")" = "True" ]; then
  ansible-galaxy collection install $collection

  # Reset ansible permissions after installing newer collections
  rm -rf /ansible
  mkdir -p /ansible
  chown 1983:1983 /ansible
else
  echo "Ansible version $ANSIBLE_VERSION does not support ansible galaxy collection installation."
fi