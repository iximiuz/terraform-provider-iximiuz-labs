// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/iximiuz/terraform-provider-iximiuz-labs/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PlayResource{}
var _ resource.ResourceWithImportState = &PlayResource{}

type PlayResource struct {
	client *client.Client
}

type PlayResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	Playground              types.String `tfsdk:"playground"`
	SafetyDisclaimerConsent types.Bool   `tfsdk:"safety_disclaimer_consent"`
	InitConditions          types.Map    `tfsdk:"init_conditions"`
	Status                  types.String `tfsdk:"status"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
	ExpiresIn               types.Int64  `tfsdk:"expires_in"`
	PageURL                 types.String `tfsdk:"page_url"`
}

func NewPlayResource() resource.Resource {
	return &PlayResource{}
}

func (r *PlayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_play"
}

func (r *PlayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iximiuz Labs play session.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Server-generated play ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"playground": schema.StringAttribute{
				MarkdownDescription: "The playground name to instantiate.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"safety_disclaimer_consent": schema.BoolAttribute{
				MarkdownDescription: "Consent for community content safety disclaimer.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"init_conditions": schema.MapAttribute{
				MarkdownDescription: "Init condition key=value pairs.",
				Optional:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
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

func (r *PlayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *PlayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PlayResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreatePlayRequest{
		Playground:              data.Playground.ValueString(),
		SafetyDisclaimerConsent: data.SafetyDisclaimerConsent.ValueBool(),
	}

	if !data.InitConditions.IsNull() && !data.InitConditions.IsUnknown() {
		var conditions map[string]string
		resp.Diagnostics.Append(data.InitConditions.ElementsAs(ctx, &conditions, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.InitConditions = conditions
	}

	tflog.Debug(ctx, "Creating play", map[string]interface{}{
		"playground": createReq.Playground,
	})

	play, err := r.client.CreatePlay(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating play", err.Error())
		return
	}

	tflog.Debug(ctx, "Play created, waiting for RUNNING state", map[string]interface{}{
		"play_id": play.ID,
	})

	play, err = waitForPlayState(ctx, r.client, play.ID, []string{client.PlayStateRunning}, 10*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for play to reach RUNNING state", err.Error())
		return
	}

	mapPlayToModel(ctx, play, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PlayResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	play, err := r.client.GetPlay(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Debug(ctx, "Play not found, removing from state", map[string]interface{}{
				"play_id": data.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading play", err.Error())
		return
	}

	// Remove from state if the play is in a terminal state
	if play.Status != nil && (play.Status.CurrentState() == client.PlayStateDestroyed || play.Status.CurrentState() == client.PlayStateFailed || play.Status.CurrentState() == client.PlayStateStopped) {
		tflog.Debug(ctx, "Play is in terminal state, removing from state", map[string]interface{}{
			"play_id": play.ID,
			"state":   play.Status.CurrentState(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	mapPlayToModel(ctx, play, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Play resources do not support in-place updates. All input attributes use ForceNew, so any change should trigger replacement.",
	)
}

func (r *PlayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PlayResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Destroying play", map[string]interface{}{
		"play_id": data.ID.ValueString(),
	})

	err := r.client.DestroyPlay(ctx, data.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Debug(ctx, "Play already destroyed or not found", map[string]interface{}{
				"play_id": data.ID.ValueString(),
			})
			return
		}
		resp.Diagnostics.AddError("Error destroying play", err.Error())
		return
	}

	_, err = waitForPlayState(ctx, r.client, data.ID.ValueString(), []string{client.PlayStateDestroyed}, 5*time.Minute)
	if err != nil {
		// If the play is not found while polling, it's already gone
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error waiting for play to reach DESTROYED state", err.Error())
		return
	}
}

func (r *PlayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapPlayToModel maps a client.Play response to the Terraform resource model.
func mapPlayToModel(_ context.Context, play *client.Play, data *PlayResourceModel, _ *diag.Diagnostics) {
	data.ID = types.StringValue(play.ID)
	data.Playground = types.StringValue(play.Playground.Name)
	data.CreatedAt = types.StringValue(play.CreatedAt)
	data.UpdatedAt = types.StringValue(play.UpdatedAt)
	data.ExpiresIn = types.Int64Value(int64(play.ExpiresIn))
	data.PageURL = types.StringValue(play.PageURL)

	if play.Status != nil {
		data.Status = types.StringValue(play.Status.CurrentState())
	} else {
		data.Status = types.StringNull()
	}
}
