package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppOpticsAlertMinimal(t *testing.T) {
	var alert appoptics.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfigMinimal(name),
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

func TestAccAppOpticsAlertBasic(t *testing.T) {
	var alert appoptics.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfigBasic(name),
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

func TestAccAppOpticsAlertFull(t *testing.T) {
	var alert appoptics.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfigFull(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertName(&alert, name),
					testAccCheckAppOpticsAlertDescription(&alert, "A Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.836525194.metric_name", "appoptics.cpu.percent.idle"),
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

func TestAccAppOpticsAlertUpdated(t *testing.T) {
	var alert appoptics.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertDescription(&alert, "A Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
				),
			},
			{
				Config: testAccCheckAppOpticsAlertConfigNewValue(name),
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

func TestAccAppOpticsAlertRename(t *testing.T) {
	var alert appoptics.Alert
	name := acctest.RandString(10)
	newName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
				),
			},
			{
				Config: testAccCheckAppOpticsAlertConfigBasic(newName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", newName),
				),
			},
		},
	})
}

func TestAccAppOpticsAlertFullUpdate(t *testing.T) {
	var alert appoptics.Alert
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsAlertConfigFullUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsAlertExists("appoptics_alert.foobar", &alert),
					testAccCheckAppOpticsAlertName(&alert, name),
					testAccCheckAppOpticsAlertDescription(&alert, "A Test Alert"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "rearm_seconds", "1200"),
					resource.TestCheckResourceAttr(
						"appoptics_alert.foobar", "condition.2524844643.metric_name", "appoptics.cpu.percent.idle"),
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
	client := testAccProvider.Meta().(*appoptics.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_alert" {
			continue
		}

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		_, err = client.AlertsService().Retrieve(int(id))

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsAlertName(alert *appoptics.Alert, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Name != name {
			return fmt.Errorf("Bad name: %s", alert.Name)
		}

		return nil
	}
}

func testAccCheckAppOpticsAlertDescription(alert *appoptics.Alert, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Description != description {
			return fmt.Errorf("Bad description: %s", alert.Description)
		}

		return nil
	}
}

func testAccCheckAppOpticsAlertExists(n string, alert *appoptics.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Alert ID is set")
		}

		client := testAccProvider.Meta().(*appoptics.Client)

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		foundAlert, err := client.AlertsService().Retrieve(int(id))

		if err != nil {
			return err
		}

		if foundAlert.ID == 0 {
			return fmt.Errorf("Alert not found")
		}

		*alert = *foundAlert

		return nil
	}
}

func testAccCheckAppOpticsAlertConfigMinimal(name string) string {
	return fmt.Sprintf(`
resource "appoptics_alert" "foobar" {
    name = "%s"
}`, name)
}

func testAccCheckAppOpticsAlertConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "appoptics_alert" "foobar" {
    name = "%s"
    description = "A Test Alert"
}`, name)
}

func testAccCheckAppOpticsAlertConfigNewValue(name string) string {
	return fmt.Sprintf(`
resource "appoptics_alert" "foobar" {
    name = "%s"
    description = "A modified Test Alert"
}`, name)
}

func testAccCheckAppOpticsAlertConfigFull(name string) string {
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
      metric_name = "appoptics.cpu.percent.idle"
    }
    attributes {
      runbook_url = "https://www.youtube.com/watch?v=oHg5SJYRHA0"
    }
    active = false
    rearm_seconds = 300
}`, name)
}

func testAccCheckAppOpticsAlertConfigFullUpdate(name string) string {
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
      metric_name = "appoptics.cpu.percent.idle"
    }
    attributes {
      runbook_url = "https://www.youtube.com/watch?v=oHg5SJYRHA0"
    }
    active = false
    rearm_seconds = 1200
}`, name)
}
