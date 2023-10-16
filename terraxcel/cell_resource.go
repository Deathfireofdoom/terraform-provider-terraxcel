package terraxcel

import (
	"context"
	"time"

	"github.com/Deathfireofdoom/excel-client-go/pkg/models"
	"github.com/Deathfireofdoom/terraxcel-client/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &cellResource{}
	_ resource.ResourceWithConfigure = &cellResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewCellResource() resource.Resource {
	return &cellResource{}
}

type cellResource struct {
	client *client.Client
}

type cellResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	WorkbookID  types.String `tfsdk:"workbook_id"`
	SheetID     types.String `tfsdk:"sheet_id"`
	Row         types.Int64  `tfsdk:"row"`
	Value       types.String `tfsdk:"value"`
	Column      types.String `tfsdk:"column"`
}

// Metadata returns the resource type name.
func (r *cellResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cell"
}

// Schema defines the schema for the resource.
func (r *cellResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"workbook_id": schema.StringAttribute{
				Required: true,
			},
			"sheet_id": schema.StringAttribute{
				Required: true,
			},
			"row": schema.Int64Attribute{
				Computed: true,
			},
			"column": schema.StringAttribute{
				Required: true,
			},
			"value": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *cellResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// creates the model, and populates it with values from the plan
	var plan cellResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// creates cell object from plan
	cell := &models.Cell{
		Row:        int(plan.Row.ValueInt64()),
		Column:     plan.Column.ValueString(),
		Value:      plan.Value.ValueString(),
		WorkbookID: plan.WorkbookID.ValueString(),
		SheetID:    plan.SheetID.ValueString(),
	}

	// creates the cell with help of the client
	cell, err := r.client.CreateCell(cell)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create cell: %s",
			err.Error(),
		)
		return
	}

	// maps the values we got from the client to the terraform model
	plan.ID = types.StringValue(cell.ID)
	plan.Row = types.Int64Value(int64(cell.Row))
	plan.Column = types.StringValue(cell.Column)

	// type assert value to string
	plan.Value = types.StringValue("test") // TODO FIX THIS

	// updates last_updated
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// sets the state with the populated model
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *cellResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cellResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed sheet value from client
	cell, err := r.client.ReadCell(state.WorkbookID.ValueString(), state.SheetID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading cell",
			"Could not read cell with ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(cell.ID)
	state.Row = types.Int64Value(int64(cell.Row))
	state.Column = types.StringValue(cell.Column)
	state.Value = types.StringValue("test")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *cellResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state cellResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create cell object from state
	cell := &models.Cell{
		ID:         state.ID.ValueString(),
		Row:        int(state.Row.ValueInt64()),
		Column:     state.Column.ValueString(),
		Value:      state.Value.ValueString(),
		WorkbookID: state.WorkbookID.ValueString(),
		SheetID:    state.SheetID.ValueString(),
	}

	// Delete existing order
	err := r.client.DeleteCell(cell)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting cell",
			"Could not delete cell, unexpected error: "+err.Error(),
		)
		return
	}
}

// update the workbook
func (r *cellResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// old state
	var state cellResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Retrieve values from plan
	var plan cellResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Converts tf-workbook-model to excel.Sheet
	cell := &models.Cell{
		ID:         state.ID.ValueString(),
		Row:        int(plan.Row.ValueInt64()),
		Column:     plan.Column.ValueString(),
		Value:      plan.Value.ValueString(),
		WorkbookID: state.WorkbookID.ValueString(),
		SheetID:    state.SheetID.ValueString(),
	}

	// Update existing order
	_, err := r.client.UpdateCell(cell)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating cell",
			"Could not update cell, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	cell, err = r.client.ReadCell(plan.WorkbookID.ValueString(), state.SheetID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading cell",
			"Could not read cell ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	plan.ID = types.StringValue(cell.ID)
	plan.Row = types.Int64Value(int64(cell.Row))
	plan.Column = types.StringValue(cell.Column)
	plan.Value = types.StringValue("test")

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *cellResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Client)
}
