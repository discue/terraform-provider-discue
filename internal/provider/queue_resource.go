// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-discue/internal/client"
	v "terraform-provider-discue/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &queueResource{}
	_ resource.ResourceWithConfigure   = &queueResource{}
	_ resource.ResourceWithImportState = &queueResource{}
)

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
	resp.TypeName = strings.Join([]string{req.ProviderTypeName, "queue"}, "_")
}

func (r *queueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A queue resource is a prerequisite for creating listeners. It acts as a container for messages, ensuring that they are delivered to the correct destination. Each queue can have multiple listeners associated with it, allowing for flexible message routing and distribution.",
		Description:         "Queue resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique id of the resource.",
				Validators: []validator.String{
					v.ValidResourceId(""),
				},
			},
			"alias": schema.StringAttribute{
				Required:    true,
				Description: "The name/alias of the resource. This should be unique.",
				Validators: []validator.String{
					v.ValidResourceAlias(""),
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

	payload, err := r.convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting queue to API model",
			"Could not convert queue, unexpected error: "+err.Error(),
		)
		return
	}

	q, err := r.client.CreateQueue(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating queue via API",
			"Could not create queue, unexpected error: "+err.Error(),
		)
		return
	}

	q, err = r.client.GetQueue(q.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading queue via API",
			"Could not read queue, unexpected error: "+err.Error(),
		)
		return
	}

	err = r.convertFromApiModel(q, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting queue received from API to internal model",
			"Could not convert queue, unexpected error: "+err.Error())

		return
	}

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
			"Error reading queue via API",
			"Could not read queue, unexpected error: "+err.Error(),
		)
		return
	}

	err = r.convertFromApiModel(q, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting queue received from API to internal model",
			"Could not convert queue, unexpected error: "+err.Error())

		return
	}

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

	payload, err := r.convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting queue to API model",
			"Could not convert queue, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.client.UpdateQueue(state.Id.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating queue via API",
			"Could not update queue, unexpected error: "+err.Error(),
		)
		return
	}

	q, err := r.client.GetQueue(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading queue via API",
			"Could not read queue, unexpected error: "+err.Error(),
		)
		return
	}

	err = r.convertFromApiModel(q, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting queue received from API to internal model",
			"Could not convert queue, unexpected error: "+err.Error())

		return
	}

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
			"Error deleting queue via API",
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
