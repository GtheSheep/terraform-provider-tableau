package tableau

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                data "tableau_datasource" "test" {
                    name = "Superstore Datasource"
                }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tableau_datasource.test", "id"),
					resource.TestCheckResourceAttrSet("data.tableau_datasource.test", "owner_id"),
					resource.TestCheckResourceAttrSet("data.tableau_datasource.test", "project_id"),
					resource.TestCheckResourceAttr("data.tableau_datasource.test", "name", "Superstore Datasource"),
					resource.TestCheckResourceAttr("data.tableau_datasource.test", "content_url", "SuperstoreDatasource"),
					resource.TestCheckResourceAttr("data.tableau_datasource.test", "type", "excel-direct"),
				),
			},
		},
	})
}
