// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/iximiuz/terraform-provider-iximiuz-labs/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &PlaygroundDataSource{}
	_ datasource.DataSourceWithConfigure = &PlaygroundDataSource{}
)

type PlaygroundDataSource struct {
	client *client.Client
}

type PlaygroundDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	PageURL     types.String `tfsdk:"page_url"`
	Published   types.Bool   `tfsdk:"published"`
	Owner       types.String `tfsdk:"owner"`
	Base        types.String `tfsdk:"base"`
	Cover       types.String `tfsdk:"cover"`
}

func NewPlaygroundDataSource() datasource.DataSource {
	return &PlaygroundDataSource{}
}

func (d *PlaygroundDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playground"
}

func (d *PlaygroundDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Looks up a single iximiuz Labs playground by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The unique name of the playground.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Server-generated playground ID.",
				Computed:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The display title of the playground.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A short description of the playground.",
				Computed:            true,
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "Web UI URL for the playground.",
				Computed:            true,
			},
			"published": schema.BoolAttribute{
				MarkdownDescription: "Whether the playground is published.",
				Computed:            true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The owner of the playground.",
				Computed:            true,
			},
			"base": schema.StringAttribute{
				MarkdownDescription: "The base playground this playground extends.",
				Computed:            true,
			},
			"cover": schema.StringAttribute{
				MarkdownDescription: "The cover image URL of the playground.",
				Computed:            true,
			},
		},
	}
}

func (d *PlaygroundDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *PlaygroundDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlaygroundDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()

	tflog.Trace(ctx, "Reading playground data source", map[string]interface{}{
		"name": name,
	})

	pg, err := d.client.GetPlayground(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Error reading playground", err.Error())
		return
	}

	data.ID = types.StringValue(pg.ID)
	data.Name = types.StringValue(pg.Name)
	data.Title = types.StringValue(pg.Title)
	data.Description = types.StringValue(pg.Description)
	data.PageURL = types.StringValue(pg.PageURL)
	data.Published = types.BoolValue(pg.Published)
	data.Owner = types.StringValue(pg.Owner)
	data.Base = types.StringValue(pg.Base)
	data.Cover = types.StringValue(pg.Cover)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
