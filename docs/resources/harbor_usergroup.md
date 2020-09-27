# Resource: harbor_usergroup

## Example Usage

```hcl
resource "harbor_usergroup" "developers" {
    name = "developers"
    type = "http"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user group to be created.

* `ldap_dn` - (Optional) The LDAP Group distingish name if `type` is `ldap`, defaults to `""`.

* `type` - (Required) The group type: `ldap`, `http` or `oidc`.

## Attributes Reference

In addition to all arguments, the following attribute is exported:

* `id` - The id of the group.

## Import

Harbor Usergroups can be imported using the `harbor_usergroup`, e.g.

```sh
terraform import harbor_usergroup.developers 5
```
