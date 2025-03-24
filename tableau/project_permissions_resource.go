package tableau

import (
	"context"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
										Validators: []validator.String{
											stringvalidator.OneOf([]string{
												"ProjectLeader",
												"Read",
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

	plan.ID = types.StringValue(projectID)

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
func (r *projectPermissionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectPermissionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := state.ID.ValueString()
	projectPermissions, err := r.client.GetProjectPermissions(projectID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	if projectPermissions == nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Project Permission",
			"GetProjectPermission returned nil with "+projectID,
		)
		return
	}
	newGranteeModel := granteeCapabilitiesMap(projectPermissions.GranteeCapabilities)
	state.GranteeCapabilities = mergeGrantees(state.GranteeCapabilities, newGranteeModel)
	state.ID = types.StringValue(projectID)
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
func (r *projectPermissionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	current, err := r.client.GetProjectPermissions(projectID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project permission",
			"Could not fetch project permission, unexpected error: "+err.Error(),
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
					err := r.client.DeleteProjectPermission(
						&fields[0],
						&fields[1],
						projectID,
						currentCapabilityName,
						currentCapabilityValue,
					)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error Deleting Tableau Project Permissions",
							"Could not delete project permissions, unexpected error: "+err.Error(),
						)
						return
					}
				}
			}
		} else { // missing grantee requires that we delete all capabilities one by one
			for currentCapabilityName, currentCapabilityValue := range currentCapability {
				fields := strings.Split(currentGrantee, "|")
				err := r.client.DeleteProjectPermission(
					&fields[0],
					&fields[1],
					projectID,
					currentCapabilityName,
					currentCapabilityValue,
				)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Deleting Tableau Project Permissions",
						"Could not delete project permissions, unexpected error: "+err.Error(),
					)
					return
				}
			}
		}
	}
	plan.ID = types.StringValue(projectID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete all grantees and their capabilities from project
func (r *projectPermissionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectPermissionsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ID.ValueString()
	for _, grantee := range state.GranteeCapabilities {
		for _, capability := range grantee.Capabilities {
			err := r.client.DeleteProjectPermission(
				grantee.UserID.ValueStringPointer(),
				grantee.GroupID.ValueStringPointer(),
				projectID,
				capability.Name.ValueString(),
				capability.Mode.ValueString(),
			)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting Tableau project Permissions",
					"Could not delete project permissions, unexpected error: "+err.Error(),
				)
				return
			}
		}
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
