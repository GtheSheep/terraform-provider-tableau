package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_user" "test" {
  name = "name"
  full_name = "full_name"
  email = "test@test.test"
  site_role = "Viewer"
  auth_setting = "SAML"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_user.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "full_name"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "email"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "site_role"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "auth_setting"),
					resource.TestCheckResourceAttr("tableau_user.test", "name", "name"),
					resource.TestCheckResourceAttr("tableau_user.test", "full_name", "full_name"),
					resource.TestCheckResourceAttr("tableau_user.test", "email", "test@test.test"),
					resource.TestCheckResourceAttr("tableau_user.test", "site_role", "Viewer"),
					resource.TestCheckResourceAttr("tableau_user.test", "auth_setting", "SAML"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "tableau_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
			resource "tableau_user" "test" {
              name = "name_update"
              full_name = "full_name_update"
              email = "test_update@test.test"
              site_role = "Viewer"
              auth_setting = "ServerDefault"
            }
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_user.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "full_name"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "email"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "site_role"),
					resource.TestCheckResourceAttrSet("tableau_user.test", "auth_setting"),
					resource.TestCheckResourceAttr("tableau_user.test", "name", "name_update"),
					resource.TestCheckResourceAttr("tableau_user.test", "full_name", "full_name_update"),
					resource.TestCheckResourceAttr("tableau_user.test", "email", "test_update@test.test"),
					resource.TestCheckResourceAttr("tableau_user.test", "site_role", "Viewer"),
					resource.TestCheckResourceAttr("tableau_user.test", "auth_setting", "ServerDefault"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
