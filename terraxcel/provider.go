package terraxcel

import (
	"context"
	"os"

	"github.com/Deathfireofdoom/terraxcel-client/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &terraxcelProvider{}
)

func NewProvider() provider.Provider {
	return &terraxcelProvider{}
}

type terraxcelProvider struct {
}

type terraxcelProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func (p *terraxcelProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "terraXcel"
}

func (p *terraxcelProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: false,
			},
			"token": schema.StringAttribute{
				Optional:  false,
				Sensitive: true,
			},
		},
	}
}

func (p *terraxcelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring TerraXcel client")

	// get provider data from config
	var config terraxcelProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating TerraXcel client")

	// check if user setup provider block or if default values should be used
	host := os.Getenv("TERRAXCEL_HOST")
	token := os.Getenv("TERRAXCEL_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// check if either provider block is configured or env var
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing TerraXcel host",
			"Missing host for TerraXcel-server, either configure provider-block or set TERRAXCEL_HOST environment variable",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing TerraXcel token",
			"Missing token for TerraXcel-server, either configure provider-block or set TERRAXCEL_TOKEN environment variable",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// setting up terraxcel-client
	clientConfig := client.ClientConfig{
		BaseURL:   host,
		AuthToken: token,
	}
	client, err := client.NewClient(&clientConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create TerraXcel API client",
			"An unexpected error occurred when creating the TerraXcel API client. "+
				"TerraXcel Client Error: "+err.Error(),
		)
		return
	}

	// make client available for resources that needs it
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *terraxcelProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p terraxcelProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
