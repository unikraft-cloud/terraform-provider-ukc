// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"unikraft.com/cloud/sdk/platform"
)

func NewServicesDataSource() datasource.DataSource {
	return &ServicesDataSource{}
}

// ServicesDataSource defines the data source implementation.
type ServicesDataSource struct {
	client platform.Client
}

// Ensure ServicesDataSource satisfies various datasource interfaces.
var _ datasource.DataSource = &ServicesDataSource{}

// ServicesDataSourceModel describes the data source data model.
type ServicesDataSourceModel struct {
	UUIDs types.List `tfsdk:"uuids"`
}

// Metadata implements datasource.DataSource.
func (d *ServicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

// Schema implements datasource.DataSource.
func (d *ServicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides UUIDs of existing Unikraft Cloud service groups.",

		Attributes: map[string]schema.Attribute{
			"uuids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of service group UUIDs.",
			},
		},
	}
}

// Configure implements datasource.DataSource.
func (d *ServicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServicesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgResp, err := d.client.GetServiceGroups(ctx, nil, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to list service groups, got error: %v", err),
		)
		return
	}

	if sgResp == nil || sgResp.Data == nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from list service groups API",
		)
		return
	}

	uuids := make([]attr.Value, 0, len(sgResp.Data.ServiceGroups))
	for _, sg := range sgResp.Data.ServiceGroups {
		if sg.Uuid != nil {
			uuids = append(uuids, types.StringValue(*sg.Uuid))
		}
	}
	var diags diag.Diagnostics
	data.UUIDs, diags = types.ListValue(types.StringType, uuids)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
