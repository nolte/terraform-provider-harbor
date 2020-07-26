#!/usr/bin/env bash
#
# This can be also run with ``
# ./tst-01-prepare-harbor.sh "10-42-0-100.sslip.io" "1.2.0"

set -e
set -o pipefail
set -o nounset

echo "==> Starting the Harbor installation..."

INGRESS_DOMAIN=${1:-"172-17-0-1.sslip.io"}
HARBOR_CHART_VERSION=${2:-"1.3.2"}

kubectl create ns harbor || true

helm repo add harbor https://helm.goharbor.io

helm upgrade -i tf-harbor-test harbor/harbor \
    --version "$HARBOR_CHART_VERSION" \
    -n harbor \
    --set expose.ingress.hosts.core="harbor.${INGRESS_DOMAIN}" \
    --set expose.ingress.hosts.notary="notary.${INGRESS_DOMAIN}" \
    --set externalURL="https://harbor.${INGRESS_DOMAIN}"

echo "==> Waiting for ready Harbor pods..."

kubectl wait --namespace harbor \
    --for=condition=ready pod \
    --selector=app=harbor \
    --timeout=680s

sleep 2

echo "==> Harbor Ready for use"
