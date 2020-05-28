
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

resource "harbor_robot_account" "master_robot" {
  name        = "god"
  description = "Robot account used to push images to harbor"
  project_id  = harbor_project.main.id
  actions     = ["docker_read", "docker_write", "helm_read", "helm_write"]
}




output "harbor_robot_account_token" {
  value = harbor_robot_account.master_robot.token
}

#
resource "harbor_registry" "dockerhub" {
  name        = "dockerhub-acc-classic"
  url         = "https://hub.docker.com"
  type        = "docker-hub"
  description = "Docker Hub Registry"
  insecure    = false
}
#
resource "harbor_registry" "helmhub" {
  name        = "helmhub-acc-classic"
  url         = "https://hub.helm.sh"
  type        = "helm-hub"
  description = "Helm Hub Registry"
  insecure    = false
}

#
resource "harbor_label" "main" {
  name        = "testlabel-acc-classic"
  description = "Test Label"
  color       = "#61717D"
  scope       = "g"
}

resource "harbor_label" "project_label" {
  name        = "projectlabel-acc-classic"
  description = "Test Label for Project"
  color       = "#333333"
  scope       = "p"
  project_id  = harbor_project.main.id
}
###


resource "harbor_replication_pull" "pull_helm_chart" {
  name                        = "helm-prometheus-operator-acc-classic"
  description                 = "Prometheus Operator Replica ACC Classic"
  source_registry_id          = harbor_registry.helmhub.id
  source_registry_filter_name = "stable/prometheus-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = harbor_project.main.name
}

resource "harbor_replication_pull" "push_helm_chart" {
  name                        = "docker-push-acc-classic"
  description                 = "Push Docker Replica ACC Classic"
  destination_registry_id     = harbor_registry.dockerhub.id
  source_registry_filter_name = "${harbor_project.main.name}/vscode-devcontainers/k8s-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = "notexisting"
}
