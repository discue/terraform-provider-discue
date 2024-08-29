// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUrlValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
	}
	tests := map[string]testCase{
		"unknown URL": {
			val: types.StringUnknown(),
		},
		"null URL": {
			val: types.StringNull(),
		},
		"valid URL": {
			val: types.StringValue("http://www.discue.io/live"),
		},
		"invalid URL wrong procotol": {
			val:         types.StringValue("ftp://www.discue.io"),
			expectError: true,
		},
		"invalid URL basic auth": {
			val:         types.StringValue("ftp://a:qb@www.discue.io"),
			expectError: true,
		},
		"invalid URL format": {
			val:         types.StringValue("ftp://a@qb:www.discue.io"),
			expectError: true,
		},
		"invalid relative URL": {
			val:         types.StringValue("/abc"),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			ValidUrl("").ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
