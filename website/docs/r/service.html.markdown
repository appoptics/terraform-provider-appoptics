---
layout: "appoptics"
page_title: "AppOptics: appoptics_service"
sidebar_current: "docs-appoptics-resource-service"
description: |-
  Provides a AppOptics service resource. This can be used to create and manage notification services on AppOptics.
---

# appoptics\_service

Provides a AppOptics Service resource. This can be used to
create and manage notification services on AppOptics.

## Example Usage

```hcl
# Create a new AppOptics service
resource "appoptics_service" "email" {
  title = "Email the admins"
  type  = "mail"

  settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}
```

## Argument Reference

The following arguments are supported. Please check the [relevant documentation](https://github.com/appoptics/appoptics-services/tree/master/services) for each type of alert.

* `type` - (Required) The type of notificaion.
* `title` - (Required) The alert title.
* `settings` - (Required) a JSON hash of settings specific to the alert type.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alert.
* `type` - The type of notificaion.
* `title` - The alert title.
* `settings` - a JSON hash of settings specific to the alert type.
