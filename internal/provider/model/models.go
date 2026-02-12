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
