package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"unikraft.com/cloud/sdk/platform"
)

// GetClientFromProviderData fetches the Unikraft SDK client instance from
// Terraform's datasource.ConfigureRequest.ProviderData
// It returns an error in the diags variable, in case the provider data is empty
// or the UKC platform client could not be found.
func GetClientFromProviderData(providerData any) (client *platform.Client, diags diag.Diagnostics) {
	if providerData == nil {
		return nil, diags
	}

	client, ok := providerData.(*platform.Client)
	if !ok {
		diags.AddError(
			"Unexpected resource Configure type",
			fmt.Sprintf("Expected *platform.Client, got: %T. Please report this issue to the provider developers.", providerData),
		)
	}

	return client, diags
}
