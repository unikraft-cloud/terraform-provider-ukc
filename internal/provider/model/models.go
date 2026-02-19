// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// svcGrpModel describes the data model for an instance's service group.
type SvcGrpModel struct {
	UUID     types.String  `tfsdk:"uuid"`
	Name     types.String  `tfsdk:"name"`
	Services []SvcModel    `tfsdk:"services"`
	Domains  []domainModel `tfsdk:"domains"`
}

// svcModel describes the data model for a service group's service.
type SvcModel struct {
	Port            types.Int64 `tfsdk:"port"`
	DestinationPort types.Int64 `tfsdk:"destination_port"`
	Handlers        types.Set   `tfsdk:"handlers"`
}

// netwIfaceModel describes the data model for an instance's network interface.
type NetwIfaceModel struct {
	UUID      types.String `tfsdk:"uuid"`
	Name      types.String `tfsdk:"name"`
	PrivateIP types.String `tfsdk:"private_ip"`
	MAC       types.String `tfsdk:"mac"`
}

var NetwIfaceModelType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uuid":       types.StringType,
		"name":       types.StringType,
		"private_ip": types.StringType,
		"mac":        types.StringType,
	},
}

type domainModel struct {
	Name        types.String      `tfsdk:"name"`
	FQDN        types.String      `tfsdk:"fqdn"`
	Certificate *certificateModel `tfsdk:"certificate"`
}

type certificateModel struct {
	UUID  types.String `tfsdk:"uuid"`
	Name  types.String `tfsdk:"name"`
	State types.String `tfsdk:"state"`
}

// SvcModelType is the object type for SvcModel, used with types.ListValueFrom.
var SvcModelType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"port":             types.Int64Type,
		"destination_port": types.Int64Type,
		"handlers":         types.SetType{ElemType: types.StringType},
	},
}

// ServiceGroupCertAttrTypes are the attribute types for a certificate nested
// inside a service group domain.
var ServiceGroupCertAttrTypes = map[string]attr.Type{
	"uuid":  types.StringType,
	"name":  types.StringType,
	"state": types.StringType,
}

// ServiceGroupDomainModel describes the data model for a service group domain.
type ServiceGroupDomainModel struct {
	FQDN        types.String `tfsdk:"fqdn"`
	Certificate types.Object `tfsdk:"certificate"`
}

// ServiceGroupDomainModelType is the object type for ServiceGroupDomainModel.
var ServiceGroupDomainModelType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"fqdn":        types.StringType,
		"certificate": types.ObjectType{AttrTypes: ServiceGroupCertAttrTypes},
	},
}

// VolumeInstanceModel describes the data model for an instance attached to a volume.
type VolumeInstanceModel struct {
	UUID types.String `tfsdk:"uuid"`
	Name types.String `tfsdk:"name"`
}

var VolumeInstanceModelType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uuid": types.StringType,
		"name": types.StringType,
	},
}

// VolumeMountModel describes the data model for a volume mount.
type VolumeMountModel struct {
	UUID     types.String `tfsdk:"uuid"`
	Name     types.String `tfsdk:"name"`
	ReadOnly types.Bool   `tfsdk:"read_only"`
}

var VolumeMountModelType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uuid":      types.StringType,
		"name":      types.StringType,
		"read_only": types.BoolType,
	},
}
