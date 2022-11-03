package lib

import "reflect"

// Pointer-dereference any struct literal
func DeRefPtr(v interface{}) interface{} {

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		return reflect.ValueOf(v).Elem().Interface()
	}
	return v
}
