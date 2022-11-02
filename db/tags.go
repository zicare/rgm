package db

import (
	"reflect"
	"strings"

	"github.com/zicare/rgm/msg"
)

// TaggedFields returns a slice of db fields sharing a tag name (tn)
// such fields must be tagged as tn:"tv"
// the tagName value is passed in the tn parameter
// tagValues options are passed in the tv parameter
// ie.:
// type User struct {
//	 UserID    int64     `db:"user_id"   auth:"id"    json:"user_id"   pk:"1"`
//	 RoleID    *int64    `db:"role_id"   auth:"role"  json:"role_id"`
//	 Email     string    `db:"email"     auth:"usr"   json:"email"`
// }
// v ar m User
// f, _ := TaggedFields(m, "auth", []string{"id","role","usr"})
// f -> []string{"user_id","role_id","email"}
func TaggedFields(m Table, tagName string, tagValues []string) ([]string, error) {

	var (
		f = make([]string, len(tagValues))
		t = reflect.Indirect(reflect.ValueOf(m))
	)

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Type().Field(i).Tag.Lookup(tagName); ok {
			if col, ok := t.Type().Field(i).Tag.Lookup("db"); ok {
				for i, v := range tagValues {
					if v == tag {
						f[i] = col
					}
				}
			}
		}
	}

	for _, col := range f {
		if col == "" {
			//%s tags not properly set
			return f, msg.Get("2").SetArgs(strings.Title(tagName)).M2E()
		}
	}

	return f, nil
}

func Pk(m Table) (f []string) {

	var t = reflect.Indirect(reflect.ValueOf(m))

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Type().Field(i).Tag.Lookup("pk"); ok {
			if col, ok := t.Type().Field(i).Tag.Lookup("db"); ok {
				if tag == "1" {
					f = append(f, col)
				}
			}
		}
	}

	return f
}
