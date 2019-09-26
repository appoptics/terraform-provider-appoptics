package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAppOpticsServiceBasic(t *testing.T) {
	var service appoptics.Service

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsServiceExists("appoptics_service.foobar", &service),
					resource.TestCheckResourceAttr(
						"appoptics_service.foobar", "title", "Foo Bar"),
				),
			},
		},
	})
}

func TestAccAppOpticsServiceUpdated(t *testing.T) {
	var service appoptics.Service

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAppOpticsServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsServiceExists("appoptics_service.foobar", &service),
					resource.TestCheckResourceAttr(
						"appoptics_service.foobar", "title", "Foo Bar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckAppOpticsServiceConfigNewValue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsServiceExists("appoptics_service.foobar", &service),
					resource.TestCheckResourceAttr(
						"appoptics_service.foobar", "title", "Bar Baz"),
				),
			},
		},
	})
}

func testAccCheckAppOpticsServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*appoptics.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_service" {
			continue
		}

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		_, err = client.ServicesService().Retrieve(int(id))

		if err == nil {
			return fmt.Errorf("Service still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsServiceExists(n string, service *appoptics.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service ID is set")
		}

		client := testAccProvider.Meta().(*appoptics.Client)

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		foundService, err := client.ServicesService().Retrieve(int(id))

		if err != nil {
			return err
		}

		if foundService.ID != int(id) {
			return fmt.Errorf("Service not found")
		}

		return nil
	}
}

const testAccCheckAppOpticsServiceConfigBasic = `
resource "appoptics_service" "foobar" {
    title = "Foo Bar"
    type = "mail"
    settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}`

const testAccCheckAppOpticsServiceConfigNewValue = `
resource "appoptics_service" "foobar" {
    title = "Bar Baz"
    type = "mail"
    settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}`
