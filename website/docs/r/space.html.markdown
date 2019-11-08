---
layout: "appoptics"
page_title: "AppOptics: appoptics_space"
sidebar_current: "docs-appoptics-resource-space"
description: |-
  Provides a AppOptics Space resource. This can be used to create and manage spaces on AppOptics.
---

# appoptics\_space

Provides a AppOptics Space resource. This can be used to
create and manage spaces on AppOptics.

## Example Usage

```hcl
# Create a new AppOptics space
resource "appoptics_space" "default" {
  name = "My New Space"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the space.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the space.
* `name` - The name of the space.
