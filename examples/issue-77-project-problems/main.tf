terraform {
  required_providers {
    #vault = {
    #  source  = "hashicorp/vault"
    #  version = "2.11.0"
    #}

    harbor = {
      source  = "nolte/harbor"
      version = "0.1.9"
    }
  }
  required_version = ">= 0.13"
}

resource "harbor_registry" "dockerhub" {
  name        = "dockerhub-acc-classic"
  url         = "https://hub.docker.com"
  type        = "docker-hub"
  description = "Docker Hub Registry"
  insecure    = false
}

resource "harbor_registry" "gcr_io" {
  name        = "k8s-gcr-io-classic"
  url         = "https://k8s.gcr.io"
  type        = "docker-hub"
  description = "Google K8S Container Registry"
  insecure    = false
}

#provider "harbor" {
#  host     = data.vault_generic_secret.harbor_credentials.data["host"]
#  username = data.vault_generic_secret.harbor_credentials.data["user"]
#  password = data.vault_generic_secret.harbor_credentials.data["password"]
#  schema   = "https"
#  insecure = false
#  basepath = "/api/v2.0"
#}
provider "harbor" {

}

resource "harbor_project" "rancher" {
  name                   = "rancher"
  public                 = true  # (Optional) Default value is false
  vulnerability_scanning = false # (Optional) Default value is true. Automatically scan images on push
}

resource "harbor_project" "metallb" {
  name                   = "metallb"
  public                 = true  # (Optional) Default value is false
  vulnerability_scanning = false # (Optional) Default value is true. Automatically scan images on push
}
resource "harbor_project" "ingress_nginx" {
  name                   = "ingress-nginx"
  public                 = true  # (Optional) Default value is false
  vulnerability_scanning = false # (Optional) Default value is true. Automatically scan images on push
}
