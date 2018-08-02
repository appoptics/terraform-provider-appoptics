package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/appoptics/go-librato/librato"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppOpticsAlert_Minimal(t *testing.T) {
	var alert librato.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfig_minimal(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertName(&alert, name),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
				),
			},
		},
	})
}

func TestAccAppOpticsAlert_Basic(t *testing.T) {
	var alert librato.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertName(&alert, name),
					testAccCheckAppOpticsAlertDescription(&alert, "A Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
				),
			},
		},
	})
}

func TestAccAppOpticsAlert_Full(t *testing.T) {
	var alert librato.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertName(&alert, name),
					testAccCheckAppOpticsAlertDescription(&alert, "A Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.836525194.metric_name", "librato.cpu.percent.idle"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.836525194.type", "above"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.836525194.threshold", "10"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.836525194.duration", "600"),
				),
			},
		},
	})
}

func TestAccAppOpticsAlert_Updated(t *testing.T) {
	var alert librato.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertDescription(&alert, "A Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
				),
			},
			{
				Config: testAccCheckAppOpticsAlertConfig_new_value(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertDescription(&alert, "A modified Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "description", "A modified Test Alert"),
				),
			},
		},
	})
}

func TestAccAppOpticsAlert_Rename(t *testing.T) {
	var alert librato.Alert
	name := acctest.RandString(10)
	newName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
				),
			},
			{
				Config: testAccCheckAppOpticsAlertConfig_basic(newName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", newName),
				),
			},
		},
	})
}

func TestAccAppOpticsAlert_FullUpdate(t *testing.T) {
	var alert librato.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfig_full_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertName(&alert, name),
					testAccCheckAppOpticsAlertDescription(&alert, "A Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "rearm_seconds", "1200"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.2524844643.metric_name", "librato.cpu.percent.idle"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.2524844643.type", "above"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.2524844643.threshold", "10"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.2524844643.duration", "60"),
				),
			},
		},
	})
}

func testAccCheckAppOpticsAlertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*librato.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_alert" {
			continue
		}

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		_, _, err = client.Alerts.Get(uint(id))

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsAlertName(alert *librato.Alert, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Name == nil || *alert.Name != name {
			return fmt.Errorf("Bad name: %s", *alert.Name)
		}

		return nil
	}
}

func testAccCheckAppOpticsAlertDescription(alert *librato.Alert, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Description == nil || *alert.Description != description {
			return fmt.Errorf("Bad description: %s", *alert.Description)
		}

		return nil
	}
}

func testAccCheckAppOpticsAlertExists(n string, alert *librato.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Alert ID is set")
		}

		client := testAccProvider.Meta().(*librato.Client)

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		foundAlert, _, err := client.Alerts.Get(uint(id))

		if err != nil {
			return err
		}

		if foundAlert.ID == nil || *foundAlert.ID != uint(id) {
			return fmt.Errorf("Alert not found")
		}

		*alert = *foundAlert

		return nil
	}
}

func testAccCheckAppOpticsAlertConfig_minimal(name string) string {
	return fmt.Sprintf(`
resource "appoptics_alert" "foobar" {
    name = "%s"
}`, name)
}

func testAccCheckAppOpticsAlertConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "appoptics_alert" "foobar" {
    name = "%s"
    description = "A Test Alert"
}`, name)
}

func testAccCheckAppOpticsAlertConfig_new_value(name string) string {
	return fmt.Sprintf(`
resource "appoptics_alert" "foobar" {
    name = "%s"
    description = "A modified Test Alert"
}`, name)
}

func testAccCheckAppOpticsAlertConfig_full(name string) string {
	return fmt.Sprintf(`
resource "appoptics_service" "foobar" {
    title = "Foo Bar"
    type = "mail"
    settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}

resource "appoptics_alert" "foobar" {
    name = "%s"
    description = "A Test Alert"
    services = [ "${appoptics_service.foobar.id}" ]
    condition {
      type = "above"
      threshold = 10
      duration = 600
      metric_name = "librato.cpu.percent.idle"
    }
    attributes {
      runbook_url = "https://www.youtube.com/watch?v=oHg5SJYRHA0"
    }
    active = false
    rearm_seconds = 300
}`, name)
}

func testAccCheckAppOpticsAlertConfig_full_update(name string) string {
	return fmt.Sprintf(`
resource "appoptics_service" "foobar" {
    title = "Foo Bar"
    type = "mail"
    settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}

resource "appoptics_alert" "foobar" {
    name = "%s"
    description = "A Test Alert"
    services = [ "${appoptics_service.foobar.id}" ]
    condition {
      type = "above"
      threshold = 10
      duration = 60
      metric_name = "librato.cpu.percent.idle"
    }
    attributes {
      runbook_url = "https://www.youtube.com/watch?v=oHg5SJYRHA0"
    }
    active = false
    rearm_seconds = 1200
}`, name)
}
