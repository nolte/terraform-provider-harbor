# Data Source: harbor_registry

## Example Usage

```hcl
data "harbor_registry" "registry_1" {
  name = "dockerhub-acc-classic"
}

data "harbor_registry" "registry_2" {
  id = 4
}

```

## Argument Reference

- `id` - (Optional, string) ID of the registry.
- `name` - (Optional, string) Name of the registry.

## Attributes Reference

- `id` - (int) Unique ID of the registry.
- `name` - (string) Name of the registry.
