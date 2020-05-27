#!/usr/bin/env bash
#
# This can be also run with ``
# ./tst-01-prepare-harbor.sh "10-42-0-100.sslip.io" "1.2.0"

set -e
set -o pipefail
set -o nounset

export TF_ACC=1
export HARBOR_ENDPOINT="$(kubectl get Ingress tf-harbor-test-harbor-ingress -n harbor -ojson | jq '.spec.rules[].host' -r | grep harbor)"
export HARBOR_USERNAME=admin
export HARBOR_PASSWORD="$(kubectl -n harbor get secret tf-harbor-test-harbor-core -ojson | jq '.data.HARBOR_ADMIN_PASSWORD' -r | base64 -d)"
export HARBOR_INSECURE="true"
go test -timeout 20m $(go list /go/src/github.com/nolte/terraform-provider-harbor/... |grep -v 'vendor') -v
