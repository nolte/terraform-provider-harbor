
provider "harbor" {
  url      = "harbor.192-168-178-51.sslip.io"
  insecure = true
  #url      = "demo.goharbor.io"
  #basepath = "/api/v2.0"
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
resource "harbor_config_email" "conf_email" {
  email_host     = "main2"
  email_port     = 25
  email_username = "main2"
  email_password = "main2"
  email_from     = "main2"
  email_ssl      = false
}
# 
# 
resource "harbor_config_auth" "oidc" {
  auth_mode          = "oidc_auth"
  oidc_name          = "azure"
  oidc_endpoint      = "https://login.microsoftonline.com/v2.0"
  oidc_client_id     = "OIDC Client ID goes here"
  oidc_client_secret = "ODDC Client Secret goes here"
  oidc_scope         = "openid,email"
  oidc_verify_cert   = true
}
# 
resource "harbor_config_system" "main" {
  project_creation_restriction = "everyone"
  robot_token_expiration       = 5259492
}


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
