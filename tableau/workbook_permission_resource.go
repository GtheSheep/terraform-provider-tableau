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
	_ resource.Resource                = &workbookPermissionResource{}
	_ resource.ResourceWithConfigure   = &workbookPermissionResource{}
	_ resource.ResourceWithImportState = &workbookPermissionResource{}
)

func NewWorkbookPermissionResource() resource.Resource {
	return &workbookPermissionResource{}
}

type workbookPermissionResource struct {
	client *Client
}

type workbookPermissionResourceModel struct {
	ID             types.String `tfsdk:"id"`
	WorkbookID     types.String `tfsdk:"workbook_id"`
	UserID         types.String `tfsdk:"user_id"`
	GroupID        types.String `tfsdk:"group_id"`
	CapabilityName types.String `tfsdk:"capability_name"`
	CapabilityMode types.String `tfsdk:"capability_mode"`
}

func (r *workbookPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workbook_permission"
}

func (r *workbookPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workbook_id": schema.StringAttribute{
				Required:    true,
				Description: "Workbook ID",
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
				Description: "The capability to assign permissions to, one of AddComment/ChangeHierarchy/ChangePermissions/CreateRefreshMetrics/Delete/ExportData/ExportImage/ExportXml/Filter/Read/RunExplainData/ShareView/ViewComments/ViewUnderlyingData/WebAuthoring/Write",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"AddComment",
						"ChangeHierarchy",
						"ChangePermissions",
						"CreateRefreshMetrics",
						"Delete",
						"ExportData",
						"ExportImage",
						"ExportXml",
						"Filter",
						"Read",
						"RunExplainData",
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

func (r *workbookPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workbookPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workbookID := plan.WorkbookID.ValueString()
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
	workbookPermissions := WorkbookPermissions{
		GranteeCapabilities: []GranteeCapability{granteeCapability},
	}

	_, err := r.client.CreateWorkbookPermissions(workbookID, workbookPermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workbook permission",
			"Could not create workbook permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(getWorkbookPermissionID(workbookID, entityType, entityID, capability.Name, capability.Mode))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *workbookPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workbookPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getWorkbookPermissionFromID(state.ID.ValueString())
	workbookPermission, err := r.client.GetWorkbookPermission(permission.WorkbookID, permission.EntityID, permission.EntityType, permission.CapabilityName, permission.CapabilityMode)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if workbookPermission.EntityType == "users" {
		state.UserID = types.StringValue(workbookPermission.EntityID)
	} else {
		state.GroupID = types.StringValue(workbookPermission.EntityID)
	}
	state.ID = types.StringValue(getWorkbookPermissionID(workbookPermission.WorkbookID, workbookPermission.EntityType, workbookPermission.EntityID, workbookPermission.CapabilityName, workbookPermission.CapabilityMode))
	state.WorkbookID = types.StringValue(workbookPermission.WorkbookID)
	state.CapabilityName = types.StringValue(workbookPermission.CapabilityName)
	state.CapabilityMode = types.StringValue(workbookPermission.CapabilityMode)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *workbookPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workbookPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *workbookPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workbookPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission := getWorkbookPermissionFromID(state.ID.ValueString())
	if permission.EntityType == "users" {
		err := r.client.DeleteWorkbookPermission(&permission.EntityID, nil, permission.WorkbookID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Workbook",
				"Could not delete workbook, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		err := r.client.DeleteWorkbookPermission(nil, &permission.EntityID, permission.WorkbookID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Workbook",
				"Could not delete workbook, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *workbookPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *workbookPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getWorkbookPermissionID(workbookID, entityType, entityID, capabilityName, capabilityMode string) string {
	return fmt.Sprintf("workbooks/%s/permissions/%s/%s/%s/%s", workbookID, entityType, entityID, capabilityName, capabilityMode)
}

func getWorkbookPermissionFromID(workbookPermissionID string) WorkbookPermission {
	parts := strings.Split(workbookPermissionID, "/")
	return WorkbookPermission{
		WorkbookID:     parts[1],
		EntityID:       parts[4],
		EntityType:     parts[3],
		CapabilityName: parts[5],
		CapabilityMode: parts[6],
	}
}
