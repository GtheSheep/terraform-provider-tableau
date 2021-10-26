package resources

import (
	"context"

	"github.com/gthesheep/terraform-provider-tableau/pkg/tableau"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	minimumSiteRoles = []string{
		"Creator",
		"Explorer",
		"ExplorerCanPublish",
		"SiteAdministratorExplorer",
		"SiteAdministratorCreator",
		"Unlicensed",
		"Viewer",
	}
)

var groupSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Group name",
	},
	"minimum_site_role": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Minimum site role for the group",
	},
	"grant_license_mode": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "onLogin",
		Description: "When to grant license for the group to users",
	},
}

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Schema: groupSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)

	var diags diag.Diagnostics

	groupID := d.Id()

	group, err := c.GetGroup(groupID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("minimum_site_role", group.Import.MinimumSiteRole); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("grant_license_mode", group.Import.GrantLicenseMode); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	minimumSiteRole := d.Get("minimum_site_role").(string)

	g, err := c.CreateGroup(name, minimumSiteRole)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*g.ID)

	resourceGroupRead(ctx, d, m)

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)
	groupID := d.Id()

	if d.HasChange("name") || d.HasChange("minimum_site_role") {
		name := d.Get("name").(string)
		minimumSiteRole := d.Get("minimum_site_role").(string)

		_, err := c.UpdateGroup(groupID, name, minimumSiteRole)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)
	groupID := d.Id()

	var diags diag.Diagnostics

	_, err := c.DeleteGroup(groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
