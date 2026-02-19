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
	models "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/model"

	"unikraft.com/cloud/sdk/platform"
)

func NewServiceDataSource() datasource.DataSource {
	return &ServiceDataSource{}
}

// ServiceDataSource defines the data source implementation.
type ServiceDataSource struct {
	client platform.Client
}

// Ensure ServiceDataSource satisfies various datasource interfaces.
var _ datasource.DataSource = &ServiceDataSource{}

// ServiceDataSourceModel describes the data source data model.
type ServiceDataSourceModel struct {
	UUID types.String `tfsdk:"uuid"`

	Name       types.String `tfsdk:"name"`
	CreatedAt  types.String `tfsdk:"created_at"`
	Persistent types.Bool   `tfsdk:"persistent"`
	Autoscale  types.Bool   `tfsdk:"autoscale"`
	SoftLimit  types.Int64  `tfsdk:"soft_limit"`
	HardLimit  types.Int64  `tfsdk:"hard_limit"`
	Services   types.List   `tfsdk:"services"`
	Domains    types.List   `tfsdk:"domains"`
	Instances  types.List   `tfsdk:"instances"`
}

// Metadata implements datasource.DataSource.
func (d *ServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

// Schema implements datasource.DataSource.
func (d *ServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides state information about a Unikraft Cloud service group.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique identifier of the service group.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable name of the service group.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time when the service group was created.",
			},
			"persistent": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the service group remains after the last instance detaches.",
			},
			"autoscale": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if autoscale is enabled for this service group.",
			},
			"soft_limit": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Soft request concurrency limit per instance used to decide when to wake standby instances.",
			},
			"hard_limit": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Hard request concurrency limit per instance; excess requests fail when no other instances are available.",
			},
			"services": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of published network ports for the service group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Public-facing port accessible from the Internet.",
						},
						"destination_port": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Internal port that traffic is forwarded to.",
						},
						"handlers": schema.SetAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "Connection handlers (tls, http, redirect).",
						},
					},
				},
			},
			"domains": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of domains associated with the service group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fqdn": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Publicly accessible domain name.",
						},
						"certificate": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Certificate associated with the domain.",
							Attributes: map[string]schema.Attribute{
								"uuid": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "UUID of the certificate.",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Name of the certificate.",
								},
								"state": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "State of the certificate (pending, valid, error).",
								},
							},
						},
					},
				},
			},
			"instances": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of instances assigned to the service group.",
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
		},
	}
}

// Configure implements datasource.DataSource.
func (d *ServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgResp, err := d.client.GetServiceGroupByUUID(ctx, data.UUID.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get service group, got error: %v", err),
		)
		return
	}

	if sgResp == nil || sgResp.Data == nil || len(sgResp.Data.ServiceGroups) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from get service group API",
		)
		return
	}
	sg := sgResp.Data.ServiceGroups[0]

	var diags diag.Diagnostics

	if sg.Name != nil {
		data.Name = types.StringValue(*sg.Name)
	}
	if sg.CreatedAt != nil {
		data.CreatedAt = types.StringValue(sg.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"))
	}
	if sg.Persistent != nil {
		data.Persistent = types.BoolValue(*sg.Persistent)
	}
	if sg.Autoscale != nil {
		data.Autoscale = types.BoolValue(*sg.Autoscale)
	}
	if sg.SoftLimit != nil {
		data.SoftLimit = types.Int64Value(int64(*sg.SoftLimit))
	}
	if sg.HardLimit != nil {
		data.HardLimit = types.Int64Value(int64(*sg.HardLimit))
	}

	svcModels := make([]models.SvcModel, len(sg.Services))
	for i, svc := range sg.Services {
		svcModels[i].Port = types.Int64Value(int64(svc.Port))
		if svc.DestinationPort != nil {
			svcModels[i].DestinationPort = types.Int64Value(int64(*svc.DestinationPort))
		}
		handlerVals := make([]attr.Value, len(svc.Handlers))
		for j, h := range svc.Handlers {
			handlerVals[j] = types.StringValue(string(h))
		}
		svcModels[i].Handlers, diags = types.SetValue(types.StringType, handlerVals)
		resp.Diagnostics.Append(diags...)
	}
	data.Services, diags = types.ListValueFrom(ctx, models.SvcModelType, svcModels)
	resp.Diagnostics.Append(diags...)

	domainModels := make([]models.ServiceGroupDomainModel, len(sg.Domains))
	for i, d := range sg.Domains {
		if d.Fqdn != nil {
			domainModels[i].FQDN = types.StringValue(*d.Fqdn)
		}
		if d.Certificate != nil {
			certAttrs := map[string]attr.Value{
				"uuid":  types.StringValue(""),
				"name":  types.StringValue(""),
				"state": types.StringValue(""),
			}
			if d.Certificate.Uuid != nil {
				certAttrs["uuid"] = types.StringValue(*d.Certificate.Uuid)
			}
			if d.Certificate.Name != nil {
				certAttrs["name"] = types.StringValue(*d.Certificate.Name)
			}
			if d.Certificate.State != nil {
				certAttrs["state"] = types.StringValue(string(*d.Certificate.State))
			}
			domainModels[i].Certificate, diags = types.ObjectValue(models.ServiceGroupCertAttrTypes, certAttrs)
			resp.Diagnostics.Append(diags...)
		} else {
			domainModels[i].Certificate = types.ObjectNull(models.ServiceGroupCertAttrTypes)
		}
	}
	data.Domains, diags = types.ListValueFrom(ctx, models.ServiceGroupDomainModelType, domainModels)
	resp.Diagnostics.Append(diags...)

	instModels := make([]models.VolumeInstanceModel, len(sg.Instances))
	for i, inst := range sg.Instances {
		if inst.Uuid != nil {
			instModels[i].UUID = types.StringValue(*inst.Uuid)
		}
		if inst.Name != nil {
			instModels[i].Name = types.StringValue(*inst.Name)
		}
	}
	data.Instances, diags = types.ListValueFrom(ctx, models.VolumeInstanceModelType, instModels)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
