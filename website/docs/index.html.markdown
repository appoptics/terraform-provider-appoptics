---
layout: "appoptics"
page_title: "Provider: AppOptics"
sidebar_current: "docs-appoptics-index"
description: |-
  The AppOptics provider is used to interact with the resources supported by AppOptics. The provider needs to be configured with the proper credentials before it can be used.
---

# AppOptics Provider

The AppOptics provider is used to interact with the
resources supported by AppOptics. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the AppOptics provider
provider "appoptics" {
  token = "${var.librato_token}"
}

# Create a new space
resource "appoptics_space" "default" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `token` - (Required) AppOptics API token. It must be provided, but it can also
  be sourced from the `APPOPTICS_TOKEN` environment variable.
