// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"terraform-provider-discue/internal/client"
	v "terraform-provider-discue/internal/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = &domainResource{}
var _ resource.ResourceWithConfigure = &domainResource{}
var _ resource.ResourceWithImportState = &domainResource{}

func NewDomainResource() resource.Resource {
	return &domainResource{}
}

type domainResource struct {
	client *client.Client
}

type DomainResourceModel struct {
	Alias        types.String          `tfsdk:"alias"`
	Id           types.String          `tfsdk:"id"`
	Hostname     types.String          `tfsdk:"hostname"`
	Port         types.Int32           `tfsdk:"port"`
	Challenge    basetypes.ObjectValue `tfsdk:"challenge"`
	Verification basetypes.ObjectValue `tfsdk:"verification"`
}

type DomainChallenge struct {
	Https HttpDomainChallenge `tfsdk:"https"`
}

type HttpDomainChallenge struct {
	FileContent types.String `tfsdk:"file_content"`
	FileName    types.String `tfsdk:"file_name"`
	ContextPath types.String `tfsdk:"context_path"`
	CreatedAt   types.Int64  `tfsdk:"created_at"`
	ExpiresAt   types.Int64  `tfsdk:"expires_at"`
}

type DomainVerification struct {
	Verified   types.Bool  `tfsdk:"verified"`
	VerifiedAt types.Int64 `tfsdk:"verified_at"`
}

func (r *domainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = strings.Join([]string{req.ProviderTypeName, "domain"}, "_")
}

func (r *domainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A domain resource is a prerequisite for receiving messages. This is a security measure to prevent messages being sent through [discue.io](https://www.discue.io) without knowledge of the recipient. The domain configuration includes a hostname and port. Both values cannot be changed after creation. The API will return instructions on how to validate the domain as a response to the creation request.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique id of the resource.",
				Validators: []validator.String{
					v.ValidResourceId(""),
				},
			},
			"alias": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name/alias of the resource. This should be unique.",
				Validators: []validator.String{
					v.ValidResourceAlias(""),
				},
			},
			"hostname": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The target hostname that will receive messages from listeners and channels. Only provide the DNS hostname portion here. The protocol will be added by the API automatically.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 253),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]{1}[a-zA-Z0-9-\\.]{0,62}[a-zA-Z0-9]{1}$`),
						"must match the pattern for valid hostnames",
					),
				},
			},
			"port": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "The target port the messages will be sent to.",
				Validators: []validator.Int32{
					int32validator.Any(
						int32validator.OneOf(80, 443),
						int32validator.Between(1024, 65535),
					),
				},
			},
			"verification": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Describes the status of the domain validation. ",
				Attributes: map[string]schema.Attribute{
					"verified": schema.BoolAttribute{
						Computed:    true,
						Description: "True if the domain was successfully verified",
					},
					"verified_at": schema.Int64Attribute{
						Computed:    true,
						Description: "Date time in MS showing since when the domain has been verified",
					},
				},
			},
			"challenge": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "A Domain challenge enables a domain to receive messages. This is a security measure to prevent other domains receiving unwanted messages. The API will send a HTTP Get request to https://{hostname}:{port}/{context_path}/{file_name} and will expect the respones to strictly equal {file_content}. If the content matches, the domain will me marked as verified.",
				Attributes: map[string]schema.Attribute{
					"https": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "The HTTP challenge is currently the only way to verify a domain. This object contains all the relevant information for the client to pass the HTTP challenge.",
						Attributes: map[string]schema.Attribute{
							"file_content": schema.StringAttribute{
								Computed:    true,
								Description: "The file content we expect for the http to succeed.",
							},
							"file_name": schema.StringAttribute{
								Computed:    true,
								Description: "The file name we will request for the http challenge.",
							},
							"context_path": schema.StringAttribute{
								Computed:    true,
								Description: "The context path we will use to proceed with the domain challenge.",
							},
							"created_at": schema.Int64Attribute{
								Computed:    true,
								Description: "A timestamp representing the date time the domain challenge was created.",
							},
							"expires_at": schema.Int64Attribute{
								Computed:    true,
								Description: "A timestamp representing the date time the domain challenge will expire.",
							},
						},
					},
				},
			},
		},
	}
}

func (r *domainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// convertStringToNumber converts a string to an integer and returns the result.
// If the conversion fails, it panics with the error message.
func convertStringToNumber(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic("Error converting string to number: " + err.Error())
	}
	return num
}

func (r *domainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := r.convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting domain to API model",
			"Could not convert domain, unexpected error: "+err.Error(),
		)
		return
	}

	d, err := r.client.CreateDomain(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating domain via API",
			"Could not create domain, unexpected error: "+err.Error(),
		)
		return
	}

	d, err = r.client.GetDomain(d.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain via API",
			"Could not read domain, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(d, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting domain received from API to internal model",
			"Could not convert domain, unexpected error: "+err.Error())

		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *domainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var d, err = r.client.GetDomain(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain via API",
			"Could not read domain, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(d, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting domain received from API to internal model",
			"Could not convert domain, unexpected error: "+err.Error())
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *domainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DomainResourceModel
	diags := req.Plan.Get(ctx, &plan)

	var state DomainResourceModel
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := r.convertToApiModel(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting domain to API model",
			"Could not convert domain, unexpected error: "+err.Error(),
		)
		return
	}
	payload.Hostname = "" // not allowed to update hostname
	payload.Port = 0      // not allowed to update port

	_, err = r.client.UpdateDomain(state.Id.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating domain via API",
			"Could not update domain, unexpected error: "+err.Error(),
		)
		return
	}

	d, err := r.client.GetDomain(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain via API",
			"Could not read domain, unexpected error: "+err.Error(),
		)
		return
	}

	_, err = r.convertFromApiModel(d, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting domain received from API to internal model",
			"Could not convert domain, unexpected error: "+err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *domainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var _, err = r.client.DeleteDomain(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting domain via API",
			"Could not delete domain, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *domainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
