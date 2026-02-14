package tableau

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &defaultPermissionsDataSource{}
	_ datasource.DataSourceWithConfigure = &defaultPermissionsDataSource{}
)

func DefaultPermissionsDataSource() datasource.DataSource {
	return &defaultPermissionsDataSource{}
}

type defaultPermissionsDataSource struct {
	client *Client
}

type CapabilityModel struct {
	Name types.String `tfsdk:"name"`
	Mode types.String `tfsdk:"mode"`
}

type GranteeCapabilityModel struct {
	UserID       types.String      `tfsdk:"user_id"`
	GroupID      types.String      `tfsdk:"group_id"`
	Capabilities []CapabilityModel `tfsdk:"capabilities"`
}

type defaultPermissionsDataSourceModel struct {
	ProjectID           types.String             `tfsdk:"project_id"`
	TargetType          types.String             `tfsdk:"target_type"`
	GranteeCapabilities []GranteeCapabilityModel `tfsdk:"grantee_capabilities"`
}

func (d *defaultPermissionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_permissions"
}

func (d *defaultPermissionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve project details",
		Attributes: map[string]schema.Attribute{
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
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the group",
						},
						"user_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the user",
						},
						"capabilities": schema.ListNestedAttribute{
							Description: "List of grantee capabilities for users and groups",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "Name of the capability",
									},
									"mode": schema.StringAttribute{
										Computed:    true,
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

func (d *defaultPermissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state defaultPermissionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	perms, err := d.client.GetDefaultPermissions(state.ProjectID.ValueString(), state.TargetType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Project Permissions",
			err.Error(),
		)
		return
	}

	for _, granteeCapability := range perms.GranteeCapabilities {
		newGranteeCapability := &GranteeCapabilityModel{}
		if granteeCapability.Group != nil {
			newGranteeCapability.GroupID = types.StringValue(granteeCapability.Group.ID)
		}
		if granteeCapability.User != nil {
			newGranteeCapability.UserID = types.StringValue(granteeCapability.User.ID)
		}
		newCapabilities := []CapabilityModel{}
		for _, capability := range granteeCapability.Capabilities.Capabilities {
			newCapabilities = append(newCapabilities, CapabilityModel{
				Name: types.StringValue(capability.Name),
				Mode: types.StringValue(capability.Mode),
			})
		}
		state.GranteeCapabilities = append(state.GranteeCapabilities, GranteeCapabilityModel{
			GroupID:      newGranteeCapability.GroupID,
			UserID:       newGranteeCapability.UserID,
			Capabilities: newCapabilities,
		})
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *defaultPermissionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
