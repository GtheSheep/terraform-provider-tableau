package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "tableau_group" "test" {
                    name = "test"
                    minimum_site_role = "Viewer"
                }
                data "tableau_group" "test" {
                    id = tableau_group.test.id
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tableau_group.test", "name", "test"),
					resource.TestCheckResourceAttr("data.tableau_group.test", "minimum_site_role", "Viewer"),
				),
			},
		},
	})
}
