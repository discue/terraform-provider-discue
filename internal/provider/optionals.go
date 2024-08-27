package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HasValue checks if the given terraform value is set.
func HasValue(val attr.Value) bool {
	return !val.IsUnknown() && !val.IsNull()
}

func BoolWithTrueDefault(tfVal types.Bool) bool {
	if HasValue(tfVal) {
		return tfVal.ValueBool()
	}
	return true
}

func BoolWithFalseDefault(tfVal types.Bool) bool {
	if HasValue(tfVal) {
		return tfVal.ValueBool()
	}
	return false
}

func OptionalInt64(tfVal types.Int64) *int64 {
	if HasValue(tfVal) {
		return tfVal.ValueInt64Pointer()
	}
	return nil
}

func OptionalString(tfVal types.String) *string {
	if HasValue(tfVal) {
		return tfVal.ValueStringPointer()
	}
	return nil
}
func OptionalMap(ctx context.Context, tfVal types.Map) (map[string]string, error) {
	if !HasValue(tfVal) {
		return nil, nil
	}
	result := make(map[string]string, len(tfVal.Elements()))
	d := tfVal.ElementsAs(ctx, &result, false)
	if d.HasError() {
		return nil, fmt.Errorf("error converting to map object %v", d.Errors()[0].Detail())
	}

	return result, nil
}

func OptionalList(tfVal types.List) []string {
	if !HasValue(tfVal) {
		return nil
	}
	result := make([]string, 0)
	for _, e := range tfVal.Elements() {
		result = append(result, e.(types.String).ValueString())
	}
	return result
}
