package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringListTypeToPlainStringArray(ctx context.Context, tfVal types.List) ([]string, diag.Diagnostics) {
	return ListTypeToPlainArray[string](ctx, tfVal)
}

func ListTypeToPlainArray[T any](ctx context.Context, tfVal types.List) ([]T, diag.Diagnostics) {
	if !HasValue(tfVal) {
		return nil, nil
	}
	result := make([]T, len(tfVal.Elements()))
	diags := tfVal.ElementsAs(ctx, &result, false)
	return result, diags
}

func PlainStringMapToMapType(stringMap map[string]string) (types.Map, error) {
	elements := map[string]attr.Value{}
	for k, v := range stringMap {
		elements[k] = types.StringValue(v)
	}
	mapValue, diags := types.MapValue(types.StringType, elements)
	if diags != nil && diags.HasError() {
		return mapValue, fmt.Errorf("failed to convert to MapType %v", diags.Errors()[0].Detail())
	}
	return mapValue, nil
}

func PlainStringArrayToListType(stringList []string) (types.List, diag.Diagnostics) {
	elements := []attr.Value{}
	for _, e := range stringList {
		elements = append(elements, types.StringValue(e))
	}
	return types.ListValue(types.StringType, elements)
}
