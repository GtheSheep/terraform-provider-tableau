package tableau

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &defaultPermissionsResource{}
	_ resource.ResourceWithConfigure   = &defaultPermissionsResource{}
	_ resource.ResourceWithImportState = &defaultPermissionsResource{}
)

func NewDefaultPermissionsResource() resource.Resource {
	return &defaultPermissionsResource{}
}

type defaultPermissionsResource struct {
	client *Client
}

type defaultPermissionsResourceModel struct {
	ID                  types.String             `tfsdk:"id"`
	ProjectID           types.String             `tfsdk:"project_id"`
	TargetType          types.String             `tfsdk:"target_type"`
	GranteeCapabilities []GranteeCapabilityModel `tfsdk:"grantee_capabilities"`
}

/*
Examples of existing capability names
- AddComment
- ChangeHierarchy
- ChangePermissions
- CreateRefreshMetrics
- Delete
- ExportData
- ExportImage"
- ExportXml
- Filter
- Read
- RunExplainData
- ShareView
- ViewComments
- ViewUnderlyingData
- WebAuthoring
- Write
*/

func (r *defaultPermissionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_permissions"
}

func (r *defaultPermissionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Combination of project_id and target type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the project",
			},
			"target_type": schema.StringAttribute{
				Required:    true,
				Description: "Permissions for: " + strings.Join(defaultPermissionTargetTypes, ","),
			},
			"grantee_capabilities": schema.ListNestedAttribute{
				Description: "List of grantee capabilities for users and groups",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of the group",
						},
						"user_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of the user",
						},
						"capabilities": schema.ListNestedAttribute{
							Description: "List of grantee capabilities for users and groups",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:    true,
										Description: "Name of the capability",
									},
									"mode": schema.StringAttribute{
										Required:    true,
										Description: "Mode of the capability (Allow/Deny)",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *defaultPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan defaultPermissionsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	targetType := plan.TargetType.ValueString()
	granteeCapabilities := []GranteeCapability{}
	for _, granteeCapability := range plan.GranteeCapabilities {
		newGranteeCapability := GranteeCapability{Capabilities: Capabilities{}}
		if groupID := granteeCapability.GroupID.ValueString(); groupID != "" {
			newGranteeCapability.Group = &Group{ID: groupID}
		}
		if userID := granteeCapability.UserID.ValueString(); userID != "" {
			newGranteeCapability.User = &User{ID: userID}
		}
		newCapabilities := []Capability{}
		for _, capability := range granteeCapability.Capabilities {
			newCapabilities = append(newCapabilities, Capability{
				Name: capability.Name.ValueString(),
				Mode: capability.Mode.ValueString(),
			})
		}
		newGranteeCapability.Capabilities.Capabilities = newCapabilities
		granteeCapabilities = append(granteeCapabilities, newGranteeCapability)
	}
	_, err := r.client.CreateDefaultPermissions(projectID, targetType, granteeCapabilities)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating default permission",
			"Could not create default permission, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(getDefaultPermissionID(projectID, targetType))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state defaultPermissionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	permission, err := getDefaultPermissionFromID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in building default permission ID",
			err.Error(),
		)
		return
	}
	if permission == nil {
		resp.Diagnostics.AddError(
			"Error in building default permission ID",
			fmt.Sprintf("getDefaultPermissionFromID(%s) returned nil", state.ID.ValueString()),
		)
		return
	}
	state.ProjectID = types.StringValue(permission.ProjectID)
	state.TargetType = types.StringValue(permission.TargetType)
	defaultPermissions, err := r.client.GetDefaultPermissions(permission.ProjectID, permission.TargetType)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if defaultPermissions == nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Default Permission",
			fmt.Sprintf("GetDefaultPermission returned nil with %#v", permission),
		)
		return
	}
	state.GranteeCapabilities = []GranteeCapabilityModel{}
	for _, granteeCapabilities := range defaultPermissions.GranteeCapabilities {
		newGranteeCapability := GranteeCapabilityModel{}
		if granteeCapabilities.Group != nil {
			newGranteeCapability.GroupID = types.StringValue(granteeCapabilities.Group.ID)
		}
		if granteeCapabilities.User != nil {
			newGranteeCapability.UserID = types.StringValue(granteeCapabilities.User.ID)
		}
		newCapabilities := []CapabilityModel{}
		for _, capabilities := range granteeCapabilities.Capabilities.Capabilities {
			newCapabilities = append(newCapabilities, CapabilityModel{
				Name: types.StringValue(capabilities.Name),
				Mode: types.StringValue(capabilities.Mode),
			})
		}
		newGranteeCapability.Capabilities = newCapabilities
		state.GranteeCapabilities = append(state.GranteeCapabilities, newGranteeCapability)
	}
	state.ID = types.StringValue(getDefaultPermissionID(permission.ProjectID, permission.TargetType))
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan defaultPermissionsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	targetType := plan.TargetType.ValueString()
	granteeCapabilities := []GranteeCapability{}
	for _, granteeCapability := range plan.GranteeCapabilities {
		newGranteeCapability := GranteeCapability{Capabilities: Capabilities{}}
		if groupID := granteeCapability.GroupID.ValueString(); groupID != "" {
			newGranteeCapability.Group = &Group{ID: groupID}
		}
		if userID := granteeCapability.UserID.ValueString(); userID != "" {
			newGranteeCapability.User = &User{ID: userID}
		}
		newCapabilities := []Capability{}
		for _, capability := range granteeCapability.Capabilities {
			newCapabilities = append(newCapabilities, Capability{
				Name: capability.Name.ValueString(),
				Mode: capability.Mode.ValueString(),
			})
		}
		newGranteeCapability.Capabilities.Capabilities = newCapabilities
		granteeCapabilities = append(granteeCapabilities, newGranteeCapability)
	}
	_, err := r.client.CreateDefaultPermissions(projectID, targetType, granteeCapabilities)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating default permission",
			"Could not create default permission, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(getDefaultPermissionID(projectID, targetType))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state defaultPermissionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := getDefaultPermissionFromID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Default Permissions",
			err.Error(),
		)
		return
	}
	if err = r.client.DeleteDefaultPermissions(permission.ProjectID, permission.TargetType); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Default Permissions",
			"Could not delete default permissions, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *defaultPermissionsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *defaultPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getDefaultPermissionID(projectID, targetType string) string {
	return fmt.Sprintf("projects/%s/default-permissions/%s", projectID, targetType)
}

func getDefaultPermissionFromID(id string) (*DefaultPermissions, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("wrong number of items in ID (%d vs. 4) in %s", len(parts), id)
	}
	perms := &DefaultPermissions{
		ProjectID:  parts[1],
		TargetType: parts[3],
	}
	return perms, nil
}
