// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/model"

	"unikraft.com/cloud/sdk/platform"
)

func NewServiceResource() resource.Resource {
	return &ServiceResource{}
}

// ServiceResource defines the resource implementation.
type ServiceResource struct {
	client platform.Client
}

// Ensure ServiceResource satisfies various resource interfaces.
var (
	_ resource.Resource                = &ServiceResource{}
	_ resource.ResourceWithImportState = &ServiceResource{}
)

// ServiceResourceModel describes the resource data model.
type ServiceResourceModel struct {
	// Configurable (user inputs)
	Name      types.String    `tfsdk:"name"`
	Services  []models.SvcModel `tfsdk:"services"`
	Domains   []ServiceDomainInputModel `tfsdk:"domains"`
	SoftLimit types.Int64     `tfsdk:"soft_limit"`
	HardLimit types.Int64     `tfsdk:"hard_limit"`

	// Computed (API-populated)
	UUID       types.String `tfsdk:"uuid"`
	CreatedAt  types.String `tfsdk:"created_at"`
	Persistent types.Bool   `tfsdk:"persistent"`
	Autoscale  types.Bool   `tfsdk:"autoscale"`

	// Computed nested lists
	ComputedDomains  types.List `tfsdk:"computed_domains"`
	Instances        types.List `tfsdk:"instances"`
}

// ServiceDomainInputModel describes a domain input for service group creation.
type ServiceDomainInputModel struct {
	Name        types.String `tfsdk:"name"`
	Certificate types.String `tfsdk:"certificate"`
}

// Metadata implements resource.Resource.
func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

// Schema implements resource.Resource.
func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Allows the creation of Unikraft Cloud service groups.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Human-readable name of the service group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"services": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Port mappings with handlers",
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
				Optional:            true,
				MarkdownDescription: "Domain names with optional certificate reference",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Domain name",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"certificate": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Certificate UUID or name",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplaceIfConfigured(),
							},
						},
					},
				},
			},
			"soft_limit": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Load balancer soft limit",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"hard_limit": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Load balancer hard limit",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
					int64planmodifier.UseStateForUnknown(),
				},
			},

			// Computed attributes
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the service group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"persistent": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"autoscale": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"computed_domains": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Domains as returned by the API",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fqdn": schema.StringAttribute{
							Computed: true,
						},
						"certificate": schema.SingleNestedAttribute{
							Computed: true,
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
			},
			"instances": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Instances attached to this service group",
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
func (r *ServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := platform.CreateServiceGroupRequest{}

	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		name := data.Name.ValueString()
		in.Name = &name
	}

	if !data.SoftLimit.IsUnknown() && !data.SoftLimit.IsNull() {
		softLimit := uint64(data.SoftLimit.ValueInt64())
		in.SoftLimit = &softLimit
	}

	if !data.HardLimit.IsUnknown() && !data.HardLimit.IsNull() {
		hardLimit := uint64(data.HardLimit.ValueInt64())
		in.HardLimit = &hardLimit
	}

	// Build services
	if len(data.Services) > 0 {
		sgServices := make([]platform.Service, len(data.Services))
		for i, svc := range data.Services {
			port := uint32(svc.Port.ValueInt64())
			sgServices[i].Port = port

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
		in.Services = sgServices
	}

	// Build domains
	if len(data.Domains) > 0 {
		sgDomains := make([]platform.CreateServiceGroupRequestDomain, len(data.Domains))
		for i, dom := range data.Domains {
			sgDomains[i].Name = dom.Name.ValueString()

			if !dom.Certificate.IsUnknown() && !dom.Certificate.IsNull() {
				sgDomains[i].Certificate = &platform.CreateServiceGroupRequestDomainCertificate{
					Uuid: dom.Certificate.ValueString(),
				}
			}
		}
		in.Domains = sgDomains
	}

	if resp.Diagnostics.HasError() {
		return
	}

	sgResp, err := r.client.CreateServiceGroup(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create service group, got error: %v", err),
		)
		return
	}

	if sgResp == nil || sgResp.Data == nil || len(sgResp.Data.ServiceGroups) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from create service group API",
		)
		return
	}
	sg := sgResp.Data.ServiceGroups[0]

	if sg.Uuid == nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Service group UUID not returned by API",
		)
		return
	}

	data.UUID = types.StringValue(*sg.Uuid)
	if sg.Name != nil {
		data.Name = types.StringValue(*sg.Name)
	}

	// Get full details
	sgFullResp, err := r.client.GetServiceGroupByUUID(ctx, *sg.Uuid, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get service group details, got error: %v", err),
		)
		return
	}

	if sgFullResp == nil || sgFullResp.Data == nil || len(sgFullResp.Data.ServiceGroups) == 0 {
		resp.Diagnostics.AddError(
			"Client Error",
			"Empty response from get service group API",
		)
		return
	}
	sgFull := sgFullResp.Data.ServiceGroups[0]

	r.populateModelFromAPI(ctx, &data, &sgFull, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource.Resource.
func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgResp, err := r.client.GetServiceGroupByUUID(ctx, data.UUID.ValueString(), true)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get service group state, got error: %v", err),
		)
		return
	}

	if sgResp == nil || sgResp.Data == nil || len(sgResp.Data.ServiceGroups) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	sg := sgResp.Data.ServiceGroups[0]

	if sg.Name != nil {
		data.Name = types.StringValue(*sg.Name)
	}

	r.populateModelFromAPI(ctx, &data, &sg, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unsupported",
		"This resource does not support updates. Configuration changes were expected to have triggered a replacement "+
			"of the resource. Please report this issue to the provider developers.",
	)
}

// Delete implements resource.Resource.
func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteServiceGroupByUUID(ctx, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete service group, got error: %v", err),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// populateModelFromAPI populates the resource model from an API ServiceGroup response.
func (r *ServiceResource) populateModelFromAPI(ctx context.Context, data *ServiceResourceModel, sg *platform.ServiceGroup, diagnostics *diag.Diagnostics) {
	var diags diag.Diagnostics

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

	// Populate services from API response
	if sg.Services != nil {
		svcModels := make([]models.SvcModel, len(sg.Services))
		for i, svc := range sg.Services {
			svcModels[i].Port = types.Int64Value(int64(svc.Port))
			if svc.DestinationPort != nil {
				svcModels[i].DestinationPort = types.Int64Value(int64(*svc.DestinationPort))
			}
			if svc.Handlers != nil {
				handlers := make([]string, len(svc.Handlers))
				for j, h := range svc.Handlers {
					handlers[j] = string(h)
				}
				svcModels[i].Handlers, diags = types.SetValueFrom(ctx, types.StringType, handlers)
				diagnostics.Append(diags...)
			}
		}
		data.Services = svcModels
	}

	// Populate computed domains
	if sg.Domains != nil {
		domainModels := make([]models.ServiceGroupDomainModel, len(sg.Domains))
		for i, dom := range sg.Domains {
			if dom.Fqdn != nil {
				domainModels[i].FQDN = types.StringValue(*dom.Fqdn)
			}
			if dom.Certificate != nil {
				domainModels[i].Certificate = &models.ServiceGroupDomainCertificateModel{}
				if dom.Certificate.Uuid != nil {
					domainModels[i].Certificate.UUID = types.StringValue(*dom.Certificate.Uuid)
				}
				if dom.Certificate.Name != nil {
					domainModels[i].Certificate.Name = types.StringValue(*dom.Certificate.Name)
				}
			}
		}
		data.ComputedDomains, diags = types.ListValueFrom(ctx, models.ServiceGroupDomainModelType, domainModels)
		diagnostics.Append(diags...)
	}

	// Populate instances
	if sg.Instances != nil {
		instanceModels := make([]models.ServiceGroupInstanceModel, len(sg.Instances))
		for i, inst := range sg.Instances {
			if inst.Uuid != nil {
				instanceModels[i].UUID = types.StringValue(*inst.Uuid)
			}
			if inst.Name != nil {
				instanceModels[i].Name = types.StringValue(*inst.Name)
			}
		}
		data.Instances, diags = types.ListValueFrom(ctx, models.ServiceGroupInstanceModelType, instanceModels)
		diagnostics.Append(diags...)
	}
}
