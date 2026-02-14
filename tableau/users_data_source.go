package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func UsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *Client
}

type usersNestedDataModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Email    types.String `tfsdk:"email"`
	SiteRole types.String `tfsdk:"site_role"`
}

type usersDataSourceModel struct {
	ID    types.String           `tfsdk:"id"`
	Users []usersNestedDataModel `tfsdk:"users"`
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve user details as a list of users available to read",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the users",
			},
			"users": schema.ListNestedAttribute{
				Description: "List of users and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the user",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name for the user",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "User email",
						},
						"site_role": schema.StringAttribute{
							Computed:    true,
							Description: "Site role for the user",
						},
					},
				},
			},
		},
	}
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	users, err := d.client.GetUsers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Users",
			err.Error(),
		)
		return
	}

	for _, user := range users {
		userDataSourceModel := usersNestedDataModel{
			ID:       types.StringValue(user.ID),
			Name:     types.StringValue(user.Name),
			Email:    types.StringValue(user.Email),
			SiteRole: types.StringValue(user.SiteRole),
		}
		state.Users = append(state.Users, userDataSourceModel)
	}

	state.ID = types.StringValue("allUsers")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
