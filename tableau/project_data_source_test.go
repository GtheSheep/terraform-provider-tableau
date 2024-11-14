package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
								resource "tableau_project" "test" {
									name = "test_project_data_source"
									description = "Test project for data source test"
									content_permissions = "ManagedByOwner"
								}
                data "tableau_project" "test" {
                    id = tableau_project.test.id
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tableau_project.test", "name", "test_project_data_source"),
					resource.TestCheckResourceAttr("data.tableau_project.test", "content_permissions", "ManagedByOwner"),
					resource.TestCheckResourceAttr("data.tableau_project.test", "description", "Test project for data source test"),
				),
			},
		},
	})
}
