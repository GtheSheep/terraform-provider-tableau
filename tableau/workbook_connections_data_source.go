package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &workbookConnectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &workbookConnectionsDataSource{}
)

func WorkbookConnectionsDataSource() datasource.DataSource {
	return &workbookConnectionsDataSource{}
}

type workbookConnectionsDataSource struct {
	client *Client
}

type workbookConnectionNestedDataModel struct {
	ID                      types.String `tfsdk:"id"`
	Type                    types.String `tfsdk:"type"`
	DatasourceID            types.String `tfsdk:"datasource_id"`
	ServerAddress           types.String `tfsdk:"server_address"`
	ServerPort              types.String `tfsdk:"server_port"`
	UserName                types.String `tfsdk:"username"`
	EmbedPassword           types.Bool   `tfsdk:"embed_password"`
	QueryTaggingEnabled     types.Bool   `tfsdk:"query_tagging_enabled"`
	AuthenticationType      types.String `tfsdk:"authentication_type"`
	UseOAuthManagedKeychain types.Bool   `tfsdk:"use_oauth_managed_keychain"`
}

type workbookConnectionsDataSourceModel struct {
	ID          types.String                        `tfsdk:"id"`
	Connections []workbookConnectionNestedDataModel `tfsdk:"connections"`
}

func (d *workbookConnectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workbook_connections"
}

func (d *workbookConnectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve workbook connections details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the workbook",
			},
			"connections": schema.ListNestedAttribute{
				Description: "List workbook connections and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the connection in Virtual Connection",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Database connection type",
						},
						"datasource_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of datasource",
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
						"embed_password": schema.BoolAttribute{
							Computed:    true,
							Description: "Embed database password into connection",
						},
						"query_tagging_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Query tagging enabled",
						},
						"authentication_type": schema.StringAttribute{
							Computed:    true,
							Description: "Authentication type",
						},
						"use_oauth_managed_keychain": schema.BoolAttribute{
							Computed:    true,
							Description: "Use OAuth managed keychain",
						},
					},
				},
			},
		},
	}
}

func (d *workbookConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state workbookConnectionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	connections, err := d.client.GetWorkbookConnections(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Workbook Connections",
			err.Error(),
		)
		return
	}
	for _, connection := range connections {
		workbookConnection := workbookConnectionNestedDataModel{
			ID:                      types.StringValue(connection.ID),
			Type:                    types.StringValue(connection.Type),
			DatasourceID:            types.StringValue(connection.DataSourceID.ID),
			ServerAddress:           types.StringValue(connection.ServerAddress),
			ServerPort:              types.StringValue(connection.ServerPort),
			UserName:                types.StringValue(connection.UserName),
			EmbedPassword:           types.BoolValue(connection.EmbedPassword),
			QueryTaggingEnabled:     types.BoolValue(connection.QueryTaggingEnabled),
			AuthenticationType:      types.StringValue(connection.AuthenticationType),
			UseOAuthManagedKeychain: types.BoolValue(connection.UseOAuthManagedKeychain),
		}
		state.Connections = append(state.Connections, workbookConnection)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *workbookConnectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
