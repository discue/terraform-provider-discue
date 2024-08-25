package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func expectScope(resourceType string, name string, access string, target string) resource.TestCheckFunc {
	return resource.TestCheckTypeSetElemNestedAttrs(resourceType, "scope.*", map[string]string{
		"resource":  name,
		"access":    access,
		"targets.#": "1",
		"targets.0": target,
	})
}

func TestAccApiKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "discue_api_key" "test_api_key" {
  alias = "my-first-api-key"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_api_key.test_api_key", "alias", "my-first-api-key"),
					resource.TestCheckResourceAttr("discue_api_key.test_api_key", "status", "enabled"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "id"),
					resource.TestCheckTypeSetElemNestedAttrs("discue_api_key.test_api_key", "scope.*", map[string]string{
						"resource":  "Domains",
						"access":    "read",
						"targets.#": "1",
						"targets.0": "*",
					}), // resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "scopes.domains.access"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "verification.verified_at"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.%"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.file_content"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.file_name"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "discue_api_key.test_api_key",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "discue_api_key" "test_api_key" {
  alias = "my-first-api-key-now"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_api_key.test_api_key", "alias", "my-first-api-key-now"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "id"),
				),
			},
			{
				Config: providerConfig + `
resource "discue_api_key" "test_api_key_disabled" {
  alias = "my-first-disabled-api-key"
  status = "disabled"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_api_key.test_api_key_disabled", "alias", "my-first-disabled-api-key"),
					resource.TestCheckResourceAttr("discue_api_key.test_api_key_disabled", "status", "disabled"),
					resource.TestCheckResourceAttrSet("discue_api_key.test_api_key_disabled", "id"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key_disabled", "scopes.domains.access"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "verification.verified_at"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.%"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.file_content"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.file_name"),
				),
			},
			{
				Config: providerConfig + `
resource "discue_api_key" "test_only_domains_scope" {
  alias = "my-first-scoped-api-key"
  scope {
	  resource = "Domains"
	  access = "read"
	  targets = ["*"]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_api_key.test_only_domains_scope", "alias", "my-first-scoped-api-key"),
					resource.TestCheckTypeSetElemNestedAttrs("discue_api_key.test_only_domains_scope", "scope.*", map[string]string{
						"resource":  "Domains",
						"access":    "read",
						"targets.#": "1",
						"targets.0": "*",
					}),
					resource.TestCheckResourceAttrSet("discue_api_key.test_only_domains_scope", "id"),
					resource.TestCheckNoResourceAttr("discue_api_key.test_only_domains_scope", "scopes.queues.access"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "verification.verified_at"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.%"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.file_content"),
					// resource.TestCheckResourceAttrSet("discue_api_key.test_api_key", "challenge.https.file_name"),
				),
			},
		},
	})
}
