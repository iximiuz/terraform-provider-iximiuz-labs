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
	_ resource.Resource                = &TutorialResource{}
	_ resource.ResourceWithImportState = &TutorialResource{}
)

func NewTutorialResource() resource.Resource {
	return &TutorialResource{}
}

type TutorialResource struct {
	client *client.Client
}

type TutorialResourceModel struct {
	NamePrefix types.String `tfsdk:"name_prefix"`
	Name       types.String `tfsdk:"name"`
	Title      types.String `tfsdk:"title"`
	PageURL    types.String `tfsdk:"page_url"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *TutorialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tutorial"
}

func (r *TutorialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iximiuz-labs tutorial.",
		Attributes: map[string]schema.Attribute{
			"name_prefix": schema.StringAttribute{
				MarkdownDescription: "A name prefix for the tutorial. The API will append a unique suffix to produce the final `name`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The tutorial name (slug), assigned by the server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The tutorial display title.",
				Computed:            true,
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "The tutorial web URL.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the tutorial was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the tutorial was last updated.",
				Computed:            true,
			},
		},
	}
}

func (r *TutorialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TutorialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TutorialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.NamePrefix.ValueString()
	tflog.Trace(ctx, "creating tutorial", map[string]interface{}{"name_prefix": name})

	tutorial, err := r.client.CreateTutorial(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tutorial, got error: %s", err))
		return
	}

	data.Name = types.StringValue(tutorial.Name)
	data.Title = types.StringValue(tutorial.Title)
	data.PageURL = types.StringValue(tutorial.PageURL)
	data.CreatedAt = types.StringValue(tutorial.CreatedAt)
	data.UpdatedAt = types.StringValue(tutorial.UpdatedAt)

	tflog.Trace(ctx, "created tutorial", map[string]interface{}{"name": tutorial.Name})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TutorialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TutorialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "reading tutorial", map[string]interface{}{"name": name})

	tutorial, err := r.client.GetTutorial(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "tutorial not found, removing from state", map[string]interface{}{"name": name})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tutorial, got error: %s", err))
		return
	}

	data.Name = types.StringValue(tutorial.Name)
	data.Title = types.StringValue(tutorial.Title)
	data.PageURL = types.StringValue(tutorial.PageURL)
	data.CreatedAt = types.StringValue(tutorial.CreatedAt)
	data.UpdatedAt = types.StringValue(tutorial.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TutorialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Tutorial resources do not support in-place updates. All input attributes require replacement.",
	)
}

func (r *TutorialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TutorialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "deleting tutorial", map[string]interface{}{"name": name})

	err := r.client.DeleteTutorial(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "tutorial already deleted or not found", map[string]interface{}{"name": name})
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tutorial, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted tutorial", map[string]interface{}{"name": name})
}

func (r *TutorialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
