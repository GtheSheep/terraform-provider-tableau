package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_project" "test" {
  name = "test"
  content_permissions = "ManagedByOwner"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_project.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "content_permissions"),
					resource.TestCheckResourceAttr("tableau_project.test", "name", "test"),
					resource.TestCheckResourceAttr("tableau_project.test", "content_permissions", "ManagedByOwner"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "tableau_project.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "tableau_project" "test_parent" {
  name = "test_parent"
  content_permissions = "ManagedByOwner"
}
resource "tableau_project" "test" {
  name = "test"
  description = "Moo"
  content_permissions = "LockedToProject"
  parent_project_id = tableau_project.test_parent.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_project.test", "id"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "last_updated"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "name"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "description"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "content_permissions"),
					resource.TestCheckResourceAttrSet("tableau_project.test", "parent_project_id"),
					resource.TestCheckResourceAttr("tableau_project.test", "name", "test"),
					resource.TestCheckResourceAttr("tableau_project.test", "description", "Moo"),
					resource.TestCheckResourceAttr("tableau_project.test", "content_permissions", "LockedToProject"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
