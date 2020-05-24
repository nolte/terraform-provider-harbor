# Resource: harbor_registry

Harbor Doc: [managing-registries](https://goharbor.io/docs/2.0.0/administration/configuring-replication/create-replication-endpoints/#managing-registries)
Harbor Api: [Create](https://demo.goharbor.io/#/Products/post_registries)

## Example Usage

```hcl
resource "harbor_registry" "helmhub" {
  name        = "helmhub"
  url         = "https://hub.helm.sh"
  type        = "helm-hub"
  description = "Helm Hub Registry"
  insecure    = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The of the project that will be created in harbor.

* `url` - (Required) The registry remote endpoint, like `https://hub.docker.com`.

* `type` - (Required) registry Type possible values are `huawei-SWR, aws-ecr, ali-acr, jfrog-artifactory, gitlab, docker-registry, docker-hub, azure-acr, quay-io, helm-hub, harbor, google-gcr`.

* `description` - (Optional) The description of the registry will be displayed in harbor.

* `insecure` - (Optional) Harbor ignores insecure external registry errors. Can be set to `true` or `false` (Default: `false`)

## Attributes Reference

In addition to all argument, the folloing attributes are exported:

* `id` - The id of the registry with harbor.

## Import

Harbor Projects can be imported using the `harbor_registry`, e.g.

```sh
terraform import harbor_registry.helmhub 1
```
