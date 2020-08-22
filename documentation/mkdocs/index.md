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
# example for harbor v2 api usage
provider "harborv2" {
  host     = "demo.goharbor.io"
  schema   = "https"
  insecure = true
  basepath = "/api/v2"
  username = "admin"
  password = "Harbor12345"
}

# example for harbor v1 api usage
provider "harborv1" {
  host     = var.harbor_endpoint
  schema   = "https"
  insecure = true
  basepath = var.harbor_base_path
  username = "admin"
  password = "Harbor12345"
}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Harbor
 `provider` block:

* `host` - (Required) Hostname from the [Harbor](https://goharbor.io) Service. like _demo.goharbor.io_, can also be specified with the `HARBOR_ENDPOINT` environment variable.

* `username` - (Required) Username for authorize at the harbor, can also be specified with the `HARBOR_USERNAME` environment variable.

* `password` - (Required) Password from given user, can also be specified with the `HARBOR_PASSWORD` environment variable.

* `schema` - (Optional) Set Used http Schema, possible values are: ```https,http```. Default: ```https```

* `insecure` - (Optional) Verify Https Certificates. Default: ```false```, can also be specified with the `HARBOR_INSECURE` environment variable.

* `basepath` - (Optional) The Harbor Api basepath, for example use ```/api``` for default HarborV1 and ```/api/v2.0``` for Harbor V2 Deployments. Default: ```/api```, can also be specified with the `HARBOR_BASEPATH` environment variable.


## Install the Custom Provider

Download the Latest Release, [![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nolte/terraform-provider-harbor)](https://github.com/nolte/terraform-provider-harbor/releases).

```bash
# for example https://github.com/nolte/terraform-provider-harbor/releases/download/release/v0.1.0/terraform-provider-harbor_v0.1.0_linux_amd64.tar.gz
LATEST_LINUX_RELEASE=$(curl -sL https://api.github.com/repos/nolte/terraform-provider-harbor/releases/latest | jq -r '.assets[].browser_download_url' | grep '_linux_amd64')
TAG_NAME=$(curl -sL https://api.github.com/repos/nolte/terraform-provider-harbor/releases/latest | jq -r '.tag_name')

# direct install to your personal plugin directory
wget -q $LATEST_LINUX_RELEASE -o /tmp/terraform-provider-harbor.zip \
    && unzip -j "/tmp/terraform-provider-harbor.zip" "terraform-provider-harbor_${TAG_NAME}" -d /home/${USERNAME}/.terraform.d/plugins/linux_amd64 \
    && rm /tmp/terraform-provider-harbor.zip
```

Here a link to the Terraform Doc how to install [third-party-plugins](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
