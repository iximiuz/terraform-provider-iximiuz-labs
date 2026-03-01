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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PlayShellResource{}
var _ resource.ResourceWithImportState = &PlayShellResource{}

type PlayShellResource struct {
	client *client.Client
}

type PlayShellResourceModel struct {
	ID       types.String `tfsdk:"id"`
	PlayID   types.String `tfsdk:"play_id"`
	Machine  types.String `tfsdk:"machine"`
	User     types.String `tfsdk:"user"`
	Access   types.String `tfsdk:"access"`
	Hostname types.String `tfsdk:"hostname"`
	URL      types.String `tfsdk:"url"`
}

func NewPlayShellResource() resource.Resource {
	return &PlayShellResource{}
}

func (r *PlayShellResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_play_shell"
}

func (r *PlayShellResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Exposes a web terminal (shell) on a play machine.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Server-generated shell ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"play_id": schema.StringAttribute{
				MarkdownDescription: "The play ID to expose the shell on.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"machine": schema.StringAttribute{
				MarkdownDescription: "The machine name to expose the shell on.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The user to open the shell as.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access": schema.StringAttribute{
				MarkdownDescription: "Access mode for the shell: `private` or `public`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname assigned to the exposed shell.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The full URL of the exposed shell.",
				Computed:            true,
			},
		},
	}
}

func (r *PlayShellResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PlayShellResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PlayShellResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	playID := data.PlayID.ValueString()

	exposeReq := client.ExposeShellRequest{
		Machine: data.Machine.ValueString(),
		User:    data.User.ValueString(),
		Access:  data.Access.ValueString(),
	}

	tflog.Trace(ctx, "Exposing shell on play machine", map[string]interface{}{
		"play_id": playID,
		"machine": exposeReq.Machine,
		"user":    exposeReq.User,
	})

	shell, err := r.client.ExposeShell(ctx, playID, exposeReq)
	if err != nil {
		resp.Diagnostics.AddError("Error exposing shell", err.Error())
		return
	}

	mapShellToModel(shell, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlayShellResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PlayShellResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	playID := data.PlayID.ValueString()
	shellID := data.ID.ValueString()

	tflog.Trace(ctx, "Reading shell", map[string]interface{}{
		"play_id":  playID,
		"shell_id": shellID,
	})

	shells, err := r.client.ListShells(ctx, playID)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "Play not found, removing shell from state", map[string]interface{}{
				"play_id":  playID,
				"shell_id": shellID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error listing shells", err.Error())
		return
	}

	var found *client.Shell
	for i := range shells {
		if shells[i].ID == shellID {
			found = &shells[i]
			break
		}
	}

	if found == nil {
		tflog.Trace(ctx, "Shell not found, removing from state", map[string]interface{}{
			"play_id":  playID,
			"shell_id": shellID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	mapShellToModel(found, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlayShellResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Play shell resources do not support in-place updates. All input attributes use ForceNew, so any change should trigger replacement.",
	)
}

func (r *PlayShellResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PlayShellResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	playID := data.PlayID.ValueString()
	shellID := data.ID.ValueString()

	tflog.Trace(ctx, "Unexposing shell", map[string]interface{}{
		"play_id":  playID,
		"shell_id": shellID,
	})

	err := r.client.UnexposeShell(ctx, playID, shellID)
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "Shell already removed or play not found", map[string]interface{}{
				"play_id":  playID,
				"shell_id": shellID,
			})
			return
		}
		resp.Diagnostics.AddError("Error unexposing shell", err.Error())
		return
	}
}

func (r *PlayShellResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format play_id/shell_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("play_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// mapShellToModel maps a client.Shell response to the Terraform resource model.
func mapShellToModel(shell *client.Shell, data *PlayShellResourceModel) {
	data.ID = types.StringValue(shell.ID)
	data.PlayID = types.StringValue(shell.PlayID)
	data.Machine = types.StringValue(shell.Machine)
	data.Access = types.StringValue(shell.AccessMode)

	if shell.Hostname != "" {
		data.Hostname = types.StringValue(shell.Hostname)
	} else {
		data.Hostname = types.StringNull()
	}

	if shell.URL != "" {
		data.URL = types.StringValue(shell.URL)
	} else {
		data.URL = types.StringNull()
	}

	if shell.User != "" {
		data.User = types.StringValue(shell.User)
	} else {
		data.User = types.StringNull()
	}
}
