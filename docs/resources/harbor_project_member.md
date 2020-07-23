# Resource: harbor_project_member

## Example Usage

```hcl
--8<--

examples/tf-acception-test/project.tf

examples/tf-acception-test/usergroup.tf

examples/tf-acception-test/project_member.tf
--8<--


```

## Argument Reference

The following arguments are supported:

* `role` - (Required) The role permission to be given to the group: `project_admin`, `master`, `developer`, `guest` or `limited_guest`.

* `group_type` - (Required) The group type: `ldap`, `http` or `oidc`.

* `project_id` - (Required) The project id of the project to be given permission to the group members.

* `group_name` - (Required) The name of the user group to be given permissions for the project.

## Attributes Reference

In addition to all arguments, the following attribute is exported:

* `id` - The id of the project account: format is `${project_id}/${group_id}`.

## Import

Harbor Project member can be imported using the `harbor_project_member`, e.g.

```sh
terraform import harbor_project_member.developers_main 12/5
```

## Known limitations

The provider currently only handles group membership, not user.
