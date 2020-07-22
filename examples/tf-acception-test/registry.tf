resource "harbor_registry" "dockerhub" {
  name        = "dockerhub-acc-classic"
  url         = "https://hub.docker.com"
  type        = "docker-hub"
  description = "Docker Hub Registry"
  insecure    = false
}

resource "harbor_registry" "helmhub" {
  name        = "helmhub-acc-classic"
  url         = "https://hub.helm.sh"
  type        = "helm-hub"
  description = "Helm Hub Registry"
  insecure    = false
}
