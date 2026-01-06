package ukc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		ResourcesMap: map[string]*schema.Resource{},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureContextFunc: providerConfigureWithDefaultUserAgent,
	}
}

func providerConfigureWithDefaultUserAgent(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return ProviderConfigure(ctx, d)
}

func ProviderConfigure(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	return nil, diags
}

type ukcProvider struct {
}

var _ provider.Provider = New()

func New() provider.Provider {
	return &ukcProvider{}
}

func (p *ukcProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ukc"
	//resp.Version = config.Version
}

func (p *ukcProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {

}

func (p *ukcProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *ukcProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ukcProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
