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

func (r *domainResource) convertDomainToApiModel(_ctx context.Context, plan *DomainResourceModel) (client.DomainRequest, error) {
	return client.DomainRequest{
		Alias:    plan.Alias.ValueString(),
		Hostname: plan.Hostname.ValueString(),
		Port:     convertStringToNumber(plan.Port.String()),
	}, nil
}

func (r *domainResource) convertDomainToInternalModel(d *client.DomainResponse, plan *DomainResourceModel) (*DomainResourceModel, error) {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)
	plan.Port = types.Int32Value(d.Port)
	plan.Hostname = types.StringValue(d.Hostname)

	var err error
	plan.Verification, err = convertDomainVerification(d)
	if err != nil {
		return plan, err
	}

	plan.Challenge, err = convertChallenges(d)
	return plan, err
}

func convertDomainVerification(d *client.DomainResponse) (basetypes.ObjectValue, error) {
	verificationAttrTypes := map[string]attr.Type{
		"verified":    types.BoolType,
		"verified_at": types.Int64Type,
	}

	verificationAttrValues := map[string]attr.Value{
		"verified":    types.BoolValue(d.Verification.Verified),
		"verified_at": types.Int64Value(d.Verification.VerifiedAt),
	}

	verificationObjectValue, diags := basetypes.NewObjectValue(verificationAttrTypes, verificationAttrValues)
	if diags.HasError() {
		return basetypes.ObjectValue{}, DiagsToStructuredError("Unable to create object types and attrs for domain challenge", diags)
	}
	return verificationObjectValue, nil
}

func convertChallenges(d *client.DomainResponse) (basetypes.ObjectValue, error) {
	domainChallengeAttrTypes := map[string]attr.Type{
		"https": basetypes.ObjectType{
			AttrTypes: map[string]attr.Type{
				"file_content": types.StringType,
				"file_name":    types.StringType,
				"context_path": types.StringType,
				"created_at":   types.Int64Type,
				"expires_at":   types.Int64Type,
			},
		},
	}

	httpChallengeAttrValues := map[string]attr.Value{
		"file_content": types.StringValue(d.Challenge.Https.FileContent),
		"file_name":    types.StringValue(d.Challenge.Https.FileName),
		"context_path": types.StringValue(d.Challenge.Https.ContextPath),
		"created_at":   types.Int64Value(d.Challenge.Https.CreatedAt),
		"expires_at":   types.Int64Value(d.Challenge.Https.ExpiresAt),
	}

	httpChallengeObjVal, diags := basetypes.NewObjectValue(domainChallengeAttrTypes["https"].(basetypes.ObjectType).AttrTypes, httpChallengeAttrValues)
	if diags.HasError() {
		return basetypes.ObjectValue{}, DiagsToStructuredError("Unable to create object value for https challenge", diags)
	}

	domainChallengeAttrValues := map[string]attr.Value{
		"https": httpChallengeObjVal,
	}

	domainChallengeObjVal, diags := basetypes.NewObjectValue(domainChallengeAttrTypes, domainChallengeAttrValues)
	if diags.HasError() {
		return basetypes.ObjectValue{}, DiagsToStructuredError(fmt.Sprintf("Unable to create doamin challenge object value"), diags)
	}

	return domainChallengeObjVal, nil
}
