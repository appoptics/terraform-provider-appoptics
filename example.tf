// Used to identify all things made in integration test runs
variable "tf-name-fragment" {
  type = "string"
  default = "tf_provider_test"
}

variable "test-metric-name" {
  type = "string"
  default = "tf_provider_test.cpu.percent.used" // can't interpolate into variables
}

provider "appoptics" {
  //
  //  Implicitly consumes APPOPTICS_TOKEN from the environment
  //
}

//
// Space
//
resource "appoptics_space" "test_space" {
  name = "${var.tf-name-fragment} Space"
}

//
// Notification Service
//
resource "appoptics_service" "test_service"{
  title = "${var.tf-name-fragment} Email Notification Service"
  type = "mail"
  settings = <<EOF
{
  "addresses": "foobar@example.com"
}
EOF
}

//
// Metric
//
resource "appoptics_metric" "test_metric"{
  name = "${var.test-metric-name}"
  display_name = "Terraform Test CPU Utilization"
  period = 60
  attributes = {
    color = "#A3BE8C"
    summarize_function = "sum"
    display_max = 100.0
    display_units_long = "CPU Utilization Percent"
    display_units_short = "cpu %"
  }
}

//
// Chart
//

//
// Alert
//
resource "appoptics_alert" "test_alert" {
  name        = "${var.tf-name-fragment}.Alert"
  description = "Managed by Terraform"
  rearm_seconds = 10800

  depends_on = ["appoptics_metric.test_metric"]

  condition {
    type        = "above"
    threshold   = 0
    metric_name = "${var.test-metric-name}"
    duration    = 60
    summary_function = "sum"

    tag {
      name    = "environment"
      grouped = true
      values  = ["staging"]
    }

    tag {
      name    = "event_status"
      grouped = true
      values  = ["test metric has test metricness"]
    }

  }
}