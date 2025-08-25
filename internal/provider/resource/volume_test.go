package resource

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	kcclient "sdk.kraft.cloud/client"
	"sdk.kraft.cloud/client/httpclient"
	"sdk.kraft.cloud/volumes"
)

type MockVolumesService struct {
	mock.Mock
}

func (m *MockVolumesService) Create(ctx context.Context, name string, sizeMB int) (*kcclient.ServiceResponse[volumes.CreateResponseItem], error) {
	args := m.Called(ctx, name, sizeMB)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kcclient.ServiceResponse[volumes.CreateResponseItem]), args.Error(1)
}

func (m *MockVolumesService) Get(ctx context.Context, ids ...string) (*kcclient.ServiceResponse[volumes.GetResponseItem], error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kcclient.ServiceResponse[volumes.GetResponseItem]), args.Error(1)
}

func (m *MockVolumesService) Delete(ctx context.Context, ids ...string) (*kcclient.ServiceResponse[volumes.DeleteResponseItem], error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kcclient.ServiceResponse[volumes.DeleteResponseItem]), args.Error(1)
}

func (m *MockVolumesService) Attach(ctx context.Context, volID, instance, at string, readOnly bool) (*kcclient.ServiceResponse[volumes.AttachResponseItem], error) {
	args := m.Called(ctx, volID, instance, at, readOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kcclient.ServiceResponse[volumes.AttachResponseItem]), args.Error(1)
}

func (m *MockVolumesService) Detach(ctx context.Context, id string, from string) (*kcclient.ServiceResponse[volumes.DetachResponseItem], error) {
	args := m.Called(ctx, id, from)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kcclient.ServiceResponse[volumes.DetachResponseItem]), args.Error(1)
}

func (m *MockVolumesService) List(ctx context.Context) (*kcclient.ServiceResponse[volumes.GetResponseItem], error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kcclient.ServiceResponse[volumes.GetResponseItem]), args.Error(1)
}

func (m *MockVolumesService) WithToken(token string) volumes.VolumesService {
	args := m.Called(token)
	return args.Get(0).(volumes.VolumesService)
}

func (m *MockVolumesService) WithHTTPClient(httpClient httpclient.HTTPClient) volumes.VolumesService {
	args := m.Called(httpClient)
	return args.Get(0).(volumes.VolumesService)
}

func (m *MockVolumesService) WithMetro(metro string) volumes.VolumesService {
	args := m.Called(metro)
	return args.Get(0).(volumes.VolumesService)
}

func (m *MockVolumesService) WithTimeout(timeout time.Duration) volumes.VolumesService {
	args := m.Called(timeout)
	return args.Get(0).(volumes.VolumesService)
}

func TestVolumeResource_Metadata(t *testing.T) {
	r := NewVolumeResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "kraftcloud",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "kraftcloud_volume", resp.TypeName)
}

func TestVolumeResource_Schema(t *testing.T) {
	r := NewVolumeResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema)
	assert.NotEmpty(t, resp.Schema.Attributes)

	assert.Contains(t, resp.Schema.Attributes, "size_mb")
	assert.Contains(t, resp.Schema.Attributes, "name")
	assert.Contains(t, resp.Schema.Attributes, "data")
	assert.Contains(t, resp.Schema.Attributes, "status")

	sizeMbAttr := resp.Schema.Attributes["size_mb"]
	assert.True(t, sizeMbAttr.IsRequired())

	nameAttr := resp.Schema.Attributes["name"]
	assert.True(t, nameAttr.IsOptional())
}

func TestVolumeResource_Configure_Success(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}

	clients := map[string]any{
		"volumes": mockClient,
	}

	req := resource.ConfigureRequest{
		ProviderData: clients,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, mockClient, r.client)
}

func TestVolumeResource_Configure_InvalidProviderData(t *testing.T) {
	r := &VolumeResource{}

	req := resource.ConfigureRequest{
		ProviderData: "invalid",
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unexpected Resource Configure Type")
}

func TestVolumeResource_Configure_MissingVolumesClient(t *testing.T) {
	r := &VolumeResource{}

	clients := map[string]any{
		"other": "client",
	}

	req := resource.ConfigureRequest{
		ProviderData: clients,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing Volumes Client")
}

func TestVolumeResource_Configure_InvalidVolumesClientType(t *testing.T) {
	r := &VolumeResource{}

	clients := map[string]any{
		"volumes": "not-a-volumes-service",
	}

	req := resource.ConfigureRequest{
		ProviderData: clients,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid Volumes Client Type")
}

func TestVolumeResource_Create_Success(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	apiResponse := &kcclient.ServiceResponse[volumes.CreateResponseItem]{
		Status:  "success",
		Message: "Volume created successfully",
		Data: kcclient.ServiceResponseData[volumes.CreateResponseItem]{
			Entries: []volumes.CreateResponseItem{
				{
					Name:   "test-volume",
					Status: "created",
					UUID:   "test-uuid-123",
				},
			},
		},
	}

	mockClient.On("Create", ctx, "test-volume", 100).Return(apiResponse, nil)

	schema := createTestSchema()
	planRaw := createPlanValue("test-volume", 100)

	req := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw:    planRaw,
			Schema: schema,
		},
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: schema,
		},
	}

	r.Create(ctx, req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	var resultData VolumeModel
	diags := resp.State.Get(ctx, &resultData)
	assert.False(t, diags.HasError())

	assert.Equal(t, "test-volume", resultData.Name.ValueString())
	assert.Equal(t, "success", resultData.Status.ValueString())
	assert.Equal(t, "Volume created successfully", resultData.Message.ValueString())

	mockClient.AssertExpectations(t)
}

func TestVolumeResource_Create_APIError(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	mockClient.On("Create", ctx, "test-volume", 100).Return(nil, errors.New("API error"))

	schema := createTestSchema()
	planRaw := createPlanValue("test-volume", 100)

	req := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw:    planRaw,
			Schema: schema,
		},
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: schema,
		},
	}

	r.Create(ctx, req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Client Error")

	mockClient.AssertExpectations(t)
}

func TestVolumeResource_Create_EmptyResponse(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	apiResponse := &kcclient.ServiceResponse[volumes.CreateResponseItem]{
		Status:  "success",
		Message: "Volume created successfully",
		Data: kcclient.ServiceResponseData[volumes.CreateResponseItem]{
			Entries: []volumes.CreateResponseItem{},
		},
	}

	mockClient.On("Create", ctx, "test-volume", 100).Return(apiResponse, nil)

	schema := createTestSchema()
	planRaw := createPlanValue("test-volume", 100)

	req := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw:    planRaw,
			Schema: schema,
		},
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: schema,
		},
	}

	r.Create(ctx, req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "API Error")

	mockClient.AssertExpectations(t)
}

func TestVolumeResource_Read_Success(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	apiResponse := &kcclient.ServiceResponse[volumes.GetResponseItem]{
		Status:  "success",
		Message: "Volume retrieved successfully",
		Data: kcclient.ServiceResponseData[volumes.GetResponseItem]{
			Entries: []volumes.GetResponseItem{
				{
					Name:   "test-volume",
					Status: "running",
					UUID:   "test-uuid-123",
					SizeMB: 100,
				},
			},
		},
	}

	mockClient.On("Get", ctx, []string{"test-uuid-123"}).Return(apiResponse, nil)

	schema := createTestSchema()
	stateRaw := createStateValue("test-volume1", "test-uuid-1234")

	req := resource.ReadRequest{
		State: tfsdk.State{
			Raw:    stateRaw,
			Schema: schema,
		},
	}
	resp := &resource.ReadResponse{
		State: tfsdk.State{
			Schema: schema,
		},
	}

	r.Read(ctx, req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	var resultData VolumeModel
	diags := resp.State.Get(ctx, &resultData)
	assert.False(t, diags.HasError())

	assert.Equal(t, "test-volume", resultData.Name.ValueString())
	assert.Equal(t, "success", resultData.Status.ValueString())

	mockClient.AssertExpectations(t)
}

func TestVolumeResource_Read_VolumeNotFound(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	apiResponse := &kcclient.ServiceResponse[volumes.GetResponseItem]{
		Status:  "success",
		Message: "Volume not found",
		Data: kcclient.ServiceResponseData[volumes.GetResponseItem]{
			Entries: []volumes.GetResponseItem{},
		},
	}

	mockClient.On("Get", ctx, []string{"test-uuid-123"}).Return(apiResponse, nil)

	schema := createTestSchema()
	stateRaw := createStateValue("test-volume", "test-uuid-123")

	req := resource.ReadRequest{
		State: tfsdk.State{
			Raw:    stateRaw,
			Schema: schema,
		},
	}
	resp := &resource.ReadResponse{
		State: tfsdk.State{
			Schema: schema,
		},
	}

	r.Read(ctx, req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	mockClient.AssertExpectations(t)
}

func TestVolumeResource_Read_InvalidState(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	schema := createTestSchema()
	stateRaw := createInvalidStateValue()

	req := resource.ReadRequest{
		State: tfsdk.State{
			Raw:    stateRaw,
			Schema: schema,
		},
	}
	resp := &resource.ReadResponse{
		State: tfsdk.State{
			Schema: schema,
		},
	}

	r.Read(ctx, req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid State")
}

func TestVolumeResource_Update(t *testing.T) {
	r := &VolumeResource{}

	req := resource.UpdateRequest{}
	resp := &resource.UpdateResponse{}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Unsupported")
}

func TestVolumeResource_Delete_Success(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	apiResponse := &kcclient.ServiceResponse[volumes.DeleteResponseItem]{
		Status:  "success",
		Message: "Volume deleted successfully",
		Data: kcclient.ServiceResponseData[volumes.DeleteResponseItem]{
			Entries: []volumes.DeleteResponseItem{
				{
					Status: "deleted",
					UUID:   "test-uuid-123",
					Name:   "test-volume",
				},
			},
		},
	}

	mockClient.On("Delete", ctx, []string{"test-uuid-123"}).Return(apiResponse, nil)

	schema := createTestSchema()
	stateRaw := createStateValue("test-volume", "test-uuid-123")

	req := resource.DeleteRequest{
		State: tfsdk.State{
			Raw:    stateRaw,
			Schema: schema,
		},
	}
	resp := &resource.DeleteResponse{}

	r.Delete(ctx, req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	mockClient.AssertExpectations(t)
}

func TestVolumeResource_Delete_APIError(t *testing.T) {
	r := &VolumeResource{}
	mockClient := &MockVolumesService{}
	r.client = mockClient

	ctx := context.Background()

	mockClient.On("Delete", ctx, []string{"test-uuid-123"}).Return(nil, errors.New("Delete failed"))

	schema := createTestSchema()
	stateRaw := createStateValue("test-volume", "test-uuid-123")

	req := resource.DeleteRequest{
		State: tfsdk.State{
			Raw:    stateRaw,
			Schema: schema,
		},
	}
	resp := &resource.DeleteResponse{}

	r.Delete(ctx, req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Client Error")

	mockClient.AssertExpectations(t)
}

// TODO(kristoffn): implement this function
// func TestVolumeResource_ImportState(t *testing.T) {
// }

func TestDataValue_ToTerraformValue(t *testing.T) {
	ctx := context.Background()

	volumeValue := VolumesValue{
		Name:   types.StringValue("test-volume"),
		Status: types.StringValue("running"),
		UUID:   types.StringValue("test-uuid"),
		state:  attr.ValueStateKnown,
	}

	volumesType := VolumesType{
		basetypes.ObjectType{
			AttrTypes: VolumesValue{}.AttributeTypes(ctx),
		},
	}

	volumesList, diags := types.ListValueFrom(ctx, volumesType, []VolumesValue{volumeValue})
	require.False(t, diags.HasError())

	dataValue := DataValue{
		Volumes: volumesList,
		state:   attr.ValueStateKnown,
	}

	tfValue, err := dataValue.ToTerraformValue(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tfValue)
}

func TestDataValue_IsNull(t *testing.T) {
	dataValue := NewDataValueNull()
	assert.True(t, dataValue.IsNull())
	assert.False(t, dataValue.IsUnknown())
}

func TestDataValue_IsUnknown(t *testing.T) {
	dataValue := NewDataValueUnknown()
	assert.False(t, dataValue.IsNull())
	assert.True(t, dataValue.IsUnknown())
}

func TestVolumesValue_Equal(t *testing.T) {
	v1 := VolumesValue{
		Name:   types.StringValue("test"),
		Status: types.StringValue("running"),
		UUID:   types.StringValue("uuid-123"),
		state:  attr.ValueStateKnown,
	}

	v2 := VolumesValue{
		Name:   types.StringValue("test"),
		Status: types.StringValue("running"),
		UUID:   types.StringValue("uuid-123"),
		state:  attr.ValueStateKnown,
	}

	v3 := VolumesValue{
		Name:   types.StringValue("different"),
		Status: types.StringValue("running"),
		UUID:   types.StringValue("uuid-123"),
		state:  attr.ValueStateKnown,
	}

	assert.True(t, v1.Equal(v2))
	assert.False(t, v1.Equal(v3))
}

func createTestSchema() schema.Schema {
	r := NewVolumeResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	return resp.Schema
}

func createPlanValue(name string, sizeMB int) tftypes.Value {
	dataType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"volumes": tftypes.List{
				ElementType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name":   tftypes.String,
						"status": tftypes.String,
						"uuid":   tftypes.String,
					},
				},
			},
		},
	}

	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"data":       dataType,
			"details":    tftypes.Bool,
			"message":    tftypes.String,
			"name":       tftypes.String,
			"op_time_us": tftypes.Number,
			"size_mb":    tftypes.Number,
			"status":     tftypes.String,
		},
	}, map[string]tftypes.Value{
		"data":       tftypes.NewValue(dataType, nil),
		"details":    tftypes.NewValue(tftypes.Bool, nil),
		"message":    tftypes.NewValue(tftypes.String, nil),
		"name":       tftypes.NewValue(tftypes.String, name),
		"op_time_us": tftypes.NewValue(tftypes.Number, nil),
		"size_mb":    tftypes.NewValue(tftypes.Number, sizeMB),
		"status":     tftypes.NewValue(tftypes.String, nil),
	})
}

func createStateValue(name, uuid string) tftypes.Value {
	volumesListValue := tftypes.NewValue(tftypes.List{
		ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":   tftypes.String,
				"status": tftypes.String,
				"uuid":   tftypes.String,
			},
		},
	}, []tftypes.Value{
		tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":   tftypes.String,
				"status": tftypes.String,
				"uuid":   tftypes.String,
			},
		}, map[string]tftypes.Value{
			"name":   tftypes.NewValue(tftypes.String, name),
			"status": tftypes.NewValue(tftypes.String, "running"),
			"uuid":   tftypes.NewValue(tftypes.String, uuid),
		}),
	})

	dataType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"volumes": tftypes.List{
				ElementType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name":   tftypes.String,
						"status": tftypes.String,
						"uuid":   tftypes.String,
					},
				},
			},
		},
	}

	dataValue := tftypes.NewValue(dataType, map[string]tftypes.Value{
		"volumes": volumesListValue,
	})

	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"data":       dataType,
			"details":    tftypes.Bool,
			"message":    tftypes.String,
			"name":       tftypes.String,
			"op_time_us": tftypes.Number,
			"size_mb":    tftypes.Number,
			"status":     tftypes.String,
		},
	}, map[string]tftypes.Value{
		"data":       dataValue,
		"details":    tftypes.NewValue(tftypes.Bool, nil),
		"message":    tftypes.NewValue(tftypes.String, "Volume message"),
		"name":       tftypes.NewValue(tftypes.String, name),
		"op_time_us": tftypes.NewValue(tftypes.Number, nil),
		"size_mb":    tftypes.NewValue(tftypes.Number, 100),
		"status":     tftypes.NewValue(tftypes.String, "success"),
	})
}

func createInvalidStateValue() tftypes.Value {
	dataType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"volumes": tftypes.List{
				ElementType: tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"name":   tftypes.String,
						"status": tftypes.String,
						"uuid":   tftypes.String,
					},
				},
			},
		},
	}

	emptyVolumesList := tftypes.NewValue(tftypes.List{
		ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":   tftypes.String,
				"status": tftypes.String,
				"uuid":   tftypes.String,
			},
		},
	}, []tftypes.Value{})

	dataValue := tftypes.NewValue(dataType, map[string]tftypes.Value{
		"volumes": emptyVolumesList,
	})

	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"data":       dataType,
			"details":    tftypes.Bool,
			"message":    tftypes.String,
			"name":       tftypes.String,
			"op_time_us": tftypes.Number,
			"size_mb":    tftypes.Number,
			"status":     tftypes.String,
		},
	}, map[string]tftypes.Value{
		"data":       dataValue,
		"details":    tftypes.NewValue(tftypes.Bool, nil),
		"message":    tftypes.NewValue(tftypes.String, nil),
		"name":       tftypes.NewValue(tftypes.String, "test-volume"),
		"op_time_us": tftypes.NewValue(tftypes.Number, nil),
		"size_mb":    tftypes.NewValue(tftypes.Number, 100),
		"status":     tftypes.NewValue(tftypes.String, nil),
	})
}
