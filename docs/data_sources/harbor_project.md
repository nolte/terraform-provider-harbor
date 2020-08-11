# Data Source: harbor_project

## Example Usage

```hcl
--8<--
examples/tf-acception-test-part-2/data_project.tf
--8<--

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
