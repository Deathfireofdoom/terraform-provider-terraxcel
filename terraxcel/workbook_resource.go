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
	_ resource.Resource              = &workbookResource
	_ resource.ResourceWithConfigure = &workbookResource
)

func NewWorkbookResource() resource.Resource {
	return &workbookResource{}
}

type workbookResource struct {
	client *client.Client
}

type workbookResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	FileName    types.String `tfsdk:"file_name"`
	Extension   types.String `tfsdk:"extension"`
	FolderPath  types.String `tfsdk:"folder_path"`
}

func (r *workbookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workbook"
}

func (r *workbookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"file_name": schema.StringAttribute{
				Required: true,
			},
			"folder_path": schema.StringAttribute{
				Required: true,
			},
			"extension": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *workbookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workbookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// convert plan to model
	newWorkbook, err := models.NewWorkbook(plan.FileName.ValueString(), models.Extension(plan.Extension.ValueString()), plan.FolderPath.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError("could not create workbook object from plan", fmt.Sprintf("could not create workbook object from plan, err: %w", err))
		return
	}

	workbook, err := r.client.CreateWorkbook(newWorkbook)
	if err != nil {
		resp.Diagnostics.AddError("could not create workbook-file", fmt.Sprintf("could not create workbook-file, err: %w", err))
		return
	}

	// map to state
	plan.ID = types.StringValue(workbook.ID)
	plan.FileName = types.StringValue(workbook.FileName)
	plan.Extension = types.StringValue(string(workbook.Extension))
	plan.FolderPath = types.StringValue(workbook.FolderPath)

	// update last updated at
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workbookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	workbook, err := r.client.ReadWorkbook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading workbook",
			fmt.Sprintf("could not read workbook with ID %s, err: %w", state.ID.ValueString(), err),
		)
		return
	}

	state.ID = types.StringValue(workbook.ID)
	state.FileName = types.StringValue(workbook.FileName)
	state.Extension = types.StringValue(string(workbook.Extension))
	state.FolderPath = types.StringValue(workbook.FolderPath)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// TODO make sure this is correct
func (r *workbookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkbook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"could not delete workbook",
			fmt.Sprintf("could not delete workbook with id %s, unexpected error: %w", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *workbookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan workbookResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create new workbook
	workbook := &models.Workbook{
		ID:         state.ID.ValueString(),
		FileName:   plan.FileName.ValueString(),
		Extension:  models.Extension(plan.Extension.ValueString()),
		FolderPath: plan.FolderPath.ValueString(),
	}

	// update existing order
	_, err := r.client.UpdateWorkbook(workbook)
	if err != nil {
		resp.Diagnostics.AddError(
			"error updating workbook",
			fmt.Sprintf("error when updating workbook with id %s, err: %w", workbook.ID, err),
		)
		return
	}

	// reading the current state of the workbook after the update
	workbook, err = r.client.ReadWorkbook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading updated workbook",
			fmt.Sprintf("error reading updated workbook with id %s, err: %w", state.ID.ValueString(), err),
		)
	}

	// update state
	plan.ID = types.StringValue(workbook.ID)
	plan.FileName = types.StringValue(workbook.FileName)
	plan.Extension = types.StringValue(string(workbook.Extension))
	plan.FolderPath = types.StringValue(workbook.FolderPath)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *workbookResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Client)
}
