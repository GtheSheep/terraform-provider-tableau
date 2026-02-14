package tableau

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &viewPermissionResource{}
	_ resource.ResourceWithConfigure   = &viewPermissionResource{}
	_ resource.ResourceWithImportState = &viewPermissionResource{}
)

func NewViewPermissionResource() resource.Resource {
	return &viewPermissionResource{}
}

type viewPermissionResource struct {
	client *Client
}

type viewPermissionResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ViewID         types.String `tfsdk:"view_id"`
	UserID         types.String `tfsdk:"user_id"`
	GroupID        types.String `tfsdk:"group_id"`
	CapabilityName types.String `tfsdk:"capability_name"`
	CapabilityMode types.String `tfsdk:"capability_mode"`
}

func (r *viewPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_view_permission"
}

func (r *viewPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"view_id": schema.StringAttribute{
				Required:    true,
				Description: "View ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Optional:    true,
				Description: "User ID to grant to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Group ID to grant to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"capability_name": schema.StringAttribute{
				Required:    true,
				Description: "The capability to assign permissions to, one of AddComment/ChangePermissions/Delete/ExportData/ExportImage/ExportXml/Filter/Read/ShareView/ViewComments/ViewUnderlyingData/WebAuthoring/Write",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"AddComment",
						"ChangePermissions",
						"Delete",
						"ExportData",
						"ExportImage",
						"ExportXml",
						"Filter",
						"Read",
						"ShareView",
						"ViewComments",
						"ViewUnderlyingData",
						"WebAuthoring",
						"Write",
					}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"capability_mode": schema.StringAttribute{
				Required:    true,
				Description: "Capability mode, Allow or Deny (case sensitive)",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Allow",
						"Deny",
					}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *viewPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan viewPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	viewID := plan.ViewID.ValueString()
	capability := Capability{
		Name: plan.CapabilityName.ValueString(),
		Mode: plan.CapabilityMode.ValueString(),
	}
	capabilities := Capabilities{
		Capabilities: []Capability{capability},
	}
	granteeCapability := GranteeCapability{
		Capabilities: capabilities,
	}

	entityType := "users"
	entityID := plan.UserID.ValueString()
	if entityID != "" {
		granteeCapability.User = &User{ID: entityID}
	} else {
		entityID = plan.GroupID.ValueString()
		entityType = "groups"
		granteeCapability.Group = &Group{ID: entityID}
	}
	viewPermissions := ViewPermissions{
		GranteeCapabilities: []GranteeCapability{granteeCapability},
	}

	_, err := r.client.CreateViewPermissions(viewID, viewPermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating view permission",
			"Could not create view permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(getViewPermissionID(viewID, entityType, entityID, capability.Name, capability.Mode))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *viewPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state viewPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getViewPermissionFromID(state.ID.ValueString())
	viewPermission, err := r.client.GetViewPermission(permission.ViewID, permission.EntityID, permission.EntityType, permission.CapabilityName, permission.CapabilityMode)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if viewPermission.EntityType == "users" {
		state.UserID = types.StringValue(viewPermission.EntityID)
	} else {
		state.GroupID = types.StringValue(viewPermission.EntityID)
	}
	state.ID = types.StringValue(getViewPermissionID(viewPermission.ViewID, viewPermission.EntityType, viewPermission.EntityID, viewPermission.CapabilityName, viewPermission.CapabilityMode))
	state.ViewID = types.StringValue(viewPermission.ViewID)
	state.CapabilityName = types.StringValue(viewPermission.CapabilityName)
	state.CapabilityMode = types.StringValue(viewPermission.CapabilityMode)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *viewPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan viewPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *viewPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state viewPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getViewPermissionFromID(state.ID.ValueString())
	if permission.EntityType == "users" {
		err := r.client.DeleteViewPermission(&permission.EntityID, nil, permission.ViewID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau View",
				"Could not delete view, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		err := r.client.DeleteViewPermission(nil, &permission.EntityID, permission.ViewID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau View",
				"Could not delete view, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *viewPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *viewPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getViewPermissionID(viewID, entityType, entityID, capabilityName, capabilityMode string) string {
	return fmt.Sprintf("views/%s/permissions/%s/%s/%s/%s", viewID, entityType, entityID, capabilityName, capabilityMode)
}

func getViewPermissionFromID(viewPermissionID string) ViewPermission {
	parts := strings.Split(viewPermissionID, "/")
	return ViewPermission{
		ViewID:         parts[1],
		EntityID:       parts[4],
		EntityType:     parts[3],
		CapabilityName: parts[5],
		CapabilityMode: parts[6],
	}
}
