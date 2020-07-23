# Data Source: harbor_registry

## Example Usage

```hcl
--8<--
examples/tf-acception-test-part-2/data_registry.tf
--8<--

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
