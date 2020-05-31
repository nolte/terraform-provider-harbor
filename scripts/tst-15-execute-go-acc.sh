#!/usr/bin/env bash
#
# This can be also run with ``
# ./tst-01-prepare-harbor.sh "10-42-0-100.sslip.io" "1.2.0"

set -e
set -o pipefail
set -o nounset

echo "==> Starting the acception tests..."

export TF_ACC=1
HARBOR_ENDPOINT="$(kubectl get Ingress tf-harbor-test-harbor-ingress -n harbor -ojson | jq '.spec.rules[].host' -r | grep harbor)"
export HARBOR_ENDPOINT
export HARBOR_USERNAME=admin
HARBOR_PASSWORD="$(kubectl -n harbor get secret tf-harbor-test-harbor-core -ojson | jq '.data.HARBOR_ADMIN_PASSWORD' -r | base64 -d)"
export HARBOR_PASSWORD
export HARBOR_INSECURE="true"
export HARBOR_BASEPATH=${1:-"/api"}

# shellcheck disable=SC2046
go test -timeout 20m $(go list /go/src/github.com/nolte/terraform-provider-harbor/... | grep -v 'vendor') -v
