package provider

import (
	"fmt"
	"reflect"
)

func getValueOf[T any](itf interface{}, key string) (T, error) {
	if !checkFieldExists(itf, key) {
		var result T
		return result, fmt.Errorf("field %s not found", key)
	}

	field := getFieldByName(itf, key)

	if field.IsValid() && field.CanInterface() {
		return field.Interface().(T), nil
	}

	var result T
	return result, fmt.Errorf("field %s inaccessible", key)
}

func setValueOf(itf interface{}, key string, value any) error {
	field := getFieldByName(itf, key)

	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(value))
		return nil
	} else {
		return fmt.Errorf("field %s not found or inaccessible", key)
	}
}

func getFieldByName(itf interface{}, key string) reflect.Value {
	v := reflect.ValueOf(itf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName(key)
	if !field.IsValid() {
		field = v.FieldByName(uppercaseFirstCharacter(key))
	}
	return field
}

func checkFieldExists(itf interface{}, fieldName string) bool {
	val := reflect.ValueOf(itf)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if !val.IsValid() {
		return false
	}

	return val.FieldByName(fieldName).IsValid()
}
