package terraxcel

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &extensionsDataSource{}
	_ datasource.DataSourceWithConfigure = &extensionsDataSource{}
)

func NewExtensionsDataSource() datasource.DataSource {
	return &extensionsDataSource{}
}

type extensionsDataSource struct {
}

type extensionsDataSourceModel struct {
	Extensions []extensionModel `tfsdk:"extensions"`
}

type extensionModel struct {
	Extension types.String `tfsdk:"extension"`
}
