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
	_ datasource.DataSource              = &PlayDataSource{}
	_ datasource.DataSourceWithConfigure = &PlayDataSource{}
)

type PlayDataSource struct {
	client *client.Client
}

type PlayDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	ExpiresIn types.Int64  `tfsdk:"expires_in"`
	PageURL   types.String `tfsdk:"page_url"`
}

func NewPlayDataSource() datasource.DataSource {
	return &PlayDataSource{}
}

func (d *PlayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_play"
}

func (d *PlayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Looks up a single iximiuz Labs play session by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the play session.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Current state of the play.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "ISO timestamp of when the play was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "ISO timestamp of when the play was last updated.",
				Computed:            true,
			},
			"expires_in": schema.Int64Attribute{
				MarkdownDescription: "Seconds until the play expires.",
				Computed:            true,
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "Web UI URL for the play.",
				Computed:            true,
			},
		},
	}
}

func (d *PlayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PlayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlayDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Trace(ctx, "Reading play data source", map[string]interface{}{
		"id": id,
	})

	play, err := d.client.GetPlay(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading play", err.Error())
		return
	}

	data.ID = types.StringValue(play.ID)
	data.CreatedAt = types.StringValue(play.CreatedAt)
	data.UpdatedAt = types.StringValue(play.UpdatedAt)
	data.ExpiresIn = types.Int64Value(int64(play.ExpiresIn))
	data.PageURL = types.StringValue(play.PageURL)

	if play.Status != nil {
		data.Status = types.StringValue(play.Status.CurrentState())
	} else {
		data.Status = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
