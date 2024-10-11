package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSiteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_site" "test" {
  name = "test"
  content_url = "moo"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_site.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_site.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_site.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_site.test", "content_url"),
					resource.TestCheckResourceAttr("tableau_site.test", "name", "test"),
					resource.TestCheckResourceAttr("tableau_site.test", "content_url", "moo"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "tableau_site.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "tableau_site" "test" {
  name = "test_new"
  content_url = "moo_new"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_site.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_site.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_site.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_site.test", "content_url"),
					resource.TestCheckResourceAttr("tableau_site.test", "name", "test_new"),
					resource.TestCheckResourceAttr("tableau_site.test", "content_url", "moo_new"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
