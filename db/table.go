package db

import (
	"github.com/go-playground/validator/v10"
	"github.com/zicare/rgm/msg"
)

// Table exported
type Table interface {

	// Must return the table name
	Name() string

	// Must attach foreign table data if available
	Dig(f ...string) []Dig

	// Must set conditions to filter out content
	// not intended for the user uid making the request
	// or not a child of optional parent t records
	Scope(uid string, t ...Table) map[string]string

	// A chance to complete extra validations, set
	// default values, etc.
	BeforeInsert(uid string, t ...Table) error

	// Must set conditions to filter out content
	// not intended for the user uid making the request
	// or not a child of optional parent t records.
	// Second return parameter indicates whether or not
	// is okay to proceed with delete action.
	BeforeDelete(uid string, t ...Table) (map[string]string, bool)

	// Must return a validation errors list
	ValidationErrors(err error) msg.MessageList
}

type BaseTable struct{}

// Dig exported
func (BaseTable) Dig(f ...string) []Dig {

	return []Dig{}
}

// Scope exported
func (BaseTable) Scope(uid string, t ...Table) map[string]string {

	return make(map[string]string)
}

// Scope exported
func (BaseTable) BeforeInsert(uid string, t ...Table) error {

	return nil
}

// Scope exported
func (BaseTable) BeforeDelete(uid string, t ...Table) (map[string]string, bool) {

	return make(map[string]string), true
}

// Return a validation errors list
func (BaseTable) ValidationErrors(err error) (ml msg.MessageList) {

	switch err.(type) {
	/*
		case *time.ParseError:
			//Time %s has a wrong format, required format is %s
			e := err.(*time.ParseError)
			m := msg.Get("22").SetArgs(lib.TrimQuotes(e.Value), "2006-01-02T15:04:05-07:00")
			eml = append(eml, m)
		case *json.UnmarshalTypeError:
			//Value is a %s, required type is %s
			e := err.(*json.UnmarshalTypeError)
			m := msg.Get("23").SetArgs(e.Value, e.Type.String()).SetField(e.Field)
			eml = append(eml, m)
	*/
	case validator.ValidationErrors:
		for _, v := range err.(validator.ValidationErrors) {
			//typ := v.Type().String()
			m := msg.Get("24").SetArgs(v.Value(), v.Tag(), v.Param()).SetField(v.Field())
			ml = append(ml, m)
		}
	}
	return ml
}
