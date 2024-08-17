// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &queueResource{}
var _ resource.ResourceWithConfigure = &queueResource{}

func NewQueueResource() resource.Resource {
	return &queueResource{}
}

type queueResource struct {
	client *client.Client
}

type QueueResourceModel struct {
	Alias types.String `tfsdk:"alias"`
	Id    types.String `tfsdk:"id"`
}

func (r *queueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_queue"
}

func (r *queueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Queue resource",
		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name/alias of the resource. This should be unique.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9.\-\\/]{4,64}$`),
						"must match the pattern for string name/alias values",
					),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(21, 22),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^{useandom26T198340PX75pxJACKVERYMINDBUSHWOLFGQZbfghjklqvwyzrict-}[21]$`),
						"must match the pattern for string id values",
					),
				},
			},
		},
	}
}

func (r *queueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *queueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan QueueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queue := client.Queue{
		Id:    plan.Id.ValueString(),
		Alias: plan.Alias.ValueString(),
	}
	var q, err = r.client.CreateQueue(queue)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Queue",
			"Could not create queue, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(q.Id)
	plan.Alias = types.StringValue(q.Alias)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *queueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state QueueResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var q, err = r.client.GetQueue(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Queue",
			fmt.Sprintf("Could not read queue %s, unexpected error: %s", state.Id.ValueString(), err.Error()),
		)
		return
	}

	state.Id = types.StringValue(q.Id)
	state.Alias = types.StringValue(q.Alias)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *queueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan QueueResourceModel
	diags := req.Plan.Get(ctx, &plan)

	var state QueueResourceModel
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queue := client.Queue{
		Alias: plan.Alias.ValueString(),
	}

	var q, err = r.client.UpdateQueue(state.Id.ValueString(), queue)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Queue",
			fmt.Sprintf("Could not update queue %s, unexpected error: %s", plan.Id.ValueString(), err.Error()),
		)
		return
	}

	plan.Id = types.StringValue(q.Id)
	plan.Alias = types.StringValue(q.Alias)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *queueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state QueueResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _, err = r.client.DeleteQueue(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Queue",
			"Could not delete queue, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *queueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
