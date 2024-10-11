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
				resource "tableau_project" "test_parent" {
                  name = "test_parent"
                  content_permissions = "ManagedByOwner"
                }
				resource "tableau_project" "test" {
                  name = "test"
                  content_permissions = "ManagedByOwner"
                  description = "Moo"
                  parent_project_id = tableau_project.test_parent.id
                }
                data "tableau_project" "test" {
                    id = tableau_project.test.id
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tableau_project.test", "name", "test"),
					resource.TestCheckResourceAttr("data.tableau_project.test", "content_permissions", "ManagedByOwner"),
					resource.TestCheckResourceAttr("data.tableau_project.test", "description", "Moo"),
					resource.TestCheckResourceAttrSet("data.tableau_project.test", "parent_project_id"),
				),
			},
		},
	})
}
