package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &virtualConnectionRevisionsDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualConnectionRevisionsDataSource{}
)

func VirtualConnectionRevisionsDataSource() datasource.DataSource {
	return &virtualConnectionRevisionsDataSource{}
}

type virtualConnectionRevisionsDataSource struct {
	client *Client
}

type virtualConnectionRevisionNestedDataModel struct {
	PublisherID    types.String `tfsdk:"publisher_id"`
	Current        types.Bool   `tfsdk:"current"`
	Deleted        types.Bool   `tfsdk:"deleted"`
	PublishedAt    types.String `tfsdk:"published_at"`
	RevisionNumber types.String `tfsdk:"revision_number"`
}

type virtualConnectionRevisionsDataSourceModel struct {
	ID        types.String                               `tfsdk:"id"`
	Revisions []virtualConnectionRevisionNestedDataModel `tfsdk:"revisions"`
}

func (d *virtualConnectionRevisionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_connection_revisions"
}

func (d *virtualConnectionRevisionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve virtual connection revisions details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the virtual connections",
			},
			"revisions": schema.ListNestedAttribute{
				Description: "List database connections of virtual connection and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"publisher_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the publisher",
						},
						"current": schema.BoolAttribute{
							Computed:    true,
							Description: "Current revision",
						},
						"deleted": schema.BoolAttribute{
							Computed:    true,
							Description: "Deleted revision",
						},
						"published_at": schema.StringAttribute{
							Computed:    true,
							Description: "Published at given date",
						},
						"revision_number": schema.StringAttribute{
							Computed:    true,
							Description: "Revision number",
						},
					},
				},
			},
		},
	}
}

func (d *virtualConnectionRevisionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state virtualConnectionRevisionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	revisions, err := d.client.GetVirtualConnectionRevisions(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Virtual Connection Revisions",
			err.Error(),
		)
		return
	}
	for _, revision := range revisions {
		virtualConnectionRevision := virtualConnectionRevisionNestedDataModel{
			PublisherID:    types.StringValue(revision.Publisher.ID),
			Current:        types.BoolValue(revision.Current),
			Deleted:        types.BoolValue(revision.Deleted),
			PublishedAt:    types.StringValue(revision.PublishedAt),
			RevisionNumber: types.StringValue(revision.RevisionNumber),
		}
		state.Revisions = append(state.Revisions, virtualConnectionRevision)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *virtualConnectionRevisionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
