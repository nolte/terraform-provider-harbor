#!/usr/bin/env bash
#
# This can be also run with `` 
# ./tst-01-prepare-harbor.sh "10-42-0-100.sslip.io" "1.2.0"

set -e
set -o pipefail
set -o nounset

INGRESS_DOMAIN=${1:-"192-168-178-51.sslip.io"}
HARBOR_CHART_VERSION=${2:-"1.4.0"}

kubectl create ns harbor || true

helm repo add harbor https://helm.goharbor.io 

helm upgrade -i tf-harbor-test harbor/harbor \
    --version "$HARBOR_CHART_VERSION" \
    -n harbor \
    --set expose.ingress.hosts.core="harbor.${INGRESS_DOMAIN}" \
    --set expose.ingress.hosts.notary="notary.${INGRESS_DOMAIN}" \
    --set externalURL="https://harbor.${INGRESS_DOMAIN}"


kubectl wait --namespace harbor \
  --for=condition=ready pod \
  --selector=component=core \
  --timeout=480s
