package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_group" "test" {
  name = "test"
  minimum_site_role = "Viewer"
}
resource "tableau_user" "test" {
  name = "test@test.test"
  full_name = "test@test.test"
  email = "test@test.test"
  site_role = "Viewer"
  auth_setting = "SAML"
}
resource "tableau_group_user" "test" {
    group_id = tableau_group.test.id
    user_id  = tableau_user.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_group_user.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_group_user.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_group_user.test", "group_id"),
					resource.TestCheckResourceAttrSet("tableau_group_user.test", "user_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "tableau_group_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// 			// Update and Read testing
			// 			{
			// 				Config: providerConfig + `
			// resource "tableau_group" "test" {
			//   name = "test"
			//   minimum_site_role = "Explorer"
			// }
			// `,
			// 				Check: resource.ComposeAggregateTestCheckFunc(
			// 					resource.TestCheckResourceAttrSet("tableau_group.test", "id"),
			// 					resource.TestCheckResourceAttrSet("tableau_group.test", "last_updated"),
			// 					resource.TestCheckResourceAttrSet("tableau_group.test", "name"),
			// 					resource.TestCheckResourceAttrSet("tableau_group.test", "minimum_site_role"),
			// 					resource.TestCheckResourceAttr("tableau_group.test", "name", "test"),
			// 					resource.TestCheckResourceAttr("tableau_group.test", "minimum_site_role", "Explorer"),
			// 				),
			// 			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
