
variable "harbor_endpoint" {
  default = "demo.goharbor.io"
}
variable "harbor_base_path" {
  default = "/api/v2.0"
}

provider "harbor" {
  host     = var.harbor_endpoint
  schema   = "https"
  insecure = true
  basepath = var.harbor_base_path
  username = "admin"
  password = "Harbor12345"
}

resource "harbor_project" "main" {
  name                   = "main"
  public                 = false # (Optional) Default value is false
  vulnerability_scanning = true  # (Optional) Default vale is true. Automatically scan images on push
}

resource "harbor_robot_account" "account" {
  name        = "myrobot"
  description = "Robot account used to push images to harbor"
  project_id  = harbor_project.main.project_id
  action      = "push"
}
# #resource "harbor_tasks" "main" {
# #  vulnerability_scan_policy = "daily"
# #}
#



# v2 problems !!!
### resource "harbor_tasks" "main" {
###   vulnerability_scan_policy = "daily"
### }

resource "harbor_registry" "main" {
  name        = "dockerhub"
  url         = "https://hub.docker.com"
  type        = "docker-hub"
  description = "Docker Hub Registry"
  insecure    = false
}

resource "harbor_registry" "helmhub" {
  name        = "helmhub"
  url         = "https://hub.helm.sh"
  type        = "helm-hub"
  description = "Helm Hub Registry"
  insecure    = false
}
##

resource "harbor_label" "main" {
  name = "testlabel"
}
