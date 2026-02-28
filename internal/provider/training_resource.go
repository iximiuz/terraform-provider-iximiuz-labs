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
	_ resource.Resource                = &TrainingResource{}
	_ resource.ResourceWithImportState = &TrainingResource{}
)

func NewTrainingResource() resource.Resource {
	return &TrainingResource{}
}

type TrainingResource struct {
	client *client.Client
}

type TrainingResourceModel struct {
	NamePrefix types.String `tfsdk:"name_prefix"`
	Name       types.String `tfsdk:"name"`
	Title      types.String `tfsdk:"title"`
	PageURL    types.String `tfsdk:"page_url"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *TrainingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_training"
}

func (r *TrainingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iximiuz-labs training.",
		Attributes: map[string]schema.Attribute{
			"name_prefix": schema.StringAttribute{
				MarkdownDescription: "A name prefix for the training. The API will append a unique suffix to produce the final `name`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The training name (slug), assigned by the server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The training display title.",
				Computed:            true,
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "The training web URL.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the training was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the training was last updated.",
				Computed:            true,
			},
		},
	}
}

func (r *TrainingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrainingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TrainingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.NamePrefix.ValueString()
	tflog.Trace(ctx, "creating training", map[string]interface{}{"name_prefix": name})

	training, err := r.client.CreateTraining(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create training, got error: %s", err))
		return
	}

	data.Name = types.StringValue(training.Name)
	data.Title = types.StringValue(training.Title)
	data.PageURL = types.StringValue(training.PageURL)
	data.CreatedAt = types.StringValue(training.CreatedAt)
	data.UpdatedAt = types.StringValue(training.UpdatedAt)

	tflog.Trace(ctx, "created training", map[string]interface{}{"name": training.Name})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrainingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TrainingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "reading training", map[string]interface{}{"name": name})

	training, err := r.client.GetTraining(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "training not found, removing from state", map[string]interface{}{"name": name})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read training, got error: %s", err))
		return
	}

	data.Name = types.StringValue(training.Name)
	data.Title = types.StringValue(training.Title)
	data.PageURL = types.StringValue(training.PageURL)
	data.CreatedAt = types.StringValue(training.CreatedAt)
	data.UpdatedAt = types.StringValue(training.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrainingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Training resources do not support in-place updates. All input attributes require replacement.",
	)
}

func (r *TrainingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TrainingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "deleting training", map[string]interface{}{"name": name})

	err := r.client.DeleteTraining(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "training already deleted or not found", map[string]interface{}{"name": name})
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete training, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted training", map[string]interface{}{"name": name})
}

func (r *TrainingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
