// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"unikraft.com/cloud/sdk/platform"
)

func NewVolumesDataSource() datasource.DataSource {
	return &VolumesDataSource{}
}

// VolumesDataSource defines the data source implementation.
type VolumesDataSource struct {
	client platform.Client
}

// Ensure VolumesDataSource satisfies various datasource interfaces.
var _ datasource.DataSource = &VolumesDataSource{}

// VolumesDataSourceModel describes the data source data model.
type VolumesDataSourceModel struct {
	States types.Set `tfsdk:"states"`

	UUIDs types.List `tfsdk:"uuids"`
}

// Metadata implements datasource.DataSource.
func (d *VolumesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volumes"
}

// Schema implements datasource.DataSource.
func (d *VolumesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides UUIDs of existing Unikraft Cloud volumes.",

		Attributes: map[string]schema.Attribute{
			"states": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Filter volumes based on their current state",
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							string(platform.VolumeStateUninitialized),
							string(platform.VolumeStateInitializing),
							string(platform.VolumeStateAvailable),
							string(platform.VolumeStateIdle),
							string(platform.VolumeStateMounted),
							string(platform.VolumeStateBusy),
							string(platform.VolumeStateBusy),
						),
					),
				},
			},
			"uuids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of volume UUIDs.",
			},
		},
	}
}

// Configure implements datasource.DataSource.
func (d *VolumesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *VolumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VolumesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	volResp, err := d.client.GetVolumes(ctx, nil, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to list volumes, got error: %v", err),
		)
		return
	}

	if volResp == nil || volResp.Data == nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from list volumes API",
		)
		return
	}

	volumesList := volResp.Data.Volumes

	// Filter by state if specified (client-side filtering)
	if len(data.States.Elements()) > 0 {
		stateVals := make([]types.String, 0, len(data.States.Elements()))
		resp.Diagnostics.Append(data.States.ElementsAs(ctx, &stateVals, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		filteredVolumes := make([]platform.Volume, 0)

		for _, vol := range volumesList {
			if vol.Uuid == nil {
				continue
			}
			volFullResp, err := d.client.GetVolumeByUUID(ctx, *vol.Uuid, false)
			if err != nil {
				resp.Diagnostics.AddError(
					"Client Error",
					fmt.Sprintf("Failed to get state of volume %s, got error: %v", *vol.Uuid, err),
				)
				return
			}

			if volFullResp == nil || volFullResp.Data == nil || len(volFullResp.Data.Volumes) == 0 {
				continue
			}

			// the number of possible states is small enough that iterating
			// them for every volume is reasonably cheap
			for _, st := range stateVals {
				if volFullResp.Data.Volumes[0].State != nil && string(*volFullResp.Data.Volumes[0].State) == st.ValueString() {
					filteredVolumes = append(filteredVolumes, vol)
					break
				}
			}
		}

		volumesList = filteredVolumes
	}

	uuids := make([]attr.Value, 0, len(volumesList))
	for _, vol := range volumesList {
		if vol.Uuid != nil {
			uuids = append(uuids, types.StringValue(*vol.Uuid))
		}
	}
	var diags diag.Diagnostics
	data.UUIDs, diags = types.ListValue(types.StringType, uuids)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
