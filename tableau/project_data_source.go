package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

func ProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *Client
}

type projectDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	ContentPermissions types.String `tfsdk:"content_permissions"`
	ParentProjectID    types.String `tfsdk:"parent_project_id"`
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve project details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the project",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name for the project",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description for the project",
			},
			"content_permissions": schema.StringAttribute{
				Computed:    true,
				Description: "Permissions for the project content - ManagedByOwner is the default",
			},
			"parent_project_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for the parent project",
			},
		},
	}
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state projectDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	project, err := d.client.GetProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Project",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(project.ID)
	state.Name = types.StringValue(project.Name)
	state.Description = types.StringValue(project.Description)
	state.ContentPermissions = types.StringValue(project.ContentPermissions)
	state.ParentProjectID = types.StringValue(project.ParentProjectID)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
