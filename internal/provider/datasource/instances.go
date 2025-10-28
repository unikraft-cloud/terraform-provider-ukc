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

func NewInstancesDataSource() datasource.DataSource {
	return &InstancesDataSource{}
}

// InstancesDataSource defines the data source implementation.
type InstancesDataSource struct {
	client platform.Client
}

// Ensure InstancesDataSource satisfies various datasource interfaces.
var _ datasource.DataSource = &InstancesDataSource{}

// InstancesDataSourceModel describes the data source data model.
type InstancesDataSourceModel struct {
	States types.Set `tfsdk:"states"`

	UUIDs types.List `tfsdk:"uuids"`
}

// Metadata implements datasource.DataSource.
func (d *InstancesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instances"
}

// Schema implements datasource.DataSource.
func (d *InstancesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides UUIDs of existing Unikraft Cloud instances.",

		Attributes: map[string]schema.Attribute{
			"states": schema.SetAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "Filter instances based on their current " +
					"[state](https://docs.kraft.cloud/002-rest-api-v1-instances.html#instance-states)",
				Optional: true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"stopped",
							"starting",
							"running",
							"draining",
							"stopping",
						),
					),
				},
			},
			"uuids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure implements datasource.DataSource.
func (d *InstancesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *InstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InstancesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	insResp, err := d.client.GetInstances(ctx, nil, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to list instances, got error: %v", err),
		)
		return
	}

	if insResp == nil || insResp.Data == nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from list instances API",
		)
		return
	}

	instancesList := insResp.Data.Instances

	// FIXME(antoineco): filtering not implemented in SDK.
	// Implemented client side for the time being (expensive operation).
	if len(data.States.Elements()) > 0 {
		stateVals := make([]types.String, 0, len(data.States.Elements()))
		resp.Diagnostics.Append(data.States.ElementsAs(ctx, &stateVals, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		filteredInstances := make([]platform.Instance, 0)

		for _, ins := range instancesList {
			if ins.Uuid == nil {
				continue
			}
			insFullResp, err := d.client.GetInstanceByUUID(ctx, *ins.Uuid, false)
			if err != nil {
				resp.Diagnostics.AddError(
					"Client Error",
					fmt.Sprintf("Failed to get state of instance %s, got error: %v", *ins.Uuid, err),
				)
				return
			}

			if insFullResp == nil || insFullResp.Data == nil || len(insFullResp.Data.Instances) == 0 {
				continue
			}

			// the number of possible states is small enough that iterating
			// them for every instance is reasonably cheap
			for _, st := range stateVals {
				if insFullResp.Data.Instances[0].State != nil && string(*insFullResp.Data.Instances[0].State) == st.ValueString() {
					filteredInstances = append(filteredInstances, ins)
					break
				}
			}
		}

		instancesList = filteredInstances
	}

	uuids := make([]attr.Value, 0, len(instancesList))
	for _, ins := range instancesList {
		if ins.Uuid != nil {
			uuids = append(uuids, types.StringValue(*ins.Uuid))
		}
	}
	var diags diag.Diagnostics
	data.UUIDs, diags = types.ListValue(types.StringType, uuids)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
