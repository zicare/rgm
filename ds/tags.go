package ds

import (
	"reflect"
	"strings"

	"github.com/zicare/rgm/msg"
)

func Meta(m IDataSource) (k, f, w []string, d map[string]int, e *TagError) {

	d = make(map[string]int)

	//r := reflect.ValueOf(m).Elem()
	r := reflect.Indirect(reflect.ValueOf(m))
	rt := r.Type()
	for i := 0; i < r.NumField(); i++ {
		tag := rt.Field(i).Tag
		if db, ok := tag.Lookup("db"); ok {
			if db == "-" {
				// diggable only if it ALSO has fk:"..."
				if fk, ok := tag.Lookup("fk"); ok && fk != "" {
					if j, ok := tag.Lookup("json"); ok {
						if j, _, _ := strings.Cut(j, ","); j != "" && j != "-" {
							d[j] = i
						}
					}
				}
				continue
			}
			// normal field
			f = append(f, db)
			if pk, ok := tag.Lookup("pk"); ok && pk == "1" {
				k = append(k, db)
			}
			if _, ok := tag.Lookup("view"); !ok {
				w = append(w, db)
			}
		}
	}

	if len(k) == 0 {
		return k, f, w, d, &TagError{Message: msg.Message{Msg: "no pk tags in model"}}
	}
	if e = ValidateTags(m, k, f); e != nil {
		return k, f, w, d, e
	}
	return k, f, w, d, nil
}

func Values(d IDataSource) map[string]interface{} {

	v := make(map[string]interface{})
	r := reflect.Indirect(reflect.ValueOf(d))
	rt := r.Type()
	for i := 0; i < r.NumField(); i++ {
		tag := rt.Field(i).Tag
		if db, ok := tag.Lookup("db"); ok && db != "-" {
			if _, ok := tag.Lookup("view"); !ok {
				v[db] = r.Field(i).Interface()
			}
		}
	}

	return v
}

// Works on structs with two tag sets, let's call them target and pivot.
//
// Returns a slice with the target tag's values, provided the pivot
// tag name and values are matched.
//
// model is the struct to work with.
//
// tagName is the target tag name.
//
// pivotTagName is the pivot tag name to be matched.
//
// pivotTagFields are the pivot tag values to be matched.
//
// Example:
//
//	type User struct {
//		 UserID    int64     `db:"user_id"   auth:"id"    json:"user_id"   pk:"1"`
//		 RoleID    *int64    `db:"role_id"   auth:"role"  json:"role_id"`
//		 Email     string    `db:"email"     auth:"usr"   json:"email"`
//	}
//
// fields, _ := TagFieldsPivoted(new(User), "auth", []string{"id","role","usr"})
// fields -> []string{"user_id","role_id","email"}
func TagValuesPivoted(dsrc IDataSource, targetTagKey string, pivotTagKey string, pivotTagValues []string) ([]string, *TagError) {

	targetTagValues := make([]string, len(pivotTagValues))

	//r := reflect.ValueOf(dsrc).Elem()
	r := reflect.Indirect(reflect.ValueOf(dsrc))
	rt := r.Type()
	for i := 0; i < r.NumField(); i++ {
		tag := rt.Field(i).Tag
		if ptv, ok := tag.Lookup(pivotTagKey); ok {
			if ttv, ok := tag.Lookup(targetTagKey); ok && ttv != "-" {
				for j, v := range pivotTagValues {
					if v == ptv {
						targetTagValues[j] = ttv
					}
				}
			}
		}
	}

	for _, f := range targetTagValues {
		if f == "" {
			return targetTagValues, &TagError{Message: msg.Message{Msg: "incorrect pivot tags in model"}}
		}
	}

	return targetTagValues, nil
}
