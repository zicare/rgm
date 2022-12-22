package lib

import (
	"reflect"
)

// Pointer-dereference any struct literal
func DeRefPtr(v interface{}) interface{} {

	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		return reflect.ValueOf(v).Elem().Interface()
	}
	return v
}

// Reset pointer to its type zero value
func Reset(v interface{}) {

	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

// Create a new pointer to a new empty value with the same type as what v points to.
// reflect.New(reflect.TypeOf(v).Elem()).Interface()

/*
// Returns true if s is a struct, false otherwise.
func IsStruct(s interface{}) bool {

	ptr := reflect.ValueOf(s)
	if ptr.Kind() != reflect.Ptr {
		return false
	}

	v := ptr.Elem() // dereference the pointer
	if v.Kind() != reflect.Struct {
		return false
	}

	return true
}
*/
