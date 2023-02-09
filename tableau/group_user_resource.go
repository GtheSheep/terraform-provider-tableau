package tableau

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &groupUserResource{}
	_ resource.ResourceWithConfigure = &groupUserResource{}
)

func NewGroupUserResource() resource.Resource {
	return &groupUserResource{}
}

type groupUserResource struct {
	client *Client
}

type groupUserResourceModel struct {
	ID          types.String `tfsdk:"id"`
	GroupID     types.String `tfsdk:"group_id"`
	UserID      types.String `tfsdk:"user_id"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *groupUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_user"
}

func (r *groupUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Group identifier",
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "User identifier",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *groupUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupUser := User{
		ID: string(plan.UserID.ValueString()),
	}

	_, err := r.client.CreateGroupUser(plan.GroupID.ValueString(), groupUser.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating groupUser",
			"Could not create groupUser, unexpected error: "+err.Error(),
		)
		return
	}

	combinedID := GetCombinedID(plan.GroupID.ValueString(), groupUser.ID)
	plan.ID = types.StringValue(combinedID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := state.GroupID.ValueString()
	userID := state.UserID.ValueString()
	if (groupID == "") || (userID == "") {
		groupID, userID = GetIDsFromCombinedID(state.ID.ValueString())
	}

	groupUser, err := r.client.GetGroupUser(groupID, userID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Group User",
			"Could not read Tableau Group ID/ User ID "+state.GroupID.ValueString()+"/"+state.UserID.ValueString()+": "+err.Error(),
		)
		return
	}

	combinedID := GetCombinedID(groupID, userID)
	state.ID = types.StringValue(combinedID)
	state.GroupID = types.StringValue(groupID)
	state.UserID = types.StringValue(groupUser.ID)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Group Users do not support updates")
	return
}

func (r *groupUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroupUser(state.GroupID.ValueString(), state.UserID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Group User",
			"Could not delete groupUser, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *groupUserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *groupUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
