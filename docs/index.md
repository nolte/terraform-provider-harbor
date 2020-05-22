---
layout: "harbor"
page_title: "Provider: Harbor"
sidebar_current: "docs-harbor-index"
description: |-
  The Harbor Registry provider is used to interact with the many resources supported by Harbor. The provider needs to be configured with the proper credentials before it can be used.
---

# Harbor Provider

Summary of what the provider is for, including use cases and links to
app/service documentation.


## Example Usage

```hcl
provider "harbor" {
  url      = "demo.goharbor.io"
  basepath = "/api/v2.0"
  username = "admin"
  password = "Harbor12345"
}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Harbor
 `provider` block:

* `url` - (Required) ddd

* `basepath` - (Optional)

* `basepath` - (Optional)

* `basepath` - (Optional)

* `basepath` - (Optional)

* `basepath` - (Optional)

* `basepath` - (Optional)

* `basepath` - (Optional)


## Install the Custom Provider

```bash

# for example https://github.com/nolte/terraform-provider-harbor/releases/download/release/v0.1.0/terraform-provider-harbor_v0.1.0_linux_amd64.tar.gz
LATEST_LINUX_RELEASE=$(curl -sL https://api.github.com/repos/nolte/terraform-provider-harbor/releases/latest | jq -r '.assets[].browser_download_url' | grep '_linux_amd64')

# direct install to your personal plugin directory
wget -qO- $LATEST_LINUX_RELEASE | tar -xvz -C ~/.terraform.d/plugins/linux_amd64/
```

Here a link to the Terraform Doc how to install [third-party-plugins](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
