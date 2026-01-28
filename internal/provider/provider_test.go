// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "release version",
			version: "1.0.0",
		},
		{
			name:    "dev version",
			version: "dev",
		},
		{
			name:    "test version",
			version: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := New(tt.version)
			assert.NotNil(t, factory)

			p := factory()
			assert.NotNil(t, p)
			assert.IsType(t, &UnikraftCloudProvider{}, p)

			ukcp := p.(*UnikraftCloudProvider)
			assert.Equal(t, tt.version, ukcp.version)
		})
	}
}

func TestUnikraftCloudProvider_Metadata(t *testing.T) {
	tests := []struct {
		name            string
		version         string
		expectedType    string
		expectedVersion string
	}{
		{
			name:            "returns correct metadata",
			version:         "test-version",
			expectedType:    "ukc",
			expectedVersion: "test-version",
		},
		{
			name:            "dev version",
			version:         "dev",
			expectedType:    "ukc",
			expectedVersion: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UnikraftCloudProvider{version: tt.version}
			req := provider.MetadataRequest{}
			resp := &provider.MetadataResponse{}

			p.Metadata(context.Background(), req, resp)

			assert.Equal(t, tt.expectedType, resp.TypeName)
			assert.Equal(t, tt.expectedVersion, resp.Version)
		})
	}
}

func TestUnikraftCloudProvider_Schema(t *testing.T) {
	p := &UnikraftCloudProvider{}
	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p.Schema(context.Background(), req, resp)

	assert.NotNil(t, resp.Schema)
	assert.Equal(t, "Manage unikernel instances on Unikraft Cloud.", resp.Schema.MarkdownDescription)

	expectedAttributes := []string{"metro", "token"}
	for _, attr := range expectedAttributes {
		t.Run("has_"+attr+"_attribute", func(t *testing.T) {
			attrSchema, ok := resp.Schema.Attributes[attr]
			assert.True(t, ok, "%s attribute should exist", attr)
			assert.NotNil(t, attrSchema)
		})
	}
}

func TestUnikraftCloudProvider_Resources(t *testing.T) {
	p := &UnikraftCloudProvider{}
	resources := p.Resources(context.Background())

	assert.Len(t, resources, 2)
	for i, rf := range resources {
		t.Run("resource_factory_"+string(rune('0'+i)), func(t *testing.T) {
			r := rf()
			assert.NotNil(t, r)
		})
	}
}

func TestUnikraftCloudProvider_DataSources(t *testing.T) {
	p := &UnikraftCloudProvider{}
	dataSources := p.DataSources(context.Background())

	assert.Len(t, dataSources, 2)
	for i, dsf := range dataSources {
		t.Run("datasource_factory_"+string(rune('0'+i)), func(t *testing.T) {
			ds := dsf()
			assert.NotNil(t, ds)
		})
	}
}

func TestUnikraftCloudProvider_Configure(t *testing.T) {
	tests := []struct {
		name           string
		envMetro       string
		envToken       string
		configMetro    any // string, nil (null), or tftypes.UnknownValue
		configToken    any
		expectError    bool
		errorSummaries []string
	}{
		{
			name:        "success with config values",
			envMetro:    "",
			envToken:    "",
			configMetro: "fra0",
			configToken: "test-token",
			expectError: false,
		},
		{
			name:        "success with env vars",
			envMetro:    "fra0",
			envToken:    "env-token",
			configMetro: nil,
			configToken: nil,
			expectError: false,
		},
		{
			name:        "config overrides env vars",
			envMetro:    "env-metro",
			envToken:    "env-token",
			configMetro: "config-metro",
			configToken: "config-token",
			expectError: false,
		},
		{
			name:           "missing metro",
			envMetro:       "",
			envToken:       "test-token",
			configMetro:    nil,
			configToken:    nil,
			expectError:    true,
			errorSummaries: []string{"Missing Unikraft Cloud API Metro"},
		},
		{
			name:           "missing token",
			envMetro:       "fra0",
			envToken:       "",
			configMetro:    nil,
			configToken:    nil,
			expectError:    true,
			errorSummaries: []string{"Missing Unikraft Cloud API Token"},
		},
		{
			name:           "missing both metro and token",
			envMetro:       "",
			envToken:       "",
			configMetro:    nil,
			configToken:    nil,
			expectError:    true,
			errorSummaries: []string{"Missing Unikraft Cloud API Metro", "Missing Unikraft Cloud API Token"},
		},
		{
			name:           "unknown metro",
			envMetro:       "",
			envToken:       "",
			configMetro:    tftypes.UnknownValue,
			configToken:    "test-token",
			expectError:    true,
			errorSummaries: []string{"Unknown Unikraft Cloud API Metro"},
		},
		{
			name:           "unknown token",
			envMetro:       "fra0",
			envToken:       "",
			configMetro:    "fra0",
			configToken:    tftypes.UnknownValue,
			expectError:    true,
			errorSummaries: []string{"Unknown Unikraft Cloud API Token"},
		},
		{
			name:           "both unknown",
			envMetro:       "",
			envToken:       "",
			configMetro:    tftypes.UnknownValue,
			configToken:    tftypes.UnknownValue,
			expectError:    true,
			errorSummaries: []string{"Unknown Unikraft Cloud API Metro", "Unknown Unikraft Cloud API Token"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore environment variables
			originalMetro := os.Getenv("UKC_METRO")
			originalToken := os.Getenv("UKC_TOKEN")
			defer func() {
				restoreEnv("UKC_METRO", originalMetro)
				restoreEnv("UKC_TOKEN", originalToken)
			}()

			// Set test environment
			setEnv("UKC_METRO", tt.envMetro)
			setEnv("UKC_TOKEN", tt.envToken)

			p := &UnikraftCloudProvider{}

			// Get the schema
			schemaReq := provider.SchemaRequest{}
			schemaResp := &provider.SchemaResponse{}
			p.Schema(context.Background(), schemaReq, schemaResp)

			// Build config value
			configValue := tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"metro": tftypes.String,
					"token": tftypes.String,
				},
			}, map[string]tftypes.Value{
				"metro": toTFValue(tt.configMetro),
				"token": toTFValue(tt.configToken),
			})

			req := provider.ConfigureRequest{
				Config: tfsdk.Config{
					Schema: schemaResp.Schema,
					Raw:    configValue,
				},
			}
			resp := &provider.ConfigureResponse{}

			p.Configure(context.Background(), req, resp)

			if tt.expectError {
				assert.True(t, resp.Diagnostics.HasError(), "Expected error but got none")
				for _, expectedSummary := range tt.errorSummaries {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if diag.Summary() == expectedSummary {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected error summary %q not found", expectedSummary)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError(), "Unexpected error: %v", resp.Diagnostics.Errors())
				assert.NotNil(t, resp.DataSourceData, "DataSourceData should be set")
				assert.NotNil(t, resp.ResourceData, "ResourceData should be set")
			}
		})
	}
}

func TestUnikraftCloudModel(t *testing.T) {
	tests := []struct {
		name        string
		metro       types.String
		token       types.String
		checkNull   bool
		checkUnknown bool
	}{
		{
			name:  "basic values",
			metro: types.StringValue("fra0"),
			token: types.StringValue("test-token"),
		},
		{
			name:      "null values",
			metro:     types.StringNull(),
			token:     types.StringNull(),
			checkNull: true,
		},
		{
			name:         "unknown values",
			metro:        types.StringUnknown(),
			token:        types.StringUnknown(),
			checkUnknown: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := UnikraftCloudModel{
				Metro: tt.metro,
				Token: tt.token,
			}

			if tt.checkNull {
				assert.True(t, model.Metro.IsNull())
				assert.True(t, model.Token.IsNull())
			} else if tt.checkUnknown {
				assert.True(t, model.Metro.IsUnknown())
				assert.True(t, model.Token.IsUnknown())
			} else {
				assert.Equal(t, tt.metro.ValueString(), model.Metro.ValueString())
				assert.Equal(t, tt.token.ValueString(), model.Token.ValueString())
			}
		})
	}
}

// Helper functions

func setEnv(key, value string) {
	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
}

func restoreEnv(key, value string) {
	if value == "" {
		os.Unsetenv(key)
	} else {
		os.Setenv(key, value)
	}
}

func toTFValue(v any) tftypes.Value {
	if v == nil {
		return tftypes.NewValue(tftypes.String, nil)
	}
	if v == tftypes.UnknownValue {
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	}
	if str, ok := v.(string); ok {
		return tftypes.NewValue(tftypes.String, str)
	}
	return tftypes.NewValue(tftypes.String, nil)
}
