#!/usr/bin/env bash
#
# This can be also run with ``
# ./tst-15-execute-classic-acc "/api/v2.0"

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
#export TF_LOG=TRACE
# shellcheck disable=SC2046
CURRENTDIR=$(pwd)
cd ./examples/test/
go test --tags=integration -timeout 20m -v
cd  "${CURRENTDIR}"
