package terraxcel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &terraxcelProvider{}
)

type terraxcelProvider struct {
}

type terraxcelProviderModel struct {
}

func (p *terraxcelProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "terraXcel"
}

func (p *terraxcelProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
}
