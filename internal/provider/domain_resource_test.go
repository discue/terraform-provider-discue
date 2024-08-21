package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "discue_domain" "test_domain" {
  alias = "my-first-domain"
  hostname = "discue.io"
  port = 443
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_domain.test_domain", "alias", "my-first-domain"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_domain.test_domain", "id"),
					resource.TestCheckResourceAttrSet("discue_domain.test_domain", "verification.verified"),
					resource.TestCheckResourceAttrSet("discue_domain.test_domain", "verification.verified_at"),
					resource.TestCheckResourceAttrSet("discue_domain.test_domain", "challenge.https.%"),
					resource.TestCheckResourceAttrSet("discue_domain.test_domain", "challenge.https.file_content"),
					resource.TestCheckResourceAttrSet("discue_domain.test_domain", "challenge.https.file_name"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "discue_domain.test_domain",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "discue_domain" "test_domain" {
  alias = "my-first-domain-with-new-alias"
  hostname = "discue.io"
  port = 443
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_domain.test_domain", "alias", "my-first-domain-with-new-alias"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_domain.test_domain", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
