package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/akahn/go-librato/librato"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppOpticsService_Basic(t *testing.T) {
	var service librato.Service

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsServiceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsServiceExists("appoptics_service.foobar", &service),
					testAccCheckAppOpticsServiceTitle(&service, "Foo Bar"),
					resource.TestCheckResourceAttr(
						"appoptics_service.foobar", "title", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsService_Updated(t *testing.T) {
	var service librato.Service

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsServiceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsServiceExists("appoptics_service.foobar", &service),
					testAccCheckAppOpticsServiceTitle(&service, "Foo Bar"),
					resource.TestCheckResourceAttr(
						"appoptics_service.foobar", "title", "Foo Bar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckAppOpticsServiceConfig_new_value,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsServiceExists("appoptics_service.foobar", &service),
					testAccCheckAppOpticsServiceTitle(&service, "Bar Baz"),
					resource.TestCheckResourceAttr(
						"appoptics_service.foobar", "title", "Bar Baz"),
				),
			},
		},
	})
}

func testAccCheckAppOpticsServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*librato.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_service" {
			continue
		}

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		_, _, err = client.Services.Get(uint(id))

		if err == nil {
			return fmt.Errorf("Service still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsServiceTitle(service *librato.Service, title string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if service.Title == nil || *service.Title != title {
			return fmt.Errorf("Bad title: %s", *service.Title)
		}

		return nil
	}
}

func testAccCheckAppOpticsServiceExists(n string, service *librato.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service ID is set")
		}

		client := testAccProvider.Meta().(*librato.Client)

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		foundService, _, err := client.Services.Get(uint(id))

		if err != nil {
			return err
		}

		if foundService.ID == nil || *foundService.ID != uint(id) {
			return fmt.Errorf("Service not found")
		}

		*service = *foundService

		return nil
	}
}

const testAccCheckAppOpticsServiceConfig_basic = `
resource "appoptics_service" "foobar" {
    title = "Foo Bar"
    type = "mail"
    settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}`

const testAccCheckAppOpticsServiceConfig_new_value = `
resource "appoptics_service" "foobar" {
    title = "Bar Baz"
    type = "mail"
    settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}`
