# Resource: harbor_registry

Harbor Doc: [managing-registries](https://goharbor.io/docs/2.0.0/administration/configuring-replication/create-replication-endpoints/#managing-registries)
Harbor Api: [Create](https://demo.goharbor.io/#/Products/post_registries)

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) of the project that will be created in harbor.

* `url` - (Required) The registry remote endpoint, like `https://hub.docker.com`.

* `type` - (Required) registry Type possible values are `huawei-SWR, aws-ecr, ali-acr, jfrog-artifactory, gitlab, docker-registry, docker-hub, azure-acr, quay-io, helm-hub, harbor, google-gcr`.

* `description` - (Optional) The description of the registry will be displayed in harbor.

* `insecure` - (Optional) Harbor ignores insecure external registry errors. Can be set to `true` or `false` (Default: `false`)

* `access_key` - (Optional) The registry access key.

* `access_secret` - (Optional) The registry access secret.

* `credential_type` - (Optional) Credential type, such as 'basic', 'oauth'.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The id of the registry with harbor.

## Import

Harbor Projects can be imported using the `harbor_registry`, e.g.

```sh
terraform import harbor_registry.helmhub 1
```
