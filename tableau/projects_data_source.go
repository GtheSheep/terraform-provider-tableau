package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &projectsDataSource{}
	_ datasource.DataSourceWithConfigure = &projectsDataSource{}
)

func ProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

type projectsDataSource struct {
	client *Client
}

type projectsNestedDataModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	ParentProjectID types.String `tfsdk:"parent_project_id"`
}

type projectsDataSourceModel struct {
	ID       types.String              `tfsdk:"id"`
	Projects []projectsNestedDataModel `tfsdk:"projects"`
}

func (d *projectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *projectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve project details as a list of projects available to read",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the list of projects",
			},
			"projects": schema.ListNestedAttribute{
				Description: "List of projects and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Project ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Project name",
							Computed:    true,
						},
						"parent_project_id": schema.StringAttribute{
							Description: "Parent project ID",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state projectsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	projects, err := d.client.GetProjects()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Projects",
			err.Error(),
		)
		return
	}

	for _, project := range projects {
		projectDataSourceModel := projectsNestedDataModel{
			ID:              types.StringValue(project.ID),
			Name:            types.StringValue(project.Name),
			ParentProjectID: types.StringValue(project.ParentProjectID),
		}
		state.Projects = append(state.Projects, projectDataSourceModel)
	}

	state.ID = types.StringValue("allProjects")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
