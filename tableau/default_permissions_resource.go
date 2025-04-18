package tableau

import (
	"context"
	"fmt"
	"maps"
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
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
										Validators: []validator.String{
											stringvalidator.OneOf([]string{
												"AddComment",
												"ChangeHierarchy",
												"ChangePermissions",
												"Connect",
												"CreateRefreshMetrics",
												"Delete",
												"Execute",
												"ExportData",
												"ExportImage",
												"ExportXml",
												"Filter",
												"PulseMetricDefine",
												"Read",
												"RunExplainData",
												"SaveAs",
												"ShareView",
												"ViewComments",
												"ViewUnderlyingData",
												"VizqlDataApiAccess",
												"WebAuthoring",
												"WebAuthoringForFlows",
												"Write",
											}...),
										},
									},
									"mode": schema.StringAttribute{
										Required:    true,
										Description: "Mode of the capability (Allow/Deny)",
										Validators: []validator.String{
											stringvalidator.OneOf([]string{
												"Allow",
												"Deny",
											}...),
										},
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

// granteeCapabilitiesModelMap turns `[]GranteeCapabilityModel` into map structure where
// 1st key is key combined from user and/or group UUID
// 2nd key is capability name
// and final value is mode of that capability
func granteeCapabilitiesModelMap(granteeCapabilityModels []GranteeCapabilityModel) map[string]map[string]string {
	mapping := map[string]map[string]string{}
	for _, granteeCapability := range granteeCapabilityModels {
		key := ""
		if granteeCapability.GroupID.ValueString() != "" {
			key = "|" + granteeCapability.GroupID.ValueString()
		} else if granteeCapability.UserID.ValueString() != "" {
			key = granteeCapability.UserID.ValueString() + "|"
		}
		mapping[key] = map[string]string{}
		for _, capability := range granteeCapability.Capabilities {
			mapping[key][capability.Name.ValueString()] = capability.Mode.ValueString()
		}
	}
	return mapping
}

// granteeCapabilitiesMap turns `[]GranteeCapability` into map structure where
// 1st key is key combined from user and/or group UUID
// 2nd key is capability name
// and final value is mode of that capability
func granteeCapabilitiesMap(granteeCapabilities []GranteeCapability) map[string]map[string]string {
	mapping := map[string]map[string]string{}
	for _, granteeCapability := range granteeCapabilities {
		key := ""
		if granteeCapability.Group != nil {
			key = "|" + granteeCapability.Group.ID
		} else if granteeCapability.User != nil {
			key = granteeCapability.User.ID + "|"
		}
		mapping[key] = map[string]string{}
		for _, capability := range granteeCapability.Capabilities.Capabilities {
			mapping[key][capability.Name] = capability.Mode
		}
	}
	return mapping
}

// mergeGrantees updates current state with the information that has been fetched from Tableau.
// This covers:
// - updating mode in current grantees capabilities
// - adding new grantees and/or grantee's capabilities
// - removing old grantees and/or grantee's capabilities
func mergeGrantees(state []GranteeCapabilityModel, readValues map[string]map[string]string) []GranteeCapabilityModel {
	grantees := len(state)
	gidx := 0
	for gidx < grantees {
		key := state[gidx].UserID.ValueString() + "|" + state[gidx].GroupID.ValueString()
		if _, ok := readValues[key]; !ok { // grantee has disappeared
			state = slices.Delete(state, gidx, gidx+1)
			grantees--
			continue
		}
		capabilities := len(state[gidx].Capabilities)
		cidx := 0
		for cidx < capabilities {
			name := state[gidx].Capabilities[cidx].Name.ValueString()
			if _, ok := readValues[key][name]; !ok { // capability in grantee has disappeared
				state[gidx].Capabilities = slices.Delete(state[gidx].Capabilities, cidx, cidx+1)
				capabilities--
				continue
			}
			mode := state[gidx].Capabilities[cidx].Mode.ValueString()
			if readValues[key][name] != mode {
				state[gidx].Capabilities[cidx].Mode = types.StringValue(readValues[key][name])
			}
			delete(readValues[key], name)
			cidx++
		}
		for name, mode := range readValues[key] { // capabilities that are missing from state
			state[gidx].Capabilities = append(state[gidx].Capabilities, CapabilityModel{
				Name: types.StringValue(name),
				Mode: types.StringValue(mode),
			})
		}
		delete(readValues, key)
		gidx++
	}
	for grantee, capabilities := range readValues { // grantees that are missing from state
		fields := strings.Split(grantee, "|")
		newGrantee := GranteeCapabilityModel{
			UserID:       types.StringValue(fields[0]),
			GroupID:      types.StringValue(fields[1]),
			Capabilities: []CapabilityModel{},
		}
		for name, mode := range capabilities { // capabilities for new grantee
			newGrantee.Capabilities = append(newGrantee.Capabilities, CapabilityModel{
				Name: types.StringValue(name),
				Mode: types.StringValue(mode),
			})
		}
		state = append(state, newGrantee)
	}
	return state
}

// Read is perhaps unnecessarily complex looking creature, but initial attempts to build something that would
// `[]GranteeCapabilityModel` and then simply assign it into `state.GranteeCapabilities` kept on making new
// `terraform plan`s which always had capabilities in wrong order or something and terraform would therefore
// say that those resources needed updating.
// Current solution reads relevant information from tableau, go through the data in same order as it is in state,
// modified if needed and adds all the new stuff at the end of capabilities/granteeCapabilities and deletes
// if it feels that something that is in state does not exist in Tableau
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
	newGranteeModel := granteeCapabilitiesMap(defaultPermissions.GranteeCapabilities)
	state.GranteeCapabilities = mergeGrantees(state.GranteeCapabilities, newGranteeModel)
	state.ID = types.StringValue(getDefaultPermissionID(permission.ProjectID, permission.TargetType))
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update operation is combination of PUT and one or more DELETE operations.
// If we modify permissions from Allow to Deny or reverse, add grantees and capabilities, we can handle it all with single
// PUT REST API call, but if we drop some capabilities from grantee or drop whole grantee, then we have to call DELETE
// method for each and every capability to clean them out.
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
			"Error updating default permission",
			"Could not update default permission, unexpected error: "+err.Error(),
		)
		return
	}
	current, err := r.client.GetDefaultPermissions(projectID, targetType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating default permission",
			"Could not fetch default permission, unexpected error: "+err.Error(),
		)
		return
	}
	// Section that involves DELETE calls
	currentGranteeModel := granteeCapabilitiesMap(current.GranteeCapabilities)
	planGranteeModel := granteeCapabilitiesModelMap(plan.GranteeCapabilities)
	planGranteeKeys := slices.Collect(maps.Keys(planGranteeModel))
	for currentGrantee, currentCapability := range currentGranteeModel {
		if slices.Contains(planGranteeKeys, currentGrantee) {
			planCapabilities := slices.Collect(maps.Keys(planGranteeModel[currentGrantee]))
			for currentCapabilityName, currentCapabilityValue := range currentCapability {
				if !slices.Contains(planCapabilities, currentCapabilityName) { // missing capability
					fields := strings.Split(currentGrantee, "|")
					err := r.client.DeleteDefaultPermission(
						&fields[0],
						&fields[1],
						projectID,
						targetType,
						currentCapabilityName,
						currentCapabilityValue,
					)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error Deleting Tableau Default Permissions",
							"Could not delete default permissions, unexpected error: "+err.Error(),
						)
						return
					}
				}
			}
		} else { // missing grantee
			for currentCapabilityName, currentCapabilityValue := range currentCapability {
				fields := strings.Split(currentGrantee, "|")
				err := r.client.DeleteDefaultPermission(
					&fields[0],
					&fields[1],
					projectID,
					targetType,
					currentCapabilityName,
					currentCapabilityValue,
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Deleting Tableau Default Permissions",
						"Could not delete default permissions, unexpected error: "+err.Error(),
					)
					return
				}
			}
		}
	}
	plan.ID = types.StringValue(getDefaultPermissionID(projectID, targetType))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes all grantees and their permission from given (project ID, target Type) combination
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
	for _, grantee := range state.GranteeCapabilities {
		for _, capability := range grantee.Capabilities {
			tflog.Info(ctx, "DeletingDefaultPermission", map[string]interface{}{
				"userID":         grantee.UserID.ValueString(),
				"groupID":        grantee.GroupID.ValueString(),
				"projectID":      permission.ProjectID,
				"targetType":     permission.TargetType,
				"capabilityName": capability.Name.ValueString(),
				"capabilityMode": capability.Mode.ValueString(),
			})
			err := r.client.DeleteDefaultPermission(
				grantee.UserID.ValueStringPointer(),
				grantee.GroupID.ValueStringPointer(),
				permission.ProjectID,
				permission.TargetType,
				capability.Name.ValueString(),
				capability.Mode.ValueString(),
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting Tableau default Permissions",
					"Could not delete default permissions, unexpected error: "+err.Error(),
				)
				return
			}
		}
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
