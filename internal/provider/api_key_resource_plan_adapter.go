// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"log"
	"strings"
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *apiKeyResource) convert(d *client.ApiKeyResponse, plan *apiKeyResourceModel) {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)
	plan.Status = types.StringValue(d.Status)
	plan.Key = types.StringValue(d.Key)

	allScopesAttrValues := make(map[string]attr.Value)
	allScopesAttrTypes := make(map[string]attr.Type)

	for _, name := range ApiResources {
		schemaName := strings.ToLower(name)

		singleScopeAttribute := map[string]attr.Type{
			"access": types.StringType,
			"targets": types.ListType{
				ElemType: types.StringType,
			},
		}

		value := getValueOf[*client.ApiKeyScope](d.Scopes, name)
		var targets = getValueOf[[]string](value, "Targets")

		targetsList, _ := StringArrayToList(targets)
		access := getValueOf[string](value, "Access")

		singleScopeValue := map[string]attr.Value{
			"access":  types.StringValue(access),
			"targets": targetsList,
		}

		if len(targets) > 0 && len(access) > 0 {
			log.Println(fmt.Sprintf("Setting scope for %s %#+v %#+v", schemaName, access, targets))

			scopeObject, _ := basetypes.NewObjectValue(singleScopeAttribute, singleScopeValue)

			allScopesAttrTypes[schemaName] = basetypes.ObjectType{
				AttrTypes: singleScopeAttribute,
			}

			allScopesAttrValues[schemaName] = scopeObject
		}
	}

	allScopes, _ := basetypes.NewObjectValue(allScopesAttrTypes, allScopesAttrValues)
	plan.Scopes = allScopes
}
