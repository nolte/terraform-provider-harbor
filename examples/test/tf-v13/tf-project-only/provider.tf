# will be configure by env vars from script
provider "harbor" {

}


terraform {
  required_providers {
    harbor = {
      #source  = "terraform.example.com/nolte/harbor"
      source  = "test.local/nolte/harbor"
      version = "0.1.6-SNAPSHOT"
    }
  }
}
