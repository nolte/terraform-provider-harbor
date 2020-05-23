#!/usr/bin/env bats
cd scripts/test/script
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

