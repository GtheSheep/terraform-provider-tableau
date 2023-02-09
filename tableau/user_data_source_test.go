package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "tableau_user" "test" {
                  name = "test@test.test"
                  full_name = "test@test.test"
                  email = "test@test.test"
                  site_role = "Viewer"
                  auth_setting = "SAML"
                }
                data "tableau_user" "test" {
                    id = tableau_user.test.id
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tableau_user.test", "name", "test@test.test"),
					resource.TestCheckResourceAttr("data.tableau_user.test", "email", "test@test.test"),
					resource.TestCheckResourceAttr("data.tableau_user.test", "full_name", "test@test.test"),
					resource.TestCheckResourceAttr("data.tableau_user.test", "site_role", "Viewer"),
					resource.TestCheckResourceAttr("data.tableau_user.test", "auth_setting", "SAML"),
				),
			},
		},
	})
}
