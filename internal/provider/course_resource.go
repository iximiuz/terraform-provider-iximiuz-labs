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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &CourseResource{}
	_ resource.ResourceWithImportState = &CourseResource{}
)

func NewCourseResource() resource.Resource {
	return &CourseResource{}
}

type CourseResource struct {
	client *client.Client
}

type CourseResourceModel struct {
	NamePrefix types.String `tfsdk:"name_prefix"`
	Name       types.String `tfsdk:"name"`
	Variant    types.String `tfsdk:"variant"`
	Sample     types.Bool   `tfsdk:"sample"`
	Title      types.String `tfsdk:"title"`
	PageURL    types.String `tfsdk:"page_url"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *CourseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_course"
}

func (r *CourseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iximiuz-labs course.",
		Attributes: map[string]schema.Attribute{
			"name_prefix": schema.StringAttribute{
				MarkdownDescription: "A name prefix for the course. The API will append a unique suffix to produce the final `name`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The course name (slug), assigned by the server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"variant": schema.StringAttribute{
				MarkdownDescription: "The course variant (`simple` or `modular`).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sample": schema.BoolAttribute{
				MarkdownDescription: "Whether to generate sample content for the course.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The course display title.",
				Computed:            true,
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "The course web URL.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the course was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the course was last updated.",
				Computed:            true,
			},
		},
	}
}

func (r *CourseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CourseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CourseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateCourseRequest{
		Name: data.NamePrefix.ValueString(),
	}
	if !data.Variant.IsNull() && !data.Variant.IsUnknown() {
		createReq.Variant = client.CourseVariant(data.Variant.ValueString())
	}
	if !data.Sample.IsNull() && !data.Sample.IsUnknown() {
		createReq.Sample = data.Sample.ValueBool()
	}

	tflog.Trace(ctx, "creating course", map[string]interface{}{"name": createReq.Name})

	course, err := r.client.CreateCourse(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create course, got error: %s", err))
		return
	}

	data.Name = types.StringValue(course.Name)
	data.Title = types.StringValue(course.Title)
	data.PageURL = types.StringValue(course.PageURL)
	data.CreatedAt = types.StringValue(course.CreatedAt)
	data.UpdatedAt = types.StringValue(course.UpdatedAt)

	tflog.Trace(ctx, "created course", map[string]interface{}{"name": course.Name})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CourseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CourseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "reading course", map[string]interface{}{"name": name})

	course, err := r.client.GetCourse(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "course not found, removing from state", map[string]interface{}{"name": name})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read course, got error: %s", err))
		return
	}

	data.Name = types.StringValue(course.Name)
	data.Title = types.StringValue(course.Title)
	data.PageURL = types.StringValue(course.PageURL)
	data.CreatedAt = types.StringValue(course.CreatedAt)
	data.UpdatedAt = types.StringValue(course.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CourseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Course resources do not support in-place updates. All input attributes require replacement.",
	)
}

func (r *CourseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CourseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	tflog.Trace(ctx, "deleting course", map[string]interface{}{"name": name})

	err := r.client.DeleteCourse(ctx, name)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "course already deleted or not found", map[string]interface{}{"name": name})
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete course, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted course", map[string]interface{}{"name": name})
}

func (r *CourseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
