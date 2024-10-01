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
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	client *Client
}

type userResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	Name        types.String `tfsdk:"name"`
	FullName    types.String `tfsdk:"full_name"`
	SiteRole    types.String `tfsdk:"site_role"`
	AuthSetting types.String `tfsdk:"auth_setting"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "User email",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for user",
			},
			"full_name": schema.StringAttribute{
				Required:    true,
				Description: "Full name for user - Note: Can't be updated due to permissioning when using SAML",
			},
			"site_role": schema.StringAttribute{
				Required:    true,
				Description: "Site role for the user",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Creator",
						"Explorer",
						"Interactor",
						"Publisher",
						"ExplorerCanPublish",
						"ServerAdministrator",
						"SiteAdministratorExplorer",
						"SiteAdministratorCreator",
						"Unlicensed",
						"Viewer",
					}...),
				},
			},
			"auth_setting": schema.StringAttribute{
				Required:    true,
				Description: "Auth setting for the user",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"SAML",
						"ServerDefault",
						"OpenID",
						"TableauIDWithMFA",
					}...),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := User{
		Email:       string(plan.Email.ValueString()),
		Name:        string(plan.Name.ValueString()),
		FullName:    string(plan.FullName.ValueString()),
		SiteRole:    string(plan.SiteRole.ValueString()),
		AuthSetting: string(plan.AuthSetting.ValueString()),
	}

	createdUser, err := r.client.CreateUser(user.Email, user.Name, user.FullName, user.SiteRole, user.AuthSetting)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}
	_, err = r.client.UpdateUser(createdUser.ID, user.Name, user.SiteRole, user.AuthSetting)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user during create",
			"Could not update user during create create, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(createdUser.ID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


	user, err := r.client.GetUser(state.ID.ValueString()) // check to see if user exists
	if err != nil {
		user, err = r.client.CreateUser(state.Email.ValueString(), state.Name.ValueString(), state.FullName.ValueString(), state.SiteRole.ValueString(), state.AuthSetting.ValueString())
		if err != nil {
			return
		}
	}

	state.ID = types.StringValue(user.ID)
	state.Email = types.StringValue(user.Email)
	state.Name = types.StringValue(user.Name)
	state.FullName = types.StringValue(user.FullName)
	state.SiteRole = types.StringValue(user.SiteRole)
	state.AuthSetting = types.StringValue(user.AuthSetting)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := User{
		Email:       string(plan.Email.ValueString()),
		Name:        string(plan.Name.ValueString()),
		SiteRole:    string(plan.SiteRole.ValueString()),
		AuthSetting: string(plan.AuthSetting.ValueString()),
	}

	_, err := r.client.UpdateUser(plan.ID.ValueString(), user.Name, user.SiteRole, user.AuthSetting)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Tableau User",
			"Could not update user, unexpected error: "+err.Error(),
		)
		return
	}

	updatedUser, err := r.client.GetUser(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tableau User",
			"Could not read Tableau user ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(updatedUser.Name)
	plan.SiteRole = types.StringValue(updatedUser.SiteRole)
	plan.AuthSetting = types.StringValue(updatedUser.AuthSetting)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Tableau User",
			"Could not delete user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
