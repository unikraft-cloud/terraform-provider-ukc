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

func NewInstanceDataSource() datasource.DataSource {
	return &InstanceDataSource{}
}

// InstanceDataSource defines the data source implementation.
type InstanceDataSource struct {
	client platform.Client
}

// Ensure InstanceDataSource satisfies various datasource interfaces.
var _ datasource.DataSource = &InstanceDataSource{}

// InstanceDataSourceModel describes the data source data model.
type InstanceDataSourceModel struct {
	UUID types.String `tfsdk:"uuid"`

	Name              types.String        `tfsdk:"name"`
	FQDN              types.String        `tfsdk:"fqdn"`
	PrivateIP         types.String        `tfsdk:"private_ip"`
	PrivateFQDN       types.String        `tfsdk:"private_fqdn"`
	State             types.String        `tfsdk:"state"`
	CreatedAt         types.String        `tfsdk:"created_at"`
	Image             types.String        `tfsdk:"image"`
	MemoryMB          types.Int64         `tfsdk:"memory_mb"`
	Args              types.List          `tfsdk:"args"`
	Env               types.Map           `tfsdk:"env"`
	ServiceGroup      *models.SvcGrpModel `tfsdk:"service_group"`
	NetworkInterfaces types.List          `tfsdk:"network_interfaces"`
	BootTimeUS        types.Int64         `tfsdk:"boot_time_us"`
}

// Metadata implements datasource.DataSource.
func (d *InstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema implements datasource.DataSource.
func (d *InstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides state information about a Unikraft Cloud instance.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Unique identifier of the " +
					"[instance](https://docs.kraft.cloud/002-rest-api-v1-instances.html)",
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
			"image": schema.StringAttribute{
				Computed: true,
			},
			"memory_mb": schema.Int64Attribute{
				Computed: true,
			},
			"args": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"env": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"service_group": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"services": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"port": schema.Int64Attribute{
									Computed: true,
								},
								"destination_port": schema.Int64Attribute{
									Computed: true,
								},
								"handlers": schema.SetAttribute{
									ElementType: types.StringType,
									Computed:    true,
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

// Configure implements datasource.DataSource.
func (d *InstanceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *InstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InstanceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	insResp, err := d.client.GetInstanceByUUID(ctx, data.UUID.ValueString(), true)
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
	if ins.Image != nil {
		data.Image = types.StringValue(*ins.Image)
	}
	if ins.MemoryMb != nil {
		data.MemoryMB = types.Int64Value(int64(*ins.MemoryMb))
	}
	if ins.BootTimeUs != nil {
		data.BootTimeUS = types.Int64Value(int64(*ins.BootTimeUs))
	}

	if ins.Args != nil {
		data.Args, diags = types.ListValueFrom(ctx, types.StringType, ins.Args)
		resp.Diagnostics.Append(diags...)
	}

	if ins.Env != nil {
		data.Env, diags = types.MapValueFrom(ctx, types.StringType, ins.Env)
		resp.Diagnostics.Append(diags...)
	}

	if ins.ServiceGroup != nil {
		sgModel := &models.SvcGrpModel{}
		if ins.ServiceGroup.Uuid != nil {
			sgModel.UUID = types.StringValue(*ins.ServiceGroup.Uuid)
		}
		if ins.ServiceGroup.Name != nil {
			sgModel.Name = types.StringValue(*ins.ServiceGroup.Name)
		}
		// Note: InstanceServiceGroup only contains Uuid, Name, and Domains
		// Services are exposed through the Domains field
		sgModel.Services = []models.SvcModel{}
		data.ServiceGroup = sgModel
	} else {
		data.ServiceGroup = &models.SvcGrpModel{}
	}

	if ins.NetworkInterfaces != nil {
		netwIfaces := make([]models.NetwIfaceModel, len(ins.NetworkInterfaces))
		for i, net := range ins.NetworkInterfaces {
			if net.Uuid != nil {
				netwIfaces[i].UUID = types.StringValue(*net.Uuid)
				netwIfaces[i].Name = types.StringValue(*net.Uuid) // No name in the API response
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
