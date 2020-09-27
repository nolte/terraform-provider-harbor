# Resource: harbor_robot_account

## Example Usage

```hcl
resource "harbor_project" "main" {
  name                    = "main"
  public                  = false # (Optional) Default value is false
  vulnerability_scanning  = true  # (Optional) Default value is true. Automatically scan images on push
  reuse_sys_cve_whitelist = false # (Optional) Default value is true.
  cve_whitelist           = ["CVE-2020-12345", "CVE-2020-54321"]
}

resource "harbor_robot_account" "master_robot" {
  name        = "god"
  description = "Robot account used to push images to harbor"
  project_id  = harbor_project.main.id
  actions     = ["docker_read", "docker_write", "helm_read", "helm_write"]
}

output "harbor_robot_account_token" {
  value = harbor_robot_account.master_robot.token
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The of the project that will be created in harbor.

* `description` - (Optional) The description of the robot account will be displayed in harbor.

* `project_id` - (Required) The project id of the project that the robot account will be associated with.

* `actions` - (Optional)

## Attributes Reference

In addition to all argument, the following attributes are exported:

* `id` - The id of the robot account.

* `token` - The token of the robot account.

## Import

Harbor Projects can be imported using the `harbor_robot_account`, e.g.

```sh
terraform import harbor_robot_account.master_robot 29
```
