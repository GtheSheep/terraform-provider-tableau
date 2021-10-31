package data_sources

import (
	"context"
	"strconv"

	"github.com/gthesheep/terraform-provider-tableau/pkg/tableau"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var groupSchema = map[string]*schema.Schema{
	"group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "ID of the group",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name for the group",
	},
}

func DatasourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceGroupRead,
		Schema:      groupSchema,
	}
}

func datasourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)

	var diags diag.Diagnostics

	groupID := strconv.Itoa(d.Get("group_id").(int))

	group, err := c.GetGroup(groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupID)

	return diags
}
