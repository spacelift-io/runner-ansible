#!/usr/bin/env ash

# This script will only install azure collection if the version of ansible is greater than 7.7 or less than 11.3
# It will install the azure collection requirements for ansible versions newer than 11.2

# Old versions of ansible use `cert_file` in HTTP calls which was removed in Python 3.12.
# Ansible galaxy collection installs will fail on these versions, so we dont install the azure collection there.
# https://github.com/ansible/ansible/pull/80751

# 11.3 includes azure.azcollection at 3.2.0 already so this step isn't needed for any version 11.3 or higher.
if [[ "$(python -c "print($ANSIBLE_VERSION > 7.7)")" == "True" && "$(python -c "print($ANSIBLE_VERSION < 11.3)")" == "True" ]]; then
  ansible-galaxy collection install 'azure.azcollection:>=3.2.0,<4.0.0'

  # Remove the uamqp dependency as its known to not build correctly. This package is no longer maintained by the
  # Azure team and older ansible versions do not require iot-hub (which depends on uamqp) so we can remove it
  # from the requirements.
  # https://github.com/Azure/azure-uamqp-python/issues/386
  sed -i '/azure-iot-hub==2.6.1;platform_machine=="x86_64"/d' /usr/share/ansible/collections/ansible_collections/azure/azcollection/requirements.txt

  # Install the requirements for the azure collection
  # this is a documented required step for the azure collection
  pip install --no-cache-dir -r /usr/share/ansible/collections/ansible_collections/azure/azcollection/requirements.txt

  # Reset ansible permissions after installing newer collections
  rm -rf /ansible
  mkdir -p /ansible
  chown 1983:1983 /ansible

  chown -R 1983:1983 "${ANSIBLE_COLLECTIONS_PATH}"
else
  if [[ "$(python -c "print($ANSIBLE_VERSION >= 11.3)")" == "True" ]]; then
    # Ansible 11.3 and higher already includes the azure.azcollection at version 3.2.0 but it does not include the pip requirements.
    # The requirements are also at a different path for this version.
    sed -i '/azure-iot-hub==2.6.1;platform_machine=="x86_64"/d' /usr/local/lib/python3.12/site-packages/ansible_collections/azure/azcollection/requirements.txt
    pip install --no-cache-dir -r /usr/local/lib/python3.12/site-packages/ansible_collections/azure/azcollection/requirements.txt
  else
    echo "Ansible version $ANSIBLE_VERSION either does not support ansible galaxy collection installation or the version of the azure collection is already installed."
  fi
fi
