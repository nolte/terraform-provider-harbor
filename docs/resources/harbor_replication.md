# Resource: harbor_replication

Harbor Doc: [configuring-replication](https://goharbor.io/docs/2.0.0/administration/configuring-replication/)


## Example Usage

```hcl
resource "harbor_project" "main" {
  name                    = "main"
  public                  = false # (Optional) Default value is false
  vulnerability_scanning  = true  # (Optional) Default value is true. Automatically scan images on push
  reuse_sys_cve_whitelist = false # (Optional) Default value is true.
  cve_whitelist           = ["CVE-2020-12345", "CVE-2020-54321"]
}

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
