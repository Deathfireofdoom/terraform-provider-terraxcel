package terraxcel

import (
	"context"
	"fmt"
	"time"

	"github.com/Deathfireofdoom/excel-client-go/pkg/models"
	"github.com/Deathfireofdoom/terraxcel-client/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &sheetResource{}
	_ resource.ResourceWithConfigure = &sheetResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewSheetResource() resource.Resource {
	return &sheetResource{}
}

type sheetResource struct {
	client *client.Client
}

type sheetResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	WorkbookID  types.String `tfsdk:"workbook_id"`
	Name        types.String `tfsdk:"name"`
	Pos         types.Int64  `tfsdk:"pos"`
}

// Metadata returns the resource type name.
func (r *sheetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sheet"
}

// Schema defines the schema for the resource.
func (r *sheetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"workbook_id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"pos": schema.Int64Attribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *sheetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// creates the model, and populates it with values from the plan
	var plan sheetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planSheet, err := models.NewSheet(plan.ID.ValueString(), int(plan.Pos.ValueInt64()), plan.Name.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create workbook",
			fmt.Sprintf("failed to create sheet: %s", err.Error()),
		)
	}

	// creates the sheet with help of the client
	sheet, err := r.client.CreateSheet(planSheet)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create workbook: %s",
			err.Error(),
		)
		return
	}

	// maps the values we got from the client to the terraform model
	plan.ID = types.StringValue(sheet.ID)
	plan.Name = types.StringValue(sheet.Name)
	plan.Pos = types.Int64Value(int64(sheet.Pos))

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
func (r *sheetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state sheetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed sheet value from client
	sheet, err := r.client.ReadSheet(state.WorkbookID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Sheet",
			"Could not read sheet with ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(sheet.ID)
	state.Name = types.StringValue(sheet.Name)
	state.Pos = types.Int64Value(int64(sheet.Pos))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sheetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sheetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// create sheet from state
	sheet := &models.Sheet{
		ID:         state.ID.ValueString(),
		WorkbookID: state.WorkbookID.ValueString(),
		Name:       state.Name.ValueString(),
		Pos:        int(state.Pos.ValueInt64()),
	}

	// Delete existing order
	err := r.client.DeleteSheet(sheet)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting sheet",
			"Could not delete sheet, unexpected error: "+err.Error(),
		)
		return
	}
}

// update the workbook
func (r *sheetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// old state
	var state sheetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Retrieve values from plan
	var plan sheetResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Converts tf-workbook-model to excel.Sheet
	sheet := &models.Sheet{
		ID:         state.ID.ValueString(),
		WorkbookID: plan.WorkbookID.ValueString(),
		Name:       plan.Name.ValueString(),
		Pos:        int(plan.Pos.ValueInt64()),
	}

	// Update existing order
	_, err := r.client.UpdateSheet(sheet)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workbook",
			"Could not update workbook, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	sheet, err = r.client.ReadSheet(plan.WorkbookID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Sheet",
			"Could not read Sheet ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	plan.ID = types.StringValue(sheet.ID)
	plan.WorkbookID = state.WorkbookID
	plan.Name = types.StringValue(sheet.Name)
	plan.Pos = types.Int64Value(int64(sheet.Pos))

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *sheetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Client)
}
