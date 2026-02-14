package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualConnectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualConnectionsDataSource{}
)

func VirtualConnectionsDataSource() datasource.DataSource {
	return &virtualConnectionsDataSource{}
}

type virtualConnectionsDataSource struct {
	client *Client
}

type virtualConnectionsNestedDataModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	HasExtracts types.Bool   `tfsdk:"has_extracts"`
	IsCertified types.Bool   `tfsdk:"is_certified"`
	WebPageURL  types.String `tfsdk:"web_page_url"`
}

type virtualConnectionsDataSourceModel struct {
	ID                 types.String                        `tfsdk:"id"`
	VirtualConnections []virtualConnectionsNestedDataModel `tfsdk:"virtual_connections"`
}

func (d *virtualConnectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_connections"
}

func (d *virtualConnectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve virtual connections details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the virtual connections",
			},
			"virtual_connections": schema.ListNestedAttribute{
				Description: "List of virtual connections and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the virtual connection",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name for the virtual connection",
						},
						"has_extracts": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether or not this virtual connection has extracts",
						},
						"is_certified": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether or not this virtual connection is certified",
						},
						"web_page_url": schema.StringAttribute{
							Computed:    true,
							Description: "Web page URL for the virtual connection",
						},
					},
				},
			},
		},
	}
}

func (d *virtualConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualConnectionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	virtualConnections, err := d.client.GetVirtualConnections()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Virtual Connection",
			err.Error(),
		)
		return
	}

	for _, virtualConnection := range virtualConnections {
		virtualConnectionsDataModel := virtualConnectionsNestedDataModel{
			ID:          types.StringValue(virtualConnection.ID),
			Name:        types.StringValue(virtualConnection.Name),
			HasExtracts: types.BoolValue(virtualConnection.HasExtracts),
			IsCertified: types.BoolValue(virtualConnection.IsCertified),
			WebPageURL:  types.StringValue(virtualConnection.WebPageURL),
		}
		state.VirtualConnections = append(state.VirtualConnections, virtualConnectionsDataModel)
	}

	state.ID = types.StringValue("allVirtualConnections")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *virtualConnectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
