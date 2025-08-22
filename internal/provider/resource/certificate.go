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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"sdk.kraft.cloud/certificates"
)

func NewCertificateResource() resource.Resource {
	return &CertificateResource{}
}

// CertificateResource defines the resource implementation.
type CertificateResource struct {
	client certificates.CertificatesService
}

var (
	_ resource.Resource                = &CertificateResource{}
	_ resource.ResourceWithImportState = &CertificateResource{}
)

type CertificateResourceModel struct {
	Chain   types.String         `tfsdk:"chain"`
	Cn      types.String         `tfsdk:"cn"`
	Data    CertificateDataValue `tfsdk:"data"`
	Message types.String         `tfsdk:"message"`
	Name    types.String         `tfsdk:"name"`
	Pkey    types.String         `tfsdk:"pkey"`
	Status  types.String         `tfsdk:"status"`
	Uuid    types.String         `tfsdk:"uuid"`
}

func (r *CertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (r *CertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"chain": schema.StringAttribute{
				Required:            true,
				Description:         "The chain of the certificate.",
				MarkdownDescription: "The chain of the certificate.",
			},
			"cn": schema.StringAttribute{
				Required:            true,
				Description:         "The common name (CN) of the certificate.",
				MarkdownDescription: "The common name (CN) of the certificate.",
			},
			"data": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"certificates": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"common_name": schema.StringAttribute{
									Computed: true,
								},
								"created_at": schema.StringAttribute{
									Computed:            true,
									Description:         "The time the certificate was created.",
									MarkdownDescription: "The time the certificate was created.",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "The name of the certificate.\n\nThis is a human-readable name that can be used to identify the certificate.\nThe name must be unique within the context of your account.  The name can\nalso be used to identify the certificate in API calls.",
									MarkdownDescription: "The name of the certificate.\n\nThis is a human-readable name that can be used to identify the certificate.\nThe name must be unique within the context of your account.  The name can\nalso be used to identify the certificate in API calls.",
								},
								"state": schema.StringAttribute{
									Computed: true,
								},
								"uuid": schema.StringAttribute{
									Computed:            true,
									Description:         "The UUID of the certificate.\n\nThis is a unique identifier for the certificate that is generated when the\ncertificate is created.  The UUID is used to reference the certificate in\nAPI calls and can be used to identify the certificate in all API calls that\nrequire an identifier.",
									MarkdownDescription: "The UUID of the certificate.\n\nThis is a unique identifier for the certificate that is generated when the\ncertificate is created.  The UUID is used to reference the certificate in\nAPI calls and can be used to identify the certificate in all API calls that\nrequire an identifier.",
								},
							},
							CustomType: CertificatesType{
								ObjectType: types.ObjectType{
									AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
								},
							},
						},
						Computed:            true,
						Description:         "The certificate which was created by this request.\n\nNote: only one certificate can be specified in the request, so this\nwill always contain a single entry.",
						MarkdownDescription: "The certificate which was created by this request.\n\nNote: only one certificate can be specified in the request, so this\nwill always contain a single entry.",
					},
				},
				CustomType: DataType{
					ObjectType: types.ObjectType{
						AttrTypes: CertificateDataValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
			"message": schema.StringAttribute{
				Computed:            true,
				Description:         "An optional message providing additional information about the response.",
				MarkdownDescription: "An optional message providing additional information about the response.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The name of the certificate (optional).",
				MarkdownDescription: "The name of the certificate (optional).",
			},
			"pkey": schema.StringAttribute{
				Required:            true,
				Description:         "The private key of the certificate.",
				MarkdownDescription: "The private key of the certificate.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The status of the response.",
				MarkdownDescription: "The status of the response.",
			},
			"uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The UUID of the certificate.",
				MarkdownDescription: "The UUID of the certificate.",
			},
		},
	}
}

func (r *CertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	certificatesClient, exists := clients["certificates"]
	if !exists {
		resp.Diagnostics.AddError(
			"Missing Certificates Client",
			"Certificates client not found in provider data",
		)
		return
	}

	r.client, ok = certificatesClient.(certificates.CertificatesService)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Certificates Client Type",
			fmt.Sprintf("Expected certificates.CertificatesService, got: %T", certificatesClient),
		)
		return
	}
}

func (r *CertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CertificateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	crt := certificates.CreateRequest{
		Name:  data.Name.ValueString(),
		CN:    data.Cn.ValueString(),
		Chain: data.Chain.ValueString(),
		PKey:  data.Pkey.ValueString(),
	}

	crtRaw, err := r.client.Create(ctx, &crt)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create certificate, got error: %v", err),
		)
		return
	}

	if len(crtRaw.Data.Entries) == 0 {
		resp.Diagnostics.AddError(
			"API Error",
			"No certificate returned from create operation",
		)
		return
	}

	crts := crtRaw.Data.Entries[0]

	// Set the basic fields from create response (CreateResponseItem only has UUID and Name)
	data.Uuid = types.StringValue(crts.UUID)
	data.Name = types.StringValue(crts.Name)

	// Set response-level fields
	data.Status = types.StringValue(crtRaw.Status)
	data.Message = types.StringValue(crtRaw.Message)

	// Get full certificate details since CreateResponseItem has limited fields
	crtFullRaw, err := r.client.Get(ctx, data.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get full certificate details, got error: %v", err),
		)
		return
	}

	if len(crtFullRaw.Data.Entries) == 0 {
		resp.Diagnostics.AddError(
			"API Error",
			"Certificate not found after creation",
		)
		return
	}

	crtFull := crtFullRaw.Data.Entries[0]

	// Create the data structure with certificates list using GetResponseItem fields
	certificatesList := []CertificatesValue{
		{
			CommonName: types.StringValue(crtFull.CommonName),
			CreatedAt:  types.StringValue(crtFull.CreatedAt),
			Name:       types.StringValue(crtFull.Name),
			State:      types.StringValue(crtFull.State),
			Uuid:       types.StringValue(crtFull.UUID),
			state:      attr.ValueStateKnown,
		},
	}

	var diags diag.Diagnostics
	certificatesListValue, diags := types.ListValueFrom(ctx, CertificatesType{
		ObjectType: types.ObjectType{
			AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
		},
	}, certificatesList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Data = CertificateDataValue{
		Certificates: certificatesListValue,
		state:        attr.ValueStateKnown,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CertificateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state from the API
	crtRaw, err := r.client.Get(ctx, data.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to read certificate, got error: %v", err),
		)
		return
	}

	// Check if certificate exists - Fix: Use Entries instead of Certificates
	if len(crtRaw.Data.Entries) == 0 {
		// Certificate no longer exists, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	crts := crtRaw.Data.Entries[0]

	// Update the model with current state using GetResponseItem fields
	data.Uuid = types.StringValue(crts.UUID)
	data.Name = types.StringValue(crts.Name)

	// Set response-level fields
	data.Status = types.StringValue(crtRaw.Status)
	data.Message = types.StringValue(crtRaw.Message)

	// Note: We typically don't update sensitive fields like chain and pkey from Read
	// These should remain as they were set in the configuration

	// Update the data structure with current certificate info
	certificatesList := []CertificatesValue{
		{
			CommonName: types.StringValue(crts.CommonName),
			CreatedAt:  types.StringValue(crts.CreatedAt),
			Name:       types.StringValue(crts.Name),
			State:      types.StringValue(crts.State),
			Uuid:       types.StringValue(crts.UUID),
			state:      attr.ValueStateKnown,
		},
	}

	var diags diag.Diagnostics
	certificatesListValue, diags := types.ListValueFrom(ctx, CertificatesType{
		ObjectType: types.ObjectType{
			AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
		},
	}, certificatesList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Data = CertificateDataValue{
		Certificates: certificatesListValue,
		state:        attr.ValueStateKnown,
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unsupported",
		"This resource does not support updates. Configuration changes were expected to have triggered a replacement "+
			"of the resource. Please report this issue to the provider developers.",
	)
}

func (r *CertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CertificateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Delete(ctx, data.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete certificate, got error: %v", err),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *CertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

var _ basetypes.ObjectTypable = DataType{}

type DataType struct {
	basetypes.ObjectType
}

func (t DataType) Equal(o attr.Type) bool {
	other, ok := o.(DataType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t DataType) String() string {
	return "DataType"
}

func (t DataType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	certificatesAttribute, ok := attributes["certificates"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`certificates is missing from object`)

		return nil, diags
	}

	certificatesVal, ok := certificatesAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`certificates expected to be basetypes.ListValue, was: %T`, certificatesAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CertificateDataValue{
		Certificates: certificatesVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewCertificateDataValueNull() CertificateDataValue {
	return CertificateDataValue{
		state: attr.ValueStateNull,
	}
}

func NewCertificateDataValueUnknown() CertificateDataValue {
	return CertificateDataValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCertificateDataValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CertificateDataValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CertificateDataValue Attribute Value",
				"While creating a CertificateDataValue value, a missing attribute value was detected. "+
					"A CertificateDataValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CertificateDataValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CertificateDataValue Attribute Type",
				"While creating a CertificateDataValue value, an invalid attribute value was detected. "+
					"A CertificateDataValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CertificateDataValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CertificateDataValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CertificateDataValue Attribute Value",
				"While creating a CertificateDataValue value, an extra attribute value was detected. "+
					"A CertificateDataValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CertificateDataValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCertificateDataValueUnknown(), diags
	}

	certificatesAttribute, ok := attributes["certificates"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`certificates is missing from object`)

		return NewCertificateDataValueUnknown(), diags
	}

	certificatesVal, ok := certificatesAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`certificates expected to be basetypes.ListValue, was: %T`, certificatesAttribute))
	}

	if diags.HasError() {
		return NewCertificateDataValueUnknown(), diags
	}

	return CertificateDataValue{
		Certificates: certificatesVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewCertificateDataValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CertificateDataValue {
	object, diags := NewCertificateDataValue(attributeTypes, attributes)

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

		panic("NewCertificateDataValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t DataType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCertificateDataValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCertificateDataValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCertificateDataValueNull(), nil
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

	return NewCertificateDataValueMust(CertificateDataValue{}.AttributeTypes(ctx), attributes), nil
}

func (t DataType) ValueType(ctx context.Context) attr.Value {
	return CertificateDataValue{}
}

var _ basetypes.ObjectValuable = CertificateDataValue{}

type CertificateDataValue struct {
	Certificates basetypes.ListValue `tfsdk:"certificates"`
	state        attr.ValueState
}

func (v CertificateDataValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 1)

	var val tftypes.Value
	var err error

	attrTypes["certificates"] = basetypes.ListType{
		ElemType: CertificatesValue{}.Type(ctx),
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 1)

		val, err = v.Certificates.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["certificates"] = val

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

func (v CertificateDataValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CertificateDataValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CertificateDataValue) String() string {
	return "CertificateDataValue"
}

func (v CertificateDataValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	certificates := types.ListValueMust(
		CertificatesType{
			basetypes.ObjectType{
				AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
			},
		},
		v.Certificates.Elements(),
	)

	if v.Certificates.IsNull() {
		certificates = types.ListNull(
			CertificatesType{
				basetypes.ObjectType{
					AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.Certificates.IsUnknown() {
		certificates = types.ListUnknown(
			CertificatesType{
				basetypes.ObjectType{
					AttrTypes: CertificatesValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	attributeTypes := map[string]attr.Type{
		"certificates": basetypes.ListType{
			ElemType: CertificatesValue{}.Type(ctx),
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
			"certificates": certificates,
		})

	return objVal, diags
}

func (v CertificateDataValue) Equal(o attr.Value) bool {
	other, ok := o.(CertificateDataValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Certificates.Equal(other.Certificates) {
		return false
	}

	return true
}

func (v CertificateDataValue) Type(ctx context.Context) attr.Type {
	return DataType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CertificateDataValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"certificates": basetypes.ListType{
			ElemType: CertificatesValue{}.Type(ctx),
		},
	}
}

var _ basetypes.ObjectTypable = CertificatesType{}

type CertificatesType struct {
	basetypes.ObjectType
}

func (t CertificatesType) Equal(o attr.Type) bool {
	other, ok := o.(CertificatesType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CertificatesType) String() string {
	return "CertificatesType"
}

func (t CertificatesType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	commonNameAttribute, ok := attributes["common_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`common_name is missing from object`)

		return nil, diags
	}

	commonNameVal, ok := commonNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`common_name expected to be basetypes.StringValue, was: %T`, commonNameAttribute))
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return nil, diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

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

	stateAttribute, ok := attributes["state"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`state is missing from object`)

		return nil, diags
	}

	stateVal, ok := stateAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`state expected to be basetypes.StringValue, was: %T`, stateAttribute))
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

	return CertificatesValue{
		CommonName: commonNameVal,
		CreatedAt:  createdAtVal,
		Name:       nameVal,
		State:      stateVal,
		Uuid:       uuidVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewCertificatesValueNull() CertificatesValue {
	return CertificatesValue{
		state: attr.ValueStateNull,
	}
}

func NewCertificatesValueUnknown() CertificatesValue {
	return CertificatesValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCertificatesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CertificatesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CertificatesValue Attribute Value",
				"While creating a CertificatesValue value, a missing attribute value was detected. "+
					"A CertificatesValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CertificatesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CertificatesValue Attribute Type",
				"While creating a CertificatesValue value, an invalid attribute value was detected. "+
					"A CertificatesValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CertificatesValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CertificatesValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CertificatesValue Attribute Value",
				"While creating a CertificatesValue value, an extra attribute value was detected. "+
					"A CertificatesValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CertificatesValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCertificatesValueUnknown(), diags
	}

	commonNameAttribute, ok := attributes["common_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`common_name is missing from object`)

		return NewCertificatesValueUnknown(), diags
	}

	commonNameVal, ok := commonNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`common_name expected to be basetypes.StringValue, was: %T`, commonNameAttribute))
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return NewCertificatesValueUnknown(), diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewCertificatesValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	stateAttribute, ok := attributes["state"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`state is missing from object`)

		return NewCertificatesValueUnknown(), diags
	}

	stateVal, ok := stateAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`state expected to be basetypes.StringValue, was: %T`, stateAttribute))
	}

	uuidAttribute, ok := attributes["uuid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`uuid is missing from object`)

		return NewCertificatesValueUnknown(), diags
	}

	uuidVal, ok := uuidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`uuid expected to be basetypes.StringValue, was: %T`, uuidAttribute))
	}

	if diags.HasError() {
		return NewCertificatesValueUnknown(), diags
	}

	return CertificatesValue{
		CommonName: commonNameVal,
		CreatedAt:  createdAtVal,
		Name:       nameVal,
		State:      stateVal,
		Uuid:       uuidVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewCertificatesValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CertificatesValue {
	object, diags := NewCertificatesValue(attributeTypes, attributes)

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

		panic("NewCertificatesValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CertificatesType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCertificatesValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCertificatesValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCertificatesValueNull(), nil
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

	return NewCertificatesValueMust(CertificatesValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CertificatesType) ValueType(ctx context.Context) attr.Value {
	return CertificatesValue{}
}

var _ basetypes.ObjectValuable = CertificatesValue{}

type CertificatesValue struct {
	CommonName basetypes.StringValue `tfsdk:"common_name"`
	CreatedAt  basetypes.StringValue `tfsdk:"created_at"`
	Name       basetypes.StringValue `tfsdk:"name"`
	State      basetypes.StringValue `tfsdk:"state"`
	Uuid       basetypes.StringValue `tfsdk:"uuid"`
	state      attr.ValueState
}

func (v CertificatesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes["common_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["state"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["uuid"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.CommonName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["common_name"] = val

		val, err = v.CreatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_at"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		val, err = v.State.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["state"] = val

		val, err = v.Uuid.ToTerraformValue(ctx)

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

func (v CertificatesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CertificatesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CertificatesValue) String() string {
	return "CertificatesValue"
}

func (v CertificatesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"common_name": basetypes.StringType{},
		"created_at":  basetypes.StringType{},
		"name":        basetypes.StringType{},
		"state":       basetypes.StringType{},
		"uuid":        basetypes.StringType{},
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
			"common_name": v.CommonName,
			"created_at":  v.CreatedAt,
			"name":        v.Name,
			"state":       v.State,
			"uuid":        v.Uuid,
		})

	return objVal, diags
}

func (v CertificatesValue) Equal(o attr.Value) bool {
	other, ok := o.(CertificatesValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.CommonName.Equal(other.CommonName) {
		return false
	}

	if !v.CreatedAt.Equal(other.CreatedAt) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.State.Equal(other.State) {
		return false
	}

	if !v.Uuid.Equal(other.Uuid) {
		return false
	}

	return true
}

func (v CertificatesValue) Type(ctx context.Context) attr.Type {
	return CertificatesType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CertificatesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"common_name": basetypes.StringType{},
		"created_at":  basetypes.StringType{},
		"name":        basetypes.StringType{},
		"state":       basetypes.StringType{},
		"uuid":        basetypes.StringType{},
	}
}
