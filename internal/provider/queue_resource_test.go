package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueueResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_queue.test_queue", "alias", "my-first-queue"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_queue.test_queue", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "discue_queue.test_queue",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue-with-new-alias"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_queue.test_queue", "alias", "my-first-queue-with-new-alias"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_queue.test_queue", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
