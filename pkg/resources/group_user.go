package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gthesheep/terraform-provider-tableau/pkg/tableau"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var groupUserSchema = map[string]*schema.Schema{
	"group_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Group ID",
	},
	"user_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "User ID",
	},
}

func ResourceGroupUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupUserCreate,
		ReadContext:   resourceGroupUserRead,
		UpdateContext: resourceGroupUserUpdate,
		DeleteContext: resourceGroupUserDelete,

		Schema: groupUserSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGroupUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)

	var diags diag.Diagnostics

	groupUserID := d.Id()
	strs := strings.Split(groupUserID, ":")
	groupID := strs[0]
	userID := strs[1]

	groupUser, err := c.GetGroupUser(groupID, userID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_id", groupUser.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGroupUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)

	var diags diag.Diagnostics

	groupID := d.Get("group_id").(string)
	userID := d.Get("user_id").(string)

	g, err := c.CreateGroupUser(groupID, userID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s", groupID, *g.ID))

	resourceGroupUserRead(ctx, d, m)

	return diags
}

func resourceGroupUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("Group Users do not support updates")
	return resourceGroupUserRead(ctx, d, m)
}

func resourceGroupUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*tableau.Client)
	groupUserID := d.Id()
	strs := strings.Split(groupUserID, ":")
	groupID := strs[0]
	userID := strs[1]

	var diags diag.Diagnostics

	_, err := c.DeleteGroupUser(groupID, userID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
