# will be configure by env vars from script
provider "harbor" {

}

terraform {
  required_providers {
    harbor = {
      source  = "registry.terraform.private/nolte/harbor"
      version = "~> 0.0.1"
    }
  }
}
