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
	_ resource.Resource                = &ChallengeResource{}
	_ resource.ResourceWithImportState = &ChallengeResource{}
)

func NewChallengeResource() resource.Resource {
	return &ChallengeResource{}
}

type ChallengeResource struct {
	client *client.Client
}

type ChallengeResourceModel struct {
	NamePrefix types.String `tfsdk:"name_prefix"`
	Name       types.String `tfsdk:"name"`
	Title      types.String `tfsdk:"title"`
	PageURL    types.String `tfsdk:"page_url"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *ChallengeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_challenge"
}

func (r *ChallengeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iximiuz-labs challenge.",
		Attributes: map[string]schema.Attribute{
			"name_prefix": schema.StringAttribute{
				MarkdownDescription: "A name prefix for the challenge. The API will append a unique suffix to produce the final `name`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The challenge name (slug), assigned by the server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The challenge display title.",
				Computed:            true,
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "The challenge web URL.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the challenge was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the challenge was last updated.",
				Computed:            true,
			},
		},
	}
}

func (r *ChallengeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ChallengeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ChallengeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.NamePrefix.ValueString()
	tflog.Trace(ctx, "creating challenge", map[string]interface{}{"name_prefix": name})

	challenge, err := r.client.CreateChallenge(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create challenge, got error: %s", err))
		return
	}

	data.Name = types.StringValue(challenge.Name)
	data.Title = types.StringValue(challenge.Title)
	data.PageURL = types.StringValue(challenge.PageURL)
	data.CreatedAt = types.StringValue(challenge.CreatedAt)
	data.UpdatedAt = types.StringValue(challenge.UpdatedAt)

	tflog.Trace(ctx, "created challenge", map[string]interface{}{"name": challenge.Name})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ChallengeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ChallengeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "reading challenge", map[string]interface{}{"name": name})

	challenge, err := r.client.GetChallenge(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "challenge not found, removing from state", map[string]interface{}{"name": name})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read challenge, got error: %s", err))
		return
	}

	data.Name = types.StringValue(challenge.Name)
	data.Title = types.StringValue(challenge.Title)
	data.PageURL = types.StringValue(challenge.PageURL)
	data.CreatedAt = types.StringValue(challenge.CreatedAt)
	data.UpdatedAt = types.StringValue(challenge.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ChallengeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Challenge resources do not support in-place updates. All input attributes require replacement.",
	)
}

func (r *ChallengeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ChallengeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "deleting challenge", map[string]interface{}{"name": name})

	err := r.client.DeleteChallenge(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "challenge already deleted or not found", map[string]interface{}{"name": name})
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete challenge, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted challenge", map[string]interface{}{"name": name})
}

func (r *ChallengeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
