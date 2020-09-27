# Resource: harbor_label

Harbor Doc: [Managing Labels](https://goharbor.io/docs/1.10/working-with-projects/working-with-images/create-labels/)
Harbor Api: [/labels](https://demo.goharbor.io/#/Products/post_labels)

## Example Usage

```hcl
resource "harbor_label" "main" {
  name        = "testlabel-acc-classic"
  description = "Test Label"
  color       = "#61717D"
  scope       = "g"
}

data "harbor_project" "project_1" {
  name = "main"
}

resource "harbor_label" "project_label" {
  name        = "projectlabel-acc-classic"
  description = "Test Label for Project"
  color       = "#333333"
  scope       = "p"
  project_id  = data.harbor_project.project_1.id
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) Name of the Project.

The following arguments are optional:

* `description` - (Optional)  The description of the label account will be displayed in harbor.

* `color` - (Optional) The colour the label.

* `scope` - (Optional) The scope the label, `p` for project and `g` for global.

* `project_id` - (Optional) The ID of project that the label belongs to, must be set if scope project.

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The id of the registry with harbor.

## Import

Harbor Projects can be imported using the `harbor_label`, e.g.

```sh
terraform import harbor_label.helmhub 1
```
