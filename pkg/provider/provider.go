package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gthesheep/terraform-provider-tableau/pkg/data_sources"
	"github.com/gthesheep/terraform-provider-tableau/pkg/resources"
	"github.com/gthesheep/terraform-provider-tableau/pkg/tableau"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TABLEAU_SERVER_URL", nil),
				Description: "URL of your Tableau server",
			},
			"server_version": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TABLEAU_SERVER_VERSION", nil),
				Description: "Version of the server identified in URL",
			},
			"username": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TABLEAU_USERNAME", nil),
				Description:   "Login Username",
				ConflictsWith: []string{"personal_access_token_name"},
			},
			"password": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TABLEAU_PASSWORD", nil),
				Description:   "Login Password",
				ConflictsWith: []string{"personal_access_token_secret"},
			},
			"personal_access_token_name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TABLEAU_PERSONAL_ACCESS_TOKEN_NAME", nil),
				Description:   "Personal access token name",
				ConflictsWith: []string{"username"},
			},
			"personal_access_token_secret": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("TABLEAU_PERSONAL_ACCESS_TOKEN_SECRET", nil),
				Description:   "Personal access token secret",
				ConflictsWith: []string{"password"},
			},
			"site": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TABLEAU_SITE_NAME", nil),
				Description: "Site name from your Tableau URL",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"tableau_group": data_sources.DatasourceGroup(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"tableau_user":       resources.ResourceUser(),
			"tableau_group":      resources.ResourceGroup(),
			"tableau_group_user": resources.ResourceGroupUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	server := d.Get("server_url").(string)
	serverVersion := d.Get("server_version").(string)
	site := d.Get("site").(string)

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	personalAccessTokenName := d.Get("personal_access_token_name").(string)
	personalAccessTokenSecret := d.Get("personal_access_token_secret").(string)

	var diags diag.Diagnostics

	if (server != "") && (site != "") && (serverVersion != "") {
		c, err := tableau.NewClient(
			&server,
			&username,
			&password,
			&personalAccessTokenName,
			&personalAccessTokenSecret,
			&site,
			&serverVersion,
		)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to login to Tableau",
				Detail:   err.Error(),
			})
			return nil, diags
		}

		return c, diags
	}

	c, err := tableau.NewClient(nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Tableau client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
