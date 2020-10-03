# Resource: harbor_webhook_policy

## Example Usage

```hcl
resource "harbor_project" "main" {
  name                    = "main"
  public                  = false # (Optional) Default value is false
  vulnerability_scanning  = true  # (Optional) Default value is true. Automatically scan images on push
  reuse_sys_cve_whitelist = false # (Optional) Default value is true.
  cve_whitelist           = ["CVE-2020-12345", "CVE-2020-54321"]
}

resource "harbor_webhook_policy" "main" {
  name         = "test-policy"
  description  = "Testing"
  project_id   = harbor_project.main.id
  endpoint_url = "https://www.googgle.com"
  event_types  = ["SCANNING_COMPLETED"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the webhook policy.

* `description` - (Optional) The description of the webhook policy.

* `project_id` - (Required) The project id of the project that will contain the webhook policy.

* `endpoint_url` - (Required) The webhook target address.

* `notify_type` - (Optional) `http` or `slack` (Default: `http`).

* `auth_header` - (Optional) The webhook auth header.

* `skip_cert_verify` - (Optional) Whether or not to skip cert verify (Default: `true`).

* `event_types` - (Required) List of events that will trigger the webhook. A full list of event types can be found [here](https://goharbor.io/docs/1.10/working-with-projects/project-configuration/configure-webhooks/).

* `enabled` - (Optional) Whether or not to enable the webhook (Default: `true`).

## Attributes Reference

In addition to all arguments, the following attribute is exported:

* `id` - The id of the webhook policy within the project: format is `${project_id}/${webhook_policy_id}`.

## Import

Harbor Project member can be imported using the `harbor_webhook_policy`, e.g.

```sh
terraform import harbor_webhook_policy.main 12/5
```
