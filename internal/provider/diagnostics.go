package provider

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func DiagsToStructuredError(message string, diags diag.Diagnostics) error {
	messages := []string{
		message,
		"\n\n",
		"Reasons:\n",
	}

	for i, error := range diags.Errors() {
		messages = append(messages, fmt.Sprintf("%d: %s \n", i, error.Detail()))
	}

	return errors.New(strings.Join(messages, ""))
}
