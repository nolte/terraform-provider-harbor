module "base_tests" {
    source = "../../../tf-acception-test"
    providers = {
        harbor = harbor
    }
}
output "harbor_project_id" {
  value = module.base_tests.harbor_project_id
}
