package terraxcel

import (
	"context"

	"github.com/Deathfireofdoom/terraxcel-client/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	client *client.Client
}

type extensionsDataSourceModel struct {
	Extensions []extensionModel `tfsdk:"extensions"`
}

type extensionModel struct {
	Extension types.String `tfsdk:"extension"`
}

func (d *extensionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extensions"
}

func (d *extensionsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"extensions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"extension": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *extensionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state extensionsDataSourceModel

	extensions, err := d.client.ReadExtensions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read extensions",
			"Unabled to read extensions: %w",
		)
		return
	}

	// maps response from client to state
	for _, extension := range extensions {
		extensionState := extensionModel{
			Extension: types.StringValue(extension),
		}
		state.Extensions = append(state.Extensions, extensionState)
	}

	// set stat
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (d *extensionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.Client)
}
