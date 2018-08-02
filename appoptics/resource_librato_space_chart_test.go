package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/akahn/go-librato/librato"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppOpticsSpaceChart_Basic(t *testing.T) {
	var spaceChart librato.SpaceChart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsSpaceChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					testAccCheckAppOpticsSpaceChartName(&spaceChart, "Foo Bar"),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsSpaceChart_Full(t *testing.T) {
	var spaceChart librato.SpaceChart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsSpaceChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfig_full,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					testAccCheckAppOpticsSpaceChartName(&spaceChart, "Foo Bar"),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsSpaceChart_Updated(t *testing.T) {
	var spaceChart librato.SpaceChart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsSpaceChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					testAccCheckAppOpticsSpaceChartName(&spaceChart, "Foo Bar"),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Foo Bar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfig_new_value,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					testAccCheckAppOpticsSpaceChartName(&spaceChart, "Bar Baz"),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Bar Baz"),
				),
			},
		},
	})
}

func testAccCheckAppOpticsSpaceChartDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*librato.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_space_chart" {
			continue
		}

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		spaceID, err := strconv.ParseUint(rs.Primary.Attributes["space_id"], 10, 0)
		if err != nil {
			return fmt.Errorf("Space ID not a number")
		}

		_, _, err = client.Spaces.GetChart(uint(spaceID), uint(id))

		if err == nil {
			return fmt.Errorf("Space Chart still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsSpaceChartName(spaceChart *librato.SpaceChart, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if spaceChart.Name == nil || *spaceChart.Name != name {
			return fmt.Errorf("Bad name: %s", *spaceChart.Name)
		}

		return nil
	}
}

func testAccCheckAppOpticsSpaceChartExists(n string, spaceChart *librato.SpaceChart) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Space Chart ID is set")
		}

		client := testAccProvider.Meta().(*librato.Client)

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		spaceID, err := strconv.ParseUint(rs.Primary.Attributes["space_id"], 10, 0)
		if err != nil {
			return fmt.Errorf("Space ID not a number")
		}

		foundSpaceChart, _, err := client.Spaces.GetChart(uint(spaceID), uint(id))

		if err != nil {
			return err
		}

		if foundSpaceChart.ID == nil || *foundSpaceChart.ID != uint(id) {
			return fmt.Errorf("Space not found")
		}

		*spaceChart = *foundSpaceChart

		return nil
	}
}

const testAccCheckAppOpticsSpaceChartConfig_basic = `
resource "appoptics_space" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_space_chart" "foobar" {
    space_id = "${appoptics_space.foobar.id}"
    name = "Foo Bar"
    type = "line"
}`

const testAccCheckAppOpticsSpaceChartConfig_new_value = `
resource "appoptics_space" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_space_chart" "foobar" {
    space_id = "${appoptics_space.foobar.id}"
    name = "Bar Baz"
    type = "line"
}`

const testAccCheckAppOpticsSpaceChartConfig_full = `
resource "appoptics_space" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_space" "barbaz" {
    name = "Bar Baz"
}

resource "appoptics_space_chart" "foobar" {
    space_id = "${appoptics_space.foobar.id}"
    name = "Foo Bar"
    type = "line"
    min = 0
    max = 100
    label = "Percent"
    related_space = "${appoptics_space.barbaz.id}"

    # Minimal metric stream
    stream {
        metric = "librato.cpu.percent.idle"
        source = "*"
    }

    # Minimal composite stream
    stream {
        composite = "s(\"cpu\", \"*\")"
    }

    # Full metric stream
    stream {
        metric = "librato.cpu.percent.idle"
        source = "*"
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
    }
}`
