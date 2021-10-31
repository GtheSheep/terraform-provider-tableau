package data_sources

import (
	"context"

	"github.com/gthesheep/terraform-provider-tableau/pkg/tableau"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userSchema = map[string]*schema.Schema{
	"user_id": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "ID of the user",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name for the user",
	},
	"full_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Full name for user",
	},
	"site_role": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Site role for the user",
	},
	"auth_setting": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Auth setting for the user",
	},
}

func DatasourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceUserRead,
		Schema:      userSchema,
	}
}

func datasourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)

	var diags diag.Diagnostics

	userID := d.Get("user_id").(string)

	user, err := c.GetUser(userID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", user.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("full_name", user.FullName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("site_role", user.SiteRole); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auth_setting", user.AuthSetting); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(userID)

	return diags
}
