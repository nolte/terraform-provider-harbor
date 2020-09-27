---
subcategory: "project"
page_title: "Harbor: harbor_tasks"
description: |-
  Manages an Harbor Task
---

# Resource: harbor_tasks

## Example Usage
```
#resource "harbor_tasks" "main" {
#  vulnerability_scan_policy = "daily"
#}
```

## Argument Reference
The following arguments are supported:

* **vulnerability_scan_policy** - (Optional) The frequency of the vulnerability scanning is done. Can be to **"hourly"**, **"daily"** or **"weekly"**
