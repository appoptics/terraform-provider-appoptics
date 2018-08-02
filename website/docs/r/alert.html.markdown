---
layout: "appoptics"
page_title: "AppOptics: appoptics_alert"
sidebar_current: "docs-appoptics-resource-alert"
description: |-
  Provides a AppOptics Alert resource. This can be used to create and manage alerts on AppOptics.
---

# appoptics\_alert

Provides a AppOptics Alert resource. This can be used to
create and manage alerts on AppOptics.

## Example Usage

```hcl
# Create a new AppOptics alert
resource "appoptics_alert" "myalert" {
  name        = "MyAlert"
  description = "A Test Alert"
  services    = ["${appoptics_service.myservice.id}"]

  condition {
    type        = "above"
    threshold   = 10
    metric_name = "appoptics.cpu.percent.idle"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the alert.
* `description` - (Required) Description of the alert.
* `active` - whether the alert is active (can be triggered). Defaults to true.
* `rearm_seconds` - minimum amount of time between sending alert notifications, in seconds.
* `services` - list of notification service IDs.
* `condition` - A trigger condition for the alert. Conditions documented below.
* `attributes` - A hash of additional attribtues for the alert. Attributes documented below.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the alert.
* `name` - The name of the alert.
* `description` - (Required) Description of the alert.
* `active` - whether the alert is active (can be triggered). Defaults to true.
* `rearm_seconds` - minimum amount of time between sending alert notifications, in seconds.
* `services` - list of notification service IDs.
* `condition` - A trigger condition for the alert. Conditions documented below.

Conditions (`condition`) support the following:

* `type` - The type of condition. Must be one of `above`, `below` or `absent`.
* `metric_name`- The name of the metric this alert condition applies to.
* `source`- A source expression which identifies which sources for the given metric to monitor.
* `detect_reset` - boolean: toggles the method used to calculate the delta from the previous sample when the summary_function is `derivative`.
* `duration` - number of seconds condition must be true to fire the alert (required for type `absent`).
* `threshold` - float: measurements over this number will fire the alert (only for `above` or `below`).
* `summary_function` - Indicates which statistic of an aggregated measurement to alert on. ((only for `above` or `below`).

Attributes (`attributes`) support the following:

* `runbook_url` - a URL for the runbook to be followed when this alert is firing. Used in the AppOptics UI if set.
