---
subcategory: "project"
page_title: "Harbor: harbor_project"
description: |-
  Manages an Harbor Project
---

# Resource: harbor_project

Handle a [Harbor Project Resource](https://goharbor.io/docs/1.10/working-with-projects/create-projects/).

## Example Usage

```hcl
resource "harbor_project" "main" {
  name                   = "main"
  public                 = false # (Optional) Default value is false
  vulnerability_scanning = true  # (Optional) Default vale is true. Automatically scan images on push
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) Name of the Project.

The following arguments are optional:

* `public` - (Optional) Handle the access to the hosted images. Default: `true`

    If `true` Any user can pull images from this project. This is a convenient way for you to share repositories with others.

    If `false` Only users who are members of the project can pull images

* `vulnerability_scanning` - (Optional) Activate [Vulnerability Scanning](https://goharbor.io/docs/1.10/administration/vulnerability-scanning/). Default: `true`


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Harbor Project ID.

## Import

Harbor Projects can be imported using the `harbor_project`, e.g.

```
$ terraform import harbor_project.main 1
```
