package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAppOpticsDashboardChartBasic(t *testing.T) {
	var dashboardChart appoptics.Chart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsDashboardChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsDashboardChartConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsDashboardChartExists("appoptics_dashboard_chart.foobar", &dashboardChart),
					resource.TestCheckResourceAttr(
						"appoptics_dashboard_chart.foobar", "name", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsDashboardChart_Full(t *testing.T) {
	var dashboardChart appoptics.Chart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsDashboardChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsDashboardChartConfigFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsDashboardChartExists("appoptics_dashboard_chart.foobar", &dashboardChart),
					resource.TestCheckResourceAttr(
						"appoptics_dashboard_chart.foobar", "name", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsDashboardChart_Updated(t *testing.T) {
	var dashboardChart appoptics.Chart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsDashboardChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsDashboardChartConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsDashboardChartExists("appoptics_dashboard_chart.foobar", &dashboardChart),
					resource.TestCheckResourceAttr(
						"appoptics_dashboard_chart.foobar", "name", "Foo Bar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckAppOpticsDashboardChartConfigNewValue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsDashboardChartExists("appoptics_dashboard_chart.foobar", &dashboardChart),
					resource.TestCheckResourceAttr(
						"appoptics_dashboard_chart.foobar", "name", "Bar Baz"),
				),
			},
		},
	})
}

func testAccCheckAppOpticsDashboardChartDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*appoptics.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_dashboard_chart" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		spaceID, err := strconv.Atoi(rs.Primary.Attributes["space_id"])
		if err != nil {
			return fmt.Errorf("Space ID not a number")
		}

		_, err = client.ChartsService().Retrieve(id, spaceID)
		if err == nil {
			return fmt.Errorf("Space Chart still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsDashboardChartExists(n string, dashboardChart *appoptics.Chart) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Space Chart ID is set")
		}

		client := testAccProvider.Meta().(*appoptics.Client)

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		spaceID, err := strconv.Atoi(rs.Primary.Attributes["space_id"])
		if err != nil {
			return fmt.Errorf("Space ID not a number")
		}

		foundDashboardChart, err := client.ChartsService().Retrieve(id, spaceID)
		if err != nil {
			return err
		}

		if foundDashboardChart.ID != id {
			return fmt.Errorf("Space not found")
		}

		*dashboardChart = *foundDashboardChart

		return nil
	}
}

const testAccCheckAppOpticsDashboardChartConfigBasic = `
resource "appoptics_dashboard" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_dashboard_chart" "foobar" {
    space_id = "${appoptics_dashboard.foobar.id}"
    name = "Foo Bar"
	type = "line"
}`

const testAccCheckAppOpticsDashboardChartConfigNewValue = `
resource "appoptics_dashboard" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_dashboard_chart" "foobar" {
    space_id = "${appoptics_dashboard.foobar.id}"
    name = "Bar Baz"
	type = "line"
	min = 0
	max = 100
}`

const testAccCheckAppOpticsDashboardChartConfigFull = `
resource "appoptics_dashboard" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_dashboard" "barbaz" {
    name = "Bar Baz"
}

resource "appoptics_dashboard_chart" "foobar" {
    space_id = "${appoptics_dashboard.foobar.id}"
    name = "Foo Bar"
    type = "line"
    min = 0
    max = 100
    label = "Percent"
    related_space = "${appoptics_dashboard.barbaz.id}"

    # Minimal metric stream
    stream {
		metric = "system.cpu.utilization"
		tags {
			name = "hostname"
			grouped = true
			values = ["host1", "host2"]
		}
    }

    # Minimal composite stream
    stream {
        composite = "s(\"cpu\", \"*\")"
    }

    # Full metric stream
    stream {
		metric = "system.cpu.utilization"
        group_function = "average"
        summary_function = "max"
        name = "CPU usage"
        color = "#990000"
        units_short = "%"
        units_long = "percent"
        min = 0
        max = 100
        transform_function = "x * 100"
		period = 60
		tags {
			name = "hostname"
			grouped = true
			values = ["host1", "host2"]
		}
    }
}`
