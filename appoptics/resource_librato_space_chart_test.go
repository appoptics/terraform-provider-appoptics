package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppOpticsSpaceChartBasic(t *testing.T) {
	var spaceChart appoptics.Chart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsSpaceChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsSpaceChart_Full(t *testing.T) {
	var spaceChart appoptics.Chart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsSpaceChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfigFull,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsSpaceChart_Updated(t *testing.T) {
	var spaceChart appoptics.Chart

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsSpaceChartDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Foo Bar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckAppOpticsSpaceChartConfigNewValue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceChartExists("appoptics_space_chart.foobar", &spaceChart),
					resource.TestCheckResourceAttr(
						"appoptics_space_chart.foobar", "name", "Bar Baz"),
				),
			},
		},
	})
}

func testAccCheckAppOpticsSpaceChartDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*appoptics.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_space_chart" {
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

func testAccCheckAppOpticsSpaceChartExists(n string, spaceChart *appoptics.Chart) resource.TestCheckFunc {
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

		foundSpaceChart, err := client.ChartsService().Retrieve(id, spaceID)
		if err != nil {
			return err
		}

		if foundSpaceChart.ID != id {
			return fmt.Errorf("Space not found")
		}

		*spaceChart = *foundSpaceChart

		return nil
	}
}

const testAccCheckAppOpticsSpaceChartConfigBasic = `
resource "appoptics_space" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_space_chart" "foobar" {
    space_id = "${appoptics_space.foobar.id}"
    name = "Foo Bar"
    type = "line"
}`

const testAccCheckAppOpticsSpaceChartConfigNewValue = `
resource "appoptics_space" "foobar" {
    name = "Foo Bar"
}

resource "appoptics_space_chart" "foobar" {
    space_id = "${appoptics_space.foobar.id}"
    name = "Bar Baz"
    type = "line"
}`

const testAccCheckAppOpticsSpaceChartConfigFull = `
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
		metric = "system.cpu.utilization"
		// TODO
        source = "*"
    }

    # Minimal composite stream
    stream {
        composite = "s(\"cpu\", \"*\")"
    }

    # Full metric stream
    stream {
		metric = "system.cpu.utilization"
		// TODO
        // source = "*"
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
