# will be configure by env vars from script
provider "harbor" {

}



terraform {

  required_providers {
    harbor = {
      #source  = "terraform.example.com/nolte/harbor"
      source  = "nolte/harbor"
      version = "~> 0.1"
    }
  }
}
