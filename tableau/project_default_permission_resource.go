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
	_                                   resource.Resource                = &projectPermissionResource{}
	_                                   resource.ResourceWithConfigure   = &projectPermissionResource{}
	_                                   resource.ResourceWithImportState = &projectPermissionResource{}
	projectDefaultPermissionTargetTypes                                  = []string{
		"databases",
		"dataroles",
		"datasources",
		"flows",
		"lenses",
		"metrics",
		"tables",
		"virtualconnections",
		"workbooks",
	}
)

func NewProjectDefaultPermissionResource() resource.Resource {
	return &projectDefaultPermissionResource{}
}

type projectDefaultPermissionResource struct {
	client *Client
}

type projectDefaultPermissionResourceModel struct {
	ID             types.String `tfsdk:"id"`
	ProjectID      types.String `tfsdk:"project_id"`
	UserID         types.String `tfsdk:"user_id"`
	GroupID        types.String `tfsdk:"group_id"`
	TargetType     types.String `tfsdk:"target_type"`
	CapabilityName types.String `tfsdk:"capability_name"`
	CapabilityMode types.String `tfsdk:"capability_mode"`
}

func (r *projectDefaultPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_default_permission"
}

func (r *projectDefaultPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"target_type": schema.StringAttribute{
				Required:    true,
				Description: fmt.Sprintf("Target type (%s)", strings.Join(projectDefaultPermissionTargetTypes, ", ")),
				Validators: []validator.String{
					stringvalidator.OneOf(projectDefaultPermissionTargetTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"capability_name": schema.StringAttribute{
				Required:    true,
				Description: "The capability to assign permissions to, one of ProjectLeader/Read/Write",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"capability_mode": schema.StringAttribute{
				Required:    true,
				Description: "Capability mode, Allow or Deny (case sensitive)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *projectDefaultPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := string(plan.ProjectID.ValueString())
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
	projectPermissions := ProjectPermissions{
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

func (r *projectDefaultPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectDefaultPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	perms, err := getProjectDefaultPermissionFromID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Project",
			err.Error(),
		)
		return
	}
	projectDefaultPermission, err := r.client.GetProjectDefaultPermission(perms.ProjectID, perms.EntityID, perms.EntityType, perms.TargetType, perms.CapabilityName, perms.CapabilityMode)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if projectDefaultPermission.EntityType == "users" {
		state.UserID = types.StringValue(projectDefaultPermission.EntityID)
	} else {
		state.GroupID = types.StringValue(projectDefaultPermission.EntityID)
	}
	state.ID = types.StringValue(getProjectDefaultPermissionID(projectDefaultPermission.ProjectID, projectDefaultPermission.EntityType, projectDefaultPermission.EntityID, projectDefaultPermission.TargetType, projectDefaultPermission.CapabilityName, projectDefaultPermission.CapabilityMode))
	state.ProjectID = types.StringValue(projectDefaultPermission.ProjectID)
	state.TargetType = types.StringValue(projectDefaultPermission.TargetType)
	state.CapabilityName = types.StringValue(projectDefaultPermission.CapabilityName)
	state.CapabilityMode = types.StringValue(projectDefaultPermission.CapabilityMode)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectDefaultPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r *projectDefaultPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	perm, err := getProjectDefaultPermissionFromID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Project",
			err.Error(),
		)
		return
	}
	if perm.EntityType == "users" {
		err := r.client.DeleteProjectDefaultPermission(&perm.EntityID, nil, perm.ProjectID, perm.TargetType, perm.CapabilityName, perm.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Project",
				"Could not delete project, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		err := r.client.DeleteProjectDefaultPermission(nil, &perm.EntityID, perm.ProjectID, perm.TargetType, perm.CapabilityName, perm.CapabilityMode)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Tableau Project",
				"Could not delete project, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *projectDefaultPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *projectDefaultPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getProjectDefaultPermissionID(projectID, entityType, entityID, targetType, capabilityName, capabilityMode string) string {
	return fmt.Sprintf(
		"projects/%s/default_permissions/%s",
		projectID,
		strings.Join([]string{entityType, entityID, targetType, capabilityName, capabilityMode}, "/"),
	)
}

func getProjectDefaultPermissionFromID(projectPermissionID string) (*ProjectDefaultPermission, error) {
	parts := strings.Split(projectPermissionID, "/")
	if len(parts) != 8 {
		return nil, fmt.Errorf("wrong number of items in ID (%d vs. 7) in %s", len(parts), projectPermissionID)
	}
	perms := &ProjectDefaultPermission{
		ProjectID:      parts[1],
		EntityID:       parts[4],
		EntityType:     parts[3],
		TargetType:     parts[5],
		CapabilityName: parts[6],
		CapabilityMode: parts[7],
	}
	entityTypes := []string{"groups", "users"}
	if !slices.Contains(entityTypes, perms.EntityType) {
		return nil, fmt.Errorf("unknown entity type (%s) not in: %s", perms.EntityType, strings.Join(entityTypes, ", "))
	}
	if !slices.Contains(projectDefaultPermissionTargetTypes, perms.TargetType) {
		return nil, fmt.Errorf("unknown target type (%s) not in: %s", perms.TargetType, strings.Join(projectDefaultPermissionTargetTypes, ", "))
	}
	return perms, nil
}
