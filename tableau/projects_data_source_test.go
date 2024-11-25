package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                data "tableau_projects" "test" {
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tableau_projects.test", "id"),
					resource.TestCheckResourceAttrSet("data.tableau_projects.test", "projects.#"),
				),
			},
		},
	})
}
