package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectPermissionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "tableau_project" "test_perm_project" {
  name = "test_project_permission"
  content_permissions = "ManagedByOwner"
}
resource "tableau_user" "new_person" {
	name = "test_person_project_perms@test.test"
  full_name = "test_person_project_perms@test.test"
  email = "test_person_project_perms@test.test"
  site_role = "Creator"
  auth_setting = "SAML"
}
resource "tableau_project_permission" "test_permission" {
  project_id = tableau_project.test_perm_project.id
	user_id = tableau_user.new_person.id
  capability_name = "Write"
	capability_mode = "Deny"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tableau_project_permission.test_permission", "id"),
					resource.TestCheckResourceAttrSet("tableau_project_permission.test_permission", "project_id"),
					resource.TestCheckResourceAttrSet("tableau_project_permission.test_permission", "user_id"),
					resource.TestCheckResourceAttrSet("tableau_project_permission.test_permission", "capability_name"),
					resource.TestCheckResourceAttrSet("tableau_project_permission.test_permission", "capability_mode"),
					resource.TestCheckResourceAttr("tableau_project_permission.test_permission", "capability_name", "Write"),
					resource.TestCheckResourceAttr("tableau_project_permission.test_permission", "capability_mode", "Deny"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "tableau_project_permission.test_permission",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			// Delete testing automatically occurs in TestCase
		},
	})
}
