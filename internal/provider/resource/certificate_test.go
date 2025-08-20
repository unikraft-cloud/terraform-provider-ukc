// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"sdk.kraft.cloud/certificates"
	"sdk.kraft.cloud/client"
	"sdk.kraft.cloud/client/httpclient"
)

// MockCertificatesService is a mock implementation of certificates.CertificatesService
type MockCertificatesService struct {
	mock.Mock
}

// Embedded ServiceClient interface methods
func (m *MockCertificatesService) WithMetro(metro string) certificates.CertificatesService {
	args := m.Called(metro)
	return args.Get(0).(certificates.CertificatesService)
}

func (m *MockCertificatesService) WithTimeout(timeout time.Duration) certificates.CertificatesService {
	args := m.Called(timeout)
	return args.Get(0).(certificates.CertificatesService)
}

func (m *MockCertificatesService) WithHTTPClient(client httpclient.HTTPClient) certificates.CertificatesService {
	args := m.Called(client)
	return args.Get(0).(certificates.CertificatesService)
}

// Create implements certificates.CertificatesService
func (m *MockCertificatesService) Create(ctx context.Context, req *certificates.CreateRequest) (*client.ServiceResponse[certificates.CreateResponseItem], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[certificates.CreateResponseItem]), args.Error(1)
}

// Get implements certificates.CertificatesService
func (m *MockCertificatesService) Get(ctx context.Context, uuid ...string) (*client.ServiceResponse[certificates.GetResponseItem], error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[certificates.GetResponseItem]), args.Error(1)
}

// Delete implements certificates.CertificatesService
func (m *MockCertificatesService) Delete(ctx context.Context, uuid ...string) (*client.ServiceResponse[certificates.DeleteResponseItem], error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[certificates.DeleteResponseItem]), args.Error(1)
}

// List implements certificates.CertificatesService
func (m *MockCertificatesService) List(ctx context.Context) (*client.ServiceResponse[certificates.GetResponseItem], error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[certificates.GetResponseItem]), args.Error(1)
}

var _ certificates.CertificatesService = (*MockCertificatesService)(nil)

func createMockCreateResponse() *client.ServiceResponse[certificates.CreateResponseItem] {
	return &client.ServiceResponse[certificates.CreateResponseItem]{
		Status:  "success",
		Message: "Certificate created successfully",
		Data: client.ServiceResponseData[certificates.CreateResponseItem]{
			Entries: []certificates.CreateResponseItem{
				{
					UUID: "test-uuid-123",
					Name: "test-cert",
				},
			},
		},
	}
}

func createMockGetResponse() *client.ServiceResponse[certificates.GetResponseItem] {
	return &client.ServiceResponse[certificates.GetResponseItem]{
		Status:  "success",
		Message: "Certificate retrieved successfully",
		Data: client.ServiceResponseData[certificates.GetResponseItem]{
			Entries: []certificates.GetResponseItem{
				{
					Status:       "active",
					UUID:         "test-uuid-123",
					Name:         "test-cert",
					CommonName:   "example.com",
					State:        "active",
					CreatedAt:    "2023-01-01T00:00:00Z",
					Subject:      "CN=example.com",
					Issuer:       "CN=Let's Encrypt Authority X3",
					SerialNumber: "12345678901234567890",
					NotBefore:    "2023-01-01T00:00:00Z",
					NotAfter:     "2023-12-31T23:59:59Z",
					Validation: &certificates.GetResponseValidation{
						Attempt: 1,
						Next:    "2023-01-02T00:00:00Z",
					},
					ServiceGroups: []certificates.GetResponseServiceGroup{
						{
							UUID: "sg-uuid-123",
							Name: "test-service-group",
						},
					},
				},
			},
		},
	}
}

func createMockDeleteResponse() *client.ServiceResponse[certificates.DeleteResponseItem] {
	return &client.ServiceResponse[certificates.DeleteResponseItem]{
		Status:  "success",
		Message: "Certificate deleted successfully",
		Data: client.ServiceResponseData[certificates.DeleteResponseItem]{
			Entries: []certificates.DeleteResponseItem{
				{
					Status: "deleted",
					UUID:   "test-uuid-123",
					Name:   "test-cert",
				},
			},
		},
	}
}

func TestNewCertificateResource(t *testing.T) {
	res := NewCertificateResource()
	assert.NotNil(t, res)
	assert.IsType(t, &CertificateResource{}, res)
}

func TestCertificateResource_Metadata(t *testing.T) {
	r := &CertificateResource{}
	req := resource.MetadataRequest{
		ProviderTypeName: "ukc",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "ukc_certificate", resp.TypeName)
}

func TestCertificateResource_Schema(t *testing.T) {
	r := &CertificateResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema)
	assert.NotEmpty(t, resp.Schema.Attributes)

	assert.True(t, resp.Schema.Attributes["chain"].IsRequired())
	assert.True(t, resp.Schema.Attributes["cn"].IsRequired())
	assert.True(t, resp.Schema.Attributes["pkey"].IsRequired())

	assert.True(t, resp.Schema.Attributes["name"].IsOptional())
	assert.True(t, resp.Schema.Attributes["uuid"].IsOptional())

	assert.True(t, resp.Schema.Attributes["data"].IsComputed())
	assert.True(t, resp.Schema.Attributes["message"].IsComputed())
	assert.True(t, resp.Schema.Attributes["status"].IsComputed())

	assert.True(t, resp.Schema.Attributes["pkey"].IsSensitive())
}

func TestCertificateResource_Configure(t *testing.T) {
	tests := []struct {
		name         string
		providerData interface{}
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "nil provider data",
			providerData: nil,
			expectError:  false,
		},
		{
			name:         "invalid provider data type",
			providerData: "invalid",
			expectError:  true,
			errorMsg:     "Expected map[string]any",
		},
		{
			name:         "missing certificates client",
			providerData: map[string]any{},
			expectError:  true,
			errorMsg:     "Certificates client not found",
		},
		{
			name: "invalid certificates client type",
			providerData: map[string]any{
				"certificates": "invalid",
			},
			expectError: true,
			errorMsg:    "Expected certificates.CertificatesService",
		},
		{
			name: "valid configuration",
			providerData: map[string]any{
				"certificates": &MockCertificatesService{},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CertificateResource{}
			req := resource.ConfigureRequest{
				ProviderData: tt.providerData,
			}
			resp := &resource.ConfigureResponse{}

			r.Configure(context.Background(), req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), tt.errorMsg)
			} else {
				assert.False(t, resp.Diagnostics.HasError())
				if tt.providerData != nil {
					assert.NotNil(t, r.client)
				}
			}
		})
	}
}

func TestCertificateResource_Create(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockCertificatesService)
		expectError    bool
		errorMsg       string
		expectedUUID   string
		expectedStatus string
	}{
		{
			name: "successful creation",
			setupMock: func(m *MockCertificatesService) {
				m.On("Create", mock.Anything, &certificates.CreateRequest{
					Name:  "test-cert",
					CN:    "example.com",
					Chain: "test-chain",
					PKey:  "test-private-key",
				}).Return(createMockCreateResponse(), nil)
				m.On("Get", mock.Anything, "test-uuid-123").Return(createMockGetResponse(), nil)
			},
			expectError:    false,
			expectedUUID:   "test-uuid-123",
			expectedStatus: "success",
		},
		{
			name: "create API error",
			setupMock: func(m *MockCertificatesService) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("API error"))
			},
			expectError: true,
			errorMsg:    "Failed to create certificate",
		},
		{
			name: "empty create response",
			setupMock: func(m *MockCertificatesService) {
				emptyResponse := &client.ServiceResponse[certificates.CreateResponseItem]{
					Data: client.ServiceResponseData[certificates.CreateResponseItem]{
						Entries: []certificates.CreateResponseItem{},
					},
				}
				m.On("Create", mock.Anything, mock.Anything).Return(emptyResponse, nil)
			},
			expectError: true,
			errorMsg:    "No certificate returned from create operation",
		},
		{
			name: "get details error after create",
			setupMock: func(m *MockCertificatesService) {
				m.On("Create", mock.Anything, mock.Anything).Return(createMockCreateResponse(), nil)
				m.On("Get", mock.Anything, "test-uuid-123").Return(nil, fmt.Errorf("Get API error"))
			},
			expectError: true,
			errorMsg:    "Failed to get full certificate details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockCertificatesService{}
			tt.setupMock(mockClient)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestCertificateResource_Read(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockCertificatesService)
		expectError  bool
		errorMsg     string
		expectRemove bool
	}{
		{
			name: "successful read",
			setupMock: func(m *MockCertificatesService) {
				m.On("Get", mock.Anything, "test-uuid-123").Return(createMockGetResponse(), nil)
			},
			expectError: false,
		},
		{
			name: "API error",
			setupMock: func(m *MockCertificatesService) {
				m.On("Get", mock.Anything, "test-uuid-123").Return(nil, fmt.Errorf("API error"))
			},
			expectError: true,
			errorMsg:    "Failed to read certificate",
		},
		{
			name: "certificate not found",
			setupMock: func(m *MockCertificatesService) {
				emptyResponse := &client.ServiceResponse[certificates.GetResponseItem]{
					Data: client.ServiceResponseData[certificates.GetResponseItem]{
						Entries: []certificates.GetResponseItem{},
					},
				}
				m.On("Get", mock.Anything, "test-uuid-123").Return(emptyResponse, nil)
			},
			expectError:  false,
			expectRemove: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockCertificatesService{}
			tt.setupMock(mockClient)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestCertificateResource_Update(t *testing.T) {
	r := &CertificateResource{}
	req := resource.UpdateRequest{}
	resp := &resource.UpdateResponse{}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unsupported")
}

func TestCertificateResource_Delete(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockCertificatesService)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful deletion",
			setupMock: func(m *MockCertificatesService) {
				m.On("Delete", mock.Anything, "test-uuid-123").Return(createMockDeleteResponse(), nil)
			},
			expectError: false,
		},
		{
			name: "API error",
			setupMock: func(m *MockCertificatesService) {
				m.On("Delete", mock.Anything, "test-uuid-123").Return(nil, fmt.Errorf("Delete API error"))
			},
			expectError: true,
			errorMsg:    "Failed to delete certificate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockCertificatesService{}
			tt.setupMock(mockClient)
			mockClient.AssertExpectations(t)
		})
	}
}

// TODO(kristoffn): implement this function
// func TestCertificateResource_ImportState(t *testing.T) {
// }

func TestDataTypeCertificate_Equal(t *testing.T) {
	ctx := context.Background()

	type1 := DataTypeCertificate{
		basetypes.ObjectType{
			AttrTypes: CertificateDataValue{}.AttributeTypes(ctx),
		},
	}
	type2 := DataTypeCertificate{
		basetypes.ObjectType{
			AttrTypes: CertificateDataValue{}.AttributeTypes(ctx),
		},
	}

	assert.True(t, type1.Equal(type2))
	assert.False(t, type1.Equal(basetypes.StringType{}))
}

func TestDataTypeCertificate_String(t *testing.T) {
	dataType := DataTypeCertificate{}
	assert.Equal(t, "DataTypeCertificate", dataType.String())
}

func TestNewCertificateDataValue(t *testing.T) {
	ctx := context.Background()

	attributeTypes := CertificateDataValue{}.AttributeTypes(ctx)
	certificatesList, _ := types.ListValueFrom(ctx, CertificatesType{
		ObjectType: types.ObjectType{
			AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
		},
	}, []CertificatesValue{})

	attributes := map[string]attr.Value{
		"certificates": certificatesList,
	}

	value, diags := NewCertificateDataValue(attributeTypes, attributes)

	assert.False(t, diags.HasError())
	assert.False(t, value.IsNull())
	assert.False(t, value.IsUnknown())
}

func TestNewCertificateDataValueNull(t *testing.T) {
	value := NewCertificateDataValueNull()
	assert.True(t, value.IsNull())
	assert.False(t, value.IsUnknown())
}

func TestNewCertificateDataValueUnknown(t *testing.T) {
	value := NewCertificateDataValueUnknown()
	assert.False(t, value.IsNull())
	assert.True(t, value.IsUnknown())
}

func TestCertificateDataValue_ToTerraformValue(t *testing.T) {
	ctx := context.Background()

	// Test null value
	nullValue := NewCertificateDataValueNull()
	tfValue, err := nullValue.ToTerraformValue(ctx)
	assert.NoError(t, err)
	assert.True(t, tfValue.IsNull())

	// Test unknown value
	unknownValue := NewCertificateDataValueUnknown()
	tfValue, err = unknownValue.ToTerraformValue(ctx)
	assert.NoError(t, err)
	assert.False(t, tfValue.IsKnown())
}

func TestCertificateDataValue_Equal(t *testing.T) {
	value1 := NewCertificateDataValueNull()
	value2 := NewCertificateDataValueNull()
	value3 := NewCertificateDataValueUnknown()

	assert.True(t, value1.Equal(value2))
	assert.False(t, value1.Equal(value3))
	assert.False(t, value1.Equal(types.StringValue("test")))
}

func TestCertificatesType_Equal(t *testing.T) {
	ctx := context.Background()

	type1 := CertificatesType{
		basetypes.ObjectType{
			AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
		},
	}
	type2 := CertificatesType{
		basetypes.ObjectType{
			AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
		},
	}

	assert.True(t, type1.Equal(type2))
	assert.False(t, type1.Equal(basetypes.StringType{}))
}

func TestCertificatesType_String(t *testing.T) {
	certType := CertificatesType{}
	assert.Equal(t, "CertificatesType", certType.String())
}

func TestNewCertificatesValue(t *testing.T) {
	ctx := context.Background()

	// Test valid creation
	attributeTypes := CertificatesValue{}.AttributeTypes(ctx)
	attributes := map[string]attr.Value{
		"common_name": types.StringValue("example.com"),
		"created_at":  types.StringValue("2023-01-01T00:00:00Z"),
		"name":        types.StringValue("test-cert"),
		"state":       types.StringValue("active"),
		"uuid":        types.StringValue("test-uuid-123"),
	}

	value, diags := NewCertificatesValue(attributeTypes, attributes)

	assert.False(t, diags.HasError())
	assert.False(t, value.IsNull())
	assert.False(t, value.IsUnknown())
	assert.Equal(t, "example.com", value.CommonName.ValueString())
}

func TestNewCertificatesValueMust_Panic(t *testing.T) {
	ctx := context.Background()

	// Test that it panics with invalid attributes
	attributeTypes := CertificatesValue{}.AttributeTypes(ctx)
	attributes := map[string]attr.Value{
		"invalid_attr": types.StringValue("test"),
	}

	assert.Panics(t, func() {
		NewCertificatesValueMust(attributeTypes, attributes)
	})
}

func TestCertificatesValue_ToTerraformValue(t *testing.T) {
	ctx := context.Background()

	nullValue := NewCertificatesValueNull()
	tfValue, err := nullValue.ToTerraformValue(ctx)
	assert.NoError(t, err)
	assert.True(t, tfValue.IsNull())

	unknownValue := NewCertificatesValueUnknown()
	tfValue, err = unknownValue.ToTerraformValue(ctx)
	assert.NoError(t, err)
	assert.False(t, tfValue.IsKnown())
}

func TestCertificatesValue_Equal(t *testing.T) {
	value1 := NewCertificatesValueNull()
	value2 := NewCertificatesValueNull()
	value3 := NewCertificatesValueUnknown()

	assert.True(t, value1.Equal(value2))
	assert.False(t, value1.Equal(value3))
	assert.False(t, value1.Equal(types.StringValue("test")))
}

func TestCertificatesValue_String(t *testing.T) {
	value := NewCertificatesValueNull()
	assert.Equal(t, "CertificatesValue", value.String())
}

func TestCertificatesValue_ToObjectValue(t *testing.T) {
	ctx := context.Background()

	nullValue := NewCertificatesValueNull()
	objValue, diags := nullValue.ToObjectValue(ctx)
	assert.False(t, diags.HasError())
	assert.True(t, objValue.IsNull())

	unknownValue := NewCertificatesValueUnknown()
	objValue, diags = unknownValue.ToObjectValue(ctx)
	assert.False(t, diags.HasError())
	assert.True(t, objValue.IsUnknown())
}

func TestCertificatesValue_AttributeTypes(t *testing.T) {
	ctx := context.Background()

	value := NewCertificatesValueNull()
	attrTypes := value.AttributeTypes(ctx)

	expectedTypes := map[string]attr.Type{
		"common_name": basetypes.StringType{},
		"created_at":  basetypes.StringType{},
		"name":        basetypes.StringType{},
		"state":       basetypes.StringType{},
		"uuid":        basetypes.StringType{},
	}

	assert.Equal(t, len(expectedTypes), len(attrTypes))
	for key, expectedType := range expectedTypes {
		assert.True(t, expectedType.Equal(attrTypes[key]), "Type mismatch for attribute %s", key)
	}
}
