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
	_ resource.Resource                = &projectPermissionsResource{}
	_ resource.ResourceWithConfigure   = &projectPermissionsResource{}
	_ resource.ResourceWithImportState = &projectPermissionsResource{}
)

func NewProjectPermissionsResource() resource.Resource {
	return &projectPermissionsResource{}
}

type projectPermissionsResource struct {
	client *Client
}

type projectPermissionsResourceModel struct {
	ID                  types.String             `tfsdk:"id"`
	GranteeCapabilities []GranteeCapabilityModel `tfsdk:"grantee_capabilities"`
}

/*
Examples of existing capability names
- Read
- Write
*/

func (r *projectPermissionsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_permissions"
}

func (r *projectPermissionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the project",
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

func (r *projectPermissionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectPermissionsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := plan.ID.ValueString()
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
	_, err := r.client.CreateProjectPermissions(projectID, GranteeCapabilities{GranteeCapabilities: granteeCapabilities})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating default permission",
			"Could not create default permission, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(getProjectPermissionsID(projectID))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectPermissionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	permission, err := getProjectPermissionsFromID(state.ID.ValueString())
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
			fmt.Sprintf("getProjectPermissionFromID(%s) returned nil", state.ID.ValueString()),
		)
		return
	}
	state.ID = types.StringValue(permission.ID)
	projectPermissions, err := r.client.GetProjectPermissions(permission.ID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if projectPermissions == nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Default Permission",
			fmt.Sprintf("GetProjectPermission returned nil with %#v", permission),
		)
		return
	}
	state.GranteeCapabilities = []GranteeCapabilityModel{}
	for _, granteeCapabilities := range projectPermissions.GranteeCapabilities {
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
	state.ID = types.StringValue(getProjectPermissionsID(permission.ID))
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectPermissionsResourceModel
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

func (r *projectPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectPermissionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := getProjectPermissionsFromID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Default Permissions",
			err.Error(),
		)
		return
	}
	if err = r.client.DeleteProjectPermissions(permission.ID); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Default Permissions",
			"Could not delete default permissions, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *projectPermissionsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *projectPermissionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getProjectPermissionsID(projectID string) string {
	return fmt.Sprintf("projects/%s/permissions", projectID)
}

func getProjectPermissionsFromID(id string) (*ProjectPermissions, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("wrong number of items in ID (%d vs. 3) in %s", len(parts), id)
	}
	perms := &ProjectPermissions{ID: parts[1]}
	return perms, nil
}
