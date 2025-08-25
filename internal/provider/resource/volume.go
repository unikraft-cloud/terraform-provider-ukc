package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"sdk.kraft.cloud/volumes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewVolumeResource() resource.Resource {
	return &VolumeResource{}
}

// VolumesService defines the resource implementation.
type VolumeResource struct {
	client volumes.VolumesService
}

// Ensure VolumesService satisfies various resource interfaces.
var (
	_ resource.Resource                = &VolumeResource{}
	_ resource.ResourceWithImportState = &VolumeResource{}
)

func (r *VolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"volumes": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "The name of the newly created volume.",
									MarkdownDescription: "The name of the newly created volume.",
								},
								"status": schema.StringAttribute{
									Computed:            true,
									Description:         "The status of the response.",
									MarkdownDescription: "The status of the response.",
								},
								"uuid": schema.StringAttribute{
									Computed:            true,
									Description:         "UUID of the newly created volume.",
									MarkdownDescription: "UUID of the newly created volume.",
								},
							},
							CustomType: VolumesType{
								ObjectType: types.ObjectType{
									AttrTypes: VolumesValue{}.AttributeTypes(ctx),
								},
							},
						},
						Computed:            true,
						Description:         "The volume(s) which were created by the request.",
						MarkdownDescription: "The volume(s) which were created by the request.",
					},
				},
				CustomType: DataTypeVolume{
					ObjectType: types.ObjectType{
						AttrTypes: DataValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
			"details": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Whether to include details about the volume in the response.  By\ndefault this is set to true, meaning that all information about the\nvolume will be included in the response.  If set to false, only the\nbasic information about the volume will be included, such as its name\nand UUID.",
				MarkdownDescription: "Whether to include details about the volume in the response.  By\ndefault this is set to true, meaning that all information about the\nvolume will be included in the response.  If set to false, only the\nbasic information about the volume will be included, such as its name\nand UUID.",
			},
			"message": schema.StringAttribute{
				Computed:            true,
				Description:         "An optional message providing additional information about the response.",
				MarkdownDescription: "An optional message providing additional information about the response.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The name of the volume.\n\nThis is a human-readable name that can be used to identify the volume. The name must be unique within the context of your account. If no name is specified, a random name of the form `vol-X` is generated for you, where `X` is a 5-character long random alphanumeric suffix. The name can also be used to identify the volume in API calls.",
				MarkdownDescription: "The name of the volume.\n\nThis is a human-readable name that can be used to identify the volume. The name must be unique within the context of your account. If no name is specified, a random name of the form `vol-X` is generated for you, where `X` is a 5-character long random alphanumeric suffix. The name can also be used to identify the volume in API calls.",
			},
			"op_time_us": schema.Int64Attribute{
				Computed:            true,
				Description:         "The operation time in microseconds.  This is the time it took to process\nthe request and generate the response.",
				MarkdownDescription: "The operation time in microseconds.  This is the time it took to process\nthe request and generate the response.",
			},
			"size_mb": schema.Int64Attribute{
				Required:            true,
				Description:         "The size of the volume in megabytes.",
				MarkdownDescription: "The size of the volume in megabytes.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The status of the response.",
				MarkdownDescription: "The status of the response.",
			},
		},
	}
}

// Metadata implements resource.Resource.
func (r *VolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Configure implements resource.Resource.
func (r *VolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Handle the map of clients from the provider
	clients, ok := req.ProviderData.(map[string]any)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected map[string]any, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	volumesClient, exists := clients["volumes"]
	if !exists {
		resp.Diagnostics.AddError(
			"Missing Volumes Client",
			"Volumes client not found in provider data",
		)
		return
	}

	r.client, ok = volumesClient.(volumes.VolumesService)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Volumes Client Type",
			fmt.Sprintf("Expected volumes.VolumesService, got: %T", volumesClient),
		)
		return
	}
}

// Create implements resource.Resource.
func (r *VolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VolumeModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API with name and sizeMB parameters
	name := data.Name.ValueString()
	sizeMB := int(data.SizeMb.ValueInt64())

	volRaw, err := r.client.Create(ctx, name, sizeMB)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create volume, got error: %v", err),
		)
		return
	}

	// The API returns a ServiceResponse[CreateResponseItem]
	// We need to check if we have a valid response
	if volRaw == nil {
		resp.Diagnostics.AddError(
			"API Error",
			"No response returned from create operation",
		)
		return
	}

	// Check if we have entries in the response
	if len(volRaw.Data.Entries) == 0 {
		resp.Diagnostics.AddError(
			"API Error",
			"No volume returned from create operation",
		)
		return
	}

	// Get the first (and should be only) entry
	createItem := volRaw.Data.Entries[0]

	// Set the basic fields from create response
	data.Name = types.StringValue(createItem.Name)
	data.Status = types.StringValue(volRaw.Status)
	data.Message = types.StringValue(volRaw.Message)
	data.OpTimeUs = types.Int64Null() // OpTimeUs not available in this response

	// Create the data structure with volumes list using the CreateResponseItem
	volumesList := []VolumesValue{
		{
			Name:   types.StringValue(createItem.Name),
			Status: types.StringValue(createItem.Status),
			UUID:   types.StringValue(createItem.UUID),
			state:  attr.ValueStateKnown,
		},
	}

	var diags diag.Diagnostics
	volumesListValue, diags := types.ListValueFrom(ctx, VolumesType{
		ObjectType: types.ObjectType{
			AttrTypes: VolumesValue{}.AttributeTypes(ctx),
		},
	}, volumesList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Data = DataValue{
		Volumes: volumesListValue,
		state:   attr.ValueStateKnown,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource.Resource.
func (r *VolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VolumeModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var volumeUUID string
	if !data.Data.Volumes.IsNull() && !data.Data.Volumes.IsUnknown() {
		elements := data.Data.Volumes.Elements()
		if len(elements) > 0 {
			volumeValue, ok := elements[0].(VolumesValue)
			if ok {
				volumeUUID = volumeValue.UUID.ValueString()
			}
		}
	}

	if volumeUUID == "" {
		resp.Diagnostics.AddError(
			"Invalid State",
			"Volume UUID not found in state",
		)
		return
	}

	volRaw, err := r.client.Get(ctx, volumeUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to read volume, got error: %v", err),
		)
		return
	}

	if len(volRaw.Data.Entries) == 0 {
		// Volume no longer exists, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	vols := volRaw.Data.Entries[0]

	// Update the model with current state
	data.Name = types.StringValue(vols.Name)
	data.Status = types.StringValue(volRaw.Status)
	data.Message = types.StringValue(volRaw.Message)
	data.OpTimeUs = types.Int64Null() // OpTimeUs not available in this response

	// Update the data structure with current volume info
	volumesList := []VolumesValue{
		{
			Name:   types.StringValue(vols.Name),
			Status: types.StringValue(vols.Status),
			UUID:   types.StringValue(vols.UUID),
			state:  attr.ValueStateKnown,
		},
	}

	var diags diag.Diagnostics
	volumesListValue, diags := types.ListValueFrom(ctx, VolumesType{
		ObjectType: types.ObjectType{
			AttrTypes: VolumesValue{}.AttributeTypes(ctx),
		},
	}, volumesList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Data = DataValue{
		Volumes: volumesListValue,
		state:   attr.ValueStateKnown,
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (r *VolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unsupported",
		"This resource does not support updates. Configuration changes were expected to have triggered a replacement "+
			"of the resource. Please report this issue to the provider developers.",
	)
}

// Delete implements resource.Resource.
func (r *VolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VolumeModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the UUID from the state data
	var volumeUUID string
	if !data.Data.Volumes.IsNull() && !data.Data.Volumes.IsUnknown() {
		elements := data.Data.Volumes.Elements()
		if len(elements) > 0 {
			volumeValue, ok := elements[0].(VolumesValue)
			if ok {
				volumeUUID = volumeValue.UUID.ValueString()
			}
		}
	}

	if volumeUUID == "" {
		resp.Diagnostics.AddError(
			"Invalid State",
			"Volume UUID not found in state",
		)
		return
	}

	_, err := r.client.Delete(ctx, volumeUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete volume, got error: %v", err),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *VolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

type VolumeModel struct {
	Data     DataValue    `tfsdk:"data"`
	Details  types.Bool   `tfsdk:"details"`
	Message  types.String `tfsdk:"message"`
	Name     types.String `tfsdk:"name"`
	OpTimeUs types.Int64  `tfsdk:"op_time_us"`
	SizeMb   types.Int64  `tfsdk:"size_mb"`
	Status   types.String `tfsdk:"status"`
}

var _ basetypes.ObjectTypable = DataTypeVolume{}

type DataTypeVolume struct {
	basetypes.ObjectType
}

func (t DataTypeVolume) Equal(o attr.Type) bool {
	other, ok := o.(DataTypeVolume)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t DataTypeVolume) String() string {
	return "DataTypeVolume"
}

func (t DataTypeVolume) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	volumesAttribute, ok := attributes["volumes"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`volumes is missing from object`)

		return nil, diags
	}

	volumesVal, ok := volumesAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`volumes expected to be basetypes.ListValue, was: %T`, volumesAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return DataValue{
		Volumes: volumesVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewDataValueNull() DataValue {
	return DataValue{
		state: attr.ValueStateNull,
	}
}

func NewDataValueUnknown() DataValue {
	return DataValue{
		state: attr.ValueStateUnknown,
	}
}

func NewDataValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (DataValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing DataValue Attribute Value",
				"While creating a DataValue value, a missing attribute value was detected. "+
					"A DataValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("DataValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid DataValue Attribute Type",
				"While creating a DataValue value, an invalid attribute value was detected. "+
					"A DataValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("DataValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("DataValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra DataValue Attribute Value",
				"While creating a DataValue value, an extra attribute value was detected. "+
					"A DataValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra DataValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewDataValueUnknown(), diags
	}

	volumesAttribute, ok := attributes["volumes"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`volumes is missing from object`)

		return NewDataValueUnknown(), diags
	}

	volumesVal, ok := volumesAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`volumes expected to be basetypes.ListValue, was: %T`, volumesAttribute))
	}

	if diags.HasError() {
		return NewDataValueUnknown(), diags
	}

	return DataValue{
		Volumes: volumesVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewDataValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) DataValue {
	object, diags := NewDataValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewDataValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t DataTypeVolume) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewDataValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewDataValueUnknown(), nil
	}

	if in.IsNull() {
		return NewDataValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)
	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewDataValueMust(DataValue{}.AttributeTypes(ctx), attributes), nil
}

func (t DataTypeVolume) ValueType(ctx context.Context) attr.Value {
	return DataValue{}
}

var _ basetypes.ObjectValuable = DataValue{}

type DataValue struct {
	Volumes basetypes.ListValue `tfsdk:"volumes"`
	state   attr.ValueState
}

func (v DataValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 1)

	var val tftypes.Value
	var err error

	attrTypes["volumes"] = basetypes.ListType{
		ElemType: VolumesValue{}.Type(ctx),
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 1)

		val, err = v.Volumes.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["volumes"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v DataValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v DataValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v DataValue) String() string {
	return "DataValue"
}

func (v DataValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	volumes := types.ListValueMust(
		VolumesType{
			basetypes.ObjectType{
				AttrTypes: VolumesValue{}.AttributeTypes(ctx),
			},
		},
		v.Volumes.Elements(),
	)

	if v.Volumes.IsNull() {
		volumes = types.ListNull(
			VolumesType{
				basetypes.ObjectType{
					AttrTypes: VolumesValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.Volumes.IsUnknown() {
		volumes = types.ListUnknown(
			VolumesType{
				basetypes.ObjectType{
					AttrTypes: VolumesValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	attributeTypes := map[string]attr.Type{
		"volumes": basetypes.ListType{
			ElemType: VolumesValue{}.Type(ctx),
		},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"volumes": volumes,
		})

	return objVal, diags
}

func (v DataValue) Equal(o attr.Value) bool {
	other, ok := o.(DataValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Volumes.Equal(other.Volumes) {
		return false
	}

	return true
}

func (v DataValue) Type(ctx context.Context) attr.Type {
	return DataTypeVolume{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v DataValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"volumes": basetypes.ListType{
			ElemType: VolumesValue{}.Type(ctx),
		},
	}
}

var _ basetypes.ObjectTypable = VolumesType{}

type VolumesType struct {
	basetypes.ObjectType
}

func (t VolumesType) Equal(o attr.Type) bool {
	other, ok := o.(VolumesType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t VolumesType) String() string {
	return "VolumesType"
}

func (t VolumesType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	statusAttribute, ok := attributes["status"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`status is missing from object`)

		return nil, diags
	}

	statusVal, ok := statusAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`status expected to be basetypes.StringValue, was: %T`, statusAttribute))
	}

	uuidAttribute, ok := attributes["uuid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`uuid is missing from object`)

		return nil, diags
	}

	uuidVal, ok := uuidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`uuid expected to be basetypes.StringValue, was: %T`, uuidAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return VolumesValue{
		Name:   nameVal,
		Status: statusVal,
		UUID:   uuidVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewVolumesValueNull() VolumesValue {
	return VolumesValue{
		state: attr.ValueStateNull,
	}
}

func NewVolumesValueUnknown() VolumesValue {
	return VolumesValue{
		state: attr.ValueStateUnknown,
	}
}

func NewVolumesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (VolumesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing VolumesValue Attribute Value",
				"While creating a VolumesValue value, a missing attribute value was detected. "+
					"A VolumesValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("VolumesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid VolumesValue Attribute Type",
				"While creating a VolumesValue value, an invalid attribute value was detected. "+
					"A VolumesValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("VolumesValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("VolumesValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra VolumesValue Attribute Value",
				"While creating a VolumesValue value, an extra attribute value was detected. "+
					"A VolumesValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra VolumesValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewVolumesValueUnknown(), diags
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewVolumesValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	statusAttribute, ok := attributes["status"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`status is missing from object`)

		return NewVolumesValueUnknown(), diags
	}

	statusVal, ok := statusAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`status expected to be basetypes.StringValue, was: %T`, statusAttribute))
	}

	uuidAttribute, ok := attributes["uuid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`uuid is missing from object`)

		return NewVolumesValueUnknown(), diags
	}

	uuidVal, ok := uuidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`uuid expected to be basetypes.StringValue, was: %T`, uuidAttribute))
	}

	if diags.HasError() {
		return NewVolumesValueUnknown(), diags
	}

	return VolumesValue{
		Name:   nameVal,
		Status: statusVal,
		UUID:   uuidVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewVolumesValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) VolumesValue {
	object, diags := NewVolumesValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewVolumesValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t VolumesType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewVolumesValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewVolumesValueUnknown(), nil
	}

	if in.IsNull() {
		return NewVolumesValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)
	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewVolumesValueMust(VolumesValue{}.AttributeTypes(ctx), attributes), nil
}

func (t VolumesType) ValueType(ctx context.Context) attr.Value {
	return VolumesValue{}
}

var _ basetypes.ObjectValuable = VolumesValue{}

type VolumesValue struct {
	Name   basetypes.StringValue `tfsdk:"name"`
	Status basetypes.StringValue `tfsdk:"status"`
	UUID   basetypes.StringValue `tfsdk:"uuid"`
	state  attr.ValueState
}

func (v VolumesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["status"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["uuid"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.Name.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		val, err = v.Status.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["status"] = val

		val, err = v.UUID.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["uuid"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v VolumesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v VolumesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v VolumesValue) String() string {
	return "VolumesValue"
}

func (v VolumesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"name":   basetypes.StringType{},
		"status": basetypes.StringType{},
		"uuid":   basetypes.StringType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"name":   v.Name,
			"status": v.Status,
			"uuid":   v.UUID,
		})

	return objVal, diags
}

func (v VolumesValue) Equal(o attr.Value) bool {
	other, ok := o.(VolumesValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.Status.Equal(other.Status) {
		return false
	}

	if !v.UUID.Equal(other.UUID) {
		return false
	}

	return true
}

func (v VolumesValue) Type(ctx context.Context) attr.Type {
	return VolumesType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v VolumesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"name":   basetypes.StringType{},
		"status": basetypes.StringType{},
		"uuid":   basetypes.StringType{},
	}
}
