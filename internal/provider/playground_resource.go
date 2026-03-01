// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

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

// Ensure provider-defined types fully satisfy framework interfaces.
var _ resource.Resource = &PlaygroundResource{}
var _ resource.ResourceWithImportState = &PlaygroundResource{}

// PlaygroundResource defines the resource implementation.
type PlaygroundResource struct {
	client *client.Client
}

// PlaygroundResourceModel describes the resource data model.
type PlaygroundResourceModel struct {
	ID            types.String              `tfsdk:"id"`
	NamePrefix    types.String              `tfsdk:"name_prefix"`
	Name          types.String              `tfsdk:"name"`
	Base          types.String              `tfsdk:"base"`
	Title         types.String              `tfsdk:"title"`
	Description   types.String              `tfsdk:"description"`
	Categories    types.List                `tfsdk:"categories"`
	Cover         types.String              `tfsdk:"cover"`
	Markdown      types.String              `tfsdk:"markdown"`
	Published     types.Bool                `tfsdk:"published"`
	PageURL       types.String              `tfsdk:"page_url"`
	Owner         types.String              `tfsdk:"owner"`
	RegistryAuth  types.String              `tfsdk:"registry_auth"`
	Networks      []NetworkModel            `tfsdk:"network"`
	Machines      []MachineModel            `tfsdk:"machine"`
	Tabs          []TabModel                `tfsdk:"tab"`
	InitTasks     []InitTaskModel           `tfsdk:"init_task"`
	InitConds     *InitConditionsBlockModel `tfsdk:"init_conditions"`
	AccessControl *AccessControlModel       `tfsdk:"access_control"`
	PortForwards  []PortForwardModel        `tfsdk:"port_forward"`
}

// NetworkModel describes a playground network.
type NetworkModel struct {
	Name    types.String `tfsdk:"name"`
	Subnet  types.String `tfsdk:"subnet"`
	Gateway types.String `tfsdk:"gateway"`
	Private types.Bool   `tfsdk:"private"`
}

// MachineModel describes a playground machine.
type MachineModel struct {
	Name         types.String          `tfsdk:"name"`
	NoSSH        types.Bool            `tfsdk:"no_ssh"`
	Users        []MachineUserModel    `tfsdk:"user"`
	Kernel       *MachineKernelModel   `tfsdk:"kernel"`
	Drives       []MachineDriveModel   `tfsdk:"drive"`
	Network      *MachineNetworkModel  `tfsdk:"network"`
	Resources    *MachineResourceModel `tfsdk:"resources"`
	StartupFiles []StartupFileModel    `tfsdk:"startup_file"`
}

// MachineUserModel describes a user on a machine.
type MachineUserModel struct {
	Name    types.String `tfsdk:"name"`
	Default types.Bool   `tfsdk:"default"`
	Welcome types.String `tfsdk:"welcome"`
}

// MachineKernelModel describes kernel configuration.
type MachineKernelModel struct {
	Source types.String `tfsdk:"source"`
}

// MachineDriveModel describes a machine drive.
type MachineDriveModel struct {
	Source     types.String `tfsdk:"source"`
	Mount      types.String `tfsdk:"mount"`
	Size       types.String `tfsdk:"size"`
	Filesystem types.String `tfsdk:"filesystem"`
	ReadOnly   types.Bool   `tfsdk:"read_only"`
	Persistent types.Bool   `tfsdk:"persistent"`
}

// MachineNetworkModel describes network config for a machine.
type MachineNetworkModel struct {
	Interfaces []MachineNetworkInterfaceModel `tfsdk:"interface"`
}

// MachineNetworkInterfaceModel describes a network interface.
type MachineNetworkInterfaceModel struct {
	Address types.String `tfsdk:"address"`
	Network types.String `tfsdk:"network"`
}

// MachineResourceModel describes resource limits for a machine.
type MachineResourceModel struct {
	CPUCount types.Int64  `tfsdk:"cpu_count"`
	RAMSize  types.String `tfsdk:"ram_size"`
}

// StartupFileModel describes a startup file.
type StartupFileModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
	Mode    types.String `tfsdk:"mode"`
	Owner   types.String `tfsdk:"owner"`
}

// TabModel describes a playground tab.
type TabModel struct {
	Kind        types.String `tfsdk:"kind"`
	Name        types.String `tfsdk:"name"`
	Machine     types.String `tfsdk:"machine"`
	Number      types.Int64  `tfsdk:"number"`
	Access      types.String `tfsdk:"access"`
	TLS         types.Bool   `tfsdk:"tls"`
	HostRewrite types.String `tfsdk:"host_rewrite"`
	PathRewrite types.String `tfsdk:"path_rewrite"`
}

// InitTaskModel describes an init task.
type InitTaskModel struct {
	Name           types.String `tfsdk:"name"`
	Machine        types.String `tfsdk:"machine"`
	Init           types.Bool   `tfsdk:"init"`
	User           types.String `tfsdk:"user"`
	TimeoutSeconds types.Int64  `tfsdk:"timeout_seconds"`
	Needs          types.List   `tfsdk:"needs"`
	Run            types.String `tfsdk:"run"`
}

// InitConditionsBlockModel describes the init_conditions block.
type InitConditionsBlockModel struct {
	Values []InitConditionValueModel `tfsdk:"value"`
}

// InitConditionValueModel describes an init condition value.
type InitConditionValueModel struct {
	Key      types.String `tfsdk:"key"`
	Default  types.String `tfsdk:"default"`
	Nullable types.Bool   `tfsdk:"nullable"`
	Options  types.List   `tfsdk:"options"`
}

// AccessControlModel describes access control settings.
type AccessControlModel struct {
	CanList  types.List `tfsdk:"can_list"`
	CanRead  types.List `tfsdk:"can_read"`
	CanStart types.List `tfsdk:"can_start"`
}

// PortForwardModel describes a port forward.
type PortForwardModel struct {
	Kind       types.String `tfsdk:"kind"`
	Machine    types.String `tfsdk:"machine"`
	LocalHost  types.String `tfsdk:"local_host"`
	LocalPort  types.Int64  `tfsdk:"local_port"`
	RemotePort types.Int64  `tfsdk:"remote_port"`
	RemoteHost types.String `tfsdk:"remote_host"`
}

func NewPlaygroundResource() resource.Resource {
	return &PlaygroundResource{}
}

func (r *PlaygroundResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playground"
}

func (r *PlaygroundResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an iximiuz Labs playground definition.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Server-assigned playground ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name_prefix": schema.StringAttribute{
				MarkdownDescription: "A name prefix for the playground. The API will append a unique suffix to produce the final `name`.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The playground slug identifier, assigned by the server.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"base": schema.StringAttribute{
				MarkdownDescription: "Parent playground to extend.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Display title.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Short description.",
				Required:            true,
			},
			"categories": schema.ListAttribute{
				MarkdownDescription: "Category tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"cover": schema.StringAttribute{
				MarkdownDescription: "Cover image URL.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"markdown": schema.StringAttribute{
				MarkdownDescription: "Long-form markdown description.",
				Optional:            true,
			},
			"published": schema.BoolAttribute{
				MarkdownDescription: "Whether the playground is publicly listed.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"page_url": schema.StringAttribute{
				MarkdownDescription: "Web UI URL.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "Owner identifier.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"registry_auth": schema.StringAttribute{
				MarkdownDescription: "Docker registry auth.",
				Optional:            true,
				Sensitive:           true,
			},
		},

		Blocks: map[string]schema.Block{
			"network": schema.ListNestedBlock{
				MarkdownDescription: "Network definitions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Network name.",
							Required:            true,
						},
						"subnet": schema.StringAttribute{
							MarkdownDescription: "Subnet CIDR.",
							Optional:            true,
						},
						"gateway": schema.StringAttribute{
							MarkdownDescription: "Gateway address.",
							Optional:            true,
						},
						"private": schema.BoolAttribute{
							MarkdownDescription: "Whether the network is private.",
							Optional:            true,
						},
					},
				},
			},

			"machine": schema.ListNestedBlock{
				MarkdownDescription: "Machine definitions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Machine name.",
							Required:            true,
						},
						"no_ssh": schema.BoolAttribute{
							MarkdownDescription: "Disable SSH access.",
							Optional:            true,
						},
					},
					Blocks: map[string]schema.Block{
						"user": schema.ListNestedBlock{
							MarkdownDescription: "Users on the machine.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "User name.",
										Required:            true,
									},
									"default": schema.BoolAttribute{
										MarkdownDescription: "Whether this is the default user.",
										Optional:            true,
									},
									"welcome": schema.StringAttribute{
										MarkdownDescription: "Welcome message.",
										Optional:            true,
									},
								},
							},
						},
						"kernel": schema.SingleNestedBlock{
							MarkdownDescription: "Kernel configuration.",
							Attributes: map[string]schema.Attribute{
								"source": schema.StringAttribute{
									MarkdownDescription: "Kernel source.",
									Optional:            true,
								},
							},
						},
						"drive": schema.ListNestedBlock{
							MarkdownDescription: "Drives attached to the machine.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"source": schema.StringAttribute{
										MarkdownDescription: "Drive source.",
										Optional:            true,
									},
									"mount": schema.StringAttribute{
										MarkdownDescription: "Mount point.",
										Optional:            true,
									},
									"size": schema.StringAttribute{
										MarkdownDescription: "Drive size.",
										Optional:            true,
									},
									"filesystem": schema.StringAttribute{
										MarkdownDescription: "Filesystem type.",
										Optional:            true,
									},
									"read_only": schema.BoolAttribute{
										MarkdownDescription: "Whether the drive is read-only.",
										Optional:            true,
									},
									"persistent": schema.BoolAttribute{
										MarkdownDescription: "Whether the drive is persistent.",
										Optional:            true,
									},
								},
							},
						},
						"network": schema.SingleNestedBlock{
							MarkdownDescription: "Network configuration for the machine.",
							Blocks: map[string]schema.Block{
								"interface": schema.ListNestedBlock{
									MarkdownDescription: "Network interfaces.",
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"address": schema.StringAttribute{
												MarkdownDescription: "Interface address.",
												Required:            true,
											},
											"network": schema.StringAttribute{
												MarkdownDescription: "Network name.",
												Required:            true,
											},
										},
									},
								},
							},
						},
						"resources": schema.SingleNestedBlock{
							MarkdownDescription: "Resource limits.",
							Attributes: map[string]schema.Attribute{
								"cpu_count": schema.Int64Attribute{
									MarkdownDescription: "Number of CPUs.",
									Optional:            true,
								},
								"ram_size": schema.StringAttribute{
									MarkdownDescription: "RAM size.",
									Optional:            true,
								},
							},
						},
						"startup_file": schema.ListNestedBlock{
							MarkdownDescription: "Files created at startup.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"path": schema.StringAttribute{
										MarkdownDescription: "File path.",
										Required:            true,
									},
									"content": schema.StringAttribute{
										MarkdownDescription: "File content.",
										Required:            true,
									},
									"mode": schema.StringAttribute{
										MarkdownDescription: "File mode.",
										Optional:            true,
									},
									"owner": schema.StringAttribute{
										MarkdownDescription: "File owner.",
										Optional:            true,
									},
								},
							},
						},
					},
				},
			},

			"tab": schema.ListNestedBlock{
				MarkdownDescription: "Tab definitions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"kind": schema.StringAttribute{
							MarkdownDescription: "Tab kind.",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Tab name.",
							Required:            true,
						},
						"machine": schema.StringAttribute{
							MarkdownDescription: "Machine name.",
							Optional:            true,
						},
						"number": schema.Int64Attribute{
							MarkdownDescription: "Port number.",
							Optional:            true,
						},
						"access": schema.StringAttribute{
							MarkdownDescription: "Access mode.",
							Optional:            true,
						},
						"tls": schema.BoolAttribute{
							MarkdownDescription: "Enable TLS.",
							Optional:            true,
						},
						"host_rewrite": schema.StringAttribute{
							MarkdownDescription: "Host rewrite rule.",
							Optional:            true,
						},
						"path_rewrite": schema.StringAttribute{
							MarkdownDescription: "Path rewrite rule.",
							Optional:            true,
						},
					},
				},
			},

			"init_task": schema.SetNestedBlock{
				MarkdownDescription: "Initialization tasks.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Task name.",
							Required:            true,
						},
						"machine": schema.StringAttribute{
							MarkdownDescription: "Machine name.",
							Optional:            true,
						},
						"init": schema.BoolAttribute{
							MarkdownDescription: "Whether this is an init task.",
							Optional:            true,
						},
						"user": schema.StringAttribute{
							MarkdownDescription: "User to run as.",
							Optional:            true,
						},
						"timeout_seconds": schema.Int64Attribute{
							MarkdownDescription: "Timeout in seconds.",
							Optional:            true,
						},
						"needs": schema.ListAttribute{
							MarkdownDescription: "Task dependencies.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"run": schema.StringAttribute{
							MarkdownDescription: "Script to run.",
							Optional:            true,
						},
					},
				},
			},

			"init_conditions": schema.SingleNestedBlock{
				MarkdownDescription: "Initialization conditions.",
				Blocks: map[string]schema.Block{
					"value": schema.ListNestedBlock{
						MarkdownDescription: "Condition values.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									MarkdownDescription: "Condition key.",
									Required:            true,
								},
								"default": schema.StringAttribute{
									MarkdownDescription: "Default value.",
									Optional:            true,
								},
								"nullable": schema.BoolAttribute{
									MarkdownDescription: "Whether the value can be null.",
									Optional:            true,
								},
								"options": schema.ListAttribute{
									MarkdownDescription: "Allowed values.",
									Optional:            true,
									ElementType:         types.StringType,
								},
							},
						},
					},
				},
			},

			"access_control": schema.SingleNestedBlock{
				MarkdownDescription: "Access control settings.",
				Attributes: map[string]schema.Attribute{
					"can_list": schema.ListAttribute{
						MarkdownDescription: "Users who can list.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"can_read": schema.ListAttribute{
						MarkdownDescription: "Users who can read.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"can_start": schema.ListAttribute{
						MarkdownDescription: "Users who can start.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},

			"port_forward": schema.ListNestedBlock{
				MarkdownDescription: "Port forward definitions.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"kind": schema.StringAttribute{
							MarkdownDescription: "Port forward kind.",
							Required:            true,
						},
						"machine": schema.StringAttribute{
							MarkdownDescription: "Machine name.",
							Required:            true,
						},
						"local_host": schema.StringAttribute{
							MarkdownDescription: "Local host.",
							Optional:            true,
						},
						"local_port": schema.Int64Attribute{
							MarkdownDescription: "Local port.",
							Optional:            true,
						},
						"remote_port": schema.Int64Attribute{
							MarkdownDescription: "Remote port.",
							Optional:            true,
						},
						"remote_host": schema.StringAttribute{
							MarkdownDescription: "Remote host.",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *PlaygroundResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.Client, got unexpected type. Please report this issue to the provider developers.",
		)
		return
	}

	r.client = c
}

func (r *PlaygroundResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PlaygroundResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating playground resource", map[string]interface{}{
		"name_prefix": data.NamePrefix.ValueString(),
	})

	createReq := playgroundModelToCreateRequest(ctx, &data)

	pg, err := r.client.CreatePlayground(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to create playground: "+err.Error())
		return
	}

	result := playgroundResponseToModel(ctx, pg)

	// Preserve fields not returned by the API.
	result.NamePrefix = data.NamePrefix
	result.RegistryAuth = data.RegistryAuth

	// Suppress server-returned inherited fields that were not in the config.
	// Blocks (network, machine, tab, access_control) and categories may be
	// populated by the API when extending a base playground, but Terraform
	// planned them as nil/null. Setting them here would cause an
	// "inconsistent result after apply" error.
	suppressInheritedFields(&result, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)

	tflog.Trace(ctx, "created playground resource", map[string]interface{}{
		"id":   result.ID.ValueString(),
		"name": result.Name.ValueString(),
	})
}

func (r *PlaygroundResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PlaygroundResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading playground resource", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	pg, err := r.client.GetPlayground(ctx, data.Name.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "playground not found, removing from state", map[string]interface{}{
				"name": data.Name.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", "Unable to read playground: "+err.Error())
		return
	}

	result := playgroundResponseToModel(ctx, pg)

	// Preserve fields not returned by the API.
	result.NamePrefix = data.NamePrefix
	result.RegistryAuth = data.RegistryAuth

	// Suppress server-returned inherited fields that were not in state.
	suppressInheritedFields(&result, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)

	tflog.Trace(ctx, "read playground resource", map[string]interface{}{
		"id":   result.ID.ValueString(),
		"name": result.Name.ValueString(),
	})
}

func (r *PlaygroundResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PlaygroundResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating playground resource", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	updateReq := playgroundModelToUpdateRequest(ctx, &data)

	pg, err := r.client.UpdatePlayground(ctx, data.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Unable to update playground: "+err.Error())
		return
	}

	result := playgroundResponseToModel(ctx, pg)

	// Preserve fields not returned by the API.
	result.NamePrefix = data.NamePrefix
	result.RegistryAuth = data.RegistryAuth

	// Suppress server-returned inherited fields that were not in the plan.
	suppressInheritedFields(&result, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)

	tflog.Trace(ctx, "updated playground resource", map[string]interface{}{
		"id":   result.ID.ValueString(),
		"name": result.Name.ValueString(),
	})
}

func (r *PlaygroundResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PlaygroundResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "deleting playground resource", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	err := r.client.DeletePlayground(ctx, data.Name.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			tflog.Trace(ctx, "playground already deleted or not found", map[string]interface{}{
				"name": data.Name.ValueString(),
			})
			return
		}
		resp.Diagnostics.AddError("Client Error", "Unable to delete playground: "+err.Error())
		return
	}

	tflog.Trace(ctx, "deleted playground resource", map[string]interface{}{
		"name": data.Name.ValueString(),
	})
}

func (r *PlaygroundResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// --------------------------------------------------------------------------
// Inherited field suppression
// --------------------------------------------------------------------------

// suppressInheritedFields keeps server-populated block fields (networks,
// machines, tabs, access_control) and the categories attribute out of the
// result when the user did not specify them in config/state.  The API always
// returns the full playground definition including inherited base fields,
// but Terraform planned those as nil/null.  Populating them would trigger
// an "inconsistent result after apply" error because blocks cannot be
// Computed.
func suppressInheritedFields(result, prior *PlaygroundResourceModel) {
	if prior.Networks == nil {
		result.Networks = nil
	}
	if prior.Machines == nil {
		result.Machines = nil
	}
	if prior.Tabs == nil {
		result.Tabs = nil
	}
	if prior.AccessControl == nil {
		result.AccessControl = nil
	}
	// Note: categories is Optional+Computed (a ListAttribute), so the
	// framework handles it correctly — no suppression needed.
}

// --------------------------------------------------------------------------
// Conversion helpers
// --------------------------------------------------------------------------

func playgroundModelToCreateRequest(ctx context.Context, m *PlaygroundResourceModel) client.CreatePlaygroundRequest {
	req := client.CreatePlaygroundRequest{
		Name:        m.NamePrefix.ValueString(),
		Title:       m.Title.ValueString(),
		Description: m.Description.ValueString(),
	}

	if !m.Base.IsNull() && !m.Base.IsUnknown() {
		req.Base = m.Base.ValueString()
	}

	if !m.Markdown.IsNull() && !m.Markdown.IsUnknown() {
		req.Markdown = m.Markdown.ValueString()
	}

	if !m.RegistryAuth.IsNull() && !m.RegistryAuth.IsUnknown() {
		req.RegistryAuth = m.RegistryAuth.ValueString()
	}

	if !m.Categories.IsNull() && !m.Categories.IsUnknown() {
		var cats []string
		m.Categories.ElementsAs(ctx, &cats, false)
		req.Categories = cats
	}

	req.Networks = convertNetworkModelsToClient(m.Networks)
	req.Machines = convertMachineModelsToClient(ctx, m.Machines)
	req.Tabs = convertTabModelsToClient(m.Tabs)
	req.InitTasks = convertInitTaskModelsToClient(ctx, m.InitTasks)
	req.PortForwards = convertPortForwardModelsToClient(m.PortForwards)

	if m.InitConds != nil {
		req.InitConditions = convertInitConditionsModelToClient(ctx, m.InitConds)
	}

	if m.AccessControl != nil {
		req.AccessControl = convertAccessControlModelToClient(ctx, m.AccessControl)
	}

	return req
}

func playgroundModelToUpdateRequest(ctx context.Context, m *PlaygroundResourceModel) client.UpdatePlaygroundRequest {
	req := client.UpdatePlaygroundRequest{
		Title:       m.Title.ValueString(),
		Description: m.Description.ValueString(),
	}

	if !m.Cover.IsNull() && !m.Cover.IsUnknown() {
		req.Cover = m.Cover.ValueString()
	}

	if !m.Markdown.IsNull() && !m.Markdown.IsUnknown() {
		req.Markdown = m.Markdown.ValueString()
	}

	if !m.RegistryAuth.IsNull() && !m.RegistryAuth.IsUnknown() {
		req.RegistryAuth = m.RegistryAuth.ValueString()
	}

	if !m.Categories.IsNull() && !m.Categories.IsUnknown() {
		var cats []string
		m.Categories.ElementsAs(ctx, &cats, false)
		req.Categories = cats
	}

	req.Networks = convertNetworkModelsToClient(m.Networks)
	req.Machines = convertMachineModelsToClient(ctx, m.Machines)
	req.Tabs = convertTabModelsToClient(m.Tabs)
	req.InitTasks = convertInitTaskModelsToClient(ctx, m.InitTasks)
	req.PortForwards = convertPortForwardModelsToClient(m.PortForwards)

	if m.InitConds != nil {
		req.InitConditions = convertInitConditionsModelToClient(ctx, m.InitConds)
	}

	if m.AccessControl != nil {
		req.AccessControl = convertAccessControlModelToClient(ctx, m.AccessControl)
	}

	return req
}

func playgroundResponseToModel(ctx context.Context, pg *client.Playground) PlaygroundResourceModel {
	m := PlaygroundResourceModel{
		ID:          types.StringValue(pg.ID),
		Name:        types.StringValue(pg.Name),
		Title:       types.StringValue(pg.Title),
		Description: types.StringValue(pg.Description),
		Published:   types.BoolValue(pg.Published),
		PageURL:     types.StringValue(pg.PageURL),
		Owner:       types.StringValue(pg.Owner),
	}

	if pg.Base != "" {
		m.Base = types.StringValue(pg.Base)
	} else {
		m.Base = types.StringNull()
	}

	if pg.Cover != "" {
		m.Cover = types.StringValue(pg.Cover)
	} else {
		m.Cover = types.StringNull()
	}

	if pg.Markdown != "" {
		m.Markdown = types.StringValue(pg.Markdown)
	} else {
		m.Markdown = types.StringNull()
	}

	if pg.RegistryAuth != "" {
		m.RegistryAuth = types.StringValue(pg.RegistryAuth)
	} else {
		m.RegistryAuth = types.StringNull()
	}

	if len(pg.Categories) > 0 {
		catList, _ := types.ListValueFrom(ctx, types.StringType, pg.Categories)
		m.Categories = catList
	} else {
		m.Categories = types.ListNull(types.StringType)
	}

	m.Networks = convertClientNetworksToModel(pg.Networks)
	m.Machines = convertClientMachinesToModel(ctx, pg.Machines)
	m.Tabs = convertClientTabsToModel(pg.Tabs)
	m.InitTasks = convertClientInitTasksToModel(ctx, pg.InitTasks)
	m.PortForwards = convertClientPortForwardsToModel(pg.PortForwards)

	if pg.InitConditions != nil {
		m.InitConds = convertClientInitConditionsToModel(ctx, pg.InitConditions)
	}

	if pg.AccessControl != nil {
		m.AccessControl = convertClientAccessControlToModel(ctx, pg.AccessControl)
	}

	return m
}

// --------------------------------------------------------------------------
// Network conversion
// --------------------------------------------------------------------------

func convertNetworkModelsToClient(models []NetworkModel) []client.PlaygroundNetwork {
	if len(models) == 0 {
		return nil
	}
	result := make([]client.PlaygroundNetwork, len(models))
	for i, n := range models {
		result[i] = client.PlaygroundNetwork{
			Name: n.Name.ValueString(),
		}
		if !n.Subnet.IsNull() && !n.Subnet.IsUnknown() {
			result[i].Subnet = n.Subnet.ValueString()
		}
		if !n.Gateway.IsNull() && !n.Gateway.IsUnknown() {
			result[i].Gateway = n.Gateway.ValueString()
		}
		if !n.Private.IsNull() && !n.Private.IsUnknown() {
			result[i].Private = n.Private.ValueBool()
		}
	}
	return result
}

func convertClientNetworksToModel(networks []client.PlaygroundNetwork) []NetworkModel {
	if len(networks) == 0 {
		return nil
	}
	result := make([]NetworkModel, len(networks))
	for i, n := range networks {
		result[i] = NetworkModel{
			Name: types.StringValue(n.Name),
		}
		if n.Subnet != "" {
			result[i].Subnet = types.StringValue(n.Subnet)
		} else {
			result[i].Subnet = types.StringNull()
		}
		if n.Gateway != "" {
			result[i].Gateway = types.StringValue(n.Gateway)
		} else {
			result[i].Gateway = types.StringNull()
		}
		if n.Private {
			result[i].Private = types.BoolValue(true)
		} else {
			result[i].Private = types.BoolNull()
		}
	}
	return result
}

// --------------------------------------------------------------------------
// Machine conversion
// --------------------------------------------------------------------------

func convertMachineModelsToClient(_ context.Context, models []MachineModel) []client.PlaygroundMachine {
	if len(models) == 0 {
		return nil
	}
	result := make([]client.PlaygroundMachine, len(models))
	for i, m := range models {
		cm := client.PlaygroundMachine{
			Name: m.Name.ValueString(),
		}
		if !m.NoSSH.IsNull() && !m.NoSSH.IsUnknown() {
			cm.NoSSH = m.NoSSH.ValueBool()
		}

		// Users
		if len(m.Users) > 0 {
			users := make([]client.MachineUser, len(m.Users))
			for j, u := range m.Users {
				users[j] = client.MachineUser{
					Name: u.Name.ValueString(),
				}
				if !u.Default.IsNull() && !u.Default.IsUnknown() {
					users[j].Default = u.Default.ValueBool()
				}
				if !u.Welcome.IsNull() && !u.Welcome.IsUnknown() {
					users[j].Welcome = u.Welcome.ValueString()
				}
			}
			cm.Users = users
		}

		// Kernel
		if m.Kernel != nil {
			k := &client.MachineKernel{}
			if !m.Kernel.Source.IsNull() && !m.Kernel.Source.IsUnknown() {
				k.Source = m.Kernel.Source.ValueString()
			}
			cm.Kernel = k
		}

		// Drives
		if len(m.Drives) > 0 {
			drives := make([]client.MachineDrive, len(m.Drives))
			for j, d := range m.Drives {
				drives[j] = client.MachineDrive{}
				if !d.Source.IsNull() && !d.Source.IsUnknown() {
					drives[j].Source = d.Source.ValueString()
				}
				if !d.Mount.IsNull() && !d.Mount.IsUnknown() {
					drives[j].Mount = d.Mount.ValueString()
				}
				if !d.Size.IsNull() && !d.Size.IsUnknown() {
					drives[j].Size = d.Size.ValueString()
				}
				if !d.Filesystem.IsNull() && !d.Filesystem.IsUnknown() {
					drives[j].Filesystem = d.Filesystem.ValueString()
				}
				if !d.ReadOnly.IsNull() && !d.ReadOnly.IsUnknown() {
					drives[j].ReadOnly = d.ReadOnly.ValueBool()
				}
				if !d.Persistent.IsNull() && !d.Persistent.IsUnknown() {
					drives[j].Persistent = d.Persistent.ValueBool()
				}
			}
			cm.Drives = drives
		}

		// Network
		if m.Network != nil && len(m.Network.Interfaces) > 0 {
			ifaces := make([]client.MachineNetworkInterface, len(m.Network.Interfaces))
			for j, iface := range m.Network.Interfaces {
				ifaces[j] = client.MachineNetworkInterface{
					Address: iface.Address.ValueString(),
					Network: iface.Network.ValueString(),
				}
			}
			cm.Network = &client.MachineNetwork{
				Interfaces: ifaces,
			}
		}

		// Resources
		if m.Resources != nil {
			res := &client.MachineResources{}
			if !m.Resources.CPUCount.IsNull() && !m.Resources.CPUCount.IsUnknown() {
				res.CPUCount = int(m.Resources.CPUCount.ValueInt64())
			}
			if !m.Resources.RAMSize.IsNull() && !m.Resources.RAMSize.IsUnknown() {
				res.RAMSize = m.Resources.RAMSize.ValueString()
			}
			cm.Resources = res
		}

		// Startup files
		if len(m.StartupFiles) > 0 {
			files := make([]client.MachineStartupFile, len(m.StartupFiles))
			for j, f := range m.StartupFiles {
				files[j] = client.MachineStartupFile{
					Path:    f.Path.ValueString(),
					Content: f.Content.ValueString(),
				}
				if !f.Mode.IsNull() && !f.Mode.IsUnknown() {
					files[j].Mode = f.Mode.ValueString()
				}
				if !f.Owner.IsNull() && !f.Owner.IsUnknown() {
					files[j].Owner = f.Owner.ValueString()
				}
			}
			cm.StartupFiles = files
		}

		result[i] = cm
	}
	return result
}

func convertClientMachinesToModel(_ context.Context, machines []client.PlaygroundMachine) []MachineModel {
	if len(machines) == 0 {
		return nil
	}
	result := make([]MachineModel, len(machines))
	for i, m := range machines {
		mm := MachineModel{
			Name: types.StringValue(m.Name),
		}
		if m.NoSSH {
			mm.NoSSH = types.BoolValue(true)
		} else {
			mm.NoSSH = types.BoolNull()
		}

		// Users
		if len(m.Users) > 0 {
			users := make([]MachineUserModel, len(m.Users))
			for j, u := range m.Users {
				users[j] = MachineUserModel{
					Name: types.StringValue(u.Name),
				}
				if u.Default {
					users[j].Default = types.BoolValue(true)
				} else {
					users[j].Default = types.BoolNull()
				}
				if u.Welcome != "" {
					users[j].Welcome = types.StringValue(u.Welcome)
				} else {
					users[j].Welcome = types.StringNull()
				}
			}
			mm.Users = users
		}

		// Kernel
		if m.Kernel != nil {
			k := &MachineKernelModel{}
			if m.Kernel.Source != "" {
				k.Source = types.StringValue(m.Kernel.Source)
			} else {
				k.Source = types.StringNull()
			}
			mm.Kernel = k
		}

		// Drives
		if len(m.Drives) > 0 {
			drives := make([]MachineDriveModel, len(m.Drives))
			for j, d := range m.Drives {
				dm := MachineDriveModel{}
				if d.Source != "" {
					dm.Source = types.StringValue(d.Source)
				} else {
					dm.Source = types.StringNull()
				}
				if d.Mount != "" {
					dm.Mount = types.StringValue(d.Mount)
				} else {
					dm.Mount = types.StringNull()
				}
				if d.Size != "" {
					dm.Size = types.StringValue(d.Size)
				} else {
					dm.Size = types.StringNull()
				}
				if d.Filesystem != "" {
					dm.Filesystem = types.StringValue(d.Filesystem)
				} else {
					dm.Filesystem = types.StringNull()
				}
				if d.ReadOnly {
					dm.ReadOnly = types.BoolValue(true)
				} else {
					dm.ReadOnly = types.BoolNull()
				}
				if d.Persistent {
					dm.Persistent = types.BoolValue(true)
				} else {
					dm.Persistent = types.BoolNull()
				}
				drives[j] = dm
			}
			mm.Drives = drives
		}

		// Network
		if m.Network != nil && len(m.Network.Interfaces) > 0 {
			ifaces := make([]MachineNetworkInterfaceModel, len(m.Network.Interfaces))
			for j, iface := range m.Network.Interfaces {
				ifaces[j] = MachineNetworkInterfaceModel{
					Address: types.StringValue(iface.Address),
					Network: types.StringValue(iface.Network),
				}
			}
			mm.Network = &MachineNetworkModel{
				Interfaces: ifaces,
			}
		}

		// Resources
		if m.Resources != nil {
			res := &MachineResourceModel{}
			if m.Resources.CPUCount > 0 {
				res.CPUCount = types.Int64Value(int64(m.Resources.CPUCount))
			} else {
				res.CPUCount = types.Int64Null()
			}
			if m.Resources.RAMSize != "" {
				res.RAMSize = types.StringValue(m.Resources.RAMSize)
			} else {
				res.RAMSize = types.StringNull()
			}
			mm.Resources = res
		}

		// Startup files
		if len(m.StartupFiles) > 0 {
			files := make([]StartupFileModel, len(m.StartupFiles))
			for j, f := range m.StartupFiles {
				fm := StartupFileModel{
					Path:    types.StringValue(f.Path),
					Content: types.StringValue(f.Content),
				}
				if f.Mode != "" {
					fm.Mode = types.StringValue(f.Mode)
				} else {
					fm.Mode = types.StringNull()
				}
				if f.Owner != "" {
					fm.Owner = types.StringValue(f.Owner)
				} else {
					fm.Owner = types.StringNull()
				}
				files[j] = fm
			}
			mm.StartupFiles = files
		}

		result[i] = mm
	}
	return result
}

// --------------------------------------------------------------------------
// Tab conversion
// --------------------------------------------------------------------------

func convertTabModelsToClient(models []TabModel) []client.PlaygroundTab {
	if len(models) == 0 {
		return nil
	}
	result := make([]client.PlaygroundTab, len(models))
	for i, t := range models {
		ct := client.PlaygroundTab{
			Kind: t.Kind.ValueString(),
			Name: t.Name.ValueString(),
		}
		if !t.Machine.IsNull() && !t.Machine.IsUnknown() {
			ct.Machine = t.Machine.ValueString()
		}
		if !t.Number.IsNull() && !t.Number.IsUnknown() {
			ct.Number = int(t.Number.ValueInt64())
		}
		if !t.Access.IsNull() && !t.Access.IsUnknown() {
			ct.Access = t.Access.ValueString()
		}
		if !t.TLS.IsNull() && !t.TLS.IsUnknown() {
			ct.TLS = t.TLS.ValueBool()
		}
		if !t.HostRewrite.IsNull() && !t.HostRewrite.IsUnknown() {
			ct.HostRewrite = t.HostRewrite.ValueString()
		}
		if !t.PathRewrite.IsNull() && !t.PathRewrite.IsUnknown() {
			ct.PathRewrite = t.PathRewrite.ValueString()
		}
		result[i] = ct
	}
	return result
}

func convertClientTabsToModel(tabs []client.PlaygroundTab) []TabModel {
	if len(tabs) == 0 {
		return nil
	}
	result := make([]TabModel, len(tabs))
	for i, t := range tabs {
		tm := TabModel{
			Kind: types.StringValue(t.Kind),
			Name: types.StringValue(t.Name),
		}
		if t.Machine != "" {
			tm.Machine = types.StringValue(t.Machine)
		} else {
			tm.Machine = types.StringNull()
		}
		if t.Number != 0 {
			tm.Number = types.Int64Value(int64(t.Number))
		} else {
			tm.Number = types.Int64Null()
		}
		if t.Access != "" {
			tm.Access = types.StringValue(t.Access)
		} else {
			tm.Access = types.StringNull()
		}
		if t.TLS {
			tm.TLS = types.BoolValue(true)
		} else {
			tm.TLS = types.BoolNull()
		}
		if t.HostRewrite != "" {
			tm.HostRewrite = types.StringValue(t.HostRewrite)
		} else {
			tm.HostRewrite = types.StringNull()
		}
		if t.PathRewrite != "" {
			tm.PathRewrite = types.StringValue(t.PathRewrite)
		} else {
			tm.PathRewrite = types.StringNull()
		}
		result[i] = tm
	}
	return result
}

// --------------------------------------------------------------------------
// InitTask conversion
// --------------------------------------------------------------------------

func convertInitTaskModelsToClient(ctx context.Context, models []InitTaskModel) map[string]client.InitTask {
	if len(models) == 0 {
		return nil
	}
	result := make(map[string]client.InitTask, len(models))
	for _, t := range models {
		ct := client.InitTask{
			Name: t.Name.ValueString(),
		}
		if !t.Machine.IsNull() && !t.Machine.IsUnknown() {
			ct.Machine = t.Machine.ValueString()
		}
		if !t.Init.IsNull() && !t.Init.IsUnknown() {
			ct.Init = t.Init.ValueBool()
		}
		if !t.User.IsNull() && !t.User.IsUnknown() {
			ct.User = t.User.ValueString()
		}
		if !t.TimeoutSeconds.IsNull() && !t.TimeoutSeconds.IsUnknown() {
			ct.TimeoutSeconds = int(t.TimeoutSeconds.ValueInt64())
		}
		if !t.Needs.IsNull() && !t.Needs.IsUnknown() {
			var needs []string
			t.Needs.ElementsAs(ctx, &needs, false)
			ct.Needs = needs
		}
		if !t.Run.IsNull() && !t.Run.IsUnknown() {
			ct.Run = t.Run.ValueString()
		}
		result[ct.Name] = ct
	}
	return result
}

func convertClientInitTasksToModel(ctx context.Context, tasks map[string]client.InitTask) []InitTaskModel {
	if len(tasks) == 0 {
		return nil
	}
	result := make([]InitTaskModel, 0, len(tasks))
	for _, t := range tasks {
		tm := InitTaskModel{
			Name: types.StringValue(t.Name),
		}
		if t.Machine != "" {
			tm.Machine = types.StringValue(t.Machine)
		} else {
			tm.Machine = types.StringNull()
		}
		if t.Init {
			tm.Init = types.BoolValue(true)
		} else {
			tm.Init = types.BoolNull()
		}
		if t.User != "" {
			tm.User = types.StringValue(t.User)
		} else {
			tm.User = types.StringNull()
		}
		if t.TimeoutSeconds != 0 {
			tm.TimeoutSeconds = types.Int64Value(int64(t.TimeoutSeconds))
		} else {
			tm.TimeoutSeconds = types.Int64Null()
		}
		if len(t.Needs) > 0 {
			needsList, _ := types.ListValueFrom(ctx, types.StringType, t.Needs)
			tm.Needs = needsList
		} else {
			tm.Needs = types.ListNull(types.StringType)
		}
		if t.Run != "" {
			tm.Run = types.StringValue(t.Run)
		} else {
			tm.Run = types.StringNull()
		}
		result = append(result, tm)
	}
	return result
}

// --------------------------------------------------------------------------
// InitConditions conversion
// --------------------------------------------------------------------------

func convertInitConditionsModelToClient(ctx context.Context, m *InitConditionsBlockModel) *client.InitConditions {
	if m == nil {
		return nil
	}
	ic := &client.InitConditions{}
	if len(m.Values) > 0 {
		vals := make([]client.InitConditionValue, len(m.Values))
		for i, v := range m.Values {
			vals[i] = client.InitConditionValue{
				Key: v.Key.ValueString(),
			}
			if !v.Default.IsNull() && !v.Default.IsUnknown() {
				vals[i].Default = v.Default.ValueString()
			}
			if !v.Nullable.IsNull() && !v.Nullable.IsUnknown() {
				vals[i].Nullable = v.Nullable.ValueBool()
			}
			if !v.Options.IsNull() && !v.Options.IsUnknown() {
				var opts []string
				v.Options.ElementsAs(ctx, &opts, false)
				vals[i].Options = opts
			}
		}
		ic.Values = vals
	}
	return ic
}

func convertClientInitConditionsToModel(ctx context.Context, ic *client.InitConditions) *InitConditionsBlockModel {
	if ic == nil {
		return nil
	}
	m := &InitConditionsBlockModel{}
	if len(ic.Values) > 0 {
		vals := make([]InitConditionValueModel, len(ic.Values))
		for i, v := range ic.Values {
			vm := InitConditionValueModel{
				Key: types.StringValue(v.Key),
			}
			if v.Default != "" {
				vm.Default = types.StringValue(v.Default)
			} else {
				vm.Default = types.StringNull()
			}
			if v.Nullable {
				vm.Nullable = types.BoolValue(true)
			} else {
				vm.Nullable = types.BoolNull()
			}
			if len(v.Options) > 0 {
				optsList, _ := types.ListValueFrom(ctx, types.StringType, v.Options)
				vm.Options = optsList
			} else {
				vm.Options = types.ListNull(types.StringType)
			}
			vals[i] = vm
		}
		m.Values = vals
	}
	return m
}

// --------------------------------------------------------------------------
// AccessControl conversion
// --------------------------------------------------------------------------

func convertAccessControlModelToClient(ctx context.Context, m *AccessControlModel) *client.PlaygroundAccessControl {
	if m == nil {
		return nil
	}
	ac := &client.PlaygroundAccessControl{}
	if !m.CanList.IsNull() && !m.CanList.IsUnknown() {
		var vals []string
		m.CanList.ElementsAs(ctx, &vals, false)
		ac.CanList = vals
	}
	if !m.CanRead.IsNull() && !m.CanRead.IsUnknown() {
		var vals []string
		m.CanRead.ElementsAs(ctx, &vals, false)
		ac.CanRead = vals
	}
	if !m.CanStart.IsNull() && !m.CanStart.IsUnknown() {
		var vals []string
		m.CanStart.ElementsAs(ctx, &vals, false)
		ac.CanStart = vals
	}
	return ac
}

func convertClientAccessControlToModel(ctx context.Context, ac *client.PlaygroundAccessControl) *AccessControlModel {
	if ac == nil {
		return nil
	}
	m := &AccessControlModel{}
	if len(ac.CanList) > 0 {
		list, _ := types.ListValueFrom(ctx, types.StringType, ac.CanList)
		m.CanList = list
	} else {
		m.CanList = types.ListNull(types.StringType)
	}
	if len(ac.CanRead) > 0 {
		list, _ := types.ListValueFrom(ctx, types.StringType, ac.CanRead)
		m.CanRead = list
	} else {
		m.CanRead = types.ListNull(types.StringType)
	}
	if len(ac.CanStart) > 0 {
		list, _ := types.ListValueFrom(ctx, types.StringType, ac.CanStart)
		m.CanStart = list
	} else {
		m.CanStart = types.ListNull(types.StringType)
	}
	return m
}

// --------------------------------------------------------------------------
// PortForward conversion
// --------------------------------------------------------------------------

func convertPortForwardModelsToClient(models []PortForwardModel) []client.PortForward {
	if len(models) == 0 {
		return nil
	}
	result := make([]client.PortForward, len(models))
	for i, p := range models {
		pf := client.PortForward{
			Kind:    p.Kind.ValueString(),
			Machine: p.Machine.ValueString(),
		}
		if !p.LocalHost.IsNull() && !p.LocalHost.IsUnknown() {
			pf.LocalHost = p.LocalHost.ValueString()
		}
		if !p.LocalPort.IsNull() && !p.LocalPort.IsUnknown() {
			pf.LocalPort = int(p.LocalPort.ValueInt64())
		}
		if !p.RemotePort.IsNull() && !p.RemotePort.IsUnknown() {
			pf.RemotePort = int(p.RemotePort.ValueInt64())
		}
		if !p.RemoteHost.IsNull() && !p.RemoteHost.IsUnknown() {
			pf.RemoteHost = p.RemoteHost.ValueString()
		}
		result[i] = pf
	}
	return result
}

func convertClientPortForwardsToModel(pfs []client.PortForward) []PortForwardModel {
	if len(pfs) == 0 {
		return nil
	}
	result := make([]PortForwardModel, len(pfs))
	for i, p := range pfs {
		pm := PortForwardModel{
			Kind:    types.StringValue(p.Kind),
			Machine: types.StringValue(p.Machine),
		}
		if p.LocalHost != "" {
			pm.LocalHost = types.StringValue(p.LocalHost)
		} else {
			pm.LocalHost = types.StringNull()
		}
		if p.LocalPort != 0 {
			pm.LocalPort = types.Int64Value(int64(p.LocalPort))
		} else {
			pm.LocalPort = types.Int64Null()
		}
		if p.RemotePort != 0 {
			pm.RemotePort = types.Int64Value(int64(p.RemotePort))
		} else {
			pm.RemotePort = types.Int64Null()
		}
		if p.RemoteHost != "" {
			pm.RemoteHost = types.StringValue(p.RemoteHost)
		} else {
			pm.RemoteHost = types.StringNull()
		}
		result[i] = pm
	}
	return result
}
