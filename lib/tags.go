package lib

import (
	"reflect"

	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
)

// TaggedFields returns a slice of db fields sharing a tag name (tn)
// such fields must be tagged as tn:"tv"
// the tn value is passed in the tn parameter
// tv options are passed in the tv parameter
// ie.:
// type User struct {
//	 UserID    int64     `db:"user_id"   auth:"id"    json:"user_id"   primary:"1"`
//	 RoleID    *int64    `db:"role_id"   auth:"role"  json:"role_id"`
//	 Email     string    `db:"email"     auth:"usr"   json:"email"`
// }
// v ar m User
// f, _ := TaggedFields(m, "auth", []string{"id","role","usr"})
// f -> []string{"user_id","role_id","email"}
func TaggedFields(m db.Table, tn string, tv []string) ([]string, error) {

	var (
		f = make([]string, len(tv))
		t = reflect.Indirect(reflect.ValueOf(m))
	)

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Type().Field(i).Tag.Lookup(tn); ok {
			if col, ok := t.Type().Field(i).Tag.Lookup("db"); ok {
				for i, v := range tv {
					if v == tag {
						f[i] = col
					}
				}
			}
		}
	}

	for _, col := range f {
		if col == "" {
			//Auth tags are not properly set
			return f, msg.Get("30").M2E()
		}
	}

	return f, nil
}
