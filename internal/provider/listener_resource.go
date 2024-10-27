// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-discue/internal/client"
	v "terraform-provider-discue/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &listenerResource{}
var _ resource.ResourceWithConfigure = &listenerResource{}
var _ resource.ResourceWithImportState = &listenerResource{}

func NewListenerResource() resource.Resource {
	return &listenerResource{}
}

type listenerResource struct {
	client *client.Client
}

type ListenerResourceModel struct {
	Alias       types.String `tfsdk:"alias"`
	Id          types.String `tfsdk:"id"`
	QueueId     types.String `tfsdk:"queue_id"`
	LivenessUrl types.String `tfsdk:"liveness_url"`
	NotifyUrl   types.String `tfsdk:"notify_url"`
}

func (r *listenerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = strings.Join([]string{req.ProviderTypeName, "listener"}, "_")
}

func (r *listenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Listener resource",
		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Required:    true,
				Description: "The name/alias of the resource. This should be unique.",
				Validators: []validator.String{
					v.ValidResourceAlias(""),
				},
			},
			"liveness_url": schema.StringAttribute{
				Required:    true,
				Description: "The URL used to check whether the listener is still live. Depends on a `",
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(4, 253),
						v.ValidUrl(""),
					),
				},
			},
			"notify_url": schema.StringAttribute{
				Required:    true,
				Description: "The URL used to send messages to the listener.",
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.LengthBetween(4, 253),
						v.ValidUrl(""),
					),
				},
			},
			"queue_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the queue this listener will receive messages from.",
				Validators: []validator.String{
					v.ValidResourceId(""),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique id of the resource.",
				Validators: []validator.String{
					v.ValidResourceId(""),
				},
			},
		},
	}
}

func (r *listenerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *listenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := r.convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting listener to API model",
			"Could not convert listener, unexpected error: "+err.Error(),
		)
		return
	}

	d, err := r.client.CreateListener(plan.QueueId.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating listener via API",
			"Could not create listener, unexpected error: "+err.Error(),
		)
		return
	}

	d, err = r.client.GetListener(plan.QueueId.ValueString(), d.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading listener via API",
			"Could not read listener, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(d, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting listener received from API to internal model",
			"Could not convert listener, unexpected error: "+err.Error())

		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *listenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var d, err = r.client.GetListener(state.QueueId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading listener via API",
			"Could not read listener, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(d, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting listener received from API to internal model",
			"Could not convert listener, unexpected error: "+err.Error())
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *listenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)

	var state ListenerResourceModel
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := r.convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting listener to API model",
			"Could not convert listener, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.client.UpdateListener(state.QueueId.ValueString(), state.Id.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating listener via API",
			"Could not update listener, unexpected error: "+err.Error(),
		)
		return
	}

	d, err := r.client.GetListener(state.QueueId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading listener via API",
			"Could not read listener, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(d, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting listener received from API to internal model",
			"Could not convert listener, unexpected error: "+err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *listenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _, err = r.client.DeleteListener(state.QueueId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting listener via API",
			"Could not delete listener, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *listenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	// Split the import ID into queue_id and listener_id
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Unexpected Format of Import ID", fmt.Sprintf("Expected format: <queue_id>/<listener_id> and got %s", req.ID))
		return
	}

	queueId := parts[0]
	listenerId := parts[1]

	state := ListenerResourceModel{}
	state.Id = types.StringValue(listenerId)
	state.QueueId = types.StringValue(queueId)

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
