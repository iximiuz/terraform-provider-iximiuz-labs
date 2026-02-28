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
	_ datasource.DataSource              = &PlaysDataSource{}
	_ datasource.DataSourceWithConfigure = &PlaysDataSource{}
)

type PlaysDataSource struct {
	client *client.Client
}

type PlaysDataSourceModel struct {
	Plays []PlaysDataSourceItemModel `tfsdk:"plays"`
}

type PlaysDataSourceItemModel struct {
	ID      types.String `tfsdk:"id"`
	Status  types.String `tfsdk:"status"`
	PageURL types.String `tfsdk:"page_url"`
}

func NewPlaysDataSource() datasource.DataSource {
	return &PlaysDataSource{}
}

func (d *PlaysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plays"
}

func (d *PlaysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all iximiuz Labs play sessions.",
		Attributes: map[string]schema.Attribute{
			"plays": schema.ListNestedAttribute{
				MarkdownDescription: "The list of play sessions.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the play session.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Current state of the play.",
							Computed:            true,
						},
						"page_url": schema.StringAttribute{
							MarkdownDescription: "Web UI URL for the play.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *PlaysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PlaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlaysDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Reading plays data source")

	plays, err := d.client.ListPlays(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing plays", err.Error())
		return
	}

	data.Plays = make([]PlaysDataSourceItemModel, len(plays))
	for i, play := range plays {
		item := PlaysDataSourceItemModel{
			ID:      types.StringValue(play.ID),
			PageURL: types.StringValue(play.PageURL),
		}

		if play.Status != nil {
			item.Status = types.StringValue(play.Status.CurrentState())
		} else {
			item.Status = types.StringNull()
		}

		data.Plays[i] = item
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
