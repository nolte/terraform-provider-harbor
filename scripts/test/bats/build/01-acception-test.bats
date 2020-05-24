#!/usr/bin/env bats
cd scripts/test/tf-acception-test
HARBOR_ENDPOINT=$(kubectl get Ingress tf-harbor-test-harbor-ingress -n harbor -ojson | jq '.spec.rules[].host' -r | grep harbor)

setup() {
   terraform init
}

#teardown() {
#   terraform destroy -force -var harbor_endpoint=${HARBOR_ENDPOINT} -var harbor_base_path='/api'
##   rm -rf .terraform
##   rm -rf terraform.tfstate*
#}

@test "Build 1: apply Terraform Script" {
  echo "Start test ${HARBOR_ENDPOINT}"
  terraform apply -auto-approve -parallelism=1 -var harbor_endpoint=${HARBOR_ENDPOINT} -var harbor_base_path='/api'
}
