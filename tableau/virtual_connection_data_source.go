package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualConnectionDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualConnectionDataSource{}
)

func VirtualConnectionDataSource() datasource.DataSource {
	return &virtualConnectionDataSource{}
}

type virtualConnectionDataSource struct {
	client *Client
}

type virtualConnectionDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	ProjectID types.String `tfsdk:"project_id"`
	OwnerID   types.String `tfsdk:"owner_id"`
	Content   types.String `tfsdk:"content"`
	Name      types.String `tfsdk:"name"`
}

func (d *virtualConnectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_connection"
}

func (d *virtualConnectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve virtual Connection details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the virtual Connection",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "Project ID of the virtual connection",
			},
			"owner_id": schema.StringAttribute{
				Computed:    true,
				Description: "Owner ID of the virtual connection",
			},
			"content": schema.StringAttribute{
				Computed:    true,
				Description: "Definition of the virtual connection as JSON",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the virtual connection",
			},
		},
	}
}

func (d *virtualConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualConnectionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	virtualConnection, err := d.client.GetVirtualConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Download Tableau Virtual Connection",
			err.Error(),
		)
		return
	}
	state.ProjectID = types.StringValue(virtualConnection.Project.ID)
	state.OwnerID = types.StringValue(virtualConnection.Owner.ID)
	state.Content = types.StringValue(virtualConnection.Content)
	state.Name = types.StringValue(virtualConnection.Name)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *virtualConnectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
