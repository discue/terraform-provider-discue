// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = idValidator{}

// idValidator validates that a string Attribute's value matches the specified regular expression.
type idValidator struct {
	regexp  *regexp.Regexp
	message string
}

// Description describes the validation in plain text formatting.
func (validator idValidator) Description(_ context.Context) string {
	if validator.message != "" {
		return validator.message
	}
	return fmt.Sprintf("value must match regular expression '%s'", validator.regexp)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator idValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v idValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !v.regexp.MatchString(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// RegexMatches returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a string.
//   - Matches the regular expression for resource ids
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
// Optionally an error message can be provided to return something friendlier
// than "value must match regular expression 'regexp'".
func ValidResourceId(message string) validator.String {
	return idValidator{
		regexp:  regexp.MustCompile(`^[useandom26T198340PX75pxJACKVERYMINDBUSHWOLFGQZbfghjklqvwyzrict-]{21}$`),
		message: message,
	}
}
