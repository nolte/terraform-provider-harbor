#!/usr/bin/env bats
cd /go/src/github.com/nolte/terraform-provider-harbor/docs/test/script
setup() {
   terraform init
}

teardown() {
   terraform destroy -force
   rm -rf .terraform
   rm -rf terraform.tfstate*
}

@test "Build 1: apply Terraform Script" {
  terraform apply -auto-approve -parallelism=1
}

