# Resource: harbor_replication_pull

Harbor Doc: [configuring-replication](https://goharbor.io/docs/2.0.0/administration/configuring-replication/)
Harbor Api: [Create](https://demo.goharbor.io/#/Products/post_registries)

## Example Usage

```hcl
resource "harbor_project" "project_replica" {
    name                   = "acc-project-replica-test"
    public                 = false
    vulnerability_scanning = false
}
resource "harbor_registry" "registry_replica_helm_hub" {
  name        = "acc-registry-replica-test"
  url         = "https://hub.helm.sh"
  type        = "helm-hub"
  description = "Helm Hub Registry"
  insecure    = false
}
resource "harbor_replication_pull" "pull_helm_chart" {
  name                        = "acc-helm-prometheus-operator-test"
  description                 = "Prometheus Operator Replica"
  source_registry_id          = harbor_registry.registry_replica_helm_hub.id
  source_registry_filter_name = "stable/prometheus-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = harbor_project.project_replica.name
}
```

## Argument Reference

The following arguments are supported:

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The id of the registry with harbor.

## Import

Harbor Projects can be imported using the `harbor_replication_pull`, e.g.

```sh
terraform import harbor_replication_pull.helmhub_prometheus 1
```
