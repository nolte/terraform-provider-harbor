# Data Source: harbor_project

## Example Usage

```hcl
data "harbor_project" "project_1" {
  name = "main"
}

data "harbor_project" "project_2" {
  id = 4
}

```

## Argument Reference

- `id` - (Optional, int) ID of the project.
- `name` - (Optional, string) Name of the project.

## Attributes Reference

- `id` - (int) Unique ID of the project.
- `name` - (string) Name of the project.
