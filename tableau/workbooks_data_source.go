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
	LocationID      types.String `tfsdk:"location_id"`
	LocationType    types.String `tfsdk:"location_type"`
	LocationName    types.String `tfsdk:"location_name"`
	OwnerID         types.String `tfsdk:"owner_id"`
	OwnerName       types.String `tfsdk:"owner_name"`
	ProjectID       types.String `tfsdk:"project_id"`
	ProjectName     types.String `tfsdk:"project_name"`
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
							Description: "Whether or not show views in tabs",
						},
						"size": schema.StringAttribute{
							Computed:    true,
							Description: "Workbook size in mega bytes",
						},
						"location_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the workbook location",
						},
						"location_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the workbook location",
						},
						"location_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the workbook location",
						},
						"owner_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the workbook owner",
						},
						"owner_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the workbook owner",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the project",
						},
						"project_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the project",
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
			OwnerID:         types.StringValue(workbook.Owner.ID),
			OwnerName:       types.StringValue(workbook.Owner.Name),
			ProjectID:       types.StringValue(workbook.Project.ID),
			ProjectName:     types.StringValue(workbook.Project.Name),
			CreatedAt:       types.StringValue(workbook.CreatedAt),
			UpdatedAt:       types.StringValue(workbook.UpdatedAt),
			ContentURL:      types.StringValue(workbook.ContentURL),
			WebPageURL:      types.StringValue(workbook.WebPageURL),
			LocationID:      types.StringValue(workbook.Location.ID),
			LocationType:    types.StringValue(workbook.Location.Type),
			LocationName:    types.StringValue(workbook.Location.Name),
		}
		state.Workbooks = append(state.Workbooks, workbooksDataModel)
	}

	state.ID = types.StringValue("allWorkbooks")

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *workbooksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*Client)
}
