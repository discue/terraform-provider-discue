// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *apiKeyResource) convertFromApiModel(d *client.ApiKeyResponse, plan *apiKeyResourceModel) (*apiKeyResourceModel, error) {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)
	plan.Status = types.StringValue(d.Status)
	plan.Key = types.StringValue(d.Key)

	scopes, err := convertScopesFromApiModel(*d.Scopes)
	if err != nil {
		var r *apiKeyResourceModel
		return r, err
	}
	plan.Scopes = scopes

	return plan, nil
}

func convertScopesFromApiModel(scopes client.ApiKeyScopes) (basetypes.ListValue, error) {
	elements := []attr.Value{}

	for _, name := range ApiResources {
		keyName := uppercaseFirstCharacter(name)
		value, _ := getValueOf[*client.ApiKeyScope](scopes, keyName)
		if value == nil {
			continue
		}

		targets := value.Targets
		access := value.Access
		targetsList, diags := PlainStringArrayToListType(targets)

		if diags.HasError() {
			var v basetypes.ListValue
			return v, DiagsToStructuredError(fmt.Sprintf("Unable to create object value for single scope %#+v", scopes), diags)
		}

		objVal, diags := types.ObjectValue(map[string]attr.Type{
			"resource": types.StringType,
			"access":   types.StringType,
			"targets":  types.ListType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"resource": types.StringValue(name),
			"access":   types.StringValue(access),
			"targets":  targetsList,
		})

		if diags.HasError() {
			var v basetypes.ListValue
			return v, DiagsToStructuredError(fmt.Sprintf("Unable to create object value for single scope %#+v", scopes), diags)
		}

		elements = append(elements, objVal)
	}

	listValue, diags := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"resource": types.StringType,
			"access":   types.StringType,
			"targets":  types.ListType{ElemType: types.StringType},
		},
	}, elements)

	if diags.HasError() {
		var v basetypes.ListValue
		return v, DiagsToStructuredError(fmt.Sprintf("Unable to create list value for scopes %#+v", scopes), diags)
	}

	return listValue, nil
}

func (r *apiKeyResource) convertToApiModel(ctx context.Context, plan *apiKeyResourceModel) (client.ApiKeyRequest, error) {
	req := client.ApiKeyRequest{
		Alias:  plan.Alias.ValueString(),
		Status: plan.Status.ValueString(),
	}

	converted, err := convertScopesToApiModel(ctx, plan)
	if err != nil {
		var r client.ApiKeyRequest
		return r, err
	}
	req.Scopes = &converted

	return req, nil
}

func convertScopesToApiModel(ctx context.Context, plan *apiKeyResourceModel) (client.ApiKeyScopes, error) {
	elements, diags := ListTypeToPlainArray[apiKeyScopeModel](ctx, plan.Scopes)
	if diags.HasError() {
		var r client.ApiKeyScopes
		return r, DiagsToStructuredError("Unable to convert to state/plan to struct", diags)
	}

	scopes := client.ApiKeyScopes{}
	for _, scope := range elements {

		targets, diags := StringListTypeToPlainStringArray(ctx, scope.Targets)
		if diags.HasError() {
			var r client.ApiKeyScopes
			return r, DiagsToStructuredError("Unable to convert scope targets to string array", diags)
		}

		err := setValueOf(&scopes, scope.Resource.ValueString(), &client.ApiKeyScope{
			Access:  scope.Access.ValueString(),
			Targets: targets,
		})
		if err != nil {
			var r client.ApiKeyScopes
			return r, DiagsToStructuredError("Unable set scope", diags)
		}
	}

	return scopes, nil
}
