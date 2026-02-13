package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualConnectionRevisionDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualConnectionRevisionDataSource{}
)

func VirtualConnectionRevisionDataSource() datasource.DataSource {
	return &virtualConnectionRevisionDataSource{}
}

type virtualConnectionRevisionDataSource struct {
	client *Client
}

type virtualConnectionRevisionDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Revision  types.Int32  `tfsdk:"revision"`
	ProjectID types.String `tfsdk:"project_id"`
	OwnerID   types.String `tfsdk:"owner_id"`
	Content   types.String `tfsdk:"content"`
	Name      types.String `tfsdk:"name"`
}

func (d *virtualConnectionRevisionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_connection_revision"
}

func (d *virtualConnectionRevisionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve virtual connection revisions details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the virtual connections",
			},
			"revision": schema.Int32Attribute{
				Required:    true,
				Description: "Revision number of virtual connection",
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

func (d *virtualConnectionRevisionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualConnectionRevisionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	revision, err := d.client.GetVirtualConnectionRevision(state.ID.ValueString(), state.Revision.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Virtual Connection Revision",
			err.Error(),
		)
		return
	}
	state.ProjectID = types.StringValue(revision.Project.ID)
	state.OwnerID = types.StringValue(revision.Owner.ID)
	state.Content = types.StringValue(revision.Content)
	state.Name = types.StringValue(revision.Name)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *virtualConnectionRevisionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
