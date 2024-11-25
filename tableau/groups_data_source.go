package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &groupsDataSource{}
	_ datasource.DataSourceWithConfigure = &groupsDataSource{}
)

func GroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

type groupsDataSource struct {
	client *Client
}

type groupsNestedDataModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	MinimumSiteRole types.String `tfsdk:"minimum_site_role"`
}

type groupsDataSourceModel struct {
	ID     types.String            `tfsdk:"id"`
	Groups []groupsNestedDataModel `tfsdk:"groups"`
}

func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve groups details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the list of groups",
			},
			"groups": schema.ListNestedAttribute{
				Description: "List of groups and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the group",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name for the groups",
						},
						"minimum_site_role": schema.StringAttribute{
							Computed:    true,
							Description: "Minimum site role for the groups",
						},
					},
				},
			},
		},
	}
}

func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	groups, err := d.client.GetGroups()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Groups",
			err.Error(),
		)
		return
	}

	for _, group := range groups {
		groupDataSourceModel := groupsNestedDataModel{
			ID:              types.StringValue(group.ID),
			Name:            types.StringValue(group.Name),
			MinimumSiteRole: types.StringValue(group.MinimumSiteRole),
		}
		state.Groups = append(state.Groups, groupDataSourceModel)
	}

	state.ID = types.StringValue("allGroups")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
