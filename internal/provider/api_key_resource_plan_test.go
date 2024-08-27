package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = debugPlan{}

type debugPlan struct{}

func (e debugPlan) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	rd, err := json.Marshal(req.Plan)
	if err != nil {
		fmt.Println("error marshalling machine-readable plan output:", err)
	}
	fmt.Printf("req.Plan - %s\n", string(rd))
}

func DebugPlan() plancheck.PlanCheck {
	return debugPlan{}
}

func Test_DebugPlan(t *testing.T) {
	t.Parallel()

	r.Test(t, r.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []r.TestStep{
			{
				Config: providerConfig + `resource "discue_api_key" "one" {
					alias = "123123aa"
					scopes = [
						{
						resource = "queues"
						access   = "read"
						targets  = ["*"]
						}
					]
				}`,
				PlanOnly: true,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
			},
			{
				Config: providerConfig + `resource "discue_api_key" "one" {
					alias = "12312311aa"
				}`,
				PlanOnly: true,
				ConfigPlanChecks: r.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
			},
		},
	})
}
