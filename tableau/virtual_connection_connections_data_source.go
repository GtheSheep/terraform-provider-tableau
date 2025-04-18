package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualConnectionConnectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualConnectionConnectionsDataSource{}
)

func VirtualConnectionConnectionsDataSource() datasource.DataSource {
	return &virtualConnectionConnectionsDataSource{}
}

type virtualConnectionConnectionsDataSource struct {
	client *Client
}

type virtualConnectionConnectionsNestedDataModel struct {
	ID            types.String `tfsdk:"id"`
	DBClass       types.String `tfsdk:"db_class"`
	ServerAddress types.String `tfsdk:"server_address"`
	ServerPort    types.String `tfsdk:"server_port"`
	UserName      types.String `tfsdk:"username"`
}
type virtualConnectionConnectionsDataSourceModel struct {
	ID          types.String                                  `tfsdk:"id"`
	Connections []virtualConnectionConnectionsNestedDataModel `tfsdk:"connections"`
}

func (d *virtualConnectionConnectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_connection_connections"
}

func (d *virtualConnectionConnectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve virtual connection's connections details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the virtual connections",
			},
			"connections": schema.ListNestedAttribute{
				Description: "List database connections of virtual connection and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the connection in Virtual Connection",
						},
						"db_class": schema.StringAttribute{
							Computed:    true,
							Description: "DB class",
						},
						"server_address": schema.StringAttribute{
							Computed:    true,
							Description: "Server address",
						},
						"server_port": schema.StringAttribute{
							Computed:    true,
							Description: "Server port",
						},
						"username": schema.StringAttribute{
							Computed:    true,
							Description: "Username",
						},
					},
				},
			},
		},
	}
}

func (d *virtualConnectionConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualConnectionConnectionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	connections, err := d.client.GetVirtualConnectionConnections(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Database Connections of Tableau Virtual Connection",
			err.Error(),
		)
		return
	}

	for _, connection := range connections {
		virtualConnectionConnection := virtualConnectionConnectionsNestedDataModel{
			ID:            types.StringValue(connection.ID),
			DBClass:       types.StringValue(connection.DBClass),
			ServerAddress: types.StringValue(connection.ServerAddress),
			ServerPort:    types.StringValue(connection.ServerPort),
			UserName:      types.StringValue(connection.UserName),
		}
		state.Connections = append(state.Connections, virtualConnectionConnection)
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *virtualConnectionConnectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
