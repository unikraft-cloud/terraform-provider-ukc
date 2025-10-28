// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/model"

	"unikraft.com/cloud/sdk/platform"
)

func NewInstanceResource() resource.Resource {
	return &InstanceResource{}
}

// InstanceResource defines the resource implementation.
type InstanceResource struct {
	client platform.Client
}

// Ensure InstanceResource satisfies various resource interfaces.
var (
	_ resource.Resource                = &InstanceResource{}
	_ resource.ResourceWithImportState = &InstanceResource{}
)

// InstanceResourceModel describes the resource data model.
type InstanceResourceModel struct {
	Image     types.String `tfsdk:"image"`
	Args      types.List   `tfsdk:"args"`
	MemoryMB  types.Int64  `tfsdk:"memory_mb"`
	Autostart types.Bool   `tfsdk:"autostart"`

	UUID              types.String        `tfsdk:"uuid"`
	Name              types.String        `tfsdk:"name"`
	FQDN              types.String        `tfsdk:"fqdn"`
	PrivateIP         types.String        `tfsdk:"private_ip"`
	PrivateFQDN       types.String        `tfsdk:"private_fqdn"`
	State             types.String        `tfsdk:"state"`
	CreatedAt         types.String        `tfsdk:"created_at"`
	Env               types.Map           `tfsdk:"env"`
	ServiceGroup      *models.SvcGrpModel `tfsdk:"service_group"`
	NetworkInterfaces types.List          `tfsdk:"network_interfaces"`
	BootTimeUS        types.Int64         `tfsdk:"boot_time_us"`
}

// Metadata implements resource.Resource.
func (r *InstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema implements resource.Resource.
func (r *InstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Allows the creation of Unikraft Cloud instances.",

		Attributes: map[string]schema.Attribute{
			"image": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"args": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"memory_mb": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(16, 256),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"autostart": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"fqdn": schema.StringAttribute{
				Computed: true,
			},
			"private_ip": schema.StringAttribute{
				Computed: true,
			},
			"private_fqdn": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"env": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"service_group": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"services": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"port": schema.Int64Attribute{
									Required: true,
									Validators: []validator.Int64{
										int64validator.Between(1, math.MaxUint16),
									},
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.RequiresReplace(),
									},
								},
								"destination_port": schema.Int64Attribute{
									Optional: true,
									Computed: true,
									Validators: []validator.Int64{
										int64validator.Between(1, math.MaxUint16),
									},
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.RequiresReplaceIfConfigured(),
										int64planmodifier.UseStateForUnknown(),
									},
								},
								"handlers": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Set{
										setplanmodifier.RequiresReplaceIfConfigured(),
										setplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"domains": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"fqdn": schema.StringAttribute{
									Computed: true,
								},
								"certificate": schema.MapNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"uuid": schema.StringAttribute{
												Computed: true,
											},
											"name": schema.StringAttribute{
												Computed: true,
											},
											"state": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"network_interfaces": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"private_ip": schema.StringAttribute{
							Computed: true,
						},
						"mac": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"boot_time_us": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

// Configure implements resource.Resource.
func (r *InstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(platform.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected platform.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create implements resource.Resource.
func (r *InstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InstanceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := platform.CreateInstanceRequest{
		Image: data.Image.ValueString(),
	}

	// New SDK properly handles optional fields with pointers
	if !data.MemoryMB.IsUnknown() && !data.MemoryMB.IsNull() {
		memoryMB := data.MemoryMB.ValueInt64()
		in.MemoryMb = &memoryMB
	}

	if !data.Autostart.IsUnknown() && !data.Autostart.IsNull() {
		autostart := data.Autostart.ValueBool()
		in.Autostart = &autostart
	}

	if !data.Args.IsNull() && !data.Args.IsUnknown() {
		argVals := make([]types.String, 0, len(data.Args.Elements()))
		resp.Diagnostics.Append(data.Args.ElementsAs(ctx, &argVals, false)...)
		for _, v := range argVals {
			in.Args = append(in.Args, v.ValueString())
		}
	}

	if data.ServiceGroup != nil && len(data.ServiceGroup.Services) > 0 {
		sgServices := make([]platform.Service, len(data.ServiceGroup.Services))
		for i, svc := range data.ServiceGroup.Services {
			port := uint32(svc.Port.ValueInt64())
			sgServices[i].Port = port

			// New SDK properly handles optional destination port with pointer
			if !svc.DestinationPort.IsUnknown() && !svc.DestinationPort.IsNull() {
				destPort := uint32(svc.DestinationPort.ValueInt64())
				sgServices[i].DestinationPort = &destPort
			}

			if !svc.Handlers.IsUnknown() {
				handlVals := make([]types.String, 0, len(svc.Handlers.Elements()))
				resp.Diagnostics.Append(svc.Handlers.ElementsAs(ctx, &handlVals, false)...)
				for _, v := range handlVals {
					sgServices[i].Handlers = append(sgServices[i].Handlers, platform.ServiceHandlers(v.ValueString()))
				}
			}
		}
		in.ServiceGroup = &platform.CreateInstanceRequestServiceGroup{
			Services: sgServices,
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	insResp, err := r.client.CreateInstance(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create instance, got error: %v", err),
		)
		return
	}

	if insResp == nil || insResp.Data == nil || len(insResp.Data.Instances) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from create instance API",
		)
		return
	}
	ins := insResp.Data.Instances[0]

	if ins.Uuid == nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Instance UUID not returned by API",
		)
		return
	}

	data.UUID = types.StringValue(*ins.Uuid)
	if ins.Name != nil {
		data.Name = types.StringValue(*ins.Name)
	}
	if ins.ServiceGroup != nil && len(ins.ServiceGroup.Domains) > 0 && ins.ServiceGroup.Domains[0].Fqdn != nil {
		data.FQDN = types.StringValue(*ins.ServiceGroup.Domains[0].Fqdn)
	}
	if ins.PrivateFqdn != nil {
		data.PrivateFQDN = types.StringValue(*ins.PrivateFqdn)
	}

	// Not all attributes are returned by CreateInstance
	insFullResp, err := r.client.GetInstanceByUUID(ctx, *ins.Uuid, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get instance state, got error: %v", err),
		)
		return
	}

	if insFullResp == nil || insFullResp.Data == nil || len(insFullResp.Data.Instances) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from get instance API",
		)
		return
	}
	insFull := insFullResp.Data.Instances[0]

	var diags diag.Diagnostics

	// NOTE(antoineco): although the Image attribute may be transformed by
	// Unikraft Cloud (e.g. replace the tag with a digest), we must not update the
	// value read from the schema, otherwise Terraform fails to apply with the
	// following error:
	//
	//   Error: Provider produced inconsistent result after apply
	//   When applying changes to unikraft-cloud_instance.xyz, provider produced an unexpected new value: .image:
	//     was cty.StringVal("myimage:latest"), but now cty.StringVal("myimage@sha256:18a381f0062...").
	//
	if insFull.State != nil {
		data.State = types.StringValue(string(*insFull.State))
	}
	if insFull.CreatedAt != nil {
		data.CreatedAt = types.StringValue(insFull.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"))
	}
	if insFull.MemoryMb != nil {
		data.MemoryMB = types.Int64Value(int64(*insFull.MemoryMb))
	}
	if insFull.BootTimeUs != nil {
		data.BootTimeUS = types.Int64Value(int64(*insFull.BootTimeUs))
	}

	if data.Args.IsNull() || data.Args.IsUnknown() {
		if insFull.Args != nil {
			data.Args, diags = types.ListValueFrom(ctx, types.StringType, insFull.Args)
			resp.Diagnostics.Append(diags...)
		}
	}

	if insFull.Env != nil {
		data.Env, diags = types.MapValueFrom(ctx, types.StringType, insFull.Env)
		resp.Diagnostics.Append(diags...)
	}

	if data.ServiceGroup == nil {
		data.ServiceGroup = &models.SvcGrpModel{}
	}

	if insFull.ServiceGroup != nil {
		if insFull.ServiceGroup.Uuid != nil {
			data.ServiceGroup.UUID = types.StringValue(*insFull.ServiceGroup.Uuid)
		}
		if insFull.ServiceGroup.Name != nil {
			data.ServiceGroup.Name = types.StringValue(*insFull.ServiceGroup.Name)
		}
	}

	if insFull.NetworkInterfaces != nil {
		netwIfaces := make([]models.NetwIfaceModel, len(insFull.NetworkInterfaces))
		for i, net := range insFull.NetworkInterfaces {
			if net.Uuid != nil {
				netwIfaces[i].UUID = types.StringValue(*net.Uuid)
				netwIfaces[i].Name = types.StringValue(*net.Uuid) // No name in the response
			}
			if net.PrivateIp != nil {
				netwIfaces[i].PrivateIP = types.StringValue(*net.PrivateIp)
			}
			if net.Mac != nil {
				netwIfaces[i].MAC = types.StringValue(*net.Mac)
			}
		}
		data.NetworkInterfaces, diags = types.ListValueFrom(ctx, models.NetwIfaceModelType, netwIfaces)
		resp.Diagnostics.Append(diags...)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource.Resource.
func (r *InstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	insResp, err := r.client.GetInstanceByUUID(ctx, data.UUID.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get instance state, got error: %v", err),
		)
		return
	}

	if insResp == nil || insResp.Data == nil || len(insResp.Data.Instances) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from get instance API",
		)
		return
	}
	ins := insResp.Data.Instances[0]

	var diags diag.Diagnostics

	// NOTE(antoineco): although the Image attribute may be transformed by
	// Unikraft Cloud (e.g. replace the tag with a digest), we must not update the
	// value read from the schema, otherwise Terraform fails to apply with the
	// following error:
	//
	//   Error: Provider produced inconsistent result after apply
	//   When applying changes to unikraft-cloud_instance.xyz, provider produced an unexpected new value: .image:
	//     was cty.StringVal("myimage:latest"), but now cty.StringVal("myimage@sha256:18a381f0062...").
	//
	// However, we must still ensure that the Image attribute is populated by
	// "terraform import".
	if data.Image.IsNull() && ins.Image != nil {
		data.Image = types.StringValue(*ins.Image)
	}
	if ins.Name != nil {
		data.Name = types.StringValue(*ins.Name)
	}
	if ins.ServiceGroup != nil && len(ins.ServiceGroup.Domains) > 0 && ins.ServiceGroup.Domains[0].Fqdn != nil {
		data.FQDN = types.StringValue(*ins.ServiceGroup.Domains[0].Fqdn)
	}
	if ins.PrivateFqdn != nil {
		data.PrivateFQDN = types.StringValue(*ins.PrivateFqdn)
	}
	if ins.State != nil {
		data.State = types.StringValue(string(*ins.State))
	}
	if ins.CreatedAt != nil {
		data.CreatedAt = types.StringValue(ins.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"))
	}
	if ins.MemoryMb != nil {
		data.MemoryMB = types.Int64Value(int64(*ins.MemoryMb))
	}
	if ins.BootTimeUs != nil {
		data.BootTimeUS = types.Int64Value(int64(*ins.BootTimeUs))
	}

	if data.Args.IsNull() || data.Args.IsUnknown() {
		if ins.Args != nil {
			data.Args, diags = types.ListValueFrom(ctx, types.StringType, ins.Args)
			resp.Diagnostics.Append(diags...)
		}
	}

	if ins.Env != nil {
		data.Env, diags = types.MapValueFrom(ctx, types.StringType, ins.Env)
		resp.Diagnostics.Append(diags...)
	}

	if data.ServiceGroup == nil {
		data.ServiceGroup = &models.SvcGrpModel{}
	}

	if ins.ServiceGroup != nil {
		if ins.ServiceGroup.Uuid != nil {
			data.ServiceGroup.UUID = types.StringValue(*ins.ServiceGroup.Uuid)
		}
		if ins.ServiceGroup.Name != nil {
			data.ServiceGroup.Name = types.StringValue(*ins.ServiceGroup.Name)
		}
	}

	if ins.NetworkInterfaces != nil {
		netwIfaces := make([]models.NetwIfaceModel, len(ins.NetworkInterfaces))
		for i, net := range ins.NetworkInterfaces {
			if net.Uuid != nil {
				netwIfaces[i].UUID = types.StringValue(*net.Uuid)
			}
			if net.Uuid != nil {
				netwIfaces[i].Name = types.StringValue(*net.Uuid) // No name in the response
			}
			if net.PrivateIp != nil {
				netwIfaces[i].PrivateIP = types.StringValue(*net.PrivateIp)
			}
			if net.Mac != nil {
				netwIfaces[i].MAC = types.StringValue(*net.Mac)
			}
		}
		data.NetworkInterfaces, diags = types.ListValueFrom(ctx, models.NetwIfaceModelType, netwIfaces)
		resp.Diagnostics.Append(diags...)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (r *InstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unsupported",
		"This resource does not support updates. Configuration changes were expected to have triggered a replacement "+
			"of the resource. Please report this issue to the provider developers.",
	)
}

// Delete implements resource.Resource.
func (r *InstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteInstanceByUUID(ctx, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete instance, got error: %v", err),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *InstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
