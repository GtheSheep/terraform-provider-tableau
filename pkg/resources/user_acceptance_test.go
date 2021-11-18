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

func TestAccTableauUserResource(t *testing.T) {

	userName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	userName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	email := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTableauUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTableauUserResourceBasicConfig(userName, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableauUserExists("tableau_user.test_user"),
					resource.TestCheckResourceAttr("tableau_user.test_user", "name", userName),
				),
			},
			// RENAME
			{
				Config: testAccTableauUserResourceBasicConfig(userName2, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableauUserExists("tableau_user.test_user"),
					resource.TestCheckResourceAttr("tableau_user.test_user", "name", userName2),
				),
			},
			// IMPORT
			{
				ResourceName:            "tableau_user.test_user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccTableauUserResourceBasicConfig(userName, email string) string {
	return fmt.Sprintf(`
resource "tableau_user" "test_user" {
  name  = "%s"
  email = "%s"
  fullName = "test user"
  site_role = "Creator"
  auth_setting = "ServerDefault"
}
`, userName, email)
}

func testAccCheckTableauUserExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*tableau.Client)
		_, err := apiClient.GetUser(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckTableauUserDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tableau.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tableau_user" {
			continue
		}
		_, err := apiClient.GetUser(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("User still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
