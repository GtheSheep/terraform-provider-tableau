package tableau

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &tableauProvider{}
)

func New() provider.Provider {
	return &tableauProvider{}
}

type tableauProvider struct{}

func (p *tableauProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tableau"
}

func (p *tableauProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Tableau",
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Optional:    true,
				Description: "URL of your Tableau server - TABLEAU_SERVER_URL env var",
			},
			"server_version": schema.StringAttribute{
				Optional:    true,
				Description: "Version of the server identified in URL - TABLEAU_SERVER_VERSION env var",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Login Username - TABLEAU_USERNAME env var",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Login Password - TABLEAU_PASSWORD env var",
			},
			"personal_access_token_name": schema.StringAttribute{
				Optional:    true,
				Description: "Personal access token name - TABLEAU_PERSONAL_ACCESS_TOKEN_NAME env var",
			},
			"personal_access_token_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Personal access token secret - TABLEAU_PERSONAL_ACCESS_TOKEN_SECRET env var",
			},
			"site": schema.StringAttribute{
				Optional:    true,
				Description: "Site name from your Tableau URL - TABLEAU_SITE_NAME env var - for Tableau Server default sites leave as ''",
			},
			"is_tcm": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true if using this provider for Tableau Cloud Manager - TABLEAU_IS_TCM env var - default false",
			},
		},
	}
}

type tableauProviderModel struct {
	ServerURL                 types.String `tfsdk:"server_url"`
	ServerVersion             types.String `tfsdk:"server_version"`
	Username                  types.String `tfsdk:"username"`
	Password                  types.String `tfsdk:"password"`
	PersonalAccessTokenName   types.String `tfsdk:"personal_access_token_name"`
	PersonalAccessTokenSecret types.String `tfsdk:"personal_access_token_secret"`
	Site                      types.String `tfsdk:"site"`
	IsTCM                     types.Bool   `tfsdk:"is_tcm"`
}

func (p *tableauProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Tableau client")

	var config tableauProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ServerURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Unknown Tableau Server URL",
			"Tableau Server URL must be provided in order to establish a connection",
		)
	}

	if config.ServerVersion.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_version"),
			"Unknown Tableau Server version",
			"Tableau Server Version must be provided in order to establish a connection, currently no default is set",
		)
	}

	if config.Username.IsUnknown() && config.PersonalAccessTokenName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Tableau Username",
			"Tableau Username or Personal Access Token Name must be provided in order to establish a connection",
		)
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_name"),
			"Unknown Tableau Personal Access Token Name",
			"Tableau Username or Personal Access Token Name must be provided in order to establish a connection",
		)
	}

	if config.Password.IsUnknown() && config.PersonalAccessTokenSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Tableau Password",
			"Tableau Password or Personal Access Token Secret must be provided in order to establish a connection",
		)
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_secret"),
			"Unknown Tableau Personal Access Token Secret",
			"Tableau Password or Personal Access Token Secret must be provided in order to establish a connection",
		)
	}

	if config.Site.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("site"),
			"Unknown Tableau Site",
			"Tableau Site must be provided in order to establish a connection",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	serverURL := os.Getenv("TABLEAU_SERVER_URL")
	serverVersion := os.Getenv("TABLEAU_SERVER_VERSION")
	username := os.Getenv("TABLEAU_USERNAME")
	password := os.Getenv("TABLEAU_PASSWORD")
	personalAccessTokenName := os.Getenv("TABLEAU_PERSONAL_ACCESS_TOKEN_NAME")
	personalAccessTokenSecret := os.Getenv("TABLEAU_PERSONAL_ACCESS_TOKEN_SECRET")
	site := os.Getenv("TABLEAU_SITE_NAME")
	isTCM := false
	isTCMString := os.Getenv("TABLEAU_IS_TCM")
	isTCM, err := strconv.ParseBool(isTCMString)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse is_tcm environment variable to boolean",
			"Tableau Client Error: "+err.Error(),
		)
		return
	}

	if !config.ServerURL.IsNull() {
		serverURL = config.ServerURL.ValueString()
	}

	if !config.ServerVersion.IsNull() {
		serverVersion = config.ServerVersion.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.PersonalAccessTokenName.IsNull() {
		personalAccessTokenName = config.PersonalAccessTokenName.ValueString()
	}

	if !config.PersonalAccessTokenSecret.IsNull() {
		personalAccessTokenSecret = config.PersonalAccessTokenSecret.ValueString()
	}

	if !config.Site.IsNull() {
		site = config.Site.ValueString()
	}

	if !config.IsTCM.IsNull() {
		if !config.IsTCM.ValueBool() {
			isTCM = config.IsTCM.ValueBool()
		}
	}

	if serverURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Missing Tableau Server URL",
			"Tableau Server URL must be provided in order to establish a connection",
		)
	}

	if serverVersion == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_version"),
			"Missing Tableau Server version",
			"Tableau Server Version must be provided in order to establish a connection, currently no default is set",
		)
	}

	if username == "" && personalAccessTokenName == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Tableau Username",
			"Tableau Username or Personal Access Token Name must be provided in order to establish a connection",
		)
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_name"),
			"Missing Tableau Personal Access Token Name",
			"Tableau Username or Personal Access Token Name must be provided in order to establish a connection",
		)
	}

	if password == "" && personalAccessTokenSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Tableau Password",
			"Tableau Password or Personal Access Token Secret must be provided in order to establish a connection",
		)
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_access_token_secret"),
			"Missing Tableau Personal Access Token Secret",
			"Tableau Password or Personal Access Token Secret must be provided in order to establish a connection",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := NewClient(
		&serverURL,
		&username,
		&password,
		&personalAccessTokenName,
		&personalAccessTokenSecret,
		&site,
		&serverVersion,
		isTCM,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Tableau API Client",
			"An unexpected error occurred when creating the Tableau API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Tableau Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Tableau client", map[string]any{"success": true})
}

func (p *tableauProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		GroupDataSource,
		GroupsDataSource,
		UserDataSource,
		UsersDataSource,
		ProjectDataSource,
		ProjectsDataSource,
		SiteDataSource,
		DatasourceDataSource,
		DatasourcesDataSource,
	}
}

func (p *tableauProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewGroupResource,
		NewGroupUserResource,
		NewProjectResource,
		NewSiteResource,
		NewDatasourcePermissionResource,
		NewProjectPermissionResource,
		NewViewPermissionResource,
		NewVirtualConnectionPermissionResource,
		NewWorkbookPermissionResource,
	}
}
