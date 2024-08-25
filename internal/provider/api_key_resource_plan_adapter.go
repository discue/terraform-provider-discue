// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *apiKeyResource) convert(d *client.ApiKeyResponse, plan *apiKeyResourceModel) {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)
	plan.Status = types.StringValue(d.Status)
	plan.Key = types.StringValue(d.Key)

	var scopes []apiKeyScopeModel

	for _, name := range ApiResources {
		value := getValueOf[*client.ApiKeyScope](d.Scopes, name)
		if value == nil {
			continue
		}

		targets := value.Targets
		access := value.Access
		targetsList, _ := StringArrayToList(targets)

		if len(targets) > 0 && len(access) > 0 {
			scopes = append(scopes, apiKeyScopeModel{
				Resource: name,
				Access:   types.StringValue(access),
				Targets:  targetsList,
			})
		}
	}

	plan.Scope = scopes
}
