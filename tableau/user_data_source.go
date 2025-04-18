package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func UserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *Client
}

type userDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	FullName    types.String `tfsdk:"full_name"`
	SiteRole    types.String `tfsdk:"site_role"`
	AuthSetting types.String `tfsdk:"auth_setting"`
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve user details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the user",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name for the user",
			},
			"full_name": schema.StringAttribute{
				Computed:    true,
				Description: "Full name for user",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "User email",
			},
			"site_role": schema.StringAttribute{
				Computed:    true,
				Description: "Site role for the user",
			},
			"auth_setting": schema.StringAttribute{
				Computed:    true,
				Description: "Auth setting for the user",
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	user, err := d.client.GetUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau User",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(user.ID)
	state.Name = types.StringValue(user.Name)
	state.Email = types.StringValue(user.Email)
	state.FullName = types.StringValue(user.FullName)
	state.SiteRole = types.StringValue(user.SiteRole)
	state.AuthSetting = types.StringValue(user.AuthSetting)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
