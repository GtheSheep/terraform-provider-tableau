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
	_ resource.Resource                = &virtualConnectionPermissionResource{}
	_ resource.ResourceWithConfigure   = &virtualConnectionPermissionResource{}
	_ resource.ResourceWithImportState = &virtualConnectionPermissionResource{}
)

func NewVirtualConnectionPermissionResource() resource.Resource {
	return &virtualConnectionPermissionResource{}
}

type virtualConnectionPermissionResource struct {
	client *Client
}

type virtualConnectionPermissionResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VirtualConnectionID types.String `tfsdk:"virtualConnection_id"`
	UserID              types.String `tfsdk:"user_id"`
	GroupID             types.String `tfsdk:"group_id"`
	CapabilityName      types.String `tfsdk:"capability_name"`
	CapabilityMode      types.String `tfsdk:"capability_mode"`
}

func (r *virtualConnectionPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_connection_permission"
}

func (r *virtualConnectionPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virtual_connection_id": schema.StringAttribute{
				Required:    true,
				Description: "Virtual connection ID",
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
				Description: "The capability to assign permissions to, one of Read/Connect/Overwrite/ChangeHierarchy/Delete/ChangePermissions",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Read",
						"Connect",
						"Overwrite",
						"ChangeHierarchy",
						"Delete",
						"ChangePermissions",
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

func (r *virtualConnectionPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualConnectionPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	virtualConnectionID := string(plan.VirtualConnectionID.ValueString())
	capability := Capability{
		Name: string(plan.CapabilityName.ValueString()),
		Mode: string(plan.CapabilityMode.ValueString()),
	}
	capabilities := Capabilities{
		Capabilities: []Capability{capability},
	}
	granteeCapability := GranteeCapability{
		Capabilities: capabilities,
	}

	entityType := "users"
	entityID := string(plan.UserID.ValueString())
	if plan.UserID.ValueString() != "" {
		granteeCapability.User = &User{ID: entityID}
	} else {
		entityID = string(plan.GroupID.ValueString())
		entityType = "groups"
		granteeCapability.Group = &Group{ID: entityID}
	}
	virtualConnectionPermissions := VirtualConnectionPermissions{
		GranteeCapabilities: []GranteeCapability{granteeCapability},
	}

	_, err := r.client.CreateVirtualConnectionPermissions(virtualConnectionID, virtualConnectionPermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating virtual connection permission",
			"Could not create virtual connection permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(getVirtualConnectionPermissionID(virtualConnectionID, entityType, entityID, capability.Name, capability.Mode))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualConnectionPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualConnectionPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getVirtualConnectionPermissionFromID(state.ID.ValueString())
	virtualConnectionPermission, err := r.client.GetVirtualConnectionPermission(permission.VirtualConnectionID, permission.EntityID, permission.EntityType, permission.CapabilityName, permission.CapabilityMode)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if virtualConnectionPermission.EntityType == "users" {
		state.UserID = types.StringValue(virtualConnectionPermission.EntityID)
	} else {
		state.GroupID = types.StringValue(virtualConnectionPermission.EntityID)
	}
	state.ID = types.StringValue(getVirtualConnectionPermissionID(virtualConnectionPermission.VirtualConnectionID, virtualConnectionPermission.EntityType, virtualConnectionPermission.EntityID, virtualConnectionPermission.CapabilityName, virtualConnectionPermission.CapabilityMode))
	state.VirtualConnectionID = types.StringValue(virtualConnectionPermission.VirtualConnectionID)
	state.CapabilityName = types.StringValue(virtualConnectionPermission.CapabilityName)
	state.CapabilityMode = types.StringValue(virtualConnectionPermission.CapabilityMode)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualConnectionPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan virtualConnectionPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualConnectionPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualConnectionPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getVirtualConnectionPermissionFromID(state.ID.ValueString())
	if permission.EntityType == "users" {
		err := r.client.DeleteVirtualConnectionPermission(&permission.EntityID, nil, permission.VirtualConnectionID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau virtual connection",
				"Could not delete virtual connection, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		err := r.client.DeleteVirtualConnectionPermission(nil, &permission.EntityID, permission.VirtualConnectionID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau virtual connection",
				"Could not delete virtual connection, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *virtualConnectionPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *virtualConnectionPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getVirtualConnectionPermissionID(virtualConnectionID, entityType, entityID, capabilityName, capabilityMode string) string {
	return fmt.Sprintf("virtualConnections/%s/permissions/%s/%s/%s/%s", virtualConnectionID, entityType, entityID, capabilityName, capabilityMode)
}

func getVirtualConnectionPermissionFromID(virtualConnectionPermissionID string) VirtualConnectionPermission {
	parts := strings.Split(virtualConnectionPermissionID, "/")
	return VirtualConnectionPermission{
		VirtualConnectionID: parts[1],
		EntityID:            parts[4],
		EntityType:          parts[3],
		CapabilityName:      parts[5],
		CapabilityMode:      parts[6],
	}
}
