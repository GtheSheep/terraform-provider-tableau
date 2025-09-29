package tableau

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &workbookResource{}
	_ resource.ResourceWithConfigure   = &workbookResource{}
	_ resource.ResourceWithImportState = &workbookResource{}
)

func NewWorkbookResource() resource.Resource {
	return &workbookResource{}
}

type workbookResource struct {
	client *Client
}

type workbookResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ProjectID        types.String `tfsdk:"project_id"`
	ShowTabs         types.String `tfsdk:"show_tabs"`
	ThumbnailsUserID types.String `tfsdk:"thumbnails_user_id"`
	WorkbookFilename types.String `tfsdk:"workbook_filename"`
	WorkbookContent  types.String `tfsdk:"workbook_content"`
	Description      types.String `tfsdk:"description"`
	EncryptExtracts  types.String `tfsdk:"encrypt_extracts"`
	OwnerID          types.String `tfsdk:"owner_id"`
}

func (r *workbookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workbook"
}

func (r *workbookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // Update
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{ // Create && Update
				Required:    true,
				Description: "Workbook name",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description for the workbook",
			},
			"encrypt_extracts": schema.StringAttribute{
				Optional:    true,
				Description: "Whether or not extracts are encrypted",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"true",
						"false",
					}...),
				},
			},
			"owner_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the workbook owner",
			},
			"project_id": schema.StringAttribute{ // Create && Update
				Required:    true,
				Description: "Workbook belongs to project with ID",
			},
			"show_tabs": schema.StringAttribute{ // Create && Update
				Required:    true,
				Description: "Whether or not show views in tabs",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"true",
						"false",
					}...),
				},
			},
			"thumbnails_user_id": schema.StringAttribute{ // ONLY Create
				Required:    true,
				Description: "Specify user for thumbnail (used when creating workbook)",
			},
			"workbook_filename": schema.StringAttribute{ // ONLY Create
				Required:    true,
				Description: "Filename of workbook file (used when creating workbook)",
			},
			"workbook_content": schema.StringAttribute{ // ONLY Create
				Required:    true,
				Description: "Content of workbook file (used when creating workbook)",
			},
		},
	}
}

func (r *workbookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workbookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	name := string(plan.Name.ValueString())
	projectID := string(plan.ProjectID.ValueString())
	showTabs := string(plan.ShowTabs.ValueString())
	thumbnailsUserID := string(plan.ThumbnailsUserID.ValueString())
	wbContent := string(plan.WorkbookContent.ValueString())
	wbFilename := string(plan.WorkbookFilename.ValueString())
	id, err := r.client.CreateWorkbook(ctx, name, projectID, showTabs, thumbnailsUserID, wbFilename, []byte(wbContent))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workbook",
			"Could not create workbook, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(id)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workbookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workbook, err := r.client.GetWorkbook(state.ID.ValueString())
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(workbook.Name)
	state.ProjectID = types.StringValue(workbook.Project.ID)
	state.ShowTabs = types.StringValue(workbook.ShowTabs)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workbookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workbookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	wb := Workbook{
		ID:              plan.ID.ValueString(),
		Name:            plan.Name.ValueString(),
		Project:         WorkbookProject{ID: plan.ProjectID.ValueString()},
		ShowTabs:        plan.ShowTabs.ValueString(),
		Description:     plan.Description.ValueString(),
		EncryptExtracts: plan.EncryptExtracts.ValueString(),
		Owner:           WorkbookOwner{ID: plan.OwnerID.ValueString()},
	}
	_, err := r.client.UpdateWorkbook(
		wb.ID, wb.Name, wb.Project.ID, wb.ShowTabs, wb.Description, wb.EncryptExtracts, wb.Owner.ID,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Tableau Workbook",
			"Could not update workbook, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workbookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteWorkbook(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Workbook",
			"Could not delete workbook, unexpected error: "+err.Error(),
		)
	}
}

func (r *workbookResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *workbookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
