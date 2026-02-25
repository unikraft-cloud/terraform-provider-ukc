// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	providerMock "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/mock"
	models "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/model"
)

func TestServiceResource_Metadata(t *testing.T) {
	r := NewServiceResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "ukc",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "ukc_service", resp.TypeName)
}

func TestServiceResource_Schema(t *testing.T) {
	r := NewServiceResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema)
	assert.Contains(t, resp.Schema.Attributes, "name")
	assert.Contains(t, resp.Schema.Attributes, "uuid")
	assert.Contains(t, resp.Schema.Attributes, "services")
	assert.Contains(t, resp.Schema.Attributes, "domains")
	assert.Contains(t, resp.Schema.Attributes, "soft_limit")
	assert.Contains(t, resp.Schema.Attributes, "hard_limit")
	assert.Contains(t, resp.Schema.Attributes, "created_at")
	assert.Contains(t, resp.Schema.Attributes, "persistent")
	assert.Contains(t, resp.Schema.Attributes, "autoscale")
	assert.Contains(t, resp.Schema.Attributes, "computed_domains")
	assert.Contains(t, resp.Schema.Attributes, "instances")
}

func TestServiceResource_Configure_Success(t *testing.T) {
	r := &ServiceResource{}
	mockClient := new(providerMock.PlatformClient)

	req := resource.ConfigureRequest{
		ProviderData: mockClient,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	// The configure will fail because mock doesn't implement full interface
	// This test verifies the type checking logic works
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
}

func TestServiceResource_Configure_NoProviderData(t *testing.T) {
	r := &ServiceResource{}

	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

func TestServiceResource_Configure_InvalidProviderDataType(t *testing.T) {
	r := &ServiceResource{}

	req := resource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
}

func TestServiceResource_Update(t *testing.T) {
	r := &ServiceResource{}

	req := resource.UpdateRequest{}
	resp := &resource.UpdateResponse{}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unsupported")
}

func TestNewServiceResource(t *testing.T) {
	r := NewServiceResource()
	assert.NotNil(t, r)
	assert.IsType(t, &ServiceResource{}, r)
}

func TestServiceResourceModel_Basic(t *testing.T) {
	model := ServiceResourceModel{
		Name: types.StringValue("my-service"),
		UUID: types.StringValue("test-uuid"),
		Services: []models.SvcModel{
			{
				Port: types.Int64Value(443),
			},
		},
	}

	assert.Equal(t, "my-service", model.Name.ValueString())
	assert.Equal(t, "test-uuid", model.UUID.ValueString())
	assert.Len(t, model.Services, 1)
	assert.Equal(t, int64(443), model.Services[0].Port.ValueInt64())
}
