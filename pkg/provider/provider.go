package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TABLEAU_USERNAME", nil),
				Description: "Login Username",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TABLEAU_PASSWORD", nil),
				Description: "Login Password",
			},
			"site": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TABLEAU_SITE_NAME", nil),
				Description: "Site name from your Tableau URL",
			},
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ResourcesMap:         map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	server := d.Get("server").(string)
	username := d.Get("username").(string)
	password := d.Get("server").(string)
	site := d.Get("site").(string)
	server_version := d.Get("server_version").(string)

	var diags diag.Diagnostics

	if (server != "") && (username != "") && (site != "") && (server_version != "") {
		c, err := tableau.NewClient(&server, &username, &password, &site, &server_version)

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

	c, err := tableau.NewClient(nil, nil, nil, nil, nil)
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
