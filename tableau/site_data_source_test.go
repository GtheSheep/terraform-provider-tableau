package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSiteDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "tableau_site" "test_site" {
                  name = "test_site"
                  content_url = "moo"
                }
                data "tableau_site" "test" {
                    id = tableau_site.test_site.id
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tableau_site.test", "name", "test"),
					resource.TestCheckResourceAttr("data.tableau_site.test", "content_url", "moo"),
				),
			},
		},
	})
}
