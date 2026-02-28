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
	_ datasource.DataSource              = &PlaygroundsDataSource{}
	_ datasource.DataSourceWithConfigure = &PlaygroundsDataSource{}
)

type PlaygroundsDataSource struct {
	client *client.Client
}

type PlaygroundsDataSourceModel struct {
	Filter      types.String                     `tfsdk:"filter"`
	Playgrounds []PlaygroundsDataSourceItemModel `tfsdk:"playgrounds"`
}

type PlaygroundsDataSourceItemModel struct {
	Name        types.String `tfsdk:"name"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	PageURL     types.String `tfsdk:"page_url"`
}

func NewPlaygroundsDataSource() datasource.DataSource {
	return &PlaygroundsDataSource{}
}

func (d *PlaygroundsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playgrounds"
}

func (d *PlaygroundsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists iximiuz Labs playgrounds with an optional filter.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				MarkdownDescription: "Optional filter string to narrow the list of playgrounds.",
				Optional:            true,
			},
			"playgrounds": schema.ListNestedAttribute{
				MarkdownDescription: "The list of playgrounds matching the filter.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The unique name of the playground.",
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
					},
				},
			},
		},
	}
}

func (d *PlaygroundsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PlaygroundsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlaygroundsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := ""
	if !data.Filter.IsNull() && !data.Filter.IsUnknown() {
		filter = data.Filter.ValueString()
	}

	tflog.Trace(ctx, "Reading playgrounds data source", map[string]interface{}{
		"filter": filter,
	})

	playgrounds, err := d.client.ListPlaygrounds(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError("Error listing playgrounds", err.Error())
		return
	}

	data.Playgrounds = make([]PlaygroundsDataSourceItemModel, len(playgrounds))
	for i, pg := range playgrounds {
		data.Playgrounds[i] = PlaygroundsDataSourceItemModel{
			Name:        types.StringValue(pg.Name),
			Title:       types.StringValue(pg.Title),
			Description: types.StringValue(pg.Description),
			PageURL:     types.StringValue(pg.PageURL),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
