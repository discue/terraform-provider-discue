// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = urlValidator{}

type urlValidator struct {
	message string
}

// Description describes the validation in plain text formatting.
func (validator urlValidator) Description(_ context.Context) string {
	if validator.message != "" {
		return validator.message
	}
	return "must be a valid URL with http or https protocol and without authentication"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator urlValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v urlValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	parsedUrl, err := url.Parse(value)
	if err != nil || parsedUrl == nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
		return
	}

	scheme := parsedUrl.Scheme
	if scheme != "http" && scheme != "https" {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}

	hostname := parsedUrl.Host
	if hostname == "" {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}

	userOrPw := parsedUrl.User.String()
	if userOrPw != "" {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// Returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a valid URL according to https://pkg.go.dev/net/url#Parse
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
// Optionally an error message can be provided
func ValidUrl(message string) validator.String {
	return urlValidator{
		message: message,
	}
}
