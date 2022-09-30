package appoptics

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/appoptics/appoptics-api-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAppOpticsDashboardBasic(t *testing.T) {
	var space appoptics.Space
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsDashboardConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsDashboardExists("appoptics_dashboard.foobar", &space),
					resource.TestCheckResourceAttr(
						"appoptics_dashboard.foobar", "name", name),
				),
			},
		},
	})
}

func testAccCheckAppOpticsDashboardDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*appoptics.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_dashboard" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		_, err = client.SpacesService().Retrieve(id)
		if err == nil {
			return fmt.Errorf("Space still exists")
		}
	}

	return nil
}

func testAccCheckAppOpticsDashboardExists(n string, space *appoptics.Space) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Space ID is set")
		}

		client := testAccProvider.Meta().(*appoptics.Client)

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		foundSpace, err := client.SpacesService().Retrieve(id)

		if err != nil {
			return err
		}

		if foundSpace.ID != id {
			return fmt.Errorf("Space not found")
		}

		return nil
	}
}

func testAccCheckAppOpticsDashboardConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "appoptics_dashboard" "foobar" {
    name = "%s"
}`, name)
}
