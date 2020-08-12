# This example is for the terraform 0.12+ compatible versions of the
# provider

// Used to identify all things made in integration test runs
variable "tf-name-fragment" {
  type    = string
  default = "tf_provider_test"
}

variable "test-metric-name" {
  type    = string
  default = "tf_provider_test.cpu.percent.used" // can't interpolate into variables
}


provider "appoptics" {
  //
  //  Implicitly consumes APPOPTICS_TOKEN from the environment
  //
}

//
// Dashboard
//
resource "appoptics_dashboard" "test_dashboard" {
  name = "${var.tf-name-fragment} Dashboard"
}


//
// Notification Service
//
resource "appoptics_notification_service" "test_service" {
  title    = "${var.tf-name-fragment} Email Notification Service"
  type     = "mail"
  settings = <<EOF
{
  "addresses": "foobar@example.com"
}
EOF
}

//
// Metric
//
resource "appoptics_metric" "test_metric" {
  name         = var.test-metric-name
  display_name = "Terraform Test CPU Utilization"
  period       = 60
  type         = "gauge"
  attributes {
    color               = "#A3BE8C"
    summarize_function  = "average"
    display_max         = 100.0
    display_units_long  = "CPU Utilization Percent"
    display_units_short = "cpu %"
  }
}

//
// Chart
//
resource "appoptics_dashboard_chart" "test_chart" {
  space_id   = appoptics_dashboard.test_dashboard.id
  name       = "Test Chart"
  depends_on = [appoptics_metric.test_metric]
  min        = 0
  max        = 100
  label      = "Used"
  type       = "line"

  stream {
    metric      = appoptics_metric.test_metric.name
    color       = "#fa7268"
    units_short = "%"
    units_long  = "Percentage used"

    tags {
      name   = "environment"
      values = ["staging"]
    }
  }
}

//
// Alert
//
resource "appoptics_alert" "test_alert" {
  name          = "${var.tf-name-fragment}.Alert"
  description   = "Managed by Terraform"
  rearm_seconds = 10800

  depends_on = [appoptics_metric.test_metric]

  condition {
    type             = "above"
    threshold        = 0
    metric_name      = var.test-metric-name
    duration         = 60
    summary_function = "sum"

    tag {
      name    = "logical_thing"
      grouped = true
      values  = ["staging"]
    }

    tag {
      name    = "event_status"
      grouped = true
      values  = ["count of scary things"]
    }

  }
}
