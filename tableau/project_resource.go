package tableau

import (
	"context"
	"time"

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
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	client *Client
}

type projectResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	ParentProjectID    types.String `tfsdk:"parent_project_id"`
	Description        types.String `tfsdk:"description"`
	ContentPermissions types.String `tfsdk:"content_permissions"`
	LastUpdated        types.String `tfsdk:"last_updated"`
}

func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for project",
			},
			"parent_project_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifier for the parent project",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description for the project",
			},
			"content_permissions": schema.StringAttribute{
				Required:    true,
				Description: "Permissions for the project content - ManagedByOwner is the default",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"LockedToProject",
						"ManagedByOwner",
						"LockedToProjectWithoutNested",
					}...),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := Project{
		Name:               string(plan.Name.ValueString()),
		ParentProjectID:    string(plan.ParentProjectID.ValueString()),
		Description:        string(plan.Description.ValueString()),
		ContentPermissions: string(plan.ContentPermissions.ValueString()),
	}

	createdProject, err := r.client.CreateProject(project.Name, project.ParentProjectID, project.Description, project.ContentPermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(createdProject.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Project",
			"Could not read Tableau project ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(project.ID)
	state.Name = types.StringValue(project.Name)
	state.ParentProjectID = types.StringValue(project.ParentProjectID)
	state.Description = types.StringValue(project.Description)
	state.ContentPermissions = types.StringValue(project.ContentPermissions)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := Project{
		Name:               string(plan.Name.ValueString()),
		ParentProjectID:    string(plan.ParentProjectID.ValueString()),
		Description:        string(plan.Description.ValueString()),
		ContentPermissions: string(plan.ContentPermissions.ValueString()),
	}

	_, err := r.client.UpdateProject(plan.ID.ValueString(), project.Name, project.ParentProjectID, project.Description, project.ContentPermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Tableau Project",
			"Could not update project, unexpected error: "+err.Error(),
		)
		return
	}

	updatedProject, err := r.client.GetProject(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau Project",
			"Could not read Tableau project ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(updatedProject.Name)
	plan.ParentProjectID = types.StringValue(updatedProject.ParentProjectID)
	plan.Description = types.StringValue(updatedProject.Description)
	plan.ContentPermissions = types.StringValue(updatedProject.ContentPermissions)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau Project",
			"Could not delete project, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
