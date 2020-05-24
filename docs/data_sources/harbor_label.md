# Data Source: harbor_label

## Example Usage

```hcl
data "harbor_label" "label_1" {
  name   = "main"
  scope  = "g"
}

data "harbor_label" "label_2" {
  id = 4
}

```

## Argument Reference

- `id` - (Optional, int) ID of the label.
- `name` - (Optional, string) Name of the label.
- `scope` - (Optional, string) Scope of the label.

## Attributes Reference

- `id` - (int) Unique ID of the label.
- `name` - (string) Name of the label.
- `scope` - (string) Scope of the label.
