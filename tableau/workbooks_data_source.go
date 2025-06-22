package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &workbooksDataSource{}
	_ datasource.DataSourceWithConfigure = &workbooksDataSource{}
)

func WorkbooksDataSource() datasource.DataSource {
	return &workbooksDataSource{}
}

type workbooksDataSource struct {
	client *Client
}

type workbooksNestedDataModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	EncryptExtracts types.String `tfsdk:"encrypt_extracts"`
	ShowTabs        types.String `tfsdk:"show_tabs"`
	Size            types.String `tfsdk:"size"`
	DefaultViewID   types.String `tfsdk:"default_view_id"`
	LocationID      types.String `tfsdk:"location_id"`
	OwnerID         types.String `tfsdk:"owner_id"`
	ProjectID       types.String `tfsdk:"project_id"`
	ContentURL      types.String `tfsdk:"content_url"`
	WebPageURL      types.String `tfsdk:"web_page_url"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	// Tags
}

type workbooksDataSourceModel struct {
	ID        types.String               `tfsdk:"id"`
	Workbooks []workbooksNestedDataModel `tfsdk:"workbooks"`
}

func (d *workbooksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workbooks"
}

func (d *workbooksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve workbooks details",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the workbooks",
			},
			"workbooks": schema.ListNestedAttribute{
				Description: "List of workbooks and their attributes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the workbook",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name for the workbook",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description for the workbook",
						},
						"encrypt_extracts": schema.StringAttribute{
							Computed:    true,
							Description: "Whether or not extracts are encrypted",
						},
						"show_tabs": schema.StringAttribute{
							Computed:    true,
							Description: "Whether or not this workbook show tabs",
						},
						"size": schema.StringAttribute{
							Computed:    true,
							Description: "Workbook size in mega bytes",
						},
						"default_view_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the workbook default view",
						},
						"location_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the workbook location",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the workbook project",
						},
						"owner_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the workbook owner",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Workbook was created at",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Workbook was updated at",
						},
						"content_url": schema.StringAttribute{
							Computed:    true,
							Description: "Content URL for the workbook",
						},
						"web_page_url": schema.StringAttribute{
							Computed:    true,
							Description: "Web page URL for the workbook",
						},
					},
				},
			},
		},
	}
}

func (d *workbooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state workbooksDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	workbooks, err := d.client.GetWorkbooks()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Tableau Workbooks",
			err.Error(),
		)
		return
	}

	for _, workbook := range workbooks {
		workbooksDataModel := workbooksNestedDataModel{
			ID:              types.StringValue(workbook.ID),
			Name:            types.StringValue(workbook.Name),
			Description:     types.StringValue(workbook.Description),
			EncryptExtracts: types.StringValue(workbook.EncryptExtracts),
			ShowTabs:        types.StringValue(workbook.ShowTabs),
			Size:            types.StringValue(workbook.Size),
			DefaultViewID:   types.StringValue(workbook.DefaultViewID),
			LocationID:      types.StringValue(workbook.Location.ID),
			OwnerID:         types.StringValue(workbook.Owner.ID),
			ProjectID:       types.StringValue(workbook.Project.ID),
			CreatedAt:       types.StringValue(workbook.CreatedAt),
			UpdatedAt:       types.StringValue(workbook.UpdatedAt),
			ContentURL:      types.StringValue(workbook.ContentURL),
			WebPageURL:      types.StringValue(workbook.WebPageURL),
		}
		state.Workbooks = append(state.Workbooks, workbooksDataModel)
	}

	state.ID = types.StringValue("allWorkbooks")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *workbooksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
