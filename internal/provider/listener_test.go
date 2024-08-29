package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccListenerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// test error when liveness_url is a invalid url
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  notify_url = "https://discue.io/live"
  liveness_url = "https://user@password:discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when liveness_url is a invalid url
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "https://discue.io/live"
  notify_url = "https://user@password:discue.io/notify"
}
`,

				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when liveness_url contains basic auth
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  notify_url = "https://discue.io/live"
  liveness_url = "https://user:password@discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when notify_url contains basic auth
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "https://discue.io/live"
  notify_url = "https://user:password@discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when liveness_url does not specificy https protocol
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "ftp://discue.io/live"
  notify_url = "https://discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when notify_url does not specificy https protocol
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "https://discue.io/live"
  notify_url = "fpt://discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when notify_url is a relative url
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "https://discue.io/live"
  notify_url = "/notify"
}
`,
				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when liveness_url is a relative url
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "live"
  notify_url = "https://discue.io//notify"
}
`,
				ExpectError: regexp.MustCompile("must be a valid URL with http or https protocol"),
			},
			{
				// test error when queue id is empty
				Config: providerConfig + `
resource "discue_listener" "test_listener_schema" {

  alias = "my-listener"
  liveness_url = "https://discue.io/live"
  notify_url = "https://discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("The argument \"queue_id\" is required, but no definition was found"),
			},
			{
				// test error when alias is empty
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  liveness_url = "https://discue.io/live"
  notify_url = "https://discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("The argument \"alias\" is required, but no definition was found"),
			},
			{
				// test error when liveness_url is empty
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  notify_url = "https://discue.io/notify"
}
`,
				ExpectError: regexp.MustCompile("The argument \"liveness_url\" is required, but no definition was found"),
			},
			{
				// test error when notify_url is empty
				Config: providerConfig + `
resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener_schema" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "https://discue.io/live"
}
`,
				ExpectError: regexp.MustCompile("The argument \"notify_url\" is required, but no definition was found"),
			},
			{
				// test create new listener
				Config: providerConfig + `
resource "discue_domain" "test_domain" {
  alias = "my-first-domain"
  hostname = "discue.io"
  port = 443
}

resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/live"
  notify_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/notify"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_listener.test_listener", "alias", "my-listener"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "liveness_url", "https://discue.io:443/live"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "notify_url", "https://discue.io:443/notify"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_listener.test_listener", "id"),
				),
			},
			// test import
			{
				ResourceName:      "discue_listener.test_listener",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					resources := s.RootModule().Resources
					return fmt.Sprintf("%s,%s", resources["discue_queue.test_queue"].Primary.ID, resources["discue_listener.test_listener"].Primary.ID), nil
				},
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			{
				// test create new listener
				Config: providerConfig + `
resource "discue_domain" "test_domain" {
  alias = "my-first-domain"
  hostname = "discue.io"
  port = 443
}

resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener" {
  queue_id = discue_queue.test_queue.id

  alias = "my-listener"
  liveness_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/live"
  notify_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/notify"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_listener.test_listener", "alias", "my-listener"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "liveness_url", "https://discue.io:443/live"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "notify_url", "https://discue.io:443/notify"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_listener.test_listener", "id"),
				),
			},
			{
				// test update alias
				Config: providerConfig + `
resource "discue_domain" "test_domain" {
  alias = "my-first-domain"
  hostname = "discue.io"
  port = 443
}

resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener" {
  queue_id = discue_queue.test_queue.id

  alias = "my-new-listener-alias"
  liveness_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/live"
  notify_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/notify"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_listener.test_listener", "alias", "my-new-listener-alias"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "liveness_url", "https://discue.io:443/live"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "notify_url", "https://discue.io:443/notify"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_listener.test_listener", "id"),
				),
			},
			{
				// test update liveness_url
				Config: providerConfig + `
resource "discue_domain" "test_domain" {
  alias = "my-first-domain"
  hostname = "discue.io"
  port = 443
}

resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}
  
resource "discue_listener" "test_listener" {
  queue_id = discue_queue.test_queue.id

  alias = "my-new-listener-alias"
  liveness_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/live11"
  notify_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/notify"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_listener.test_listener", "alias", "my-new-listener-alias"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "liveness_url", "https://discue.io:443/live11"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "notify_url", "https://discue.io:443/notify"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_listener.test_listener", "id"),
				),
			},
			{
				// test update notify_url
				Config: providerConfig + `
resource "discue_domain" "test_domain" {
  alias = "my-first-domain"
  hostname = "discue.io"
  port = 443
}

resource "discue_queue" "test_queue" {
  alias = "my-first-queue"
}

resource "discue_listener" "test_listener" {
  queue_id = discue_queue.test_queue.id

  alias = "my-new-listener-alias"
  liveness_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/live11"
  notify_url = "https://${discue_domain.test_domain.hostname}:${discue_domain.test_domain.port}/notify1123"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					resource.TestCheckResourceAttr("discue_listener.test_listener", "alias", "my-new-listener-alias"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "liveness_url", "https://discue.io:443/live11"),
					resource.TestCheckResourceAttr("discue_listener.test_listener", "notify_url", "https://discue.io:443/notify1123"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("discue_listener.test_listener", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
