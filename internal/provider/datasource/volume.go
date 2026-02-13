// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/model"

	"unikraft.com/cloud/sdk/platform"
)

func NewVolumeDataSource() datasource.DataSource {
	return &VolumeDataSource{}
}

// VolumeDataSource defines the data source implementation.
type VolumeDataSource struct {
	client platform.Client
}

// Ensure VolumeDataSource satisfies various datasource interfaces.
var _ datasource.DataSource = &VolumeDataSource{}

// VolumeDataSourceModel describes the data source data model.
type VolumeDataSourceModel struct {
	UUID types.String `tfsdk:"uuid"`

	Name       types.String `tfsdk:"name"`
	State      types.String `tfsdk:"state"`
	SizeMB     types.Int64  `tfsdk:"size_mb"`
	CreatedAt  types.String `tfsdk:"created_at"`
	Persistent types.Bool   `tfsdk:"persistent"`
	AttachedTo types.List   `tfsdk:"attached_to"`
	MountedBy  types.List   `tfsdk:"mounted_by"`
}

// Metadata implements datasource.DataSource.
func (d *VolumeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema implements datasource.DataSource.
func (d *VolumeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides state information about a Unikraft Cloud volume.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Unique identifier of the " +
					"[volume](https://docs.kraft.cloud/006-rest-api-v1-volumes.html)",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable name of the volume.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current state of the volume (uninitialized, initializing, available, idle, mounted, busy, error).",
			},
			"size_mb": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Size of the volume in megabytes.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time when the volume was created.",
			},
			"persistent": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the volume will stay alive when the last instance it is attached to is deleted.",
			},
			"attached_to": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of instances that this volume is attached to.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "UUID of the instance.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the instance.",
						},
					},
				},
			},
			"mounted_by": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of instances that have this volume mounted.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "UUID of the instance.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the instance.",
						},
						"read_only": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the volume is mounted read-only.",
						},
					},
				},
			},
		},
	}
}

// Configure implements datasource.DataSource.
func (d *VolumeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(platform.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected platform.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read implements datasource.DataSource.
func (d *VolumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VolumeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	volResp, err := d.client.GetVolumeByUUID(ctx, data.UUID.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get volume state, got error: %v", err),
		)
		return
	}

	if volResp == nil || volResp.Data == nil || len(volResp.Data.Volumes) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from get volume API",
		)
		return
	}
	vol := volResp.Data.Volumes[0]

	var diags diag.Diagnostics

	if vol.Name != nil {
		data.Name = types.StringValue(*vol.Name)
	}
	if vol.State != nil {
		data.State = types.StringValue(string(*vol.State))
	}
	if vol.SizeMb != nil {
		data.SizeMB = types.Int64Value(int64(*vol.SizeMb))
	}
	if vol.CreatedAt != nil {
		data.CreatedAt = types.StringValue(vol.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"))
	}
	if vol.Persistent != nil {
		data.Persistent = types.BoolValue(*vol.Persistent)
	}

	if vol.AttachedTo != nil {
		attachedTo := make([]models.VolumeInstanceModel, len(vol.AttachedTo))
		for i, inst := range vol.AttachedTo {
			if inst.Uuid != nil {
				attachedTo[i].UUID = types.StringValue(*inst.Uuid)
			}
			if inst.Name != nil {
				attachedTo[i].Name = types.StringValue(*inst.Name)
			}
		}
		data.AttachedTo, diags = types.ListValueFrom(ctx, models.VolumeInstanceModelType, attachedTo)
		resp.Diagnostics.Append(diags...)
	} else {
		data.AttachedTo = types.ListNull(models.VolumeInstanceModelType)
	}

	if vol.MountedBy != nil {
		mountedBy := make([]models.VolumeMountModel, len(vol.MountedBy))
		for i, mount := range vol.MountedBy {
			if mount.Uuid != nil {
				mountedBy[i].UUID = types.StringValue(*mount.Uuid)
			}
			if mount.Name != nil {
				mountedBy[i].Name = types.StringValue(*mount.Name)
			}
			if mount.ReadOnly != nil {
				mountedBy[i].ReadOnly = types.BoolValue(*mount.ReadOnly)
			}
		}
		data.MountedBy, diags = types.ListValueFrom(ctx, models.VolumeMountModelType, mountedBy)
		resp.Diagnostics.Append(diags...)
	} else {
		data.MountedBy = types.ListNull(models.VolumeMountModelType)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
