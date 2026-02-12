// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/model"

	"unikraft.com/cloud/sdk/platform"
)

func NewVolumeResource() resource.Resource {
	return &VolumeResource{}
}

// VolumeResource defines the resource implementation.
type VolumeResource struct {
	client platform.Client
}

// Ensure VolumeResource satisfies various resource interfaces.
var (
	_ resource.Resource                = &VolumeResource{}
	_ resource.ResourceWithImportState = &VolumeResource{}
)

// VolumeResourceModel describes the resource data model.
type VolumeResourceModel struct {
	Name       types.String `tfsdk:"name"`
	SizeMB     types.Int64  `tfsdk:"size_mb"`
	UUID       types.String `tfsdk:"uuid"`
	State      types.String `tfsdk:"state"`
	Persistent types.Bool   `tfsdk:"persistent"`
	CreatedAt  types.String `tfsdk:"created_at"`
	AttachedTo types.List   `tfsdk:"attached_to"`
}

// Metadata implements resource.Resource.
func (r *VolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema implements resource.Resource.
func (r *VolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Allows the creation of Unikraft Cloud volumes.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the volume. If not specified, a random name is generated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"size_mb": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The size of the volume in megabytes.",
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current state of the volume.",
			},
			"persistent": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the volume survives instance deletion.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The time the volume was created.",
			},
			"attached_to": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of instances this volume is attached to.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Configure implements resource.Resource.
func (r *VolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VolumeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := platform.CreateVolumeRequest{
		SizeMb: uint64(data.SizeMB.ValueInt64()),
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		name := data.Name.ValueString()
		in.Name = &name
	}

	volResp, err := r.client.CreateVolume(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create volume, got error: %v", err),
		)
		return
	}

	if volResp == nil || volResp.Data == nil || len(volResp.Data.Volumes) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from create volume API",
		)
		return
	}
	vol := volResp.Data.Volumes[0]

	if vol.Uuid == nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Volume UUID not returned by API",
		)
		return
	}

	data.UUID = types.StringValue(*vol.Uuid)
	if vol.Name != nil {
		data.Name = types.StringValue(*vol.Name)
	}

	// Get full volume state
	resp.Diagnostics.Append(r.readVolumeState(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource.Resource.
func (r *VolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VolumeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readVolumeState(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (r *VolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VolumeResourceModel
	var state VolumeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.SizeMB.Equal(state.SizeMB) {
		newSize := plan.SizeMB.ValueInt64()
		val := any(newSize)
		_, err := r.client.UpdateVolumeByUUID(ctx, state.UUID.ValueString(), platform.UpdateVolumeByUUIDRequestBody{
			Prop:  platform.UpdateVolumeByUUIDRequestBodyPropSize_mb,
			Op:    platform.UpdateVolumeByUUIDRequestBodyOpSet,
			Value: &val,
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Failed to update volume size, got error: %v", err),
			)
			return
		}
	}

	// Re-read full state after update
	data := plan
	data.UUID = state.UUID
	resp.Diagnostics.Append(r.readVolumeState(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (r *VolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VolumeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteVolumeByUUID(ctx, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete volume, got error: %v", err),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *VolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// readVolumeState fetches the current volume state from the API and populates
// computed fields in the model.
func (r *VolumeResource) readVolumeState(ctx context.Context, data *VolumeResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	volResp, err := r.client.GetVolumeByUUID(ctx, data.UUID.ValueString(), true)
	if err != nil {
		diags.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get volume state, got error: %v", err),
		)
		return diags
	}

	if volResp == nil || volResp.Data == nil || len(volResp.Data.Volumes) == 0 {
		diags.AddError(
			"Client Error",
			"Empty response from get volume API",
		)
		return diags
	}
	vol := volResp.Data.Volumes[0]

	// Don't overwrite user-configured name
	if data.Name.IsNull() && vol.Name != nil {
		data.Name = types.StringValue(*vol.Name)
	}
	if vol.SizeMb != nil {
		data.SizeMB = types.Int64Value(int64(*vol.SizeMb))
	}
	if vol.State != nil {
		data.State = types.StringValue(string(*vol.State))
	}
	if vol.Persistent != nil {
		data.Persistent = types.BoolValue(*vol.Persistent)
	}
	if vol.CreatedAt != nil {
		data.CreatedAt = types.StringValue(vol.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"))
	}

	if vol.AttachedTo != nil {
		attachedInstances := make([]models.VolumeInstanceModel, len(vol.AttachedTo))
		for i, inst := range vol.AttachedTo {
			if inst.Uuid != nil {
				attachedInstances[i].UUID = types.StringValue(*inst.Uuid)
			}
			if inst.Name != nil {
				attachedInstances[i].Name = types.StringValue(*inst.Name)
			}
		}
		var d diag.Diagnostics
		data.AttachedTo, d = types.ListValueFrom(ctx, models.VolumeInstanceModelType, attachedInstances)
		diags.Append(d...)
	} else {
		data.AttachedTo = types.ListValueMust(models.VolumeInstanceModelType, []attr.Value{})
	}

	return diags
}
