package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use environment variables instead
	providerConfig = `
provider "discue" {
  api_key="6vkK9NpVfWjRDbJ5Tob3uNwslDTPcm0Iag91XXHmP335es4dFTzHuLQxkzfYim9v" 
}

`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"discue": providerserver.NewProtocol6WithError(New("test")()),
	}
)
