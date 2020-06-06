# Resource: harbor_replication

Harbor Doc: [configuring-replication](https://goharbor.io/docs/2.0.0/administration/configuring-replication/)  


## Example Usage

```hcl
data "harbor_project" "project_replica" {
  name = "main"
}

data "harbor_registry" "registry_helm_hub" {
  name = "helmhub"
}

data "harbor_registry" "registry_docker_hub" {
  name = "dockerhub"
}

resource "harbor_replication" "pull_based_helm" {
  name                        = "acc-helm-prometheus-operator-test"
  description                 = "Prometheus Operator Replica"
  source_registry_id          = data.harbor_registry.registry_replica_helm_hub.id
  source_registry_filter_name = "stable/prometheus-operator"
  source_registry_filter_tag  = "**"
  destination_namespace       = data.harbor_project.project_replica.name
}

resource "harbor_replication_pull" "push_based_docker" {
  name                        = "docker-push"
  description                 = "Push Docker"
  destination_registry_id     = data.harbor_registry.registry_docker_hub.id
  destination_namespace       = "notexisting"
  source_registry_filter_name = "${data.harbor_project.project_replica.name}/vscode-devcontainers/k8s-operator"
  source_registry_filter_tag  = "**"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) of the replication that will be created in harbor.

* `description` - (Optional) of replication that will be displayed in harbor.

* `source_registry_id` - (Optional) Used for pull the resources from the remote registry to the local Harbor.

* `source_registry_filter_name` - (Optional) Filter the name of the resource. Leave empty or use '\*\*' to match all. 'library/\*\*' only matches resources under 'library'.

* `source_registry_filter_tag` - (Optional) Filter the tag/version part of the resources. Leave empty or use '\*\*' to match all. '1.0*' only matches the tags that starts with '1.0'.

* `destination_namespace` - (Optional) Destination namespace Specify the destination namespace. If empty, the resources will be put under the same namespace as the source.

* `destination_registry_id` - (Optional) The target Registry ID, used only for `push-based` replications.

* `trigger_mode` - (Optional) Can be `manual`,`scheduled` and for push-based addition `event_based`, Default: `manual`

* `trigger_cron` - (Optional) Used cron for `scheduled` trigger mode, like `* * 5 * * *`

* `override` - (Optional) Specify whether to override the resources at the destination if a resource with the same name exists. Default: `false`

* `enabled` - (Optional) 


## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The id of the registry policy with harbor.

## Import

Harbor Projects can be imported using the `harbor_replication`, e.g.

```sh
terraform import harbor_replication.helmhub_prometheus 1
```
