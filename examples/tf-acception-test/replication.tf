resource "harbor_replication" "pull_helm_chart" {
  name                        = "helm-prometheus-operator-acc-classic"
  description                 = "Prometheus Operator Replica ACC Classic"
  source_registry_id          = harbor_registry.helmhub.id
  source_registry_filter_name = "stable/prometheus-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = harbor_project.main.name
}

resource "harbor_replication" "push_helm_chart" {
  name                        = "docker-push-acc-classic"
  description                 = "Push Docker Replica ACC Classic"
  destination_registry_id     = harbor_registry.dockerhub.id
  source_registry_filter_name = "${harbor_project.main.name}/vscode-devcontainers/k8s-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = "notexisting"
}
