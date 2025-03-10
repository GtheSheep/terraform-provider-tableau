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
	_ resource.Resource                = &datasourcePermissionResource{}
	_ resource.ResourceWithConfigure   = &datasourcePermissionResource{}
	_ resource.ResourceWithImportState = &datasourcePermissionResource{}
)

func NewDatasourcePermissionResource() resource.Resource {
	return &datasourcePermissionResource{}
}

type datasourcePermissionResource struct {
	client *Client
}

type datasourcePermissionResourceModel struct {
	ID             types.String `tfsdk:"id"`
	DatasourceID   types.String `tfsdk:"datasource_id"`
	UserID         types.String `tfsdk:"user_id"`
	GroupID        types.String `tfsdk:"group_id"`
	CapabilityName types.String `tfsdk:"capability_name"`
	CapabilityMode types.String `tfsdk:"capability_mode"`
}

func (r *datasourcePermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datasource_permission"
}

func (r *datasourcePermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"datasource_id": schema.StringAttribute{
				Required:    true,
				Description: "Datasource ID",
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
				Description: "The capability to assign permissions to, one of ChangePermissions/Connect/Delete/ExportXml/Read/Write/SaveAs",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"ChangePermissions",
						"Connect",
						"Delete",
						"ExportXml",
						"Read",
						"Write",
						"SaveAs",
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

func (r *datasourcePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan datasourcePermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	datasourceID := plan.DatasourceID.ValueString()
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
	datasourcePermissions := DatasourcePermissions{
		GranteeCapabilities: []GranteeCapability{granteeCapability},
	}

	_, err := r.client.CreateDatasourcePermissions(datasourceID, datasourcePermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating datasource permission",
			"Could not create datasource permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(getDatasourcePermissionID(datasourceID, entityType, entityID, capability.Name, capability.Mode))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *datasourcePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state datasourcePermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getDatasourcePermissionFromID(state.ID.ValueString())
	datasourcePermission, err := r.client.GetDatasourcePermission(permission.DatasourceID, permission.EntityID, permission.EntityType, permission.CapabilityName, permission.CapabilityMode)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if datasourcePermission.EntityType == "users" {
		state.UserID = types.StringValue(datasourcePermission.EntityID)
	} else {
		state.GroupID = types.StringValue(datasourcePermission.EntityID)
	}
	state.ID = types.StringValue(getDatasourcePermissionID(datasourcePermission.DatasourceID, datasourcePermission.EntityType, datasourcePermission.EntityID, datasourcePermission.CapabilityName, datasourcePermission.CapabilityMode))
	state.DatasourceID = types.StringValue(datasourcePermission.DatasourceID)
	state.CapabilityName = types.StringValue(datasourcePermission.CapabilityName)
	state.CapabilityMode = types.StringValue(datasourcePermission.CapabilityMode)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *datasourcePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan datasourcePermissionResourceModel
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

func (r *datasourcePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state datasourcePermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getDatasourcePermissionFromID(state.ID.ValueString())
	if permission.EntityType == "users" {
		err := r.client.DeleteDatasourcePermission(&permission.EntityID, nil, permission.DatasourceID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Datasource",
				"Could not delete datasource, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		err := r.client.DeleteDatasourcePermission(nil, &permission.EntityID, permission.DatasourceID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Datasource",
				"Could not delete datasource, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *datasourcePermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *datasourcePermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getDatasourcePermissionID(datasourceID, entityType, entityID, capabilityName, capabilityMode string) string {
	return fmt.Sprintf("datasources/%s/permissions/%s/%s/%s/%s", datasourceID, entityType, entityID, capabilityName, capabilityMode)
}

func getDatasourcePermissionFromID(datasourcePermissionID string) DatasourcePermission {
	parts := strings.Split(datasourcePermissionID, "/")
	return DatasourcePermission{
		DatasourceID:   parts[1],
		EntityID:       parts[4],
		EntityType:     parts[3],
		CapabilityName: parts[5],
		CapabilityMode: parts[6],
	}
}
