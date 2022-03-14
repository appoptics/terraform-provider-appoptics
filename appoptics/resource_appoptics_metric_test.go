package appoptics

import (
	"fmt"
	"strings"
	"testing"

	"github.com/appoptics/appoptics-api-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAppOpticsMetrics(t *testing.T) {
	var metric appoptics.Metric

	name := fmt.Sprintf("tftest-metric-%s", acctest.RandString(10))
	typ := "gauge"
	desc1 := fmt.Sprintf("A test %s metric", typ)
	desc2 := fmt.Sprintf("An updated test %s metric", typ)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: gaugeMetricConfig(name, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsMetricExists("appoptics_metric.foobar", &metric),
					testAccCheckAppOpticsMetricName(&metric, name),
					testAccCheckAppOpticsMetricType(&metric, typ),
					resource.TestCheckResourceAttr(
						"appoptics_metric.foobar", "name", name),
				),
			},
			{
				PreConfig: sleep(t, 5),
				Config:    gaugeMetricConfig(name, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsMetricExists("appoptics_metric.foobar", &metric),
					testAccCheckAppOpticsMetricName(&metric, name),
					testAccCheckAppOpticsMetricType(&metric, typ),
					testAccCheckAppOpticsMetricDescription(&metric, desc2),
					resource.TestCheckResourceAttr(
						"appoptics_metric.foobar", "name", name),
				),
			},
		},
	})

	name = fmt.Sprintf("tftest-metric-%s", acctest.RandString(10))
	typ = "composite"
	desc1 = fmt.Sprintf("A test %s metric", typ)
	desc2 = fmt.Sprintf("An updated test %s metric", typ)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsMetricDestroy,
		Steps: []resource.TestStep{
			{
				Config: compositeMetricConfig(name, typ, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsMetricExists("appoptics_metric.foobar", &metric),
					testAccCheckAppOpticsMetricName(&metric, name),
					testAccCheckAppOpticsMetricType(&metric, typ),
					resource.TestCheckResourceAttr(
						"appoptics_metric.foobar", "name", name),
				),
			},
			{
				PreConfig: sleep(t, 5),
				Config:    compositeMetricConfig(name, typ, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsMetricExists("appoptics_metric.foobar", &metric),
					testAccCheckAppOpticsMetricName(&metric, name),
					testAccCheckAppOpticsMetricType(&metric, typ),
					testAccCheckAppOpticsMetricDescription(&metric, desc2),
					resource.TestCheckResourceAttr(
						"appoptics_metric.foobar", "name", name),
				),
			},
		},
	})
}

func testAccCheckAppOpticsMetricDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*appoptics.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_metric" {
			continue
		}

		_, err := client.MetricsService().Retrieve(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Metric still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsMetricName(metric *appoptics.Metric, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if metric.Name == "" || metric.Name != name {
			return fmt.Errorf("Bad name: %s", metric.Name)
		}

		return nil
	}
}

func testAccCheckAppOpticsMetricDescription(metric *appoptics.Metric, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if metric.Description == "" || metric.Description != desc {
			return fmt.Errorf("Bad description: %s", metric.Description)
		}

		return nil
	}
}

func testAccCheckAppOpticsMetricType(metric *appoptics.Metric, wantType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if metric.Type == "" || metric.Type != wantType {
			return fmt.Errorf("Bad metric type: %s. Expected: %s", metric.Type, wantType)
		}

		return nil
	}
}

func testAccCheckAppOpticsMetricExists(n string, metric *appoptics.Metric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Metric ID is set")
		}

		client := testAccProvider.Meta().(*appoptics.Client)

		foundMetric, err := client.MetricsService().Retrieve(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundMetric.Name == "" || foundMetric.Name != rs.Primary.ID {
			return fmt.Errorf("Metric not found")
		}

		*metric = *foundMetric

		return nil
	}
}

func gaugeMetricConfig(name, desc string) string {
	return strings.TrimSpace(fmt.Sprintf(`
    resource "appoptics_metric" "foobar" {
        type = "gauge"
        name = "%s"
        description = "%s"
        attributes {
          display_stacked = true
        }
    }`, name, desc))
}

func compositeMetricConfig(name, typ, desc string) string {
	return strings.TrimSpace(fmt.Sprintf(`
    resource "appoptics_metric" "foobar" {
        type = "composite"
        name = "%s"
        description = "%s"
        composite = "s(\"librato.cpu.percent.user\", {\"environment\" : \"prod\", \"service\": \"api\"})"
        attributes {
          display_stacked = true
        }
    }`, name, desc))
}
