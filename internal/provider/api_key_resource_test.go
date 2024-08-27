package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testCheckScope(resourceType string, name string, access string, target string) resource.TestCheckFunc {
	return resource.TestCheckTypeSetElemNestedAttrs(resourceType, "scopes.*", map[string]string{
		"resource":  name,
		"access":    access,
		"targets.#": "1",
		"targets.0": target,
	})
}

func testCheckNoScope(resourceName string, scopeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		scopesAttr := rs.Primary.Attributes["scopes.#"]
		scopesCount, _ := strconv.Atoi(scopesAttr)

		for i := 0; i < scopesCount; i++ {
			resourceAttr := rs.Primary.Attributes[fmt.Sprintf("scopes.%d.resource", i)]
			if resourceAttr == scopeName {
				return fmt.Errorf("Found %s scope when it was not expected", scopeName)
			}
		}

		return nil
	}
}

func TestAccApiKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "discue_api_key" "test_alias" {
  alias = "my-first-api-key"
  scopes = [{
	  resource = "topics"
	  access = "read"
	  targets = ["*"]
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discue_api_key.test_alias", "alias", "my-first-api-key"),
					resource.TestCheckResourceAttr("discue_api_key.test_alias", "status", "enabled"),
					resource.TestCheckResourceAttrSet("discue_api_key.test_alias", "id"),

					testCheckScope("discue_api_key.test_alias", "topics", "read", "*"),
					testCheckNoScope("discue_api_key.test_alias", "channels"),
					testCheckNoScope("discue_api_key.test_alias", "domains"),
					testCheckNoScope("discue_api_key.test_alias", "events"),
					testCheckNoScope("discue_api_key.test_alias", "listeners"),
					testCheckNoScope("discue_api_key.test_alias", "messages"),
					testCheckNoScope("discue_api_key.test_alias", "queues"),
					testCheckNoScope("discue_api_key.test_alias", "schemas"),
					testCheckNoScope("discue_api_key.test_alias", "stats"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "discue_api_key.test_alias",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "discue_api_key" "test_alias" {
  alias = "my-first-api-key-now"
  scopes = [{
	  resource = "topics"
	  access = "read"
	  targets = ["*"]
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discue_api_key.test_alias", "alias", "my-first-api-key-now"),
					resource.TestCheckResourceAttrSet("discue_api_key.test_alias", "id"),
					resource.TestCheckResourceAttr("discue_api_key.test_alias", "status", "enabled"),

					testCheckScope("discue_api_key.test_alias", "topics", "read", "*"),
					testCheckNoScope("discue_api_key.test_alias", "channels"),
					testCheckNoScope("discue_api_key.test_alias", "domains"),
					testCheckNoScope("discue_api_key.test_alias", "events"),
					testCheckNoScope("discue_api_key.test_alias", "listeners"),
					testCheckNoScope("discue_api_key.test_alias", "messages"),
					testCheckNoScope("discue_api_key.test_alias", "queues"),
					testCheckNoScope("discue_api_key.test_alias", "schemas"),
					testCheckNoScope("discue_api_key.test_alias", "stats"),
				),
			},
			{
				Config: providerConfig + `
resource "discue_api_key" "test_status" {
  alias = "my-first-disabled-api-key"
  status = "disabled"
  scopes = [{
	  resource = "topics"
	  access = "read"
	  targets = ["*"]
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discue_api_key.test_status", "alias", "my-first-disabled-api-key"),
					resource.TestCheckResourceAttr("discue_api_key.test_status", "status", "disabled"),
					resource.TestCheckResourceAttrSet("discue_api_key.test_status", "id"),
				),
			},
			{
				Config: providerConfig + `
resource "discue_api_key" "test_status" {
  alias = "my-first-disabled-api-key"
  status = "enabled"
  scopes = [{
	  resource = "topics"
	  access = "read"
	  targets = ["*"]
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discue_api_key.test_status", "alias", "my-first-disabled-api-key"),
					resource.TestCheckResourceAttr("discue_api_key.test_status", "status", "enabled"),
					resource.TestCheckResourceAttrSet("discue_api_key.test_status", "id"),
				),
			},
			{
				Config: providerConfig + `
resource "discue_api_key" "test_scopes" {
  alias = "my-first-scoped-api-key"
  scopes = [{
	  resource = "domains"
	  access = "read"
	  targets = ["*"]
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discue_api_key.test_scopes", "alias", "my-first-scoped-api-key"),
					resource.TestCheckResourceAttrSet("discue_api_key.test_scopes", "id"),

					testCheckScope("discue_api_key.test_scopes", "domains", "read", "*"),
					testCheckNoScope("discue_api_key.test_scopes", "channels"),
					testCheckNoScope("discue_api_key.test_scopes", "events"),
					testCheckNoScope("discue_api_key.test_scopes", "listeners"),
					testCheckNoScope("discue_api_key.test_scopes", "messages"),
					testCheckNoScope("discue_api_key.test_scopes", "queues"),
					testCheckNoScope("discue_api_key.test_scopes", "schemas"),
					testCheckNoScope("discue_api_key.test_scopes", "stats"),
					testCheckNoScope("discue_api_key.test_scopes", "topics"),
				),
			},
			{
				Config: providerConfig + `
resource "discue_api_key" "test_scopes" {
  alias = "my-first-scoped-api-key"
  scopes = [{
	  resource = "channels"
	  access = "write"
	  targets = ["*"]
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("discue_api_key.test_scopes", "alias", "my-first-scoped-api-key"),
					resource.TestCheckResourceAttrSet("discue_api_key.test_scopes", "id"),

					testCheckScope("discue_api_key.test_scopes", "channels", "write", "*"),
					testCheckNoScope("discue_api_key.test_scopes", "domains"),
					testCheckNoScope("discue_api_key.test_scopes", "events"),
					testCheckNoScope("discue_api_key.test_scopes", "listeners"),
					testCheckNoScope("discue_api_key.test_scopes", "messages"),
					testCheckNoScope("discue_api_key.test_scopes", "queues"),
					testCheckNoScope("discue_api_key.test_scopes", "schemas"),
					testCheckNoScope("discue_api_key.test_scopes", "stats"),
					testCheckNoScope("discue_api_key.test_scopes", "topics"),
				),
			},
		},
	})
}
