#!/usr/bin/env bash
# Bash3 Boilerplate. Copyright (c) 2014, kvz.io

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

# Set magic variables for current file & dir
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__file="${__dir}/$(basename "${BASH_SOURCE[0]}")"
__base="$(basename ${__file} .sh)"
__root="$(cd "$(dirname "${__dir}")" && pwd)" # <-- change this as it depends on your app

arg1="${1:-}"

kubectl create ns harbor

helm repo add harbor https://helm.goharbor.io

export INGRESSDOMAIN=192-168-178-51.sslip.io

helm upgrade -i tf-harbor-test harbor/harbor \
    -n harbor \
    --set expose.ingress.hosts.core=harbor.${INGRESSDOMAIN},expose.ingress.hosts.notary=notary.${INGRESSDOMAIN},externalURL=https://harbor.${INGRESSDOMAIN}


kubectl wait --namespace harbor \
  --for=condition=ready pod \
  --selector=component=core \
  --timeout=480s
