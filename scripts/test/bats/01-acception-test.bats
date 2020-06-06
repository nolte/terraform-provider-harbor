#!/usr/bin/env bats
cd scripts/test/tf-acception-test || exit

# used for terraform provider config
export HARBOR_ENDPOINT=$(kubectl get Ingress tf-harbor-test-harbor-ingress -n harbor -ojson | jq '.spec.rules[].host' -r | grep harbor)
export HARBOR_USERNAME=admin
export HARBOR_PASSWORD=Harbor12345
export HARBOR_BASEPATH="/api/v2.0"
export HARBOR_INSECURE="true"

setup() {
   terraform init
}

teardown() {
   terraform destroy -force
   rm -rf .terraform
   rm -rf terraform.tfstate*
}

@test "Build 1: apply Terraform Script" {
  echo "Start test ${HARBOR_ENDPOINT}"
  terraform apply -auto-approve -parallelism=1
}

