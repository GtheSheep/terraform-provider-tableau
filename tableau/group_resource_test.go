package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_group" "test" {
  name = "test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_group.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "name"),
					resource.TestCheckResourceAttr("tableau_group.test", "name", "test"),
				),
			},
			// Add minimum_site_role
			{
				Config: providerConfig + `
resource "tableau_group" "test" {
  name = "test"
  minimum_site_role = "Viewer"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_group.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "minimum_site_role"),
					resource.TestCheckResourceAttr("tableau_group.test", "name", "test"),
					resource.TestCheckResourceAttr("tableau_group.test", "minimum_site_role", "Viewer"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "tableau_group.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "tableau_group" "test" {
  name = "test"
  minimum_site_role = "Explorer"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_group.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_group.test", "minimum_site_role"),
					resource.TestCheckResourceAttr("tableau_group.test", "name", "test"),
					resource.TestCheckResourceAttr("tableau_group.test", "minimum_site_role", "Explorer"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
