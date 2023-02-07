package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

func GroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	client *Client
}

type groupDataSourceModel struct {
	GroupID types.String `tfsdk:"group_id"`
	Name    types.String `tfsdk:"name"`
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve group details",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the group",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name for the group",
			},
		},
	}
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	group, err := d.client.GetGroup(state.GroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Group",
			err.Error(),
		)
		return
	}

	state.GroupID = types.StringValue(group.ID)
	state.Name = types.StringValue(group.Name)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
