package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &datasourceDataSource{}
	_ datasource.DataSourceWithConfigure = &datasourceDataSource{}
)

func DatasourceDataSource() datasource.DataSource {
	return &datasourceDataSource{}
}

type datasourceDataSource struct {
	client *Client
}

type datasourceDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	CertificationNote   types.String `tfsdk:"certification_note"`
	Type                types.String `tfsdk:"type"`
	ContentURL          types.String `tfsdk:"content_url"`
	EncryptExtracts     types.String `tfsdk:"encrypt_extracts"`
	HasExtracts         types.Bool   `tfsdk:"has_extracts"`
	IsCertified         types.Bool   `tfsdk:"is_certified"`
	UseRemoteQueryAgent types.Bool   `tfsdk:"use_remote_query_agent"`
	WebPageURL          types.String `tfsdk:"web_page_url"`
	OwnerID             types.String `tfsdk:"owner_id"`
	ProjectID           types.String `tfsdk:"project_id"`
	Tags                types.List   `tfsdk:"tags"`
}

func (d *datasourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datasource"
}

func (d *datasourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve datasource details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the datasource",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name for the datasource",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Datasource description",
			},
			"certification_note": schema.StringAttribute{
				Computed:    true,
				Description: "Certification note",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of datasource",
			},
			"content_url": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the datasource content",
			},
			"encrypt_extracts": schema.StringAttribute{
				Computed:    true,
				Description: "Whether or not this datasource encrypts extracts",
			},
			"has_extracts": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether or not this datasource has extracts",
			},
			"is_certified": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether or not this datasource is certified",
			},
			"use_remote_query_agent": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether or not this datasource uses a remote query agent",
			},
			"web_page_url": schema.StringAttribute{
				Computed:    true,
				Description: "Web page URL for the datasource",
			},
			"owner_id": schema.StringAttribute{
				Computed:    true,
				Description: "Datasource Owner ID",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "Datasource Project ID",
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of tags on the datasource",
			},
		},
	}
}

func (d *datasourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state datasourceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	datasource, err := d.client.GetDatasource(state.ID.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Datasource",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(datasource.ID)
	state.Name = types.StringValue(datasource.Name)
	state.Description = types.StringValue(datasource.Description)
	state.CertificationNote = types.StringValue(datasource.CertificationNote)
	state.Type = types.StringValue(datasource.Type)
	state.ContentURL = types.StringValue(datasource.ContentURL)
	state.EncryptExtracts = types.StringValue(datasource.EncryptExtracts)
	state.HasExtracts = types.BoolValue(datasource.HasExtracts)
	state.IsCertified = types.BoolValue(datasource.IsCertified)
	state.UseRemoteQueryAgent = types.BoolValue(datasource.UseRemoteQueryAgent)
	state.WebPageURL = types.StringValue(datasource.WebPageURL)
	state.OwnerID = types.StringValue(datasource.Owner.ID)
	state.ProjectID = types.StringValue(datasource.Project.ID)

	tags := make([]attr.Value, 0, len(datasource.Tags.Tags))
	for _, tag := range datasource.Tags.Tags {
		tags = append(tags, types.StringValue(tag.Label))
	}
	state.Tags, _ = types.ListValue(types.StringType, tags)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *datasourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
