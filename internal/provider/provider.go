// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &discueProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &discueProvider{
			version: version,
		}
	}
}

type discueProvider struct {
	version string
}

func (p *discueProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "discue"
	resp.Version = p.version
}

func (p *discueProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Sensitive:   true,
				Optional:    true,
				Description: "The API key used to access discue.io resources.",
			},
			"api_endpoint": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "The API endpoint used to access discue.io resources.",
			},
		},
	}
}

type discueProviderModel struct {
	ApiKey      types.String `tfsdk:"api_key"`
	ApiEndpoint types.String `tfsdk:"api_endpoint"`
}

func (p *discueProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring discue client")

	var config discueProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("DISCUE_API_KEY")
	if config.ApiKey.ValueString() != "" {
		apiKey = config.ApiKey.ValueString()
	}

	apiEndpoint := os.Getenv("DISCUE_API_ENDPOINT")
	if apiEndpoint == "" {
		apiEndpoint = "http://localhost:3000"
	}
	if config.ApiEndpoint.ValueString() != "" {
		apiEndpoint = config.ApiEndpoint.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing discue API key",
			"The provider cannot create the discue API client as there is a missing or empty value for the discue API key.",
		)
	}

	if apiEndpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_endpoint"),
			"Missing discue API endpoint",
			"The provider cannot create the discue API client as there is a missing or empty value for the discue API endpoint.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new client using the configuration values
	client, err := client.NewClient(apiEndpoint, &apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create discue API Client",
			"An unexpected error occurred when creating the API client.",
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured discue client", map[string]any{"success": true})
}

func (p *discueProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *discueProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQueueResource,
	}
}

func (p *discueProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}
