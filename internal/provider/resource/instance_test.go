// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/unikraft-cloud/terraform-provider-unikraft-cloud/internal/provider/model"
	"sdk.kraft.cloud/client"
	"sdk.kraft.cloud/client/httpclient"
	"sdk.kraft.cloud/instances"
	"sdk.kraft.cloud/services"
)

// MockInstancesService is a mock implementation of instances.InstancesService
type MockInstancesService struct {
	mock.Mock
}

// Embedded ServiceClient interface methods
func (m *MockInstancesService) WithMetro(metro string) instances.InstancesService {
	args := m.Called(metro)
	return args.Get(0).(instances.InstancesService)
}

func (m *MockInstancesService) WithTimeout(timeout time.Duration) instances.InstancesService {
	args := m.Called(timeout)
	return args.Get(0).(instances.InstancesService)
}

func (m *MockInstancesService) WithHTTPClient(client httpclient.HTTPClient) instances.InstancesService {
	args := m.Called(client)
	return args.Get(0).(instances.InstancesService)
}

// Create implements instances.InstancesService
func (m *MockInstancesService) Create(ctx context.Context, req instances.CreateRequest) (*client.ServiceResponse[instances.CreateResponseItem], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.CreateResponseItem]), args.Error(1)
}

// Get implements instances.InstancesService
func (m *MockInstancesService) Get(ctx context.Context, uuids ...string) (*client.ServiceResponse[instances.GetResponseItem], error) {
	args := m.Called(ctx, uuids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.GetResponseItem]), args.Error(1)
}

// Delete implements instances.InstancesService
func (m *MockInstancesService) Delete(ctx context.Context, uuids ...string) (*client.ServiceResponse[instances.DeleteResponseItem], error) {
	args := m.Called(ctx, uuids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.DeleteResponseItem]), args.Error(1)
}

// List implements instances.InstancesService
func (m *MockInstancesService) List(ctx context.Context) (*client.ServiceResponse[instances.GetResponseItem], error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.GetResponseItem]), args.Error(1)
}

// Start implements instances.InstancesService
func (m *MockInstancesService) Start(ctx context.Context, waitTimeoutMs int, uuids ...string) (*client.ServiceResponse[instances.StartResponseItem], error) {
	args := m.Called(ctx, waitTimeoutMs, uuids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.StartResponseItem]), args.Error(1)
}

// Stop implements instances.InstancesService
func (m *MockInstancesService) Stop(ctx context.Context, waitTimeoutMs int, force bool, uuids ...string) (*client.ServiceResponse[instances.StopResponseItem], error) {
	args := m.Called(ctx, waitTimeoutMs, force, uuids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.StopResponseItem]), args.Error(1)
}

// Wait implements instances.InstancesService
func (m *MockInstancesService) Wait(ctx context.Context, state instances.State, timeoutMs int, ids ...string) (*client.ServiceResponse[instances.WaitResponseItem], error) {
	args := m.Called(ctx, state, timeoutMs, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.WaitResponseItem]), args.Error(1)
}

// Log implements instances.InstancesService
func (m *MockInstancesService) Log(ctx context.Context, uuid string, start, end int) (*client.ServiceResponse[instances.LogResponseItem], error) {
	args := m.Called(ctx, uuid, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.LogResponseItem]), args.Error(1)
}

// Metrics implements instances.InstancesService
func (m *MockInstancesService) Metrics(ctx context.Context, uuids ...string) (*client.ServiceResponse[instances.MetricsResponseItem], error) {
	args := m.Called(ctx, uuids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.ServiceResponse[instances.MetricsResponseItem]), args.Error(1)
}

// TailLogs implements instances.InstancesService
func (m *MockInstancesService) TailLogs(ctx context.Context, id string, follow bool, tail int, delay time.Duration) (chan string, chan error, error) {
	args := m.Called(ctx, id, follow, tail, delay)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(chan string), args.Get(1).(chan error), args.Error(2)
}

var _ instances.InstancesService = (*MockInstancesService)(nil)

func TestInstanceResource_Metadata(t *testing.T) {
	r := &InstanceResource{}
	req := resource.MetadataRequest{
		ProviderTypeName: "unikraft-cloud",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "unikraft-cloud_instance", resp.TypeName)
}

func TestInstanceResource_Schema(t *testing.T) {
	r := &InstanceResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	schema := resp.Schema
	assert.NotNil(t, schema)
	assert.Equal(t, "Allows the creation of Unikraft Cloud instances.", schema.MarkdownDescription)

	// Test required attributes
	imageAttr, exists := schema.Attributes["image"]
	assert.True(t, exists)
	assert.True(t, imageAttr.IsRequired())

	serviceGroupAttr, exists := schema.Attributes["service_group"]
	assert.True(t, exists)
	assert.True(t, serviceGroupAttr.IsRequired())

	// Test computed attributes
	uuidAttr, exists := schema.Attributes["uuid"]
	assert.True(t, exists)
	assert.True(t, uuidAttr.IsComputed())

	// Test optional attributes
	argsAttr, exists := schema.Attributes["args"]
	assert.True(t, exists)
	assert.True(t, argsAttr.IsOptional())
	assert.True(t, argsAttr.IsComputed())
}

func TestInstanceResource_Configure_Success(t *testing.T) {
	r := &InstanceResource{}
	mockClient := &MockInstancesService{}

	clients := map[string]interface{}{
		"instances": mockClient,
	}

	req := resource.ConfigureRequest{
		ProviderData: clients,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, mockClient, r.client)
}

func TestInstanceResource_Configure_NoProviderData(t *testing.T) {
	r := &InstanceResource{}
	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

func TestInstanceResource_Configure_InvalidProviderDataType(t *testing.T) {
	r := &InstanceResource{}
	req := resource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics[0].Summary(), "Unexpected Resource Configure Type")
}

func TestInstanceResource_Configure_MissingInstancesClient(t *testing.T) {
	r := &InstanceResource{}
	clients := map[string]interface{}{
		"other": "client",
	}

	req := resource.ConfigureRequest{
		ProviderData: clients,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics[0].Summary(), "Missing Instances Client")
}

func TestInstanceResource_Configure_InvalidInstancesClientType(t *testing.T) {
	r := &InstanceResource{}
	clients := map[string]interface{}{
		"instances": "invalid",
	}

	req := resource.ConfigureRequest{
		ProviderData: clients,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics[0].Summary(), "Invalid Instances Client Type")
}

func TestInstanceResource_Create_ValidatesCreateRequest(t *testing.T) {
	mockClient := &MockInstancesService{}

	expectedImage := "nginx:latest"
	expectedMemory := 256
	expectedArgs := []string{"--port", "8080"}
	expectedAutostart := true

	createResp := &client.ServiceResponse[instances.CreateResponseItem]{
		Data: client.ServiceResponseData[instances.CreateResponseItem]{
			Entries: []instances.CreateResponseItem{
				{
					UUID:        "test-uuid",
					Name:        "test-instance",
					PrivateIP:   "10.0.0.1",
					PrivateFQDN: "test.internal",
				},
			},
		},
	}

	getResp := &client.ServiceResponse[instances.GetResponseItem]{
		Data: client.ServiceResponseData[instances.GetResponseItem]{
			Entries: []instances.GetResponseItem{
				{
					UUID:       "test-uuid",
					Name:       "test-instance",
					Image:      expectedImage,
					Args:       expectedArgs,
					MemoryMB:   uint(expectedMemory),
					State:      instances.InstanceStateRunning,
					CreatedAt:  "2023-01-01T00:00:00Z",
					BootTimeUs: 1000,
					Env:        map[string]string{},
					ServiceGroup: &instances.GetCreateResponseServiceGroup{
						UUID: "sg-uuid",
						Name: "sg-name",
					},
					NetworkInterfaces: []instances.GetResponseNetworkInterface{},
				},
			},
		},
	}

	// Set up mock expectations
	mockClient.On("Create", mock.Anything, mock.MatchedBy(func(req instances.CreateRequest) bool {
		return req.Image == expectedImage &&
			*req.MemoryMB == expectedMemory &&
			len(req.Args) == 2 &&
			req.Args[0] == "--port" &&
			req.Args[1] == "8080" &&
			*req.Autostart == expectedAutostart &&
			req.ServiceGroup != nil &&
			len(req.ServiceGroup.Services) == 1 &&
			req.ServiceGroup.Services[0].Port == 80
	})).Return(createResp, nil)

	mockClient.On("Get", mock.Anything, []string{"test-uuid"}).Return(getResp, nil)

	// Create the resource model that would come from Terraform plan
	data := InstanceResourceModel{
		Image:     types.StringValue(expectedImage),
		Args:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("--port"), types.StringValue("8080")}),
		MemoryMB:  types.Int64Value(int64(expectedMemory)),
		Autostart: types.BoolValue(expectedAutostart),
		ServiceGroup: &models.SvcGrpModel{
			Services: []models.SvcModel{
				{Port: types.Int64Value(80)},
			},
		},
	}

	mockClient.AssertNotCalled(t, "Create")
	mockClient.AssertNotCalled(t, "Get")

	assert.Equal(t, expectedImage, data.Image.ValueString())
	assert.Equal(t, int64(expectedMemory), data.MemoryMB.ValueInt64())
}

func TestInstanceResource_Create_DefaultMemoryHandling(t *testing.T) {
	data := InstanceResourceModel{
		Image:    types.StringValue("nginx:latest"),
		MemoryMB: types.Int64Unknown(),
		ServiceGroup: &models.SvcGrpModel{
			Services: []models.SvcModel{
				{Port: types.Int64Value(80)},
			},
		},
	}

	if data.MemoryMB.IsUnknown() || data.MemoryMB.IsNull() {
		data.MemoryMB = types.Int64Value(128)
	}

	assert.Equal(t, int64(128), data.MemoryMB.ValueInt64())
}

func TestInstanceResource_Create_ServiceGroupMapping(t *testing.T) {
	serviceGroupServices := []models.SvcModel{
		{
			Port:            types.Int64Value(80),
			DestinationPort: types.Int64Value(8080),
			Handlers:        types.SetValueMust(types.StringType, []attr.Value{types.StringValue("http")}),
		},
	}

	data := InstanceResourceModel{
		Image: types.StringValue("nginx:latest"),
		ServiceGroup: &models.SvcGrpModel{
			Services: serviceGroupServices,
		},
	}

	createReq := instances.CreateRequest{
		Image: data.Image.ValueString(),
		ServiceGroup: &instances.CreateRequestServiceGroup{
			Services: make([]services.CreateRequestService, len(data.ServiceGroup.Services)),
		},
	}

	for i, svc := range data.ServiceGroup.Services {
		createReq.ServiceGroup.Services[i].Port = int(svc.Port.ValueInt64())
		createReq.ServiceGroup.Services[i].DestinationPort = ptr(int(svc.DestinationPort.ValueInt64()))

		if !svc.Handlers.IsUnknown() {
			createReq.ServiceGroup.Services[i].Handlers = append(createReq.ServiceGroup.Services[i].Handlers, services.Handler("http"))
		}
	}

	assert.Equal(t, "nginx:latest", createReq.Image)
	assert.Equal(t, 80, createReq.ServiceGroup.Services[0].Port)
	assert.Equal(t, 8080, *createReq.ServiceGroup.Services[0].DestinationPort)
	assert.Len(t, createReq.ServiceGroup.Services[0].Handlers, 1)
}

func TestInstanceResource_Create_Error(t *testing.T) {
	mockClient := &MockInstancesService{}

	mockClient.On("Create", mock.Anything, mock.AnythingOfType("instances.CreateRequest")).Return(nil, errors.New("API error"))

	ctx := context.Background()
	req := instances.CreateRequest{
		Image: "nginx:latest",
		ServiceGroup: &instances.CreateRequestServiceGroup{
			Services: []services.CreateRequestService{{Port: 80}},
		},
	}

	result, err := mockClient.Create(ctx, req)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
}

func TestInstanceResource_Read_Success(t *testing.T) {
	mockClient := &MockInstancesService{}

	getResp := &client.ServiceResponse[instances.GetResponseItem]{
		Data: client.ServiceResponseData[instances.GetResponseItem]{
			Entries: []instances.GetResponseItem{
				{
					UUID:        "test-uuid",
					Name:        "test-instance",
					Image:       "nginx@sha256:abc123",
					MemoryMB:    256,
					PrivateIP:   "10.0.0.1",
					PrivateFQDN: "test.internal",
					State:       instances.InstanceStateRunning,
					CreatedAt:   "2023-01-01T00:00:00Z",
					BootTimeUs:  1000,
					ServiceGroup: &instances.GetCreateResponseServiceGroup{
						UUID: "sg-uuid",
						Name: "sg-name",
						Domains: []services.GetCreateResponseDomain{
							{FQDN: "test.example.com"},
						},
					},
					NetworkInterfaces: []instances.GetResponseNetworkInterface{
						{
							UUID:      "net-uuid",
							PrivateIP: "10.0.0.1",
							MAC:       "00:00:00:00:00:01",
						},
					},
				},
			},
		},
	}

	mockClient.On("Get", mock.Anything, []string{"test-uuid"}).Return(getResp, nil)

	result, err := mockClient.Get(context.Background(), "test-uuid")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data.Entries, 1)

	ins := result.Data.Entries[0]
	assert.Equal(t, "test-uuid", ins.UUID)
	assert.Equal(t, "test-instance", ins.Name)
	assert.Equal(t, "nginx@sha256:abc123", ins.Image)
	assert.Equal(t, uint(256), ins.MemoryMB)

	mockClient.AssertExpectations(t)
}

func TestInstanceResource_Read_Error(t *testing.T) {
	mockClient := &MockInstancesService{}

	mockClient.On("Get", mock.Anything, []string{"test-uuid"}).Return(nil, errors.New("API error"))

	result, err := mockClient.Get(context.Background(), "test-uuid")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")

	mockClient.AssertExpectations(t)
}

func TestInstanceResource_Update(t *testing.T) {
	r := &InstanceResource{}
	req := resource.UpdateRequest{}
	resp := &resource.UpdateResponse{}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics[0].Summary(), "Unsupported")
	assert.Contains(t, resp.Diagnostics[0].Detail(), "This resource does not support updates")
}

func TestInstanceResource_Delete_Success(t *testing.T) {
	mockClient := &MockInstancesService{}

	deleteResp := &client.ServiceResponse[instances.DeleteResponseItem]{
		Data: client.ServiceResponseData[instances.DeleteResponseItem]{
			Entries: []instances.DeleteResponseItem{
				{
					UUID:   "test-uuid",
					Name:   "test-instance",
					Status: "success",
				},
			},
		},
	}
	mockClient.On("Delete", mock.Anything, []string{"test-uuid"}).Return(deleteResp, nil)

	result, err := mockClient.Delete(context.Background(), "test-uuid")

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockClient.AssertExpectations(t)
}

func TestInstanceResource_Delete_Error(t *testing.T) {
	mockClient := &MockInstancesService{}

	mockClient.On("Delete", mock.Anything, []string{"test-uuid"}).Return(nil, errors.New("API error"))

	result, err := mockClient.Delete(context.Background(), "test-uuid")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")

	mockClient.AssertExpectations(t)
}

// TODO(kristoffn): implement this function
// func TestInstanceResource_ImportState(t *testing.T) {
// }

func TestNewInstanceResource(t *testing.T) {
	resource := NewInstanceResource()
	assert.NotNil(t, resource)
	assert.IsType(t, &InstanceResource{}, resource)
}

func TestPtr(t *testing.T) {
	value := 42
	result := ptr(value)
	assert.Equal(t, &value, result)
	assert.Equal(t, 42, *result)

	stringValue := "test"
	stringResult := ptr(stringValue)
	assert.Equal(t, &stringValue, stringResult)
	assert.Equal(t, "test", *stringResult)

	boolValue := true
	boolResult := ptr(boolValue)
	assert.Equal(t, &boolValue, boolResult)
	assert.Equal(t, true, *boolResult)
}

func TestInstanceResourceModel_DataMapping(t *testing.T) {
	model := InstanceResourceModel{
		Image:     types.StringValue("nginx:latest"),
		Args:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("--port"), types.StringValue("8080")}),
		MemoryMB:  types.Int64Value(256),
		Autostart: types.BoolValue(true),
		UUID:      types.StringValue("test-uuid"),
		Name:      types.StringValue("test-instance"),
		State:     types.StringValue("running"),
	}

	assert.Equal(t, "nginx:latest", model.Image.ValueString())
	assert.Equal(t, int64(256), model.MemoryMB.ValueInt64())
	assert.True(t, model.Autostart.ValueBool())
	assert.Equal(t, "test-uuid", model.UUID.ValueString())
	assert.Equal(t, "test-instance", model.Name.ValueString())
	assert.Equal(t, "running", model.State.ValueString())

	assert.Equal(t, 2, len(model.Args.Elements()))
}

func TestInstanceResourceModel_ServiceGroupMapping(t *testing.T) {
	serviceGroup := &models.SvcGrpModel{
		UUID: types.StringValue("sg-uuid"),
		Name: types.StringValue("sg-name"),
		Services: []models.SvcModel{
			{
				Port:            types.Int64Value(80),
				DestinationPort: types.Int64Value(8080),
				Handlers:        types.SetValueMust(types.StringType, []attr.Value{types.StringValue("http")}),
			},
		},
	}

	assert.Equal(t, "sg-uuid", serviceGroup.UUID.ValueString())
	assert.Equal(t, "sg-name", serviceGroup.Name.ValueString())
	assert.Len(t, serviceGroup.Services, 1)
	assert.Equal(t, int64(80), serviceGroup.Services[0].Port.ValueInt64())
	assert.Equal(t, int64(8080), serviceGroup.Services[0].DestinationPort.ValueInt64())
	assert.Equal(t, 1, len(serviceGroup.Services[0].Handlers.Elements()))
}
