package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/gthesheep/terraform-provider-tableau/pkg/tableau"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTableauGroupResource(t *testing.T) {

	groupName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	groupName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTableauGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTableauGroupResourceBasicConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableauGroupExists("tableau_group.test_group"),
					resource.TestCheckResourceAttr("tableau_group.test_group", "name", groupName),
				),
			},
			// RENAME
			{
				Config: testAccTableauGroupResourceBasicConfig(groupName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableauGroupExists("tableau_group.tableau_group"),
					resource.TestCheckResourceAttr("tableau_group.tableau_group", "name", groupName2),
				),
			},
			// UPDATE
			{
				Config: testAccTableauGroupResourceFullConfig(groupName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableauGroupExists("tableau_group.tableau_group"),
					resource.TestCheckResourceAttr("tableau_group.tableau_group", "name", groupName2),
					resource.TestCheckResourceAttr("tableau_group.tableau_group", "minimum_site_role", "ExplorerCanPublish"),
				),
			},
			// IMPORT
			{
				ResourceName:            "tableau_group.test_group",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccTableauGroupResourceBasicConfig(groupName string) string {
	return fmt.Sprintf(`
resource "tableau_group" "test_group" {
  name  = "%s"
}
`, groupName)
}

func testAccTableauGroupResourceFullConfig(groupName string) string {
	return fmt.Sprintf(`
resource "tableau_group" "test_group" {
  name  = "%s"
  minimum_site_role = "ExplorerCanPublish"
}
`, groupName)
}

func testAccCheckTableauGroupExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*tableau.Client)
		_, err := apiClient.GetGroup(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckTableauGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tableau.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tableau_group" {
			continue
		}
		_, err := apiClient.GetGroup(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Group still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
