package tableau

import (
	"context"
	"fmt"
	"slices"
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
	_ resource.Resource                = &projectPermissionResource{}
	_ resource.ResourceWithConfigure   = &projectPermissionResource{}
	_ resource.ResourceWithImportState = &projectPermissionResource{}
)

func NewProjectPermissionResource() resource.Resource {
	return &projectPermissionResource{}
}

type projectPermissionResource struct {
	client *Client
}

type projectPermissionResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ProjectID      types.String `tfsdk:"project_id"`
	UserID         types.String `tfsdk:"user_id"`
	GroupID        types.String `tfsdk:"group_id"`
	CapabilityName types.String `tfsdk:"capability_name"`
	CapabilityMode types.String `tfsdk:"capability_mode"`
}

func (r *projectPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_permission"
}

func (r *projectPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "Project ID",
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
				Description: "The capability to assign permissions to, one of ProjectLeader/Read/Write",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"ProjectLeader",
						"Read",
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

func (r *projectPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
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
	if plan.UserID.ValueString() != "" {
		granteeCapability.User = &User{ID: entityID}
	} else {
		entityID = plan.GroupID.ValueString()
		entityType = "groups"
		granteeCapability.Group = &Group{ID: entityID}
	}
	projectPermissions := GranteeCapabilities{
		GranteeCapabilities: []GranteeCapability{granteeCapability},
	}

	_, err := r.client.CreateProjectPermissions(projectID, projectPermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project permission",
			"Could not create project permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(getProjectPermissionID(projectID, entityType, entityID, capability.Name, capability.Mode))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := getProjectPermissionFromID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Project",
			err.Error(),
		)
		return
	}
	projectPermission, err := r.client.GetProjectPermission(permission.ProjectID, permission.EntityID, permission.EntityType, permission.CapabilityName, permission.CapabilityMode)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if projectPermission.EntityType == "users" {
		state.UserID = types.StringValue(projectPermission.EntityID)
	} else {
		state.GroupID = types.StringValue(projectPermission.EntityID)
	}
	state.ID = types.StringValue(getProjectPermissionID(projectPermission.ProjectID, projectPermission.EntityType, projectPermission.EntityID, projectPermission.CapabilityName, projectPermission.CapabilityMode))
	state.ProjectID = types.StringValue(projectPermission.ProjectID)
	state.CapabilityName = types.StringValue(projectPermission.CapabilityName)
	state.CapabilityMode = types.StringValue(projectPermission.CapabilityMode)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectPermissionResourceModel
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

func (r *projectPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := getProjectPermissionFromID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Project",
			err.Error(),
		)
		return
	}
	if permission.EntityType == "users" {
		err := r.client.DeleteProjectPermission(&permission.EntityID, nil, permission.ProjectID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Project",
				"Could not delete project, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		err := r.client.DeleteProjectPermission(nil, &permission.EntityID, permission.ProjectID, permission.CapabilityName, permission.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Project",
				"Could not delete project, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *projectPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *projectPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getProjectPermissionID(projectID, entityType, entityID, capabilityName, capabilityMode string) string {
	return fmt.Sprintf("projects/%s/permissions/%s/%s/%s/%s", projectID, entityType, entityID, capabilityName, capabilityMode)
}

func getProjectPermissionFromID(projectPermissionID string) (*ProjectPermission, error) {
	parts := strings.Split(projectPermissionID, "/")
	if len(parts) != 7 {
		return nil, fmt.Errorf("wrong number of items in ID (%d vs. 7) in %s", len(parts), projectPermissionID)
	}
	perms := &ProjectPermission{
		ProjectID:      parts[1],
		EntityID:       parts[4],
		EntityType:     parts[3],
		CapabilityName: parts[5],
		CapabilityMode: parts[6],
	}
	entityTypes := []string{"groups", "users"}
	if !slices.Contains(entityTypes, perms.EntityType) {
		return nil, fmt.Errorf("unknown entity type (%s) not in: %s", perms.EntityType, strings.Join(entityTypes, ", "))
	}
	return perms, nil
}
