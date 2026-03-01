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
	_ datasource.DataSource              = &MeDataSource{}
	_ datasource.DataSourceWithConfigure = &MeDataSource{}
)

type MeDataSource struct {
	client *client.Client
}

type MeDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	GithubProfileID types.String `tfsdk:"github_profile_id"`
	HasPremium      types.Bool   `tfsdk:"has_premium"`
	PremiumLifetime types.Bool   `tfsdk:"premium_lifetime"`
	PremiumTrial    types.Bool   `tfsdk:"premium_trial"`
	PremiumUntil    types.String `tfsdk:"premium_until"`
}

func NewMeDataSource() datasource.DataSource {
	return &MeDataSource{}
}

func (d *MeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_me"
}

func (d *MeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about the current authenticated iximiuz Labs user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique user ID.",
				Computed:            true,
			},
			"github_profile_id": schema.StringAttribute{
				MarkdownDescription: "The user's GitHub profile ID.",
				Computed:            true,
			},
			"has_premium": schema.BoolAttribute{
				MarkdownDescription: "Whether the user has an active premium subscription.",
				Computed:            true,
			},
			"premium_lifetime": schema.BoolAttribute{
				MarkdownDescription: "Whether the user has a lifetime premium subscription.",
				Computed:            true,
			},
			"premium_trial": schema.BoolAttribute{
				MarkdownDescription: "Whether the user is on a premium trial.",
				Computed:            true,
			},
			"premium_until": schema.StringAttribute{
				MarkdownDescription: "ISO timestamp of when the premium subscription expires.",
				Computed:            true,
			},
		},
	}
}

func (d *MeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Reading me data source")

	me, err := d.client.GetMe(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading current user", err.Error())
		return
	}

	data.ID = types.StringValue(me.ID)
	data.GithubProfileID = types.StringValue(me.GithubProfileID)

	if me.PremiumAccess != nil {
		data.HasPremium = types.BoolValue(true)
		data.PremiumLifetime = types.BoolValue(me.PremiumAccess.Lifetime)
		data.PremiumTrial = types.BoolValue(me.PremiumAccess.Trial)
		data.PremiumUntil = types.StringValue(me.PremiumAccess.Until)
	} else {
		data.HasPremium = types.BoolValue(false)
		data.PremiumLifetime = types.BoolValue(false)
		data.PremiumTrial = types.BoolValue(false)
		data.PremiumUntil = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
