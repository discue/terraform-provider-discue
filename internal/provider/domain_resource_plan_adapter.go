// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"terraform-provider-discue/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *domainResource) convert(d *client.DomainResponse, plan *DomainResourceModel) *DomainResourceModel {
	plan.Id = types.StringValue(d.Id)
	plan.Alias = types.StringValue(d.Alias)
	plan.Port = types.Int32Value(d.Port)
	plan.Hostname = types.StringValue(d.Hostname)

	// Convert DomainVerification to basetypes.ObjectValue
	verificationAttrTypes := map[string]attr.Type{
		"verified":    types.BoolType,
		"verified_at": types.Int64Type,
	}

	verificationAttrValues := map[string]attr.Value{
		"verified":    types.BoolValue(d.Verification.Verified),
		"verified_at": types.Int64Value(d.Verification.VerifiedAt),
	}

	verificationObjectValue, _ := basetypes.NewObjectValue(verificationAttrTypes, verificationAttrValues)
	plan.Verification = verificationObjectValue

	plan.Challenge = func() basetypes.ObjectValue {
		// Define the attribute types for DomainChallenge
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

		// Define the attribute values for HttpDomainChallenge
		httpChallengeAttrValues := map[string]attr.Value{
			"file_content": types.StringValue(d.Challenge.Https.FileContent),
			"file_name":    types.StringValue(d.Challenge.Https.FileName),
			"context_path": types.StringValue(d.Challenge.Https.ContextPath),
			"created_at":   types.Int64Value(d.Challenge.Https.CreatedAt),
			"expires_at":   types.Int64Value(d.Challenge.Https.ExpiresAt),
		}

		// Create a new ObjectValue for HttpDomainChallenge
		httpChallengeObjVal, err := basetypes.NewObjectValue(domainChallengeAttrTypes["https"].(basetypes.ObjectType).AttrTypes, httpChallengeAttrValues)
		if err != nil {
			// Handle the error appropriately
			return basetypes.NewObjectNull(nil)
		}

		// Define the attribute values for DomainChallenge
		domainChallengeAttrValues := map[string]attr.Value{
			"https": httpChallengeObjVal,
		}

		// Create a new ObjectValue for DomainChallenge
		domainChallengeObjVal, err := basetypes.NewObjectValue(domainChallengeAttrTypes, domainChallengeAttrValues)
		if err != nil {
			// Handle the error appropriately
			return basetypes.NewObjectNull(nil)
		}

		return domainChallengeObjVal
	}()

	return plan
}
