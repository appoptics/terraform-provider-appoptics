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

func TestAccAppOpticsSpaceBasic(t *testing.T) {
	var space appoptics.Space
	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOpticsSpaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppOpticsSpaceConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOpticsSpaceExists("appoptics_space.foobar", &space),
					testAccCheckAppOpticsSpaceAttributes(&space, name),
					resource.TestCheckResourceAttr(
						"appoptics_space.foobar", "name", name),
				),
			},
		},
	})
}

func testAccCheckAppOpticsSpaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*appoptics.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appoptics_space" {
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

func testAccCheckAppOpticsSpaceAttributes(space *appoptics.Space, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if space.Name == "" || space.Name != name {
			return fmt.Errorf("Bad name: %s", space.Name)
		}

		return nil
	}
}

func testAccCheckAppOpticsSpaceExists(n string, space *appoptics.Space) resource.TestCheckFunc {
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

func testAccCheckAppOpticsSpaceConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "appoptics_space" "foobar" {
    name = "%s"
}`, name)
}
