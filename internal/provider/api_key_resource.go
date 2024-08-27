// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &apiKeyResource{}
var _ resource.ResourceWithConfigure = &apiKeyResource{}
var _ resource.ResourceWithImportState = &apiKeyResource{}

var ApiResources = []string{"channels", "domains", "events", "listeners", "messages", "queues", "schemas", "stats", "topics"}

func NewApiKeyResource() resource.Resource {
	return &apiKeyResource{}
}

type apiKeyResource struct {
	client *client.Client
}

type apiKeyResourceModel struct {
	Key    types.String `tfsdk:"key"`
	Id     types.String `tfsdk:"id"`
	Alias  types.String `tfsdk:"alias"`
	Status types.String `tfsdk:"status"`
	Scopes types.List   `tfsdk:"scopes"`
}

type apiKeyScopeModel struct {
	Resource types.String `tfsdk:"resource"`
	Access   types.String `tfsdk:"access"`
	Targets  types.List   `tfsdk:"targets"`
}

func (r *apiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = strings.Join([]string{req.ProviderTypeName, "api_key"}, "_")
}

func (r *apiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Api Key resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:  true,
				Sensitive: false,
				Validators: []validator.String{
					stringvalidator.LengthBetween(21, 22),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[useandom26T198340PX75pxJACKVERYMINDBUSHWOLFGQZbfghjklqvwyzrict-]{21}$`),
						"must match the pattern for string id values",
					),
				},
			},
			"alias": schema.StringAttribute{
				Required:            true,
				Sensitive:           false,
				MarkdownDescription: "The name/alias of the resource. This should be unique.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9.\-\\/]{4,64}$`),
						"must match the pattern for string name/alias values",
					),
				},
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Default:             stringdefault.StaticString("enabled"),
				MarkdownDescription: "The status of the api key. Default is\"enabled\".",
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
			},
			"key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"scopes": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.List{listvalidator.All(
					listvalidator.IsRequired(),
					listvalidator.SizeAtLeast(1),
				)},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(ApiResources...),
							},
						},
						"access": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("write"),
							Validators: []validator.String{
								stringvalidator.OneOf("read", "write"),
							},
						},
						"targets": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("*")})),
							Validators: []validator.List{
								listvalidator.ValueStringsAre(
									stringvalidator.Any(
										stringvalidator.OneOf("*"),
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^[useandom26T198340PX75pxJACKVERYMINDBUSHWOLFGQZbfghjklqvwyzrict-]{21}$`),
											"must match the pattern for string id values",
										),
									),
								),
							},
						},
					},
				},
			},
		},
	}
}

func (r *apiKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting Key",
			"Could not convert api key, unexpected error: "+err.Error(),
		)
		return
	}

	k, err := r.client.CreateApiKey(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Api Key",
			"Could not create api key, unexpected error: "+err.Error(),
		)
		return
	}

	k, err = r.client.GetApiKey(k.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Api Key",
			"Could not read api key, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(k, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting API response", "Unexpected error: "+err.Error())
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Done Reading api client %s", plan))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state apiKeyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	k, err := r.client.GetApiKey(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Api Key",
			"Could not read api key, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(k, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error converting API response", "Unexpected error: "+err.Error())
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Done Reading api client %s", state))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state apiKeyResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting Key",
			"Could not convert api key, unexpected error: "+err.Error(),
		)
		return
	}

	k, err := r.client.UpdateApiKey(state.Id.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Api Key",
			"Could not create api key, unexpected error: "+err.Error(),
		)
		return
	}

	k, err = r.client.GetApiKey(k.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Api Key",
			"Could not read api key, unexpected error: "+err.Error(),
		)
		return
	}
	_, err = r.convertFromApiModel(k, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error converting API response", "Unexpected error: "+err.Error())
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Done Reading api client %s", state))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteApiKey(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Api Key",
			"Could not delete api key, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
