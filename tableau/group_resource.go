package tableau

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	client *Client
}

type groupResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	MinimumSiteRole types.String `tfsdk:"minimum_site_role"`
	OnDemandAccess  types.Bool   `tfsdk:"on_demand_access"`
	LastUpdated     types.String `tfsdk:"last_updated"`
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for group",
			},
			"minimum_site_role": schema.StringAttribute{
				Optional:    true,
				Description: "Minimum site role for the group",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Creator",
						"Explorer",
						"Interactor",
						"Publisher",
						"ExplorerCanPublish",
						"ServerAdministrator",
						"SiteAdministratorExplorer",
						"SiteAdministratorCreator",
						"Unlicensed",
						"Viewer",
					}...),
				},
			},
			"on_demand_access": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable on-demand access for embedded Tableau content.",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var onDemandAccess *bool
	if !plan.OnDemandAccess.IsNull() {
		value := plan.OnDemandAccess.ValueBool()
		onDemandAccess = &value
	}

	var minimumSiteRole *string
	if !plan.MinimumSiteRole.IsNull() {
		value := plan.MinimumSiteRole.ValueString()
		if value != "" {
			minimumSiteRole = &value
		}
	}

	createdGroup, err := r.client.CreateGroup(
		plan.Name.ValueString(),
		minimumSiteRole,
		onDemandAccess,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			"Could not create group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(createdGroup.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetGroup(state.ID.ValueString())
	if err != nil {
		// Handle the case where the resource does not exist
		if strings.Contains(err.Error(), "Did not find group") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Tableau Group",
			"Could not read group: "+err.Error(),
		)
		return
	}

	if group == nil {
		// If the group doesn't exist, remove it from the state
		resp.State.RemoveResource(ctx)
		return
	}

	// Update the state with the current resource details
	state.ID = types.StringValue(group.ID)
	state.Name = types.StringValue(group.Name)
	if group.OnDemandAccess != nil {
		state.OnDemandAccess = types.BoolValue(*group.OnDemandAccess)
	} else {
		state.OnDemandAccess = types.BoolNull()
	}
	if group.Import != nil && group.Import.MinimumSiteRole != nil {
		state.MinimumSiteRole = types.StringValue(*group.Import.MinimumSiteRole)
	} else {
		state.MinimumSiteRole = types.StringNull()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var onDemandAccess *bool
	if !plan.OnDemandAccess.IsNull() {
		value := plan.OnDemandAccess.ValueBool()
		onDemandAccess = &value
	}

	var minimumSiteRole *string
	if !plan.MinimumSiteRole.IsNull() {
		value := plan.MinimumSiteRole.ValueString()
		if value != "" {
			minimumSiteRole = &value
		}
	}

	updatedGroup, err := r.client.UpdateGroup(
		plan.ID.ValueString(),
		plan.Name.ValueString(),
		minimumSiteRole,
		onDemandAccess,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Tableau Group",
			"Could not update group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(updatedGroup.Name)

	if updatedGroup.MinimumSiteRole != "" {
		plan.MinimumSiteRole = types.StringValue(updatedGroup.MinimumSiteRole)
	} else {
		plan.MinimumSiteRole = types.StringNull()
	}

	if updatedGroup.OnDemandAccess != nil {
		plan.OnDemandAccess = types.BoolValue(*updatedGroup.OnDemandAccess)
	} else {
		plan.OnDemandAccess = types.BoolNull()
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Group",
			"Could not delete group, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
