package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &datasourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &datasourcesDataSource{}
)

func DatasourcesDataSource() datasource.DataSource {
	return &datasourcesDataSource{}
}

type datasourcesDataSource struct {
	client *Client
}

type datasourcesNestedDataModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	OwnerID   types.String `tfsdk:"owner_id"`
	ProjectID types.String `tfsdk:"project_id"`
}

type datasourcesDataSourceModel struct {
	ID          types.String                 `tfsdk:"id"`
	Datasources []datasourcesNestedDataModel `tfsdk:"datasources"`
}

func (d *datasourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datasources"
}

func (d *datasourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve datasource details as a list of datasources available to read",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the project",
			},
			"datasources": schema.ListNestedAttribute{
				Description: "List of datasources and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the datasource",
						},
						"name": schema.StringAttribute{
							Description: "Datasource name",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of datasource",
						},
						"owner_id": schema.StringAttribute{
							Computed:    true,
							Description: "Datasource Owner ID",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "Datasource Project ID",
						},
					},
				},
			},
		},
	}
}

func (d *datasourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state datasourcesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	datasources, err := d.client.GetDatasources()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Datasource",
			err.Error(),
		)
		return
	}

	for _, datasource := range datasources {
		datasourceDataSourceModel := datasourcesNestedDataModel{
			ID:        types.StringValue(datasource.ID),
			Name:      types.StringValue(datasource.Name),
			Type:      types.StringValue(datasource.Type),
			OwnerID:   types.StringValue(datasource.Owner.ID),
			ProjectID: types.StringValue(datasource.Project.ID),
		}
		state.Datasources = append(state.Datasources, datasourceDataSourceModel)
	}

	state.ID = types.StringValue("allDatasources")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *datasourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
