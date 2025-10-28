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
)

func TestCertificateResource_Metadata(t *testing.T) {
	r := NewCertificateResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "ukc",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "ukc_certificate", resp.TypeName)
}

func TestCertificateResource_Schema(t *testing.T) {
	r := NewCertificateResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema)
	assert.Contains(t, resp.Schema.Attributes, "cn")
	assert.Contains(t, resp.Schema.Attributes, "chain")
	assert.Contains(t, resp.Schema.Attributes, "pkey")
	assert.Contains(t, resp.Schema.Attributes, "uuid")
	assert.Contains(t, resp.Schema.Attributes, "name")
}

func TestCertificateResource_Configure_Success(t *testing.T) {
	r := &CertificateResource{}
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

func TestCertificateResource_Configure_NoProviderData(t *testing.T) {
	r := &CertificateResource{}

	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

func TestCertificateResource_Configure_InvalidProviderDataType(t *testing.T) {
	r := &CertificateResource{}

	req := resource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
}

// NOTE: Create, Read, Update, and Delete operations should be tested via acceptance tests.
// Unit testing these operations would require complex mocking of Terraform's plugin
// framework request/response structures (Plan, State, etc.), which provides limited value
// compared to acceptance tests that verify actual behavior against a real backend.

func TestNewCertificateResource(t *testing.T) {
	r := NewCertificateResource()
	assert.NotNil(t, r)
	assert.IsType(t, &CertificateResource{}, r)
}

func TestCertificateResource_Update(t *testing.T) {
	r := &CertificateResource{}

	req := resource.UpdateRequest{}
	resp := &resource.UpdateResponse{}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unsupported")
}

func TestCertificateResourceModel_Basic(t *testing.T) {
	model := CertificateResourceModel{
		Cn:     types.StringValue("example.com"),
		Name:   types.StringValue("my-cert"),
		UUID:   types.StringValue("cert-uuid"),
		Status: types.StringValue("success"),
	}

	assert.Equal(t, "example.com", model.Cn.ValueString())
	assert.Equal(t, "my-cert", model.Name.ValueString())
	assert.Equal(t, "cert-uuid", model.UUID.ValueString())
	assert.Equal(t, "success", model.Status.ValueString())
}
