output "harbor_project_id" {
  value = harbor_project.main.id
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
#
