// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
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

var ApiResources = []string{"Channels", "Domains", "Events", "Listeners", "Messages", "Queues", "Schemas", "Stats", "Topics"}

func NewApiKeyResource() resource.Resource {
	return &apiKeyResource{}
}

type apiKeyResource struct {
	client *client.Client
}

type apiKeyResourceModel struct {
	Key    types.String       `tfsdk:"key"`
	Id     types.String       `tfsdk:"id"`
	Alias  types.String       `tfsdk:"alias"`
	Status types.String       `tfsdk:"status"`
	Scope  []apiKeyScopeModel `tfsdk:"scope"`
}

type apiKeyScopeModel struct {
	Resource string       `tfsdk:"resource"`
	Access   types.String `tfsdk:"access"`
	Targets  types.List   `tfsdk:"targets"`
}

func (r *apiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
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
				Sensitive: false,
			},
		},
		Blocks: map[string]schema.Block{
			"scope": schema.ListNestedBlock{ // Use ListNestedBlock or SetNestedBlock
				NestedObject: schema.NestedBlockObject{
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

func convertScopes(ctx context.Context, plan apiKeyResourceModel) (*client.ApiKeyScopes, error) {
	result := client.ApiKeyScopes{}

	for _, scope := range plan.Scope {
		access := scope.Access.ValueString()
		targets, _ := TfTypesValueToList(scope.Targets)

		setValueOf(&result, scope.Resource, &client.ApiKeyScope{
			Access:  access,
			Targets: targets,
		})
	}

	return &result, nil
}

func (r *apiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	key := client.ApiKeyRequest{
		Alias:  plan.Alias.ValueString(),
		Status: plan.Status.ValueString(),
	}

	if len(plan.Scope) > 0 {
		scopes, convertErr := convertScopes(ctx, plan)
		if convertErr != nil {
			resp.Diagnostics.AddError(
				"Error Creating Api Key",
				"Could not create api key, unexpected error: "+convertErr.Error(),
			)
			return
		}
		key.Scopes = scopes
	}
	tflog.Info(ctx, fmt.Sprintf("API key to be created %#v", key))

	var cK, createErr = r.client.CreateApiKey(key)
	if createErr != nil {
		resp.Diagnostics.AddError(
			"Error Creating Api Key",
			"Could not create api key, unexpected error: "+createErr.Error(),
		)
		return
	}

	var k, gErr = r.client.GetApiKey(cK.Id)
	if gErr != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Api Key",
			"Could not read api key, unexpected error: "+gErr.Error(),
		)
		return
	}

	r.convert(k, &plan)
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

	var d, err = r.client.GetApiKey(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Api Key",
			"Could not read api key, unexpected error: "+err.Error(),
		)
		return
	}

	r.convert(d, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan apiKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)

	var state apiKeyResourceModel
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	key := client.ApiKeyRequest{
		Alias: plan.Alias.ValueString(),
	}

	var k, err = r.client.UpdateApiKey(state.Id.ValueString(), key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Api Key",
			fmt.Sprintf("Could not update api key %s, unexpected error: %s", plan.Id.ValueString(), err.Error()),
		)
		return
	}

	k, err = r.client.GetApiKey(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Created Api Key",
			"Could not read api key, unexpected error: "+err.Error(),
		)
		return
	}

	r.convert(k, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *apiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state apiKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _, err = r.client.DeleteApiKey(state.Id.ValueString())
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
