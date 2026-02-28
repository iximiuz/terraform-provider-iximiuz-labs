// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/iximiuz/terraform-provider-iximiuz-labs/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PlayPortResource{}
var _ resource.ResourceWithImportState = &PlayPortResource{}

type PlayPortResource struct {
	client *client.Client
}

type PlayPortResourceModel struct {
	ID          types.String `tfsdk:"id"`
	PlayID      types.String `tfsdk:"play_id"`
	Machine     types.String `tfsdk:"machine"`
	Number      types.Int64  `tfsdk:"number"`
	Access      types.String `tfsdk:"access"`
	TLS         types.Bool   `tfsdk:"tls"`
	HostRewrite types.String `tfsdk:"host_rewrite"`
	PathRewrite types.String `tfsdk:"path_rewrite"`
	Hostname    types.String `tfsdk:"hostname"`
	URL         types.String `tfsdk:"url"`
}

func NewPlayPortResource() resource.Resource {
	return &PlayPortResource{}
}

func (r *PlayPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_play_port"
}

func (r *PlayPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Exposes an HTTP port on a play machine.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Server-generated port ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"play_id": schema.StringAttribute{
				MarkdownDescription: "The play ID to expose the port on.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"machine": schema.StringAttribute{
				MarkdownDescription: "The machine name to expose the port on.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: "The port number to expose.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"access": schema.StringAttribute{
				MarkdownDescription: "Access mode for the port: `private` or `public`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("private"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tls": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable TLS for the exposed port.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"host_rewrite": schema.StringAttribute{
				MarkdownDescription: "Host header rewrite value.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"path_rewrite": schema.StringAttribute{
				MarkdownDescription: "Path prefix rewrite value.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname assigned to the exposed port.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The full URL of the exposed port.",
				Computed:            true,
			},
		},
	}
}

func (r *PlayPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PlayPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PlayPortResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	playID := data.PlayID.ValueString()

	exposeReq := client.ExposePortRequest{
		Machine:     data.Machine.ValueString(),
		Number:      int(data.Number.ValueInt64()),
		Access:      data.Access.ValueString(),
		TLS:         data.TLS.ValueBool(),
		HostRewrite: data.HostRewrite.ValueString(),
		PathRewrite: data.PathRewrite.ValueString(),
	}

	tflog.Trace(ctx, "Exposing port on play machine", map[string]interface{}{
		"play_id": playID,
		"machine": exposeReq.Machine,
		"number":  exposeReq.Number,
	})

	port, err := r.client.ExposePort(ctx, playID, exposeReq)
	if err != nil {
		resp.Diagnostics.AddError("Error exposing port", err.Error())
		return
	}

	mapPortToModel(port, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlayPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PlayPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	playID := data.PlayID.ValueString()
	portID := data.ID.ValueString()

	tflog.Trace(ctx, "Reading port", map[string]interface{}{
		"play_id": playID,
		"port_id": portID,
	})

	ports, err := r.client.ListPorts(ctx, playID)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "Play not found, removing port from state", map[string]interface{}{
				"play_id": playID,
				"port_id": portID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error listing ports", err.Error())
		return
	}

	var found *client.Port
	for i := range ports {
		if ports[i].ID == portID {
			found = &ports[i]
			break
		}
	}

	if found == nil {
		tflog.Trace(ctx, "Port not found, removing from state", map[string]interface{}{
			"play_id": playID,
			"port_id": portID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	mapPortToModel(found, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlayPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Play port resources do not support in-place updates. All input attributes use ForceNew, so any change should trigger replacement.",
	)
}

func (r *PlayPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PlayPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	playID := data.PlayID.ValueString()
	portID := data.ID.ValueString()

	tflog.Trace(ctx, "Unexposing port", map[string]interface{}{
		"play_id": playID,
		"port_id": portID,
	})

	err := r.client.UnexposePort(ctx, playID, portID)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "Port already removed or play not found", map[string]interface{}{
				"play_id": playID,
				"port_id": portID,
			})
			return
		}
		resp.Diagnostics.AddError("Error unexposing port", err.Error())
		return
	}
}

func (r *PlayPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format play_id/port_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("play_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// mapPortToModel maps a client.Port response to the Terraform resource model.
func mapPortToModel(port *client.Port, data *PlayPortResourceModel) {
	data.ID = types.StringValue(port.ID)
	data.PlayID = types.StringValue(port.PlayID)
	data.Machine = types.StringValue(port.Machine)
	data.Number = types.Int64Value(int64(port.Number))
	data.Access = types.StringValue(port.AccessMode)
	data.TLS = types.BoolValue(port.TLS)

	if port.Hostname != "" {
		data.Hostname = types.StringValue(port.Hostname)
	} else {
		data.Hostname = types.StringNull()
	}

	if port.URL != "" {
		data.URL = types.StringValue(port.URL)
	} else {
		data.URL = types.StringNull()
	}

	if port.HostRewrite != "" {
		data.HostRewrite = types.StringValue(port.HostRewrite)
	} else {
		data.HostRewrite = types.StringNull()
	}

	if port.PathRewrite != "" {
		data.PathRewrite = types.StringValue(port.PathRewrite)
	} else {
		data.PathRewrite = types.StringNull()
	}
}
