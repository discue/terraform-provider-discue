package provider

import (
	"fmt"
	"reflect"
)

func getValueOf[T any](itf interface{}, key string) T {
	if !checkFieldExists(itf, key) {
		var result T
		return result
	}

	field := getFieldByName(itf, key, true)

	if field.IsValid() && field.CanInterface() {
		return field.Interface().(T)
	} else {
		fmt.Printf("Field %s not found or inaccessible\n", key)
	}
	var result T
	return result
}

func setValueOf(itf interface{}, key string, value any) {
	field := getFieldByName(itf, key, false)

	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		fmt.Printf("Field %s not found or inaccessible\n", key)
	}
}

func getFieldByName(itf interface{}, key string, ignoreInvalid bool) reflect.Value {
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
