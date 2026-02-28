// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/iximiuz/terraform-provider-iximiuz-labs/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &SkillPathResource{}
	_ resource.ResourceWithImportState = &SkillPathResource{}
)

func NewSkillPathResource() resource.Resource {
	return &SkillPathResource{}
}

type SkillPathResource struct {
	client *client.Client
}

type SkillPathResourceModel struct {
	NamePrefix types.String `tfsdk:"name_prefix"`
	Name       types.String `tfsdk:"name"`
	Title      types.String `tfsdk:"title"`
	PageURL    types.String `tfsdk:"page_url"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *SkillPathResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_path"
}

func (r *SkillPathResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iximiuz-labs skill path.",
		Attributes: map[string]schema.Attribute{
			"name_prefix": schema.StringAttribute{
				MarkdownDescription: "A name prefix for the skill path. The API will append a unique suffix to produce the final `name`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The skill path name (slug), assigned by the server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The skill path display title.",
				Computed:            true,
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "The skill path web URL.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the skill path was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the skill path was last updated.",
				Computed:            true,
			},
		},
	}
}

func (r *SkillPathResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *SkillPathResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SkillPathResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.NamePrefix.ValueString()
	tflog.Trace(ctx, "creating skill path", map[string]interface{}{"name_prefix": name})

	skillPath, err := r.client.CreateSkillPath(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create skill path, got error: %s", err))
		return
	}

	data.Name = types.StringValue(skillPath.Name)
	data.Title = types.StringValue(skillPath.Title)
	data.PageURL = types.StringValue(skillPath.PageURL)
	data.CreatedAt = types.StringValue(skillPath.CreatedAt)
	data.UpdatedAt = types.StringValue(skillPath.UpdatedAt)

	tflog.Trace(ctx, "created skill path", map[string]interface{}{"name": skillPath.Name})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SkillPathResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SkillPathResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "reading skill path", map[string]interface{}{"name": name})

	skillPath, err := r.client.GetSkillPath(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "skill path not found, removing from state", map[string]interface{}{"name": name})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read skill path, got error: %s", err))
		return
	}

	data.Name = types.StringValue(skillPath.Name)
	data.Title = types.StringValue(skillPath.Title)
	data.PageURL = types.StringValue(skillPath.PageURL)
	data.CreatedAt = types.StringValue(skillPath.CreatedAt)
	data.UpdatedAt = types.StringValue(skillPath.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SkillPathResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Skill path resources do not support in-place updates. All input attributes require replacement.",
	)
}

func (r *SkillPathResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SkillPathResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "deleting skill path", map[string]interface{}{"name": name})

	err := r.client.DeleteSkillPath(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "skill path already deleted or not found", map[string]interface{}{"name": name})
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete skill path, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted skill path", map[string]interface{}{"name": name})
}

func (r *SkillPathResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
