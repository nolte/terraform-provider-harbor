# Data Source: harbor_registry

## Example Usage

```hcl
data "harbor_registry" "registry_1" {
  name = "main"
}

data "harbor_registry" "registry_2" {
  id = 4
}

```

## Argument Reference

- `id` - (Optional, string) ID of the datacenter.
- `name` - (Optional, string) Name of the datacenter.

## Attributes Reference

- `id` - (int) Unique ID of the datacenter.
- `name` - (string) Name of the datacenter.
